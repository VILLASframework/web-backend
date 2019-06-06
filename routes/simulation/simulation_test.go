package simulation

import (
	"encoding/json"
	"testing"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/user"

	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
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

var user_A = common.UserResponse{
	Username: "User_A",
	Role:     "User",
	Mail:     "",
}

var user_B = common.UserResponse{
	Username: "User_B",
	Role:     "User",
	Mail:     "",
}

var myUsers = []common.UserResponse{
	user_A,
	user_B,
}

var myUserA = []common.UserResponse{
	user_A,
}

var msgUsers = common.ResponseMsgUsers{
	Users: myUsers,
}

var msgUserA = common.ResponseMsgUsers{
	Users: myUserA,
}

var simulationA = common.SimulationResponse{
	Name:    "Simulation_A",
	ID:      1,
	Running: false,
}

var simulationB = common.SimulationResponse{
	Name:    "Simulation_B",
	ID:      2,
	Running: false,
}

var simulationC = common.Simulation{
	Name:            "Simulation_C",
	Running:         false,
	StartParameters: "test",
}

var simulationC_response = common.SimulationResponse{
	ID:          3,
	Name:        simulationC.Name,
	Running:     simulationC.Running,
	StartParams: simulationC.StartParameters,
}

var mySimulations = []common.SimulationResponse{
	simulationA,
	simulationB,
}

var msgSimulations = common.ResponseMsgSimulations{
	Simulations: mySimulations,
}

var msgSimulation = common.ResponseMsgSimulation{
	Simulation: simulationC_response,
}

// Test /simulation endpoints
func TestSimulationEndpoints(t *testing.T) {

	db := common.DummyInitDB()
	defer db.Close()
	common.DummyPopulateDB(db)

	router := gin.Default()
	api := router.Group("/api")

	// All endpoints require authentication except when someone wants to
	// login (POST /authenticate)
	user.VisitorAuthenticate(api.Group("/authenticate"))

	api.Use(user.Authentication(true))

	RegisterSimulationEndpoints(api.Group("/simulations"))

	credjson, err := json.Marshal(cred)

	msgOKjson, err := json.Marshal(msgOK)
	if err != nil {
		panic(err)
	}

	msgUsersjson, err := json.Marshal(msgUsers)
	if err != nil {
		panic(err)
	}

	msgUserAjson, err := json.Marshal(msgUserA)
	if err != nil {
		panic(err)
	}

	msgSimulationsjson, err := json.Marshal(msgSimulations)
	if err != nil {
		panic(err)
	}

	msgSimulationjson, err := json.Marshal(msgSimulation)
	if err != nil {
		panic(err)
	}

	simulationCjson, err := json.Marshal(simulationC)
	if err != nil {
		panic(err)
	}

	token = common.AuthenticateForTest(t, router, "/api/authenticate", "POST", credjson, 200)

	// test GET simulations/
	common.TestEndpoint(t, router, token, "/api/simulations", "GET", nil, 200, string(msgSimulationsjson))

	// test POST simulations/
	common.TestEndpoint(t, router, token, "/api/simulations", "POST", simulationCjson, 200, string(msgOKjson))

	// test GET simulations/:SimulationID
	common.TestEndpoint(t, router, token, "/api/simulations/3", "GET", nil, 200, string(msgSimulationjson))

	// test DELETE simulations/:SimulationID
	common.TestEndpoint(t, router, token, "/api/simulations/3", "DELETE", nil, 200, string(msgOKjson))
	common.TestEndpoint(t, router, token, "/api/simulations", "GET", nil, 200, string(msgSimulationsjson))

	// test GET simulations/:SimulationID/users
	common.TestEndpoint(t, router, token, "/api/simulations/1/users", "GET", nil, 200, string(msgUsersjson))

	// test DELETE simulations/:SimulationID/user
	common.TestEndpoint(t, router, token, "/api/simulations/1/user?username=User_B", "DELETE", nil, 200, string(msgOKjson))
	common.TestEndpoint(t, router, token, "/api/simulations/1/users", "GET", nil, 200, string(msgUserAjson))

	// test PUT simulations/:SimulationID/user
	common.TestEndpoint(t, router, token, "/api/simulations/1/user?username=User_B", "PUT", nil, 200, string(msgOKjson))
	common.TestEndpoint(t, router, token, "/api/simulations/1/users", "GET", nil, 200, string(msgUsersjson))

	// test DELETE simulations/:SimulationID/user for logged in user User_A
	common.TestEndpoint(t, router, token, "/api/simulations/1/user?username=User_A", "DELETE", nil, 200, string(msgOKjson))
	common.TestEndpoint(t, router, token, "/api/simulations/1/users", "GET", nil, 422, "\"Access denied (for simulation ID).\"")

	// TODO add tests for other return codes
}
