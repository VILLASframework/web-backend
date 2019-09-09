package widget

import (
	"fmt"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/dashboard"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/scenario"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/user"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var router *gin.Engine
var db *gorm.DB

type WidgetRequest struct {
	Name             string         `json:"name,omitempty"`
	Type             string         `json:"type,omitempty"`
	Width            uint           `json:"width,omitempty"`
	Height           uint           `json:"height,omitempty"`
	MinWidth         uint           `json:"minWidth,omitempty"`
	MinHeight        uint           `json:"minHeight,omitempty"`
	X                int            `json:"x,omitempty"`
	Y                int            `json:"y,omitempty"`
	Z                int            `json:"z,omitempty"`
	DashboardID      uint           `json:"dashboardID,omitempty"`
	IsLocked         bool           `json:"isLocked,omitempty"`
	CustomProperties postgres.Jsonb `json:"customProperties,omitempty"`
}

type DashboardRequest struct {
	Name       string `json:"name,omitempty"`
	Grid       int    `json:"grid,omitempty"`
	ScenarioID uint   `json:"scenarioID,omitempty"`
}

type ScenarioRequest struct {
	Name            string         `json:"name,omitempty"`
	Running         bool           `json:"running,omitempty"`
	StartParameters postgres.Jsonb `json:"startParameters,omitempty"`
}

func addScenarioAndDashboard(token string) (scenarioID uint, dashboardID uint) {

	// POST $newScenario
	newScenario := ScenarioRequest{
		Name:            common.ScenarioA.Name,
		Running:         common.ScenarioA.Running,
		StartParameters: common.ScenarioA.StartParameters,
	}
	_, resp, _ := common.TestEndpoint(router, token,
		"/api/scenarios", "POST", common.KeyModels{"scenario": newScenario})

	// Read newScenario's ID from the response
	newScenarioID, _ := common.GetResponseID(resp)

	// test POST dashboards/ $newDashboard
	newDashboard := DashboardRequest{
		Name:       common.DashboardA.Name,
		Grid:       common.DashboardA.Grid,
		ScenarioID: uint(newScenarioID),
	}
	_, resp, _ = common.TestEndpoint(router, token,
		"/api/dashboards", "POST", common.KeyModels{"dashboard": newDashboard})

	// Read newDashboard's ID from the response
	newDashboardID, _ := common.GetResponseID(resp)

	return uint(newScenarioID), uint(newDashboardID)
}

func TestMain(m *testing.M) {

	db = common.InitDB(common.DB_TEST)
	defer db.Close()

	router = gin.Default()
	api := router.Group("/api")

	user.RegisterAuthenticate(api.Group("/authenticate"))
	api.Use(user.Authentication(true))
	RegisterWidgetEndpoints(api.Group("/widgets"))
	// scenario endpoints required here to first add a scenario to the DB
	// that can be associated with a new dashboard
	scenario.RegisterScenarioEndpoints(api.Group("/scenarios"))
	// dashboard endpoints required here to first add a dashboard to the DB
	// that can be associated with a new widget
	dashboard.RegisterDashboardEndpoints(api.Group("/dashboards"))

	os.Exit(m.Run())
}

func TestAddWidget(t *testing.T) {
	common.DropTables(db)
	common.MigrateModels(db)
	common.DBAddAdminAndUser(db)

	// authenticate as normal user
	token, err := common.AuthenticateForTest(router,
		"/api/authenticate", "POST", common.UserACredentials)
	assert.NoError(t, err)

	_, dashboardID := addScenarioAndDashboard(token)

	// test POST widgets/ $newWidget
	newWidget := WidgetRequest{
		Name:             common.WidgetA.Name,
		Type:             common.WidgetA.Type,
		Width:            common.WidgetA.Width,
		Height:           common.WidgetA.Height,
		MinWidth:         common.WidgetA.MinWidth,
		MinHeight:        common.WidgetA.MinHeight,
		X:                common.WidgetA.X,
		Y:                common.WidgetA.Y,
		Z:                common.WidgetA.Z,
		IsLocked:         common.WidgetA.IsLocked,
		CustomProperties: common.WidgetA.CustomProperties,
		DashboardID:      dashboardID,
	}
	code, resp, err := common.TestEndpoint(router, token,
		"/api/widgets", "POST", common.KeyModels{"widget": newWidget})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare POST's response with the newWidget
	err = common.CompareResponse(resp, common.KeyModels{"widget": newWidget})
	assert.NoError(t, err)

	// Read newWidget's ID from the response
	newWidgetID, err := common.GetResponseID(resp)
	assert.NoError(t, err)

	// Get the newWidget
	code, resp, err = common.TestEndpoint(router, token,
		fmt.Sprintf("/api/widgets/%v", newWidgetID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare GET's response with the newWidget
	err = common.CompareResponse(resp, common.KeyModels{"widget": newWidget})
	assert.NoError(t, err)

	// try to POST a malformed widget
	// Required fields are missing
	malformedNewWidget := WidgetRequest{
		Name: "ThisIsAMalformedDashboard",
	}
	// this should NOT work and return a unprocessable entity 442 status code
	code, resp, err = common.TestEndpoint(router, token,
		"/api/widgets", "POST", common.KeyModels{"widget": malformedNewWidget})
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)
}

func TestUpdateWidget(t *testing.T) {
	common.DropTables(db)
	common.MigrateModels(db)
	common.DBAddAdminAndUser(db)

	// authenticate as normal user
	token, err := common.AuthenticateForTest(router,
		"/api/authenticate", "POST", common.UserACredentials)
	assert.NoError(t, err)

	_, dashboardID := addScenarioAndDashboard(token)

	// test POST widgets/ $newWidget
	newWidget := WidgetRequest{
		Name:             common.WidgetA.Name,
		Type:             common.WidgetA.Type,
		Width:            common.WidgetA.Width,
		Height:           common.WidgetA.Height,
		MinWidth:         common.WidgetA.MinWidth,
		MinHeight:        common.WidgetA.MinHeight,
		X:                common.WidgetA.X,
		Y:                common.WidgetA.Y,
		Z:                common.WidgetA.Z,
		IsLocked:         common.WidgetA.IsLocked,
		CustomProperties: common.WidgetA.CustomProperties,
		DashboardID:      dashboardID,
	}
	code, resp, err := common.TestEndpoint(router, token,
		"/api/widgets", "POST", common.KeyModels{"widget": newWidget})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Read newWidget's ID from the response
	newWidgetID, err := common.GetResponseID(resp)
	assert.NoError(t, err)

	updatedWidget := WidgetRequest{
		Name:             common.WidgetB.Name,
		Type:             common.WidgetB.Type,
		Width:            common.WidgetB.Width,
		Height:           common.WidgetB.Height,
		MinWidth:         common.WidgetB.MinWidth,
		MinHeight:        common.WidgetB.MinHeight,
		CustomProperties: common.WidgetA.CustomProperties,
	}

	code, resp, err = common.TestEndpoint(router, token,
		fmt.Sprintf("/api/widgets/%v", newWidgetID), "PUT", common.KeyModels{"widget": updatedWidget})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare PUT's response with the updatedWidget
	err = common.CompareResponse(resp, common.KeyModels{"widget": updatedWidget})
	assert.NoError(t, err)

	// Get the updatedWidget
	code, resp, err = common.TestEndpoint(router, token,
		fmt.Sprintf("/api/widgets/%v", newWidgetID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare GET's response with the updatedWidget
	err = common.CompareResponse(resp, common.KeyModels{"widget": updatedWidget})
	assert.NoError(t, err)

	// try to update a widget that does not exist (should return not found 404 status code)
	code, resp, err = common.TestEndpoint(router, token,
		fmt.Sprintf("/api/widgets/%v", newWidgetID+1), "PUT", common.KeyModels{"widget": updatedWidget})
	assert.NoError(t, err)
	assert.Equalf(t, 404, code, "Response body: \n%v\n", resp)

}

func TestDeleteWidget(t *testing.T) {
	common.DropTables(db)
	common.MigrateModels(db)
	common.DBAddAdminAndUser(db)

	// authenticate as normal user
	token, err := common.AuthenticateForTest(router,
		"/api/authenticate", "POST", common.UserACredentials)
	assert.NoError(t, err)

	_, dashboardID := addScenarioAndDashboard(token)

	// test POST widgets/ $newWidget
	newWidget := WidgetRequest{
		Name:             common.WidgetA.Name,
		Type:             common.WidgetA.Type,
		Width:            common.WidgetA.Width,
		Height:           common.WidgetA.Height,
		MinWidth:         common.WidgetA.MinWidth,
		MinHeight:        common.WidgetA.MinHeight,
		X:                common.WidgetA.X,
		Y:                common.WidgetA.Y,
		Z:                common.WidgetA.Z,
		IsLocked:         common.WidgetA.IsLocked,
		CustomProperties: common.WidgetA.CustomProperties,
		DashboardID:      dashboardID,
	}
	code, resp, err := common.TestEndpoint(router, token,
		"/api/widgets", "POST", common.KeyModels{"widget": newWidget})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Read newWidget's ID from the response
	newWidgetID, err := common.GetResponseID(resp)
	assert.NoError(t, err)

	// Count the number of all the widgets returned for dashboard
	initialNumber, err := common.LengthOfResponse(router, token,
		fmt.Sprintf("/api/widgets?dashboardID=%v", dashboardID), "GET", nil)
	assert.NoError(t, err)

	// Delete the added newWidget
	code, resp, err = common.TestEndpoint(router, token,
		fmt.Sprintf("/api/widgets/%v", newWidgetID), "DELETE", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare DELETE's response with the newWidget
	err = common.CompareResponse(resp, common.KeyModels{"widget": newWidget})
	assert.NoError(t, err)

	// Again count the number of all the widgets returned for dashboard
	finalNumber, err := common.LengthOfResponse(router, token,
		fmt.Sprintf("/api/widgets?dashboardID=%v", dashboardID), "GET", nil)
	assert.NoError(t, err)

	assert.Equal(t, initialNumber-1, finalNumber)
}

func TestGetAllWidgetsOfDashboard(t *testing.T) {
	common.DropTables(db)
	common.MigrateModels(db)
	common.DBAddAdminAndUser(db)

	// authenticate as normal user
	token, err := common.AuthenticateForTest(router,
		"/api/authenticate", "POST", common.UserACredentials)
	assert.NoError(t, err)

	_, dashboardID := addScenarioAndDashboard(token)

	// Count the number of all the widgets returned for dashboard
	initialNumber, err := common.LengthOfResponse(router, token,
		fmt.Sprintf("/api/widgets?dashboardID=%v", dashboardID), "GET", nil)
	assert.NoError(t, err)

	// test POST widgets/ $newWidget
	newWidgetA := WidgetRequest{
		Name:             common.WidgetA.Name,
		Type:             common.WidgetA.Type,
		Width:            common.WidgetA.Width,
		Height:           common.WidgetA.Height,
		MinWidth:         common.WidgetA.MinWidth,
		MinHeight:        common.WidgetA.MinHeight,
		X:                common.WidgetA.X,
		Y:                common.WidgetA.Y,
		Z:                common.WidgetA.Z,
		IsLocked:         common.WidgetA.IsLocked,
		CustomProperties: common.WidgetA.CustomProperties,
		DashboardID:      dashboardID,
	}
	code, resp, err := common.TestEndpoint(router, token,
		"/api/widgets", "POST", common.KeyModels{"widget": newWidgetA})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	newWidgetB := WidgetRequest{
		Name:             common.WidgetB.Name,
		Type:             common.WidgetB.Type,
		Width:            common.WidgetB.Width,
		Height:           common.WidgetB.Height,
		MinWidth:         common.WidgetB.MinWidth,
		MinHeight:        common.WidgetB.MinHeight,
		X:                common.WidgetB.X,
		Y:                common.WidgetB.Y,
		Z:                common.WidgetB.Z,
		IsLocked:         common.WidgetB.IsLocked,
		CustomProperties: common.WidgetB.CustomProperties,
		DashboardID:      dashboardID,
	}
	code, resp, err = common.TestEndpoint(router, token,
		"/api/widgets", "POST", common.KeyModels{"widget": newWidgetB})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Again count the number of all the widgets returned for dashboard
	finalNumber, err := common.LengthOfResponse(router, token,
		fmt.Sprintf("/api/widgets?dashboardID=%v", dashboardID), "GET", nil)
	assert.NoError(t, err)

	assert.Equal(t, initialNumber+2, finalNumber)
}
