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

package component_configuration

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"git.rwth-aachen.de/acs/public/villas/web-backend-go/configuration"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/helper"
	infrastructure_component "git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/infrastructure-component"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/scenario"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/user"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/stretchr/testify/assert"
)

var router *gin.Engine
var baseAPIConfigs = "/api/v2/configs"

type ConfigRequest struct {
	Name            string         `json:"name,omitempty"`
	ScenarioID      uint           `json:"scenarioID,omitempty"`
	ICID            uint           `json:"icID,omitempty"`
	StartParameters postgres.Jsonb `json:"startParameters,omitempty"`
	FileIDs         []int64        `json:"fileIDs,omitempty"`
}

type ICRequest struct {
	UUID                 string         `json:"uuid,omitempty"`
	WebsocketURL         string         `json:"websocketurl,omitempty"`
	Type                 string         `json:"type,omitempty"`
	Name                 string         `json:"name,omitempty"`
	Category             string         `json:"category,omitempty"`
	State                string         `json:"state,omitempty"`
	Location             string         `json:"location,omitempty"`
	Description          string         `json:"description,omitempty"`
	StartParameterSchema postgres.Jsonb `json:"startparameterschema,omitempty"`
	ManagedExternally    *bool          `json:"managedexternally,omitempty"`
}

type ScenarioRequest struct {
	Name            string         `json:"name,omitempty"`
	StartParameters postgres.Jsonb `json:"startParameters,omitempty"`
}

var newConfig1 = ConfigRequest{
	Name:            "Example for Signal generator",
	StartParameters: postgres.Jsonb{RawMessage: json.RawMessage(`{"parameter1" : "testValue1A", "parameter2" : "testValue2A", "parameter3" : 42}`)},
	FileIDs:         []int64{},
}

func newFalse() *bool {
	b := false
	return &b
}

func addScenarioAndIC() (scenarioID uint, ICID uint) {

	// authenticate as admin
	token, _ := helper.AuthenticateForTest(router, database.AdminCredentials)

	// POST $newICA
	newICA := ICRequest{
		UUID:         "7be0322d-354e-431e-84bd-ae4c9633138b",
		WebsocketURL: "https://villas.k8s.eonerc.rwth-aachen.de/ws/ws_sig",
		Type:         "villas-node",
		Name:         "ACS Demo Signals",
		Category:     "gateway",
		State:        "idle",
		Location:     "k8s",
		Description:  "A signal generator for testing purposes",
		StartParameterSchema: postgres.Jsonb{
			RawMessage: json.RawMessage(`{"prop1" : "a nice prop"}`),
		},
		ManagedExternally: newFalse(),
	}

	code, resp, err := helper.TestEndpoint(router, token,
		"/api/v2/ic", "POST", helper.KeyModels{"ic": newICA})
	if code != 200 || err != nil {
		fmt.Println("Adding IC returned code", code, err, resp)
	}

	// Read newIC's ID from the response
	newICID, _ := helper.GetResponseID(resp)

	// POST a second IC to change to that IC during testing
	code, resp, err = helper.TestEndpoint(router, token,
		"/api/v2/ic", "POST", helper.KeyModels{"ic": newICA})
	if code != 200 || err != nil {
		fmt.Println("Adding IC returned code", code, err, resp)
	}

	// authenticate as normal user
	token, _ = helper.AuthenticateForTest(router, database.UserACredentials)

	// POST $newScenario
	newScenario := ScenarioRequest{
		Name: "Scenario1",
		StartParameters: postgres.Jsonb{
			RawMessage: json.RawMessage(`{"parameter1" : "testValue1A", "parameter2" : "testValue2A", "parameter3" : 42}`),
		},
	}
	code, resp, err = helper.TestEndpoint(router, token,
		"/api/v2/scenarios", "POST", helper.KeyModels{"scenario": newScenario})
	if code != 200 || err != nil {
		fmt.Println("Adding Scenario returned code", code, err, resp)
	}

	// Read newScenario's ID from the response
	newScenarioID, _ := helper.GetResponseID(resp)

	// add the guest user to the new scenario
	_, _, _ = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/scenarios/%v/user?username=User_C", newScenarioID), "PUT", nil)

	return uint(newScenarioID), uint(newICID)
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
	RegisterComponentConfigurationEndpoints(api.Group("/configs"))
	// scenario endpoints required here to first add a scenario to the DB
	// that can be associated with a new component configuration
	scenario.RegisterScenarioEndpoints(api.Group("/scenarios"))
	// IC endpoints required here to first add a IC to the DB
	// that can be associated with a new component configuration
	infrastructure_component.RegisterICEndpoints(api.Group("/ic"))

	os.Exit(m.Run())
}

