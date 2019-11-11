package file

import (
	"bytes"
	"fmt"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/helper"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/dashboard"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/scenario"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/simulationmodel"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/simulator"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/user"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/widget"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
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
var db *gorm.DB

type SimulationModelRequest struct {
	Name            string         `json:"name,omitempty"`
	ScenarioID      uint           `json:"scenarioID,omitempty"`
	SimulatorID     uint           `json:"simulatorID,omitempty"`
	StartParameters postgres.Jsonb `json:"startParameters,omitempty"`
}

type SimulatorRequest struct {
	UUID       string         `json:"uuid,omitempty"`
	Host       string         `json:"host,omitempty"`
	Modeltype  string         `json:"modelType,omitempty"`
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

func addScenarioAndSimulatorAndSimulationModelAndDashboardAndWidget() (scenarioID uint, simulatorID uint, simulationModelID uint, dashboardID uint, widgetID uint) {

	// authenticate as admin
	token, _ := helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.AdminCredentials)

	// POST $newSimulatorA
	newSimulatorA := SimulatorRequest{
		UUID:       database.SimulatorA.UUID,
		Host:       database.SimulatorA.Host,
		Modeltype:  database.SimulatorA.Modeltype,
		State:      database.SimulatorA.State,
		Properties: database.SimulatorA.Properties,
	}
	_, resp, _ := helper.TestEndpoint(router, token,
		"/api/simulators", "POST", helper.KeyModels{"simulator": newSimulatorA})

	// Read newSimulator's ID from the response
	newSimulatorID, _ := helper.GetResponseID(resp)

	// authenticate as normal user
	token, _ = helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserACredentials)

	// POST $newScenario
	newScenario := ScenarioRequest{
		Name:            database.ScenarioA.Name,
		Running:         database.ScenarioA.Running,
		StartParameters: database.ScenarioA.StartParameters,
	}
	_, resp, _ = helper.TestEndpoint(router, token,
		"/api/scenarios", "POST", helper.KeyModels{"scenario": newScenario})

	// Read newScenario's ID from the response
	newScenarioID, _ := helper.GetResponseID(resp)

	// POST new simulation model
	newSimulationModel := SimulationModelRequest{
		Name:            database.SimulationModelA.Name,
		ScenarioID:      uint(newScenarioID),
		SimulatorID:     uint(newSimulatorID),
		StartParameters: database.SimulationModelA.StartParameters,
	}
	_, resp, _ = helper.TestEndpoint(router, token,
		"/api/models", "POST", helper.KeyModels{"model": newSimulationModel})

	// Read newSimulationModel's ID from the response
	newSimulationModelID, _ := helper.GetResponseID(resp)

	// POST new dashboard
	newDashboard := DashboardRequest{
		Name:       database.DashboardA.Name,
		Grid:       database.DashboardA.Grid,
		ScenarioID: uint(newScenarioID),
	}
	_, resp, _ = helper.TestEndpoint(router, token,
		"/api/dashboards", "POST", helper.KeyModels{"dashboard": newDashboard})

	// Read newDashboard's ID from the response
	newDashboardID, _ := helper.GetResponseID(resp)

	// POST new widget
	newWidget := WidgetRequest{
		Name:             database.WidgetA.Name,
		Type:             database.WidgetA.Type,
		Width:            database.WidgetA.Width,
		Height:           database.WidgetA.Height,
		MinWidth:         database.WidgetA.MinWidth,
		MinHeight:        database.WidgetA.MinHeight,
		X:                database.WidgetA.X,
		Y:                database.WidgetA.Y,
		Z:                database.WidgetA.Z,
		IsLocked:         database.WidgetA.IsLocked,
		CustomProperties: database.WidgetA.CustomProperties,
		DashboardID:      uint(newDashboardID),
	}
	_, resp, _ = helper.TestEndpoint(router, token,
		"/api/widgets", "POST", helper.KeyModels{"widget": newWidget})

	// Read newWidget's ID from the response
	newWidgetID, _ := helper.GetResponseID(resp)

	// add the guest user to the new scenario
	_, resp, _ = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/scenarios/%v/user?username=User_C", newScenarioID), "PUT", nil)

	return uint(newScenarioID), uint(newSimulatorID), uint(newSimulationModelID), uint(newDashboardID), uint(newWidgetID)
}

