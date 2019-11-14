package signal

import (
	"fmt"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/configuration"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/helper"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/scenario"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/simulationmodel"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/simulator"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/user"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var router *gin.Engine
var db *gorm.DB

type SignalRequest struct {
	Name              string `json:"name,omitempty"`
	Unit              string `json:"unit,omitempty"`
	Index             uint   `json:"index,omitempty"`
	Direction         string `json:"direction,omitempty"`
	SimulationModelID uint   `json:"simulationModelID,omitempty"`
}

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

	// test POST models/ $newSimulationModel
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

	// add the guest user to the new scenario
	_, resp, _ = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/scenarios/%v/user?username=User_C", newScenarioID), "PUT", nil)

	return uint(newScenarioID), uint(newSimulatorID), uint(newSimulationModelID)
}

func TestMain(m *testing.M) {
	err := configuration.InitConfig()
	if err != nil {
		panic(m)
	}

	db, err = database.InitDB(configuration.GolbalConfig)
	if err != nil {
		panic(m)
	}
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
	RegisterSignalEndpoints(api.Group("/signals"))

	os.Exit(m.Run())
}

func TestAddSignal(t *testing.T) {
	database.DropTables(db)
	database.MigrateModels(db)
	assert.NoError(t, database.DBAddAdminAndUserAndGuest(db))

	// prepare the content of the DB for testing
	// by adding a scenario and a simulator to the DB
	// using the respective endpoints of the API
	_, _, simulationModelID := addScenarioAndSimulatorAndSimulationModel()

	// authenticate as normal user
	token, err := helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserACredentials)
	assert.NoError(t, err)

	newSignal := SignalRequest{
		Name:              database.InSignalA.Name,
		Unit:              database.InSignalA.Unit,
		Direction:         database.InSignalA.Direction,
		Index:             1,
		SimulationModelID: simulationModelID,
	}

	// authenticate as normal userB who has no access to new scenario
	token, err = helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserBCredentials)
	assert.NoError(t, err)

	// try to POST to simulation model without access
	// should result in unprocessable entity
	code, resp, err := helper.TestEndpoint(router, token,
		"/api/signals", "POST", helper.KeyModels{"signal": newSignal})
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

	// authenticate as normal user
	token, err = helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserACredentials)
	assert.NoError(t, err)

	// try to POST a signal with non JSON body
	// should result in a bad request
	code, resp, err = helper.TestEndpoint(router, token,
		"/api/signals", "POST", "this is not a JSON")
	assert.NoError(t, err)
	assert.Equalf(t, 400, code, "Response body: \n%v\n", resp)

	// test POST signals/ $newSignal
	code, resp, err = helper.TestEndpoint(router, token,
		"/api/signals", "POST", helper.KeyModels{"signal": newSignal})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare POST's response with the newSignal
	err = helper.CompareResponse(resp, helper.KeyModels{"signal": newSignal})
	assert.NoError(t, err)

	// Read newSignal's ID from the response
	newSignalID, err := helper.GetResponseID(resp)
	assert.NoError(t, err)

	// Get the newSignal
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/signals/%v", newSignalID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare GET's response with the newSSignal
	err = helper.CompareResponse(resp, helper.KeyModels{"signal": newSignal})
	assert.NoError(t, err)

	// try to POST a malformed signal
	// Required fields are missing
	malformedNewSignal := SignalRequest{
		Name: "ThisIsAMalformedRequest",
	}
	// this should NOT work and return a unprocessable entity 442 status code
	code, resp, err = helper.TestEndpoint(router, token,
		"/api/signals", "POST", helper.KeyModels{"model": malformedNewSignal})
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

	// authenticate as normal userB who has no access to new scenario
	token, err = helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserBCredentials)
	assert.NoError(t, err)

	// Try to Get the newSignal as user B
	// should result in unprocessable entity
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/signals/%v", newSignalID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)
}

