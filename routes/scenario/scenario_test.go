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

	token := common.AuthenticateForTest(t, router, "/api/authenticate",
		"POST", credjson, 200)

	// test GET scenarios/
	err := common.NewTestEndpoint(router, token, "/api/scenarios",
		"GET", nil, 200, msgScenarios)
	assert.NoError(t, err)

	// test POST scenarios/
	err = common.NewTestEndpoint(router, token, "/api/scenarios",
		"POST", msgScenario, 200, common.MsgOK)
	assert.NoError(t, err)

	// test GET scenarios/:ScenarioID
	err = common.NewTestEndpoint(router, token, "/api/scenarios/3",
		"GET", nil, 200, msgScenario)
	assert.NoError(t, err)

	// test PUT scenarios/:ScenarioID
	err = common.NewTestEndpoint(router, token, "/api/scenarios/3",
		"PUT", msgScenarioUpdated, 200, common.MsgOK)
	assert.NoError(t, err)
	err = common.NewTestEndpoint(router, token, "/api/scenarios/3",
		"GET", nil, 200, msgScenarioUpdated)
	assert.NoError(t, err)

	// test DELETE scenarios/:ScenarioID
	err = common.NewTestEndpoint(router, token, "/api/scenarios/3",
		"DELETE", nil, 200, common.MsgOK)
	assert.NoError(t, err)
	err = common.NewTestEndpoint(router, token, "/api/scenarios", "GET",
		nil, 200, msgScenarios)
	assert.NoError(t, err)

	// test GET scenarios/:ScenarioID/users
	err = common.NewTestEndpoint(router, token,
		"/api/scenarios/1/users", "GET", nil,
		200, common.KeyModels{"users": []common.User{common.UserA, common.UserB}})
	assert.NoError(t, err)

	// test DELETE scenarios/:ScenarioID/user
	err = common.NewTestEndpoint(router, token,
		"/api/scenarios/1/user?username=User_B", "DELETE", nil, 200, common.MsgOK)
	assert.NoError(t, err)
	err = common.NewTestEndpoint(router, token,
		"/api/scenarios/1/users", "GET", nil,
		200, common.KeyModels{"users": []common.User{common.UserA}})
	assert.NoError(t, err)

	// test PUT scenarios/:ScenarioID/user
	err = common.NewTestEndpoint(router, token,
		"/api/scenarios/1/user?username=User_B", "PUT", nil, 200, common.MsgOK)
	assert.NoError(t, err)
	err = common.NewTestEndpoint(router, token,
		"/api/scenarios/1/users", "GET", nil,
		200, common.KeyModels{"users": []common.User{common.UserA, common.UserB}})
	assert.NoError(t, err)

	// test DELETE scenarios/:ScenarioID/user for logged in user User_A
	err = common.NewTestEndpoint(router, token,
		"/api/scenarios/1/user?username=User_A", "DELETE", nil, 200, common.MsgOK)
	assert.NoError(t, err)

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
