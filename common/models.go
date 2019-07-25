package common

import "github.com/jinzhu/gorm/dialects/postgres"

// User data model
type User struct {
	// ID of user
	ID uint `gorm:"primary_key;auto_increment"`
	// Username of user
	Username string `gorm:"unique;not null"`
	// Password of user
	Password string `gorm:"not null"`
	// Mail of user
	Mail string `gorm:"default:''"`
	// Role of user
	Role string `gorm:"default:'user'"`
	// Scenarios to which user has access
	Scenarios []*Scenario `gorm:"many2many:user_scenarios"`
}

// Scenario data model
type Scenario struct {
	// ID of scenario
	ID uint `gorm:"primary_key;auto_increment"`
	// Name of scenario
	Name string `gorm:"not null"`
	// Running state of scenario
	Running bool `gorm:"default:false"`
	// Start parameters of scenario as JSON
	StartParameters postgres.Jsonb
	// Users that have access to the scenario
	Users []*User `gorm:"not null;many2many:user_scenarios"`
	// SimulationModels that belong to the scenario
	SimulationModels []SimulationModel `gorm:"foreignkey:ScenarioID"`
	// Dashboards that belong to the Scenario
	Dashboards []Dashboard `gorm:"foreignkey:ScenarioID"`
}

// SimulationModel data model
type SimulationModel struct {
	// ID of simulation model
	ID uint `gorm:"primary_key;auto_increment"`
	// Name of simulation model
	Name string `gorm:"not null"`
	// Number of output signals
	OutputLength int `gorm:"default:1"`
	// Number of input signals
	InputLength int `gorm:"default:1"`
	// Start parameters of simulation model as JSON
	StartParameters postgres.Jsonb
	// ID of Scenario to which simulation model belongs
	ScenarioID uint
	// ID of simulator associated with simulation model
	SimulatorID uint
	// Mapping of output signals of the simulation model, order of signals is important
	OutputMapping []Signal `gorm:"foreignkey:SimulationModelID"`
	// Mapping of input signals of the simulation model, order of signals is important
	InputMapping []Signal `gorm:"foreignkey:SimulationModelID"`
	// Files of simulation model (can be CIM and other simulation model file formats)
	Files []File `gorm:"foreignkey:SimulationModelID"`
}

// Signal data model
type Signal struct {
	// ID of simulation model
	ID uint `gorm:"primary_key;auto_increment"`
	// Name of Signal
	Name string
	// Unit of Signal
	Unit string
	// Index of the Signal in the mapping
	Index uint
	// Direction of the signal (in or out)
	Direction string
	// ID of simulation model
	SimulationModelID uint
}

// Simulator data model
type Simulator struct {
	// ID of the simulator
	ID uint `gorm:"primary_key;auto_increment",json:"id"`
	// UUID of the simulator
	UUID string `gorm:"not null",json:"uuid"`
	// Host if the simulator
	Host string `gorm:"default:''",json:"host"`
	// Model type supported by the simulator
	Modeltype string `gorm:"default:''",json:"modelType"`
	// Uptime of the simulator
	Uptime int `gorm:"default:0",json:"uptime"`
	// State of the simulator
	State string `gorm:"default:''",json:"state"`
	// Time of last state update
	StateUpdateAt string `gorm:"default:''",json:"stateUpdateAt"`
	// Properties of simulator as JSON string
	Properties postgres.Jsonb `json:"properties"`
	// Raw properties of simulator as JSON string
	RawProperties postgres.Jsonb `json:"rawProperties"`
	// SimulationModels in which the simulator is used
	SimulationModels []SimulationModel `gorm:"foreignkey:SimulatorID"`
}

// Dashboard data model
type Dashboard struct {
	// ID of dashboard
	ID uint `gorm:"primary_key;auto_increment"`
	// Name of dashboard
	Name string `gorm:"not null"`
	// Grid of dashboard
	Grid int `gorm:"default:15"`
	// ID of scenario to which dashboard belongs
	ScenarioID uint
	// Widgets that belong to dashboard
	Widgets []Widget `gorm:"foreignkey:DashboardID"`
}

// Widget data model
type Widget struct {
	// ID of widget
	ID uint `gorm:"primary_key;auto_increment"`
	// Name of widget
	Name string `gorm:"not null"`
	// Type of widget
	Type string `gorm:"not null"`
	// Width of widget
	Width uint `gorm:"not null"`
	// Height of widget
	Height uint `gorm:"not null"`
	// Minimal width of widget
	MinWidth uint `gorm:"not null"`
	// Minimal height of widget
	MinHeight uint `gorm:"not null"`
	// X position of widget
	X int `gorm:"not null"`
	// Y position of widget
	Y int `gorm:"not null"`
	// Z position of widget
	Z int `gorm:"not null"`
	// Locked state of widget
	IsLocked bool `gorm:"default:false"`
	// Custom properties of widget as JSON string
	CustomProperties postgres.Jsonb
	// ID of dashboard to which widget belongs
	DashboardID uint
	// Files that belong to widget (for example images)
	Files []File `gorm:"foreignkey:WidgetID"`
}

// File data model
type File struct {
	// ID of file
	ID uint `gorm:"primary_key;auto_increment"`
	// Name of file
	Name string `gorm:"not null"`
	// Type of file (MIME type)
	Type string `gorm:"not null"`
	// Size of file (in byte)
	Size uint `gorm:"not null"`
	// Height of image (only needed in case of image)
	ImageHeight uint
	// Width of image (only needed in case of image)
	ImageWidth uint
	// Last modification time of file
	Date string
	// ID of model to which file belongs
	SimulationModelID uint
	// ID of widget to which file belongs
	WidgetID uint
	// TODO Add file itself here??
	FileData []byte `gorm:"column:FileData"`
}

// Credentials type (not for DB)
type credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
