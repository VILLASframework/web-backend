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

var newUserGroupNoMapping = UserGroupRequest{
	Name: "UserGroupNoMapping",
}

var ug_AddScenario1 = UserGroupRequest{
	Name: "UserGroup1",
	ScenarioMappings: []ScenarioMappingRequest{
		{
			ScenarioID: 1,
			Duplicate:  false,
		},
	},
}

var ug_AddScenario1_DuplicateScenario2 = UserGroupRequest{
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

var initUpdateTestUg = UserGroupRequest{
	Name: "UserGroup4",
	ScenarioMappings: []ScenarioMappingRequest{
		{
			ScenarioID: 1,
			Duplicate:  true,
		},
		{
			ScenarioID: 2,
			Duplicate:  false,
		},
		{
			ScenarioID: 3,
			Duplicate:  true,
		},
		{
			ScenarioID: 4,
			Duplicate:  false,
		},
		{
			ScenarioID: 5,
			Duplicate:  true,
		},
		{
			ScenarioID: 6,
			Duplicate:  false,
		},
	},
}

var updateTestUg = UserGroupRequest{
	Name: "UserGroup4",
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
			ScenarioID: 5,
			Duplicate:  true,
		},
		{
			ScenarioID: 6,
			Duplicate:  false,
		},
		{
			ScenarioID: 7,
			Duplicate:  true,
		},
		{
			ScenarioID: 8,
			Duplicate:  false,
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

	//Test with inexistent scenario
	code, resp, err = helper.TestEndpoint(router, token, "/api/v2/usergroups",
		"POST", helper.KeyModels{"usergroup": ug_AddScenario1})
	assert.NoError(t, err)
	assert.Equalf(t, 404, code, "Response body: \n%v\n", resp)

	// Test with valid user group with no scenario mappings
	code, resp, err = helper.TestEndpoint(router, token, "/api/v2/usergroups",
		"POST", helper.KeyModels{"usergroup": newUserGroupNoMapping})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Test with valid user group with one scenario mapping
	helper.TestEndpoint(router, token, "/api/v2/scenarios", "POST", helper.KeyModels{"scenario": database.Scenario{Name: "scenario1"}})
	helper.TestEndpoint(router, token, "/api/v2/scenarios", "POST", helper.KeyModels{"scenario": database.Scenario{Name: "scenario2"}})
	code, resp, err = helper.TestEndpoint(router, token, "/api/v2/usergroups",
		"POST", helper.KeyModels{"usergroup": ug_AddScenario1})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Test with valid user group and two scenario mappings
	code, resp, err = helper.TestEndpoint(router, token, "/api/v2/usergroups",
		"POST", helper.KeyModels{"usergroup": ug_AddScenario1_DuplicateScenario2})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)
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
	helper.TestEndpoint(router, token, "/api/v2/usergroups", "POST", helper.KeyModels{"usergroup": ug_AddScenario1_DuplicateScenario2})

	n_users := 3
	for n := 1; n <= n_users; n++ {
		//Add user to usergroup
		n := strconv.Itoa(n)
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

	//Actually check whether the users are in the scenarios
	assert.Equal(t, 2+n_users, len(scenarios)) // 2 original scenarios + 1 duplicate for each user

	for _, scenario := range scenarios {
		path := fmt.Sprintf("/api/v2/scenarios/%d/users", scenario.ID)
		var usersMap map[string]([]database.User)
		_, res, _ = helper.TestEndpoint(router, token, path, "GET", struct{}{})
		json.Unmarshal(res.Bytes(), &usersMap)
		users := usersMap["users"]
		switch scenario.ID {
		case 1: // Scenario "scenarioNoDups"
			assert.Equal(t, "scenarioNoDups", scenario.Name)
			assert.Equal(t, 4, len(users)) // 3 users + admin
		case 2: // Scenario "scenarioDups"
			assert.Equal(t, "scenarioDups", scenario.Name)
			assert.Equal(t, 1, len(users))
			assert.Equal(t, "admin", users[0].Username)
		default: // duplicated scenarios
			usr := "usr" + strconv.Itoa(int(scenario.ID-2)) // shift ids by the first two scenarios
			assert.Equal(t, "scenarioDups "+usr, scenario.Name)
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
	helper.TestEndpoint(router, token, "/api/v2/usergroups", "POST", helper.KeyModels{"usergroup": ug_AddScenario1_DuplicateScenario2})

	// Create users and add them to the scenarios via usergroup
	n_users := 2
	for n := 1; n <= n_users; n++ {
		//Add user
		n := strconv.Itoa(n)
		usr := UserRequest{Username: "usr" + n, Password: "legendre" + n, Role: "User", Mail: "usr" + n + "@harmonics.de"}
		helper.TestEndpoint(router, token, "/api/v2/users", "POST", helper.KeyModels{"user": usr})
		code, _, err := helper.TestEndpoint(router, token, "/api/v2/usergroups/1/user?username=usr"+n, "PUT", struct{}{})
		assert.Equal(t, 200, code)
		assert.NoError(t, err)
	}
	// add usr1 to usergroup 2, which doubles the access to scenario 1
	helper.TestEndpoint(router, token, "/api/v2/usergroups", "POST", helper.KeyModels{"usergroup": ug_AddScenario1})
	helper.TestEndpoint(router, token, "/api/v2/usergroups/2/user?username=usr1", "PUT", struct{}{})

	//Delete user usr1 from usergroup 1, this will also delete the user's duplicate of scenario 2
	helper.TestEndpoint(router, token, "/api/v2/usergroups/1/user?username=usr1", "DELETE", struct{}{})

	//get all scenarios
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
		case 1: // scenario "scenarioNoDups" should still contain usr1 through usergroup 2
			assert.Equal(t, "scenarioNoDups", v.Name)
			assert.Equal(t, 3, len(users))
		case 2: // scenario "scenarioDups" should only contain admin, gets duplicated for users in usergroup 1
			assert.Equal(t, "scenarioDups", v.Name)
			assert.Equal(t, 1, len(users))
			assert.Equal(t, "admin", users[0].Username)
		default: // remaining duplicated scenario
			assert.Equal(t, "scenarioDups usr2", v.Name)
			assert.Equal(t, 1, len(users))
			assert.Equal(t, "usr2", users[0].Username)
		}
	}
	// delete usr1 from usergroup 2
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

	// make sure that usr1 is not in the scenario anymore
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

	//we add usr1 to a group that doubles its right of access to scenario 1
	helper.TestEndpoint(router, token, "/api/v2/usergroups", "POST", helper.KeyModels{"usergroup": ug_AddScenario1})
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
			assert.Equal(t, 2, len(users)) // admin and usr1
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

func TestUpdateUserGroup(t *testing.T) {
	// Prep DB
	database.DropTables()
	database.MigrateModels()
	adminpw, _ := database.AddAdminUser(configuration.GlobalConfig)

	//Auth
	token, _ := helper.AuthenticateForTest(router, database.Credentials{Username: "admin", Password: adminpw})

	//Post necessities
	helper.TestEndpoint(router, token, "/api/v2/scenarios", "POST", helper.KeyModels{"scenario": database.Scenario{Name: "changeDups"}})
	helper.TestEndpoint(router, token, "/api/v2/scenarios", "POST", helper.KeyModels{"scenario": database.Scenario{Name: "changeNoDups"}})
	helper.TestEndpoint(router, token, "/api/v2/scenarios", "POST", helper.KeyModels{"scenario": database.Scenario{Name: "removeDups"}})
	helper.TestEndpoint(router, token, "/api/v2/scenarios", "POST", helper.KeyModels{"scenario": database.Scenario{Name: "removeNoDups"}})
	helper.TestEndpoint(router, token, "/api/v2/scenarios", "POST", helper.KeyModels{"scenario": database.Scenario{Name: "unchangedDups"}})
	helper.TestEndpoint(router, token, "/api/v2/scenarios", "POST", helper.KeyModels{"scenario": database.Scenario{Name: "unchangedNoDups"}})
	helper.TestEndpoint(router, token, "/api/v2/scenarios", "POST", helper.KeyModels{"scenario": database.Scenario{Name: "addDups"}})
	helper.TestEndpoint(router, token, "/api/v2/scenarios", "POST", helper.KeyModels{"scenario": database.Scenario{Name: "addNoDups"}})
	helper.TestEndpoint(router, token, "/api/v2/usergroups", "POST", helper.KeyModels{"usergroup": initUpdateTestUg})

	//Add user
	usr := UserRequest{Username: "usr1", Password: "legendre1", Role: "User", Mail: "usr1@harmonics.de"}
	helper.TestEndpoint(router, token, "/api/v2/users", "POST", helper.KeyModels{"user": usr})
	helper.TestEndpoint(router, token, "/api/v2/usergroups/1/user?username=usr1", "PUT", struct{}{})

	// update group
	helper.TestEndpoint(router, token, "/api/v2/usergroups/1", "PUT", helper.KeyModels{"usergroup": updateTestUg})

	//get scenarios
	_, res, _ := helper.TestEndpoint(router, token, "/api/v2/scenarios", "GET", struct{}{})
	var scenariosRes map[string]([]database.Scenario)
	json.Unmarshal(res.Bytes(), &scenariosRes)
	scenarios := scenariosRes["scenarios"]
	assert.Equal(t, 11, len(scenarios))
	var scenariosMap map[string](database.Scenario) = make(map[string](database.Scenario))
	for _, s := range scenarios {
		scenariosMap[s.Name] = s
	}

	//scenarios that transformed into/remained/got added as duplicated (6)
	for _, name := range []string{"addDups", "unchangedDups", "changeNoDups"} {
		sc, exists := scenariosMap[name]
		assert.True(t, exists)
		path := fmt.Sprintf("/api/v2/scenarios/%d/users", sc.ID)
		var usersMap map[string]([]database.User)
		_, res, _ = helper.TestEndpoint(router, token, path, "GET", struct{}{})
		json.Unmarshal(res.Bytes(), &usersMap)
		users := usersMap["users"]
		assert.Equal(t, 1, len(users))

		sc, exists = scenariosMap[name+" usr1"]
		path = fmt.Sprintf("/api/v2/scenarios/%d/users", sc.ID)
		assert.True(t, exists)
		_, res, _ = helper.TestEndpoint(router, token, path, "GET", struct{}{})
		json.Unmarshal(res.Bytes(), &usersMap)
		users = usersMap["users"]
		assert.Equal(t, 1, len(users))
	}

	//scenarios that transformed into/remained/got added as single (+3)
	for _, name := range []string{"addNoDups", "unchangedNoDups", "changeDups"} {
		sc, exists := scenariosMap[name]
		assert.True(t, exists)
		path := fmt.Sprintf("/api/v2/scenarios/%d/users", sc.ID)
		var usersMap map[string]([]database.User)
		_, res, _ = helper.TestEndpoint(router, token, path, "GET", struct{}{})
		json.Unmarshal(res.Bytes(), &usersMap)
		users := usersMap["users"]
		assert.Equal(t, 2, len(users))

		_, exists = scenariosMap[name+" usr1"]
		assert.False(t, exists)
	}

	//scenarios that got removed (+2 = 11)
	for _, name := range []string{"removeDups", "removeNoDups"} {
		sc, exists := scenariosMap[name]
		assert.True(t, exists)
		path := fmt.Sprintf("/api/v2/scenarios/%d/users", sc.ID)
		var usersMap map[string]([]database.User)
		_, res, _ = helper.TestEndpoint(router, token, path, "GET", struct{}{})
		json.Unmarshal(res.Bytes(), &usersMap)
		users := usersMap["users"]
		assert.Equal(t, 1, len(users))

		_, exists = scenariosMap[name+" usr1"]
		assert.False(t, exists)
	}

}
