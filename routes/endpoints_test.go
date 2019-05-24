package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/file"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/model"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/simulation"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/simulator"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/user"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/visualization"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/widget"
)


var msgOK = common.ResponseMsg{
	Message: "OK.",
}


var user_A = common.UserResponse{
	Username: "User_A",
	Role: "user",
	Mail: "",
}

var user_B = common.UserResponse{
	Username: "User_B",
	Role: "user",
	Mail: "",
}

var myUsers = []common.UserResponse{
	user_A,
	user_B,
}

var msgUsers = common.ResponseMsgUsers{
	Users: myUsers,
}


var simulationA = common.SimulationResponse{
	Name: "Simulation_A",
	ID: 1,
	Running: false,
}

var simulationB = common.SimulationResponse{
	Name: "Simulation_B",
	ID: 2,
	Running: false,
}

var mySimulations = []common.SimulationResponse{
	simulationA,
	simulationB,
}

var msgSimulations = common.ResponseMsgSimulations{
	Simulations: mySimulations,
}

var msgSimulation = common.ResponseMsgSimulation{
	Simulation: simulationA,
}

// Test /simulation endpoints
func TestSimulationEndpoints(t *testing.T) {

	db := common.DummyInitDB()
	defer db.Close()
	common.DummyPopulateDB(db)


	router := gin.Default()
	api := router.Group("/api")
	simulation.RegisterSimulationEndpoints(api.Group("/simulations"))
	file.RegisterFileEndpoints(api.Group("/simulations"))
	model.RegisterModelEndpoints(api.Group("/simulations"))
	visualization.RegisterVisualizationEndpoints(api.Group("/simulations"))
	widget.RegisterWidgetEndpoints(api.Group("/simulations"))
	user.RegisterUserEndpointsForSimulation(api.Group("/simulations")) // TODO ugly structure

	user.RegisterUserEndpoints(api.Group("/users"))
	simulator.RegisterSimulatorEndpoints(api.Group("/simulators"))

	msgOKjson, err := json.Marshal(msgOK)
	if err !=nil {
		panic(err)
	}

	msgUsersjson, err := json.Marshal(msgUsers)
	if err !=nil {
		panic(err)
	}

	msgSimulationsjson, err := json.Marshal(msgSimulations)
	if err !=nil {
		panic(err)
	}

	msgSimulationjson, err := json.Marshal(msgSimulation)
	if err !=nil {
		panic(err)
	}

	// test GET simulations/
	testEndpoint(t, router, "/api/simulations/", "GET", "", 200, string(msgSimulationsjson))

	// test GET simulations/:SimulationID
	testEndpoint(t, router, "/api/simulations/1", "GET", "", 200, string(msgSimulationjson))

	// test GET simulations/:SimulationID/users
	testEndpoint(t, router, "/api/simulations/1/users", "GET", "", 200, string(msgUsersjson))

	// test DELETE simulations/:SimulationID/user/:username
	testEndpoint(t, router, "/api/simulations/1/user/User_A", "DELETE", "", 200, string(msgOKjson))

	// test PUT simulations/:SimulationID/user/:username
	testEndpoint(t, router, "/api/simulations/1/user/User_A", "PUT", "", 200, string(msgOKjson))


	// TODO add more tests

}


func testEndpoint(t *testing.T, router *gin.Engine, url string, method string, body string, expected_code int, expected_response string ) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, url, nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, expected_code, w.Code)
	fmt.Println(w.Body.String())
	assert.Equal(t, expected_response, w.Body.String())
}