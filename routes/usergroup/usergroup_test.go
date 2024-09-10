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
	"os"
	"testing"

	"git.rwth-aachen.de/acs/public/villas/web-backend-go/configuration"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/helper"
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

	// Test with valid user group with one scenario mapping
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
