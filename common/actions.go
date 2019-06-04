package common

// The types ModelName and CRUD exist in order to have const variables
// that will be used by every model's package. That way we avoid
// hardcoding the same names all over and spelling mistakes that will
// lead IsActionAllowed to fail. We do not typedef the different roles
// because they will be read as strings from the context.

type ModelName string

const ModelUser = ModelName("user")
const ModelSimulation = ModelName("simulation")
const ModelSimulator = ModelName("simulator")
const ModelVisualization = ModelName("visualization")

type CRUD string

const Create = CRUD("create")
const Read = CRUD("read")
const Update = CRUD("update")
const Delete = CRUD("delete")

// Type Modes maps a CRUD operation to true or false
type Modes map[CRUD]bool

// Type ModelActions maps a model name to a map of Modes for every model
type ModelActions map[ModelName]Modes

// Type RoleActions maps a role to a map of ModelActions for every role
type RoleActions map[string]ModelActions

// Predefined CRUD operations permissions to be used in Roles
var crud = Modes{Create: true, Read: true, Update: true, Delete: true}
var _ru_ = Modes{Create: false, Read: true, Update: true, Delete: false}
var _r__ = Modes{Create: false, Read: true, Update: false, Delete: false}

// Roles is used as a look up variable to determine if a certain user is
// allowed to do a certain action on a given model based on his role
var Roles = RoleActions{
	"Admin": {
		ModelUser:       crud,
		ModelSimulation: crud,
		ModelSimulator:  crud,
	},
	"User": {
		ModelUser:       _ru_,
		ModelSimulation: crud,
		ModelSimulator:  _r__,
	},
	"Guest": {
		ModelVisualization: _r__,
	},
}
