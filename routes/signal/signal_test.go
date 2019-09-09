package signal

import (
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
	RegisterSignalEndpoints(api.Group("/signals"))

	os.Exit(m.Run())
}

func TestAddSignal(t *testing.T) {
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

	// test POST signals/ $newSignal
	newSignal := SignalRequest{
		Name:              common.InSignalA.Name,
		Unit:              common.InSignalA.Unit,
		Direction:         common.InSignalA.Direction,
		Index:             1,
		SimulationModelID: simulationModelID,
	}
	code, resp, err := common.TestEndpoint(router, token,
		"/api/signals", "POST", common.KeyModels{"signal": newSignal})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare POST's response with the newSignal
	err = common.CompareResponse(resp, common.KeyModels{"signal": newSignal})
	assert.NoError(t, err)

	// Read newSignal's ID from the response
	newSignalID, err := common.GetResponseID(resp)
	assert.NoError(t, err)

	// Get the newSignal
	code, resp, err = common.TestEndpoint(router, token,
		fmt.Sprintf("/api/signals/%v", newSignalID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare GET's response with the newSSignal
	err = common.CompareResponse(resp, common.KeyModels{"signal": newSignal})
	assert.NoError(t, err)

	// try to POST a malformed signal
	// Required fields are missing
	malformedNewSignal := SignalRequest{
		Name: "ThisIsAMalformedRequest",
	}
	// this should NOT work and return a unprocessable entity 442 status code
	code, resp, err = common.TestEndpoint(router, token,
		"/api/signals", "POST", common.KeyModels{"model": malformedNewSignal})
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)
}

func TestUpdateSignal(t *testing.T) {
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

	// test POST signals/ $newSignal
	newSignal := SignalRequest{
		Name:              common.InSignalA.Name,
		Unit:              common.InSignalA.Unit,
		Direction:         common.InSignalA.Direction,
		Index:             1,
		SimulationModelID: simulationModelID,
	}
	code, resp, err := common.TestEndpoint(router, token,
		"/api/signals", "POST", common.KeyModels{"signal": newSignal})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Read newSignal's ID from the response
	newSignalID, err := common.GetResponseID(resp)
	assert.NoError(t, err)

	updatedSignal := SignalRequest{
		Name:  common.InSignalB.Name,
		Unit:  common.InSignalB.Unit,
		Index: 1,
	}
	code, resp, err = common.TestEndpoint(router, token,
		fmt.Sprintf("/api/signals/%v", newSignalID), "PUT", common.KeyModels{"signal": updatedSignal})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare PUT's response with the updatedSignal
	err = common.CompareResponse(resp, common.KeyModels{"signal": updatedSignal})
	assert.NoError(t, err)

	// Get the updatedSignal
	code, resp, err = common.TestEndpoint(router, token,
		fmt.Sprintf("/api/signals/%v", newSignalID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare GET's response with the updatedSignal
	err = common.CompareResponse(resp, common.KeyModels{"signal": updatedSignal})
	assert.NoError(t, err)

	// try to update a signal that does not exist (should return not found 404 status code)
	code, resp, err = common.TestEndpoint(router, token,
		fmt.Sprintf("/api/signals/%v", newSignalID+1), "PUT", common.KeyModels{"signal": updatedSignal})
	assert.NoError(t, err)
	assert.Equalf(t, 404, code, "Response body: \n%v\n", resp)

}

func TestDeleteSignal(t *testing.T) {
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

	// test POST signals/ $newSignal
	newSignal := SignalRequest{
		Name:              common.InSignalA.Name,
		Unit:              common.InSignalA.Unit,
		Direction:         common.InSignalA.Direction,
		Index:             1,
		SimulationModelID: simulationModelID,
	}
	code, resp, err := common.TestEndpoint(router, token,
		"/api/signals", "POST", common.KeyModels{"signal": newSignal})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Read newSignal's ID from the response
	newSignalID, err := common.GetResponseID(resp)
	assert.NoError(t, err)

	// Count the number of all the input signals returned for simulation model
	initialNumber, err := common.LengthOfResponse(router, token,
		fmt.Sprintf("/api/signals?modelID=%v&direction=in", simulationModelID), "GET", nil)
	assert.NoError(t, err)

	// add an output signal to make sure that counting of input signals works
	newSignalout := SignalRequest{
		Name:              common.OutSignalA.Name,
		Unit:              common.OutSignalA.Unit,
		Direction:         common.OutSignalA.Direction,
		Index:             1,
		SimulationModelID: simulationModelID,
	}
	code, resp, err = common.TestEndpoint(router, token,
		"/api/signals", "POST", common.KeyModels{"signal": newSignalout})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Delete the added newSignal
	code, resp, err = common.TestEndpoint(router, token,
		fmt.Sprintf("/api/signals/%v", newSignalID), "DELETE", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare DELETE's response with the newSignal
	err = common.CompareResponse(resp, common.KeyModels{"signal": newSignal})
	assert.NoError(t, err)

	// Again count the number of all the input signals returned for simulation model
	finalNumber, err := common.LengthOfResponse(router, token,
		fmt.Sprintf("/api/signals?modelID=%v&direction=in", simulationModelID), "GET", nil)
	assert.NoError(t, err)

	assert.Equal(t, initialNumber-1, finalNumber)
}

func TestGetAllInputSignalsOfSimulationModel(t *testing.T) {
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

	// Count the number of all the input signals returned for simulation model
	initialNumber, err := common.LengthOfResponse(router, token,
		fmt.Sprintf("/api/signals?modelID=%v&direction=in", simulationModelID), "GET", nil)
	assert.NoError(t, err)

	// test POST signals/ $newSignal
	newSignalA := SignalRequest{
		Name:              common.InSignalA.Name,
		Unit:              common.InSignalA.Unit,
		Direction:         common.InSignalA.Direction,
		Index:             1,
		SimulationModelID: simulationModelID,
	}
	code, resp, err := common.TestEndpoint(router, token,
		"/api/signals", "POST", common.KeyModels{"signal": newSignalA})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// add a second input signal
	newSignalB := SignalRequest{
		Name:              common.InSignalB.Name,
		Unit:              common.InSignalB.Unit,
		Direction:         common.InSignalB.Direction,
		Index:             2,
		SimulationModelID: simulationModelID,
	}
	code, resp, err = common.TestEndpoint(router, token,
		"/api/signals", "POST", common.KeyModels{"signal": newSignalB})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// add an output signal
	newSignalAout := SignalRequest{
		Name:              common.OutSignalA.Name,
		Unit:              common.OutSignalA.Unit,
		Direction:         common.OutSignalA.Direction,
		Index:             1,
		SimulationModelID: simulationModelID,
	}
	code, resp, err = common.TestEndpoint(router, token,
		"/api/signals", "POST", common.KeyModels{"signal": newSignalAout})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// add a second output signal
	newSignalBout := SignalRequest{
		Name:              common.OutSignalB.Name,
		Unit:              common.OutSignalB.Unit,
		Direction:         common.OutSignalB.Direction,
		Index:             1,
		SimulationModelID: simulationModelID,
	}
	code, resp, err = common.TestEndpoint(router, token,
		"/api/signals", "POST", common.KeyModels{"signal": newSignalBout})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Again count the number of all the input signals returned for simulation model
	finalNumber, err := common.LengthOfResponse(router, token,
		fmt.Sprintf("/api/signals?modelID=%v&direction=in", simulationModelID), "GET", nil)
	assert.NoError(t, err)

	assert.Equal(t, initialNumber+2, finalNumber)

}