func TestUpdateSignal(t *testing.T) {
	database.DropTables(db)
	database.MigrateModels(db)
	assert.NoError(t, database.DBAddAdminAndUserAndGuest(db))

	// prepare the content of the DB for testing
	// by adding a scenario and a simulator to the DB
	// using the respective endpoints of the API
	_, _, simulationModelID := addScenarioAndSimulatorAndSimulationModel()

	// authenticate as normal user
	token, err := helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserACredentials)
	assert.NoError(t, err)

	// test POST signals/ $newSignal
	newSignal := SignalRequest{
		Name:              database.InSignalA.Name,
		Unit:              database.InSignalA.Unit,
		Direction:         database.InSignalA.Direction,
		Index:             1,
		SimulationModelID: simulationModelID,
	}
	code, resp, err := helper.TestEndpoint(router, token,
		"/api/signals", "POST", helper.KeyModels{"signal": newSignal})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Read newSignal's ID from the response
	newSignalID, err := helper.GetResponseID(resp)
	assert.NoError(t, err)

	updatedSignal := SignalRequest{
		Name:  database.InSignalB.Name,
		Unit:  database.InSignalB.Unit,
		Index: 1,
	}

	// authenticate as normal userB who has no access to new scenario
	token, err = helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserBCredentials)
	assert.NoError(t, err)

	// try to PUT signal without access
	// should result in unprocessable entity
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/signals/%v", newSignalID), "PUT", helper.KeyModels{"signal": updatedSignal})
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

	// authenticate as guest user
	token, err = helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.GuestCredentials)
	assert.NoError(t, err)

	// try to update signal as guest who has access to scenario
	// should result in unprocessable entity
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/signals/%v", newSignalID), "PUT", helper.KeyModels{"signal": updatedSignal})
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

	// authenticate as normal user
	token, err = helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserACredentials)
	assert.NoError(t, err)

	// try to PUT with non JSON body
	// should result in bad request
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/signals/%v", newSignalID), "PUT", "This is not JSON")
	assert.NoError(t, err)
	assert.Equalf(t, 400, code, "Response body: \n%v\n", resp)

	// test PUT
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/signals/%v", newSignalID), "PUT", helper.KeyModels{"signal": updatedSignal})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare PUT's response with the updatedSignal
	err = helper.CompareResponse(resp, helper.KeyModels{"signal": updatedSignal})
	assert.NoError(t, err)

	// Get the updatedSignal
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/signals/%v", newSignalID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare GET's response with the updatedSignal
	err = helper.CompareResponse(resp, helper.KeyModels{"signal": updatedSignal})
	assert.NoError(t, err)

	// try to update a signal that does not exist (should return not found 404 status code)
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/signals/%v", newSignalID+1), "PUT", helper.KeyModels{"signal": updatedSignal})
	assert.NoError(t, err)
	assert.Equalf(t, 404, code, "Response body: \n%v\n", resp)

}

func TestDeleteSignal(t *testing.T) {
	database.DropTables(db)
	database.MigrateModels(db)
	assert.NoError(t, database.DBAddAdminAndUserAndGuest(db))

	// prepare the content of the DB for testing
	// by adding a scenario and a simulator to the DB
	// using the respective endpoints of the API
	_, _, simulationModelID := addScenarioAndSimulatorAndSimulationModel()

	// authenticate as normal user
	token, err := helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserACredentials)
	assert.NoError(t, err)

	// test POST signals/ $newSignal
	newSignal := SignalRequest{
		Name:              database.InSignalA.Name,
		Unit:              database.InSignalA.Unit,
		Direction:         database.InSignalA.Direction,
		Index:             1,
		SimulationModelID: simulationModelID,
	}
	code, resp, err := helper.TestEndpoint(router, token,
		"/api/signals", "POST", helper.KeyModels{"signal": newSignal})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Read newSignal's ID from the response
	newSignalID, err := helper.GetResponseID(resp)
	assert.NoError(t, err)

	// authenticate as normal userB who has no access to new scenario
	token, err = helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserBCredentials)
	assert.NoError(t, err)

	// Try to DELETE signal with no access
	// should result in unprocessable entity
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/signals/%v", newSignalID), "DELETE", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

	// authenticate as normal user
	token, err = helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserACredentials)
	assert.NoError(t, err)

	// Count the number of all the input signals returned for simulation model
	initialNumber, err := helper.LengthOfResponse(router, token,
		fmt.Sprintf("/api/signals?modelID=%v&direction=in", simulationModelID), "GET", nil)
	assert.NoError(t, err)

	// add an output signal to make sure that counting of input signals works
	newSignalout := SignalRequest{
		Name:              database.OutSignalA.Name,
		Unit:              database.OutSignalA.Unit,
		Direction:         database.OutSignalA.Direction,
		Index:             1,
		SimulationModelID: simulationModelID,
	}
	code, resp, err = helper.TestEndpoint(router, token,
		"/api/signals", "POST", helper.KeyModels{"signal": newSignalout})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Read newSignalout's ID from the response
	newSignaloutID, err := helper.GetResponseID(resp)
	assert.NoError(t, err)

	// Delete the added newSignal
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/signals/%v", newSignalID), "DELETE", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare DELETE's response with the newSignal
	err = helper.CompareResponse(resp, helper.KeyModels{"signal": newSignal})
	assert.NoError(t, err)

	// Again count the number of all the input signals returned for simulation model
	finalNumber, err := helper.LengthOfResponse(router, token,
		fmt.Sprintf("/api/signals?modelID=%v&direction=in", simulationModelID), "GET", nil)
	assert.NoError(t, err)

	assert.Equal(t, initialNumber-1, finalNumber)

	// Delete the output signal
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/signals/%v", newSignaloutID), "DELETE", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)
}