func TestAddConfig(t *testing.T) {
	database.DropTables()
	database.MigrateModels()
	assert.NoError(t, database.AddTestUsers())

	// prepare the content of the DB for testing
	// by adding a scenario and a IC to the DB
	// using the respective endpoints of the API
	scenarioID, ICID := addScenarioAndIC()

	newConfig1.ScenarioID = scenarioID
	// authenticate as normal userB who has no access to new scenario
	token, err := helper.AuthenticateForTest(router, database.UserBCredentials)
	assert.NoError(t, err)

	// try to POST with no access
	// should result in unprocessable entity
	code, resp, err := helper.TestEndpoint(router, token,
		baseAPIConfigs, "POST", helper.KeyModels{"config": newConfig1})
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

	// authenticate as normal user
	token, err = helper.AuthenticateForTest(router, database.UserACredentials)
	assert.NoError(t, err)

	// try to POST non JSON body
	code, resp, err = helper.TestEndpoint(router, token,
		baseAPIConfigs, "POST", "this is not JSON")
	assert.NoError(t, err)
	assert.Equalf(t, 400, code, "Response body: \n%v\n", resp)

	// authenticate as normal user
	token, err = helper.AuthenticateForTest(router, database.UserACredentials)
	assert.NoError(t, err)

	// test POST newConfig without IC ID specified
	// should work
	code, resp, err = helper.TestEndpoint(router, token,
		baseAPIConfigs, "POST", helper.KeyModels{"config": newConfig1})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare POST's response with the newConfig
	err = helper.CompareResponse(resp, helper.KeyModels{"config": newConfig1})
	assert.NoError(t, err)

	// test POST newConfig WITH IC ID specified
	// should work
	newConfig1.ICID = ICID
	code, resp, err = helper.TestEndpoint(router, token,
		baseAPIConfigs, "POST", helper.KeyModels{"config": newConfig1})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare POST's response with the newConfig
	err = helper.CompareResponse(resp, helper.KeyModels{"config": newConfig1})
	assert.NoError(t, err)

	// Read newConfig's ID from the response
	newConfigID, err := helper.GetResponseID(resp)
	assert.NoError(t, err)

	// Get the newConfig
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("%v/%v", baseAPIConfigs, newConfigID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare GET's response with the newConfig
	err = helper.CompareResponse(resp, helper.KeyModels{"config": newConfig1})
	assert.NoError(t, err)

	// try to POST a malformed component config
	// Required fields are missing
	malformedNewConfig := ConfigRequest{
		Name: "ThisIsAMalformedRequest",
	}
	// this should NOT work and return a unprocessable entity 442 status code
	code, resp, err = helper.TestEndpoint(router, token,
		baseAPIConfigs, "POST", helper.KeyModels{"config": malformedNewConfig})
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

	// authenticate as normal userB who has no access to new scenario
	token, err = helper.AuthenticateForTest(router, database.UserBCredentials)
	assert.NoError(t, err)

	// Try to GET the newConfig with no access
	// Should result in unprocessable entity
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("%v/%v", baseAPIConfigs, newConfigID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

}

func TestUpdateConfig(t *testing.T) {

	database.DropTables()
	database.MigrateModels()
	assert.NoError(t, database.AddTestUsers())

	// prepare the content of the DB for testing
	// by adding a scenario and a IC to the DB
	// using the respective endpoints of the API
	scenarioID, ICID := addScenarioAndIC()

	// authenticate as normal user
	token, err := helper.AuthenticateForTest(router, database.UserACredentials)
	assert.NoError(t, err)

	// test POST newConfig
	newConfig1.ScenarioID = scenarioID
	newConfig1.ICID = ICID
	code, resp, err := helper.TestEndpoint(router, token,
		baseAPIConfigs, "POST", helper.KeyModels{"config": newConfig1})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Read newConfig's ID from the response
	newConfigID, err := helper.GetResponseID(resp)
	assert.NoError(t, err)

	updatedConfig := ConfigRequest{
		Name:            "VILLASnode gateway X",
		StartParameters: postgres.Jsonb{RawMessage: json.RawMessage(`{"parameter1" : "testValue1B", "parameter2" : "testValue2B", "parameter3" : 43}`)},
	}

	// authenticate as normal userB who has no access to new scenario
	token, err = helper.AuthenticateForTest(router, database.UserBCredentials)
	assert.NoError(t, err)

	// try to PUT with no access
	// should result in unprocessable entity
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("%v/%v", baseAPIConfigs, newConfigID), "PUT", helper.KeyModels{"config": updatedConfig})
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

	// authenticate as guest user who has access to component config
	token, err = helper.AuthenticateForTest(router, database.GuestCredentials)
	assert.NoError(t, err)

	// try to PUT as guest
	// should NOT work and result in unprocessable entity
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("%v/%v", baseAPIConfigs, newConfigID), "PUT", helper.KeyModels{"config": updatedConfig})
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

	// authenticate as normal user
	token, err = helper.AuthenticateForTest(router, database.UserACredentials)
	assert.NoError(t, err)

	// try to PUT a non JSON body
	// should result in a bad request
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("%v/%v", baseAPIConfigs, newConfigID), "PUT", "This is not JSON")
	assert.NoError(t, err)
	assert.Equalf(t, 400, code, "Response body: \n%v\n", resp)

	// test PUT
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("%v/%v", baseAPIConfigs, newConfigID), "PUT", helper.KeyModels{"config": updatedConfig})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare PUT's response with the updateConfig
	err = helper.CompareResponse(resp, helper.KeyModels{"config": updatedConfig})
	assert.NoError(t, err)

	//Change IC ID to use second IC available in DB
	updatedConfig.ICID = ICID + 1
	// test PUT again
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("%v/%v", baseAPIConfigs, newConfigID), "PUT", helper.KeyModels{"config": updatedConfig})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare PUT's response with the updateConfig
	err = helper.CompareResponse(resp, helper.KeyModels{"config": updatedConfig})
	assert.NoError(t, err)

	// Get the updateConfig
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("%v/%v", baseAPIConfigs, newConfigID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare GET's response with the updateConfig
	err = helper.CompareResponse(resp, helper.KeyModels{"config": updatedConfig})
	assert.NoError(t, err)

	// try to update a component config that does not exist (should return not found 404 status code)
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("%v/%v", baseAPIConfigs, newConfigID+1), "PUT", helper.KeyModels{"config": updatedConfig})
	assert.NoError(t, err)
	assert.Equalf(t, 404, code, "Response body: \n%v\n", resp)
}

