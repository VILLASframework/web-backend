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
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/scenario"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/user"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/stretchr/testify/assert"
)

var router *gin.Engine

type ScenarioRequest struct {
	Name            string         `json:"name,omitempty"`
	Running         bool           `json:"running,omitempty"`
	StartParameters postgres.Jsonb `json:"startParameters,omitempty"`
}

func addScenario() (scenarioID uint) {

	// authenticate as admin
	token, _ := helper.AuthenticateForTest(router, helper.AdminCredentials)

	// authenticate as normal user
	token, _ = helper.AuthenticateForTest(router, helper.UserACredentials)

	// POST $newScenario
	newScenario := ScenarioRequest{
		Name:            "Scenario1",
		Running:         true,
		StartParameters: postgres.Jsonb{RawMessage: json.RawMessage(`{"parameter1" : "testValue1A", "parameter2" : "testValue2A", "parameter3" : 42}`)},
	}
	_, resp, _ := helper.TestEndpoint(router, token,
		"/api/v2/scenarios", "POST", helper.KeyModels{"scenario": newScenario})

	// Read newScenario's ID from the response
	newScenarioID, _ := helper.GetResponseID(resp)

	// add the guest user to the new scenario
	_, resp, _ = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/scenarios/%v/user?username=User_C", newScenarioID), "PUT", nil)

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
	api := router.Group("/api/v2")

	user.RegisterAuthenticate(api.Group("/authenticate"))
	api.Use(user.Authentication())
	// scenario endpoints required here to first add a scenario to the DB
	scenario.RegisterScenarioEndpoints(api.Group("/scenarios"))

	RegisterFileEndpoints(api.Group("/files"))

	os.Exit(m.Run())
}

