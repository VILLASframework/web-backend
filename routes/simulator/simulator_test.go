package simulator

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"

	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/user"
)

var router *gin.Engine
var db *gorm.DB

type SimulatorRequest struct {
	UUID       string         `json:"uuid,omitempty"`
	Host       string         `json:"host,omitempty"`
	Modeltype  string         `json:"modelType,omitempty"`
	State      string         `json:"state,omitempty"`
	Properties postgres.Jsonb `json:"properties,omitempty"`
}

func TestMain(m *testing.M) {

	db = common.DummyInitDB()
	defer db.Close()

	router = gin.Default()
	api := router.Group("/api")

	user.RegisterAuthenticate(api.Group("/authenticate"))
	api.Use(user.Authentication(true))
	RegisterSimulatorEndpoints(api.Group("/simulators"))

	os.Exit(m.Run())
}

func TestAddSimulatorAsAdmin(t *testing.T) {
	common.DropTables(db)
	common.MigrateModels(db)
	common.DummyAddOnlyUserTableWithAdminAndUsersDB(db)

	// authenticate as admin
	token, err := common.NewAuthenticateForTest(router,
		"/api/authenticate", "POST", common.AdminCredentials)
	assert.NoError(t, err)

	// test POST simulators/ $newSimulator
	newSimulator := SimulatorRequest{
		UUID:       common.SimulatorA.UUID,
		Host:       common.SimulatorA.Host,
		Modeltype:  common.SimulatorA.Modeltype,
		State:      common.SimulatorA.State,
		Properties: common.SimulatorA.Properties,
	}
	code, resp, err := common.NewTestEndpoint(router, token,
		"/api/simulators", "POST", common.KeyModels{"simulator": newSimulator})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare POST's response with the newSimulator
	err = common.CompareResponse(resp, common.KeyModels{"simulator": newSimulator})
	assert.NoError(t, err)

	// Read newSimulator's ID from the response
	newSimulatorID, err := common.GetResponseID(resp)
	assert.NoError(t, err)

	// Get the newSimulator
	code, resp, err = common.NewTestEndpoint(router, token,
		fmt.Sprintf("/api/simulators/%v", newSimulatorID), "GET", nil)
	assert.NoError(t, err)

	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare GET's response with the newSimulator
	err = common.CompareResponse(resp, common.KeyModels{"simulator": newSimulator})
	assert.NoError(t, err)
}

func TestAddSimulatorAsUser(t *testing.T) {
	common.DropTables(db)
	common.MigrateModels(db)
	common.DummyAddOnlyUserTableWithAdminAndUsersDB(db)

	// authenticate as user
	token, err := common.NewAuthenticateForTest(router,
		"/api/authenticate", "POST", common.UserACredentials)
	assert.NoError(t, err)

	// test POST simulators/ $newSimulator
	newSimulator := SimulatorRequest{
		UUID:       common.SimulatorA.UUID,
		Host:       common.SimulatorA.Host,
		Modeltype:  common.SimulatorA.Modeltype,
		State:      common.SimulatorA.State,
		Properties: common.SimulatorA.Properties,
	}

	// This should fail with unprocessable entity 422 error code
	// Normal users are not allowed to add simulators
	code, resp, err := common.NewTestEndpoint(router, token,
		"/api/simulators", "POST", common.KeyModels{"simulator": newSimulator})
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)
}

func TestUpdateSimulatorAsAdmin(t *testing.T) {
	common.DropTables(db)
	common.MigrateModels(db)
	common.DummyAddOnlyUserTableWithAdminAndUsersDB(db)

	// authenticate as admin
	token, err := common.NewAuthenticateForTest(router,
		"/api/authenticate", "POST", common.AdminCredentials)
	assert.NoError(t, err)

	// test POST simulators/ $newSimulator
	newSimulator := SimulatorRequest{
		UUID:       common.SimulatorA.UUID,
		Host:       common.SimulatorA.Host,
		Modeltype:  common.SimulatorA.Modeltype,
		State:      common.SimulatorA.State,
		Properties: common.SimulatorA.Properties,
	}
	code, resp, err := common.NewTestEndpoint(router, token,
		"/api/simulators", "POST", common.KeyModels{"simulator": newSimulator})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare POST's response with the newSimulator
	err = common.CompareResponse(resp, common.KeyModels{"simulator": newSimulator})
	assert.NoError(t, err)

	// Read newSimulator's ID from the response
	newSimulatorID, err := common.GetResponseID(resp)
	assert.NoError(t, err)

	// Test PUT simulators
	newSimulator.Host = "ThisIsMyNewHost"
	code, resp, err = common.NewTestEndpoint(router, token,
		fmt.Sprintf("/api/simulators/%v", newSimulatorID), "PUT", common.KeyModels{"simulator": newSimulator})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare PUT's response with the updated newSimulator
	err = common.CompareResponse(resp, common.KeyModels{"simulator": newSimulator})
	assert.NoError(t, err)

	// Get the updated newSimulator
	code, resp, err = common.NewTestEndpoint(router, token,
		fmt.Sprintf("/api/simulators/%v", newSimulatorID), "GET", nil)
	assert.NoError(t, err)

	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare GET's response with the updated newSimulator
	err = common.CompareResponse(resp, common.KeyModels{"simulator": newSimulator})
	assert.NoError(t, err)

}

