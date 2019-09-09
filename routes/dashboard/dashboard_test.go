package dashboard

import (
	"fmt"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
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
		Name:            common.ScenarioA.Name,
		Running:         common.ScenarioA.Running,
		StartParameters: common.ScenarioA.StartParameters,
	}
	_, resp, _ := common.TestEndpoint(router, token,
		"/api/scenarios", "POST", common.KeyModels{"scenario": newScenario})

	// Read newScenario's ID from the response
	newScenarioID, _ := common.GetResponseID(resp)

	return uint(newScenarioID)
}

func TestMain(m *testing.M) {

	db = common.DummyInitDB()
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
	common.DropTables(db)
	common.MigrateModels(db)
	common.DummyAddOnlyUserTableWithAdminAndUsersDB(db)

	// authenticate as normal user
	token, err := common.AuthenticateForTest(router,
		"/api/authenticate", "POST", common.UserACredentials)
	assert.NoError(t, err)

	scenarioID := addScenario(token)

	// test POST dashboards/ $newDashboard
	newDashboard := DashboardRequest{
		Name:       common.DashboardA.Name,
		Grid:       common.DashboardA.Grid,
		ScenarioID: scenarioID,
	}
	code, resp, err := common.TestEndpoint(router, token,
		"/api/dashboards", "POST", common.KeyModels{"dashboard": newDashboard})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare POST's response with the newDashboard
	err = common.CompareResponse(resp, common.KeyModels{"dashboard": newDashboard})
	assert.NoError(t, err)

	// Read newDashboard's ID from the response
	newDashboardID, err := common.GetResponseID(resp)
	assert.NoError(t, err)

	// Get the newDashboard
	code, resp, err = common.TestEndpoint(router, token,
		fmt.Sprintf("/api/dashboards/%v", newDashboardID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare GET's response with the newDashboard
	err = common.CompareResponse(resp, common.KeyModels{"dashboard": newDashboard})
	assert.NoError(t, err)

	// try to POST a malformed dashboard
	// Required fields are missing
	malformedNewDashboard := DashboardRequest{
		Name: "ThisIsAMalformedDashboard",
	}
	// this should NOT work and return a unprocessable entity 442 status code
	code, resp, err = common.TestEndpoint(router, token,
		"/api/dashboards", "POST", common.KeyModels{"dashboard": malformedNewDashboard})
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

}

func TestUpdateDashboard(t *testing.T) {
	common.DropTables(db)
	common.MigrateModels(db)
	common.DummyAddOnlyUserTableWithAdminAndUsersDB(db)

	// authenticate as normal user
	token, err := common.AuthenticateForTest(router,
		"/api/authenticate", "POST", common.UserACredentials)
	assert.NoError(t, err)

	scenarioID := addScenario(token)

	// test POST dashboards/ $newDashboard
	newDashboard := DashboardRequest{
		Name:       common.DashboardA.Name,
		Grid:       common.DashboardA.Grid,
		ScenarioID: scenarioID,
	}
	code, resp, err := common.TestEndpoint(router, token,
		"/api/dashboards", "POST", common.KeyModels{"dashboard": newDashboard})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Read newDashboard's ID from the response
	newDashboardID, err := common.GetResponseID(resp)
	assert.NoError(t, err)

	updatedDashboard := DashboardRequest{
		Name: common.DashboardB.Name,
		Grid: common.DashboardB.Grid,
	}

	code, resp, err = common.TestEndpoint(router, token,
		fmt.Sprintf("/api/dashboards/%v", newDashboardID), "PUT", common.KeyModels{"dashboard": updatedDashboard})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare PUT's response with the updatedDashboard
	err = common.CompareResponse(resp, common.KeyModels{"dashboard": updatedDashboard})
	assert.NoError(t, err)

	// Get the updatedDashboard
	code, resp, err = common.TestEndpoint(router, token,
		fmt.Sprintf("/api/dashboards/%v", newDashboardID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare GET's response with the updatedDashboard
	err = common.CompareResponse(resp, common.KeyModels{"dashboard": updatedDashboard})
	assert.NoError(t, err)

	// try to update a dashboard that does not exist (should return not found 404 status code)
	code, resp, err = common.TestEndpoint(router, token,
		fmt.Sprintf("/api/dashboards/%v", newDashboardID+1), "PUT", common.KeyModels{"dashboard": updatedDashboard})
	assert.NoError(t, err)
	assert.Equalf(t, 404, code, "Response body: \n%v\n", resp)

}

func TestDeleteDashboard(t *testing.T) {
	common.DropTables(db)
	common.MigrateModels(db)
	common.DummyAddOnlyUserTableWithAdminAndUsersDB(db)

	// authenticate as normal user
	token, err := common.AuthenticateForTest(router,
		"/api/authenticate", "POST", common.UserACredentials)
	assert.NoError(t, err)

	scenarioID := addScenario(token)
	fmt.Println(scenarioID)

	// test POST dashboards/ $newDashboard
	newDashboard := DashboardRequest{
		Name:       common.DashboardA.Name,
		Grid:       common.DashboardA.Grid,
		ScenarioID: scenarioID,
	}
	code, resp, err := common.TestEndpoint(router, token,
		"/api/dashboards", "POST", common.KeyModels{"dashboard": newDashboard})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Read newDashboard's ID from the response
	newDashboardID, err := common.GetResponseID(resp)
	assert.NoError(t, err)

	// Count the number of all the dashboards returned for scenario
	initialNumber, err := common.LengthOfResponse(router, token,
		fmt.Sprintf("/api/dashboards?scenarioID=%v", scenarioID), "GET", nil)
	assert.NoError(t, err)

	// Delete the added newDashboard
	code, resp, err = common.TestEndpoint(router, token,
		fmt.Sprintf("/api/dashboards/%v", newDashboardID), "DELETE", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare DELETE's response with the newDashboard
	err = common.CompareResponse(resp, common.KeyModels{"dashboard": newDashboard})
	assert.NoError(t, err)

	// Again count the number of all the dashboards returned for scenario
	finalNumber, err := common.LengthOfResponse(router, token,
		fmt.Sprintf("/api/dashboards?scenarioID=%v", scenarioID), "GET", nil)
	assert.NoError(t, err)

	assert.Equal(t, finalNumber, initialNumber-1)

}

func TestGetAllDashboardsOfScenario(t *testing.T) {
	common.DropTables(db)
	common.MigrateModels(db)
	common.DummyAddOnlyUserTableWithAdminAndUsersDB(db)

	// authenticate as normal user
	token, err := common.AuthenticateForTest(router,
		"/api/authenticate", "POST", common.UserACredentials)
	assert.NoError(t, err)

	scenarioID := addScenario(token)
	fmt.Println(scenarioID)

	// Count the number of all the dashboards returned for scenario
	initialNumber, err := common.LengthOfResponse(router, token,
		fmt.Sprintf("/api/dashboards?scenarioID=%v", scenarioID), "GET", nil)
	assert.NoError(t, err)

	// test POST dashboards/ $newDashboard
	newDashboardA := DashboardRequest{
		Name:       common.DashboardA.Name,
		Grid:       common.DashboardA.Grid,
		ScenarioID: scenarioID,
	}
	code, resp, err := common.TestEndpoint(router, token,
		"/api/dashboards", "POST", common.KeyModels{"dashboard": newDashboardA})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// POST a second dashboard for the same scenario
	newDashboardB := DashboardRequest{
		Name:       common.DashboardB.Name,
		Grid:       common.DashboardB.Grid,
		ScenarioID: scenarioID,
	}
	code, resp, err = common.TestEndpoint(router, token,
		"/api/dashboards", "POST", common.KeyModels{"dashboard": newDashboardB})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Count again the number of all the dashboards returned for scenario
	finalNumber, err := common.LengthOfResponse(router, token,
		fmt.Sprintf("/api/dashboards?scenarioID=%v", scenarioID), "GET", nil)
	assert.NoError(t, err)

	assert.Equal(t, initialNumber+2, finalNumber)
}
