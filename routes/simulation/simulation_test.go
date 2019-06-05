package simulation

import (
	"bytes"
	"encoding/json"
	"fmt"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/user"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

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

	authenticate(t, router, "/api/authenticate", "POST", credjson, 200)

	// test GET simulations/
	testEndpoint(t, router, "/api/simulations/", "GET", nil, 200, string(msgSimulationsjson))

	// test POST simulations/
	testEndpoint(t, router, "/api/simulations/", "POST", simulationCjson, 200, string(msgOKjson))

	// test GET simulations/:SimulationID
	testEndpoint(t, router, "/api/simulations/3", "GET", nil, 200, string(msgSimulationjson))

	// test DELETE simulations/:SimulationID
	testEndpoint(t, router, "/api/simulations/3", "DELETE", nil, 200, string(msgOKjson))
	testEndpoint(t, router, "/api/simulations/", "GET", nil, 200, string(msgSimulationsjson))

	// test GET simulations/:SimulationID/users
	testEndpoint(t, router, "/api/simulations/1/users", "GET", nil, 200, string(msgUsersjson))

	// test DELETE simulations/:SimulationID/user
	testEndpoint(t, router, "/api/simulations/1/user?username=User_B", "DELETE", nil, 200, string(msgOKjson))
	testEndpoint(t, router, "/api/simulations/1/users", "GET", nil, 200, string(msgUserAjson))

	// test PUT simulations/:SimulationID/user
	testEndpoint(t, router, "/api/simulations/1/user?username=User_B", "PUT", nil, 200, string(msgOKjson))
	testEndpoint(t, router, "/api/simulations/1/users", "GET", nil, 200, string(msgUsersjson))

	// test DELETE simulations/:SimulationID/user for logged in user User_A
	testEndpoint(t, router, "/api/simulations/1/user?username=User_A", "DELETE", nil, 200, string(msgOKjson))
	testEndpoint(t, router, "/api/simulations/1/users", "GET", nil, 422, "\"Access denied (entity ID).\"")

	// TODO add tests for other return codes
}

func testEndpoint(t *testing.T, router *gin.Engine, url string, method string, body []byte, expected_code int, expected_response string) {
	w := httptest.NewRecorder()

	if body != nil {
		req, _ := http.NewRequest(method, url, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)
	} else {
		req, _ := http.NewRequest(method, url, nil)
		req.Header.Add("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)
	}

	assert.Equal(t, expected_code, w.Code)
	fmt.Println(w.Body.String())
	assert.Equal(t, expected_response, w.Body.String())
}

func authenticate(t *testing.T, router *gin.Engine, url string, method string, body []byte, expected_code int) {
	w := httptest.NewRecorder()

	req, _ := http.NewRequest(method, url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, expected_code, w.Code)

	var body_data map[string]interface{}

	err := json.Unmarshal([]byte(w.Body.String()), &body_data)
	if err != nil {
		panic(err)
	}

	success := body_data["success"].(bool)
	if !success {
		panic(-1)
	}
	// save token to global variable
	token = body_data["token"].(string)

	fmt.Println(w.Body.String())
}
