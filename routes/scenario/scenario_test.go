/** Scenario package, testing.
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
package scenario

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"git.rwth-aachen.de/acs/public/villas/web-backend-go/configuration"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/helper"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/user"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/stretchr/testify/assert"
)

var router *gin.Engine

type ScenarioRequest struct {
	Name            string         `json:"name,omitempty"`
	IsLocked        bool           `json:"isLocked,omitempty"`
	StartParameters postgres.Jsonb `json:"startParameters,omitempty"`
}

type UserRequest struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Mail     string `json:"mail,omitempty"`
	Role     string `json:"role,omitempty"`
	Active   string `json:"active,omitempty"`
}

var newScenario1 = ScenarioRequest{
	Name:            "Scenario1",
	StartParameters: postgres.Jsonb{json.RawMessage(`{"parameter1" : "testValue1A", "parameter2" : "testValue2A", "parameter3" : 42}`)},
}

var newScenario2 = ScenarioRequest{
	Name:            "Scenario2",
	StartParameters: postgres.Jsonb{json.RawMessage(`{"parameter1" : "testValue1B", "parameter2" : "testValue2B", "parameter3" : 55}`)},
}

func TestMain(m *testing.M) {
	err := configuration.InitConfig()
	if err != nil {
		panic(m)
	}

	err = database.InitDB(configuration.GlobalConfig, "true")
	if err != nil {
		panic(m)
	}
	defer database.DBpool.Close()

	router = gin.Default()
	api := router.Group("/api/v2")

	user.RegisterAuthenticate(api.Group("/authenticate"))
	api.Use(user.Authentication())

	// user endpoints required to set user to inactive
	user.RegisterUserEndpoints(api.Group("/users"))
	RegisterScenarioEndpoints(api.Group("/scenarios"))

	os.Exit(m.Run())
}

func TestAddScenario(t *testing.T) {

	database.DropTables()
	database.MigrateModels()
	assert.NoError(t, helper.AddTestUsers())

	// authenticate as normal user
	token, err := helper.AuthenticateForTest(router, helper.UserACredentials)
	assert.NoError(t, err)

	// try to POST with non JSON body
	// should return a bad request error
	code, resp, err := helper.TestEndpoint(router, token,
		"/api/v2/scenarios", "POST", "this is not a JSON")
	assert.NoError(t, err)
	assert.Equalf(t, 400, code, "Response body: \n%v\n", resp)

	// test POST scenarios/ $newScenario as normal user
	code, resp, err = helper.TestEndpoint(router, token,
		"/api/v2/scenarios", "POST", helper.KeyModels{"scenario": newScenario1})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare POST's response with the newScenario
	err = helper.CompareResponse(resp, helper.KeyModels{"scenario": newScenario1})
	assert.NoError(t, err)

	// Read newScenario's ID from the response
	newScenarioID, err := helper.GetResponseID(resp)
	assert.NoError(t, err)

	// Get the newScenario
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/scenarios/%v", newScenarioID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare GET's response with the newScenario
	err = helper.CompareResponse(resp, helper.KeyModels{"scenario": newScenario1})
	assert.NoError(t, err)

	// try to POST a malformed scenario
	// Required fields are missing
	malformedNewScenario := ScenarioRequest{
		IsLocked: false,
	}
	// this should NOT work and return a unprocessable entity 442 status code
	code, resp, err = helper.TestEndpoint(router, token,
		"/api/v2/scenarios", "POST", helper.KeyModels{"scenario": malformedNewScenario})
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

	// try to GET a non-existing scenario
	// should return a not found 404
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/scenarios/%v", newScenarioID+1), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 404, code, "Response body: \n%v\n", resp)

	// authenticate as guest user
	token, err = helper.AuthenticateForTest(router, helper.GuestCredentials)
	assert.NoError(t, err)

	// try to add scenario as guest user
	// should return an unprocessable entity error
	code, resp, err = helper.TestEndpoint(router, token,
		"/api/v2/scenarios", "POST", helper.KeyModels{"scenario": newScenario1})
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

	// authenticate as userB who has no access to the added scenario
	token, err = helper.AuthenticateForTest(router, helper.UserBCredentials)
	assert.NoError(t, err)

	// try to GET a scenario to which user B has no access
	// should return an unprocessable entity error
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/scenarios/%v", newScenarioID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

	// authenticate as admin user who has no access to everything
	token, err = helper.AuthenticateForTest(router, helper.AdminCredentials)
	assert.NoError(t, err)

	// try to GET a scenario that is not created by admin user; should work anyway
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/scenarios/%v", newScenarioID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)
}

func TestUpdateScenario(t *testing.T) {

	database.DropTables()
	database.MigrateModels()
	assert.NoError(t, helper.AddTestUsers())

	// authenticate as normal user
	token, err := helper.AuthenticateForTest(router, helper.UserACredentials)
	assert.NoError(t, err)

	// test POST scenarios/ $newScenario
	code, resp, err := helper.TestEndpoint(router, token,
		"/api/v2/scenarios", "POST", helper.KeyModels{"scenario": newScenario1})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare POST's response with the newScenario
	err = helper.CompareResponse(resp, helper.KeyModels{"scenario": newScenario1})
	assert.NoError(t, err)

	// Read newScenario's ID from the response
	newScenarioID, err := helper.GetResponseID(resp)
	assert.NoError(t, err)

	// try to update with non JSON body
	// should return a bad request error
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/scenarios/%v", newScenarioID), "PUT", "This is not a JSON body")
	assert.NoError(t, err)
	assert.Equalf(t, 400, code, "Response body: \n%v\n", resp)

	updatedScenario := ScenarioRequest{
		Name:            "Updated name",
		IsLocked:        true,
		StartParameters: postgres.Jsonb{RawMessage: json.RawMessage(`{"parameter1" : "testValue1A", "parameter2" : "testValue2A", "parameter3" : 42}`)},
	}

	// try to change locked state as non admin user
	// should return 200 but locked state not updated
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/scenarios/%v", newScenarioID), "PUT", helper.KeyModels{"scenario": updatedScenario})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare PUT's response with the updatedScenario (should result in error)
	err = helper.CompareResponse(resp, helper.KeyModels{"scenario": updatedScenario})
	assert.Error(t, err)

	updatedScenario.IsLocked = false
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/scenarios/%v", newScenarioID), "PUT", helper.KeyModels{"scenario": updatedScenario})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare PUT's response with the updatedScenario
	err = helper.CompareResponse(resp, helper.KeyModels{"scenario": updatedScenario})
	assert.NoError(t, err)

	// Get the updatedScenario
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/scenarios/%v", newScenarioID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare GET's response with the newScenario
	err = helper.CompareResponse(resp, helper.KeyModels{"scenario": updatedScenario})
	assert.NoError(t, err)

	// try to update a scenario that does not exist (should return not found 404 status code)
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/scenarios/%v", newScenarioID+1), "PUT", helper.KeyModels{"scenario": updatedScenario})
	assert.NoError(t, err)
	assert.Equalf(t, 404, code, "Response body: \n%v\n", resp)

	// authenticate as admin user who has no access to everything
	token, err = helper.AuthenticateForTest(router, helper.AdminCredentials)
	assert.NoError(t, err)

	// changed locked state of scenario as admin user (should work)
	updatedScenario.IsLocked = true
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/scenarios/%v", newScenarioID), "PUT", helper.KeyModels{"scenario": updatedScenario})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare PUT's response with the updatedScenario
	err = helper.CompareResponse(resp, helper.KeyModels{"scenario": updatedScenario})
	assert.NoError(t, err)

	// Get the updatedScenario
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/scenarios/%v", newScenarioID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare GET's response with the newScenario
	err = helper.CompareResponse(resp, helper.KeyModels{"scenario": updatedScenario})
	assert.NoError(t, err)

	// change a locked scenario as admin user (should work)
	updatedScenario.Name = "Updated as admin"
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/scenarios/%v", newScenarioID), "PUT", helper.KeyModels{"scenario": updatedScenario})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare GET's response with the updatedScenario
	err = helper.CompareResponse(resp, helper.KeyModels{"scenario": updatedScenario})
	assert.NoError(t, err)

	// authenticate as normal user
	token, err = helper.AuthenticateForTest(router, helper.UserACredentials)
	assert.NoError(t, err)

	// Get the updatedScenario
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/scenarios/%v", newScenarioID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare GET's response with the updatedScenario
	err = helper.CompareResponse(resp, helper.KeyModels{"scenario": updatedScenario})
	assert.NoError(t, err)

	// try to change a locked scenario as normal user (should result in unprocessable entity error)
	updatedScenario.Name = "another new name"
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/scenarios/%v", newScenarioID), "PUT", helper.KeyModels{"scenario": updatedScenario})
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

}

func TestGetAllScenariosAsAdmin(t *testing.T) {

	database.DropTables()
	database.MigrateModels()
	assert.NoError(t, helper.AddTestUsers())

	// authenticate as admin
	token, err := helper.AuthenticateForTest(router, helper.AdminCredentials)
	assert.NoError(t, err)

	// get the length of the GET all scenarios response for admin
	initialNumber, err := helper.LengthOfResponse(router, token,
		"/api/v2/scenarios", "GET", nil)
	assert.NoError(t, err)

	// authenticate as normal userB
	token, err = helper.AuthenticateForTest(router, helper.UserBCredentials)
	assert.NoError(t, err)

	// test POST scenarios/ $newScenario1
	code, resp, err := helper.TestEndpoint(router, token,
		"/api/v2/scenarios", "POST", helper.KeyModels{"scenario": newScenario1})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// authenticate as normal userA
	token, err = helper.AuthenticateForTest(router, helper.UserACredentials)
	assert.NoError(t, err)

	// test POST scenarios/ $newScenario2
	code, resp, err = helper.TestEndpoint(router, token,
		"/api/v2/scenarios", "POST", helper.KeyModels{"scenario": newScenario2})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// authenticate as admin
	token, err = helper.AuthenticateForTest(router, helper.AdminCredentials)
	assert.NoError(t, err)

	// get the length of the GET all scenarios response again
	finalNumber, err := helper.LengthOfResponse(router, token,
		"/api/v2/scenarios", "GET", nil)
	assert.NoError(t, err)

	assert.Equal(t, finalNumber, initialNumber+2)
}

func TestGetAllScenariosAsUser(t *testing.T) {

	database.DropTables()
	database.MigrateModels()
	assert.NoError(t, helper.AddTestUsers())

	// authenticate as normal userB
	token, err := helper.AuthenticateForTest(router, helper.UserBCredentials)
	assert.NoError(t, err)

	// get the length of the GET all scenarios response for userB
	initialNumber, err := helper.LengthOfResponse(router, token,
		"/api/v2/scenarios", "GET", nil)
	assert.NoError(t, err)

	// test POST scenarios/ $newScenario2
	code, resp, err := helper.TestEndpoint(router, token,
		"/api/v2/scenarios", "POST", helper.KeyModels{"scenario": newScenario2})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// authenticate as normal userA
	token, err = helper.AuthenticateForTest(router, helper.UserACredentials)
	assert.NoError(t, err)

	// test POST scenarios/ $newScenario1
	code, resp, err = helper.TestEndpoint(router, token,
		"/api/v2/scenarios", "POST", helper.KeyModels{"scenario": newScenario1})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// authenticate as normal userB
	token, err = helper.AuthenticateForTest(router, helper.UserBCredentials)
	assert.NoError(t, err)

	// get the length of the GET all scenarios response again
	finalNumber, err := helper.LengthOfResponse(router, token,
		"/api/v2/scenarios", "GET", nil)
	assert.NoError(t, err)

	assert.Equal(t, finalNumber, initialNumber+1)
}

func TestDeleteScenario(t *testing.T) {

	database.DropTables()
	database.MigrateModels()
	assert.NoError(t, helper.AddTestUsers())

	// authenticate as normal user
	token, err := helper.AuthenticateForTest(router, helper.UserACredentials)
	assert.NoError(t, err)

	// test POST scenarios/ $newScenario
	code, resp, err := helper.TestEndpoint(router, token,
		"/api/v2/scenarios", "POST", helper.KeyModels{"scenario": newScenario1})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Read newScenario's ID from the response
	newScenarioID, err := helper.GetResponseID(resp)
	assert.NoError(t, err)

	// add guest user to new scenario
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/scenarios/%v/user?username=User_C", newScenarioID), "PUT", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// authenticate as guest user
	token, err = helper.AuthenticateForTest(router, helper.GuestCredentials)
	assert.NoError(t, err)

	// try to delete scenario as guest
	// should return an unprocessable entity error
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/scenarios/%v", newScenarioID), "DELETE", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

	// authenticate as normal user
	token, err = helper.AuthenticateForTest(router, helper.UserACredentials)
	assert.NoError(t, err)

	// Count the number of all the scenarios returned for userA
	initialNumber, err := helper.LengthOfResponse(router, token,
		"/api/v2/scenarios", "GET", nil)
	assert.NoError(t, err)

	// Delete the added newScenario
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/scenarios/%v", newScenarioID), "DELETE", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare DELETE's response with the newScenario
	err = helper.CompareResponse(resp, helper.KeyModels{"scenario": newScenario1})
	assert.NoError(t, err)

	// Again count the number of all the scenarios returned
	finalNumber, err := helper.LengthOfResponse(router, token,
		"/api/v2/scenarios", "GET", nil)
	assert.NoError(t, err)

	assert.Equal(t, finalNumber, initialNumber-1)
}

func TestAddUserToScenario(t *testing.T) {

	database.DropTables()
	database.MigrateModels()
	assert.NoError(t, helper.AddTestUsers())

	// authenticate as normal user
	token, err := helper.AuthenticateForTest(router, helper.UserACredentials)
	assert.NoError(t, err)

	// test POST scenarios/ $newScenario
	code, resp, err := helper.TestEndpoint(router, token,
		"/api/v2/scenarios", "POST", helper.KeyModels{"scenario": newScenario1})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Read newScenario's ID from the response
	newScenarioID, err := helper.GetResponseID(resp)
	assert.NoError(t, err)

	// authenticate as normal userB who has no access to new scenario
	token, err = helper.AuthenticateForTest(router, helper.UserBCredentials)
	assert.NoError(t, err)

	// try to add new user User_C to scenario as userB
	// should return an unprocessable entity error
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/scenarios/%v/user?username=User_C", newScenarioID), "PUT", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

	// try to add new user UserB to scenario as userB
	// should return an unprocessable entity error
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/scenarios/%v/user?username=User_B", newScenarioID), "PUT", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

	// authenticate as normal user
	token, err = helper.AuthenticateForTest(router, helper.UserACredentials)
	assert.NoError(t, err)

	// Count the number of all the users returned for newScenario
	initialNumber, err := helper.LengthOfResponse(router, token,
		fmt.Sprintf("/api/v2/scenarios/%v/users", newScenarioID), "GET", nil)
	assert.NoError(t, err)
	assert.Equal(t, initialNumber, 1)

	// add userB to newScenario
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/scenarios/%v/user?username=User_B", newScenarioID), "PUT", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare resp to userB
	userB := UserRequest{
		Username: helper.UserB.Username,
		Mail:     helper.UserB.Mail,
		Role:     helper.UserB.Role,
	}
	err = helper.CompareResponse(resp, helper.KeyModels{"user": userB})
	assert.NoError(t, err)

	// Count AGAIN the number of all the users returned for newScenario
	finalNumber, err := helper.LengthOfResponse(router, token,
		fmt.Sprintf("/api/v2/scenarios/%v/users", newScenarioID), "GET", nil)
	assert.NoError(t, err)
	assert.Equal(t, finalNumber, initialNumber+1)

	// try to add a non-existing user to newScenario, should return a not found 404
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/scenarios/%v/user?username=User_D", newScenarioID), "PUT", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 404, code, "Response body: \n%v\n", resp)

}

func TestGetAllUsersOfScenario(t *testing.T) {

	database.DropTables()
	database.MigrateModels()
	assert.NoError(t, helper.AddTestUsers())

	// authenticate as normal user
	token, err := helper.AuthenticateForTest(router, helper.UserACredentials)
	assert.NoError(t, err)

	// test POST scenarios/ $newScenario
	code, resp, err := helper.TestEndpoint(router, token,
		"/api/v2/scenarios", "POST", helper.KeyModels{"scenario": newScenario1})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Read newScenario's ID from the response
	newScenarioID, err := helper.GetResponseID(resp)
	assert.NoError(t, err)

	// authenticate as normal userB who has no access to new scenario
	token, err = helper.AuthenticateForTest(router, helper.UserBCredentials)
	assert.NoError(t, err)

	// try to get all users of new scenario with userB
	// should return an unprocessable entity error
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/scenarios/%v/users", newScenarioID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

	// authenticate as normal user
	token, err = helper.AuthenticateForTest(router, helper.UserACredentials)
	assert.NoError(t, err)

	// Count the number of all the users returned for newScenario
	initialNumber, err := helper.LengthOfResponse(router, token,
		fmt.Sprintf("/api/v2/scenarios/%v/users", newScenarioID), "GET", nil)
	assert.NoError(t, err)
	assert.Equal(t, initialNumber, 1)

	// add userB to newScenario
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/scenarios/%v/user?username=User_B", newScenarioID), "PUT", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Count AGAIN the number of all the users returned for newScenario
	finalNumber, err := helper.LengthOfResponse(router, token,
		fmt.Sprintf("/api/v2/scenarios/%v/users", newScenarioID), "GET", nil)
	assert.NoError(t, err)
	assert.Equal(t, finalNumber, initialNumber+1)

	// authenticate as admin
	token, err = helper.AuthenticateForTest(router, helper.AdminCredentials)
	assert.NoError(t, err)

	// set userB as inactive
	modUserB := UserRequest{Active: "no"}
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/users/%v", 3), "PUT", helper.KeyModels{"user": modUserB})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// authenticate as normal user
	token, err = helper.AuthenticateForTest(router, helper.UserACredentials)
	assert.NoError(t, err)

	// Count AGAIN the number of all the users returned for newScenario
	finalNumber2, err := helper.LengthOfResponse(router, token,
		fmt.Sprintf("/api/v2/scenarios/%v/users", newScenarioID), "GET", nil)
	assert.NoError(t, err)
	assert.Equal(t, finalNumber2, initialNumber)
}

func TestRemoveUserFromScenario(t *testing.T) {

	database.DropTables()
	database.MigrateModels()
	assert.NoError(t, helper.AddTestUsers())

	// authenticate as normal user
	token, err := helper.AuthenticateForTest(router, helper.UserACredentials)
	assert.NoError(t, err)

	// test POST scenarios/ $newScenario
	code, resp, err := helper.TestEndpoint(router, token,
		"/api/v2/scenarios", "POST", helper.KeyModels{"scenario": newScenario1})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Read newScenario's ID from the response
	newScenarioID, err := helper.GetResponseID(resp)
	assert.NoError(t, err)

	// add userC to newScenario
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/scenarios/%v/user?username=User_C", newScenarioID), "PUT", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// authenticate as normal userB who has no access to new scenario
	token, err = helper.AuthenticateForTest(router, helper.UserBCredentials)
	assert.NoError(t, err)

	// try to delete userC from new scenario
	// should return an unprocessable entity error
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/scenarios/%v/user?username=User_C", newScenarioID), "DELETE", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

	// authenticate as normal user
	token, err = helper.AuthenticateForTest(router, helper.UserACredentials)
	assert.NoError(t, err)

	// Count the number of all the users returned for newScenario
	initialNumber, err := helper.LengthOfResponse(router, token,
		fmt.Sprintf("/api/v2/scenarios/%v/users", newScenarioID), "GET", nil)
	assert.NoError(t, err)
	assert.Equal(t, 2, initialNumber)

	// remove userC from newScenario
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/scenarios/%v/user?username=User_C", newScenarioID), "DELETE", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare DELETE's response with UserC's data
	userC := UserRequest{
		Username: helper.UserC.Username,
		Mail:     helper.UserC.Mail,
		Role:     helper.UserC.Role,
	}
	err = helper.CompareResponse(resp, helper.KeyModels{"user": userC})
	assert.NoError(t, err)

	// Count AGAIN the number of all the users returned for newScenario
	finalNumber, err := helper.LengthOfResponse(router, token,
		fmt.Sprintf("/api/v2/scenarios/%v/users", newScenarioID), "GET", nil)
	assert.NoError(t, err)
	assert.Equal(t, initialNumber-1, finalNumber)

	// Try to remove userA from new scenario
	// This should fail since User_A is the last user of newScenario
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/scenarios/%v/user?username=User_A", newScenarioID), "DELETE", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 500, code, "Response body: \n%v\n", resp)

	// Try to remove a user that does not exist in DB
	// This should fail with not found 404 status code
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/scenarios/%v/user?username=User_D", newScenarioID), "DELETE", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 404, code, "Response body: \n%v\n", resp)

	// Try to remove an admin user that is not explicitly a user of the scenario
	// This should fail with not found 404 status code
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/scenarios/%v/user?username=User_0", newScenarioID), "DELETE", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 404, code, "Response body: \n%v\n", resp)

}
