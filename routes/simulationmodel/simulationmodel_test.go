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
	SimulationID: 1,
	SimulatorID:  1,
	StartParams:  "",
}

var modelAUpdated = common.SimulationModelResponse{
	ID:           1,
	Name:         "SimulationModel_A",
	OutputLength: 1,
	InputLength:  3,
	SimulationID: 1,
	SimulatorID:  1,
	StartParams:  "",
}

var modelAUpdated2 = common.SimulationModelResponse{
	ID:           1,
	Name:         "SimulationModel_A",
	OutputLength: 1,
	InputLength:  0,
	SimulationID: 1,
	SimulatorID:  1,
	StartParams:  "",
}

var modelB = common.SimulationModelResponse{
	ID:           2,
	Name:         "SimulationModel_B",
	OutputLength: 1,
	InputLength:  1,
	SimulationID: 1,
	SimulatorID:  1,
	StartParams:  "",
}

var modelC = common.SimulationModel{
	ID:              3,
	Name:            "SimulationModel_C",
	OutputLength:    1,
	InputLength:     1,
	SimulationID:    1,
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
	SimulationID:    modelC.SimulationID,
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
	SimulationID: modelC.SimulationID,
	SimulatorID:  modelC.SimulatorID,
	StartParams:  modelC.StartParameters,
}

var modelC_responseUpdated = common.SimulationModelResponse{
	ID:           modelC.ID,
	Name:         modelCupdated.Name,
	InputLength:  modelC.InputLength,
	OutputLength: modelC.OutputLength,
	SimulationID: modelC.SimulationID,
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

var msgModelAUpdated = common.ResponseMsgSimulationModel{
	SimulationModel: modelAUpdated,
}

var msgModelAUpdated2 = common.ResponseMsgSimulationModel{
	SimulationModel: modelAUpdated2,
}

var msgModelupdated = common.ResponseMsgSimulationModel{
	SimulationModel: modelC_responseUpdated,
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

var myInSignalsUpdated = []common.SignalResponse{
	inSignalA,
	inSignalB,
	inSignalC,
}

var myOutSignals = []common.SignalResponse{
	outSignalA,
	outSignalB,
}

var msgSignalsEmpty = common.ResponseMsgSignals{
	Signals: []common.SignalResponse{},
}

var msgInSignals = common.ResponseMsgSignals{
	Signals: myInSignals,
}

var msgInSignalsUpdated = common.ResponseMsgSignals{
	Signals: myInSignalsUpdated,
}

var msgOutSignals = common.ResponseMsgSignals{
	Signals: myOutSignals,
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

	msgModelAUpdatedjson, err := json.Marshal(msgModelAUpdated)
	if err != nil {
		panic(err)
	}

	msgModelAUpdated2json, err := json.Marshal(msgModelAUpdated2)
	if err != nil {
		panic(err)
	}

	msgSignalsEmptyjson, err := json.Marshal(msgSignalsEmpty)
	if err != nil {
		panic(err)
	}

	msgInSignalsjson, err := json.Marshal(msgInSignals)
	if err != nil {
		panic(err)
	}

	msgInSignalsUpdatedjson, err := json.Marshal(msgInSignalsUpdated)
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

	token = common.AuthenticateForTest(t, router, "/api/authenticate", "POST", credjson, 200)

	// test GET models
	common.TestEndpoint(t, router, token, "/api/models?simulationID=1", "GET", nil, 200, string(msgModelsjson))

	// test POST models
	common.TestEndpoint(t, router, token, "/api/models", "POST", modelCjson, 200, string(msgOKjson))

	// test GET models/:ModelID to check if previous POST worked correctly
	common.TestEndpoint(t, router, token, "/api/models/3", "GET", nil, 200, string(msgModeljson))

	// test PUT models/:ModelID
	common.TestEndpoint(t, router, token, "/api/models/3", "PUT", modelCupdatedjson, 200, string(msgOKjson))
	common.TestEndpoint(t, router, token, "/api/models/3", "GET", nil, 200, string(msgModelupdatedjson))

	// test DELETE models/:ModelID
	common.TestEndpoint(t, router, token, "/api/models/3", "DELETE", nil, 200, string(msgOKjson))
	common.TestEndpoint(t, router, token, "/api/models?simulationID=1", "GET", nil, 200, string(msgModelsjson))

	// test GET models/:ModelID/signals
	common.TestEndpoint(t, router, token, "/api/models/1/signals?direction=in", "GET", nil, 200, string(msgInSignalsjson))
	common.TestEndpoint(t, router, token, "/api/models/1/signals?direction=out", "GET", nil, 200, string(msgOutSignalsjson))

	// test PUT models/:ModelID/signals
	common.TestEndpoint(t, router, token, "/api/models/1/signals", "PUT", inSignalCjson, 200, string(msgOKjson))
	common.TestEndpoint(t, router, token, "/api/models/1/signals?direction=in", "GET", nil, 200, string(msgInSignalsUpdatedjson))

	// test GET models/:ModelID to check if PUT adapted InputLength correctly
	common.TestEndpoint(t, router, token, "/api/models/1", "GET", nil, 200, string(msgModelAUpdatedjson))

	// test DELETE models/:ModelID/signals
	common.TestEndpoint(t, router, token, "/api/models/1/signals?direction=in", "DELETE", nil, 200, string(msgOKjson))
	common.TestEndpoint(t, router, token, "/api/models/1/signals?direction=in", "GET", nil, 200, string(msgSignalsEmptyjson))
	common.TestEndpoint(t, router, token, "/api/models/1/signals?direction=out", "GET", nil, 200, string(msgOutSignalsjson))

	// test GET models/:ModelID to check if DELETE adapted InputLength correctly
	common.TestEndpoint(t, router, token, "/api/models/1", "GET", nil, 200, string(msgModelAUpdated2json))

	// TODO add testing for other return codes

}