func TestMain(m *testing.M) {

	db = database.InitDB(database.DB_NAME, true)
	defer db.Close()

	router = gin.Default()
	api := router.Group("/api")

	user.RegisterAuthenticate(api.Group("/authenticate"))
	api.Use(user.Authentication(true))
	// simulationmodel endpoints required here to first add a simulation to the DB
	// that can be associated with a new file
	simulationmodel.RegisterSimulationModelEndpoints(api.Group("/models"))
	// scenario endpoints required here to first add a scenario to the DB
	// that can be associated with a new simulation model
	scenario.RegisterScenarioEndpoints(api.Group("/scenarios"))
	// simulator endpoints required here to first add a simulator to the DB
	// that can be associated with a new simulation model
	simulator.RegisterSimulatorEndpoints(api.Group("/simulators"))
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
	database.DropTables(db)
	database.MigrateModels(db)
	assert.NoError(t, database.DBAddAdminAndUserAndGuest(db))

	// prepare the content of the DB for testing
	// using the respective endpoints of the API
	_, _, simulationModelID, _, widgetID := addScenarioAndSimulatorAndSimulationModelAndDashboardAndWidget()

	// authenticate as userB who has no access to the elements in the DB
	token, err := helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserBCredentials)
	assert.NoError(t, err)

	emptyBuf := &bytes.Buffer{}

	// try to POST to a simulation model to which UserB has no access
	// should return a 422 unprocessable entity error
	code, resp, err := helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/files?objectID=%v&objectType=model", simulationModelID), "POST", emptyBuf)
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
		fmt.Sprintf("/api/files?objectType=model"), "POST", emptyBuf)
	assert.NoError(t, err)
	assert.Equalf(t, 400, code, "Response body: \n%v\n", resp)

	// try to POST an invalid file
	// should return a bad request
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/files?objectID=%v&objectType=model", simulationModelID), "POST", emptyBuf)
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
	req, err := http.NewRequest("POST", fmt.Sprintf("/api/files?objectID=%v&objectType=model", simulationModelID), bodyBuf)
	assert.NoError(t, err, "create request")

	req.Header.Set("Content-Type", contentType)
	req.Header.Add("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	assert.Equalf(t, 200, w.Code, "Response body: \n%v\n", w.Body)
	//fmt.Println(w.Body)

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

	database.DropTables(db)
	database.MigrateModels(db)
	assert.NoError(t, database.DBAddAdminAndUserAndGuest(db))

	// prepare the content of the DB for testing
	// using the respective endpoints of the API
	_, _, simulationModelID, _, _ := addScenarioAndSimulatorAndSimulationModelAndDashboardAndWidget()

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
	req, err := http.NewRequest("POST", fmt.Sprintf("/api/files?objectID=%v&objectType=model", simulationModelID), bodyBuf)
	assert.NoError(t, err, "create request")

	req.Header.Set("Content-Type", contentType)
	req.Header.Add("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	assert.Equalf(t, 200, w.Code, "Response body: \n%v\n", w.Body)
	//fmt.Println(w.Body)

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
	//fmt.Println(w_updated.Body)
	//assert.Equal(t, c2, w_updated.Header().Get("file"))

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
	database.DropTables(db)
	database.MigrateModels(db)
	assert.NoError(t, database.DBAddAdminAndUserAndGuest(db))

	// prepare the content of the DB for testing
	// using the respective endpoints of the API
	_, _, simulationModelID, _, widgetID := addScenarioAndSimulatorAndSimulationModelAndDashboardAndWidget()

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
	//req, err := http.NewRequest("POST", "/api/files?objectID=1&objectType=widget", bodyBuf)

	// Create the request
	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", fmt.Sprintf("/api/files?objectID=%v&objectType=model", simulationModelID), bodyBuf)
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

	// try to DELETE file of simulation model to which userB has no access
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

	// Count the number of all files returned for simulation model
	initialNumber, err := helper.LengthOfResponse(router, token,
		fmt.Sprintf("/api/files?objectID=%v&objectType=model", simulationModelID), "GET", nil)
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

	// try to DELETE file of simulation model as guest
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

	// Again count the number of all the files returned for simulation model
	finalNumber, err := helper.LengthOfResponse(router, token,
		fmt.Sprintf("/api/files?objectID=%v&objectType=model", simulationModelID), "GET", nil)
	assert.NoError(t, err)

	assert.Equal(t, initialNumber-1, finalNumber)
}

func TestGetAllFilesOfSimulationModel(t *testing.T) {

	database.DropTables(db)
	database.MigrateModels(db)
	assert.NoError(t, database.DBAddAdminAndUserAndGuest(db))

	// prepare the content of the DB for testing
	// using the respective endpoints of the API
	_, _, simulationModelID, _, widgetID := addScenarioAndSimulatorAndSimulationModelAndDashboardAndWidget()

	// authenticate as userB who has no access to the elements in the DB
	token, err := helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserBCredentials)
	assert.NoError(t, err)

	// try to get all files for simulation model to which userB has not access
	// should return unprocessable entity error
	code, resp, err := helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/files?objectID=%v&objectType=model", simulationModelID), "GET", nil)
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
		fmt.Sprintf("/api/files?objectID=%v&objectType=wrongtype", simulationModelID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 400, code, "Response body: \n%v\n", resp)

	//try to get all files with missing object ID; should return a bad request error
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/files?objectType=model"), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 400, code, "Response body: \n%v\n", resp)

	// Count the number of all files returned for simulation model
	initialNumberModel, err := helper.LengthOfResponse(router, token,
		fmt.Sprintf("/api/files?objectID=%v&objectType=model", simulationModelID), "GET", nil)
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

	// test POST a file to simulation model and widget
	bodyBufModel1 := &bytes.Buffer{}
	bodyBufWidget1 := &bytes.Buffer{}
	bodyWriterModel1 := multipart.NewWriter(bodyBufModel1)
	bodyWriterWidget1 := multipart.NewWriter(bodyBufWidget1)
	fileWriterModel1, err := bodyWriterModel1.CreateFormFile("file", "testuploadfile.txt")
	assert.NoError(t, err, "writing to buffer")
	fileWriterWidget1, err := bodyWriterWidget1.CreateFormFile("file", "testuploadfile.txt")
	assert.NoError(t, err, "writing to buffer")
	// io copy
	_, err = io.Copy(fileWriterModel1, fh)
	assert.NoError(t, err, "IO copy")
	_, err = io.Copy(fileWriterWidget1, fh)
	assert.NoError(t, err, "IO copy")
	contentTypeModel1 := bodyWriterModel1.FormDataContentType()
	contentTypeWidget1 := bodyWriterWidget1.FormDataContentType()
	bodyWriterModel1.Close()
	bodyWriterWidget1.Close()

	// Create the request for simulation model
	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", fmt.Sprintf("/api/files?objectID=%v&objectType=model", simulationModelID), bodyBufModel1)
	assert.NoError(t, err, "create request")
	req.Header.Set("Content-Type", contentTypeModel1)
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

	// POST a second file to simulation model and widget

	// open a second file handle
	fh2, err := os.Open("testfile.txt")
	assert.NoError(t, err, "opening file")
	defer fh2.Close()

	bodyBufModel2 := &bytes.Buffer{}
	bodyBufWidget2 := &bytes.Buffer{}
	bodyWriterModel2 := multipart.NewWriter(bodyBufModel2)
	bodyWriterWidget2 := multipart.NewWriter(bodyBufWidget2)
	fileWriterModel2, err := bodyWriterModel2.CreateFormFile("file", "testuploadfile2.txt")
	assert.NoError(t, err, "writing to buffer")
	fileWriterWidget2, err := bodyWriterWidget2.CreateFormFile("file", "testuploadfile2.txt")
	assert.NoError(t, err, "writing to buffer")

	// io copy
	_, err = io.Copy(fileWriterModel2, fh2)
	assert.NoError(t, err, "IO copy")
	_, err = io.Copy(fileWriterWidget2, fh2)
	assert.NoError(t, err, "IO copy")
	contentTypeModel2 := bodyWriterModel2.FormDataContentType()
	contentTypeWidget2 := bodyWriterWidget2.FormDataContentType()
	bodyWriterModel2.Close()
	bodyWriterWidget2.Close()

	w3 := httptest.NewRecorder()
	req3, err := http.NewRequest("POST", fmt.Sprintf("/api/files?objectID=%v&objectType=model", simulationModelID), bodyBufModel2)
	assert.NoError(t, err, "create request")
	req3.Header.Set("Content-Type", contentTypeModel2)
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

	// Again count the number of all the files returned for simulation model
	finalNumberModel, err := helper.LengthOfResponse(router, token,
		fmt.Sprintf("/api/files?objectID=%v&objectType=model", simulationModelID), "GET", nil)
	assert.NoError(t, err)
	assert.Equal(t, initialNumberModel+2, finalNumberModel)

	// Again count the number of all the files returned for widget
	finalNumberWidget, err := helper.LengthOfResponse(router, token,
		fmt.Sprintf("/api/files?objectID=%v&objectType=widget", widgetID), "GET", nil)
	assert.NoError(t, err)
	assert.Equal(t, initialNumberWidget+2, finalNumberWidget)
}
