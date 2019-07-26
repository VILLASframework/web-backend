package signal

import (
	"encoding/json"
	"testing"

	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/user"
)

// Test /models endpoints
func TestSignalEndpoints(t *testing.T) {

	var token string

	var myInSignals = []common.SignalResponse{common.InSignalA_response, common.InSignalB_response}
	var myOutSignals = []common.SignalResponse{common.OutSignalA_response, common.OutSignalB_response}
	var msgInSignals = common.ResponseMsgSignals{Signals: myInSignals}
	var msgInSignalC = common.ResponseMsgSignal{Signal: common.InSignalC_response}
	var msgInSignalCupdated = common.ResponseMsgSignal{Signal: common.InSignalCUpdated_response}
	var msgOutSignals = common.ResponseMsgSignals{Signals: myOutSignals}

	db := common.DummyInitDB()
	defer db.Close()
	common.DummyPopulateDB(db)

	router := gin.Default()
	api := router.Group("/api")

	// All endpoints require authentication except when someone wants to
	// login (POST /authenticate)
	user.VisitorAuthenticate(api.Group("/authenticate"))

	api.Use(user.Authentication(true))

	RegisterSignalEndpoints(api.Group("/signals"))

	credjson, _ := json.Marshal(common.CredUser)
	msgOKjson, _ := json.Marshal(common.MsgOK)
	msgInSignalsjson, _ := json.Marshal(msgInSignals)
	msgOutSignalsjson, _ := json.Marshal(msgOutSignals)
	inSignalCjson, _ := json.Marshal(msgInSignalC)
	inSignalCupdatedjson, _ := json.Marshal(msgInSignalCupdated)

	token = common.AuthenticateForTest(t, router, "/api/authenticate", "POST", credjson, 200)

	// test GET signals
	common.TestEndpoint(t, router, token, "/api/signals?modelID=1&direction=in", "GET", nil, 200, msgInSignalsjson)
	common.TestEndpoint(t, router, token, "/api/signals?modelID=1&direction=out", "GET", nil, 200, msgOutSignalsjson)

	// test POST signals
	common.TestEndpoint(t, router, token, "/api/signals", "POST", inSignalCjson, 200, msgOKjson)

	// test GET signals/:signalID
	common.TestEndpoint(t, router, token, "/api/signals/5", "GET", nil, 200, inSignalCjson)

	// test PUT signals/:signalID
	common.TestEndpoint(t, router, token, "/api/signals/5", "PUT", inSignalCupdatedjson, 200, msgOKjson)
	common.TestEndpoint(t, router, token, "/api/signals/5", "GET", nil, 200, inSignalCupdatedjson)

	// test DELETE signals/:signalID
	common.TestEndpoint(t, router, token, "/api/signals/5", "DELETE", nil, 200, msgOKjson)
	common.TestEndpoint(t, router, token, "/api/signals?modelID=1&direction=in", "GET", nil, 200, msgInSignalsjson)
	common.TestEndpoint(t, router, token, "/api/signals?modelID=1&direction=out", "GET", nil, 200, msgOutSignalsjson)

	// TODO test GET models/:ModelID to check if POST and DELETE adapt InputLength correctly??
	//common.TestEndpoint(t, router, token, "/api/models/1", "GET", nil, 200, string(msgModelAUpdated2json))

	// TODO add testing for other return codes

}
