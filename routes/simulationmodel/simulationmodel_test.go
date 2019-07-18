package simulationmodel

import (
	"encoding/json"
	"testing"

	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/user"
)

var token string

type credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var cred = credentials{
	Username: "User_A",
	Password: "abc123",
}

var msgOK = common.ResponseMsg{
	Message: "OK.",
}

var modelA = common.SimulationModelResponse{
	ID:           1,
	Name:         "SimulationModel_A",
	OutputLength: 1,
	InputLength:  1,
	ScenarioID:   1,
	SimulatorID:  1,
	StartParams:  "",
}

var modelB = common.SimulationModelResponse{
	ID:           2,
	Name:         "SimulationModel_B",
	OutputLength: 1,
	InputLength:  1,
	ScenarioID:   1,
	SimulatorID:  1,
	StartParams:  "",
}

var modelC = common.SimulationModel{
	ID:              3,
	Name:            "SimulationModel_C",
	OutputLength:    1,
	InputLength:     1,
	ScenarioID:      1,
	SimulatorID:     1,
	StartParameters: "test",
	InputMapping:    nil,
	OutputMapping:   nil,
}

var modelCupdated = common.SimulationModel{
	ID:              modelC.ID,
	Name:            "SimulationModel_CUpdated",
	OutputLength:    modelC.OutputLength,
	InputLength:     modelC.InputLength,
	ScenarioID:      modelC.ScenarioID,
	SimulatorID:     2,
	StartParameters: modelC.StartParameters,
	InputMapping:    modelC.InputMapping,
	OutputMapping:   modelC.OutputMapping,
}

var modelC_response = common.SimulationModelResponse{
	ID:           modelC.ID,
	Name:         modelC.Name,
	InputLength:  modelC.InputLength,
	OutputLength: modelC.OutputLength,
	ScenarioID:   modelC.ScenarioID,
	SimulatorID:  modelC.SimulatorID,
	StartParams:  modelC.StartParameters,
}

var modelC_responseUpdated = common.SimulationModelResponse{
	ID:           modelC.ID,
	Name:         modelCupdated.Name,
	InputLength:  modelC.InputLength,
	OutputLength: modelC.OutputLength,
	ScenarioID:   modelC.ScenarioID,
	SimulatorID:  modelCupdated.SimulatorID,
	StartParams:  modelC.StartParameters,
}

var myModels = []common.SimulationModelResponse{
	modelA,
	modelB,
}

var msgModels = common.ResponseMsgSimulationModels{
	SimulationModels: myModels,
}

var msgModel = common.ResponseMsgSimulationModel{
	SimulationModel: modelC_response,
}

var msgModelupdated = common.ResponseMsgSimulationModel{
	SimulationModel: modelC_responseUpdated,
}

// Test /models endpoints
func TestSimulationModelEndpoints(t *testing.T) {

	db := common.DummyInitDB()
	defer db.Close()
	common.DummyPopulateDB(db)

	router := gin.Default()
	api := router.Group("/api")

	// All endpoints require authentication except when someone wants to
	// login (POST /authenticate)
	user.VisitorAuthenticate(api.Group("/authenticate"))

	api.Use(user.Authentication(true))

	RegisterSimulationModelEndpoints(api.Group("/models"))

	credjson, err := json.Marshal(cred)
	if err != nil {
		panic(err)
	}

	msgOKjson, err := json.Marshal(msgOK)
	if err != nil {
		panic(err)
	}

	msgModelsjson, err := json.Marshal(msgModels)
	if err != nil {
		panic(err)
	}

	msgModeljson, err := json.Marshal(msgModel)
	if err != nil {
		panic(err)
	}

	msgModelupdatedjson, err := json.Marshal(msgModelupdated)
	if err != nil {
		panic(err)
	}

	modelCjson, err := json.Marshal(modelC)
	if err != nil {
		panic(err)
	}

	modelCupdatedjson, err := json.Marshal(modelCupdated)
	if err != nil {
		panic(err)
	}

	token = common.AuthenticateForTest(t, router, "/api/authenticate", "POST", credjson, 200)

	// test GET models
	common.TestEndpoint(t, router, token, "/api/models?scenarioID=1", "GET", nil, 200, string(msgModelsjson))

	// test POST models
	common.TestEndpoint(t, router, token, "/api/models", "POST", modelCjson, 200, string(msgOKjson))

	// test GET models/:ModelID to check if previous POST worked correctly
	common.TestEndpoint(t, router, token, "/api/models/3", "GET", nil, 200, string(msgModeljson))

	// test PUT models/:ModelID
	common.TestEndpoint(t, router, token, "/api/models/3", "PUT", modelCupdatedjson, 200, string(msgOKjson))
	common.TestEndpoint(t, router, token, "/api/models/3", "GET", nil, 200, string(msgModelupdatedjson))

	// test DELETE models/:ModelID
	common.TestEndpoint(t, router, token, "/api/models/3", "DELETE", nil, 200, string(msgOKjson))
	common.TestEndpoint(t, router, token, "/api/models?scenarioID=1", "GET", nil, 200, string(msgModelsjson))

	// TODO add testing for other return codes

}
