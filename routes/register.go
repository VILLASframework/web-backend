/**
* This file is part of VILLASweb-backend-go
*
* This program is free software: you can redistribute it and/or modify
* it under the terms of the GNU General Public License as published by
* the Free Software Foundation, either version 3 of the License, or
* any later version.
*
* This program is distributed in the hope that it will be useful,
* but WITHOUT ANY WARRANTY; without even the implied warranty of
* MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
* GNU General Public License for more details.
*
* You should have received a copy of the GNU General Public License
* along with this program.  If not, see <http://www.gnu.org/licenses/>.
*********************************************************************************/

package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"

	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/helper"
	component_configuration "git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/component-configuration"
	config_route "git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/config"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/dashboard"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/file"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/healthz"
	infrastructure_component "git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/infrastructure-component"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/metrics"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/openapi"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/result"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/scenario"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/signal"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/user"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/usergroup"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/widget"
	"github.com/gin-gonic/gin"
	"github.com/zpatrick/go-config"
)

// #######################################################################
// #################### Structures for test data #########################
// #######################################################################

type jsonUser struct {
	Username string
	Password string
	Mail     string
	Role     string
}

var GlobalTestData struct {
	Users      []jsonUser
	ICs        []database.InfrastructureComponent
	Scenarios  []database.Scenario
	Results    []database.Result
	Configs    []database.ComponentConfiguration
	Dashboards []database.Dashboard
	Widgets    []database.Widget
	Signals    []database.Signal
}

// register all backend endpoints; to be called after DB is initialized
func RegisterEndpoints(router *gin.Engine, api *gin.RouterGroup) {

	healthz.RegisterHealthzEndpoint(api.Group("/healthz"))
	metrics.RegisterMetricsEndpoint(api.Group("/metrics"))
	openapi.RegisterOpenAPIEndpoint(api.Group("/openapi"))
	config_route.RegisterConfigEndpoint(api.Group("/config"))
	user.RegisterAuthenticate(api.Group("/authenticate"))

	// The following endpoints require authentication

	api.Use(user.Authentication())

	scenario.RegisterScenarioEndpoints(api.Group("/scenarios"))
	usergroup.RegisterUserGroupEndpoints(api.Group("/usergroups"))
	component_configuration.RegisterComponentConfigurationEndpoints(api.Group("/configs"))
	signal.RegisterSignalEndpoints(api.Group("/signals"))
	dashboard.RegisterDashboardEndpoints(api.Group("/dashboards"))
	widget.RegisterWidgetEndpoints(api.Group("/widgets"))
	file.RegisterFileEndpoints(api.Group("/files"))
	user.RegisterUserEndpoints(api.Group("/users"))
	infrastructure_component.RegisterICEndpoints(api.Group("/ic"))
	result.RegisterResultEndpoints(api.Group("/results"))

	metrics.InitCounters()

}

// ReadTestDataFromJson Reads test data from JSON file (path set by ENV variable or command line param)
func ReadTestDataFromJson(path string) error {

	_, err := os.Stat(path)

	if err == nil {

		jsonFile, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("error opening json file: %v", err)
		}
		log.Println("Successfully opened json data file", path)

		defer jsonFile.Close()

		byteValue, _ := io.ReadAll(jsonFile)

		err = json.Unmarshal(byteValue, &GlobalTestData)
		if err != nil {
			return fmt.Errorf("error unmarshalling json: %v", err)
		}
	} else if os.IsNotExist(err) {
		log.Println("Test data file does not exist, no test data added to DB:", path)
		return nil
	} else {
		log.Println("Something is wrong with this file path:", path)
		return nil
	}

	return nil
}

