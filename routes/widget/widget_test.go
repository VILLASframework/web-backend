package widget

import (
	"encoding/json"
	"testing"

	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/user"
)

// Test /widgets endpoints
func TestWidgetEndpoints(t *testing.T) {

	var token string

	var myWidgets = []common.WidgetResponse{common.WidgetA_response, common.WidgetB_response}
	var msgWidgets = common.ResponseMsgWidgets{Widgets: myWidgets}
	var msgWdg = common.ResponseMsgWidget{Widget: common.WidgetC_response}
	var msgWdgupdated = common.ResponseMsgWidget{Widget: common.WidgetCUpdated_response}

	db := common.DummyInitDB()
	defer db.Close()
	common.DummyPopulateDB(db)

	router := gin.Default()
	api := router.Group("/api")

	// All endpoints require authentication except when someone wants to
	// login (POST /authenticate)
	user.VisitorAuthenticate(api.Group("/authenticate"))

	api.Use(user.Authentication(true))

	RegisterWidgetEndpoints(api.Group("/widgets"))

	credjson, _ := json.Marshal(common.CredUser)
	msgOKjson, _ := json.Marshal(common.MsgOK)
	msgWidgetsjson, _ := json.Marshal(msgWidgets)
	msgWdgjson, _ := json.Marshal(msgWdg)
	msgWdgupdatedjson, _ := json.Marshal(msgWdgupdated)

	token = common.AuthenticateForTest(t, router, "/api/authenticate", "POST", credjson, 200)

	// test GET widgets
	common.TestEndpoint(t, router, token, "/api/widgets?dashboardID=1", "GET", nil, 200, msgWidgetsjson)

	// test POST widgets
	common.TestEndpoint(t, router, token, "/api/widgets", "POST", msgWdgjson, 200, msgOKjson)

	// test GET widgets/:widgetID to check if previous POST worked correctly
	common.TestEndpoint(t, router, token, "/api/widgets/3", "GET", nil, 200, msgWdgjson)

	// test PUT widgets/:widgetID
	common.TestEndpoint(t, router, token, "/api/widgets/3", "PUT", msgWdgupdatedjson, 200, msgOKjson)
	common.TestEndpoint(t, router, token, "/api/widgets/3", "GET", nil, 200, msgWdgupdatedjson)

	// test DELETE widgets/:widgetID
	common.TestEndpoint(t, router, token, "/api/widgets/3", "DELETE", nil, 200, msgOKjson)
	common.TestEndpoint(t, router, token, "/api/widgets?dashboardID=1", "GET", nil, 200, msgWidgetsjson)

	// TODO add testing for other return codes

}
