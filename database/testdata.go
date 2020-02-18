/** Database package, test data.
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
	"encoding/json"
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// #######################################################################
// #################### Data used for testing ############################
// #######################################################################

// Users
var StrPassword0 = "xyz789"
var StrPasswordA = "abc123"
var StrPasswordB = "bcd234"
var StrPasswordC = "guestpw"

// Hash passwords with bcrypt algorithm
var bcryptCost = 10
var pw0, _ = bcrypt.GenerateFromPassword([]byte(StrPassword0), bcryptCost)
var pwA, _ = bcrypt.GenerateFromPassword([]byte(StrPasswordA), bcryptCost)
var pwB, _ = bcrypt.GenerateFromPassword([]byte(StrPasswordB), bcryptCost)
var pwC, _ = bcrypt.GenerateFromPassword([]byte(StrPasswordC), bcryptCost)

var User0 = User{Username: "User_0", Password: string(pw0),
	Role: "Admin", Mail: "User_0@example.com", Active: true}
var UserA = User{Username: "User_A", Password: string(pwA),
	Role: "User", Mail: "User_A@example.com", Active: true}
var UserB = User{Username: "User_B", Password: string(pwB),
	Role: "User", Mail: "User_B@example.com", Active: true}
var UserC = User{Username: "User_C", Password: string(pwC),
	Role: "Guest", Mail: "User_C@example.com", Active: true}

// Simulators

var propertiesA = json.RawMessage(`{"name" : "TestNameA", "category" : "CategoryA", "location" : "anywhere on earth", "type": "dummy"}`)
var propertiesB = json.RawMessage(`{"name" : "TestNameB", "category" : "CategoryB", "location" : "where ever you want", "type": "generic"}`)

var SimulatorA = Simulator{
	UUID:          "4854af30-325f-44a5-ad59-b67b2597de68",
	Host:          "Host_A",
	Modeltype:     "ModelTypeA",
	Uptime:        0,
	State:         "running",
	StateUpdateAt: time.Now().String(),
	Properties:    postgres.Jsonb{propertiesA},
	RawProperties: postgres.Jsonb{propertiesA},
}

var SimulatorB = Simulator{
	UUID:          "7be0322d-354e-431e-84bd-ae4c9633138b",
	Host:          "Host_B",
	Modeltype:     "ModelTypeB",
	Uptime:        0,
	State:         "idle",
	StateUpdateAt: time.Now().String(),
	Properties:    postgres.Jsonb{propertiesB},
	RawProperties: postgres.Jsonb{propertiesB},
}

// Scenarios

var startParametersA = json.RawMessage(`{"parameter1" : "testValue1A", "parameter2" : "testValue2A", "parameter3" : 42}`)
var startParametersB = json.RawMessage(`{"parameter1" : "testValue1B", "parameter2" : "testValue2B", "parameter3" : 43}`)

var ScenarioA = Scenario{
	Name:            "Scenario_A",
	Running:         true,
	StartParameters: postgres.Jsonb{startParametersA},
}
var ScenarioB = Scenario{
	Name:            "Scenario_B",
	Running:         false,
	StartParameters: postgres.Jsonb{startParametersB},
}

// Simulation Models

var SimulationModelA = SimulationModel{
	Name:                "SimulationModel_A",
	StartParameters:     postgres.Jsonb{startParametersA},
	SelectedModelFileID: 1,
}

var SimulationModelB = SimulationModel{
	Name:                "SimulationModel_B",
	StartParameters:     postgres.Jsonb{startParametersB},
	SelectedModelFileID: 2,
}

// Signals

var OutSignalA = Signal{
	Name:      "outSignal_A",
	Direction: "out",
	Index:     0,
	Unit:      "V",
}

var OutSignalB = Signal{
	Name:      "outSignal_B",
	Direction: "out",
	Index:     1,
	Unit:      "V",
}

var InSignalA = Signal{
	Name:      "inSignal_A",
	Direction: "in",
	Index:     0,
	Unit:      "A",
}

var InSignalB = Signal{
	Name:      "inSignal_B",
	Direction: "in",
	Index:     1,
	Unit:      "A",
}

// Dashboards

var DashboardA = Dashboard{
	Name: "Dashboard_A",
	Grid: 15,
}
var DashboardB = Dashboard{
	Name: "Dashboard_B",
	Grid: 10,
}

// Files

var FileA = File{
	Name:        "File_A",
	Type:        "text/plain",
	Size:        42,
	ImageHeight: 333,
	ImageWidth:  111,
	Date:        time.Now().String(),
}

var FileB = File{
	Name:        "File_B",
	Type:        "text/plain",
	Size:        1234,
	ImageHeight: 55,
	ImageWidth:  22,
	Date:        time.Now().String(),
}

var FileC = File{
	Name:        "File_C",
	Type:        "text/plain",
	Size:        32,
	ImageHeight: 10,
	ImageWidth:  10,
	Date:        time.Now().String(),
}
var FileD = File{
	Name:        "File_D",
	Type:        "text/plain",
	Size:        5000,
	ImageHeight: 400,
	ImageWidth:  800,
	Date:        time.Now().String(),
}

// Widgets
var customPropertiesA = json.RawMessage(`{"property1" : "testValue1A", "property2" : "testValue2A", "property3" : 42}`)
var customPropertiesBox = json.RawMessage(`{"border_color" : "18", "background_color" : "2", "background_color_opacity" : 0.3}`)
var customPropertiesSlider = json.RawMessage(`{"default_value" : 10, "orientation" : 0, "rangeMin" : 0, "rangeMax": 500, "step" : 1}`)
var customPropertiesLabel = json.RawMessage(`{"textSize" : "20", "fontColor" : 5}`)
var customPropertiesButton = json.RawMessage(`{"toggle" : "Value1", "on_value" : "Value2", "off_value" : Value3}`)
var customPropertiesCustomActions = json.RawMessage(`{"actions" : "Value1", "icon" : "Value2"}`)
var customPropertiesGauge = json.RawMessage(`{ "valueMin": 0, "valueMax": 1}`)
var customPropertiesLamp = json.RawMessage(`{"simulationModel" : "null", "signal" : 0, "on_color" : 4, "off_color": 2 , "threshold" : 0.5}`)

var WidgetA = Widget{
	Name:             "Label",
	Type:             "Label",
	Width:            100,
	Height:           50,
	MinHeight:        40,
	MinWidth:         80,
	X:                10,
	Y:                10,
	Z:                200,
	IsLocked:         false,
	CustomProperties: postgres.Jsonb{customPropertiesLabel},
}

var WidgetB = Widget{
	Name:             "Slider",
	Type:             "Slider",
	Width:            200,
	Height:           20,
	MinHeight:        10,
	MinWidth:         50,
	X:                70,
	Y:                10,
	Z:                0,
	IsLocked:         false,
	CustomProperties: postgres.Jsonb{customPropertiesSlider},
}

var WidgetC = Widget{
	Name:             "Box",
	Type:             "Box",
	Width:            200,
	Height:           200,
	MinHeight:        10,
	MinWidth:         50,
	X:                300,
	Y:                10,
	Z:                0,
	IsLocked:         false,
	CustomProperties: postgres.Jsonb{customPropertiesBox},
}

var WidgetD = Widget{
	Name:             "Action",
	Type:             "Action",
	Width:            20,
	Height:           20,
	MinHeight:        10,
	MinWidth:         50,
	X:                10,
	Y:                50,
	Z:                0,
	IsLocked:         false,
	CustomProperties: postgres.Jsonb{customPropertiesBox},
}

var WidgetE = Widget{
	Name:             "Lamp",
	Type:             "Lamp",
	Width:            200,
	Height:           20,
	MinHeight:        10,
	MinWidth:         50,
	X:                50,
	Y:                300,
	Z:                0,
	IsLocked:         false,
	CustomProperties: postgres.Jsonb{customPropertiesLamp},
}

func DBAddAdminUser(db *gorm.DB) error {
	db.AutoMigrate(&User{})

	//create a copy of global test data
	user0 := User0
	// add admin user to DB
	err := db.Create(&user0).Error
	return err
}

func DBAddAdminAndUserAndGuest(db *gorm.DB) error {
	db.AutoMigrate(&User{})

	//create a copy of global test data
	user0 := User0
	userA := UserA
	userB := UserB
	userC := UserC

	// add admin user to DB
	err := db.Create(&user0).Error
	// add normal users to DB
	err = db.Create(&userA).Error
	err = db.Create(&userB).Error
	// add guest user to DB
	err = db.Create(&userC).Error
	return err
}

// Populates DB with test data
func DBAddTestData(db *gorm.DB) error {

	MigrateModels(db)
	// Create entries of each model (data defined in testdata.go)

	//create a copy of global test data
	user0 := User0
	userA := UserA
	userB := UserB
	userC := UserC

	simulatorA := SimulatorA
	simulatorB := SimulatorB

	scenarioA := ScenarioA
	scenarioB := ScenarioB

	outSignalA := OutSignalA
	outSignalB := OutSignalB
	inSignalA := InSignalA
	inSignalB := InSignalB

	modelA := SimulationModelA
	modelB := SimulationModelB

	dashboardA := DashboardA
	dashboardB := DashboardB

	fileA := FileA
	fileB := FileB
	fileC := FileC
	fileD := FileD

	widgetA := WidgetA
	widgetB := WidgetB
	widgetC := WidgetC
	widgetD := WidgetD
	widgetE := WidgetE

	// Users
	err := db.Create(&user0).Error

	// add normal users to DB
	err = db.Create(&userA).Error
	err = db.Create(&userB).Error

	// add Guest user to DB
	err = db.Create(&userC).Error

	// Simulators
	err = db.Create(&simulatorA).Error
	err = db.Create(&simulatorB).Error

	// Scenarios
	err = db.Create(&scenarioA).Error
	err = db.Create(&scenarioB).Error

	// Signals
	err = db.Create(&inSignalA).Error
	err = db.Create(&inSignalB).Error
	err = db.Create(&outSignalA).Error
	err = db.Create(&outSignalB).Error

	// Simulation Models
	err = db.Create(&modelA).Error
	err = db.Create(&modelB).Error

	// Dashboards
	err = db.Create(&dashboardA).Error
	err = db.Create(&dashboardB).Error

	// Files
	err = db.Create(&fileA).Error
	err = db.Create(&fileB).Error
	err = db.Create(&fileC).Error
	err = db.Create(&fileD).Error

	// Widgets
	err = db.Create(&widgetA).Error
	err = db.Create(&widgetB).Error
	err = db.Create(&widgetC).Error
	err = db.Create(&widgetD).Error
	err = db.Create(&widgetE).Error

	// Associations between models
	// For `belongs to` use the model with id=1
	// For `has many` use the models with id=1 and id=2

	// User HM Scenarios, Scenario HM Users (Many-to-Many)
	err = db.Model(&scenarioA).Association("Users").Append(&userA).Error
	err = db.Model(&scenarioB).Association("Users").Append(&userA).Error
	err = db.Model(&scenarioB).Association("Users").Append(&userB).Error
	err = db.Model(&scenarioA).Association("Users").Append(&userC).Error
	err = db.Model(&scenarioA).Association("Users").Append(&user0).Error

	// Scenario HM SimulationModels
	err = db.Model(&scenarioA).Association("SimulationModels").Append(&modelA).Error
	err = db.Model(&scenarioA).Association("SimulationModels").Append(&modelB).Error

	// Scenario HM Dashboards
	err = db.Model(&scenarioA).Association("Dashboards").Append(&dashboardA).Error
	err = db.Model(&scenarioA).Association("Dashboards").Append(&dashboardB).Error

	// Dashboard HM Widget
	err = db.Model(&dashboardA).Association("Widgets").Append(&widgetA).Error
	err = db.Model(&dashboardA).Association("Widgets").Append(&widgetB).Error
	err = db.Model(&dashboardA).Association("Widgets").Append(&widgetC).Error
	err = db.Model(&dashboardA).Association("Widgets").Append(&widgetD).Error
	err = db.Model(&dashboardA).Association("Widgets").Append(&widgetE).Error

	// SimulationModel HM Signals
	err = db.Model(&modelA).Association("InputMapping").Append(&inSignalA).Error
	err = db.Model(&modelA).Association("InputMapping").Append(&inSignalB).Error
	err = db.Model(&modelA).Association("InputMapping").Append(&outSignalA).Error
	err = db.Model(&modelA).Association("InputMapping").Append(&outSignalB).Error

	// SimulationModel HM Files
	err = db.Model(&modelA).Association("Files").Append(&fileC).Error
	err = db.Model(&modelA).Association("Files").Append(&fileD).Error

	// Simulator HM SimulationModels
	err = db.Model(&simulatorA).Association("SimulationModels").Append(&modelA).Error
	err = db.Model(&simulatorA).Association("SimulationModels").Append(&modelB).Error

	// Widget HM Files
	err = db.Model(&widgetA).Association("Files").Append(&fileA).Error
	err = db.Model(&widgetA).Association("Files").Append(&fileB).Error

	return err
}
