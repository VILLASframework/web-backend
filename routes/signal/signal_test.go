/** Signal package, testing.
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
package signal

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"git.rwth-aachen.de/acs/public/villas/web-backend-go/configuration"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/helper"
	component_configuration "git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/component-configuration"
	infrastructure_component "git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/infrastructure-component"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/scenario"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/user"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/stretchr/testify/assert"
)

var router *gin.Engine

type SignalRequest struct {
	Name          string  `json:"name,omitempty"`
	Unit          string  `json:"unit,omitempty"`
	Index         uint    `json:"index,omitempty"`
	Direction     string  `json:"direction,omitempty"`
	ScalingFactor float32 `json:"scalingFactor,omitempty"`
	ConfigID      uint    `json:"configID,omitempty"`
}

type ConfigRequest struct {
	Name            string         `json:"name,omitempty"`
	ScenarioID      uint           `json:"scenarioID,omitempty"`
	ICID            uint           `json:"icID,omitempty"`
	StartParameters postgres.Jsonb `json:"startParameters,omitempty"`
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
	Running         bool           `json:"running,omitempty"`
	StartParameters postgres.Jsonb `json:"startParameters,omitempty"`
}

var newSignal1 = SignalRequest{
	Name:      "outSignal_A",
	Unit:      "V",
	Direction: "out",
	Index:     1,
}

func newFalse() *bool {
	b := false
	return &b
}

func addScenarioAndICAndConfig() (scenarioID uint, ICID uint, configID uint) {

	// authenticate as admin
	token, _ := helper.AuthenticateForTest(router, helper.AdminCredentials)

	// POST $newICA
	newICA := ICRequest{
		UUID:                 "7be0322d-354e-431e-84bd-ae4c9633138b",
		WebsocketURL:         "https://villas.k8s.eonerc.rwth-aachen.de/ws/ws_sig",
		Type:                 "villas-node",
		Name:                 "ACS Demo Signals",
		Category:             "gateway",
		State:                "idle",
		Location:             "k8s",
		Description:          "A signal generator for testing purposes",
		StartParameterSchema: postgres.Jsonb{json.RawMessage(`{"prop1" : "a nice prop"}`)},
		ManagedExternally:    newFalse(),
	}
	_, resp, _ := helper.TestEndpoint(router, token,
		"/api/v2/ic", "POST", helper.KeyModels{"ic": newICA})

	// Read newIC's ID from the response
	newICID, _ := helper.GetResponseID(resp)

	// authenticate as normal user
	token, _ = helper.AuthenticateForTest(router, helper.UserACredentials)

	// POST $newScenario
	newScenario := ScenarioRequest{
		Name:            "Scenario1",
		Running:         true,
		StartParameters: postgres.Jsonb{json.RawMessage(`{"parameter1" : "testValue1A", "parameter2" : "testValue2A", "parameter3" : 42}`)},
	}
	_, resp, _ = helper.TestEndpoint(router, token,
		"/api/v2/scenarios", "POST", helper.KeyModels{"scenario": newScenario})

	// Read newScenario's ID from the response
	newScenarioID, _ := helper.GetResponseID(resp)

	// test POST newConfig
	newConfig := ConfigRequest{
		Name:            "Example for Signal generator",
		StartParameters: postgres.Jsonb{RawMessage: json.RawMessage(`{"parameter1" : "testValue1A", "parameter2" : "testValue2A", "parameter3" : 42}`)},
		ScenarioID:      uint(newScenarioID),
		ICID:            uint(newICID),
	}
	_, resp, _ = helper.TestEndpoint(router, token,
		"/api/v2/configs", "POST", helper.KeyModels{"config": newConfig})

	// Read newConfig's ID from the response
	newConfigID, _ := helper.GetResponseID(resp)

	// add the guest user to the new scenario
	_, resp, _ = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/scenarios/%v/user?username=User_C", newScenarioID), "PUT", nil)

	return uint(newScenarioID), uint(newICID), uint(newConfigID)
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
	// component-configuration endpoints required here to first add a component config to the DB
	// that can be associated with a new signal
	component_configuration.RegisterComponentConfigurationEndpoints(api.Group("/configs"))
	// scenario endpoints required here to first add a scenario to the DB
	// that can be associated with a new component config
	scenario.RegisterScenarioEndpoints(api.Group("/scenarios"))
	// IC endpoints required here to first add a IC to the DB
	// that can be associated with a new component config
	infrastructure_component.RegisterICEndpoints(api.Group("/ic"))
	RegisterSignalEndpoints(api.Group("/signals"))

	os.Exit(m.Run())
}

func TestAddSignal(t *testing.T) {
	database.DropTables()
	database.MigrateModels()
	assert.NoError(t, helper.AddTestUsers())

	// prepare the content of the DB for testing
	// by adding a scenario and a IC to the DB
	// using the respective endpoints of the API
	_, _, configID := addScenarioAndICAndConfig()

	// authenticate as normal user
	token, err := helper.AuthenticateForTest(router, helper.UserACredentials)
	assert.NoError(t, err)

	newSignal1.ConfigID = configID
	// authenticate as normal userB who has no access to new scenario
	token, err = helper.AuthenticateForTest(router, helper.UserBCredentials)
	assert.NoError(t, err)

	// try to POST to component config without access
	// should result in unprocessable entity
	code, resp, err := helper.TestEndpoint(router, token,
		"/api/v2/signals", "POST", helper.KeyModels{"signal": newSignal1})
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

	// authenticate as normal user
	token, err = helper.AuthenticateForTest(router, helper.UserACredentials)
	assert.NoError(t, err)

	// try to POST a signal with non JSON body
	// should result in a bad request
	code, resp, err = helper.TestEndpoint(router, token,
		"/api/v2/signals", "POST", "this is not a JSON")
	assert.NoError(t, err)
	assert.Equalf(t, 400, code, "Response body: \n%v\n", resp)

	// test POST signals/ $newSignal
	code, resp, err = helper.TestEndpoint(router, token,
		"/api/v2/signals", "POST", helper.KeyModels{"signal": newSignal1})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare POST's response with the newSignal
	err = helper.CompareResponse(resp, helper.KeyModels{"signal": newSignal1})
	assert.NoError(t, err)

	// Read newSignal's ID from the response
	newSignalID, err := helper.GetResponseID(resp)
	assert.NoError(t, err)

	// Get the newSignal
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/signals/%v", newSignalID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare GET's response with the newSSignal
	err = helper.CompareResponse(resp, helper.KeyModels{"signal": newSignal1})
	assert.NoError(t, err)

	// try to POST a malformed signal
	// Required fields are missing
	malformedNewSignal := SignalRequest{
		Name: "ThisIsAMalformedRequest",
	}
	// this should NOT work and return a unprocessable entity 442 status code
	code, resp, err = helper.TestEndpoint(router, token,
		"/api/v2/signals", "POST", helper.KeyModels{"signal": malformedNewSignal})
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

	// authenticate as normal userB who has no access to new scenario
	token, err = helper.AuthenticateForTest(router, helper.UserBCredentials)
	assert.NoError(t, err)

	// Try to Get the newSignal as user B
	// should result in unprocessable entity
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/signals/%v", newSignalID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)
}

func TestUpdateSignal(t *testing.T) {
	database.DropTables()
	database.MigrateModels()
	assert.NoError(t, helper.AddTestUsers())

	// prepare the content of the DB for testing
	// by adding a scenario and a IC to the DB
	// using the respective endpoints of the API
	_, _, configID := addScenarioAndICAndConfig()

	// authenticate as normal user
	token, err := helper.AuthenticateForTest(router, helper.UserACredentials)
	assert.NoError(t, err)

	// test POST signals/ $newSignal
	newSignal1.ConfigID = configID
	code, resp, err := helper.TestEndpoint(router, token,
		"/api/v2/signals", "POST", helper.KeyModels{"signal": newSignal1})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Read newSignal's ID from the response
	newSignalID, err := helper.GetResponseID(resp)
	assert.NoError(t, err)

	updatedSignal := SignalRequest{
		Name:  "outSignal_B",
		Unit:  "A",
		Index: 1,
	}

	// authenticate as normal userB who has no access to new scenario
	token, err = helper.AuthenticateForTest(router, helper.UserBCredentials)
	assert.NoError(t, err)

	// try to PUT signal without access
	// should result in unprocessable entity
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/signals/%v", newSignalID), "PUT", helper.KeyModels{"signal": updatedSignal})
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

	// authenticate as guest user
	token, err = helper.AuthenticateForTest(router, helper.GuestCredentials)
	assert.NoError(t, err)

	// try to update signal as guest who has access to scenario
	// should result in unprocessable entity
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/signals/%v", newSignalID), "PUT", helper.KeyModels{"signal": updatedSignal})
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

	// authenticate as normal user
	token, err = helper.AuthenticateForTest(router, helper.UserACredentials)
	assert.NoError(t, err)

	// try to PUT with non JSON body
	// should result in bad request
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/signals/%v", newSignalID), "PUT", "This is not JSON")
	assert.NoError(t, err)
	assert.Equalf(t, 400, code, "Response body: \n%v\n", resp)

	// test PUT
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/signals/%v", newSignalID), "PUT", helper.KeyModels{"signal": updatedSignal})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare PUT's response with the updatedSignal
	err = helper.CompareResponse(resp, helper.KeyModels{"signal": updatedSignal})
	assert.NoError(t, err)

	// Get the updatedSignal
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/signals/%v", newSignalID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare GET's response with the updatedSignal
	err = helper.CompareResponse(resp, helper.KeyModels{"signal": updatedSignal})
	assert.NoError(t, err)

	// try to update a signal that does not exist (should return not found 404 status code)
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/signals/%v", newSignalID+1), "PUT", helper.KeyModels{"signal": updatedSignal})
	assert.NoError(t, err)
	assert.Equalf(t, 404, code, "Response body: \n%v\n", resp)

}

func TestDeleteSignal(t *testing.T) {
	database.DropTables()
	database.MigrateModels()
	assert.NoError(t, helper.AddTestUsers())

	// prepare the content of the DB for testing
	// by adding a scenario and a IC to the DB
	// using the respective endpoints of the API
	_, _, configID := addScenarioAndICAndConfig()

	// authenticate as normal user
	token, err := helper.AuthenticateForTest(router, helper.UserACredentials)
	assert.NoError(t, err)

	// test POST signals/ $newSignal
	newSignal := SignalRequest{
		Name:      "insignalA",
		Unit:      "s",
		Direction: "in",
		Index:     1,
		ConfigID:  configID,
	}
	code, resp, err := helper.TestEndpoint(router, token,
		"/api/v2/signals", "POST", helper.KeyModels{"signal": newSignal})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Read newSignal's ID from the response
	newSignalID, err := helper.GetResponseID(resp)
	assert.NoError(t, err)

	// authenticate as normal userB who has no access to new scenario
	token, err = helper.AuthenticateForTest(router, helper.UserBCredentials)
	assert.NoError(t, err)

	// Try to DELETE signal with no access
	// should result in unprocessable entity
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/signals/%v", newSignalID), "DELETE", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

	// authenticate as normal user
	token, err = helper.AuthenticateForTest(router, helper.UserACredentials)
	assert.NoError(t, err)

	// Count the number of all the input signals returned for component config
	initialNumber, err := helper.LengthOfResponse(router, token,
		fmt.Sprintf("/api/v2/signals?configID=%v&direction=in", configID), "GET", nil)
	assert.NoError(t, err)

	// add an output signal to make sure that counting of input signals works
	newSignal1.ConfigID = configID
	code, resp, err = helper.TestEndpoint(router, token,
		"/api/v2/signals", "POST", helper.KeyModels{"signal": newSignal1})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Read newSignalout's ID from the response
	newSignaloutID, err := helper.GetResponseID(resp)
	assert.NoError(t, err)

	// Delete the added newSignal
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/signals/%v", newSignalID), "DELETE", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare DELETE's response with the newSignal
	err = helper.CompareResponse(resp, helper.KeyModels{"signal": newSignal})
	assert.NoError(t, err)

	// Again count the number of all the input signals returned for component config
	finalNumber, err := helper.LengthOfResponse(router, token,
		fmt.Sprintf("/api/v2/signals?configID=%v&direction=in", configID), "GET", nil)
	assert.NoError(t, err)

	assert.Equal(t, initialNumber-1, finalNumber)

	// Delete the output signal
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/signals/%v", newSignaloutID), "DELETE", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)
}

func TestGetAllInputSignalsOfConfig(t *testing.T) {
	database.DropTables()
	database.MigrateModels()
	assert.NoError(t, helper.AddTestUsers())

	// prepare the content of the DB for testing
	// by adding a scenario and a IC to the DB
	// using the respective endpoints of the API
	_, _, configID := addScenarioAndICAndConfig()

	// authenticate as normal user
	token, err := helper.AuthenticateForTest(router, helper.UserACredentials)
	assert.NoError(t, err)

	// Count the number of all the input signals returned for component config
	initialNumber, err := helper.LengthOfResponse(router, token,
		fmt.Sprintf("/api/v2/signals?configID=%v&direction=in", configID), "GET", nil)
	assert.NoError(t, err)

	// test POST signals/ $newSignal
	newSignalA := SignalRequest{
		Name:      "inA",
		Unit:      "s",
		Direction: "in",
		Index:     1,
		ConfigID:  configID,
	}
	code, resp, err := helper.TestEndpoint(router, token,
		"/api/v2/signals", "POST", helper.KeyModels{"signal": newSignalA})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// add a second input signal
	newSignalB := SignalRequest{
		Name:      "inB",
		Unit:      "m",
		Direction: "in",
		Index:     2,
		ConfigID:  configID,
	}
	code, resp, err = helper.TestEndpoint(router, token,
		"/api/v2/signals", "POST", helper.KeyModels{"signal": newSignalB})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// add an output signal
	newSignal1.ConfigID = configID
	code, resp, err = helper.TestEndpoint(router, token,
		"/api/v2/signals", "POST", helper.KeyModels{"signal": newSignal1})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// add a second output signal
	newSignalBout := SignalRequest{
		Name:      "outB",
		Unit:      "A",
		Direction: "out",
		Index:     1,
		ConfigID:  configID,
	}
	code, resp, err = helper.TestEndpoint(router, token,
		"/api/v2/signals", "POST", helper.KeyModels{"signal": newSignalBout})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Again count the number of all the input signals returned for component config
	finalNumber, err := helper.LengthOfResponse(router, token,
		fmt.Sprintf("/api/v2/signals?configID=%v&direction=in", configID), "GET", nil)
	assert.NoError(t, err)

	assert.Equal(t, initialNumber+2, finalNumber)

	// Get the number of output signals
	outputNumber, err := helper.LengthOfResponse(router, token,
		fmt.Sprintf("/api/v2/signals?configID=%v&direction=out", configID), "GET", nil)
	assert.NoError(t, err)
	assert.Equal(t, initialNumber+2, outputNumber)

	// Try to get all signals for non-existing direction
	// should result in bad request
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/signals?configID=%v&direction=thisiswrong", configID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 400, code, "Response body: \n%v\n", resp)

	// authenticate as normal userB who has no access to new scenario
	token, err = helper.AuthenticateForTest(router, helper.UserBCredentials)
	assert.NoError(t, err)

	// try to get all input signals
	// should result in unprocessable entity
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/signals?configID=%v&direction=in", configID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

}
