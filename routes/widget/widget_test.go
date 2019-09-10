package widget

import (
	"fmt"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/database"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/helper"
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
		Name:            database.ScenarioA.Name,
		Running:         database.ScenarioA.Running,
		StartParameters: database.ScenarioA.StartParameters,
	}
	_, resp, _ := helper.TestEndpoint(router, token,
		"/api/scenarios", "POST", helper.KeyModels{"scenario": newScenario})

	// Read newScenario's ID from the response
	newScenarioID, _ := helper.GetResponseID(resp)

	// test POST dashboards/ $newDashboard
	newDashboard := DashboardRequest{
		Name:       database.DashboardA.Name,
		Grid:       database.DashboardA.Grid,
		ScenarioID: uint(newScenarioID),
	}
	_, resp, _ = helper.TestEndpoint(router, token,
		"/api/dashboards", "POST", helper.KeyModels{"dashboard": newDashboard})

	// Read newDashboard's ID from the response
	newDashboardID, _ := helper.GetResponseID(resp)

	return uint(newScenarioID), uint(newDashboardID)
}

func TestMain(m *testing.M) {

	db = database.InitDB(database.DB_TEST)
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
	database.DropTables(db)
	database.MigrateModels(db)
	assert.NoError(t, database.DBAddAdminAndUser(db))

	// authenticate as normal user
	token, err := helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserACredentials)
	assert.NoError(t, err)

	_, dashboardID := addScenarioAndDashboard(token)

	// test POST widgets/ $newWidget
	newWidget := WidgetRequest{
		Name:             database.WidgetA.Name,
		Type:             database.WidgetA.Type,
		Width:            database.WidgetA.Width,
		Height:           database.WidgetA.Height,
		MinWidth:         database.WidgetA.MinWidth,
		MinHeight:        database.WidgetA.MinHeight,
		X:                database.WidgetA.X,
		Y:                database.WidgetA.Y,
		Z:                database.WidgetA.Z,
		IsLocked:         database.WidgetA.IsLocked,
		CustomProperties: database.WidgetA.CustomProperties,
		DashboardID:      dashboardID,
	}
	code, resp, err := helper.TestEndpoint(router, token,
		"/api/widgets", "POST", helper.KeyModels{"widget": newWidget})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare POST's response with the newWidget
	err = helper.CompareResponse(resp, helper.KeyModels{"widget": newWidget})
	assert.NoError(t, err)

	// Read newWidget's ID from the response
	newWidgetID, err := helper.GetResponseID(resp)
	assert.NoError(t, err)

	// Get the newWidget
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/widgets/%v", newWidgetID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare GET's response with the newWidget
	err = helper.CompareResponse(resp, helper.KeyModels{"widget": newWidget})
	assert.NoError(t, err)

	// try to POST a malformed widget
	// Required fields are missing
	malformedNewWidget := WidgetRequest{
		Name: "ThisIsAMalformedDashboard",
	}
	// this should NOT work and return a unprocessable entity 442 status code
	code, resp, err = helper.TestEndpoint(router, token,
		"/api/widgets", "POST", helper.KeyModels{"widget": malformedNewWidget})
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)
}

func TestUpdateWidget(t *testing.T) {
	database.DropTables(db)
	database.MigrateModels(db)
	assert.NoError(t, database.DBAddAdminAndUser(db))

	// authenticate as normal user
	token, err := helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserACredentials)
	assert.NoError(t, err)

	_, dashboardID := addScenarioAndDashboard(token)

	// test POST widgets/ $newWidget
	newWidget := WidgetRequest{
		Name:             database.WidgetA.Name,
		Type:             database.WidgetA.Type,
		Width:            database.WidgetA.Width,
		Height:           database.WidgetA.Height,
		MinWidth:         database.WidgetA.MinWidth,
		MinHeight:        database.WidgetA.MinHeight,
		X:                database.WidgetA.X,
		Y:                database.WidgetA.Y,
		Z:                database.WidgetA.Z,
		IsLocked:         database.WidgetA.IsLocked,
		CustomProperties: database.WidgetA.CustomProperties,
		DashboardID:      dashboardID,
	}
	code, resp, err := helper.TestEndpoint(router, token,
		"/api/widgets", "POST", helper.KeyModels{"widget": newWidget})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Read newWidget's ID from the response
	newWidgetID, err := helper.GetResponseID(resp)
	assert.NoError(t, err)

	updatedWidget := WidgetRequest{
		Name:             database.WidgetB.Name,
		Type:             database.WidgetB.Type,
		Width:            database.WidgetB.Width,
		Height:           database.WidgetB.Height,
		MinWidth:         database.WidgetB.MinWidth,
		MinHeight:        database.WidgetB.MinHeight,
		CustomProperties: database.WidgetA.CustomProperties,
	}

	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/widgets/%v", newWidgetID), "PUT", helper.KeyModels{"widget": updatedWidget})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare PUT's response with the updatedWidget
	err = helper.CompareResponse(resp, helper.KeyModels{"widget": updatedWidget})
	assert.NoError(t, err)

	// Get the updatedWidget
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/widgets/%v", newWidgetID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare GET's response with the updatedWidget
	err = helper.CompareResponse(resp, helper.KeyModels{"widget": updatedWidget})
	assert.NoError(t, err)

	// try to update a widget that does not exist (should return not found 404 status code)
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/widgets/%v", newWidgetID+1), "PUT", helper.KeyModels{"widget": updatedWidget})
	assert.NoError(t, err)
	assert.Equalf(t, 404, code, "Response body: \n%v\n", resp)

}

