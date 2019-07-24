package simulator

import (
	"encoding/json"
	"github.com/jinzhu/gorm/dialects/postgres"
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
	Username: "User_0",
	Password: "xyz789",
}

var msgOK = common.ResponseMsg{
	Message: "OK.",
}

var model_A = common.SimulationModelResponse{
	ID:           1,
	Name:         "SimulationModel_A",
	OutputLength: 1,
	InputLength:  1,
	ScenarioID:   1,
	SimulatorID:  1,
	StartParams:  "",
}

var model_B = common.SimulationModelResponse{
	ID:           2,
	Name:         "SimulationModel_B",
	OutputLength: 1,
	InputLength:  1,
	ScenarioID:   1,
	SimulatorID:  1,
	StartParams:  "",
}

var myModels = []common.SimulationModelResponse{
	model_A,
	model_B,
}

var msgModels = common.ResponseMsgSimulationModels{
	SimulationModels: myModels,
}

var simulatorA = common.SimulatorResponse{
	ID:            1,
	UUID:          "4854af30-325f-44a5-ad59-b67b2597de68",
	Host:          "Host_A",
	Modeltype:     "ModelTypeA",
	Uptime:        0,
	State:         "running",
	StateUpdateAt: "placeholder",
	Properties:    postgres.Jsonb{json.RawMessage(`{"name" : "TestNameA", "category" : "CategoryA", "location" : "anywhere on earth", "type": "dummy"}`)},
	RawProperties: postgres.Jsonb{json.RawMessage(`{"name" : "TestNameA", "category" : "CategoryA", "location" : "anywhere on earth", "type": "dummy"}`)},
}

var simulatorB = common.SimulatorResponse{
	ID:            2,
	UUID:          "7be0322d-354e-431e-84bd-ae4c9633138b",
	Host:          "Host_B",
	Modeltype:     "ModelTypeB",
	Uptime:        0,
	State:         "idle",
	StateUpdateAt: "placeholder",
	Properties:    postgres.Jsonb{json.RawMessage(`{"name" : "TestNameB", "category" : "CategoryB", "location" : "where ever you want", "type": "generic"}`)},
	RawProperties: postgres.Jsonb{json.RawMessage(`{"name" : "TestNameB", "category" : "CategoryB", "location" : "where ever you want", "type": "generic"}`)},
}

var simulatorC = common.Simulator{
	ID:            3,
	UUID:          "6d9776bf-b693-45e8-97b6-4c13d151043f",
	Host:          "Host_C",
	Modeltype:     "ModelTypeC",
	Uptime:        0,
	State:         "idle",
	StateUpdateAt: "placeholder",
	Properties:    postgres.Jsonb{json.RawMessage(`{"name" : "TestNameC", "category" : "CategoryC", "location" : "my desk", "type": "blubb"}`)},
	RawProperties: postgres.Jsonb{json.RawMessage(`{"name" : "TestNameC", "category" : "CategoryC", "location" : "my desk", "type": "blubb"}`)},
}

var simulatorCupdated = common.Simulator{
	ID:            3,
	UUID:          "6d9776bf-b693-45e8-97b6-4c13d151043f",
	Host:          "Host_Cupdated",
	Modeltype:     "ModelTypeCUpdated",
	Uptime:        0,
	State:         "running",
	StateUpdateAt: "placeholder",
	Properties:    postgres.Jsonb{json.RawMessage(`{"name" : "TestNameCUpdate", "category" : "CategoryC", "location" : "my desk", "type": "blubb"}`)},
	RawProperties: postgres.Jsonb{json.RawMessage(`{"name" : "TestNameCUpdate", "category" : "CategoryC", "location" : "my desk", "type": "blubb"}`)},
}