// AddTestData Uses API endpoints to add test data to the backend; All endpoints have to be registered before invoking this function.
func AddTestData(cfg *config.Config, router *gin.Engine) (*bytes.Buffer, error) {

	adminPW, errPW := cfg.String("admin.pass")
	adminName, errName := cfg.String("admin.user")
	if errPW != nil || errName != nil {
		if errName != nil {
			log.Println("WARNING:", errName)
		}
		if errPW != nil {
			log.Println("WARNING:", errPW)
		}

		log.Println("WARNING: cannot add test data because of missing admin config, continue without it")
		return nil, nil
	}

	var Admin = database.Credentials{
		Username: adminName,
		Password: adminPW,
	}

	// authenticate as admin
	token, err := helper.AuthenticateForTest(router, Admin)
	if err != nil {
		return nil, err
	}

	basePath := "/api/v2"

	db := database.GetDB()

	// add users
	for _, u := range GlobalTestData.Users {

		var x []user.User
		err = db.Find(&x, "Username = ?", u.Username).Error
		if err != nil {
			return nil, err
		}

		if len(x) == 0 {
			code, resp, err := helper.TestEndpoint(router, token, basePath+"/users", "POST", helper.KeyModels{"user": u})
			if code != http.StatusOK {
				return resp, fmt.Errorf("error adding user %v: %v", u.Username, err)
			}
		}
	}

	// add infrastructure components
	amqphost, _ := cfg.String("amqp.host")
	counterICs := 0
	for _, i := range GlobalTestData.ICs {

		if (i.ManagedExternally && amqphost != "") || !i.ManagedExternally {

			var x []infrastructure_component.InfrastructureComponent
			err = db.Find(&x, "Name = ?", i.Name).Error
			if err != nil {
				return nil, err
			}

			if len(x) == 0 {
				code, resp, err := helper.TestEndpoint(router, token, basePath+"/ic", "POST", helper.KeyModels{"ic": i})
				if code != http.StatusOK {
					return resp, fmt.Errorf("error adding IC %v: %v", i.Name, err)
				}
				counterICs++
			}
		}
	}

	// add scenarios
	for _, s := range GlobalTestData.Scenarios {

		var x []scenario.Scenario
		err = db.Find(&x, "Name = ?", s.Name).Error
		if err != nil {
			return nil, err
		}

		if len(x) == 0 {
			code, resp, err := helper.TestEndpoint(router, token, basePath+"/scenarios", "POST", helper.KeyModels{"scenario": s})
			if code != http.StatusOK {
				return resp, fmt.Errorf("error adding Scenario %v: %v", s.Name, err)
			}

			// add all users to the scenario
			for _, u := range GlobalTestData.Users {
				code, resp, err := helper.TestEndpoint(router, token, fmt.Sprintf("%v/scenarios/1/user?username="+u.Username, basePath), "PUT", nil)
				if code != http.StatusOK {
					return resp, fmt.Errorf("error adding user %v to scenario %v: %v", u.Username, s.Name, err)
				}
			}
		}
	}

	// If there is at least one scenario and one IC in the test data, add component configs
	configCounter := 0
	if len(GlobalTestData.Scenarios) > 0 && counterICs > 0 {

		for _, c := range GlobalTestData.Configs {

			var x []component_configuration.ComponentConfiguration
			err = db.Find(&x, "Name = ?", c.Name).Error
			if err != nil {
				return nil, err
			}

			if len(x) == 0 {
				c.ScenarioID = 1
				c.ICID = 1
				code, resp, err := helper.TestEndpoint(router, token, basePath+"/configs", "POST", helper.KeyModels{"config": c})
				if code != http.StatusOK {
					return resp, fmt.Errorf("error adding Config %v: %v", c.Name, err)
				}
			}
			configCounter++
		}
	}

	// If there is at least one scenario, add dashboards, results, and 2 test files
	dashboardCounter := 0
	if len(GlobalTestData.Scenarios) > 0 {
		for _, d := range GlobalTestData.Dashboards {

			var x []dashboard.Dashboard
			err = db.Find(&x, "Name = ?", d.Name).Error
			if err != nil {
				return nil, err
			}

			if len(x) == 0 {
				d.ScenarioID = 1
				code, resp, err := helper.TestEndpoint(router, token, basePath+"/dashboards", "POST", helper.KeyModels{"dashboard": d})
				if code != http.StatusOK {
					return resp, fmt.Errorf("error adding Dashboard %v: %v", d.Name, err)
				}
			}
			dashboardCounter++
		}

		for _, r := range GlobalTestData.Results {

			var x []result.Result
			err = db.Find(&x, "Description = ?", r.Description).Error
			if err != nil {
				return nil, err
			}

			if len(x) == 0 {
				r.ScenarioID = 1
				r.ResultFileIDs = []int64{}
				code, resp, err := helper.TestEndpoint(router, token, basePath+"/results", "POST", helper.KeyModels{"result": r})
				if code != http.StatusOK {
					return resp, fmt.Errorf("error adding Result %v: %v", r.Description, err)
				}
			}
		}

		// upload files

		var x []file.File
		err = db.Find(&x, "Name = ?", "Readme.md").Error
		if err != nil {
			return nil, err
		}

		if len(x) == 0 {
			// upload readme file
			bodyBuf := &bytes.Buffer{}
			bodyWriter := multipart.NewWriter(bodyBuf)
			fileWriter, _ := bodyWriter.CreateFormFile("file", "Readme.md")
			fh, _ := os.Open("README.md")
			defer fh.Close()

			// io copy
			_, _ = io.Copy(fileWriter, fh)
			contentType := bodyWriter.FormDataContentType()
			bodyWriter.Close()

			// Create the request and add file to scenario
			w1 := httptest.NewRecorder()
			req1, _ := http.NewRequest("POST", basePath+"/files?scenarioID=1", bodyBuf)
			req1.Header.Set("Content-Type", contentType)
			req1.Header.Add("Authorization", "Bearer "+token)
			router.ServeHTTP(w1, req1)
		}

		var y []file.File
		err = db.Find(&y, "Name = ?", "logo.png").Error
		if err != nil {
			return nil, err
		}

		if len(y) == 0 {
			// upload image file
			bodyBuf := &bytes.Buffer{}
			bodyWriter := multipart.NewWriter(bodyBuf)
			fileWriter, _ := bodyWriter.CreateFormFile("file", "logo.png")
			fh, _ := os.Open("doc/pictures/villas_web.png")
			defer fh.Close()

			// io copy
			_, _ = io.Copy(fileWriter, fh)
			contentType := bodyWriter.FormDataContentType()
			bodyWriter.Close()

			// Create the request and add a second file to scenario
			w2 := httptest.NewRecorder()
			req2, _ := http.NewRequest("POST", basePath+"/files?scenarioID=1", bodyBuf)
			req2.Header.Set("Content-Type", contentType)
			req2.Header.Add("Authorization", "Bearer "+token)
			router.ServeHTTP(w2, req2)
		}
	}

	// If there is at least one dashboard, add widgets
	if dashboardCounter > 0 {
		for _, w := range GlobalTestData.Widgets {

			var x []widget.Widget
			err = db.Find(&x, "Name = ?", w.Name).Error
			if err != nil {
				return nil, err
			}

			if len(x) == 0 {
				w.DashboardID = 1
				code, resp, err := helper.TestEndpoint(router, token, basePath+"/widgets", "POST", helper.KeyModels{"widget": w})
				if code != http.StatusOK {
					return resp, fmt.Errorf("error adding Widget %v: %v", w.Name, err)
				}
			}
		}
	}

	// If there is at least one config, add signals
	if configCounter > 0 {
		for _, s := range GlobalTestData.Signals {

			var x []signal.Signal
			err = db.Find(&x, "Name = ?", s.Name).Error
			if err != nil {
				return nil, err
			}

			if len(x) == 0 {
				s.ConfigID = 1
				code, resp, err := helper.TestEndpoint(router, token, basePath+"/signals", "POST", helper.KeyModels{"signal": s})
				if code != http.StatusOK {
					return resp, fmt.Errorf("error adding Signal %v: %v", s.Name, err)
				}
			}
		}
	}

	return nil, nil
}
