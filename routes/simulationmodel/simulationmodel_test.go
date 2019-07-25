package simulationmodel

import (
	"encoding/json"
	"testing"

	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/user"
)

// Test /models endpoints
func TestSimulationModelEndpoints(t *testing.T) {

	var token string

	var myModels = []common.SimulationModelResponse{common.SimulationModelA_response, common.SimulationModelB_response}
	var msgModels = common.ResponseMsgSimulationModels{SimulationModels: myModels}
	var msgModel = common.ResponseMsgSimulationModel{SimulationModel: common.SimulationModelC_response}
	var msgModelupdated = common.ResponseMsgSimulationModel{SimulationModel: common.SimulationModelCUpdated_response}

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

	credjson, _ := json.Marshal(common.CredUser)
	msgOKjson, _ := json.Marshal(common.MsgOK)
	msgModelsjson, _ := json.Marshal(msgModels)
	msgModeljson, _ := json.Marshal(msgModel)
	msgModelupdatedjson, _ := json.Marshal(msgModelupdated)

	token = common.AuthenticateForTest(t, router, "/api/authenticate", "POST", credjson, 200)

	// test GET models
	common.TestEndpoint(t, router, token, "/api/models?scenarioID=1", "GET", nil, 200, msgModelsjson)

	// test POST models
	common.TestEndpoint(t, router, token, "/api/models", "POST", msgModeljson, 200, msgOKjson)

	// test GET models/:ModelID to check if previous POST worked correctly
	common.TestEndpoint(t, router, token, "/api/models/3", "GET", nil, 200, msgModeljson)

	// test PUT models/:ModelID
	common.TestEndpoint(t, router, token, "/api/models/3", "PUT", msgModelupdatedjson, 200, msgOKjson)
	common.TestEndpoint(t, router, token, "/api/models/3", "GET", nil, 200, msgModelupdatedjson)

	// test DELETE models/:ModelID
	common.TestEndpoint(t, router, token, "/api/models/3", "DELETE", nil, 200, msgOKjson)
	common.TestEndpoint(t, router, token, "/api/models?scenarioID=1", "GET", nil, 200, msgModelsjson)

	// TODO add testing for other return codes

}
