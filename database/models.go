/** Database package, models.
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
	"time"

	"github.com/jinzhu/gorm/dialects/postgres"
)

// The type Model is exactly the same with gorm.Model (see jinzhu/gorm)
// except the json tags that are needed for serializing the models
type Model struct {
	ID        uint       `json:"id,omitempty" gorm:"primary_key:true"`
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
	// Indicating status of user (false means user is inactive and should not be able to login)
	Active bool `json:"active" gorm:"default:true"`
	// Scenarios to which user has access
	Scenarios []*Scenario `json:"-" gorm:"many2many:user_scenarios;"`
}

// Scenario data model
type Scenario struct {
	Model
	// Name of scenario
	Name string `json:"name" gorm:"not null"`
	// Running state of scenario
	Running bool `json:"running" gorm:"default:false" `
	// Start parameters of scenario as JSON
	StartParameters postgres.Jsonb `json:"startParameters"`
	// Users that have access to the scenario
	Users []*User `json:"-" gorm:"many2many:user_scenarios;"`
	// SimulationModels that belong to the scenario
	SimulationModels []SimulationModel `json:"-" gorm:"foreignkey:ScenarioID" `
	// Dashboards that belong to the Scenario
	Dashboards []Dashboard `json:"-" gorm:"foreignkey:ScenarioID" `
}

// SimulationModel data model
type SimulationModel struct {
	Model
	// Name of simulation model
	Name string `json:"name" gorm:"not null"`
	// Number of output signals
	OutputLength int `json:"outputLength" gorm:"default:0"`
	// Number of input signals
	InputLength int `json:"inputLength" gorm:"default:0"`
	// Start parameters of simulation model as JSON
	StartParameters postgres.Jsonb `json:"startParameters"`
	// ID of Scenario to which simulation model belongs
	ScenarioID uint `json:"scenarioID"`
	// ID of simulator associated with simulation model
	SimulatorID uint `json:"simulatorID"`
	// Mapping of output signals of the simulation model, order of signals is important
	OutputMapping []Signal `json:"-" gorm:"foreignkey:SimulationModelID"`
	// Mapping of input signals of the simulation model, order of signals is important
	InputMapping []Signal `json:"-" gorm:"foreignkey:SimulationModelID"`
	// Files of simulation model (can be CIM and other simulation model file formats)
	Files []File `json:"-" gorm:"foreignkey:SimulationModelID"`
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
	UUID string `json:"uuid" gorm:"not null"`
	// Host if the simulator
	Host string `json:"host" gorm:"default:''"`
	// Model type supported by the simulator
	Modeltype string `json:"modelType" gorm:"default:''"`
	// Uptime of the simulator
	Uptime int `json:"uptime" gorm:"default:0"`
	// State of the simulator
	State string `json:"state" gorm:"default:''"`
	// Time of last state update
	StateUpdateAt string `json:"stateUpdateAt" gorm:"default:''"`
	// Properties of simulator as JSON string
	Properties postgres.Jsonb `json:"properties"`
	// Raw properties of simulator as JSON string
	RawProperties postgres.Jsonb `json:"rawProperties"`
	// SimulationModels in which the simulator is used
	SimulationModels []SimulationModel `json:"-" gorm:"foreignkey:SimulatorID"`
}

// Dashboard data model
type Dashboard struct {
	Model
	// Name of dashboard
	Name string `json:"name" gorm:"not null"`
	// Grid of dashboard
	Grid int `json:"grid" gorm:"default:15"`
	// ID of scenario to which dashboard belongs
	ScenarioID uint `json:"scenarioID"`
	// Widgets that belong to dashboard
	Widgets []Widget `json:"-" gorm:"foreignkey:DashboardID"`
}

// Widget data model
type Widget struct {
	Model
	// Name of widget
	Name string `json:"name" gorm:"not null"`
	// Type of widget
	Type string `json:"type" gorm:"not null"`
	// Width of widget
	Width uint `json:"width" gorm:"not null"`
	// Height of widget
	Height uint `json:"height" gorm:"not null"`
	// Minimal width of widget
	MinWidth uint `json:"minWidth" gorm:"not null"`
	// Minimal height of widget
	MinHeight uint `json:"minHeight" gorm:"not null"`
	// X position of widget
	X int `json:"x" gorm:"not null"`
	// Y position of widget
	Y int `json:"y" gorm:"not null"`
	// Z position of widget
	Z int `json:"z" gorm:"not null"`
	// Locked state of widget
	IsLocked bool `json:"isLocked" gorm:"default:false"`
	// Custom properties of widget as JSON string
	CustomProperties postgres.Jsonb `json:"customProperties"`
	// ID of dashboard to which widget belongs
	DashboardID uint `json:"dashboardID"`
	// Files that belong to widget (for example images)
	Files []File `json:"-" gorm:"foreignkey:WidgetID"`
}

// File data model
type File struct {
	Model
	// Name of file
	Name string `json:"name" gorm:"not null"`
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
	FileData []byte `json:"-" gorm:"column:FileData"`
}
