package dashboard

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

var dabA = common.DashboardResponse{
	ID:         1,
	Name:       "Dashboard_A",
	Grid:       15,
	ScenarioID: 1,
}

var dabB = common.DashboardResponse{
	ID:         2,
	Name:       "Dashboard_B",
	Grid:       15,
	ScenarioID: 1,
}

var dabC = common.Dashboard{
	ID:         3,
	Name:       "Dashboard_C",
	Grid:       99,
	ScenarioID: 1,
}

var dabCupdated = common.Dashboard{
	ID:         dabC.ID,
	Name:       "Dashboard_CUpdated",
	ScenarioID: dabC.ScenarioID,
	Grid:       dabC.Grid,
}

var dabC_response = common.DashboardResponse{
	ID:         dabC.ID,
	Name:       dabC.Name,
	Grid:       dabC.Grid,
	ScenarioID: dabC.ScenarioID,
}

var dabC_responseUpdated = common.DashboardResponse{
	ID:         dabCupdated.ID,
	Name:       dabCupdated.Name,
	Grid:       dabCupdated.Grid,
	ScenarioID: dabCupdated.ScenarioID,
}

var myDashboards = []common.DashboardResponse{
	dabA,
	dabB,
}

var msgDashboards = common.ResponseMsgDashboards{
	Dashboards: myDashboards,
}

var msgDab = common.ResponseMsgDashboard{
	Dashboard: dabC_response,
}

var msgDabupdated = common.ResponseMsgDashboard{
	Dashboard: dabC_responseUpdated,
}

// Test /dashboards endpoints
func TestEndpoints(t *testing.T) {

	db := common.DummyInitDB()
	defer db.Close()
	common.DummyPopulateDB(db)

	router := gin.Default()
	api := router.Group("/api")

	// All endpoints require authentication except when someone wants to
	// login (POST /authenticate)
	user.VisitorAuthenticate(api.Group("/authenticate"))

	api.Use(user.Authentication(true))

	RegisterDashboardEndpoints(api.Group("/dashboards"))

	credjson, err := json.Marshal(cred)
	if err != nil {
		panic(err)
	}

	msgOKjson, err := json.Marshal(msgOK)
	if err != nil {
		panic(err)
	}

	msgDashboardsjson, err := json.Marshal(msgDashboards)
	if err != nil {
		panic(err)
	}

	msgDabjson, err := json.Marshal(msgDab)
	if err != nil {
		panic(err)
	}

	msgDabupdatedjson, err := json.Marshal(msgDabupdated)
	if err != nil {
		panic(err)
	}

	dabCjson, err := json.Marshal(dabC)
	if err != nil {
		panic(err)
	}

	dabCupdatedjson, err := json.Marshal(dabCupdated)
	if err != nil {
		panic(err)
	}

	token = common.AuthenticateForTest(t, router, "/api/authenticate", "POST", credjson, 200)

	// test GET dashboards
	common.TestEndpoint(t, router, token, "/api/dashboards?scenarioID=1", "GET", nil, 200, string(msgDashboardsjson))

	// test POST dashboards
	common.TestEndpoint(t, router, token, "/api/dashboards", "POST", dabCjson, 200, string(msgOKjson))

	// test GET dashboards/:dashboardID to check if previous POST worked correctly
	common.TestEndpoint(t, router, token, "/api/dashboards/3", "GET", nil, 200, string(msgDabjson))

	// test PUT dashboards/:dashboardID
	common.TestEndpoint(t, router, token, "/api/dashboards/3", "PUT", dabCupdatedjson, 200, string(msgOKjson))
	common.TestEndpoint(t, router, token, "/api/dashboards/3", "GET", nil, 200, string(msgDabupdatedjson))

	// test DELETE dashboards/:dashboardID
	common.TestEndpoint(t, router, token, "/api/dashboards/3", "DELETE", nil, 200, string(msgOKjson))
	common.TestEndpoint(t, router, token, "/api/dashboards?scenarioID=1", "GET", nil, 200, string(msgDashboardsjson))

	// TODO add testing for other return codes

}
