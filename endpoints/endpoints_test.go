package endpoints

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
)


type responseMsg struct{
	Message string `json:"message"`
}

var msgOK = responseMsg{
	Message: "OK.",
}

type user struct{
	Username string `json:"Username"`
	Role string `json:"Role"`
	Mail string `json:"Mail"`
}

type responseUsers struct {
	Users []user `json:"users"`
}

var users = []user{
	{
		Username: "User_A",
		Role: "user",
		Mail: "",
	},
	{
		Username: "User_B",
		Role: "user",
		Mail: "",
	},
}

var msgUsers = responseUsers{
	Users: users,
}

type simulation struct{
	Name string `json:"Name"`
	SimulationID uint `json:"SimulationID"`
	Running bool `json:"Running"`
}

type responseSimulations struct {
	Simulations []simulation `json:"simulations"`
}

type responseSimulation struct {
	Simulation simulation `json:"simulation"`
}

var simulationA = simulation{
	Name: "Simulation_A",
	SimulationID: 1,
	Running: false,
}

var simulationB = simulation{
	Name: "Simulation_B",
	SimulationID: 2,
	Running: false,
}

var simulations = []simulation{
	simulationA,
	simulationB,
}

var msgSimulations = responseSimulations{
	Simulations: simulations,
}

// Test /simulation endpoints
func TestSimulationEndpoints(t *testing.T) {

	db := common.DummyInitDB()
	defer db.Close()
	common.DummyPopulateDB(db)


	router := gin.Default()
	api := router.Group("/api")
	SimulationsRegister(api.Group("/simulations"))

	msgOKjson, err := json.Marshal(msgOK)
	if err !=nil {
		panic(err)
	}

	msgUsersjson, err := json.Marshal(msgUsers)
	if err !=nil {
		panic(err)
	}

	// msgSimulationsjson, err := json.Marshal(msgSimulations)
	// if err !=nil {
	// 	panic(err)
	// }

	// test GET simulations/
	var expected_response = "{\"simulations\":[{\"Name\":\"Simulation_A\",\"SimulationID\":1,\"Running\":false,\"Starting Parameters\":null},{\"Name\":\"Simulation_B\",\"SimulationID\":2,\"Running\":false,\"Starting Parameters\":null}]}"
	testEndpoint(t, router, "/api/simulations/", "GET", "", 200, expected_response)

	// test GET simulations/:SimulationID
	expected_response = "{\"simulation\":{\"Name\":\"Simulation_A\",\"SimulationID\":1,\"Running\":false,\"Starting Parameters\":null}}"
	testEndpoint(t, router, "/api/simulations/1", "GET", "", 200, expected_response)

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