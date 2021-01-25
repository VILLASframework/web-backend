/** Result package, testing.
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

package result

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"git.rwth-aachen.de/acs/public/villas/web-backend-go/configuration"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/helper"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/file"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/scenario"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/user"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/stretchr/testify/assert"
)

var router *gin.Engine
var base_api_results = "/api/results"
var base_api_auth = "/api/authenticate"

type ScenarioRequest struct {
	Name            string         `json:"name,omitempty"`
	Running         bool           `json:"running,omitempty"`
	StartParameters postgres.Jsonb `json:"startParameters,omitempty"`
}

type ResultRequest struct {
	Description     string         `json:"description,omitempty"`
	ScenarioID      uint           `json:"scenarioID,omitempty"`
	ConfigSnapshots postgres.Jsonb `json:"configSnapshots,omitempty"`
}

type ResponseResult struct {
	Result database.Result `json:"result"`
}

var newResult = ResultRequest{
	Description: "This is a test result.",
}

func addScenario() (scenarioID uint) {

	// authenticate as admin
	token, _ := helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.AdminCredentials)

	// authenticate as normal user
	token, _ = helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserACredentials)

	// POST $newScenario
	newScenario := ScenarioRequest{
		Name:            "Scenario1",
		Running:         true,
		StartParameters: postgres.Jsonb{RawMessage: json.RawMessage(`{"parameter1" : "testValue1A", "parameter2" : "testValue2A", "parameter3" : 42}`)},
	}
	_, resp, _ := helper.TestEndpoint(router, token,
		"/api/scenarios", "POST", helper.KeyModels{"scenario": newScenario})

	// Read newScenario's ID from the response
	newScenarioID, _ := helper.GetResponseID(resp)

	// add the guest user to the new scenario
	_, resp, _ = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/scenarios/%v/user?username=User_C", newScenarioID), "PUT", nil)

	return uint(newScenarioID)
}

func TestMain(m *testing.M) {
	err := configuration.InitConfig()
	if err != nil {
		panic(m)
	}
	err = database.InitDB(configuration.GlobalConfig)
	if err != nil {
		panic(m)
	}
	defer database.DBpool.Close()

	router = gin.Default()
	api := router.Group("/api")

	user.RegisterAuthenticate(api.Group("/authenticate"))
	api.Use(user.Authentication())
	// scenario endpoints required here to first add a scenario to the DB
	scenario.RegisterScenarioEndpoints(api.Group("/scenarios"))
	// file endpoints required to download result file
	file.RegisterFileEndpoints(api.Group("/files"))

	RegisterResultEndpoints(api.Group("/results"))

	os.Exit(m.Run())
}

func TestGetAllResultsOfScenario(t *testing.T) {

	database.DropTables()
	database.MigrateModels()
	assert.NoError(t, helper.AddTestUsers())

	// prepare the content of the DB for testing
	// by adding a scenario
	scenarioID := addScenario()

	// authenticate as normal user
	token, err := helper.AuthenticateForTest(router,
		base_api_auth, "POST", helper.UserACredentials)
	assert.NoError(t, err)

	// test POST newResult
	configSnapshot1 := json.RawMessage(`{"configs": [ {"Name" : "conf1", "scenarioID" : 1}, {"Name" : "conf2", "scenarioID" : 1}]}`)
	confSnapshots := postgres.Jsonb{configSnapshot1}

	newResult.ScenarioID = scenarioID
	newResult.ConfigSnapshots = confSnapshots
	code, resp, err := helper.TestEndpoint(router, token,
		base_api_results, "POST", helper.KeyModels{"result": newResult})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Count the number of all the results returned for scenario
	NumberOfConfigs, err := helper.LengthOfResponse(router, token,
		fmt.Sprintf("%v?scenarioID=%v", base_api_results, scenarioID), "GET", nil)
	assert.NoError(t, err)

	assert.Equal(t, 1, NumberOfConfigs)

	// authenticate as normal userB who has no access to scenario
	token, err = helper.AuthenticateForTest(router,
		base_api_auth, "POST", helper.UserBCredentials)
	assert.NoError(t, err)

	// try to get results without access
	// should result in unprocessable entity
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("%v?scenarioID=%v", base_api_results, scenarioID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

}

func TestAddGetUpdateDeleteResult(t *testing.T) {

	database.DropTables()
	database.MigrateModels()
	assert.NoError(t, helper.AddTestUsers())

	// prepare the content of the DB for testing
	// by adding a scenario
	scenarioID := addScenario()
	configSnapshot1 := json.RawMessage(`{"configs": [ {"Name" : "conf1", "scenarioID" : 1}, {"Name" : "conf2", "scenarioID" : 1}]}`)
	confSnapshots := postgres.Jsonb{configSnapshot1}
	newResult.ScenarioID = scenarioID
	newResult.ConfigSnapshots = confSnapshots
	// authenticate as normal userB who has no access to new scenario
	token, err := helper.AuthenticateForTest(router,
		base_api_auth, "POST", helper.UserBCredentials)
	assert.NoError(t, err)

	// try to POST with no access
	// should result in unprocessable entity
	code, resp, err := helper.TestEndpoint(router, token,
		base_api_results, "POST", helper.KeyModels{"result": newResult})
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

	// authenticate as normal user
	token, err = helper.AuthenticateForTest(router,
		base_api_auth, "POST", helper.UserACredentials)
	assert.NoError(t, err)

	// try to POST non JSON body
	code, resp, err = helper.TestEndpoint(router, token,
		base_api_results, "POST", "this is not JSON")
	assert.NoError(t, err)
	assert.Equalf(t, 400, code, "Response body: \n%v\n", resp)

	// test POST newResult
	code, resp, err = helper.TestEndpoint(router, token,
		base_api_results, "POST", helper.KeyModels{"result": newResult})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare POST's response with the newResult
	err = helper.CompareResponse(resp, helper.KeyModels{"result": newResult})
	assert.NoError(t, err)

	// Read newResults's ID from the response
	newResultID, err := helper.GetResponseID(resp)
	assert.NoError(t, err)

	// Get the newResult
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("%v/%v", base_api_results, newResultID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare GET's response with the newResult
	err = helper.CompareResponse(resp, helper.KeyModels{"result": newResult})
	assert.NoError(t, err)

	// try to POST a malformed result
	// Required fields are missing
	malformedNewResult := ResultRequest{
		Description: "ThisIsAMalformedRequest",
	}
	// this should NOT work and return a unprocessable entity 442 status code
	code, resp, err = helper.TestEndpoint(router, token,
		base_api_results, "POST", helper.KeyModels{"result": malformedNewResult})
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

	// authenticate as normal userB who has no access to new scenario
	token, err = helper.AuthenticateForTest(router,
		base_api_auth, "POST", helper.UserBCredentials)
	assert.NoError(t, err)

	// Try to GET the newResult with no access
	// Should result in unprocessable entity
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("%v/%v", base_api_results, newResultID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

	// Test UPDATE/ PUT

	updatedResult := ResultRequest{
		Description:     "This is an updated description",
		ConfigSnapshots: confSnapshots,
	}

	// try to PUT with no access
	// should result in unprocessable entity
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("%v/%v", base_api_results, newResultID), "PUT", helper.KeyModels{"result": updatedResult})
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

	// authenticate as guest user who has access to result
	token, err = helper.AuthenticateForTest(router,
		base_api_auth, "POST", helper.GuestCredentials)
	assert.NoError(t, err)

	// try to PUT as guest
	// should NOT work and result in unprocessable entity
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("%v/%v", base_api_results, newResultID), "PUT", helper.KeyModels{"result": updatedResult})
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

	// authenticate as normal user
	token, err = helper.AuthenticateForTest(router,
		base_api_auth, "POST", helper.UserACredentials)
	assert.NoError(t, err)

	// try to PUT a non JSON body
	// should result in a bad request
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("%v/%v", base_api_results, newResultID), "PUT", "This is not JSON")
	assert.NoError(t, err)
	assert.Equalf(t, 400, code, "Response body: \n%v\n", resp)

	// test PUT
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("%v/%v", base_api_results, newResultID), "PUT", helper.KeyModels{"result": updatedResult})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare PUT's response with the updatedResult
	err = helper.CompareResponse(resp, helper.KeyModels{"result": updatedResult})
	assert.NoError(t, err)

	// try to update a result that does not exist (should return not found 404 status code)
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("%v/%v", base_api_results, newResultID+1), "PUT", helper.KeyModels{"result": updatedResult})
	assert.NoError(t, err)
	assert.Equalf(t, 404, code, "Response body: \n%v\n", resp)

	// Test DELETE
	newResult.Description = updatedResult.Description

	// authenticate as normal userB who has no access to new scenario
	token, err = helper.AuthenticateForTest(router,
		base_api_auth, "POST", helper.UserBCredentials)
	assert.NoError(t, err)

	// try to DELETE with no access
	// should result in unprocessable entity
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("%v/%v", base_api_results, newResultID), "DELETE", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

	// authenticate as normal user
	token, err = helper.AuthenticateForTest(router,
		base_api_auth, "POST", helper.UserACredentials)
	assert.NoError(t, err)

	// Count the number of all the results returned for scenario
	initialNumber, err := helper.LengthOfResponse(router, token,
		fmt.Sprintf("%v?scenarioID=%v", base_api_results, scenarioID), "GET", nil)
	assert.NoError(t, err)

	// Delete the added newResult
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("%v/%v", base_api_results, newResultID), "DELETE", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare DELETE's response with the newResult
	err = helper.CompareResponse(resp, helper.KeyModels{"result": newResult})
	assert.NoError(t, err)

	// Again count the number of all the results returned
	finalNumber, err := helper.LengthOfResponse(router, token,
		fmt.Sprintf("%v?scenarioID=%v", base_api_results, scenarioID), "GET", nil)
	assert.NoError(t, err)

	assert.Equal(t, initialNumber-1, finalNumber)

}

func TestAddDeleteResultFile(t *testing.T) {
	database.DropTables()
	database.MigrateModels()
	assert.NoError(t, helper.AddTestUsers())

	// prepare the content of the DB for testing
	// by adding a scenario
	scenarioID := addScenario()
	configSnapshot1 := json.RawMessage(`{"configs": [ {"Name" : "conf1", "scenarioID" : 1}, {"Name" : "conf2", "scenarioID" : 1}]}`)
	confSnapshots := postgres.Jsonb{configSnapshot1}

	newResult.ScenarioID = scenarioID
	newResult.ConfigSnapshots = confSnapshots
	// authenticate as normal user
	token, err := helper.AuthenticateForTest(router,
		base_api_auth, "POST", helper.UserACredentials)
	assert.NoError(t, err)

	// test POST newResult
	code, resp, err := helper.TestEndpoint(router, token,
		base_api_results, "POST", helper.KeyModels{"result": newResult})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare POST's response with the newResult
	err = helper.CompareResponse(resp, helper.KeyModels{"result": newResult})
	assert.NoError(t, err)

	// Read newResults's ID from the response
	newResultID, err := helper.GetResponseID(resp)
	assert.NoError(t, err)

	// test POST result file

	// create a testfile.txt in local folder
	c1 := []byte("a,few,values\n1,2,3\n")
	err = ioutil.WriteFile("testfile.csv", c1, 0644)
	assert.NoError(t, err)

	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	fileWriter, err := bodyWriter.CreateFormFile("file", "testuploadfile.csv")
	assert.NoError(t, err, "writing to buffer")

	// open file handle
	fh, err := os.Open("testfile.csv")
	assert.NoError(t, err, "opening file")
	defer fh.Close()

	// io copy
	_, err = io.Copy(fileWriter, fh)
	assert.NoError(t, err, "IO copy")

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	// Create the request
	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", fmt.Sprintf("%v/%v/file", base_api_results, newResultID), bodyBuf)
	assert.NoError(t, err, "create request")

	req.Header.Set("Content-Type", contentType)
	req.Header.Add("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	assert.Equalf(t, 200, w.Code, "Response body: \n%v\n", w.Body)
	err = helper.CompareResponse(w.Body, helper.KeyModels{"result": newResult})

	// extract file ID from response body
	var respResult ResponseResult
	err = json.Unmarshal(w.Body.Bytes(), &respResult)
	assert.NoError(t, err, "unmarshal response body")

	assert.Equal(t, 1, len(respResult.Result.ResultFileIDs))
	fileID := respResult.Result.ResultFileIDs[0]

	// DELETE the file

	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("%v/%v/file/%v", base_api_results, newResultID, fileID), "DELETE", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	var respResult2 ResponseResult
	err = json.Unmarshal(resp.Bytes(), &respResult2)
	assert.NoError(t, err, "unmarshal response body")
	assert.Equal(t, 0, len(respResult2.Result.ResultFileIDs))

	// ADD the file again

	bodyBuf2 := &bytes.Buffer{}
	bodyWriter2 := multipart.NewWriter(bodyBuf2)
	fileWriter2, err := bodyWriter2.CreateFormFile("file", "testuploadfile.csv")
	assert.NoError(t, err, "writing to buffer")

	// open file handle
	fh2, err := os.Open("testfile.csv")
	assert.NoError(t, err, "opening file")
	defer fh2.Close()

	// io copy
	_, err = io.Copy(fileWriter2, fh2)
	assert.NoError(t, err, "IO copy")

	contentType2 := bodyWriter2.FormDataContentType()
	bodyWriter2.Close()

	// Create the request
	w2 := httptest.NewRecorder()
	req2, err := http.NewRequest("POST", fmt.Sprintf("%v/%v/file", base_api_results, newResultID), bodyBuf2)
	assert.NoError(t, err, "create request")

	req2.Header.Set("Content-Type", contentType2)
	req2.Header.Add("Authorization", "Bearer "+token)
	router.ServeHTTP(w2, req2)

	assert.Equalf(t, 200, w2.Code, "Response body: \n%v\n", w2.Body)
	err = helper.CompareResponse(w2.Body, helper.KeyModels{"result": newResult})

	// extract file ID from response body
	var respResult3 ResponseResult
	err = json.Unmarshal(w2.Body.Bytes(), &respResult3)
	assert.NoError(t, err, "unmarshal response body")

	assert.Equal(t, 1, len(respResult3.Result.ResultFileIDs))

	// DELETE result inlc. file
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("%v/%v", base_api_results, newResultID), "DELETE", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare DELETE's response with the newResult
	err = helper.CompareResponse(resp, helper.KeyModels{"result": newResult})
	assert.NoError(t, err)

	// Again count the number of all the results returned
	finalNumber, err := helper.LengthOfResponse(router, token,
		fmt.Sprintf("%v?scenarioID=%v", base_api_results, scenarioID), "GET", nil)
	assert.NoError(t, err)

	assert.Equal(t, 0, finalNumber)

}
