/** Widget package, testing.
*
* @author Sonja Happ <sonja.happ@eonerc.rwth-aachen.de>
* @copyright 2014-2019, Institute for Automation of Complex Power Systems, EONERC
* @license GNU General Public License (version 3)
*
* VILLASweb-backend-go
*
* This program is free software: you can redistribute it and/or modify
* it under the terms of the GNU General Public License as published by
* the Free Software Foundation, either version 3 of the License, or
* any later version.
*
* This program is distributed in the hope that it will be useful,
* but WITHOUT ANY WARRANTY; without even the implied warranty of
* MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
* GNU General Public License for more details.
*
* You should have received a copy of the GNU General Public License
* along with this program.  If not, see <http://www.gnu.org/licenses/>.
*********************************************************************************/
package widget

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"git.rwth-aachen.de/acs/public/villas/web-backend-go/configuration"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/helper"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/dashboard"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/scenario"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/user"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/stretchr/testify/assert"
)

var router *gin.Engine

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
	SignalIDs        []int64        `json:"signalIDs,omitempty"`
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

var newWidget = WidgetRequest{
	Name:             "My label",
	Type:             "Label",
	Width:            100,
	Height:           50,
	MinWidth:         40,
	MinHeight:        80,
	X:                10,
	Y:                10,
	Z:                200,
	IsLocked:         false,
	CustomProperties: postgres.Jsonb{RawMessage: json.RawMessage(`{"textSize" : "20", "fontColor" : "#4287f5", "fontColor_opacity": 1}`)},
	SignalIDs:        []int64{},
}

func addScenarioAndDashboard(token string) (scenarioID uint, dashboardID uint) {

	// POST $newScenario
	newScenario := ScenarioRequest{
		Name:            "Scenario1",
		Running:         true,
		StartParameters: postgres.Jsonb{json.RawMessage(`{"parameter1" : "testValue1A", "parameter2" : "testValue2A", "parameter3" : 42}`)},
	}
	_, resp, _ := helper.TestEndpoint(router, token,
		"/api/scenarios", "POST", helper.KeyModels{"scenario": newScenario})

	// Read newScenario's ID from the response
	newScenarioID, _ := helper.GetResponseID(resp)

	// test POST dashboards/ $newDashboard
	newDashboard := DashboardRequest{
		Name:       "DashboardA",
		Grid:       15,
		ScenarioID: uint(newScenarioID),
	}
	_, resp, _ = helper.TestEndpoint(router, token,
		"/api/dashboards", "POST", helper.KeyModels{"dashboard": newDashboard})

	// Read newDashboard's ID from the response
	newDashboardID, _ := helper.GetResponseID(resp)

	// add the guest user to the new scenario
	_, resp, _ = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/scenarios/%v/user?username=User_C", newScenarioID), "PUT", nil)

	return uint(newScenarioID), uint(newDashboardID)
}

func TestMain(m *testing.M) {
	err := configuration.InitConfig()
	if err != nil {
		panic(m)
	}

	err = database.InitDB(configuration.GlobalConfig)
	if err != nil {
		panic(m)
	}
	defer database.DBpool.Close()

	router = gin.Default()
	api := router.Group("/api")

	user.RegisterAuthenticate(api.Group("/authenticate"))
	api.Use(user.Authentication())
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
	database.DropTables()
	database.MigrateModels()
	assert.NoError(t, helper.AddTestUsers())

	// authenticate as normal user
	token, err := helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserACredentials)
	assert.NoError(t, err)

	_, dashboardID := addScenarioAndDashboard(token)

	newWidget.DashboardID = dashboardID
	// authenticate as userB who has no access to scenario
	token, err = helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserBCredentials)
	assert.NoError(t, err)

	// try to POST the newWidget with no access to the scenario
	// should result in unprocessable entity
	code, resp, err := helper.TestEndpoint(router, token,
		"/api/widgets", "POST", helper.KeyModels{"widget": newWidget})
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

	// authenticate as normal user
	token, err = helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserACredentials)
	assert.NoError(t, err)

	// try to POST non JSON body
	// should result in bad request
	code, resp, err = helper.TestEndpoint(router, token,
		"/api/widgets", "POST", "This is no JSON")
	assert.NoError(t, err)
	assert.Equalf(t, 400, code, "Response body: \n%v\n", resp)

	// test POST widgets/ $newWidget
	code, resp, err = helper.TestEndpoint(router, token,
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

	// authenticate as userB who has no access to scenario
	token, err = helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserBCredentials)
	assert.NoError(t, err)

	// try to GET the newWidget with no access to the scenario
	// should result in unprocessable entity
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/widgets/%v", newWidgetID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)
}

