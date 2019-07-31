package common

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
const ModelSimulator = ModelName("simulator")
const ModelSimulatorAction = ModelName("simulatoraction")
const ModelDashboard = ModelName("dashboard")
const ModelWidget = ModelName("widget")
const ModelSimulationModel = ModelName("simulationmodel")
const ModelSignal = ModelName("signal")
const ModelFile = ModelName("file")

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
var __u_ = Permission{Create: false, Read: false, Update: true, Delete: false}
var _r__ = Permission{Create: false, Read: true, Update: false, Delete: false}
var none = Permission{Create: false, Read: false, Update: false, Delete: false}

// Roles is used as a look up variable to determine if a certain user is
// allowed to do a certain action on a given model based on his role
var Roles = RoleActions{
	"Admin": {
		ModelUser:            crud,
		ModelUsers:           crud,
		ModelScenario:        crud,
		ModelSimulationModel: crud,
		ModelSimulator:       crud,
		ModelSimulatorAction: crud,
		ModelWidget:          crud,
		ModelDashboard:       crud,
		ModelSignal:          crud,
		ModelFile:            crud,
	},
	"User": {
		ModelUser:            _ru_,
		ModelUsers:           none,
		ModelScenario:        crud,
		ModelSimulationModel: crud,
		ModelSimulator:       _r__,
		ModelSimulatorAction: _ru_,
		ModelWidget:          crud,
		ModelDashboard:       crud,
		ModelSignal:          crud,
		ModelFile:            crud,
	},
	"Guest": {
		ModelScenario:        _r__,
		ModelSimulationModel: _r__,
		ModelDashboard:       _r__,
		ModelWidget:          _r__,
		ModelSimulator:       _r__,
		ModelSimulatorAction: _r__,
		ModelUser:            _ru_,
		ModelUsers:           none,
		ModelSignal:          _r__,
		ModelFile:            _r__,
	},
}

func ValidateRole(c *gin.Context, model ModelName, action CRUD) error {
	// Extracts and validates the role which is saved in the context for
	// executing a specific CRUD operation on a specific model. In case
	// of invalid role return an error.

	// Get user's role from context
	role, exists := c.Get("user_role")
	if !exists {
		return fmt.Errorf("Request does not contain user's role")
	}

	// Check if the role can execute the action on the model
	if !Roles[role.(string)][model][action] {
		return fmt.Errorf("Action not allowed for role %v", role)
	}

	return nil
}
