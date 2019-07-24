package signal

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

var inSignalA = common.SignalResponse{
	Name:              "inSignal_A",
	Direction:         "in",
	Index:             0,
	Unit:              "A",
	SimulationModelID: 1,
}

var inSignalB = common.SignalResponse{
	Name:              "inSignal_B",
	Direction:         "in",
	Index:             1,
	Unit:              "A",
	SimulationModelID: 1,
}

var inSignalC = common.SignalResponse{
	Name:              "inSignal_C",
	Direction:         "in",
	Index:             2,
	Unit:              "A",
	SimulationModelID: 1,
}

var inSignalCupdated = common.Signal{
	Name:              "inSignalupdated_C",
	Direction:         "in",
	Index:             2,
	Unit:              "Ohm",
	SimulationModelID: 1,
}

var inSignalCupdatedResp = common.SignalResponse{
	Name:              inSignalCupdated.Name,
	Direction:         inSignalCupdated.Direction,
	Index:             inSignalCupdated.Index,
	Unit:              inSignalCupdated.Unit,
	SimulationModelID: inSignalCupdated.SimulationModelID,
}

var outSignalA = common.SignalResponse{
	Name:              "outSignal_A",
	Direction:         "out",
	Index:             0,
	Unit:              "V",
	SimulationModelID: 1,
}

var outSignalB = common.SignalResponse{
	Name:              "outSignal_B",
	Direction:         "out",
	Index:             1,
	Unit:              "V",
	SimulationModelID: 1,
}

var myInSignals = []common.SignalResponse{
	inSignalA,
	inSignalB,
}

var myOutSignals = []common.SignalResponse{
	outSignalA,
	outSignalB,
}

var msgInSignals = common.ResponseMsgSignals{
	Signals: myInSignals,
}

var msgInSignalCupdated = common.ResponseMsgSignal{
	Signal: inSignalCupdatedResp,
}

var msgOutSignals = common.ResponseMsgSignals{
	Signals: myOutSignals,
}

var msgInSignalC = common.ResponseMsgSignal{
	Signal: inSignalC,
}

// Test /models endpoints
func TestSignalEndpoints(t *testing.T) {

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

	credjson, err := json.Marshal(cred)
	if err != nil {
		panic(err)
	}

	msgOKjson, err := json.Marshal(msgOK)
	if err != nil {
		panic(err)
	}

	msgInSignalsjson, err := json.Marshal(msgInSignals)
	if err != nil {
		panic(err)
	}

	msgOutSignalsjson, err := json.Marshal(msgOutSignals)
	if err != nil {
		panic(err)
	}

	inSignalCjson, err := json.Marshal(inSignalC)
	if err != nil {
		panic(err)
	}

	msgInSignalCjson, err := json.Marshal(msgInSignalC)
	if err != nil {
		panic(err)
	}

	msgInSignalCupdatedjson, err := json.Marshal(msgInSignalCupdated)

	inSignalCupdatedjson, err := json.Marshal(inSignalCupdated)

	token = common.AuthenticateForTest(t, router, "/api/authenticate", "POST", credjson, 200)

	// test GET signals
	common.TestEndpoint(t, router, token, "/api/signals?modelID=1&direction=in", "GET", nil, 200, msgInSignalsjson)
	common.TestEndpoint(t, router, token, "/api/signals?modelID=1&direction=out", "GET", nil, 200, msgOutSignalsjson)

	// test POST signals
	common.TestEndpoint(t, router, token, "/api/signals", "POST", inSignalCjson, 200, msgOKjson)

	// test GET signals/:signalID
	common.TestEndpoint(t, router, token, "/api/signals/5", "GET", nil, 200, msgInSignalCjson)

	// test PUT signals/:signalID
	common.TestEndpoint(t, router, token, "/api/signals/5", "PUT", inSignalCupdatedjson, 200, msgOKjson)
	common.TestEndpoint(t, router, token, "/api/signals/5", "GET", nil, 200, msgInSignalCupdatedjson)

	// test DELETE signals/:signalID
	common.TestEndpoint(t, router, token, "/api/signals/5", "DELETE", nil, 200, msgOKjson)
	common.TestEndpoint(t, router, token, "/api/signals?modelID=1&direction=in", "GET", nil, 200, msgInSignalsjson)
	common.TestEndpoint(t, router, token, "/api/signals?modelID=1&direction=out", "GET", nil, 200, msgOutSignalsjson)

	// TODO test GET models/:ModelID to check if POST and DELETE adapt InputLength correctly??
	//common.TestEndpoint(t, router, token, "/api/models/1", "GET", nil, 200, string(msgModelAUpdated2json))

	// TODO add testing for other return codes

}
