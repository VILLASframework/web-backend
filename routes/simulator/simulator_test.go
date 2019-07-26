package simulator

import (
	"encoding/json"
	"testing"

	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/user"
)

// Test /simulator endpoints
func TestSimulatorEndpoints(t *testing.T) {

	var token string

	var myModels = []common.SimulationModelResponse{common.SimulationModelA_response, common.SimulationModelB_response}
	var msgModels = common.ResponseMsgSimulationModels{SimulationModels: myModels}
	var simulatorC_msg = common.ResponseMsgSimulator{Simulator: common.SimulatorC_response}
	var simulatorCupdated_msg = common.ResponseMsgSimulator{Simulator: common.SimulatorCUpdated_response}
	var mySimulators = []common.SimulatorResponse{common.SimulatorA_response, common.SimulatorB_response}
	var msgSimulators = common.ResponseMsgSimulators{Simulators: mySimulators}

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

	credjson, _ := json.Marshal(common.CredAdmin)
	msgOKjson, _ := json.Marshal(common.MsgOK)
	msgModelsjson, _ := json.Marshal(msgModels)
	msgSimulatorsjson, _ := json.Marshal(msgSimulators)
	msgSimulatorjson, _ := json.Marshal(simulatorC_msg)
	simulatorCjson, _ := json.Marshal(simulatorC_msg)
	simulatorCupdatedjson, _ := json.Marshal(simulatorCupdated_msg)

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
