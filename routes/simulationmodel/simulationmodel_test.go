package simulationmodel

import (
	"fmt"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
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
	token, _ := common.NewAuthenticateForTest(router,
		"/api/authenticate", "POST", common.AdminCredentials)

	// POST $newSimulatorA
	newSimulatorA := SimulatorRequest{
		UUID:       common.SimulatorA.UUID,
		Host:       common.SimulatorA.Host,
		Modeltype:  common.SimulatorA.Modeltype,
		State:      common.SimulatorA.State,
		Properties: common.SimulatorA.Properties,
	}
	_, resp, _ := common.NewTestEndpoint(router, token,
		"/api/simulators", "POST", common.KeyModels{"simulator": newSimulatorA})

	// Read newSimulator's ID from the response
	newSimulatorID, _ := common.GetResponseID(resp)

	// authenticate as normal user
	token, _ = common.NewAuthenticateForTest(router,
		"/api/authenticate", "POST", common.UserACredentials)

	// POST $newScenario
	newScenario := ScenarioRequest{
		Name:            common.ScenarioA.Name,
		Running:         common.ScenarioA.Running,
		StartParameters: common.ScenarioA.StartParameters,
	}
	_, resp, _ = common.NewTestEndpoint(router, token,
		"/api/scenarios", "POST", common.KeyModels{"scenario": newScenario})

	// Read newScenario's ID from the response
	newScenarioID, _ := common.GetResponseID(resp)

	return uint(newScenarioID), uint(newSimulatorID)
}

