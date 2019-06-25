package common

import (
	"time"
)

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
	// Simulations to which user has access
	Simulations []*Simulation `gorm:"many2many:user_simulations"`
}

// Simulation data model
type Simulation struct {
	// ID of simulation
	ID uint `gorm:"primary_key;auto_increment"`
	// Name of simulation
	Name string `gorm:"not null"`
	// Running state of simulation
	Running bool `gorm:"default:false"`
	// Start parameters of simulation as JSON string
	StartParameters string
	// Users that have access to the simulation
	Users []*User `gorm:"not null;many2many:user_simulations"`
	// SimulationModels that belong to the simulation
	SimulationModels []SimulationModel `gorm:"foreignkey:SimulationID"`
	// Visualizations that belong to the simulation
	Visualizations []Visualization `gorm:"foreignkey:SimulationID"`
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
	// Start parameters of simulation model as JSON string
	StartParameters string
	// ID of simulation to which simulation model belongs
	SimulationID uint
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
	ID uint `gorm:"primary_key;auto_increment"`
	// UUID of the simulator
	UUID string `gorm:"unique;not null"`
	// Host if the simulator
	Host string `gorm:"default:''"`
	// Model type supported by the simulator
	Modeltype string `gorm:"default:''"`
	// Uptime of the simulator
	Uptime int `gorm:"default:0"`
	// State of the simulator
	State string `gorm:"default:''"`
	// Time of last state update
	StateUpdateAt string `gorm:"default:''"`
	// Properties of simulator as JSON string
	Properties string `gorm:"default:''"`
	// Raw properties of simulator as JSON string
	RawProperties string `gorm:"default:''"`
	// SimulationModels in which the simulator is used
	SimulationModels []SimulationModel `gorm:"foreignkey:SimulatorID"`
}

// Visualization data model
type Visualization struct {
	// ID of visualization
	ID uint `gorm:"primary_key;auto_increment"`
	// Name of visualization
	Name string `gorm:"not null"`
	// Grid of visualization
	Grid int `gorm:"default:15"`
	// ID of simulation to which visualization belongs
	SimulationID uint
	// Widgets that belong to visualization
	Widgets []Widget `gorm:"foreignkey:VisualizationID"`
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
	CustomProperties string
	// ID of visualization to which widget belongs
	VisualizationID uint
	// Files that belong to widget (for example images)
	Files []File `gorm:"foreignkey:WidgetID"`
}

// File data model
type File struct {
	// ID of file
	ID uint `gorm:"primary_key;auto_increment"`
	// Name of file
	Name string `gorm:"not null"`
	// Path at which file is saved at server side
	Path string `gorm:"not null"`
	// Type of file (MIME type)
	Type string `gorm:"not null"`
	// Size of file (in byte)
	Size uint `gorm:"not null"`
	// Height of image (only needed in case of image)
	ImageHeight uint
	// Width of image (only needed in case of image)
	ImageWidth uint
	// Last modification time of file
	Date time.Time
	// ID of model to which file belongs
	SimulationModelID uint `gorm:""`
	// ID of widget to which file belongs
	WidgetID uint `gorm:""`
}
