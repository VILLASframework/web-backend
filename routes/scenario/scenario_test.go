package scenario

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/user"
)

var router *gin.Engine
var db *gorm.DB

type ScenarioRequest struct {
	Name            string         `json:"name,omitempty"`
	Running         bool           `json:"running,omitempty"`
	StartParameters postgres.Jsonb `json:"startParameters,omitempty"`
}

func TestMain(m *testing.M) {

	db = common.DummyInitDB()
	defer db.Close()

	router = gin.Default()
	api := router.Group("/api")

	user.RegisterAuthenticate(api.Group("/authenticate"))
	api.Use(user.Authentication(true))
	RegisterScenarioEndpoints(api.Group("/scenarios"))

	os.Exit(m.Run())
}

func TestAddScenario(t *testing.T) {

	common.DropTables(db)
	common.MigrateModels(db)
	common.DummyAddOnlyUserTableWithAdminAndUsersDB(db)

	// authenticate as normal user
	token, err := common.NewAuthenticateForTest(router,
		"/api/authenticate", "POST", common.UserACredentials)
	assert.NoError(t, err)

	// test POST scenarios/ $newScenario
	newScenario := ScenarioRequest{
		Name:            common.ScenarioA.Name,
		Running:         common.ScenarioA.Running,
		StartParameters: common.ScenarioA.StartParameters,
	}
	code, resp, err := common.NewTestEndpoint(router, token,
		"/api/scenarios", "POST", common.KeyModels{"scenario": newScenario})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare POST's response with the newScenario
	err = common.CompareResponse(resp, common.KeyModels{"scenario": newScenario})
	assert.NoError(t, err)

	// Read newScenario's ID from the response
	newScenarioID, err := common.GetResponseID(resp)
	assert.NoError(t, err)

	// Get the newScenario
	code, resp, err = common.NewTestEndpoint(router, token,
		fmt.Sprintf("/api/scenarios/%v", newScenarioID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare GET's response with the newScenario
	err = common.CompareResponse(resp, common.KeyModels{"scenario": newScenario})
	assert.NoError(t, err)

	// try to POST a malformed scenario
	// Required fields are missing
	malformedNewScenario := ScenarioRequest{
		Running: false,
	}
	// this should NOT work and return a unprocessable entity 442 status code
	code, resp, err = common.NewTestEndpoint(router, token,
		"/api/scenarios", "POST", common.KeyModels{"scenario": malformedNewScenario})
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)
}

func TestUpdateScenario(t *testing.T) {

	common.DropTables(db)
	common.MigrateModels(db)
	common.DummyAddOnlyUserTableWithAdminAndUsersDB(db)

	// authenticate as normal user
	token, err := common.NewAuthenticateForTest(router,
		"/api/authenticate", "POST", common.UserACredentials)
	assert.NoError(t, err)

	// test POST scenarios/ $newScenario
	newScenario := ScenarioRequest{
		Name:            common.ScenarioA.Name,
		Running:         common.ScenarioA.Running,
		StartParameters: common.ScenarioA.StartParameters,
	}
	code, resp, err := common.NewTestEndpoint(router, token,
		"/api/scenarios", "POST", common.KeyModels{"scenario": newScenario})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare POST's response with the newScenario
	err = common.CompareResponse(resp, common.KeyModels{"scenario": newScenario})
	assert.NoError(t, err)

	// Read newScenario's ID from the response
	newScenarioID, err := common.GetResponseID(resp)
	assert.NoError(t, err)

	updatedScenario := ScenarioRequest{
		Name:            "Updated name",
		Running:         !common.ScenarioA.Running,
		StartParameters: common.ScenarioA.StartParameters,
	}

	code, resp, err = common.NewTestEndpoint(router, token,
		fmt.Sprintf("/api/scenarios/%v", newScenarioID), "PUT", common.KeyModels{"scenario": updatedScenario})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare PUT's response with the updatedScenario
	err = common.CompareResponse(resp, common.KeyModels{"scenario": updatedScenario})
	assert.NoError(t, err)

	// Get the updatedScenario
	code, resp, err = common.NewTestEndpoint(router, token,
		fmt.Sprintf("/api/scenarios/%v", newScenarioID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare GET's response with the newScenario
	err = common.CompareResponse(resp, common.KeyModels{"scenario": updatedScenario})
	assert.NoError(t, err)

	// try to update a scenario that does not exist (should return not found 404 status code)
	code, resp, err = common.NewTestEndpoint(router, token,
		fmt.Sprintf("/api/scenarios/%v", newScenarioID+1), "PUT", common.KeyModels{"scenario": updatedScenario})
	assert.NoError(t, err)
	assert.Equalf(t, 404, code, "Response body: \n%v\n", resp)

}

func TestGetAllScenariosAsAdmin(t *testing.T) {

	common.DropTables(db)
	common.MigrateModels(db)
	common.DummyAddOnlyUserTableWithAdminAndUsersDB(db)

	// authenticate as admin
	token, err := common.NewAuthenticateForTest(router,
		"/api/authenticate", "POST", common.AdminCredentials)
	assert.NoError(t, err)

	// get the length of the GET all scenarios response for admin
	initialNumber, err := common.LengthOfResponse(router, token,
		"/api/scenarios", "GET", nil)
	assert.NoError(t, err)

	// authenticate as normal userB
	token, err = common.NewAuthenticateForTest(router,
		"/api/authenticate", "POST", common.UserBCredentials)
	assert.NoError(t, err)

	// test POST scenarios/ $newScenarioB
	newScenarioB := ScenarioRequest{
		Name:            common.ScenarioB.Name,
		Running:         common.ScenarioB.Running,
		StartParameters: common.ScenarioB.StartParameters,
	}
	code, resp, err := common.NewTestEndpoint(router, token,
		"/api/scenarios", "POST", common.KeyModels{"scenario": newScenarioB})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// authenticate as normal userA
	token, err = common.NewAuthenticateForTest(router,
		"/api/authenticate", "POST", common.UserACredentials)
	assert.NoError(t, err)

	// test POST scenarios/ $newScenarioA
	newScenarioA := ScenarioRequest{
		Name:            common.ScenarioA.Name,
		Running:         common.ScenarioA.Running,
		StartParameters: common.ScenarioA.StartParameters,
	}
	code, resp, err = common.NewTestEndpoint(router, token,
		"/api/scenarios", "POST", common.KeyModels{"scenario": newScenarioA})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// authenticate as admin
	token, err = common.NewAuthenticateForTest(router,
		"/api/authenticate", "POST", common.AdminCredentials)
	assert.NoError(t, err)

	// get the length of the GET all scenarios response again
	finalNumber, err := common.LengthOfResponse(router, token,
		"/api/scenarios", "GET", nil)
	assert.NoError(t, err)

	assert.Equal(t, finalNumber, initialNumber+2)
}

func TestGetAllScenariosAsUser(t *testing.T) {

	common.DropTables(db)
	common.MigrateModels(db)
	common.DummyAddOnlyUserTableWithAdminAndUsersDB(db)

	// authenticate as normal userB
	token, err := common.NewAuthenticateForTest(router,
		"/api/authenticate", "POST", common.UserBCredentials)
	assert.NoError(t, err)

	// get the length of the GET all scenarios response for userB
	initialNumber, err := common.LengthOfResponse(router, token,
		"/api/scenarios", "GET", nil)
	assert.NoError(t, err)

	// test POST scenarios/ $newScenarioB
	newScenarioB := ScenarioRequest{
		Name:            common.ScenarioB.Name,
		Running:         common.ScenarioB.Running,
		StartParameters: common.ScenarioB.StartParameters,
	}

	code, resp, err := common.NewTestEndpoint(router, token,
		"/api/scenarios", "POST", common.KeyModels{"scenario": newScenarioB})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// authenticate as normal userA
	token, err = common.NewAuthenticateForTest(router,
		"/api/authenticate", "POST", common.UserACredentials)
	assert.NoError(t, err)

	// test POST scenarios/ $newScenarioA
	newScenarioA := ScenarioRequest{
		Name:            common.ScenarioA.Name,
		Running:         common.ScenarioA.Running,
		StartParameters: common.ScenarioA.StartParameters,
	}
	code, resp, err = common.NewTestEndpoint(router, token,
		"/api/scenarios", "POST", common.KeyModels{"scenario": newScenarioA})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// authenticate as normal userB
	token, err = common.NewAuthenticateForTest(router,
		"/api/authenticate", "POST", common.UserBCredentials)
	assert.NoError(t, err)

	// get the length of the GET all scenarios response again
	finalNumber, err := common.LengthOfResponse(router, token,
		"/api/scenarios", "GET", nil)
	assert.NoError(t, err)

	assert.Equal(t, finalNumber, initialNumber+1)
}

func TestDeleteScenario(t *testing.T) {

	common.DropTables(db)
	common.MigrateModels(db)
	common.DummyAddOnlyUserTableWithAdminAndUsersDB(db)

	// authenticate as normal user
	token, err := common.NewAuthenticateForTest(router,
		"/api/authenticate", "POST", common.UserACredentials)
	assert.NoError(t, err)

	// test POST scenarios/ $newScenario
	newScenario := ScenarioRequest{
		Name:            common.ScenarioA.Name,
		Running:         common.ScenarioA.Running,
		StartParameters: common.ScenarioA.StartParameters,
	}
	code, resp, err := common.NewTestEndpoint(router, token,
		"/api/scenarios", "POST", common.KeyModels{"scenario": newScenario})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Read newScenario's ID from the response
	newScenarioID, err := common.GetResponseID(resp)
	assert.NoError(t, err)

	// Count the number of all the scenarios returned for userA
	initialNumber, err := common.LengthOfResponse(router, token,
		"/api/scenarios", "GET", nil)
	assert.NoError(t, err)

	// Delete the added newScenario
	code, resp, err = common.NewTestEndpoint(router, token,
		fmt.Sprintf("/api/scenarios/%v", newScenarioID), "DELETE", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare DELETE's response with the newScenario
	err = common.CompareResponse(resp, common.KeyModels{"scenario": newScenario})
	assert.NoError(t, err)

	// Again count the number of all the users returned
	finalNumber, err := common.LengthOfResponse(router, token,
		"/api/scenarios", "GET", nil)
	assert.NoError(t, err)

	assert.Equal(t, finalNumber, initialNumber-1)
}

func TestAddUserToScenario(t *testing.T) {

	common.DropTables(db)
	common.MigrateModels(db)
	common.DummyAddOnlyUserTableWithAdminAndUsersDB(db)

	// authenticate as normal user
	token, err := common.NewAuthenticateForTest(router,
		"/api/authenticate", "POST", common.UserACredentials)
	assert.NoError(t, err)

	// test POST scenarios/ $newScenario
	newScenario := ScenarioRequest{
		Name:            common.ScenarioA.Name,
		Running:         common.ScenarioA.Running,
		StartParameters: common.ScenarioA.StartParameters,
	}
	code, resp, err := common.NewTestEndpoint(router, token,
		"/api/scenarios", "POST", common.KeyModels{"scenario": newScenario})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Read newScenario's ID from the response
	newScenarioID, err := common.GetResponseID(resp)
	assert.NoError(t, err)

	// Count the number of all the users returned for newScenario
	initialNumber, err := common.LengthOfResponse(router, token,
		fmt.Sprintf("/api/scenarios/%v/users", newScenarioID), "GET", nil)
	assert.NoError(t, err)
	assert.Equal(t, initialNumber, 1)

	// add userB to newScenario
	code, resp, err = common.NewTestEndpoint(router, token,
		fmt.Sprintf("/api/scenarios/%v/user?username=User_B", newScenarioID), "PUT", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare resp to userB
	err = common.CompareResponse(resp, common.KeyModels{"user": common.UserB})
	assert.NoError(t, err)

	// Count AGAIN the number of all the users returned for newScenario
	finalNumber, err := common.LengthOfResponse(router, token,
		fmt.Sprintf("/api/scenarios/%v/users", newScenarioID), "GET", nil)
	assert.NoError(t, err)
	assert.Equal(t, finalNumber, initialNumber+1)

	// try to add a non-existing user to newScenario, should return a not found 404
	code, resp, err = common.NewTestEndpoint(router, token,
		fmt.Sprintf("/api/scenarios/%v/user?username=User_C", newScenarioID), "PUT", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 404, code, "Response body: \n%v\n", resp)
}

func TestGetAllUsersOfScenario(t *testing.T) {

	common.DropTables(db)
	common.MigrateModels(db)
	common.DummyAddOnlyUserTableWithAdminAndUsersDB(db)

	// authenticate as normal user
	token, err := common.NewAuthenticateForTest(router,
		"/api/authenticate", "POST", common.UserACredentials)
	assert.NoError(t, err)

	// test POST scenarios/ $newScenario
	newScenario := ScenarioRequest{
		Name:            common.ScenarioA.Name,
		Running:         common.ScenarioA.Running,
		StartParameters: common.ScenarioA.StartParameters,
	}
	code, resp, err := common.NewTestEndpoint(router, token,
		"/api/scenarios", "POST", common.KeyModels{"scenario": newScenario})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Read newScenario's ID from the response
	newScenarioID, err := common.GetResponseID(resp)
	assert.NoError(t, err)

	// Count the number of all the users returned for newScenario
	initialNumber, err := common.LengthOfResponse(router, token,
		fmt.Sprintf("/api/scenarios/%v/users", newScenarioID), "GET", nil)
	assert.NoError(t, err)
	assert.Equal(t, initialNumber, 1)

	// add userB to newScenario
	code, resp, err = common.NewTestEndpoint(router, token,
		fmt.Sprintf("/api/scenarios/%v/user?username=User_B", newScenarioID), "PUT", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Count AGAIN the number of all the users returned for newScenario
	finalNumber, err := common.LengthOfResponse(router, token,
		fmt.Sprintf("/api/scenarios/%v/users", newScenarioID), "GET", nil)
	assert.NoError(t, err)
	assert.Equal(t, finalNumber, initialNumber+1)
}

func TestRemoveUserFromScenario(t *testing.T) {

	common.DropTables(db)
	common.MigrateModels(db)
	common.DummyAddOnlyUserTableWithAdminAndUsersDB(db)

	// authenticate as normal user
	token, err := common.NewAuthenticateForTest(router,
		"/api/authenticate", "POST", common.UserACredentials)
	assert.NoError(t, err)

	// test POST scenarios/ $newScenario
	newScenario := ScenarioRequest{
		Name:            common.ScenarioA.Name,
		Running:         common.ScenarioA.Running,
		StartParameters: common.ScenarioA.StartParameters,
	}
	code, resp, err := common.NewTestEndpoint(router, token,
		"/api/scenarios", "POST", common.KeyModels{"scenario": newScenario})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Read newScenario's ID from the response
	newScenarioID, err := common.GetResponseID(resp)
	assert.NoError(t, err)

	// add userB to newScenario
	code, resp, err = common.NewTestEndpoint(router, token,
		fmt.Sprintf("/api/scenarios/%v/user?username=User_B", newScenarioID), "PUT", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Count the number of all the users returned for newScenario
	initialNumber, err := common.LengthOfResponse(router, token,
		fmt.Sprintf("/api/scenarios/%v/users", newScenarioID), "GET", nil)
	assert.NoError(t, err)
	assert.Equal(t, 2, initialNumber)

	// remove userB from newScenario
	code, resp, err = common.NewTestEndpoint(router, token,
		fmt.Sprintf("/api/scenarios/%v/user?username=User_B", newScenarioID), "DELETE", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare DELETE's response with UserB
	err = common.CompareResponse(resp, common.KeyModels{"user": common.UserB})
	assert.NoError(t, err)

	// Count AGAIN the number of all the users returned for newScenario
	finalNumber, err := common.LengthOfResponse(router, token,
		fmt.Sprintf("/api/scenarios/%v/users", newScenarioID), "GET", nil)
	assert.NoError(t, err)
	assert.Equal(t, initialNumber-1, finalNumber)

	// Try to remove userA from new scenario
	// This should fail since User_A is the last user of newScenario
	code, resp, err = common.NewTestEndpoint(router, token,
		fmt.Sprintf("/api/scenarios/%v/user?username=User_A", newScenarioID), "DELETE", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 500, code, "Response body: \n%v\n", resp)

	// Try to remove a user that does not exist in DB
	// This should fail with not found 404 status code
	code, resp, err = common.NewTestEndpoint(router, token,
		fmt.Sprintf("/api/scenarios/%v/user?username=User_C", newScenarioID), "DELETE", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 404, code, "Response body: \n%v\n", resp)

	// Try to remove an admin user that is not explicitly a user of the scenario
	// This should fail with not found 404 status code
	code, resp, err = common.NewTestEndpoint(router, token,
		fmt.Sprintf("/api/scenarios/%v/user?username=User_0", newScenarioID), "DELETE", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 404, code, "Response body: \n%v\n", resp)

}