func TestUpdateSimulatorAsUser(t *testing.T) {
	common.DropTables(db)
	common.MigrateModels(db)
	common.DummyAddOnlyUserTableWithAdminAndUsersDB(db)

	// authenticate as admin
	token, err := common.NewAuthenticateForTest(router,
		"/api/authenticate", "POST", common.AdminCredentials)
	assert.NoError(t, err)

	// test POST simulators/ $newSimulator
	newSimulator := SimulatorRequest{
		UUID:       common.SimulatorA.UUID,
		Host:       common.SimulatorA.Host,
		Modeltype:  common.SimulatorA.Modeltype,
		State:      common.SimulatorA.State,
		Properties: common.SimulatorA.Properties,
	}
	code, resp, err := common.NewTestEndpoint(router, token,
		"/api/simulators", "POST", common.KeyModels{"simulator": newSimulator})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Read newSimulator's ID from the response
	newSimulatorID, err := common.GetResponseID(resp)
	assert.NoError(t, err)

	// authenticate as user
	token, err = common.NewAuthenticateForTest(router,
		"/api/authenticate", "POST", common.UserACredentials)
	assert.NoError(t, err)

	// Test PUT simulators
	// This should fail with unprocessable entity status code 422
	newSimulator.Host = "ThisIsMyNewHost"
	code, resp, err = common.NewTestEndpoint(router, token,
		fmt.Sprintf("/api/simulators/%v", newSimulatorID), "PUT", common.KeyModels{"simulator": newSimulator})
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

}