func TestGetAllInputSignalsOfSimulationModel(t *testing.T) {
	database.DropTables(db)
	database.MigrateModels(db)
	assert.NoError(t, database.DBAddAdminAndUserAndGuest(db))

	// prepare the content of the DB for testing
	// by adding a scenario and a simulator to the DB
	// using the respective endpoints of the API
	_, _, simulationModelID := addScenarioAndSimulatorAndSimulationModel()

	// authenticate as normal user
	token, err := helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserACredentials)
	assert.NoError(t, err)

	// Count the number of all the input signals returned for simulation model
	initialNumber, err := helper.LengthOfResponse(router, token,
		fmt.Sprintf("/api/signals?modelID=%v&direction=in", simulationModelID), "GET", nil)
	assert.NoError(t, err)

	// test POST signals/ $newSignal
	newSignalA := SignalRequest{
		Name:              database.InSignalA.Name,
		Unit:              database.InSignalA.Unit,
		Direction:         database.InSignalA.Direction,
		Index:             1,
		SimulationModelID: simulationModelID,
	}
	code, resp, err := helper.TestEndpoint(router, token,
		"/api/signals", "POST", helper.KeyModels{"signal": newSignalA})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// add a second input signal
	newSignalB := SignalRequest{
		Name:              database.InSignalB.Name,
		Unit:              database.InSignalB.Unit,
		Direction:         database.InSignalB.Direction,
		Index:             2,
		SimulationModelID: simulationModelID,
	}
	code, resp, err = helper.TestEndpoint(router, token,
		"/api/signals", "POST", helper.KeyModels{"signal": newSignalB})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// add an output signal
	newSignalAout := SignalRequest{
		Name:              database.OutSignalA.Name,
		Unit:              database.OutSignalA.Unit,
		Direction:         database.OutSignalA.Direction,
		Index:             1,
		SimulationModelID: simulationModelID,
	}
	code, resp, err = helper.TestEndpoint(router, token,
		"/api/signals", "POST", helper.KeyModels{"signal": newSignalAout})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// add a second output signal
	newSignalBout := SignalRequest{
		Name:              database.OutSignalB.Name,
		Unit:              database.OutSignalB.Unit,
		Direction:         database.OutSignalB.Direction,
		Index:             1,
		SimulationModelID: simulationModelID,
	}
	code, resp, err = helper.TestEndpoint(router, token,
		"/api/signals", "POST", helper.KeyModels{"signal": newSignalBout})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Again count the number of all the input signals returned for simulation model
	finalNumber, err := helper.LengthOfResponse(router, token,
		fmt.Sprintf("/api/signals?modelID=%v&direction=in", simulationModelID), "GET", nil)
	assert.NoError(t, err)

	assert.Equal(t, initialNumber+2, finalNumber)

	// Get the number of output signals
	outputNumber, err := helper.LengthOfResponse(router, token,
		fmt.Sprintf("/api/signals?modelID=%v&direction=out", simulationModelID), "GET", nil)
	assert.NoError(t, err)
	assert.Equal(t, initialNumber+2, outputNumber)

	// Try to get all signals for non-existing direction
	// should result in bad request
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/signals?modelID=%v&direction=thisiswrong", simulationModelID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 400, code, "Response body: \n%v\n", resp)

	// authenticate as normal userB who has no access to new scenario
	token, err = helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserBCredentials)
	assert.NoError(t, err)

	// try to get all input signals
	// should result in unprocessable entity
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/signals?modelID=%v&direction=in", simulationModelID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

}
