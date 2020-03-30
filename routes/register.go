/** Routes package, registration function
*
* @author Sonja Happ <sonja.happ@eonerc.rwth-aachen.de>
* @copyright 2014-2019, Institute for Automation of Complex Power Systems, EONERC
* @license GNU General Public License (version 3)
*
* VILLASweb-backend-go
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
	"fmt"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/helper"
	component_configuration "git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/component-configuration"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/dashboard"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/file"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/healthz"
	infrastructure_component "git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/infrastructure-component"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/metrics"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/scenario"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/signal"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/user"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/widget"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
)

// register all backend endpoints; to be called after DB is initialized
func RegisterEndpoints(router *gin.Engine, api *gin.RouterGroup) {

	healthz.RegisterHealthzEndpoint(api.Group("/healthz"))
	metrics.RegisterMetricsEndpoint(api.Group("/metrics"))
	// All endpoints (except for /healthz and /metrics) require authentication except when someone wants to
	// login (POST /authenticate)
	user.RegisterAuthenticate(api.Group("/authenticate"))

	api.Use(user.Authentication(true))

	scenario.RegisterScenarioEndpoints(api.Group("/scenarios"))
	component_configuration.RegisterComponentConfigurationEndpoints(api.Group("/configs"))
	signal.RegisterSignalEndpoints(api.Group("/signals"))
	dashboard.RegisterDashboardEndpoints(api.Group("/dashboards"))
	widget.RegisterWidgetEndpoints(api.Group("/widgets"))
	file.RegisterFileEndpoints(api.Group("/files"))
	user.RegisterUserEndpoints(api.Group("/users"))
	infrastructure_component.RegisterICEndpoints(api.Group("/ic"))

	router.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	metrics.InitCounters()

}

// Uses API endpoints to add test data to the backend; All endpoints have to be registered before invoking this function.
func AddTestData(basePath string, router *gin.Engine) (*bytes.Buffer, error) {

	database.MigrateModels()
	// Create entries of each model (data defined in test_data.go)
	// add Admin user
	err := helper.DBAddAdminUser()
	if err != nil {
		return nil, err
	}

	// authenticate as admin
	token, err := helper.AuthenticateForTest(router, basePath+"/authenticate", "POST", helper.AdminCredentials)
	if err != nil {
		return nil, err
	}

	// add 2 normal and 1 guest user
	code, resp, err := helper.TestEndpoint(router, token, basePath+"/users", "POST", helper.KeyModels{"user": helper.NewUserA})
	if code != http.StatusOK {
		return resp, fmt.Errorf("error adding User_A")
	}
	code, resp, err = helper.TestEndpoint(router, token, basePath+"/users", "POST", helper.KeyModels{"user": helper.NewUserB})
	if code != http.StatusOK {
		return resp, fmt.Errorf("error adding User_B")
	}
	code, resp, err = helper.TestEndpoint(router, token, basePath+"/users", "POST", helper.KeyModels{"user": helper.NewUserC})
	if code != http.StatusOK {
		return resp, fmt.Errorf("error adding User_C")
	}

	// add infrastructure components
	code, resp, err = helper.TestEndpoint(router, token, basePath+"/ic", "POST", helper.KeyModels{"ic": helper.ICA})
	if code != http.StatusOK {
		return resp, fmt.Errorf("error adding IC A")
	}
	code, resp, err = helper.TestEndpoint(router, token, basePath+"/ic", "POST", helper.KeyModels{"ic": helper.ICB})
	if code != http.StatusOK {
		return resp, fmt.Errorf("error adding IC B")
	}

	// add scenarios
	code, resp, err = helper.TestEndpoint(router, token, basePath+"/scenarios", "POST", helper.KeyModels{"scenario": helper.ScenarioA})
	if code != http.StatusOK {
		return resp, fmt.Errorf("error adding Scenario A")
	}
	code, resp, err = helper.TestEndpoint(router, token, basePath+"/scenarios", "POST", helper.KeyModels{"scenario": helper.ScenarioB})
	if code != http.StatusOK {
		return resp, fmt.Errorf("error adding Scenario B")
	}

	// add users to scenario
	code, resp, err = helper.TestEndpoint(router, token, fmt.Sprintf("%v/scenarios/1/user?username=User_A", basePath), "PUT", nil)
	if code != http.StatusOK {
		return resp, fmt.Errorf("error adding User_A to Scenario A")
	}
	code, resp, err = helper.TestEndpoint(router, token, fmt.Sprintf("%v/scenarios/2/user?username=User_A", basePath), "PUT", nil)
	if code != http.StatusOK {
		return resp, fmt.Errorf("error adding User_A to Scenario B")
	}
	code, resp, err = helper.TestEndpoint(router, token, fmt.Sprintf("%v/scenarios/2/user?username=User_B", basePath), "PUT", nil)
	if code != http.StatusOK {
		return resp, fmt.Errorf("error adding User_B to Scenario B")
	}
	code, resp, err = helper.TestEndpoint(router, token, fmt.Sprintf("%v/scenarios/1/user?username=User_C", basePath), "PUT", nil)
	if code != http.StatusOK {
		return resp, fmt.Errorf("error adding User_C to Scenario A")
	}

	// add component configurations
	configA := helper.ConfigA
	configB := helper.ConfigB
	configA.ScenarioID = 1
	configB.ScenarioID = 1
	configA.ICID = 2
	configB.ICID = 1
	code, resp, err = helper.TestEndpoint(router, token, basePath+"/configs", "POST", helper.KeyModels{"config": configA})
	if code != http.StatusOK {
		return resp, fmt.Errorf("error adding Config A")
	}
	code, resp, err = helper.TestEndpoint(router, token, basePath+"/configs", "POST", helper.KeyModels{"config": configB})
	if code != http.StatusOK {
		return resp, fmt.Errorf("error adding Config B")
	}

	// add dashboards
	dashboardA := helper.DashboardA
	dashboardB := helper.DashboardB
	dashboardA.ScenarioID = 1
	dashboardB.ScenarioID = 1
	code, resp, err = helper.TestEndpoint(router, token, basePath+"/dashboards", "POST", helper.KeyModels{"dashboard": dashboardA})
	if code != http.StatusOK {
		return resp, fmt.Errorf("error adding Dashboard B")
	}
	code, resp, err = helper.TestEndpoint(router, token, basePath+"/dashboards", "POST", helper.KeyModels{"dashboard": dashboardB})
	if code != http.StatusOK {
		return resp, fmt.Errorf("error adding Dashboard B")
	}

	// add widgets
	widgetA := helper.WidgetA
	widgetB := helper.WidgetB
	widgetC := helper.WidgetC
	widgetD := helper.WidgetD
	widgetE := helper.WidgetE
	widgetA.DashboardID = 1
	widgetB.DashboardID = 1
	widgetC.DashboardID = 1
	widgetD.DashboardID = 1
	widgetE.DashboardID = 1
	code, resp, err = helper.TestEndpoint(router, token, basePath+"/widgets", "POST", helper.KeyModels{"widget": widgetA})
	if code != http.StatusOK {
		return resp, fmt.Errorf("error adding Widget A")
	}
	code, resp, err = helper.TestEndpoint(router, token, basePath+"/widgets", "POST", helper.KeyModels{"widget": widgetB})
	if code != http.StatusOK {
		return resp, fmt.Errorf("error adding Widget B")
	}
	code, resp, err = helper.TestEndpoint(router, token, basePath+"/widgets", "POST", helper.KeyModels{"widget": widgetC})
	if code != http.StatusOK {
		return resp, fmt.Errorf("error adding Widget C")
	}
	code, resp, err = helper.TestEndpoint(router, token, basePath+"/widgets", "POST", helper.KeyModels{"widget": widgetD})
	if code != http.StatusOK {
		return resp, fmt.Errorf("error adding Widget D")
	}
	code, resp, err = helper.TestEndpoint(router, token, basePath+"/widgets", "POST", helper.KeyModels{"widget": widgetE})
	if code != http.StatusOK {
		return resp, fmt.Errorf("error adding Widget E")
	}

	// add signals
	outSignalA := helper.OutSignalA
	outSignalB := helper.OutSignalB
	inSignalA := helper.InSignalA
	inSignalB := helper.InSignalB
	outSignalC := helper.OutSignalC
	outSignalD := helper.OutSignalD
	outSignalE := helper.OutSignalE
	outSignalA.ConfigID = 1
	outSignalB.ConfigID = 1
	outSignalC.ConfigID = 1
	outSignalD.ConfigID = 1
	outSignalE.ConfigID = 1
	inSignalA.ConfigID = 1
	inSignalB.ConfigID = 2

	code, resp, err = helper.TestEndpoint(router, token, basePath+"/signals", "POST", helper.KeyModels{"signal": outSignalB})
	if code != http.StatusOK {
		return resp, fmt.Errorf("error adding outSignalB")
	}
	code, resp, err = helper.TestEndpoint(router, token, basePath+"/signals", "POST", helper.KeyModels{"signal": outSignalA})
	if code != http.StatusOK {
		return resp, fmt.Errorf("error adding outSignalA")
	}
	code, resp, err = helper.TestEndpoint(router, token, basePath+"/signals", "POST", helper.KeyModels{"signal": outSignalC})
	if code != http.StatusOK {
		return resp, fmt.Errorf("error adding outSignalC")
	}
	code, resp, err = helper.TestEndpoint(router, token, basePath+"/signals", "POST", helper.KeyModels{"signal": outSignalD})
	if code != http.StatusOK {
		return resp, fmt.Errorf("error adding outSignalD")
	}
	code, resp, err = helper.TestEndpoint(router, token, basePath+"/signals", "POST", helper.KeyModels{"signal": outSignalE})
	if code != http.StatusOK {
		return resp, fmt.Errorf("error adding outSignalE")
	}

	code, resp, err = helper.TestEndpoint(router, token, basePath+"/signals", "POST", helper.KeyModels{"signal": inSignalA})
	if code != http.StatusOK {
		return resp, fmt.Errorf("error adding inSignalA")
	}
	code, resp, err = helper.TestEndpoint(router, token, basePath+"/signals", "POST", helper.KeyModels{"signal": inSignalB})
	if code != http.StatusOK {
		return resp, fmt.Errorf("error adding inSignalB")
	}

	// upload files

	// upload readme file
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	fileWriter, _ := bodyWriter.CreateFormFile("file", "Readme.md")
	fh, _ := os.Open("README.md")
	defer fh.Close()

	// io copy
	_, err = io.Copy(fileWriter, fh)
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	// Create the request and add file to component config
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("POST", basePath+"/files?objectID=1&objectType=config", bodyBuf)
	req1.Header.Set("Content-Type", contentType)
	req1.Header.Add("Authorization", "Bearer "+token)
	router.ServeHTTP(w1, req1)

	// upload image file
	bodyBuf = &bytes.Buffer{}
	bodyWriter = multipart.NewWriter(bodyBuf)
	fileWriter, _ = bodyWriter.CreateFormFile("file", "logo.png")
	fh, _ = os.Open("doc/pictures/villas_web.png")
	defer fh.Close()

	// io copy
	_, err = io.Copy(fileWriter, fh)
	contentType = bodyWriter.FormDataContentType()
	bodyWriter.Close()

	// Create the request and add file to widget
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", basePath+"/files?objectID=1&objectType=widget", bodyBuf)
	req2.Header.Set("Content-Type", contentType)
	req2.Header.Add("Authorization", "Bearer "+token)
	router.ServeHTTP(w2, req2)

	return nil, nil
}
