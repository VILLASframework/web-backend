/** Dashboard package, testing.
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
package dashboard

import (
	"fmt"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/configuration"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/helper"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/scenario"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/user"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

var router *gin.Engine
var db *gorm.DB

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

func addScenario(token string) (scenarioID uint) {

	// POST $newScenario
	newScenario := ScenarioRequest{
		Name:            database.ScenarioA.Name,
		Running:         database.ScenarioA.Running,
		StartParameters: database.ScenarioA.StartParameters,
	}
	_, resp, err := helper.TestEndpoint(router, token,
		"/api/scenarios", "POST", helper.KeyModels{"scenario": newScenario})
	if err != nil {
		log.Panic("The following error happend on POSTing a scenario: ", err.Error())
	}

	// Read newScenario's ID from the response
	newScenarioID, _ := helper.GetResponseID(resp)

	// add the guest user to the new scenario
	_, resp, _ = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/scenarios/%v/user?username=User_C", newScenarioID), "PUT", nil)

	return uint(newScenarioID)
}

func TestMain(m *testing.M) {
	err := configuration.InitConfig()
	if err != nil {
		panic(m)
	}
	db, err = database.InitDB(configuration.GolbalConfig)
	if err != nil {
		panic(m)
	}
	defer db.Close()

	router = gin.Default()
	api := router.Group("/api")

	user.RegisterAuthenticate(api.Group("/authenticate"))
	api.Use(user.Authentication(true))
	RegisterDashboardEndpoints(api.Group("/dashboards"))
	// scenario endpoints required here to first add a scenario to the DB
	// that can be associated with a new dashboard
	scenario.RegisterScenarioEndpoints(api.Group("/scenarios"))

	os.Exit(m.Run())
}

func TestAddDashboard(t *testing.T) {
	database.DropTables(db)
	database.MigrateModels(db)
	assert.NoError(t, database.DBAddAdminAndUserAndGuest(db))

	// authenticate as normal user
	token, err := helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserACredentials)
	assert.NoError(t, err)

	scenarioID := addScenario(token)

	// test POST dashboards/ $newDashboard
	newDashboard := DashboardRequest{
		Name:       database.DashboardA.Name,
		Grid:       database.DashboardA.Grid,
		ScenarioID: scenarioID,
	}
	code, resp, err := helper.TestEndpoint(router, token,
		"/api/dashboards", "POST", helper.KeyModels{"dashboard": newDashboard})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare POST's response with the newDashboard
	err = helper.CompareResponse(resp, helper.KeyModels{"dashboard": newDashboard})
	assert.NoError(t, err)

	// Read newDashboard's ID from the response
	newDashboardID, err := helper.GetResponseID(resp)
	assert.NoError(t, err)

	// Get the newDashboard
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/dashboards/%v", newDashboardID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare GET's response with the newDashboard
	err = helper.CompareResponse(resp, helper.KeyModels{"dashboard": newDashboard})
	assert.NoError(t, err)

	// try to POST a malformed dashboard
	// Required fields are missing (validation should fail)
	malformedNewDashboard := DashboardRequest{
		Name: "ThisIsAMalformedDashboard",
	}
	// this should NOT work and return a unprocessable entity 442 status code
	code, resp, err = helper.TestEndpoint(router, token,
		"/api/dashboards", "POST", helper.KeyModels{"dashboard": malformedNewDashboard})
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

	// this should NOT work and return a bad request 400 status code
	code, resp, err = helper.TestEndpoint(router, token,
		"/api/dashboards", "POST", "This is a test using plain text as body")
	assert.NoError(t, err)
	assert.Equalf(t, 400, code, "Response body: \n%v\n", resp)

	// try to add a dashboard to a scenario that does not exist
	// should return not found error
	newDashboard.ScenarioID = scenarioID + 1
	code, resp, err = helper.TestEndpoint(router, token,
		"/api/dashboards", "POST", helper.KeyModels{"dashboard": newDashboard})
	assert.NoError(t, err)
	assert.Equalf(t, 404, code, "Response body: \n%v\n", resp)

	// try to get dashboard as a user that is not in the scenario (userB)
	token, err = helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserBCredentials)
	assert.NoError(t, err)

	// this should fail with unprocessable entity
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/dashboards/%v", newDashboardID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

	// try to add a dashboard to a scenario to which the user has no access
	// this should give an unprocessable entity error
	newDashboard.ScenarioID = scenarioID
	code, resp, err = helper.TestEndpoint(router, token,
		"/api/dashboards", "POST", helper.KeyModels{"dashboard": newDashboard})
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)
}

func TestUpdateDashboard(t *testing.T) {
	database.DropTables(db)
	database.MigrateModels(db)
	assert.NoError(t, database.DBAddAdminAndUserAndGuest(db))

	// authenticate as normal user
	token, err := helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserACredentials)
	assert.NoError(t, err)

	scenarioID := addScenario(token)

	// test POST dashboards/ $newDashboard
	newDashboard := DashboardRequest{
		Name:       database.DashboardA.Name,
		Grid:       database.DashboardA.Grid,
		ScenarioID: scenarioID,
	}
	code, resp, err := helper.TestEndpoint(router, token,
		"/api/dashboards", "POST", helper.KeyModels{"dashboard": newDashboard})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Read newDashboard's ID from the response
	newDashboardID, err := helper.GetResponseID(resp)
	assert.NoError(t, err)

	updatedDashboard := DashboardRequest{
		Name: database.DashboardB.Name,
		Grid: database.DashboardB.Grid,
	}

	// authenticate as guest user
	token, err = helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.GuestCredentials)
	assert.NoError(t, err)

	// try to update a dashboard as guest
	// should return an unprocessable entity error
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/dashboards/%v", newDashboardID), "PUT", helper.KeyModels{"dashboard": updatedDashboard})
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

	// authenticate as normal user
	token, err = helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserACredentials)
	assert.NoError(t, err)

	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/dashboards/%v", newDashboardID), "PUT", helper.KeyModels{"dashboard": updatedDashboard})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare PUT's response with the updatedDashboard
	err = helper.CompareResponse(resp, helper.KeyModels{"dashboard": updatedDashboard})
	assert.NoError(t, err)

	// Get the updatedDashboard
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/dashboards/%v", newDashboardID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare GET's response with the updatedDashboard
	err = helper.CompareResponse(resp, helper.KeyModels{"dashboard": updatedDashboard})
	assert.NoError(t, err)

	// try to update a dashboard that does not exist (should return not found 404 status code)
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/dashboards/%v", newDashboardID+1), "PUT", helper.KeyModels{"dashboard": updatedDashboard})
	assert.NoError(t, err)
	assert.Equalf(t, 404, code, "Response body: \n%v\n", resp)

	// try to update with a malformed body, should return a bad request error
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/dashboards/%v", newDashboardID), "PUT", "This is the body of a malformed update request.")
	assert.NoError(t, err)
	assert.Equalf(t, 400, code, "Response body: \n%v\n", resp)
}

func TestDeleteDashboard(t *testing.T) {
	database.DropTables(db)
	database.MigrateModels(db)
	assert.NoError(t, database.DBAddAdminAndUserAndGuest(db))

	// authenticate as normal user
	token, err := helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserACredentials)
	assert.NoError(t, err)

	scenarioID := addScenario(token)

	// test POST dashboards/ $newDashboard
	newDashboard := DashboardRequest{
		Name:       database.DashboardA.Name,
		Grid:       database.DashboardA.Grid,
		ScenarioID: scenarioID,
	}
	code, resp, err := helper.TestEndpoint(router, token,
		"/api/dashboards", "POST", helper.KeyModels{"dashboard": newDashboard})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Read newDashboard's ID from the response
	newDashboardID, err := helper.GetResponseID(resp)
	assert.NoError(t, err)

	// try to delete a dashboard from a scenario to which the user has no access
	token, err = helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserBCredentials)
	assert.NoError(t, err)

	// this should fail with unprocessable entity
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/dashboards/%v", newDashboardID), "DELETE", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

	// authenticate as normal user
	token, err = helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserACredentials)
	assert.NoError(t, err)

	// try to delete a dashboard that does not exist; should return a not found error
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/dashboards/%v", newDashboardID+1), "DELETE", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 404, code, "Response body: \n%v\n", resp)

	// Count the number of all the dashboards returned for scenario
	initialNumber, err := helper.LengthOfResponse(router, token,
		fmt.Sprintf("/api/dashboards?scenarioID=%v", scenarioID), "GET", nil)
	assert.NoError(t, err)

	// Delete the added newDashboard
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/dashboards/%v", newDashboardID), "DELETE", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare DELETE's response with the newDashboard
	err = helper.CompareResponse(resp, helper.KeyModels{"dashboard": newDashboard})
	assert.NoError(t, err)

	// Again count the number of all the dashboards returned for scenario
	finalNumber, err := helper.LengthOfResponse(router, token,
		fmt.Sprintf("/api/dashboards?scenarioID=%v", scenarioID), "GET", nil)
	assert.NoError(t, err)

	assert.Equal(t, finalNumber, initialNumber-1)

}

func TestGetAllDashboardsOfScenario(t *testing.T) {
	database.DropTables(db)
	database.MigrateModels(db)
	assert.NoError(t, database.DBAddAdminAndUserAndGuest(db))

	// authenticate as normal user
	token, err := helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserACredentials)
	assert.NoError(t, err)

	scenarioID := addScenario(token)

	// Count the number of all the dashboards returned for scenario
	initialNumber, err := helper.LengthOfResponse(router, token,
		fmt.Sprintf("/api/dashboards?scenarioID=%v", scenarioID), "GET", nil)
	assert.NoError(t, err)

	// test POST dashboards/ $newDashboard
	newDashboardA := DashboardRequest{
		Name:       database.DashboardA.Name,
		Grid:       database.DashboardA.Grid,
		ScenarioID: scenarioID,
	}
	code, resp, err := helper.TestEndpoint(router, token,
		"/api/dashboards", "POST", helper.KeyModels{"dashboard": newDashboardA})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// POST a second dashboard for the same scenario
	newDashboardB := DashboardRequest{
		Name:       database.DashboardB.Name,
		Grid:       database.DashboardB.Grid,
		ScenarioID: scenarioID,
	}
	code, resp, err = helper.TestEndpoint(router, token,
		"/api/dashboards", "POST", helper.KeyModels{"dashboard": newDashboardB})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Count again the number of all the dashboards returned for scenario
	finalNumber, err := helper.LengthOfResponse(router, token,
		fmt.Sprintf("/api/dashboards?scenarioID=%v", scenarioID), "GET", nil)
	assert.NoError(t, err)

	assert.Equal(t, initialNumber+2, finalNumber)

	// try to get all dashboards of a scenario that does not exist (should fail with not found)
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/dashboards?scenarioID=%v", scenarioID+1), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 404, code, "Response body: \n%v\n", resp)

	// try to get all dashboards as a user that does not belong to scenario
	token, err = helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserBCredentials)
	assert.NoError(t, err)

	// this should fail with unprocessable entity
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/dashboards?scenarioID=%v", scenarioID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

}
