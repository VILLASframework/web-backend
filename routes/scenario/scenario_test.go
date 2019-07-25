package scenario

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/user"
)

// Test /scenarios endpoints
func TestScenarioEndpoints(t *testing.T) {

	var token string

	var myUsers = []common.UserResponse{common.UserA_response, common.UserB_response}
	var myUserA = []common.UserResponse{common.UserA_response}
	var msgUsers = common.ResponseMsgUsers{Users: myUsers}
	var msgUserA = common.ResponseMsgUsers{Users: myUserA}
	var myScenarios = []common.ScenarioResponse{common.ScenarioA_response, common.ScenarioB_response}
	var msgScenarios = common.ResponseMsgScenarios{Scenarios: myScenarios}
	var msgScenario = common.ResponseMsgScenario{Scenario: common.ScenarioC_response}
	var msgScenarioUpdated = common.ResponseMsgScenario{Scenario: common.ScenarioCUpdated_response}

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

	credjson, _ := json.Marshal(common.CredUser)
	msgOKjson, _ := json.Marshal(common.MsgOK)
	msgScenariosjson, _ := json.Marshal(msgScenarios)
	msgScenariojson, _ := json.Marshal(msgScenario)
	msgScenarioUpdatedjson, _ := json.Marshal(msgScenarioUpdated)

	msgUsersjson, _ := json.Marshal(msgUsers)
	msgUserAjson, _ := json.Marshal(msgUserA)

	token = common.AuthenticateForTest(t, router, "/api/authenticate", "POST", credjson, 200)

	// test GET scenarios/
	common.TestEndpoint(t, router, token, "/api/scenarios", "GET", nil, 200, msgScenariosjson)

	// test POST scenarios/
	common.TestEndpoint(t, router, token, "/api/scenarios", "POST", msgScenariojson, 200, msgOKjson)

	// test GET scenarios/:ScenarioID
	common.TestEndpoint(t, router, token, "/api/scenarios/3", "GET", nil, 200, msgScenariojson)

	// test PUT scenarios/:ScenarioID
	common.TestEndpoint(t, router, token, "/api/scenarios/3", "PUT", msgScenarioUpdatedjson, 200, msgOKjson)
	common.TestEndpoint(t, router, token, "/api/scenarios/3", "GET", nil, 200, msgScenarioUpdatedjson)

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
