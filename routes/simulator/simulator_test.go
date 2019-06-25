package simulator

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
	SimulationID: 1,
	SimulatorID:  1,
	StartParams:  "",
}

var model_B = common.SimulationModelResponse{
	ID:           2,
	Name:         "SimulationModel_B",
	OutputLength: 1,
	InputLength:  1,
	SimulationID: 1,
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
	UUID:          "1",
	Host:          "Host_A",
	ModelType:     "",
	Uptime:        0,
	State:         "",
	StateUpdateAt: "",
	Properties:    "",
	RawProperties: "",
}

var simulatorB = common.SimulatorResponse{
	ID:            2,
	UUID:          "2",
	Host:          "Host_B",
	ModelType:     "",
	Uptime:        0,
	State:         "",
	StateUpdateAt: "",
	Properties:    "",
	RawProperties: "",
}

var simulatorC = common.Simulator{
	ID:            3,
	UUID:          "3",
	Host:          "Host_C",
	Modeltype:     "",
	Uptime:        0,
	State:         "",
	StateUpdateAt: "",
	Properties:    "",
	RawProperties: "",
}

var simulatorCupdated = common.Simulator{
	ID:            3,
	UUID:          "3",
	Host:          "Host_Cupdated",
	Modeltype:     "",
	Uptime:        0,
	State:         "running",
	StateUpdateAt: "",
	Properties:    "",
	RawProperties: "",
}

var simulatorC_response = common.SimulatorResponse{
	ID:            3,
	UUID:          simulatorC.UUID,
	Host:          simulatorC.Host,
	ModelType:     simulatorC.Modeltype,
	Uptime:        simulatorC.Uptime,
	State:         simulatorC.State,
	StateUpdateAt: simulatorC.StateUpdateAt,
	Properties:    simulatorC.Properties,
	RawProperties: simulatorC.RawProperties,
}

var simulatorCupdated_response = common.SimulatorResponse{
	ID:            simulatorCupdated.ID,
	UUID:          simulatorCupdated.UUID,
	Host:          simulatorCupdated.Host,
	ModelType:     simulatorCupdated.Modeltype,
	Uptime:        simulatorCupdated.Uptime,
	State:         simulatorCupdated.State,
	StateUpdateAt: simulatorCupdated.StateUpdateAt,
	Properties:    simulatorCupdated.Properties,
	RawProperties: simulatorCupdated.RawProperties,
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

var msgSimulatorUpdated = common.ResponseMsgSimulator{
	Simulator: simulatorCupdated_response,
}

// Test /simulation endpoints
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

	msgSimulatorUpdatedjson, err := json.Marshal(msgSimulatorUpdated)
	if err != nil {
		panic(err)
	}

	simulatorCjson, err := json.Marshal(simulatorC)
	if err != nil {
		panic(err)
	}

	simulatorCupdatedjson, err := json.Marshal(simulatorCupdated)
	if err != nil {
		panic(err)
	}

	token = common.AuthenticateForTest(t, router, "/api/authenticate", "POST", credjson, 200)

	// test GET simulators/
	common.TestEndpoint(t, router, token, "/api/simulators", "GET", nil, 200, string(msgSimulatorsjson))

	// test POST simulators/
	common.TestEndpoint(t, router, token, "/api/simulators", "POST", simulatorCjson, 200, string(msgOKjson))

	// test GET simulators/:SimulatorID
	common.TestEndpoint(t, router, token, "/api/simulators/3", "GET", nil, 200, string(msgSimulatorjson))

	// test PUT simulators/:SimulatorID
	common.TestEndpoint(t, router, token, "/api/simulators/3", "PUT", simulatorCupdatedjson, 200, string(msgOKjson))
	common.TestEndpoint(t, router, token, "/api/simulators/3", "GET", nil, 200, string(msgSimulatorUpdatedjson))

	// test DELETE simulators/:SimulatorID
	common.TestEndpoint(t, router, token, "/api/simulators/3", "DELETE", nil, 200, string(msgOKjson))
	common.TestEndpoint(t, router, token, "/api/simulators", "GET", nil, 200, string(msgSimulatorsjson))

	// test GET simulators/:SimulatorID/models
	common.TestEndpoint(t, router, token, "/api/simulators/1/models", "GET", nil, 200, string(msgModelsjson))

	// TODO add tests for other return codes
}
