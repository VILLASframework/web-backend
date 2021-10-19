/** Database package, roles.
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
package database

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// The types ModelName and CRUD exist in order to have const variables
// that will be used by every model's package. That way we avoid
// hardcoding the same names all over and spelling mistakes that will
// lead IsActionAllowed to fail. We do not typedef the different roles
// because they will be read as strings from the context.

type ModelName string

const ModelUser = ModelName("user")
const ModelUsers = ModelName("users")
const ModelScenario = ModelName("scenario")
const ModelInfrastructureComponent = ModelName("ic")
const ModelInfrastructureComponentAction = ModelName("icaction")
const ModelDashboard = ModelName("dashboard")
const ModelWidget = ModelName("widget")
const ModelComponentConfiguration = ModelName("component-configuration")
const ModelSignal = ModelName("signal")
const ModelFile = ModelName("file")
const ModelResult = ModelName("result")

type CRUD string

const Create = CRUD("create")
const Read = CRUD("read")
const Update = CRUD("update")
const Delete = CRUD("delete")

// Type Permission maps a CRUD operation to true or false
type Permission map[CRUD]bool

// Type ModelActions maps a model name to a map of Permission for every model
type ModelActions map[ModelName]Permission

// Type RoleActions maps a role to a map of ModelActions for every role
type RoleActions map[string]ModelActions

// Predefined CRUD operations permissions to be used in Roles
var crud = Permission{Create: true, Read: true, Update: true, Delete: true}
var _ru_ = Permission{Create: false, Read: true, Update: true, Delete: false}
var _r__ = Permission{Create: false, Read: true, Update: false, Delete: false}
var none = Permission{Create: false, Read: false, Update: false, Delete: false}

// var __u_ = Permission{Create: false, Read: false, Update: true, Delete: false}

// Roles is used as a look up variable to determine if a certain user is
// allowed to do a certain action on a given model based on his role
var Roles = RoleActions{
	"Admin": {
		ModelUser:                          crud,
		ModelUsers:                         crud,
		ModelScenario:                      crud,
		ModelComponentConfiguration:        crud,
		ModelInfrastructureComponent:       crud,
		ModelInfrastructureComponentAction: crud,
		ModelWidget:                        crud,
		ModelDashboard:                     crud,
		ModelSignal:                        crud,
		ModelFile:                          crud,
		ModelResult:                        crud,
	},
	"User": {
		ModelUser:                          _ru_,
		ModelUsers:                         none,
		ModelScenario:                      crud,
		ModelComponentConfiguration:        crud,
		ModelInfrastructureComponent:       _r__,
		ModelInfrastructureComponentAction: _ru_,
		ModelWidget:                        crud,
		ModelDashboard:                     crud,
		ModelSignal:                        crud,
		ModelFile:                          crud,
		ModelResult:                        crud,
	},
	"Guest": {
		ModelScenario:                      _r__,
		ModelComponentConfiguration:        _r__,
		ModelDashboard:                     _r__,
		ModelWidget:                        _r__,
		ModelInfrastructureComponent:       _r__,
		ModelInfrastructureComponentAction: _r__,
		ModelUser:                          _ru_,
		ModelUsers:                         none,
		ModelSignal:                        _r__,
		ModelFile:                          _r__,
		ModelResult:                        none,
	},
	"Download": {
		ModelScenario:                      none,
		ModelComponentConfiguration:        none,
		ModelDashboard:                     none,
		ModelWidget:                        none,
		ModelInfrastructureComponent:       none,
		ModelInfrastructureComponentAction: none,
		ModelUser:                          none,
		ModelUsers:                         none,
		ModelSignal:                        none,
		ModelFile:                          _r__,
		ModelResult:                        none,
	},
}

func ValidateRole(c *gin.Context, model ModelName, action CRUD) error {
	// Extracts and validates the role which is saved in the context for
	// executing a specific CRUD operation on a specific model. In case
	// of invalid role return an error.

	// Get user's role from context
	role, exists := c.Get(UserRoleCtx)
	if !exists {
		return fmt.Errorf("request does not contain user's role")
	}

	// Check if the role can execute the action on the model
	if !Roles[role.(string)][model][action] {
		return fmt.Errorf("action not allowed for role %v", role)
	}

	return nil
}

// elements added to context about a user
const UserIDCtx = "user_id"
const UserRoleCtx = "user_role"
