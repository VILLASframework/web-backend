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

package scenario

import (
	"bytes"
	"encoding/json"
	"fmt"
	component_configuration "git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/component-configuration"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/dashboard"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/file"
	infrastructure_component "git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/infrastructure-component"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/result"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/signal"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/widget"
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
	Name: "Scenario1",
	StartParameters: postgres.Jsonb{
		RawMessage: json.RawMessage(`{"parameter1" : "testValue1A", "parameter2" : "testValue2A", "parameter3" : 42}`),
	},
	IsLocked: false,
}

var newScenario2 = ScenarioRequest{
	Name: "Scenario2",
	StartParameters: postgres.Jsonb{
		RawMessage: json.RawMessage(`{"parameter1" : "testValue1B", "parameter2" : "testValue2B", "parameter3" : 55}`),
	},
	IsLocked: false,
}

func TestMain(m *testing.M) {
	err := configuration.InitConfig()
	if err != nil {
		panic(m)
	}

	err = database.InitDB(configuration.GlobalConfig, true)
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
	file.RegisterFileEndpoints(api.Group("/files"))
	component_configuration.RegisterComponentConfigurationEndpoints(api.Group("/configs"))
	signal.RegisterSignalEndpoints(api.Group("/signals"))
	dashboard.RegisterDashboardEndpoints(api.Group("/dashboards"))
	widget.RegisterWidgetEndpoints(api.Group("/widgets"))
	result.RegisterResultEndpoints(api.Group("/results"))
	infrastructure_component.RegisterICEndpoints(api.Group("/ic"))
	RegisterScenarioEndpoints(api.Group("/scenarios"))

	os.Exit(m.Run())
}