func TestAddFile(t *testing.T) {
	database.DropTables()
	database.MigrateModels()
	assert.NoError(t, helper.AddTestUsers())

	// prepare the content of the DB for testing
	// using the respective endpoints of the API
	scenarioID := addScenario()

	// authenticate as userB who has no access to the elements in the DB
	token, err := helper.AuthenticateForTest(router, helper.UserBCredentials)
	assert.NoError(t, err)

	emptyBuf := &bytes.Buffer{}

	// try to POST to a scenario to which UserB has no access
	// should return a 422 unprocessable entity error
	code, resp, err := helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/files?scenarioID=%v", scenarioID), "POST", emptyBuf)
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

	// authenticate as normal userA
	token, err = helper.AuthenticateForTest(router, helper.UserACredentials)
	assert.NoError(t, err)

	// try to POST without a scenario ID
	// should return a bad request error
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/files"), "POST", emptyBuf)
	assert.NoError(t, err)
	assert.Equalf(t, 400, code, "Response body: \n%v\n", resp)

	// try to POST an invalid file
	// should return a bad request
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/files?scenarioID=%v", scenarioID), "POST", emptyBuf)
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

	// Create the request
	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", fmt.Sprintf("/api/v2/files?scenarioID=%v", scenarioID), bodyBuf)
	assert.NoError(t, err, "create request")

	req.Header.Set("Content-Type", contentType)
	req.Header.Add("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	assert.Equalf(t, 200, w.Code, "Response body: \n%v\n", w.Body)

	newFileID, err := helper.GetResponseID(w.Body)
	assert.NoError(t, err)

	// Get the new file
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/files/%v", newFileID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)
	assert.Equalf(t, string(c1), resp.String(), "Response body: \n%v\n", resp)
}

func TestUpdateFile(t *testing.T) {

	database.DropTables()
	database.MigrateModels()
	assert.NoError(t, helper.AddTestUsers())

	// prepare the content of the DB for testing
	// using the respective endpoints of the API
	scenarioID := addScenario()

	// authenticate as normal user
	token, err := helper.AuthenticateForTest(router, helper.UserACredentials)
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

	// Create the POST request
	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", fmt.Sprintf("/api/v2/files?scenarioID=%v", scenarioID), bodyBuf)
	assert.NoError(t, err, "create request")

	req.Header.Set("Content-Type", contentType)
	req.Header.Add("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	assert.Equalf(t, 200, w.Code, "Response body: \n%v\n", w.Body)

	newFileID, err := helper.GetResponseID(w.Body)
	assert.NoError(t, err)

	// authenticate as userB who has no access to the elements in the DB
	token, err = helper.AuthenticateForTest(router, helper.UserBCredentials)
	assert.NoError(t, err)

	emptyBuf := &bytes.Buffer{}

	// try to PUT to a file to which UserB has no access
	// should return a 422 unprocessable entity error
	code, resp, err := helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/files/%v", newFileID), "PUT", emptyBuf)
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

	// authenticate as guest user C
	token, err = helper.AuthenticateForTest(router, helper.GuestCredentials)
	assert.NoError(t, err)

	// try to PUT as guest
	// should return an unprocessable entity error
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/files/%v", newFileID), "PUT", emptyBuf)
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

	// Prepare update
	// authenticate as normal user
	token, err = helper.AuthenticateForTest(router, helper.UserACredentials)
	assert.NoError(t, err)

	// try to PUT with empty body
	// should return bad request
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/files/%v", newFileID), "PUT", emptyBuf)
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
	req, err = http.NewRequest("PUT", fmt.Sprintf("/api/v2/files/%v", newFileID), bodyBufUpdated)
	assert.NoError(t, err, "create request")

	req.Header.Set("Content-Type", contentType)
	req.Header.Add("Authorization", "Bearer "+token)
	router.ServeHTTP(w_updated, req)
	assert.Equalf(t, 200, w_updated.Code, "Response body: \n%v\n", w_updated.Body)

	newFileIDUpdated, err := helper.GetResponseID(w_updated.Body)

	assert.Equal(t, newFileID, newFileIDUpdated)

	// Get the updated file
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/files/%v", newFileIDUpdated), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)
	assert.Equalf(t, string(c2), resp.String(), "Response body: \n%v\n", resp)
}

func TestDeleteFile(t *testing.T) {
	database.DropTables()
	database.MigrateModels()
	assert.NoError(t, helper.AddTestUsers())

	// prepare the content of the DB for testing
	// using the respective endpoints of the API
	scenarioID := addScenario()

	// authenticate as normal user
	token, err := helper.AuthenticateForTest(router, helper.UserACredentials)
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
	req, err := http.NewRequest("POST", fmt.Sprintf("/api/v2/files?scenarioID=%v", scenarioID), bodyBuf)
	assert.NoError(t, err, "create request")
	req.Header.Set("Content-Type", contentType)
	req.Header.Add("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)
	assert.Equalf(t, 200, w.Code, "Response body: \n%v\n", w.Body)

	newFileID, err := helper.GetResponseID(w.Body)
	assert.NoError(t, err)

	// add a second file to a scenario
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
	req2, err := http.NewRequest("POST", fmt.Sprintf("/api/v2/files?scenarioID=%v", scenarioID), bodyBuf2)
	assert.NoError(t, err, "create request")
	req2.Header.Set("Content-Type", contentType2)
	req2.Header.Add("Authorization", "Bearer "+token)
	router.ServeHTTP(w2, req2)
	assert.Equalf(t, 200, w2.Code, "Response body: \n%v\n", w2.Body)

	newFileID2, err := helper.GetResponseID(w2.Body)
	assert.NoError(t, err)

	// authenticate as userB who has no access to the elements in the DB
	token, err = helper.AuthenticateForTest(router, helper.UserBCredentials)
	assert.NoError(t, err)

	// try to DELETE file from scenario to which userB has no access
	// should return an unprocessable entity error
	code, resp, err := helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/files/%v", newFileID), "DELETE", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

	// authenticate as normal user
	token, err = helper.AuthenticateForTest(router, helper.UserACredentials)
	assert.NoError(t, err)

	// Count the number of all files returned for scenario
	initialNumber, err := helper.LengthOfResponse(router, token,
		fmt.Sprintf("/api/v2/files?scenarioID=%v", scenarioID), "GET", nil)
	assert.NoError(t, err)

	// try to DELETE non-existing fileID
	// should return not found
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/files/5"), "DELETE", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 404, code, "Response body: \n%v\n", resp)

	// authenticate as guest user C
	token, err = helper.AuthenticateForTest(router, helper.GuestCredentials)
	assert.NoError(t, err)

	// try to DELETE file of scenario as guest
	// should return an unprocessable entity error
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/files/%v", newFileID), "DELETE", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

	// authenticate as normal user
	token, err = helper.AuthenticateForTest(router, helper.UserACredentials)
	assert.NoError(t, err)

	// Delete the added file 1
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/files/%v", newFileID), "DELETE", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Delete the added file 2
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/files/%v", newFileID2), "DELETE", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Again count the number of all the files returned for scenario
	finalNumber, err := helper.LengthOfResponse(router, token,
		fmt.Sprintf("/api/v2/files?scenarioID=%v", scenarioID), "GET", nil)
	assert.NoError(t, err)

	assert.Equal(t, initialNumber-2, finalNumber)
}

func TestGetAllFilesOfScenario(t *testing.T) {

	database.DropTables()
	database.MigrateModels()
	assert.NoError(t, helper.AddTestUsers())

	// prepare the content of the DB for testing
	// using the respective endpoints of the API
	scenarioID := addScenario()

	// authenticate as userB who has no access to the elements in the DB
	token, err := helper.AuthenticateForTest(router, helper.UserBCredentials)
	assert.NoError(t, err)

	// try to get all files for scenario to which userB has not access
	// should return unprocessable entity error
	code, resp, err := helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/files?scenarioID=%v", scenarioID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

	// authenticate as normal userA
	token, err = helper.AuthenticateForTest(router, helper.UserACredentials)
	assert.NoError(t, err)

	//try to get all files with missing scenario ID; should return a bad request error
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/files"), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 400, code, "Response body: \n%v\n", resp)

	// Count the number of all files returned for scenario
	initialNumber, err := helper.LengthOfResponse(router, token,
		fmt.Sprintf("/api/v2/files?scenarioID=%v", scenarioID), "GET", nil)
	assert.NoError(t, err)

	// create a testfile.txt in local folder
	c1 := []byte("This is my testfile\n")
	err = ioutil.WriteFile("testfile.txt", c1, 0644)
	assert.NoError(t, err)

	// open file handle
	fh, err := os.Open("testfile.txt")
	assert.NoError(t, err, "opening file")
	defer fh.Close()

	// test POST a file to scenario
	bodyBuf1 := &bytes.Buffer{}
	bodyWriter1 := multipart.NewWriter(bodyBuf1)
	fileWriter1, err := bodyWriter1.CreateFormFile("file", "testuploadfile.txt")
	assert.NoError(t, err, "writing to buffer")
	// io copy
	_, err = io.Copy(fileWriter1, fh)
	assert.NoError(t, err, "IO copy")
	contentType1 := bodyWriter1.FormDataContentType()
	bodyWriter1.Close()

	// Create the request
	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", fmt.Sprintf("/api/v2/files?scenarioID=%v", scenarioID), bodyBuf1)
	assert.NoError(t, err, "create request")
	req.Header.Set("Content-Type", contentType1)
	req.Header.Add("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)
	assert.Equalf(t, 200, w.Code, "Response body: \n%v\n", w.Body)

	// POST a second file to scenario

	// open a second file handle
	fh2, err := os.Open("testfile.txt")
	assert.NoError(t, err, "opening file")
	defer fh2.Close()

	bodyBuf2 := &bytes.Buffer{}
	bodyWriter2 := multipart.NewWriter(bodyBuf2)
	fileWriter2, err := bodyWriter2.CreateFormFile("file", "testuploadfile2.txt")
	assert.NoError(t, err, "writing to buffer")

	// io copy
	_, err = io.Copy(fileWriter2, fh2)
	assert.NoError(t, err, "IO copy")
	contentType2 := bodyWriter2.FormDataContentType()
	bodyWriter2.Close()

	w2 := httptest.NewRecorder()
	req2, err := http.NewRequest("POST", fmt.Sprintf("/api/v2/files?scenarioID=%v", scenarioID), bodyBuf2)
	assert.NoError(t, err, "create request")
	req2.Header.Set("Content-Type", contentType2)
	req2.Header.Add("Authorization", "Bearer "+token)
	router.ServeHTTP(w2, req2)
	assert.Equalf(t, 200, w2.Code, "Response body: \n%v\n", w2.Body)

	// Again count the number of all the files returned for scenario
	finalNumber, err := helper.LengthOfResponse(router, token,
		fmt.Sprintf("/api/v2/files?scenarioID=%v", scenarioID), "GET", nil)
	assert.NoError(t, err)
	assert.Equal(t, initialNumber+2, finalNumber)
}
