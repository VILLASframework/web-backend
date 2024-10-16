/**
* This file is part of VILLASweb-backend-go
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

package usergroup

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"testing"

	"git.rwth-aachen.de/acs/public/villas/web-backend-go/configuration"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/helper"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/scenario"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/user"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var router *gin.Engine

type ScenarioMappingRequest struct {
	ScenarioID uint `json:"scenarioID"`
	Duplicate  bool `json:"duplicate"`
}

type UserGroupRequest struct {
	Name             string                   `json:"name"`
	ScenarioMappings []ScenarioMappingRequest `json:"scenarioMappings"`
}

type UserRequest struct {
	Username    string `json:"username,omitempty"`
	Password    string `json:"password,omitempty"`
	OldPassword string `json:"oldPassword,omitempty"`
	Mail        string `json:"mail,omitempty"`
	Role        string `json:"role,omitempty"`
	Active      string `json:"active,omitempty"`
}

var newUserGroupOneMapping = UserGroupRequest{
	Name: "UserGroup1",
	ScenarioMappings: []ScenarioMappingRequest{
		{
			ScenarioID: 1,
			Duplicate:  false,
		},
	},
}

var newUserGroupTwoMappings = UserGroupRequest{
	Name: "UserGroup2",
	ScenarioMappings: []ScenarioMappingRequest{
		{
			ScenarioID: 1,
			Duplicate:  false,
		},
		{
			ScenarioID: 2,
			Duplicate:  true,
		},
	},
}

var deleteTestUg = UserGroupRequest{
	Name: "UserGroup3",
	ScenarioMappings: []ScenarioMappingRequest{
		{
			ScenarioID: 1,
			Duplicate:  false,
		},
		{
			ScenarioID: 2,
			Duplicate:  true,
		},
		{
			ScenarioID: 3,
			Duplicate:  true,
		},
	},
}

func TestMain(m *testing.M) {
	err := configuration.InitConfig()
	if err != nil {
		panic(m)
	}

	err = database.InitDB(configuration.GlobalConfig, true)
	if err != nil {
		panic(m)
	}
	defer database.DBpool.Close()

	router = gin.Default()
	api := router.Group("/api/v2")

	user.RegisterAuthenticate(api.Group("/authenticate"))
	api.Use(user.Authentication())
	scenario.RegisterScenarioEndpoints(api.Group("/scenarios"))
	// user endpoints required to set user to inactive
	user.RegisterUserEndpoints(api.Group("/users"))
	RegisterUserGroupEndpoints(api.Group("/usergroups"))

	os.Exit(m.Run())
}

func TestAddUserGroup(t *testing.T) {

	database.DropTables()
	database.MigrateModels()
	assert.NoError(t, database.AddTestUsers())

	token, err := helper.AuthenticateForTest(router, database.AdminCredentials)
	assert.NoError(t, err)

	// try to POST with non JSON body
	// should return a bad request error
	code, resp, err := helper.TestEndpoint(router, token,
		"/api/v2/usergroups", "POST", "this is not a JSON")
	assert.NoError(t, err)
	assert.Equalf(t, 400, code, "Response body: \n%v\n", resp)

	//Test with inexitent scenario
	code, resp, err = helper.TestEndpoint(router, token, "/api/v2/usergroups",
		"POST", helper.KeyModels{"usergroup": newUserGroupOneMapping})
	assert.NoError(t, err)
	assert.Equalf(t, 404, code, "Response body: \n%v\n", resp)

	// Test with valid user group with one scenario mapping
	helper.TestEndpoint(router, token, "/api/v2/scenarios", "POST", helper.KeyModels{"scenario": database.Scenario{Name: "scenario1"}})
	helper.TestEndpoint(router, token, "/api/v2/scenarios", "POST", helper.KeyModels{"scenario": database.Scenario{Name: "scenario2"}})
	code, resp, err = helper.TestEndpoint(router, token, "/api/v2/usergroups",
		"POST", helper.KeyModels{"usergroup": newUserGroupOneMapping})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Test with valid user group and two scenario mappings
	code, resp, err = helper.TestEndpoint(router, token, "/api/v2/usergroups",
		"POST", helper.KeyModels{"usergroup": newUserGroupTwoMappings})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Test with valid user group and multiple mappings
	// Test with invalid user group
	// Test with invalid user group and one mapping
	// Test with invalid user group and multiple mappings
}

func TestAddUserToGroup(t *testing.T) {
	// Prep DB
	database.DropTables()
	database.MigrateModels()
	adminpw, _ := database.AddAdminUser(configuration.GlobalConfig)

	//Auth
	token, _ := helper.AuthenticateForTest(router, database.Credentials{Username: "admin", Password: adminpw})

	//Post necessities
	helper.TestEndpoint(router, token, "/api/v2/scenarios", "POST", helper.KeyModels{"scenario": database.Scenario{Name: "scenarioNoDups"}})
	helper.TestEndpoint(router, token, "/api/v2/scenarios", "POST", helper.KeyModels{"scenario": database.Scenario{Name: "scenarioDups"}})
	helper.TestEndpoint(router, token, "/api/v2/usergroups", "POST", helper.KeyModels{"usergroup": newUserGroupTwoMappings})

	for n_users := 0; n_users < 3; n_users++ {
		//Add user
		n := strconv.Itoa(n_users + 1)
		usr := UserRequest{Username: "usr" + n, Password: "legendre" + n, Role: "User", Mail: "usr" + n + "@harmonics.de"}
		helper.TestEndpoint(router, token, "/api/v2/users", "POST", helper.KeyModels{"user": usr})
		code, _, err := helper.TestEndpoint(router, token, "/api/v2/usergroups/1/user?username=usr"+n, "PUT", struct{}{})
		assert.Equal(t, 200, code)
		assert.NoError(t, err)
	}

	//get scenarios
	_, res, _ := helper.TestEndpoint(router, token, "/api/v2/scenarios", "GET", struct{}{})
	var scenariosMap map[string]([]database.Scenario)
	json.Unmarshal(res.Bytes(), &scenariosMap)
	scenarios := scenariosMap["scenarios"]

	//Actual checks
	assert.Equal(t, 5, len(scenarios))

	for _, v := range scenarios {
		path := fmt.Sprintf("/api/v2/scenarios/%d/users", v.ID)
		var usersMap map[string]([]database.User)
		_, res, _ = helper.TestEndpoint(router, token, path, "GET", struct{}{})
		json.Unmarshal(res.Bytes(), &usersMap)
		users := usersMap["users"]
		switch v.ID {
		case 1: //no dups
			assert.Equal(t, "scenarioNoDups", v.Name)
			assert.Equal(t, 4, len(users))
		case 2: // with dups
			assert.Equal(t, "scenarioDups", v.Name)
			assert.Equal(t, 1, len(users))
			assert.Equal(t, "admin", users[0].Username)
		default:
			usr := "usr" + strconv.Itoa(int(v.ID-2)) // shift ids by the first two scenarios
			assert.Equal(t, "scenarioDups "+usr, v.Name)
			assert.Equal(t, 1, len(users))
			assert.Equal(t, usr, users[0].Username)
		}
	}

}

func TestDeleteUserFromGroup(t *testing.T) {
	// Prep DB
	database.DropTables()
	database.MigrateModels()
	adminpw, _ := database.AddAdminUser(configuration.GlobalConfig)

	//Auth
	token, _ := helper.AuthenticateForTest(router, database.Credentials{Username: "admin", Password: adminpw})

	//Post necessities
	helper.TestEndpoint(router, token, "/api/v2/scenarios", "POST", helper.KeyModels{"scenario": database.Scenario{Name: "scenarioNoDups"}})
	helper.TestEndpoint(router, token, "/api/v2/scenarios", "POST", helper.KeyModels{"scenario": database.Scenario{Name: "scenarioDups"}})
	helper.TestEndpoint(router, token, "/api/v2/usergroups", "POST", helper.KeyModels{"usergroup": newUserGroupTwoMappings})

	//Add 2 users
	for n_users := 0; n_users < 2; n_users++ {
		//Add user
		n := strconv.Itoa(n_users + 1)
		usr := UserRequest{Username: "usr" + n, Password: "legendre" + n, Role: "User", Mail: "usr" + n + "@harmonics.de"}
		helper.TestEndpoint(router, token, "/api/v2/users", "POST", helper.KeyModels{"user": usr})
		code, _, err := helper.TestEndpoint(router, token, "/api/v2/usergroups/1/user?username=usr"+n, "PUT", struct{}{})
		assert.Equal(t, 200, code)
		assert.NoError(t, err)
	}
	//we add usr1 to a group that doubles its right of acces to scenario 1
	helper.TestEndpoint(router, token, "/api/v2/usergroups", "POST", helper.KeyModels{"usergroup": newUserGroupOneMapping})
	helper.TestEndpoint(router, token, "/api/v2/usergroups/2/user?username=usr1", "PUT", struct{}{})

	//Delete one
	helper.TestEndpoint(router, token, "/api/v2/usergroups/1/user?username=usr1", "DELETE", struct{}{})

	//get scenarios
	_, res, _ := helper.TestEndpoint(router, token, "/api/v2/scenarios", "GET", struct{}{})
	var scenariosMap map[string]([]database.Scenario)
	json.Unmarshal(res.Bytes(), &scenariosMap)
	scenarios := scenariosMap["scenarios"]

	//Actual checks
	assert.Equal(t, 3, len(scenarios))
	for _, v := range scenarios {
		path := fmt.Sprintf("/api/v2/scenarios/%d/users", v.ID)
		var usersMap map[string]([]database.User)
		_, res, _ = helper.TestEndpoint(router, token, path, "GET", struct{}{})
		json.Unmarshal(res.Bytes(), &usersMap)
		users := usersMap["users"]
		switch v.ID {
		case 1: //no dups should still contain usr1 through ug2
			assert.Equal(t, "scenarioNoDups", v.Name)
			assert.Equal(t, 3, len(users))
		case 2: // with dups
			assert.Equal(t, "scenarioDups", v.Name)
			assert.Equal(t, 1, len(users))
			assert.Equal(t, "admin", users[0].Username)
		default: // remaining duped scenario
			assert.Equal(t, "scenarioDups usr2", v.Name)
			assert.Equal(t, 1, len(users))
			assert.Equal(t, "usr2", users[0].Username)
		}
	}
	//Delete from other
	helper.TestEndpoint(router, token, "/api/v2/usergroups/2/user?username=usr1", "DELETE", struct{}{})

	_, res, _ = helper.TestEndpoint(router, token, "/api/v2/scenarios/1", "GET", struct{}{})
	var scenarioMap map[string](database.Scenario)
	json.Unmarshal(res.Bytes(), &scenarioMap)
	scenario := scenarioMap["scenario"]

	path := fmt.Sprintf("/api/v2/scenarios/%d/users", scenario.ID)
	var usersMap map[string]([]database.User)
	_, res, _ = helper.TestEndpoint(router, token, path, "GET", struct{}{})
	json.Unmarshal(res.Bytes(), &usersMap)
	users := usersMap["users"]

	assert.Equal(t, 2, len(users))
	for _, u := range users {
		assert.NotEqual(t, "usr1", u.Username)
	}
}

func TestDeleteUserGroup(t *testing.T) {
	// Prep DB
	database.DropTables()
	database.MigrateModels()
	adminpw, _ := database.AddAdminUser(configuration.GlobalConfig)

	//Auth
	token, _ := helper.AuthenticateForTest(router, database.Credentials{Username: "admin", Password: adminpw})

	//Post necessities
	helper.TestEndpoint(router, token, "/api/v2/scenarios", "POST", helper.KeyModels{"scenario": database.Scenario{Name: "scenarioNoDups"}})
	helper.TestEndpoint(router, token, "/api/v2/scenarios", "POST", helper.KeyModels{"scenario": database.Scenario{Name: "scenarioDups1"}})
	helper.TestEndpoint(router, token, "/api/v2/scenarios", "POST", helper.KeyModels{"scenario": database.Scenario{Name: "scenarioDups2"}})
	helper.TestEndpoint(router, token, "/api/v2/usergroups", "POST", helper.KeyModels{"usergroup": deleteTestUg})

	//Add 2 users
	for n_users := 0; n_users < 2; n_users++ {
		//Add user
		n := strconv.Itoa(n_users + 1)
		usr := UserRequest{Username: "usr" + n, Password: "legendre" + n, Role: "User", Mail: "usr" + n + "@harmonics.de"}
		helper.TestEndpoint(router, token, "/api/v2/users", "POST", helper.KeyModels{"user": usr})
		helper.TestEndpoint(router, token, "/api/v2/usergroups/1/user?username=usr"+n, "PUT", struct{}{})
	}

	//we add usr1 to a group that doubles its right of acces to scenario 1
	helper.TestEndpoint(router, token, "/api/v2/usergroups", "POST", helper.KeyModels{"usergroup": newUserGroupOneMapping})
	helper.TestEndpoint(router, token, "/api/v2/usergroups/2/user?username=usr1", "PUT", struct{}{})

	//delete usergroup
	code, _, err := helper.TestEndpoint(router, token, "/api/v2/usergroups/1", "DELETE", struct{}{})
	assert.Equal(t, 200, code)
	assert.NoError(t, err)

	//get scenarios
	_, res, _ := helper.TestEndpoint(router, token, "/api/v2/scenarios", "GET", struct{}{})
	var scenariosMap map[string]([]database.Scenario)
	json.Unmarshal(res.Bytes(), &scenariosMap)
	scenarios := scenariosMap["scenarios"]

	//Actual checks
	assert.Equal(t, 3, len(scenarios))
	var exists []string = []string{"scenarioNoDups", "scenarioDups1", "scenarioDups2"}
	var deleted []string = []string{"scenarioDups1 usr1", "scenarioDups2 usr1", "scenarioDups1 usr2", "scenarioDups2 usr2"}
	for _, sc := range scenarios {
		assert.Contains(t, exists, sc.Name)
		assert.NotContains(t, deleted, sc.Name)
		if sc.Name == "scenarioNoDups" {
			_, res, _ = helper.TestEndpoint(router, token, "/api/v2/scenarios/1", "GET", struct{}{})
			var scenarioMap map[string](database.Scenario)
			json.Unmarshal(res.Bytes(), &scenarioMap)
			scenario := scenarioMap["scenario"]

			path := fmt.Sprintf("/api/v2/scenarios/%d/users", scenario.ID)
			var usersMap map[string]([]database.User)
			_, res, _ = helper.TestEndpoint(router, token, path, "GET", struct{}{})
			json.Unmarshal(res.Bytes(), &usersMap)
			users := usersMap["users"]
			assert.Equal(t, 2, len(users))
		}
	}
	//Delete from other
	helper.TestEndpoint(router, token, "/api/v2/usergroups/2/user?username=usr1", "DELETE", struct{}{})

	_, res, _ = helper.TestEndpoint(router, token, "/api/v2/scenarios/1", "GET", struct{}{})
	var scenarioMap map[string](database.Scenario)
	json.Unmarshal(res.Bytes(), &scenarioMap)
	scenario := scenarioMap["scenario"]

	path := fmt.Sprintf("/api/v2/scenarios/%d/users", scenario.ID)
	var usersMap map[string]([]database.User)
	_, res, _ = helper.TestEndpoint(router, token, path, "GET", struct{}{})
	json.Unmarshal(res.Bytes(), &usersMap)
	users := usersMap["users"]

	assert.Equal(t, 1, len(users))
	for _, u := range users {
		assert.NotEqual(t, "usr1", u.Username)
		assert.NotEqual(t, "usr2", u.Username)
	}
}
