package scenario

import (
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
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
	ID:       2,
}

var user_B = common.UserResponse{
	Username: "User_B",
	Role:     "User",
	Mail:     "",
	ID:       3,
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

var scenarioA = common.ScenarioResponse{
	Name:    "Scenario_A",
	ID:      1,
	Running: false,
}

var scenarioB = common.ScenarioResponse{
	Name:    "Scenario_B",
	ID:      2,
	Running: false,
}

var scenarioC = common.Scenario{
	Name:            "Scenario_C",
	Running:         false,
	StartParameters: "test",
}

var scenarioC_response = common.ScenarioResponse{
	ID:          3,
	Name:        scenarioC.Name,
	Running:     scenarioC.Running,
	StartParams: scenarioC.StartParameters,
}

var myScenarios = []common.ScenarioResponse{
	scenarioA,
	scenarioB,
}

var msgScenarios = common.ResponseMsgScenarios{
	Scenarios: myScenarios,
}

var msgScenario = common.ResponseMsgScenario{
	Scenario: scenarioC_response,
}

// Test /scenarios endpoints
func TestScenarioEndpoints(t *testing.T) {

	db := common.DummyInitDB()
	defer db.Close()
	common.DummyPopulateDB(db)

	router := gin.Default()
	api := router.Group("/api")

	// All endpoints require authentication except when someone wants to
	// login (POST /authenticate)
	user.VisitorAuthenticate(api.Group("/authenticate"))

	api.Use(user.Authentication(true))

	RegisterScenarioEndpoints(api.Group("/scenarios"))

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

	msgScenariosjson, err := json.Marshal(msgScenarios)
	if err != nil {
		panic(err)
	}

	msgScenariojson, err := json.Marshal(msgScenario)
	if err != nil {
		panic(err)
	}

	scenarioCjson, err := json.Marshal(scenarioC)
	if err != nil {
		panic(err)
	}

	token = common.AuthenticateForTest(t, router, "/api/authenticate", "POST", credjson, 200)

	// test GET scenarios/
	common.TestEndpoint(t, router, token, "/api/scenarios", "GET", nil, 200, msgScenariosjson)

	// test POST scenarios/
	common.TestEndpoint(t, router, token, "/api/scenarios", "POST", scenarioCjson, 200, msgOKjson)

	// test GET scenarios/:ScenarioID
	common.TestEndpoint(t, router, token, "/api/scenarios/3", "GET", nil, 200, msgScenariojson)

	// test DELETE scenarios/:ScenarioID
	common.TestEndpoint(t, router, token, "/api/scenarios/3", "DELETE", nil, 200, msgOKjson)
	common.TestEndpoint(t, router, token, "/api/scenarios", "GET", nil, 200, msgScenariosjson)

	// test GET scenarios/:ScenarioID/users
	common.TestEndpoint(t, router, token, "/api/scenarios/1/users", "GET", nil, 200, msgUsersjson)

	// test DELETE scenarios/:ScenarioID/user
	common.TestEndpoint(t, router, token, "/api/scenarios/1/user?username=User_B", "DELETE", nil, 200, msgOKjson)
	common.TestEndpoint(t, router, token, "/api/scenarios/1/users", "GET", nil, 200, msgUserAjson)

	// test PUT scenarios/:ScenarioID/user
	common.TestEndpoint(t, router, token, "/api/scenarios/1/user?username=User_B", "PUT", nil, 200, msgOKjson)
	common.TestEndpoint(t, router, token, "/api/scenarios/1/users", "GET", nil, 200, msgUsersjson)

	// test DELETE scenarios/:ScenarioID/user for logged in user User_A
	common.TestEndpoint(t, router, token, "/api/scenarios/1/user?username=User_A", "DELETE", nil, 200, msgOKjson)

	// test if deletion of user from scenario has worked
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("GET", "/api/scenarios/1/users", nil)
	req2.Header.Add("Authorization", "Bearer "+token)
	router.ServeHTTP(w2, req2)

	assert.Equal(t, 422, w2.Code)
	fmt.Println(w2.Body.String())
	assert.Equal(t, "\"Access denied (for scenario ID).\"", w2.Body.String())

	// TODO add tests for other return codes
}
