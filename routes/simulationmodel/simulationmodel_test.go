package simulationmodel

import (
	"fmt"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/database"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/helper"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/scenario"
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

func addScenarioAndSimulator() (scenarioID uint, simulatorID uint) {

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

	return uint(newScenarioID), uint(newSimulatorID)
}

func TestMain(m *testing.M) {

	db = database.InitDB(database.DB_TEST)
	defer db.Close()

	router = gin.Default()
	api := router.Group("/api")

	user.RegisterAuthenticate(api.Group("/authenticate"))
	api.Use(user.Authentication(true))
	RegisterSimulationModelEndpoints(api.Group("/models"))
	// scenario endpoints required here to first add a scenario to the DB
	// that can be associated with a new simulation model
	scenario.RegisterScenarioEndpoints(api.Group("/scenarios"))
	// simulator endpoints required here to first add a simulator to the DB
	// that can be associated with a new simulation model
	simulator.RegisterSimulatorEndpoints(api.Group("/simulators"))

	os.Exit(m.Run())
}

func TestAddSimulationModel(t *testing.T) {
	database.DropTables(db)
	database.MigrateModels(db)
	assert.NoError(t, database.DBAddAdminAndUser(db))

	// prepare the content of the DB for testing
	// by adding a scenario and a simulator to the DB
	// using the respective endpoints of the API
	scenarioID, simulatorID := addScenarioAndSimulator()

	// authenticate as normal user
	token, err := helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserACredentials)
	assert.NoError(t, err)

	// test POST models/ $newSimulationModel
	newSimulationModel := SimulationModelRequest{
		Name:            database.SimulationModelA.Name,
		ScenarioID:      scenarioID,
		SimulatorID:     simulatorID,
		StartParameters: database.SimulationModelA.StartParameters,
	}
	code, resp, err := helper.TestEndpoint(router, token,
		"/api/models", "POST", helper.KeyModels{"model": newSimulationModel})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare POST's response with the newSimulationModel
	err = helper.CompareResponse(resp, helper.KeyModels{"model": newSimulationModel})
	assert.NoError(t, err)

	// Read newSimulationModel's ID from the response
	newSimulationModelID, err := helper.GetResponseID(resp)
	assert.NoError(t, err)

	// Get the newSimulationModel
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/models/%v", newSimulationModelID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare GET's response with the newSimulationModel
	err = helper.CompareResponse(resp, helper.KeyModels{"model": newSimulationModel})
	assert.NoError(t, err)

	// try to POST a malformed simulation model
	// Required fields are missing
	malformedNewSimulationModel := SimulationModelRequest{
		Name: "ThisIsAMalformedRequest",
	}
	// this should NOT work and return a unprocessable entity 442 status code
	code, resp, err = helper.TestEndpoint(router, token,
		"/api/models", "POST", helper.KeyModels{"model": malformedNewSimulationModel})
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

}

func TestUpdateSimulationModel(t *testing.T) {

	database.DropTables(db)
	database.MigrateModels(db)
	assert.NoError(t, database.DBAddAdminAndUser(db))

	// prepare the content of the DB for testing
	// by adding a scenario and a simulator to the DB
	// using the respective endpoints of the API
	scenarioID, simulatorID := addScenarioAndSimulator()

	// authenticate as normal user
	token, err := helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserACredentials)
	assert.NoError(t, err)

	// test POST models/ $newSimulationModel
	newSimulationModel := SimulationModelRequest{
		Name:            database.SimulationModelA.Name,
		ScenarioID:      scenarioID,
		SimulatorID:     simulatorID,
		StartParameters: database.SimulationModelA.StartParameters,
	}
	code, resp, err := helper.TestEndpoint(router, token,
		"/api/models", "POST", helper.KeyModels{"model": newSimulationModel})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Read newSimulationModel's ID from the response
	newSimulationModelID, err := helper.GetResponseID(resp)
	assert.NoError(t, err)

	updatedSimulationModel := SimulationModelRequest{
		Name:            database.SimulationModelB.Name,
		StartParameters: database.SimulationModelB.StartParameters,
	}

	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/models/%v", newSimulationModelID), "PUT", helper.KeyModels{"model": updatedSimulationModel})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare PUT's response with the updatedSimulationModel
	err = helper.CompareResponse(resp, helper.KeyModels{"model": updatedSimulationModel})
	assert.NoError(t, err)

	// Get the updatedSimulationModel
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/models/%v", newSimulationModelID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare GET's response with the updatedSimulationModel
	err = helper.CompareResponse(resp, helper.KeyModels{"model": updatedSimulationModel})
	assert.NoError(t, err)

	// try to update a simulation model that does not exist (should return not found 404 status code)
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/models/%v", newSimulationModelID+1), "PUT", helper.KeyModels{"model": updatedSimulationModel})
	assert.NoError(t, err)
	assert.Equalf(t, 404, code, "Response body: \n%v\n", resp)
}

