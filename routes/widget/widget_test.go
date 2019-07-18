package widget

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

var wdgA = common.WidgetResponse{
	ID:               1,
	Name:             "Widget_A",
	Type:             "",
	Height:           0,
	Width:            0,
	MinHeight:        0,
	MinWidth:         0,
	X:                0,
	Y:                0,
	Z:                0,
	IsLocked:         false,
	CustomProperties: "",
	DashboardID:      1,
}

var wdgB = common.WidgetResponse{
	ID:               2,
	Name:             "Widget_B",
	Type:             "",
	Height:           0,
	Width:            0,
	MinHeight:        0,
	MinWidth:         0,
	X:                0,
	Y:                0,
	Z:                0,
	IsLocked:         false,
	CustomProperties: "",
	DashboardID:      1,
}

var wdgC = common.Widget{
	ID:               3,
	Name:             "Widget_C",
	Type:             "",
	Height:           30,
	Width:            100,
	MinHeight:        20,
	MinWidth:         50,
	X:                11,
	Y:                12,
	Z:                13,
	IsLocked:         false,
	CustomProperties: "",
	DashboardID:      1,
}

var wdgCupdated = common.Widget{
	ID:               wdgC.ID,
	Name:             "Widget_CUpdated",
	Type:             wdgC.Type,
	Height:           wdgC.Height,
	Width:            wdgC.Width,
	MinHeight:        wdgC.MinHeight,
	MinWidth:         wdgC.MinWidth,
	X:                wdgC.X,
	Y:                wdgC.Y,
	Z:                wdgC.Z,
	IsLocked:         wdgC.IsLocked,
	CustomProperties: wdgC.CustomProperties,
	DashboardID:      wdgC.DashboardID,
}

var wdgC_response = common.WidgetResponse{
	ID:               wdgC.ID,
	Name:             wdgC.Name,
	Type:             wdgC.Type,
	Height:           wdgC.Height,
	Width:            wdgC.Width,
	MinHeight:        wdgC.MinHeight,
	MinWidth:         wdgC.MinWidth,
	X:                wdgC.X,
	Y:                wdgC.Y,
	Z:                wdgC.Z,
	IsLocked:         wdgC.IsLocked,
	CustomProperties: wdgC.CustomProperties,
	DashboardID:      wdgC.DashboardID,
}

var wdgC_responseUpdated = common.WidgetResponse{
	ID:               wdgC.ID,
	Name:             "Widget_CUpdated",
	Type:             wdgC.Type,
	Height:           wdgC.Height,
	Width:            wdgC.Width,
	MinHeight:        wdgC.MinHeight,
	MinWidth:         wdgC.MinWidth,
	X:                wdgC.X,
	Y:                wdgC.Y,
	Z:                wdgC.Z,
	IsLocked:         wdgC.IsLocked,
	CustomProperties: wdgC.CustomProperties,
	DashboardID:      wdgC.DashboardID,
}

var myWidgets = []common.WidgetResponse{
	wdgA,
	wdgB,
}

var msgWidgets = common.ResponseMsgWidgets{
	Widgets: myWidgets,
}

var msgWdg = common.ResponseMsgWidget{
	Widget: wdgC_response,
}

var msgWdgupdated = common.ResponseMsgWidget{
	Widget: wdgC_responseUpdated,
}

// Test /widgets endpoints
func TestWidgetEndpoints(t *testing.T) {

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

	credjson, err := json.Marshal(cred)
	if err != nil {
		panic(err)
	}

	msgOKjson, err := json.Marshal(msgOK)
	if err != nil {
		panic(err)
	}

	msgWidgetsjson, err := json.Marshal(msgWidgets)
	if err != nil {
		panic(err)
	}

	msgWdgjson, err := json.Marshal(msgWdg)
	if err != nil {
		panic(err)
	}

	msgWdgupdatedjson, err := json.Marshal(msgWdgupdated)
	if err != nil {
		panic(err)
	}

	wdgCjson, err := json.Marshal(wdgC)
	if err != nil {
		panic(err)
	}

	wdgCupdatedjson, err := json.Marshal(wdgCupdated)
	if err != nil {
		panic(err)
	}

	token = common.AuthenticateForTest(t, router, "/api/authenticate", "POST", credjson, 200)

	// test GET widgets
	common.TestEndpoint(t, router, token, "/api/widgets?dashboardID=1", "GET", nil, 200, string(msgWidgetsjson))

	// test POST widgets
	common.TestEndpoint(t, router, token, "/api/widgets", "POST", wdgCjson, 200, string(msgOKjson))

	// test GET widgets/:widgetID to check if previous POST worked correctly
	common.TestEndpoint(t, router, token, "/api/widgets/3", "GET", nil, 200, string(msgWdgjson))

	// test PUT widgets/:widgetID
	common.TestEndpoint(t, router, token, "/api/widgets/3", "PUT", wdgCupdatedjson, 200, string(msgOKjson))
	common.TestEndpoint(t, router, token, "/api/widgets/3", "GET", nil, 200, string(msgWdgupdatedjson))

	// test DELETE widgets/:widgetID
	common.TestEndpoint(t, router, token, "/api/widgets/3", "DELETE", nil, 200, string(msgOKjson))
	common.TestEndpoint(t, router, token, "/api/widgets?dashboardID=1", "GET", nil, 200, string(msgWidgetsjson))

	// TODO add testing for other return codes

}