func TestDeleteConfig(t *testing.T) {
	database.DropTables()
	database.MigrateModels()
	assert.NoError(t, database.AddTestUsers())

	// prepare the content of the DB for testing
	// by adding a scenario and a IC to the DB
	// using the respective endpoints of the API
	scenarioID, ICID := addScenarioAndIC()
	newConfig1.ScenarioID = scenarioID
	newConfig1.ICID = ICID

	// authenticate as normal user
	token, err := helper.AuthenticateForTest(router, database.UserACredentials)
	assert.NoError(t, err)

	// test POST newConfig
	code, resp, err := helper.TestEndpoint(router, token,
		baseAPIConfigs, "POST", helper.KeyModels{"config": newConfig1})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Read newConfig's ID from the response
	newConfigID, err := helper.GetResponseID(resp)
	assert.NoError(t, err)

	// authenticate as normal userB who has no access to new scenario
	token, err = helper.AuthenticateForTest(router, database.UserBCredentials)
	assert.NoError(t, err)

	// try to DELETE with no access
	// should result in unprocessable entity
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("%v/%v", baseAPIConfigs, newConfigID), "DELETE", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

	// authenticate as normal user
	token, err = helper.AuthenticateForTest(router, database.UserACredentials)
	assert.NoError(t, err)

	// Count the number of all the component config returned for scenario
	initialNumber, err := helper.LengthOfResponse(router, token,
		fmt.Sprintf("%v?scenarioID=%v", baseAPIConfigs, scenarioID), "GET", nil)
	assert.NoError(t, err)

	// Count the number of configs associated with IC, should be 1 at this point
	initialConfigsofIC, err := helper.LengthOfResponse(router, token,
		fmt.Sprintf("/api/v2/ic/%v/configs", ICID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 1, initialConfigsofIC, "IC should have 1 associated config, has %v", initialConfigsofIC)

	// Delete the added newConfig
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("%v/%v", baseAPIConfigs, newConfigID), "DELETE", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare DELETE's response with the newConfig
	err = helper.CompareResponse(resp, helper.KeyModels{"config": newConfig1})
	assert.NoError(t, err)

	// Again count the number of all the component configs returned
	finalNumber, err := helper.LengthOfResponse(router, token,
		fmt.Sprintf("%v?scenarioID=%v", baseAPIConfigs, scenarioID), "GET", nil)
	assert.NoError(t, err)

	assert.Equal(t, initialNumber-1, finalNumber)

	// Count final the number of configs associated with IC, should be 0 at this point
	finalConfigsofIC, err := helper.LengthOfResponse(router, token,
		fmt.Sprintf("/api/v2/ic/%v/configs", ICID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 0, finalConfigsofIC, "IC should have 0 associated config, has %v", finalConfigsofIC)

	// authenticate as admin user to be able to delete IC
	token, err = helper.AuthenticateForTest(router, database.AdminCredentials)
	assert.NoError(t, err)

	// Delete IC (should work because there is no config associated to it)
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/ic/%v", ICID), "DELETE", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

}

func TestGetAllConfigsOfScenario(t *testing.T) {
	database.DropTables()
	database.MigrateModels()
	assert.NoError(t, database.AddTestUsers())

	// prepare the content of the DB for testing
	// by adding a scenario and a IC to the DB
	// using the respective endpoints of the API
	scenarioID, ICID := addScenarioAndIC()
	newConfig1.ScenarioID = scenarioID
	newConfig1.ICID = ICID

	// authenticate as normal user
	token, err := helper.AuthenticateForTest(router, database.UserACredentials)
	assert.NoError(t, err)

	// test POST newConfig
	code, resp, err := helper.TestEndpoint(router, token,
		baseAPIConfigs, "POST", helper.KeyModels{"config": newConfig1})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Count the number of all the component config returned for scenario
	NumberOfConfigs, err := helper.LengthOfResponse(router, token,
		fmt.Sprintf("%v?scenarioID=%v", baseAPIConfigs, scenarioID), "GET", nil)
	assert.NoError(t, err)

	assert.Equal(t, 1, NumberOfConfigs)

	// authenticate as normal userB who has no access to scenario
	token, err = helper.AuthenticateForTest(router, database.UserBCredentials)
	assert.NoError(t, err)

	// try to get configs without access
	// should result in unprocessable entity
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("%v?scenarioID=%v", baseAPIConfigs, scenarioID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

}