func TestAddScenario(t *testing.T) {

	database.DropTables()
	database.MigrateModels()
	assert.NoError(t, database.AddTestUsers())

	// authenticate as normal user
	token, err := helper.AuthenticateForTest(router, database.UserACredentials)
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
	token, err = helper.AuthenticateForTest(router, database.GuestCredentials)
	assert.NoError(t, err)

	// try to add scenario as guest user
	// should return an unprocessable entity error
	code, resp, err = helper.TestEndpoint(router, token,
		"/api/v2/scenarios", "POST", helper.KeyModels{"scenario": newScenario1})
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

	// authenticate as userB who has no access to the added scenario
	token, err = helper.AuthenticateForTest(router, database.UserBCredentials)
	assert.NoError(t, err)

	// try to GET a scenario to which user B has no access
	// should return an unprocessable entity error
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/scenarios/%v", newScenarioID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

	// authenticate as admin user who has no access to everything
	token, err = helper.AuthenticateForTest(router, database.AdminCredentials)
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
	assert.NoError(t, database.AddTestUsers())

	// authenticate as normal user
	token, err := helper.AuthenticateForTest(router, database.UserACredentials)
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
	token, err = helper.AuthenticateForTest(router, database.AdminCredentials)
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
	token, err = helper.AuthenticateForTest(router, database.UserACredentials)
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
	assert.NoError(t, database.AddTestUsers())

	// authenticate as admin
	token, err := helper.AuthenticateForTest(router, database.AdminCredentials)
	assert.NoError(t, err)

	// get the length of the GET all scenarios response for admin
	initialNumber, err := helper.LengthOfResponse(router, token,
		"/api/v2/scenarios", "GET", nil)
	assert.NoError(t, err)

	// authenticate as normal userB
	token, err = helper.AuthenticateForTest(router, database.UserBCredentials)
	assert.NoError(t, err)

	// test POST scenarios/ $newScenario1
	code, resp, err := helper.TestEndpoint(router, token,
		"/api/v2/scenarios", "POST", helper.KeyModels{"scenario": newScenario1})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// authenticate as normal userA
	token, err = helper.AuthenticateForTest(router, database.UserACredentials)
	assert.NoError(t, err)

	// test POST scenarios/ $newScenario2
	code, resp, err = helper.TestEndpoint(router, token,
		"/api/v2/scenarios", "POST", helper.KeyModels{"scenario": newScenario2})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// authenticate as admin
	token, err = helper.AuthenticateForTest(router, database.AdminCredentials)
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
	assert.NoError(t, database.AddTestUsers())

	// authenticate as normal userB
	token, err := helper.AuthenticateForTest(router, database.UserBCredentials)
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
	token, err = helper.AuthenticateForTest(router, database.UserACredentials)
	assert.NoError(t, err)

	// test POST scenarios/ $newScenario1
	code, resp, err = helper.TestEndpoint(router, token,
		"/api/v2/scenarios", "POST", helper.KeyModels{"scenario": newScenario1})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// authenticate as normal userB
	token, err = helper.AuthenticateForTest(router, database.UserBCredentials)
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
	assert.NoError(t, database.AddTestUsers())

	// authenticate as admin user to add ICs
	token, err := helper.AuthenticateForTest(router, database.AdminCredentials)
	assert.NoError(t, err)

	ic1ID, ic2ID := addICs(t, token)

	// authenticate as normal user
	token, err = helper.AuthenticateForTest(router, database.UserACredentials)
	assert.NoError(t, err)

	// test POST scenarios/ $newScenario
	code, resp, err := helper.TestEndpoint(router, token,
		"/api/v2/scenarios", "POST", helper.KeyModels{"scenario": newScenario1})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Read newScenario's ID from the response
	newScenarioID, err := helper.GetResponseID(resp)
	assert.NoError(t, err)

	// add file to the scenario
	fileID := addFile(t, token, newScenarioID)

	// add result to the scenario
	resultID := addResult(t, token, newScenarioID)

	// add dashboard to the scenario
	dashboardID := addDashboard(t, token, newScenarioID)

	// add widget to the dashboard
	widgetID := addWidget(t, token, dashboardID)

	// add component config to the scenario
	componentConfig1ID, componentConfig2ID := addComponentConfigs(t, token, newScenarioID, ic1ID, ic2ID)

	// add signal to the component config
	signalInID, signalOutID := addSignals(t, token, componentConfig1ID)

	// add guest user to new scenario
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/scenarios/%v/user?username=User_C", newScenarioID), "PUT", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// authenticate as guest user
	token, err = helper.AuthenticateForTest(router, database.GuestCredentials)
	assert.NoError(t, err)

	// try to delete scenario as guest
	// should return an unprocessable entity error
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/scenarios/%v", newScenarioID), "DELETE", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

	// authenticate as normal user
	token, err = helper.AuthenticateForTest(router, database.UserACredentials)
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

	// check if dashboard, result, file, etc. still exists
	// make sure everything is properly deleted

	// Get the file
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/files/%v", fileID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 404, code, "Response body: \n%v\n", resp)

	// Get the result
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/results/%v", resultID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 404, code, "Response body: \n%v\n", resp)

	// Get the dashboard
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/dashboards/%v", dashboardID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 404, code, "Response body: \n%v\n", resp)

	// Get the widget
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/widgets/%v", widgetID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 404, code, "Response body: \n%v\n", resp)

	// Get the configs
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/configs/%v", componentConfig1ID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 404, code, "Response body: \n%v\n", resp)
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/configs/%v", componentConfig2ID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 404, code, "Response body: \n%v\n", resp)

	// Get the signals
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/signals/%v", signalInID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 404, code, "Response body: \n%v\n", resp)
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/signals/%v", signalOutID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 404, code, "Response body: \n%v\n", resp)

	// Get IC1 (should be in DB)
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/ic/%v", ic1ID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Get number of configs of IC1 (should be zero)
	numberOfConfigs, err := helper.LengthOfResponse(router, token,
		fmt.Sprintf("/api/v2/ic/%v/configs", ic1ID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)
	assert.Equal(t, 0, numberOfConfigs)

	// Get IC2 (should be deleted)
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/ic/%v", ic2ID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 404, code, "Response body: \n%v\n", resp)
	assert.Equalf(t, 404, code, "Response body: \n%v\n", resp)
}

func TestAddUserToScenario(t *testing.T) {

	database.DropTables()
	database.MigrateModels()
	assert.NoError(t, database.AddTestUsers())

	// authenticate as normal user
	token, err := helper.AuthenticateForTest(router, database.UserACredentials)
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
	token, err = helper.AuthenticateForTest(router, database.UserBCredentials)
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
	token, err = helper.AuthenticateForTest(router, database.UserACredentials)
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
		Username: database.UserB.Username,
		Mail:     database.UserB.Mail,
		Role:     database.UserB.Role,
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
	assert.NoError(t, database.AddTestUsers())

	// authenticate as normal user
	token, err := helper.AuthenticateForTest(router, database.UserACredentials)
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
	token, err = helper.AuthenticateForTest(router, database.UserBCredentials)
	assert.NoError(t, err)

	// try to get all users of new scenario with userB
	// should return an unprocessable entity error
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/scenarios/%v/users", newScenarioID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

	// authenticate as normal user
	token, err = helper.AuthenticateForTest(router, database.UserACredentials)
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
	token, err = helper.AuthenticateForTest(router, database.AdminCredentials)
	assert.NoError(t, err)

	// set userB as inactive
	modUserB := UserRequest{Active: "no"}
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/users/%v", 3), "PUT", helper.KeyModels{"user": modUserB})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// authenticate as normal user
	token, err = helper.AuthenticateForTest(router, database.UserACredentials)
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
	assert.NoError(t, database.AddTestUsers())

	// authenticate as normal user
	token, err := helper.AuthenticateForTest(router, database.UserACredentials)
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
	token, err = helper.AuthenticateForTest(router, database.UserBCredentials)
	assert.NoError(t, err)

	// try to delete userC from new scenario
	// should return an unprocessable entity error
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/scenarios/%v/user?username=User_C", newScenarioID), "DELETE", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

	// authenticate as normal user
	token, err = helper.AuthenticateForTest(router, database.UserACredentials)
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
		Username: database.UserC.Username,
		Mail:     database.UserC.Mail,
		Role:     database.UserC.Role,
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

func addFile(t *testing.T, token string, scenarioID int) int {
	// create a testfile.txt in local folder
	c1 := []byte("This is my testfile\n")
	err := ioutil.WriteFile("testfile.txt", c1, 0644)
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

	return newFileID
}

func addResult(t *testing.T, token string, scenarioID int) int {

	type ResultRequest struct {
		Description     string         `json:"description,omitempty"`
		ScenarioID      uint           `json:"scenarioID,omitempty"`
		ConfigSnapshots postgres.Jsonb `json:"configSnapshots,omitempty"`
	}

	var newResult = ResultRequest{
		Description: "This is a test result.",
	}

	configSnapshot1 := json.RawMessage(`{"configs": [ {"Name" : "conf1", "scenarioID" : 1}, {"Name" : "conf2", "scenarioID" : 1}]}`)
	confSnapshots := postgres.Jsonb{
		RawMessage: configSnapshot1,
	}
	newResult.ScenarioID = uint(scenarioID)
	newResult.ConfigSnapshots = confSnapshots

	code, resp, err := helper.TestEndpoint(router, token,
		"/api/v2/results", "POST", helper.KeyModels{"result": newResult})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Read newResults's ID from the response
	newResultID, err := helper.GetResponseID(resp)
	assert.NoError(t, err)

	return newResultID
}

func addDashboard(t *testing.T, token string, scenarioID int) int {

	type DashboardRequest struct {
		Name       string `json:"name,omitempty"`
		Grid       int    `json:"grid,omitempty"`
		Height     int    `json:"height,omitempty"`
		ScenarioID uint   `json:"scenarioID,omitempty"`
	}

	var newDashboard = DashboardRequest{
		Name: "Dashboard_A",
		Grid: 15,
	}

	// test POST dashboards/ $newDashboad
	newDashboard.ScenarioID = uint(scenarioID)
	code, resp, err := helper.TestEndpoint(router, token,
		"/api/v2/dashboards", "POST", helper.KeyModels{"dashboard": newDashboard})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Read newDashboard's ID from the response
	newDashboardID, err := helper.GetResponseID(resp)
	assert.NoError(t, err)

	return newDashboardID
}

func addWidget(t *testing.T, token string, dashboardID int) int {

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
		SignalIDs        []int64        `json:"signalIDs,omitempty"`
	}

	var newWidget = WidgetRequest{
		Name:             "My label",
		Type:             "Label",
		Width:            100,
		Height:           50,
		MinWidth:         40,
		MinHeight:        80,
		X:                10,
		Y:                10,
		Z:                200,
		IsLocked:         false,
		CustomProperties: postgres.Jsonb{RawMessage: json.RawMessage(`{"textSize" : "20", "fontColor" : "#4287f5", "fontColor_opacity": 1}`)},
		SignalIDs:        []int64{},
	}

	newWidget.DashboardID = uint(dashboardID)

	code, resp, err := helper.TestEndpoint(router, token,
		"/api/v2/widgets", "POST", helper.KeyModels{"widget": newWidget})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Read newWidget's ID from the response
	newWidgetID, err := helper.GetResponseID(resp)
	assert.NoError(t, err)

	return newWidgetID
}

func addComponentConfigs(t *testing.T, token string, scenarioID int, ic1ID int, ic2ID int) (int, int) {
	type ConfigRequest struct {
		Name            string         `json:"name,omitempty"`
		ScenarioID      uint           `json:"scenarioID,omitempty"`
		ICID            uint           `json:"icID,omitempty"`
		StartParameters postgres.Jsonb `json:"startParameters,omitempty"`
		FileIDs         []int64        `json:"fileIDs,omitempty"`
	}

	var newConfig1 = ConfigRequest{
		Name:            "Example for Signal generator",
		StartParameters: postgres.Jsonb{RawMessage: json.RawMessage(`{"parameter1" : "testValue1A", "parameter2" : "testValue2A", "parameter3" : 42}`)},
		FileIDs:         []int64{},
	}

	var newConfig2 = ConfigRequest{
		Name:            "Config example 2",
		StartParameters: postgres.Jsonb{RawMessage: json.RawMessage(`{"parameter1" : "testValue1A"}`)},
		FileIDs:         []int64{},
	}

	newConfig1.ScenarioID = uint(scenarioID)
	newConfig1.ICID = uint(ic1ID)
	newConfig2.ScenarioID = uint(scenarioID)
	newConfig2.ICID = uint(ic2ID)

	code, resp, err := helper.TestEndpoint(router, token,
		"/api/v2/configs", "POST", helper.KeyModels{"config": newConfig1})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Read newConfig's ID from the response
	newConfig1ID, err := helper.GetResponseID(resp)
	assert.NoError(t, err)

	code, resp, err = helper.TestEndpoint(router, token,
		"/api/v2/configs", "POST", helper.KeyModels{"config": newConfig2})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Read newConfig's ID from the response
	newConfig2ID, err := helper.GetResponseID(resp)
	assert.NoError(t, err)

	return newConfig1ID, newConfig2ID
}

func addSignals(t *testing.T, token string, componentConfigID int) (int, int) {

	type SignalRequest struct {
		Name          string  `json:"name,omitempty"`
		Unit          string  `json:"unit,omitempty"`
		Index         *uint   `json:"index,omitempty"`
		Direction     string  `json:"direction,omitempty"`
		ScalingFactor float32 `json:"scalingFactor,omitempty"`
		ConfigID      uint    `json:"configID,omitempty"`
	}

	var signalIndex0 uint = 0

	var newSignalOut = SignalRequest{
		Name:      "outSignal_A",
		Unit:      "V",
		Direction: "out",
		Index:     &signalIndex0,
	}

	var newSignalIn = SignalRequest{
		Name:      "inSignal_A",
		Unit:      "V",
		Direction: "in",
		Index:     &signalIndex0,
	}

	newSignalOut.ConfigID = uint(componentConfigID)
	newSignalIn.ConfigID = uint(componentConfigID)

	code, resp, err := helper.TestEndpoint(router, token,
		"/api/v2/signals", "POST", helper.KeyModels{"signal": newSignalOut})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Read newSignal's ID from the response
	newSignalIDOut, err := helper.GetResponseID(resp)
	assert.NoError(t, err)

	code, resp, err = helper.TestEndpoint(router, token,
		"/api/v2/signals", "POST", helper.KeyModels{"signal": newSignalIn})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Read newSignal's ID from the response
	newSignalIDIn, err := helper.GetResponseID(resp)
	assert.NoError(t, err)

	return newSignalIDIn, newSignalIDOut

}

func addICs(t *testing.T, token string) (int, int) {

	type ICRequest struct {
		UUID                  string         `json:"uuid,omitempty"`
		WebsocketURL          string         `json:"websocketurl,omitempty"`
		APIURL                string         `json:"apiurl,omitempty"`
		Type                  string         `json:"type,omitempty"`
		Name                  string         `json:"name,omitempty"`
		Category              string         `json:"category,omitempty"`
		State                 string         `json:"state,omitempty"`
		Location              string         `json:"location,omitempty"`
		Description           string         `json:"description,omitempty"`
		StartParameterSchema  postgres.Jsonb `json:"startparameterschema,omitempty"`
		CreateParameterSchema postgres.Jsonb `json:"createparameterschema,omitempty"`
		ManagedExternally     *bool          `json:"managedexternally"`
		Manager               string         `json:"manager,omitempty"`
	}

	var newIC1 = ICRequest{
		UUID:                  "7be0322d-354e-431e-84bd-ae4c9633138b",
		WebsocketURL:          "https://villas.k8s.eonerc.rwth-aachen.de/ws/ws_sig",
		APIURL:                "https://villas.k8s.eonerc.rwth-aachen.de/ws/api/v2",
		Type:                  "villas-node",
		Name:                  "ACS Demo Signals",
		Category:              "gateway",
		State:                 "idle",
		Location:              "k8s",
		Description:           "A signal generator for testing purposes",
		StartParameterSchema:  postgres.Jsonb{RawMessage: json.RawMessage(`{"startprop1" : "a nice prop"}`)},
		CreateParameterSchema: postgres.Jsonb{RawMessage: json.RawMessage(`{"createprop1" : "a really nice prop"}`)},
		ManagedExternally:     newFalse(),
		Manager:               "7be0322d-354e-431e-84bd-ae4c9633beef",
	}

	// IC with state gone
	var newIC2 = ICRequest{
		UUID:                  "4854af30-325f-44a5-ad59-b67b2597de68",
		WebsocketURL:          "xxx.yyy.zzz.aaa",
		APIURL:                "https://villas.k8s.eonerc.rwth-aachen.de/ws/api/v2",
		Type:                  "dpsim",
		Name:                  "Test DPsim Simulator",
		Category:              "simulator",
		State:                 "gone",
		Location:              "k8s",
		Description:           "This is a test description",
		StartParameterSchema:  postgres.Jsonb{RawMessage: json.RawMessage(`{"startprop1" : "a nice prop"}`)},
		CreateParameterSchema: postgres.Jsonb{RawMessage: json.RawMessage(`{"createprop1" : "a really nice prop"}`)},
		ManagedExternally:     newFalse(),
		Manager:               "4854af30-325f-44a5-ad59-b67b2597de99",
	}

	code, resp, err := helper.TestEndpoint(router, token,
		"/api/v2/ic", "POST", helper.KeyModels{"ic": newIC1})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Read newIC's ID from the response
	newIC1ID, err := helper.GetResponseID(resp)
	assert.NoError(t, err)

	code, resp, err = helper.TestEndpoint(router, token,
		"/api/v2/ic", "POST", helper.KeyModels{"ic": newIC2})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Read newIC's ID from the response
	newIC2ID, err := helper.GetResponseID(resp)
	assert.NoError(t, err)

	return newIC1ID, newIC2ID
}

func newFalse() *bool {
	b := false
	return &b
}
