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

package database

import (
	"time"

	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/lib/pq"
)

// The type Model is exactly the same with gorm.Model (see jinzhu/gorm)
// except the json tags that are needed for serializing the models
type Model struct {
	ID        uint       `json:"id,omitempty" gorm:"primary_key:true"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
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
	// Groups of user
	UserGroups []*UserGroup `json:"-" gorm:"many2many:user_groups_users;"`
}

// ScenarioMapping data model
type ScenarioMapping struct {
	Model
	ScenarioID  uint      `json:"scenarioID" ` // Foreign key to Scenario
	Scenario    Scenario  `json:"-" gorm:"foreignkey:ScenarioID"`
	UserGroupID uint      `json:"-"` // Foreign key to UserGroup
	UserGroup   UserGroup `json:"-" gorm:"foreignkey:UserGroupID"`
	// Whether to duplicate Scenario or add users to existing Scenario
	Duplicate bool `json:"duplicate"`
}

// UserGroup data model
type UserGroup struct {
	Model
	// Name of user group
	Name string `json:"name" gorm:"unique;not null"`
	// Scenarios that belong to the user group
	ScenarioMappings []ScenarioMapping `json:"scenarioMappings" gorm:"foreignkey:UserGroupID"`
	// Users that belong to the user group
	Users []*User `json:"users" gorm:"many2many:user_groups_users;"`
}

// Scenario data model
type Scenario struct {
	Model
	// Name of scenario
	Name string `json:"name" gorm:"not null"`
	// IsLocked state of scenario (true if scenario is locked by administrator)
	IsLocked bool `json:"isLocked" gorm:"default:false" `
	// Start parameters of scenario as JSON
	StartParameters postgres.Jsonb `json:"startParameters"`
	// Users that have access to the scenario
	Users []*User `json:"-" gorm:"many2many:user_scenarios;"`
	// ComponentConfigurations that belong to the scenario
	ComponentConfigurations []ComponentConfiguration `json:"-" gorm:"foreignkey:ScenarioID" `
	// Dashboards that belong to the Scenario
	Dashboards []Dashboard `json:"-" gorm:"foreignkey:ScenarioID" `
	// Files that belong to the Scenario (for example images, models, etc.)
	Files []File `json:"-" gorm:"foreignkey:ScenarioID"`
	// Results that belong to the Scenario
	Results []Result `json:"-" gorm:"foreignkey:ScenarioID"`
}

// ComponentConfiguration data model
type ComponentConfiguration struct {
	Model
	// Name of Component Configuration
	Name string `json:"name" gorm:"not null"`
	// Start parameters of Component Configuration as JSON
	StartParameters postgres.Jsonb `json:"startParameters"`
	// ID of Scenario to which Component Configuration belongs
	ScenarioID uint `json:"scenarioID"`
	// ID of IC associated with Component Configuration
	ICID uint `json:"icID"`
	// Mapping of output signals of the ComponentConfiguration, order of signals is important
	OutputMapping []Signal `json:"-" gorm:"foreignkey:ConfigID"`
	// Mapping of input signals of the Component Configuration, order of signals is important
	InputMapping []Signal `json:"-" gorm:"foreignkey:ConfigID"`
	// Array of file IDs used by the component configuration
	FileIDs pq.Int64Array `json:"fileIDs" gorm:"type:integer[]"`
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
	// Scaling factor for the signal raw value (defaults to 1.0)
	ScalingFactor float32 `json:"scalingFactor" gorm:"default:1"`
	// ID of Component Configuration
	ConfigID uint `json:"configID"`
}

// InfrastructureComponent data model
type InfrastructureComponent struct {
	Model
	// UUID of the IC
	UUID string `json:"uuid" gorm:"not null"`
	// Name of the IC
	Name string `json:"name" gorm:"default:''"`
	// WebsocketURL if the IC
	WebsocketURL string `json:"websocketurl" gorm:"default:''"`
	// API URL of API for IC
	APIURL string `json:"apiurl" gorm:"default:''"`
	// Category of IC (simulator, gateway, database, etc.)
	Category string `json:"category" gorm:"default:''"`
	// Type of IC (RTDS, VILLASnode, RTDS, etc.)
	Type string `json:"type" gorm:"default:''"`
	// Uptime of the IC
	Uptime float64 `json:"uptime" gorm:"default:-1"`
	// State of the IC
	State string `json:"state" gorm:"default:''"`
	// Time of last state update
	StateUpdateAt string `json:"stateUpdateAt" gorm:"default:''"`
	// Location of the IC
	Location string `json:"location" gorm:"default:''"`
	// Description of the IC
	Description string `json:"description" gorm:"default:''"`
	// JSON schema of start parameters for IC
	StartParameterSchema postgres.Jsonb `json:"startparameterschema"`
	// JSON schema of create parameters for IC
	CreateParameterSchema postgres.Jsonb `json:"createparameterschema"`
	// raw JSON of last status update
	StatusUpdateRaw postgres.Jsonb `json:"statusupdateraw"`
	// Boolean indicating if IC is managed externally (via AMQP/ VILLAScontroller)
	ManagedExternally bool `json:"managedexternally" gorm:"default:false"`
	// UUID of IC that manages this IC
	Manager string `json:"manager" gorm:"default:''"`
	// ComponentConfigurations in which the IC is used
	ComponentConfigurations []ComponentConfiguration `json:"-" gorm:"foreignkey:ICID"`
}

// Dashboard data model
type Dashboard struct {
	Model
	// Name of dashboard
	Name string `json:"name" gorm:"not null"`
	// Grid of dashboard
	Grid int `json:"grid" gorm:"default:15"`
	// Height of dashboard
	Height int `json:"height"`
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
	// IDs of signals that widget uses
	SignalIDs pq.Int64Array `json:"signalIDs" gorm:"type:integer[]"`
}

// File data model
type File struct {
	Model
	// Name of file
	Name string `json:"name" gorm:"not null"`
	// Key of file in S3 bucket
	Key string `json:"key"`
	// Type of file (MIME type)
	Type string `json:"type"`
	// Size of file (in byte)
	Size uint `json:"size"`
	// Last modification time of file
	Date string `json:"date"`
	// ID of Scenario to which file belongs
	ScenarioID uint `json:"scenarioID"`
	// File itself
	FileData []byte `json:"-" gorm:"column:FileData"`
	// Height of an image file in pixels (optional)
	ImageHeight int `json:"imageHeight" gorm:"default:0"`
	// Width of an image file in pixels (optional)
	ImageWidth int `json:"imageWidth" gorm:"default:0"`
}

// Result data model
type Result struct {
	Model
	// JSON snapshots of component configurations used to generate results
	ConfigSnapshots postgres.Jsonb `json:"configSnapshots"`
	// Description of results
	Description string `json:"description"`
	// ID of Scenario to which result belongs
	ScenarioID uint `json:"scenarioID"`
	// File IDs associated with result
	ResultFileIDs pq.Int64Array `json:"resultFileIDs" gorm:"type:integer[]"`
}