func TestDeleteSimulationModel(t *testing.T) {
	database.DropTables(db)
	database.MigrateModels(db)
	assert.NoError(t, database.DBAddAdminAndUser(db))

	// prepare the content of the DB for testing
	// by adding a scenario and a simulator to the DB
	// using the respective endpoints of the API
	scenarioID, simulatorID := addScenarioAndSimulator()

	// authenticate as normal user
	token, err := helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserACredentials)
	assert.NoError(t, err)

	// test POST models/ $newSimulationModel
	newSimulationModel := SimulationModelRequest{
		Name:            database.SimulationModelA.Name,
		ScenarioID:      scenarioID,
		SimulatorID:     simulatorID,
		StartParameters: database.SimulationModelA.StartParameters,
	}
	code, resp, err := helper.TestEndpoint(router, token,
		"/api/models", "POST", helper.KeyModels{"model": newSimulationModel})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Read newSimulationModel's ID from the response
	newSimulationModelID, err := helper.GetResponseID(resp)
	assert.NoError(t, err)

	// Count the number of all the simulation models returned for scenario
	initialNumber, err := helper.LengthOfResponse(router, token,
		fmt.Sprintf("/api/models?scenarioID=%v", scenarioID), "GET", nil)
	assert.NoError(t, err)

	// Delete the added newSimulationModel
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/models/%v", newSimulationModelID), "DELETE", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare DELETE's response with the newSimulationModel
	err = helper.CompareResponse(resp, helper.KeyModels{"model": newSimulationModel})
	assert.NoError(t, err)

	// Again count the number of all the simulation models returned
	finalNumber, err := helper.LengthOfResponse(router, token,
		fmt.Sprintf("/api/models?scenarioID=%v", scenarioID), "GET", nil)
	assert.NoError(t, err)

	assert.Equal(t, initialNumber-1, finalNumber)
}

func TestGetAllSimulationModelsOfScenario(t *testing.T) {
	database.DropTables(db)
	database.MigrateModels(db)
	assert.NoError(t, database.DBAddAdminAndUser(db))

	// prepare the content of the DB for testing
	// by adding a scenario and a simulator to the DB
	// using the respective endpoints of the API
	scenarioID, simulatorID := addScenarioAndSimulator()

	// authenticate as normal user
	token, err := helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserACredentials)
	assert.NoError(t, err)

	// test POST models/ $newSimulationModel
	newSimulationModel := SimulationModelRequest{
		Name:            database.SimulationModelA.Name,
		ScenarioID:      scenarioID,
		SimulatorID:     simulatorID,
		StartParameters: database.SimulationModelA.StartParameters,
	}
	code, resp, err := helper.TestEndpoint(router, token,
		"/api/models", "POST", helper.KeyModels{"model": newSimulationModel})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Count the number of all the simulation models returned for scenario
	NumberOfSimulationModels, err := helper.LengthOfResponse(router, token,
		fmt.Sprintf("/api/models?scenarioID=%v", scenarioID), "GET", nil)
	assert.NoError(t, err)

	assert.Equal(t, 1, NumberOfSimulationModels)

}