var simulatorC_response = common.SimulatorResponse{
	ID:            simulatorC.ID,
	UUID:          simulatorC.UUID,
	Host:          simulatorC.Host,
	Modeltype:     simulatorC.Modeltype,
	Uptime:        simulatorC.Uptime,
	State:         simulatorC.State,
	StateUpdateAt: simulatorC.StateUpdateAt,
	Properties:    simulatorC.Properties,
	RawProperties: simulatorC.RawProperties,
}

var simulatorC_msg = common.ResponseMsgSimulator{
	Simulator: simulatorC_response,
}

var simulatorCupdated_response = common.SimulatorResponse{
	ID:            simulatorCupdated.ID,
	UUID:          simulatorCupdated.UUID,
	Host:          simulatorCupdated.Host,
	Modeltype:     simulatorCupdated.Modeltype,
	Uptime:        simulatorCupdated.Uptime,
	State:         simulatorCupdated.State,
	StateUpdateAt: simulatorCupdated.StateUpdateAt,
	Properties:    simulatorCupdated.Properties,
	RawProperties: simulatorCupdated.RawProperties,
}

var simulatorCupdated_msg = common.ResponseMsgSimulator{
	Simulator: simulatorCupdated_response,
}

var mySimulators = []common.SimulatorResponse{
	simulatorA,
	simulatorB,
}

var msgSimulators = common.ResponseMsgSimulators{
	Simulators: mySimulators,
}

var msgSimulator = common.ResponseMsgSimulator{
	Simulator: simulatorC_response,
}

// Test /simulator endpoints
func TestSimulatorEndpoints(t *testing.T) {

	db := common.DummyInitDB()
	defer db.Close()
	common.DummyPopulateDB(db)

	router := gin.Default()
	api := router.Group("/api")

	// All endpoints require authentication except when someone wants to
	// login (POST /authenticate)
	user.VisitorAuthenticate(api.Group("/authenticate"))

	api.Use(user.Authentication(true))

	RegisterSimulatorEndpoints(api.Group("/simulators"))

	credjson, err := json.Marshal(cred)

	msgOKjson, err := json.Marshal(msgOK)
	if err != nil {
		panic(err)
	}

	msgModelsjson, err := json.Marshal(msgModels)
	if err != nil {
		panic(err)
	}

	msgSimulatorsjson, err := json.Marshal(msgSimulators)
	if err != nil {
		panic(err)
	}

	msgSimulatorjson, err := json.Marshal(msgSimulator)
	if err != nil {
		panic(err)
	}

	simulatorCjson, err := json.Marshal(simulatorC_msg)
	if err != nil {
		panic(err)
	}

	simulatorCupdatedjson, err := json.Marshal(simulatorCupdated_msg)
	if err != nil {
		panic(err)
	}

	token = common.AuthenticateForTest(t, router, "/api/authenticate", "POST", credjson, 200)

	// test GET simulators/
	common.TestEndpoint(t, router, token, "/api/simulators", "GET", nil, 200, msgSimulatorsjson)

	// test POST simulators/
	common.TestEndpoint(t, router, token, "/api/simulators", "POST", simulatorCjson, 200, msgOKjson)

	// test GET simulators/:SimulatorID
	common.TestEndpoint(t, router, token, "/api/simulators/3", "GET", nil, 200, msgSimulatorjson)

	// test PUT simulators/:SimulatorID
	common.TestEndpoint(t, router, token, "/api/simulators/3", "PUT", simulatorCupdatedjson, 200, msgOKjson)
	common.TestEndpoint(t, router, token, "/api/simulators/3", "GET", nil, 200, simulatorCupdatedjson)

	// test DELETE simulators/:SimulatorID
	common.TestEndpoint(t, router, token, "/api/simulators/3", "DELETE", nil, 200, msgOKjson)
	common.TestEndpoint(t, router, token, "/api/simulators", "GET", nil, 200, msgSimulatorsjson)

	// test GET simulators/:SimulatorID/models
	common.TestEndpoint(t, router, token, "/api/simulators/1/models", "GET", nil, 200, msgModelsjson)

	// TODO add tests for other return codes
}