func TestDeleteWidget(t *testing.T) {
	database.DropTables(db)
	database.MigrateModels(db)
	assert.NoError(t, database.DBAddAdminAndUser(db))

	// authenticate as normal user
	token, err := helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserACredentials)
	assert.NoError(t, err)

	_, dashboardID := addScenarioAndDashboard(token)

	// test POST widgets/ $newWidget
	newWidget := WidgetRequest{
		Name:             database.WidgetA.Name,
		Type:             database.WidgetA.Type,
		Width:            database.WidgetA.Width,
		Height:           database.WidgetA.Height,
		MinWidth:         database.WidgetA.MinWidth,
		MinHeight:        database.WidgetA.MinHeight,
		X:                database.WidgetA.X,
		Y:                database.WidgetA.Y,
		Z:                database.WidgetA.Z,
		IsLocked:         database.WidgetA.IsLocked,
		CustomProperties: database.WidgetA.CustomProperties,
		DashboardID:      dashboardID,
	}
	code, resp, err := helper.TestEndpoint(router, token,
		"/api/widgets", "POST", helper.KeyModels{"widget": newWidget})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Read newWidget's ID from the response
	newWidgetID, err := helper.GetResponseID(resp)
	assert.NoError(t, err)

	// Count the number of all the widgets returned for dashboard
	initialNumber, err := helper.LengthOfResponse(router, token,
		fmt.Sprintf("/api/widgets?dashboardID=%v", dashboardID), "GET", nil)
	assert.NoError(t, err)

	// Delete the added newWidget
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/widgets/%v", newWidgetID), "DELETE", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare DELETE's response with the newWidget
	err = helper.CompareResponse(resp, helper.KeyModels{"widget": newWidget})
	assert.NoError(t, err)

	// Again count the number of all the widgets returned for dashboard
	finalNumber, err := helper.LengthOfResponse(router, token,
		fmt.Sprintf("/api/widgets?dashboardID=%v", dashboardID), "GET", nil)
	assert.NoError(t, err)

	assert.Equal(t, initialNumber-1, finalNumber)
}

func TestGetAllWidgetsOfDashboard(t *testing.T) {
	database.DropTables(db)
	database.MigrateModels(db)
	assert.NoError(t, database.DBAddAdminAndUser(db))

	// authenticate as normal user
	token, err := helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserACredentials)
	assert.NoError(t, err)

	_, dashboardID := addScenarioAndDashboard(token)

	// Count the number of all the widgets returned for dashboard
	initialNumber, err := helper.LengthOfResponse(router, token,
		fmt.Sprintf("/api/widgets?dashboardID=%v", dashboardID), "GET", nil)
	assert.NoError(t, err)

	// test POST widgets/ $newWidget
	newWidgetA := WidgetRequest{
		Name:             database.WidgetA.Name,
		Type:             database.WidgetA.Type,
		Width:            database.WidgetA.Width,
		Height:           database.WidgetA.Height,
		MinWidth:         database.WidgetA.MinWidth,
		MinHeight:        database.WidgetA.MinHeight,
		X:                database.WidgetA.X,
		Y:                database.WidgetA.Y,
		Z:                database.WidgetA.Z,
		IsLocked:         database.WidgetA.IsLocked,
		CustomProperties: database.WidgetA.CustomProperties,
		DashboardID:      dashboardID,
	}
	code, resp, err := helper.TestEndpoint(router, token,
		"/api/widgets", "POST", helper.KeyModels{"widget": newWidgetA})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	newWidgetB := WidgetRequest{
		Name:             database.WidgetB.Name,
		Type:             database.WidgetB.Type,
		Width:            database.WidgetB.Width,
		Height:           database.WidgetB.Height,
		MinWidth:         database.WidgetB.MinWidth,
		MinHeight:        database.WidgetB.MinHeight,
		X:                database.WidgetB.X,
		Y:                database.WidgetB.Y,
		Z:                database.WidgetB.Z,
		IsLocked:         database.WidgetB.IsLocked,
		CustomProperties: database.WidgetB.CustomProperties,
		DashboardID:      dashboardID,
	}
	code, resp, err = helper.TestEndpoint(router, token,
		"/api/widgets", "POST", helper.KeyModels{"widget": newWidgetB})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Again count the number of all the widgets returned for dashboard
	finalNumber, err := helper.LengthOfResponse(router, token,
		fmt.Sprintf("/api/widgets?dashboardID=%v", dashboardID), "GET", nil)
	assert.NoError(t, err)

	assert.Equal(t, initialNumber+2, finalNumber)
}
