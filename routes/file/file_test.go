/** File package, testing.
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
package file

import (
	"bytes"
	"fmt"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/configuration"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/helper"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/component-configuration"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/dashboard"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/infrastructure-component"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/scenario"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/user"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/widget"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var router *gin.Engine

type ConfigRequest struct {
	Name            string         `json:"name,omitempty"`
	ScenarioID      uint           `json:"scenarioID,omitempty"`
	ICID            uint           `json:"icID,omitempty"`
	StartParameters postgres.Jsonb `json:"startParameters,omitempty"`
}

type ICRequest struct {
	UUID       string         `json:"uuid,omitempty"`
	Host       string         `json:"host,omitempty"`
	Type       string         `json:"type,omitempty"`
	Name       string         `json:"name,omitempty"`
	Category   string         `json:"category,omitempty"`
	State      string         `json:"state,omitempty"`
	Properties postgres.Jsonb `json:"properties,omitempty"`
}

type ScenarioRequest struct {
	Name            string         `json:"name,omitempty"`
	Running         bool           `json:"running,omitempty"`
	StartParameters postgres.Jsonb `json:"startParameters,omitempty"`
}

type DashboardRequest struct {
	Name       string `json:"name,omitempty"`
	Grid       int    `json:"grid,omitempty"`
	ScenarioID uint   `json:"scenarioID,omitempty"`
}

type WidgetRequest struct {
	Name             string         `json:"name,omitempty"`
	Type             string         `json:"type,omitempty"`
	Width            uint           `json:"width,omitempty"`
	Height           uint           `json:"height,omitempty"`
	MinWidth         uint           `json:"minWidth,omitempty"`
	MinHeight        uint           `json:"minHeight,omitempty"`
	X                int            `json:"x,omitempty"`
	Y                int            `json:"y,omitempty"`
	Z                int            `json:"z,omitempty"`
	DashboardID      uint           `json:"dashboardID,omitempty"`
	IsLocked         bool           `json:"isLocked,omitempty"`
	CustomProperties postgres.Jsonb `json:"customProperties,omitempty"`
}

func addScenarioAndICAndConfigAndDashboardAndWidget() (scenarioID uint, ICID uint, configID uint, dashboardID uint, widgetID uint) {

	// authenticate as admin
	token, _ := helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.AdminCredentials)

	// POST $newICA
	newICA := ICRequest{
		UUID:       helper.ICA.UUID,
		Host:       helper.ICA.Host,
		Type:       helper.ICA.Type,
		Name:       helper.ICA.Name,
		Category:   helper.ICA.Category,
		State:      helper.ICA.State,
		Properties: helper.ICA.Properties,
	}
	_, resp, _ := helper.TestEndpoint(router, token,
		"/api/ic", "POST", helper.KeyModels{"ic": newICA})

	// Read newIC's ID from the response
	newICID, _ := helper.GetResponseID(resp)

	// authenticate as normal user
	token, _ = helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserACredentials)

	// POST $newScenario
	newScenario := ScenarioRequest{
		Name:            helper.ScenarioA.Name,
		Running:         helper.ScenarioA.Running,
		StartParameters: helper.ScenarioA.StartParameters,
	}
	_, resp, _ = helper.TestEndpoint(router, token,
		"/api/scenarios", "POST", helper.KeyModels{"scenario": newScenario})

	// Read newScenario's ID from the response
	newScenarioID, _ := helper.GetResponseID(resp)

	// POST new component config
	newConfig := ConfigRequest{
		Name:            helper.ConfigA.Name,
		ScenarioID:      uint(newScenarioID),
		ICID:            uint(newICID),
		StartParameters: helper.ConfigA.StartParameters,
	}
	_, resp, _ = helper.TestEndpoint(router, token,
		"/api/configs", "POST", helper.KeyModels{"config": newConfig})

	// Read newConfig's ID from the response
	newConfigID, _ := helper.GetResponseID(resp)

	// POST new dashboard
	newDashboard := DashboardRequest{
		Name:       helper.DashboardA.Name,
		Grid:       helper.DashboardA.Grid,
		ScenarioID: uint(newScenarioID),
	}
	_, resp, _ = helper.TestEndpoint(router, token,
		"/api/dashboards", "POST", helper.KeyModels{"dashboard": newDashboard})

	// Read newDashboard's ID from the response
	newDashboardID, _ := helper.GetResponseID(resp)

	// POST new widget
	newWidget := WidgetRequest{
		Name:             helper.WidgetA.Name,
		Type:             helper.WidgetA.Type,
		Width:            helper.WidgetA.Width,
		Height:           helper.WidgetA.Height,
		MinWidth:         helper.WidgetA.MinWidth,
		MinHeight:        helper.WidgetA.MinHeight,
		X:                helper.WidgetA.X,
		Y:                helper.WidgetA.Y,
		Z:                helper.WidgetA.Z,
		IsLocked:         helper.WidgetA.IsLocked,
		CustomProperties: helper.WidgetA.CustomProperties,
		DashboardID:      uint(newDashboardID),
	}
	_, resp, _ = helper.TestEndpoint(router, token,
		"/api/widgets", "POST", helper.KeyModels{"widget": newWidget})

	// Read newWidget's ID from the response
	newWidgetID, _ := helper.GetResponseID(resp)

	// add the guest user to the new scenario
	_, resp, _ = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/scenarios/%v/user?username=User_C", newScenarioID), "PUT", nil)

	return uint(newScenarioID), uint(newICID), uint(newConfigID), uint(newDashboardID), uint(newWidgetID)
}

func TestMain(m *testing.M) {
	err := configuration.InitConfig()
	if err != nil {
		panic(m)
	}
	err = database.InitDB(configuration.GolbalConfig)
	if err != nil {
		panic(m)
	}
	defer database.DBpool.Close()

	router = gin.Default()
	api := router.Group("/api")

	user.RegisterAuthenticate(api.Group("/authenticate"))
	api.Use(user.Authentication(true))
	// component-configuration endpoints required here to first add a config to the DB
	// that can be associated with a new file
	component_configuration.RegisterComponentConfigurationEndpoints(api.Group("/configs"))
	// scenario endpoints required here to first add a scenario to the DB
	// that can be associated with a new component config
	scenario.RegisterScenarioEndpoints(api.Group("/scenarios"))
	// IC endpoints required here to first add a IC to the DB
	// that can be associated with a new component config
	infrastructure_component.RegisterICEndpoints(api.Group("/ic"))
	// dashboard endpoints required here to first add a dashboard to the DB
	// that can be associated with a new widget
	dashboard.RegisterDashboardEndpoints(api.Group("/dashboards"))
	// widget endpoints required here to first add a widget to the DB
	// that can be associated with a new file
	widget.RegisterWidgetEndpoints(api.Group("/widgets"))

	RegisterFileEndpoints(api.Group("/files"))

	os.Exit(m.Run())
}

func TestAddFile(t *testing.T) {
	database.DropTables()
	database.MigrateModels()
	assert.NoError(t, helper.DBAddAdminAndUserAndGuest())

	// prepare the content of the DB for testing
	// using the respective endpoints of the API
	_, _, configID, _, widgetID := addScenarioAndICAndConfigAndDashboardAndWidget()

	// authenticate as userB who has no access to the elements in the DB
	token, err := helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserBCredentials)
	assert.NoError(t, err)

	emptyBuf := &bytes.Buffer{}

	// try to POST to a component config to which UserB has no access
	// should return a 422 unprocessable entity error
	code, resp, err := helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/files?objectID=%v&objectType=config", configID), "POST", emptyBuf)
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

	// try to POST to a widget to which UserB has no access
	// should return a 422 unprocessable entity error
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/files?objectID=%v&objectType=widget", widgetID), "POST", emptyBuf)
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

	// authenticate as normal userA
	token, err = helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserACredentials)
	assert.NoError(t, err)

	// try to POST to an invalid object type
	// should return a bad request error
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/files?objectID=%v&objectType=wrongtype", widgetID), "POST", emptyBuf)
	assert.NoError(t, err)
	assert.Equalf(t, 400, code, "Response body: \n%v\n", resp)

	// try to POST without an object ID
	// should return a bad request error
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/files?objectType=config"), "POST", emptyBuf)
	assert.NoError(t, err)
	assert.Equalf(t, 400, code, "Response body: \n%v\n", resp)

	// try to POST an invalid file
	// should return a bad request
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/files?objectID=%v&objectType=config", configID), "POST", emptyBuf)
	assert.NoError(t, err)
	assert.Equalf(t, 400, code, "Response body: \n%v\n", resp)

	// create a testfile.txt in local folder
	c1 := []byte("This is my testfile\n")
	err = ioutil.WriteFile("testfile.txt", c1, 0644)
	assert.NoError(t, err)

	// test POST files
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	fileWriter, err := bodyWriter.CreateFormFile("file", "testuploadfile.txt")
	assert.NoError(t, err, "writing to buffer")

	// open file handle
	fh, err := os.Open("testfile.txt")
	assert.NoError(t, err, "opening file")
	defer fh.Close()

	// io copy
	_, err = io.Copy(fileWriter, fh)
	assert.NoError(t, err, "IO copy")

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()
	//req, err := http.NewRequest("POST", "/api/files?objectID=1&objectType=widget", bodyBuf)

	// Create the request
	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", fmt.Sprintf("/api/files?objectID=%v&objectType=config", configID), bodyBuf)
	assert.NoError(t, err, "create request")

	req.Header.Set("Content-Type", contentType)
	req.Header.Add("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	assert.Equalf(t, 200, w.Code, "Response body: \n%v\n", w.Body)

	newFileID, err := helper.GetResponseID(w.Body)
	assert.NoError(t, err)

	// Get the new file
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/files/%v", newFileID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)
	assert.Equalf(t, string(c1), resp.String(), "Response body: \n%v\n", resp)

	// authenticate as userB who has no access to the elements in the DB
	token, err = helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserBCredentials)
	assert.NoError(t, err)

	// try to get a file to which user has no access
	// should return unprocessable entity
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/files/%v", newFileID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

}

func TestUpdateFile(t *testing.T) {

	database.DropTables()
	database.MigrateModels()
	assert.NoError(t, helper.DBAddAdminAndUserAndGuest())

	// prepare the content of the DB for testing
	// using the respective endpoints of the API
	_, _, configID, _, _ := addScenarioAndICAndConfigAndDashboardAndWidget()

	// authenticate as normal user
	token, err := helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserACredentials)
	assert.NoError(t, err)

	// create a testfile.txt in local folder
	c1 := []byte("This is my testfile\n")
	err = ioutil.WriteFile("testfile.txt", c1, 0644)
	assert.NoError(t, err)

	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	fileWriter, err := bodyWriter.CreateFormFile("file", "testfile.txt")
	assert.NoError(t, err, "writing to buffer")

	// open file handle
	fh, err := os.Open("testfile.txt")
	assert.NoError(t, err, "opening file")
	defer fh.Close()

	// io copy
	_, err = io.Copy(fileWriter, fh)
	assert.NoError(t, err, "IO copy")

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()
	//req, err := http.NewRequest("POST", "/api/files?objectID=1&objectType=widget", bodyBuf)

	// Create the POST request
	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", fmt.Sprintf("/api/files?objectID=%v&objectType=config", configID), bodyBuf)
	assert.NoError(t, err, "create request")

	req.Header.Set("Content-Type", contentType)
	req.Header.Add("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	assert.Equalf(t, 200, w.Code, "Response body: \n%v\n", w.Body)

	newFileID, err := helper.GetResponseID(w.Body)
	assert.NoError(t, err)

	// authenticate as userB who has no access to the elements in the DB
	token, err = helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserBCredentials)
	assert.NoError(t, err)

	emptyBuf := &bytes.Buffer{}

	// try to PUT to a file to which UserB has no access
	// should return a 422 unprocessable entity error
	code, resp, err := helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/files/%v", newFileID), "PUT", emptyBuf)
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

	// authenticate as guest user C
	token, err = helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.GuestCredentials)
	assert.NoError(t, err)

	// try to PUT as guest
	// should return an unprocessable entity error
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/files/%v", newFileID), "PUT", emptyBuf)
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

	// Prepare update
	// authenticate as normal user
	token, err = helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserACredentials)
	assert.NoError(t, err)

	// try to PUT with empty body
	// should return bad request
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/files/%v", newFileID), "PUT", emptyBuf)
	assert.NoError(t, err)
	assert.Equalf(t, 400, code, "Response body: \n%v\n", resp)

	// create a testfile_updated.txt in local folder
	c2 := []byte("This is my updated testfile\n")
	err = ioutil.WriteFile("testfileupdated.txt", c2, 0644)
	assert.NoError(t, err)

	bodyBufUpdated := &bytes.Buffer{}
	bodyWriterUpdated := multipart.NewWriter(bodyBufUpdated)
	fileWriterUpdated, err := bodyWriterUpdated.CreateFormFile("file", "testfileupdated.txt")
	assert.NoError(t, err, "writing to buffer")

	// open file handle for updated file
	fh_updated, err := os.Open("testfileupdated.txt")
	assert.NoError(t, err, "opening file")
	defer fh_updated.Close()

	// io copy
	_, err = io.Copy(fileWriterUpdated, fh_updated)
	assert.NoError(t, err, "IO copy")

	contentType = bodyWriterUpdated.FormDataContentType()
	bodyWriterUpdated.Close()

	// Create the PUT request
	w_updated := httptest.NewRecorder()
	req, err = http.NewRequest("PUT", fmt.Sprintf("/api/files/%v", newFileID), bodyBufUpdated)
	assert.NoError(t, err, "create request")

	req.Header.Set("Content-Type", contentType)
	req.Header.Add("Authorization", "Bearer "+token)
	router.ServeHTTP(w_updated, req)
	assert.Equalf(t, 200, w_updated.Code, "Response body: \n%v\n", w_updated.Body)

	newFileIDUpdated, err := helper.GetResponseID(w_updated.Body)

	assert.Equal(t, newFileID, newFileIDUpdated)

	// Get the updated file
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/files/%v", newFileIDUpdated), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)
	assert.Equalf(t, string(c2), resp.String(), "Response body: \n%v\n", resp)
}

func TestDeleteFile(t *testing.T) {
	database.DropTables()
	database.MigrateModels()
	assert.NoError(t, helper.DBAddAdminAndUserAndGuest())

	// prepare the content of the DB for testing
	// using the respective endpoints of the API
	_, _, configID, _, widgetID := addScenarioAndICAndConfigAndDashboardAndWidget()

	// authenticate as normal user
	token, err := helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserACredentials)
	assert.NoError(t, err)

	// create a testfile.txt in local folder
	c1 := []byte("This is my testfile\n")
	err = ioutil.WriteFile("testfile.txt", c1, 0644)
	assert.NoError(t, err)

	// open file handle
	fh, err := os.Open("testfile.txt")
	assert.NoError(t, err, "opening file")
	defer fh.Close()

	// test POST files
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	fileWriter, err := bodyWriter.CreateFormFile("file", "testuploadfile.txt")
	assert.NoError(t, err, "writing to buffer")
	// io copy
	_, err = io.Copy(fileWriter, fh)
	assert.NoError(t, err, "IO copy")
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	// Create the request
	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", fmt.Sprintf("/api/files?objectID=%v&objectType=config", configID), bodyBuf)
	assert.NoError(t, err, "create request")
	req.Header.Set("Content-Type", contentType)
	req.Header.Add("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)
	assert.Equalf(t, 200, w.Code, "Response body: \n%v\n", w.Body)

	newFileID, err := helper.GetResponseID(w.Body)
	assert.NoError(t, err)

	// add a second file to a widget, this time to a widget
	bodyBuf2 := &bytes.Buffer{}
	bodyWriter2 := multipart.NewWriter(bodyBuf2)
	fileWriter2, err := bodyWriter2.CreateFormFile("file", "testuploadfile.txt")
	assert.NoError(t, err, "writing to buffer")
	// io copy
	_, err = io.Copy(fileWriter2, fh)
	assert.NoError(t, err, "IO copy")
	contentType2 := bodyWriter2.FormDataContentType()
	bodyWriter2.Close()

	// Create the request
	w2 := httptest.NewRecorder()
	req2, err := http.NewRequest("POST", fmt.Sprintf("/api/files?objectID=%v&objectType=widget", widgetID), bodyBuf2)
	assert.NoError(t, err, "create request")
	req2.Header.Set("Content-Type", contentType2)
	req2.Header.Add("Authorization", "Bearer "+token)
	router.ServeHTTP(w2, req2)
	assert.Equalf(t, 200, w2.Code, "Response body: \n%v\n", w2.Body)

	newFileID2, err := helper.GetResponseID(w2.Body)
	assert.NoError(t, err)

	// authenticate as userB who has no access to the elements in the DB
	token, err = helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserBCredentials)
	assert.NoError(t, err)

	// try to DELETE file of component config to which userB has no access
	// should return an unprocessable entity error
	code, resp, err := helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/files/%v", newFileID), "DELETE", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

	// try to DELETE file of widget to which userB has no access
	// should return an unprocessable entity error
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/files/%v", newFileID2), "DELETE", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

	// authenticate as normal user
	token, err = helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserACredentials)
	assert.NoError(t, err)

	// Count the number of all files returned for component config
	initialNumber, err := helper.LengthOfResponse(router, token,
		fmt.Sprintf("/api/files?objectID=%v&objectType=config", configID), "GET", nil)
	assert.NoError(t, err)

	// try to DELETE non-existing fileID
	// should return not found
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/files/5"), "DELETE", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 404, code, "Response body: \n%v\n", resp)

	// authenticate as guest user C
	token, err = helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.GuestCredentials)
	assert.NoError(t, err)

	// try to DELETE file of component config as guest
	// should return an unprocessable entity error
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/files/%v", newFileID), "DELETE", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

	// try to DELETE file of widget as guest
	// should return an unprocessable entity error
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/files/%v", newFileID2), "DELETE", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

	// authenticate as normal user
	token, err = helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserACredentials)
	assert.NoError(t, err)

	// Delete the added file 1
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/files/%v", newFileID), "DELETE", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Delete the added file 2
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/files/%v", newFileID2), "DELETE", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Again count the number of all the files returned for component config
	finalNumber, err := helper.LengthOfResponse(router, token,
		fmt.Sprintf("/api/files?objectID=%v&objectType=config", configID), "GET", nil)
	assert.NoError(t, err)

	assert.Equal(t, initialNumber-1, finalNumber)
}

func TestGetAllFilesOfConfig(t *testing.T) {

	database.DropTables()
	database.MigrateModels()
	assert.NoError(t, helper.DBAddAdminAndUserAndGuest())

	// prepare the content of the DB for testing
	// using the respective endpoints of the API
	_, _, ConfigID, _, widgetID := addScenarioAndICAndConfigAndDashboardAndWidget()

	// authenticate as userB who has no access to the elements in the DB
	token, err := helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserBCredentials)
	assert.NoError(t, err)

	// try to get all files for component config to which userB has not access
	// should return unprocessable entity error
	code, resp, err := helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/files?objectID=%v&objectType=config", ConfigID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

	// try to get all files for widget to which userB has not access
	// should return unprocessable entity error
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/files?objectID=%v&objectType=widget", widgetID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

	// authenticate as normal userA
	token, err = helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserACredentials)
	assert.NoError(t, err)

	//try to get all files for unsupported object type; should return a bad request error
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/files?objectID=%v&objectType=wrongtype", ConfigID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 400, code, "Response body: \n%v\n", resp)

	//try to get all files with missing object ID; should return a bad request error
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/files?objectType=config"), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 400, code, "Response body: \n%v\n", resp)

	// Count the number of all files returned for component config
	initialNumberConfig, err := helper.LengthOfResponse(router, token,
		fmt.Sprintf("/api/files?objectID=%v&objectType=config", ConfigID), "GET", nil)
	assert.NoError(t, err)

	// Count the number of all files returned for widget
	initialNumberWidget, err := helper.LengthOfResponse(router, token,
		fmt.Sprintf("/api/files?objectID=%v&objectType=widget", widgetID), "GET", nil)
	assert.NoError(t, err)

	// create a testfile.txt in local folder
	c1 := []byte("This is my testfile\n")
	err = ioutil.WriteFile("testfile.txt", c1, 0644)
	assert.NoError(t, err)

	// open file handle
	fh, err := os.Open("testfile.txt")
	assert.NoError(t, err, "opening file")
	defer fh.Close()

	// test POST a file to component config and widget
	bodyBufConfig1 := &bytes.Buffer{}
	bodyBufWidget1 := &bytes.Buffer{}
	bodyWriterConfig1 := multipart.NewWriter(bodyBufConfig1)
	bodyWriterWidget1 := multipart.NewWriter(bodyBufWidget1)
	fileWriterConfig1, err := bodyWriterConfig1.CreateFormFile("file", "testuploadfile.txt")
	assert.NoError(t, err, "writing to buffer")
	fileWriterWidget1, err := bodyWriterWidget1.CreateFormFile("file", "testuploadfile.txt")
	assert.NoError(t, err, "writing to buffer")
	// io copy
	_, err = io.Copy(fileWriterConfig1, fh)
	assert.NoError(t, err, "IO copy")
	_, err = io.Copy(fileWriterWidget1, fh)
	assert.NoError(t, err, "IO copy")
	contentTypeConfig1 := bodyWriterConfig1.FormDataContentType()
	contentTypeWidget1 := bodyWriterWidget1.FormDataContentType()
	bodyWriterConfig1.Close()
	bodyWriterWidget1.Close()

	// Create the request for component config
	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", fmt.Sprintf("/api/files?objectID=%v&objectType=config", ConfigID), bodyBufConfig1)
	assert.NoError(t, err, "create request")
	req.Header.Set("Content-Type", contentTypeConfig1)
	req.Header.Add("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)
	assert.Equalf(t, 200, w.Code, "Response body: \n%v\n", w.Body)

	// Create the request for widget
	w2 := httptest.NewRecorder()
	req2, err := http.NewRequest("POST", fmt.Sprintf("/api/files?objectID=%v&objectType=widget", widgetID), bodyBufWidget1)
	assert.NoError(t, err, "create request")
	req2.Header.Set("Content-Type", contentTypeWidget1)
	req2.Header.Add("Authorization", "Bearer "+token)
	router.ServeHTTP(w2, req2)
	assert.Equalf(t, 200, w2.Code, "Response body: \n%v\n", w2.Body)

	// POST a second file to component config and widget

	// open a second file handle
	fh2, err := os.Open("testfile.txt")
	assert.NoError(t, err, "opening file")
	defer fh2.Close()

	bodyBufConfig2 := &bytes.Buffer{}
	bodyBufWidget2 := &bytes.Buffer{}
	bodyWriterConfig2 := multipart.NewWriter(bodyBufConfig2)
	bodyWriterWidget2 := multipart.NewWriter(bodyBufWidget2)
	fileWriterConfig2, err := bodyWriterConfig2.CreateFormFile("file", "testuploadfile2.txt")
	assert.NoError(t, err, "writing to buffer")
	fileWriterWidget2, err := bodyWriterWidget2.CreateFormFile("file", "testuploadfile2.txt")
	assert.NoError(t, err, "writing to buffer")

	// io copy
	_, err = io.Copy(fileWriterConfig2, fh2)
	assert.NoError(t, err, "IO copy")
	_, err = io.Copy(fileWriterWidget2, fh2)
	assert.NoError(t, err, "IO copy")
	contentTypeConfig2 := bodyWriterConfig2.FormDataContentType()
	contentTypeWidget2 := bodyWriterWidget2.FormDataContentType()
	bodyWriterConfig2.Close()
	bodyWriterWidget2.Close()

	w3 := httptest.NewRecorder()
	req3, err := http.NewRequest("POST", fmt.Sprintf("/api/files?objectID=%v&objectType=config", ConfigID), bodyBufConfig2)
	assert.NoError(t, err, "create request")
	req3.Header.Set("Content-Type", contentTypeConfig2)
	req3.Header.Add("Authorization", "Bearer "+token)
	router.ServeHTTP(w3, req3)
	assert.Equalf(t, 200, w3.Code, "Response body: \n%v\n", w3.Body)

	w4 := httptest.NewRecorder()
	req4, err := http.NewRequest("POST", fmt.Sprintf("/api/files?objectID=%v&objectType=widget", widgetID), bodyBufWidget2)
	assert.NoError(t, err, "create request")
	req4.Header.Set("Content-Type", contentTypeWidget2)
	req4.Header.Add("Authorization", "Bearer "+token)
	router.ServeHTTP(w4, req4)
	assert.Equalf(t, 200, w4.Code, "Response body: \n%v\n", w4.Body)

	// Again count the number of all the files returned for component config
	finalNumberConfig, err := helper.LengthOfResponse(router, token,
		fmt.Sprintf("/api/files?objectID=%v&objectType=config", ConfigID), "GET", nil)
	assert.NoError(t, err)
	assert.Equal(t, initialNumberConfig+2, finalNumberConfig)

	// Again count the number of all the files returned for widget
	finalNumberWidget, err := helper.LengthOfResponse(router, token,
		fmt.Sprintf("/api/files?objectID=%v&objectType=widget", widgetID), "GET", nil)
	assert.NoError(t, err)
	assert.Equal(t, initialNumberWidget+2, finalNumberWidget)
}