func TestDeleteSimulatorAsAdmin(t *testing.T) {
	common.DropTables(db)
	common.MigrateModels(db)
	common.DummyAddOnlyUserTableWithAdminAndUsersDB(db)

	// authenticate as admin
	token, err := common.NewAuthenticateForTest(router,
		"/api/authenticate", "POST", common.AdminCredentials)
	assert.NoError(t, err)

	// test POST simulators/ $newSimulator
	newSimulator := SimulatorRequest{
		UUID:       common.SimulatorA.UUID,
		Host:       common.SimulatorA.Host,
		Modeltype:  common.SimulatorA.Modeltype,
		State:      common.SimulatorA.State,
		Properties: common.SimulatorA.Properties,
	}
	code, resp, err := common.NewTestEndpoint(router, token,
		"/api/simulators", "POST", common.KeyModels{"simulator": newSimulator})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Read newSimulator's ID from the response
	newSimulatorID, err := common.GetResponseID(resp)
	assert.NoError(t, err)

	// Count the number of all the simulators returned for admin
	initialNumber, err := common.LengthOfResponse(router, token,
		"/api/simulators", "GET", nil)
	assert.NoError(t, err)

	// Delete the added newSimulator
	code, resp, err = common.NewTestEndpoint(router, token,
		fmt.Sprintf("/api/simulators/%v", newSimulatorID), "DELETE", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare DELETE's response with the newSimulator
	err = common.CompareResponse(resp, common.KeyModels{"simulator": newSimulator})
	assert.NoError(t, err)

	// Again count the number of all the simulators returned
	finalNumber, err := common.LengthOfResponse(router, token,
		"/api/simulators", "GET", nil)
	assert.NoError(t, err)

	assert.Equal(t, finalNumber, initialNumber-1)
}

func TestDeleteSimulatorAsUser(t *testing.T) {
	common.DropTables(db)
	common.MigrateModels(db)
	common.DummyAddOnlyUserTableWithAdminAndUsersDB(db)

	// authenticate as admin
	token, err := common.NewAuthenticateForTest(router,
		"/api/authenticate", "POST", common.AdminCredentials)
	assert.NoError(t, err)

	// test POST simulators/ $newSimulator
	newSimulator := SimulatorRequest{
		UUID:       common.SimulatorA.UUID,
		Host:       common.SimulatorA.Host,
		Modeltype:  common.SimulatorA.Modeltype,
		State:      common.SimulatorA.State,
		Properties: common.SimulatorA.Properties,
	}
	code, resp, err := common.NewTestEndpoint(router, token,
		"/api/simulators", "POST", common.KeyModels{"simulator": newSimulator})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Read newSimulator's ID from the response
	newSimulatorID, err := common.GetResponseID(resp)
	assert.NoError(t, err)

	// authenticate as user
	token, err = common.NewAuthenticateForTest(router,
		"/api/authenticate", "POST", common.UserACredentials)
	assert.NoError(t, err)

	// Test DELETE simulators
	// This should fail with unprocessable entity status code 422
	newSimulator.Host = "ThisIsMyNewHost"
	code, resp, err = common.NewTestEndpoint(router, token,
		fmt.Sprintf("/api/simulators/%v", newSimulatorID), "DELETE", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)
}

func TestGetAllSimulators(t *testing.T) {
	common.DropTables(db)
	common.MigrateModels(db)
	common.DummyAddOnlyUserTableWithAdminAndUsersDB(db)

	// authenticate as admin
	token, err := common.NewAuthenticateForTest(router,
		"/api/authenticate", "POST", common.AdminCredentials)
	assert.NoError(t, err)

	// get the length of the GET all simulators response for user
	initialNumber, err := common.LengthOfResponse(router, token,
		"/api/simulators", "GET", nil)
	assert.NoError(t, err)

	// test POST simulators/ $newSimulatorA
	newSimulatorA := SimulatorRequest{
		UUID:       common.SimulatorA.UUID,
		Host:       common.SimulatorA.Host,
		Modeltype:  common.SimulatorA.Modeltype,
		State:      common.SimulatorA.State,
		Properties: common.SimulatorA.Properties,
	}
	code, resp, err := common.NewTestEndpoint(router, token,
		"/api/simulators", "POST", common.KeyModels{"simulator": newSimulatorA})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// test POST simulators/ $newSimulatorB
	newSimulatorB := SimulatorRequest{
		UUID:       common.SimulatorB.UUID,
		Host:       common.SimulatorB.Host,
		Modeltype:  common.SimulatorB.Modeltype,
		State:      common.SimulatorB.State,
		Properties: common.SimulatorB.Properties,
	}
	code, resp, err = common.NewTestEndpoint(router, token,
		"/api/simulators", "POST", common.KeyModels{"simulator": newSimulatorB})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// get the length of the GET all simulators response again
	finalNumber, err := common.LengthOfResponse(router, token,
		"/api/simulators", "GET", nil)
	assert.NoError(t, err)

	assert.Equal(t, finalNumber, initialNumber+2)

	// authenticate as normal user
	token, err = common.NewAuthenticateForTest(router,
		"/api/authenticate", "POST", common.UserACredentials)
	assert.NoError(t, err)

	// get the length of the GET all simulators response again
	finalNumber2, err := common.LengthOfResponse(router, token,
		"/api/simulators", "GET", nil)
	assert.NoError(t, err)

	assert.Equal(t, finalNumber2, initialNumber+2)
}

func TestGetSimulationModelsOfSimulator(t *testing.T) {
	common.DropTables(db)
	common.MigrateModels(db)
	common.DummyAddOnlyUserTableWithAdminAndUsersDB(db)

	// authenticate as admin
	token, err := common.NewAuthenticateForTest(router,
		"/api/authenticate", "POST", common.AdminCredentials)
	assert.NoError(t, err)

	// test POST simulators/ $newSimulatorA
	newSimulatorA := SimulatorRequest{
		UUID:       common.SimulatorA.UUID,
		Host:       common.SimulatorA.Host,
		Modeltype:  common.SimulatorA.Modeltype,
		State:      common.SimulatorA.State,
		Properties: common.SimulatorA.Properties,
	}
	code, resp, err := common.NewTestEndpoint(router, token,
		"/api/simulators", "POST", common.KeyModels{"simulator": newSimulatorA})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Read newSimulator's ID from the response
	newSimulatorID, err := common.GetResponseID(resp)
	assert.NoError(t, err)

	// test GET simulators/ID/models
	// TODO how to properly test this without using simulation model endpoints?
	numberOfModels, err := common.LengthOfResponse(router, token,
		fmt.Sprintf("/api/simulators/%v/models", newSimulatorID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	assert.Equal(t, 0, numberOfModels)

	// authenticate as normal user
	token, err = common.NewAuthenticateForTest(router,
		"/api/authenticate", "POST", common.UserACredentials)
	assert.NoError(t, err)

	// test GET simulators/ID/models
	// TODO how to properly test this without using simulation model endpoints?
	numberOfModels, err = common.LengthOfResponse(router, token,
		fmt.Sprintf("/api/simulators/%v/models", newSimulatorID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	assert.Equal(t, 0, numberOfModels)
}
