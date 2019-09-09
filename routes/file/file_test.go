package file

import (
	"bytes"
	"fmt"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/scenario"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/simulationmodel"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/simulator"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/user"
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

func addScenarioAndSimulatorAndSimulationModel() (scenarioID uint, simulatorID uint, simulationModelID uint) {

	// authenticate as admin
	token, _ := common.AuthenticateForTest(router,
		"/api/authenticate", "POST", common.AdminCredentials)

	// POST $newSimulatorA
	newSimulatorA := SimulatorRequest{
		UUID:       common.SimulatorA.UUID,
		Host:       common.SimulatorA.Host,
		Modeltype:  common.SimulatorA.Modeltype,
		State:      common.SimulatorA.State,
		Properties: common.SimulatorA.Properties,
	}
	_, resp, _ := common.TestEndpoint(router, token,
		"/api/simulators", "POST", common.KeyModels{"simulator": newSimulatorA})

	// Read newSimulator's ID from the response
	newSimulatorID, _ := common.GetResponseID(resp)

	// authenticate as normal user
	token, _ = common.AuthenticateForTest(router,
		"/api/authenticate", "POST", common.UserACredentials)

	// POST $newScenario
	newScenario := ScenarioRequest{
		Name:            common.ScenarioA.Name,
		Running:         common.ScenarioA.Running,
		StartParameters: common.ScenarioA.StartParameters,
	}
	_, resp, _ = common.TestEndpoint(router, token,
		"/api/scenarios", "POST", common.KeyModels{"scenario": newScenario})

	// Read newScenario's ID from the response
	newScenarioID, _ := common.GetResponseID(resp)

	// test POST models/ $newSimulationModel
	newSimulationModel := SimulationModelRequest{
		Name:            common.SimulationModelA.Name,
		ScenarioID:      uint(newScenarioID),
		SimulatorID:     uint(newSimulatorID),
		StartParameters: common.SimulationModelA.StartParameters,
	}
	_, resp, _ = common.TestEndpoint(router, token,
		"/api/models", "POST", common.KeyModels{"model": newSimulationModel})

	// Read newSimulationModel's ID from the response
	newSimulationModelID, _ := common.GetResponseID(resp)

	return uint(newScenarioID), uint(newSimulatorID), uint(newSimulationModelID)
}

func TestMain(m *testing.M) {

	db = common.DummyInitDB()
	defer db.Close()

	router = gin.Default()
	api := router.Group("/api")

	user.RegisterAuthenticate(api.Group("/authenticate"))
	api.Use(user.Authentication(true))
	// simulationmodel endpoints required here to first add a simulation to the DB
	// that can be associated with a new signal model
	simulationmodel.RegisterSimulationModelEndpoints(api.Group("/models"))
	// scenario endpoints required here to first add a scenario to the DB
	// that can be associated with a new simulation model
	scenario.RegisterScenarioEndpoints(api.Group("/scenarios"))
	// simulator endpoints required here to first add a simulator to the DB
	// that can be associated with a new simulation model
	simulator.RegisterSimulatorEndpoints(api.Group("/simulators"))
	RegisterFileEndpoints(api.Group("/files"))

	os.Exit(m.Run())
}

func TestAddFile(t *testing.T) {
	common.DropTables(db)
	common.MigrateModels(db)
	common.DummyAddOnlyUserTableWithAdminAndUsersDB(db)

	// prepare the content of the DB for testing
	// by adding a scenario and a simulator to the DB
	// using the respective endpoints of the API
	_, _, simulationModelID := addScenarioAndSimulatorAndSimulationModel()

	// authenticate as normal user
	token, err := common.AuthenticateForTest(router,
		"/api/authenticate", "POST", common.UserACredentials)
	assert.NoError(t, err)

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
	fmt.Println(w.Body)

	newFileID, err := common.GetResponseID(w.Body)
	assert.NoError(t, err)

	// Get the new file
	code, resp, err := common.TestEndpoint(router, token,
		fmt.Sprintf("/api/files/%v", newFileID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

}

func TestUpdateFile(t *testing.T) {

	common.DropTables(db)
	common.MigrateModels(db)
	common.DummyAddOnlyUserTableWithAdminAndUsersDB(db)

	// prepare the content of the DB for testing
	// by adding a scenario and a simulator to the DB
	// using the respective endpoints of the API
	_, _, simulationModelID := addScenarioAndSimulatorAndSimulationModel()

	// authenticate as normal user
	token, err := common.AuthenticateForTest(router,
		"/api/authenticate", "POST", common.UserACredentials)
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
	fmt.Println(w.Body)

	newFileID, err := common.GetResponseID(w.Body)
	assert.NoError(t, err)

	// Prepare update

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
	fmt.Println(w_updated.Body)

	newFileIDUpdated, err := common.GetResponseID(w_updated.Body)

	assert.Equal(t, newFileID, newFileIDUpdated)

	// Get the updated file
	code, resp, err := common.TestEndpoint(router, token,
		fmt.Sprintf("/api/files/%v", newFileIDUpdated), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

}

func TestDeleteFile(t *testing.T) {
	common.DropTables(db)
	common.MigrateModels(db)
	common.DummyAddOnlyUserTableWithAdminAndUsersDB(db)

	// prepare the content of the DB for testing
	// by adding a scenario and a simulator to the DB
	// using the respective endpoints of the API
	_, _, simulationModelID := addScenarioAndSimulatorAndSimulationModel()

	// authenticate as normal user
	token, err := common.AuthenticateForTest(router,
		"/api/authenticate", "POST", common.UserACredentials)
	assert.NoError(t, err)

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

	newFileID, err := common.GetResponseID(w.Body)
	assert.NoError(t, err)

	// Count the number of all files returned for simulation model
	initialNumber, err := common.LengthOfResponse(router, token,
		fmt.Sprintf("/api/files?objectID=%v&objectType=model", simulationModelID), "GET", nil)
	assert.NoError(t, err)

	// Delete the added file
	code, resp, err := common.TestEndpoint(router, token,
		fmt.Sprintf("/api/files/%v", newFileID), "DELETE", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Again count the number of all the files returned for simulation model
	finalNumber, err := common.LengthOfResponse(router, token,
		fmt.Sprintf("/api/files?objectID=%v&objectType=model", simulationModelID), "GET", nil)
	assert.NoError(t, err)

	assert.Equal(t, initialNumber-1, finalNumber)
}

func TestGetAllFilesOfSimulationModel(t *testing.T) {

	common.DropTables(db)
	common.MigrateModels(db)
	common.DummyAddOnlyUserTableWithAdminAndUsersDB(db)

	// prepare the content of the DB for testing
	// by adding a scenario and a simulator to the DB
	// using the respective endpoints of the API
	_, _, simulationModelID := addScenarioAndSimulatorAndSimulationModel()

	// authenticate as normal user
	token, err := common.AuthenticateForTest(router,
		"/api/authenticate", "POST", common.UserACredentials)
	assert.NoError(t, err)

	// Count the number of all files returned for simulation model
	initialNumber, err := common.LengthOfResponse(router, token,
		fmt.Sprintf("/api/files?objectID=%v&objectType=model", simulationModelID), "GET", nil)
	assert.NoError(t, err)

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

	// POST a second file

	bodyBuf2 := &bytes.Buffer{}
	bodyWriter2 := multipart.NewWriter(bodyBuf2)
	fileWriter2, err := bodyWriter2.CreateFormFile("file", "testuploadfile2.txt")
	assert.NoError(t, err, "writing to buffer")

	// open file handle
	fh2, err := os.Open("testfile.txt")
	assert.NoError(t, err, "opening file")
	defer fh2.Close()

	// io copy
	_, err = io.Copy(fileWriter2, fh2)
	assert.NoError(t, err, "IO copy")

	contentType = bodyWriter2.FormDataContentType()
	bodyWriter2.Close()

	w2 := httptest.NewRecorder()
	req, err = http.NewRequest("POST", fmt.Sprintf("/api/files?objectID=%v&objectType=model", simulationModelID), bodyBuf2)
	assert.NoError(t, err, "create request")

	req.Header.Set("Content-Type", contentType)
	req.Header.Add("Authorization", "Bearer "+token)
	router.ServeHTTP(w2, req)
	assert.Equalf(t, 200, w2.Code, "Response body: \n%v\n", w2.Body)

	// Again count the number of all the files returned for simulation model
	finalNumber, err := common.LengthOfResponse(router, token,
		fmt.Sprintf("/api/files?objectID=%v&objectType=model", simulationModelID), "GET", nil)
	assert.NoError(t, err)

	assert.Equal(t, initialNumber+2, finalNumber)
}
