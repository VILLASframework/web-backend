package common

import (
	"time"

	"github.com/jinzhu/gorm/dialects/postgres"
)

// The type Model is exactly the same with gorm.Model (see jinzhu/gorm)
// except the json tags that are needed for serializing the models
type Model struct {
	ID        uint       `json:"id,omitempty" gorm:"unique;primary_key:true"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `json:"-" sql:"index"`
}

// User data model
type User struct {
	Model
	// Username of user
	Username string `json:"username" gorm:"unique;not null"`
	// Password of user
	Password string `json:"-" gorm:"not null"`
	// Mail of user
	Mail string `json:"mail" gorm:"default:''"`
	// Role of user
	Role string `json:"role" gorm:"default:'user'"`
	// Scenarios to which user has access
	Scenarios []*Scenario `json:"-" gorm:"many2many:user_scenarios"`
}

// Scenario data model
type Scenario struct {
	Model
	// Name of scenario
	Name string `gorm:"not null",json:"name"`
	// Running state of scenario
	Running bool `gorm:"default:false",json:"running"`
	// Start parameters of scenario as JSON
	StartParameters postgres.Jsonb `json:"startParameters"`
	// Users that have access to the scenario
	Users []*User `gorm:"not null;many2many:user_scenarios",json:"-"`
	// SimulationModels that belong to the scenario
	SimulationModels []SimulationModel `gorm:"foreignkey:ScenarioID",json:"-"`
	// Dashboards that belong to the Scenario
	Dashboards []Dashboard `gorm:"foreignkey:ScenarioID",json:"-"`
}

// SimulationModel data model
type SimulationModel struct {
	Model
	// Name of simulation model
	Name string `gorm:"not null",json:"name"`
	// Number of output signals
	OutputLength int `gorm:"default:1",json:"outputLength"`
	// Number of input signals
	InputLength int `gorm:"default:1",json:"inputLength"`
	// Start parameters of simulation model as JSON
	StartParameters postgres.Jsonb `json:"startParameters"`
	// ID of Scenario to which simulation model belongs
	ScenarioID uint `json:"scenarioID"`
	// ID of simulator associated with simulation model
	SimulatorID uint `json:"simulatorID"`
	// Mapping of output signals of the simulation model, order of signals is important
	OutputMapping []Signal `gorm:"foreignkey:SimulationModelID",json:"-"`
	// Mapping of input signals of the simulation model, order of signals is important
	InputMapping []Signal `gorm:"foreignkey:SimulationModelID",json:"-"`
	// Files of simulation model (can be CIM and other simulation model file formats)
	Files []File `gorm:"foreignkey:SimulationModelID",json:"-"`
}

// Signal data model
type Signal struct {
	Model
	// Name of Signal
	Name string `json:"name"`
	// Unit of Signal
	Unit string `json:"unit"`
	// Index of the Signal in the mapping
	Index uint `json:"index"`
	// Direction of the signal (in or out)
	Direction string `json:"direction"`
	// ID of simulation model
	SimulationModelID uint `json:"simulationModelID"`
}

// Simulator data model
type Simulator struct {
	Model
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
	SimulationModels []SimulationModel `gorm:"foreignkey:SimulatorID",json:"-"`
}

// Dashboard data model
type Dashboard struct {
	Model
	// Name of dashboard
	Name string `gorm:"not null",json:"name"`
	// Grid of dashboard
	Grid int `gorm:"default:15",json:"grid"`
	// ID of scenario to which dashboard belongs
	ScenarioID uint `json:"scenarioID"`
	// Widgets that belong to dashboard
	Widgets []Widget `gorm:"foreignkey:DashboardID",json:"-"`
}

// Widget data model
type Widget struct {
	Model
	// Name of widget
	Name string `gorm:"not null",json:"name"`
	// Type of widget
	Type string `gorm:"not null",json:"type"`
	// Width of widget
	Width uint `gorm:"not null",json:"width"`
	// Height of widget
	Height uint `gorm:"not null",json:"height"`
	// Minimal width of widget
	MinWidth uint `gorm:"not null",json:"minWidth"`
	// Minimal height of widget
	MinHeight uint `gorm:"not null",json:"minHeight"`
	// X position of widget
	X int `gorm:"not null",json:"x"`
	// Y position of widget
	Y int `gorm:"not null",json:"y"`
	// Z position of widget
	Z int `gorm:"not null",json:"z"`
	// Locked state of widget
	IsLocked bool `gorm:"default:false",json:"isLocked"`
	// Custom properties of widget as JSON string
	CustomProperties postgres.Jsonb `json:"customProperties"`
	// ID of dashboard to which widget belongs
	DashboardID uint `json:"dashboardID"`
	// Files that belong to widget (for example images)
	Files []File `gorm:"foreignkey:WidgetID",json:"-"`
}

// File data model
type File struct {
	Model
	// Name of file
	Name string `gorm:"not null",json:"name"`
	// Type of file (MIME type)
	Type string `json:"type"`
	// Size of file (in byte)
	Size uint `json:"size"`
	// Height of image (only needed in case of image)
	ImageHeight uint `json:"imageHeight"`
	// Width of image (only needed in case of image)
	ImageWidth uint `json:"imageWidth"`
	// Last modification time of file
	Date string `json:"date"`
	// ID of model to which file belongs
	SimulationModelID uint `json:"simulationModelID"`
	// ID of widget to which file belongs
	WidgetID uint `json:"widgetID"`
	// File itself
	FileData []byte `gorm:"column:FileData",json:"-"`
}