func TestUpdateWidget(t *testing.T) {
	database.DropTables()
	database.MigrateModels()
	assert.NoError(t, helper.AddTestUsers())

	// authenticate as normal user
	token, err := helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserACredentials)
	assert.NoError(t, err)

	_, dashboardID := addScenarioAndDashboard(token)

	// test POST widgets/ $newWidget
	newWidget.DashboardID = dashboardID
	code, resp, err := helper.TestEndpoint(router, token,
		"/api/widgets", "POST", helper.KeyModels{"widget": newWidget})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Read newWidget's ID from the response
	newWidgetID, err := helper.GetResponseID(resp)
	assert.NoError(t, err)

	updatedWidget := WidgetRequest{
		Name:             "My slider",
		Type:             "Slider",
		Width:            400,
		Height:           50,
		MinWidth:         30,
		MinHeight:        380,
		CustomProperties: postgres.Jsonb{RawMessage: json.RawMessage(`{"default_value" : 0, "orientation" : 0, "rangeUseMinMax": false, "rangeMin" : 0, "rangeMax": 200, "rangeUseMinMax" : true, "showUnit": true, "continous_update": false, "value": "", "resizeLeftRightLock": false, "resizeTopBottomLock": true, "step": 0.1 }`)},
		SignalIDs:        []int64{},
	}

	// authenticate as userB who has no access to scenario
	token, err = helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserBCredentials)
	assert.NoError(t, err)

	// try to PUT the updatedWidget with no access to the scenario
	// should result in unprocessable entity
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/widgets/%v", newWidgetID), "PUT", helper.KeyModels{"widget": updatedWidget})
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

	// authenticate as guest user who has access to scenario
	token, err = helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.GuestCredentials)
	assert.NoError(t, err)

	// try to PUT as guest
	// should NOT work and result in unprocessable entity
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/widgets/%v", newWidgetID), "PUT", helper.KeyModels{"widget": updatedWidget})
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

	// authenticate as normal user
	token, err = helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserACredentials)
	assert.NoError(t, err)

	// try to PUT non JSON body
	// should result in bad request
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/widgets/%v", newWidgetID), "PUT", "This is no JSON")
	assert.NoError(t, err)
	assert.Equalf(t, 400, code, "Response body: \n%v\n", resp)

	// test PUT
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
	database.DropTables()
	database.MigrateModels()
	assert.NoError(t, helper.AddTestUsers())

	// authenticate as normal user
	token, err := helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserACredentials)
	assert.NoError(t, err)

	_, dashboardID := addScenarioAndDashboard(token)

	// test POST widgets/ $newWidget
	newWidget.DashboardID = dashboardID
	code, resp, err := helper.TestEndpoint(router, token,
		"/api/widgets", "POST", helper.KeyModels{"widget": newWidget})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Read newWidget's ID from the response
	newWidgetID, err := helper.GetResponseID(resp)
	assert.NoError(t, err)

	// authenticate as userB who has no access to scenario
	token, err = helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserBCredentials)
	assert.NoError(t, err)

	// try to DELETE the newWidget with no access to the scenario
	// should result in unprocessable entity
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/widgets/%v", newWidgetID), "DELETE", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

	// authenticate as normal user
	token, err = helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserACredentials)
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
	database.DropTables()
	database.MigrateModels()
	assert.NoError(t, helper.AddTestUsers())

	// authenticate as normal user
	token, err := helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserACredentials)
	assert.NoError(t, err)

	_, dashboardID := addScenarioAndDashboard(token)

	// authenticate as userB who has no access to scenario
	token, err = helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserBCredentials)
	assert.NoError(t, err)

	// try to GET all widgets of dashboard
	// should result in unprocessable entity
	code, resp, err := helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/widgets?dashboardID=%v", dashboardID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

	// authenticate as normal user
	token, err = helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserACredentials)
	assert.NoError(t, err)

	// Count the number of all the widgets returned for dashboard
	initialNumber, err := helper.LengthOfResponse(router, token,
		fmt.Sprintf("/api/widgets?dashboardID=%v", dashboardID), "GET", nil)
	assert.NoError(t, err)

	// test POST widgets/ $newWidget
	newWidget.DashboardID = dashboardID
	code, resp, err = helper.TestEndpoint(router, token,
		"/api/widgets", "POST", helper.KeyModels{"widget": newWidget})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	newWidgetB := WidgetRequest{
		Name:             "My slider",
		Type:             "Slider",
		Width:            400,
		Height:           50,
		MinWidth:         30,
		MinHeight:        380,
		X:                70,
		Y:                400,
		Z:                0,
		IsLocked:         true,
		CustomProperties: postgres.Jsonb{RawMessage: json.RawMessage(`{"default_value" : 0, "orientation" : 0, "rangeUseMinMax": false, "rangeMin" : 0, "rangeMax": 200, "rangeUseMinMax" : true, "showUnit": true, "continous_update": false, "value": "", "resizeLeftRightLock": false, "resizeTopBottomLock": true, "step": 0.1 }`)},
		DashboardID:      dashboardID,
		SignalIDs:        []int64{},
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
