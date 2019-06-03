package common

// Type Modes maps a CRUD operation to true or false
type Modes map[string]bool

// Type ModelActions maps a model to a map of Modes for every model
type ModelActions map[string]Modes

// Type RoleActions maps a role to a map of ModelActions for every role
type RoleActions map[string]ModelActions

// Predefined CRUD operations permissions to be used in Roles
var crud = Modes{"create": true, "read": true, "update": true, "delete": true}
var cru_ = Modes{"create": true, "read": true, "update": true, "delete": false}
var _r__ = Modes{"create": false, "read": true, "update": false, "delete": false}

// Roles is used as a look up variable to determine if a certain user is
// allowed to do a certain action on a given model based on his role
var Roles = RoleActions{
	"Admin": {
		"user":       crud,
		"simulation": crud,
		"simulator":  crud,
	},
	"User": {
		"user":       cru_,
		"simulation": crud,
		"simulator":  _r__,
	},
	"Guest": {
		"visualization": _r__,
	},
}
