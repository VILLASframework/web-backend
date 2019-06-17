package visualization

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

var visA = common.VisualizationResponse{
	ID:           1,
	Name:         "Visualization_A",
	Grid:         15,
	SimulationID: 1,
}

var visB = common.VisualizationResponse{
	ID:           2,
	Name:         "Visualization_B",
	Grid:         15,
	SimulationID: 1,
}

var visC = common.Visualization{
	ID:           3,
	Name:         "Visualization_C",
	Grid:         99,
	SimulationID: 1,
}

var visCupdated = common.Visualization{
	ID:           visC.ID,
	Name:         "Visualization_CUpdated",
	SimulationID: visC.SimulationID,
	Grid:         visC.Grid,
}

var visC_response = common.VisualizationResponse{
	ID:           visC.ID,
	Name:         visC.Name,
	Grid:         visC.Grid,
	SimulationID: visC.SimulationID,
}

var visC_responseUpdated = common.VisualizationResponse{
	ID:           visCupdated.ID,
	Name:         visCupdated.Name,
	Grid:         visCupdated.Grid,
	SimulationID: visCupdated.SimulationID,
}

var myVisualizations = []common.VisualizationResponse{
	visA,
	visB,
}

var msgVisualizations = common.ResponseMsgVisualizations{
	Visualizations: myVisualizations,
}

var msgVis = common.ResponseMsgVisualization{
	Visualization: visC_response,
}

var msgVisupdated = common.ResponseMsgVisualization{
	Visualization: visC_responseUpdated,
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

	RegisterVisualizationEndpoints(api.Group("/visualizations"))

	credjson, err := json.Marshal(cred)
	if err != nil {
		panic(err)
	}

	msgOKjson, err := json.Marshal(msgOK)
	if err != nil {
		panic(err)
	}

	msgVisualizationsjson, err := json.Marshal(msgVisualizations)
	if err != nil {
		panic(err)
	}

	msgVisjson, err := json.Marshal(msgVis)
	if err != nil {
		panic(err)
	}

	msgVisupdatedjson, err := json.Marshal(msgVisupdated)
	if err != nil {
		panic(err)
	}

	visCjson, err := json.Marshal(visC)
	if err != nil {
		panic(err)
	}

	visCupdatedjson, err := json.Marshal(visCupdated)
	if err != nil {
		panic(err)
	}

	token = common.AuthenticateForTest(t, router, "/api/authenticate", "POST", credjson, 200)

	// test GET models
	common.TestEndpoint(t, router, token, "/api/visualizations?simulationID=1", "GET", nil, 200, string(msgVisualizationsjson))

	// test POST models
	common.TestEndpoint(t, router, token, "/api/visualizations", "POST", visCjson, 200, string(msgOKjson))

	// test GET models/:ModelID to check if previous POST worked correctly
	common.TestEndpoint(t, router, token, "/api/visualizations/3", "GET", nil, 200, string(msgVisjson))

	// test PUT models/:ModelID
	common.TestEndpoint(t, router, token, "/api/visualizations/3", "PUT", visCupdatedjson, 200, string(msgOKjson))
	common.TestEndpoint(t, router, token, "/api/visualizations/3", "GET", nil, 200, string(msgVisupdatedjson))

	// test DELETE models/:ModelID
	common.TestEndpoint(t, router, token, "/api/visualizations/3", "DELETE", nil, 200, string(msgOKjson))
	common.TestEndpoint(t, router, token, "/api/visualizations?simulationID=1", "GET", nil, 200, string(msgVisualizationsjson))

	// TODO add testing for other return codes

}