func TestMain(m *testing.M) {

	db = common.DummyInitDB()
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
	common.DropTables(db)
	common.MigrateModels(db)
	common.DummyAddOnlyUserTableWithAdminAndUsersDB(db)

	// prepare the content of the DB for testing
	// by adding a scenario and a simulator to the DB
	// using the respective endpoints of the API
	scenarioID, simulatorID := addScenarioAndSimulator()

	// authenticate as normal user
	token, err := common.NewAuthenticateForTest(router,
		"/api/authenticate", "POST", common.UserACredentials)
	assert.NoError(t, err)

	// test POST models/ $newSimulationModel
	newSimulationModel := SimulationModelRequest{
		Name:            common.SimulationModelA.Name,
		ScenarioID:      scenarioID,
		SimulatorID:     simulatorID,
		StartParameters: common.SimulationModelA.StartParameters,
	}
	code, resp, err := common.NewTestEndpoint(router, token,
		"/api/models", "POST", common.KeyModels{"model": newSimulationModel})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare POST's response with the newSimulationModel
	err = common.CompareResponse(resp, common.KeyModels{"model": newSimulationModel})
	assert.NoError(t, err)

	// Read newSimulationModel's ID from the response
	newSimulationModelID, err := common.GetResponseID(resp)
	assert.NoError(t, err)

	// Get the newSimulationModel
	code, resp, err = common.NewTestEndpoint(router, token,
		fmt.Sprintf("/api/models/%v", newSimulationModelID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare GET's response with the newSimulationModel
	err = common.CompareResponse(resp, common.KeyModels{"model": newSimulationModel})
	assert.NoError(t, err)

	// try to POST a malformed simulation model
	// Required fields are missing
	malformedNewSimulationModel := SimulationModelRequest{
		Name: "ThisIsAMalformedRequest",
	}
	// this should NOT work and return a unprocessable entity 442 status code
	code, resp, err = common.NewTestEndpoint(router, token,
		"/api/models", "POST", common.KeyModels{"model": malformedNewSimulationModel})
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

}

func TestUpdateSimulationModel(t *testing.T) {

	common.DropTables(db)
	common.MigrateModels(db)
	common.DummyAddOnlyUserTableWithAdminAndUsersDB(db)

	// prepare the content of the DB for testing
	// by adding a scenario and a simulator to the DB
	// using the respective endpoints of the API
	scenarioID, simulatorID := addScenarioAndSimulator()

	// authenticate as normal user
	token, err := common.NewAuthenticateForTest(router,
		"/api/authenticate", "POST", common.UserACredentials)
	assert.NoError(t, err)

	// test POST models/ $newSimulationModel
	newSimulationModel := SimulationModelRequest{
		Name:            common.SimulationModelA.Name,
		ScenarioID:      scenarioID,
		SimulatorID:     simulatorID,
		StartParameters: common.SimulationModelA.StartParameters,
	}
	code, resp, err := common.NewTestEndpoint(router, token,
		"/api/models", "POST", common.KeyModels{"model": newSimulationModel})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Read newSimulationModel's ID from the response
	newSimulationModelID, err := common.GetResponseID(resp)
	assert.NoError(t, err)

	updatedSimulationModel := SimulationModelRequest{
		Name:            common.SimulationModelB.Name,
		StartParameters: common.SimulationModelB.StartParameters,
	}

	code, resp, err = common.NewTestEndpoint(router, token,
		fmt.Sprintf("/api/models/%v", newSimulationModelID), "PUT", common.KeyModels{"model": updatedSimulationModel})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare PUT's response with the updatedSimulationModel
	err = common.CompareResponse(resp, common.KeyModels{"model": updatedSimulationModel})
	assert.NoError(t, err)

	// Get the updatedSimulationModel
	code, resp, err = common.NewTestEndpoint(router, token,
		fmt.Sprintf("/api/models/%v", newSimulationModelID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare GET's response with the updatedSimulationModel
	err = common.CompareResponse(resp, common.KeyModels{"model": updatedSimulationModel})
	assert.NoError(t, err)

	// try to update a simulation model that does not exist (should return not found 404 status code)
	code, resp, err = common.NewTestEndpoint(router, token,
		fmt.Sprintf("/api/models/%v", newSimulationModelID+1), "PUT", common.KeyModels{"scenario": updatedSimulationModel})
	assert.NoError(t, err)
	assert.Equalf(t, 404, code, "Response body: \n%v\n", resp)
}

func TestDeleteSimulationModel(t *testing.T) {
	common.DropTables(db)
	common.MigrateModels(db)
	common.DummyAddOnlyUserTableWithAdminAndUsersDB(db)

	// prepare the content of the DB for testing
	// by adding a scenario and a simulator to the DB
	// using the respective endpoints of the API
	scenarioID, simulatorID := addScenarioAndSimulator()

	// authenticate as normal user
	token, err := common.NewAuthenticateForTest(router,
		"/api/authenticate", "POST", common.UserACredentials)
	assert.NoError(t, err)

	// test POST models/ $newSimulationModel
	newSimulationModel := SimulationModelRequest{
		Name:            common.SimulationModelA.Name,
		ScenarioID:      scenarioID,
		SimulatorID:     simulatorID,
		StartParameters: common.SimulationModelA.StartParameters,
	}
	code, resp, err := common.NewTestEndpoint(router, token,
		"/api/models", "POST", common.KeyModels{"model": newSimulationModel})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Read newSimulationModel's ID from the response
	newSimulationModelID, err := common.GetResponseID(resp)
	assert.NoError(t, err)

	// Count the number of all the simulation models returned for scenario
	initialNumber, err := common.LengthOfResponse(router, token,
		fmt.Sprintf("/api/models?scenarioID=%v", scenarioID), "GET", nil)
	assert.NoError(t, err)

	// Delete the added newSimulationModel
	code, resp, err = common.NewTestEndpoint(router, token,
		fmt.Sprintf("/api/models/%v", newSimulationModelID), "DELETE", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare DELETE's response with the newSimulationModel
	err = common.CompareResponse(resp, common.KeyModels{"model": newSimulationModel})
	assert.NoError(t, err)

	// Again count the number of all the simulation models returned
	finalNumber, err := common.LengthOfResponse(router, token,
		fmt.Sprintf("/api/models?scenarioID=%v", scenarioID), "GET", nil)
	assert.NoError(t, err)

	assert.Equal(t, initialNumber-1, finalNumber)
}

func TestGetAllSimulationModelsOfScenario(t *testing.T) {
	common.DropTables(db)
	common.MigrateModels(db)
	common.DummyAddOnlyUserTableWithAdminAndUsersDB(db)

	// prepare the content of the DB for testing
	// by adding a scenario and a simulator to the DB
	// using the respective endpoints of the API
	scenarioID, simulatorID := addScenarioAndSimulator()

	// authenticate as normal user
	token, err := common.NewAuthenticateForTest(router,
		"/api/authenticate", "POST", common.UserACredentials)
	assert.NoError(t, err)

	// test POST models/ $newSimulationModel
	newSimulationModel := SimulationModelRequest{
		Name:            common.SimulationModelA.Name,
		ScenarioID:      scenarioID,
		SimulatorID:     simulatorID,
		StartParameters: common.SimulationModelA.StartParameters,
	}
	code, resp, err := common.NewTestEndpoint(router, token,
		"/api/models", "POST", common.KeyModels{"model": newSimulationModel})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Count the number of all the simulation models returned for scenario
	NumberOfSimulationModels, err := common.LengthOfResponse(router, token,
		fmt.Sprintf("/api/models?scenarioID=%v", scenarioID), "GET", nil)
	assert.NoError(t, err)

	assert.Equal(t, 1, NumberOfSimulationModels)

}

//// Test /models endpoints
//func TestSimulationModelEndpoints(t *testing.T) {
//
//	var token string
//
//	var myModels = []common.SimulationModelResponse{common.SimulationModelA_response, common.SimulationModelB_response}
//	var msgModels = common.ResponseMsgSimulationModels{SimulationModels: myModels}
//	var msgModel = common.ResponseMsgSimulationModel{SimulationModel: common.SimulationModelC_response}
//	var msgModelupdated = common.ResponseMsgSimulationModel{SimulationModel: common.SimulationModelCUpdated_response}
//
//	db := common.DummyInitDB()
//	defer db.Close()
//	common.DummyPopulateDB(db)
//
//	router := gin.Default()
//	api := router.Group("/api")
//
//	// All endpoints require authentication except when someone wants to
//	// login (POST /authenticate)
//	user.RegisterAuthenticate(api.Group("/authenticate"))
//
//	api.Use(user.Authentication(true))
//
//	RegisterSimulationModelEndpoints(api.Group("/models"))
//
//	credjson, _ := json.Marshal(common.CredUser)
//	msgOKjson, _ := json.Marshal(common.MsgOK)
//	msgModelsjson, _ := json.Marshal(msgModels)
//	msgModeljson, _ := json.Marshal(msgModel)
//	msgModelupdatedjson, _ := json.Marshal(msgModelupdated)
//
//	token = common.AuthenticateForTest(t, router, "/api/authenticate", "POST", credjson, 200)
//
//	// test GET models
//	common.TestEndpoint(t, router, token, "/api/models?scenarioID=1", "GET", nil, 200, msgModelsjson)
//
//	// test POST models
//	common.TestEndpoint(t, router, token, "/api/models", "POST", msgModeljson, 200, msgOKjson)
//
//	// test GET models/:ModelID to check if previous POST worked correctly
//	common.TestEndpoint(t, router, token, "/api/models/3", "GET", nil, 200, msgModeljson)
//
//	// test PUT models/:ModelID
//	common.TestEndpoint(t, router, token, "/api/models/3", "PUT", msgModelupdatedjson, 200, msgOKjson)
//	common.TestEndpoint(t, router, token, "/api/models/3", "GET", nil, 200, msgModelupdatedjson)
//
//	// test DELETE models/:ModelID
//	common.TestEndpoint(t, router, token, "/api/models/3", "DELETE", nil, 200, msgOKjson)
//	common.TestEndpoint(t, router, token, "/api/models?scenarioID=1", "GET", nil, 200, msgModelsjson)
//
//	// TODO add testing for other return codes
//
//}
