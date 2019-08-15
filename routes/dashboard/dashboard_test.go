package dashboard

import (
	"encoding/json"
	"testing"

	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/user"
)

// Test /dashboards endpoints
func TestEndpoints(t *testing.T) {

	var token string

	var myDashboards = []common.DashboardResponse{common.DashboardA_response, common.DashboardB_response}
	var msgDashboards = common.ResponseMsgDashboards{Dashboards: myDashboards}
	var msgDab = common.ResponseMsgDashboard{Dashboard: common.DashboardC_response}
	var msgDabupdated = common.ResponseMsgDashboard{Dashboard: common.DashboardCUpdated_response}

	db := common.DummyInitDB()
	defer db.Close()
	common.DummyPopulateDB(db)

	router := gin.Default()
	api := router.Group("/api")

	// All endpoints require authentication except when someone wants to
	// login (POST /authenticate)
	user.RegisterAuthenticate(api.Group("/authenticate"))

	api.Use(user.Authentication(true))

	RegisterDashboardEndpoints(api.Group("/dashboards"))

	credjson, _ := json.Marshal(common.CredUser)
	msgOKjson, _ := json.Marshal(common.MsgOK)
	msgDashboardsjson, _ := json.Marshal(msgDashboards)
	msgDabjson, _ := json.Marshal(msgDab)
	msgDabupdatedjson, _ := json.Marshal(msgDabupdated)

	token = common.AuthenticateForTest(t, router, "/api/authenticate", "POST", credjson, 200)

	// test GET dashboards
	common.TestEndpoint(t, router, token, "/api/dashboards?scenarioID=1", "GET", nil, 200, msgDashboardsjson)

	// test POST dashboards
	common.TestEndpoint(t, router, token, "/api/dashboards", "POST", msgDabjson, 200, msgOKjson)

	// test GET dashboards/:dashboardID to check if previous POST worked correctly
	common.TestEndpoint(t, router, token, "/api/dashboards/3", "GET", nil, 200, msgDabjson)

	// test PUT dashboards/:dashboardID
	common.TestEndpoint(t, router, token, "/api/dashboards/3", "PUT", msgDabupdatedjson, 200, msgOKjson)
	common.TestEndpoint(t, router, token, "/api/dashboards/3", "GET", nil, 200, msgDabupdatedjson)

	// test DELETE dashboards/:dashboardID
	common.TestEndpoint(t, router, token, "/api/dashboards/3", "DELETE", nil, 200, msgOKjson)
	common.TestEndpoint(t, router, token, "/api/dashboards?scenarioID=1", "GET", nil, 200, msgDashboardsjson)

	// TODO add testing for other return codes

}
