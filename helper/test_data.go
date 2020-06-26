/** Helper package, test data.
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
package helper

import (
	"encoding/json"
	"fmt"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
	"github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// #######################################################################
// #################### Data used for testing ############################
// #######################################################################

// Credentials
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

var User0 = database.User{Username: "User_0", Password: string(pw0),
	Role: "Admin", Mail: "User_0@example.com", Active: true}
var UserA = database.User{Username: "User_A", Password: string(pwA),
	Role: "User", Mail: "User_A@example.com", Active: true}
var UserB = database.User{Username: "User_B", Password: string(pwB),
	Role: "User", Mail: "User_B@example.com", Active: true}
var UserC = database.User{Username: "User_C", Password: string(pwC),
	Role: "Guest", Mail: "User_C@example.com", Active: true}

type UserRequest struct {
	Username    string `json:"username,omitempty"`
	Password    string `json:"password,omitempty"`
	OldPassword string `json:"oldPassword,omitempty"`
	Mail        string `json:"mail,omitempty"`
	Role        string `json:"role,omitempty"`
	Active      string `json:"active,omitempty"`
}

var NewUserA = UserRequest{
	Username: UserA.Username,
	Password: StrPasswordA,
	Role:     UserA.Role,
	Mail:     UserA.Mail,
}

var NewUserB = UserRequest{
	Username: UserB.Username,
	Password: StrPasswordB,
	Role:     UserB.Role,
	Mail:     UserB.Mail,
}

var NewUserC = UserRequest{
	Username: UserC.Username,
	Password: StrPasswordC,
	Role:     UserC.Role,
	Mail:     UserC.Mail,
}

// Infrastructure components

var propertiesA = json.RawMessage(`{"location" : "ACSlab"}`)
var propertiesB = json.RawMessage(`{"location" : "ACSlab"}`)

var ICA = database.InfrastructureComponent{
	UUID:          "4854af30-325f-44a5-ad59-b67b2597de68",
	Host:          "xxx.yyy.zzz.aaa",
	Type:          "DPsim",
	Category:      "Simulator",
	Name:          "Test DPsim Simulator",
	Uptime:        0,
	State:         "running",
	StateUpdateAt: time.Now().String(),
	Properties:    postgres.Jsonb{propertiesA},
	RawProperties: postgres.Jsonb{propertiesA},
}

var ICB = database.InfrastructureComponent{
	UUID:          "7be0322d-354e-431e-84bd-ae4c9633138b",
	Host:          "https://villas.k8s.eonerc.rwth-aachen.de/ws/ws_sig",
	APIHost:       "https://villas.k8s.eonerc.rwth-aachen.de/ws/api/v1",
	Type:          "VILLASnode Signal Generator",
	Category:      "Signal Generator",
	Name:          "ACS Demo Signals",
	Uptime:        0,
	State:         "idle",
	StateUpdateAt: time.Now().String(),
	Properties:    postgres.Jsonb{propertiesB},
	RawProperties: postgres.Jsonb{propertiesB},
}

// Scenarios

var startParametersA = json.RawMessage(`{"parameter1" : "testValue1A", "parameter2" : "testValue2A", "parameter3" : 42}`)
var startParametersB = json.RawMessage(`{"parameter1" : "testValue1B", "parameter2" : "testValue2B", "parameter3" : 43}`)

var ScenarioA = database.Scenario{
	Name:            "Scenario_A",
	Running:         true,
	StartParameters: postgres.Jsonb{startParametersA},
}
var ScenarioB = database.Scenario{
	Name:            "Scenario_B",
	Running:         false,
	StartParameters: postgres.Jsonb{startParametersB},
}

// Component Configuration

var ConfigA = database.ComponentConfiguration{
	Name:            "Example for Signal generator",
	StartParameters: postgres.Jsonb{startParametersA},
	SelectedFileID:  -1,
}

var ConfigB = database.ComponentConfiguration{
	Name:            "VILLASnode gateway X",
	StartParameters: postgres.Jsonb{startParametersB},
	SelectedFileID:  -1,
}

// Signals

var OutSignalA = database.Signal{
	Name:      "outSignal_A",
	Direction: "out",
	Index:     1,
	Unit:      "V",
}

var OutSignalB = database.Signal{
	Name:      "outSignal_B",
	Direction: "out",
	Index:     2,
	Unit:      "V",
}

var OutSignalC = database.Signal{
	Name:      "outSignal_C",
	Direction: "out",
	Index:     3,
	Unit:      "---",
}

var OutSignalD = database.Signal{
	Name:      "outSignal_D",
	Direction: "out",
	Index:     4,
	Unit:      "---",
}

var OutSignalE = database.Signal{
	Name:      "outSignal_E",
	Direction: "out",
	Index:     5,
	Unit:      "---",
}

var InSignalA = database.Signal{
	Name:      "inSignal_A",
	Direction: "in",
	Index:     1,
	Unit:      "---",
}

var InSignalB = database.Signal{
	Name:      "inSignal_B",
	Direction: "in",
	Index:     2,
	Unit:      "---",
}

// Dashboards

var DashboardA = database.Dashboard{
	Name: "Dashboard_A",
	Grid: 15,
}
var DashboardB = database.Dashboard{
	Name: "Dashboard_B",
	Grid: 10,
}

// Files

var FileA = database.File{
	Name: "File_A",
	Type: "text/plain",
	Size: 42,
	Date: time.Now().String(),
}

var FileB = database.File{
	Name: "File_B",
	Type: "text/plain",
	Size: 1234,
	Date: time.Now().String(),
}

var FileC = database.File{
	Name: "File_C",
	Type: "text/plain",
	Size: 32,
	Date: time.Now().String(),
}
var FileD = database.File{
	Name: "File_D",
	Type: "text/plain",
	Size: 5000,
	Date: time.Now().String(),
}

// Widgets
var customPropertiesA = json.RawMessage(`{"property1" : "testValue1A", "property2" : "testValue2A", "property3" : 42}`)
var customPropertiesBox = json.RawMessage(`{"border_color" : "18", "background_color" : "2", "background_color_opacity" : 0.3}`)
var customPropertiesSlider = json.RawMessage(`{"default_value" : 10, "orientation" : 0, "rangeUseMinMax": false, "rangeMin" : 0, "rangeMax": 500, "step" : 1}`)
var customPropertiesLabel = json.RawMessage(`{"textSize" : "20", "fontColor" : 5}`)
var customPropertiesButton = json.RawMessage(`{"toggle" : "Value1", "on_value" : "Value2", "off_value" : Value3}`)
var customPropertiesCustomActions = json.RawMessage(`{"actions" : "Value1", "icon" : "Value2"}`)
var customPropertiesGauge = json.RawMessage(`{ "valueMin": 0, "valueMax": 1}`)
var customPropertiesLamp = json.RawMessage(`{"signal" : 0, "on_color" : 4, "off_color": 2 , "threshold" : 0.5}`)

var WidgetA = database.Widget{
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
	SignalIDs:        []int64{1},
}

var WidgetB = database.Widget{
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
	SignalIDs:        []int64{1},
}

var WidgetC = database.Widget{
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
	SignalIDs:        []int64{3},
}

var WidgetD = database.Widget{
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
	CustomProperties: postgres.Jsonb{customPropertiesCustomActions},
	SignalIDs:        []int64{2},
}

var WidgetE = database.Widget{
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
	SignalIDs:        []int64{4},
}

// add a default admin user to the DB
func DBAddAdminUser() error {
	database.DBpool.AutoMigrate(&database.User{})

	// Check if admin user exists in DB
	var users []database.User
	err := database.DBpool.Where("Role = ?", "Admin").Find(&users).Error

	if len(users) == 0 {
		fmt.Println("No admin user found in DB, adding default admin user.")
		//create a copy of global test data
		user0 := User0
		// add admin user to DB
		err = database.DBpool.Create(&user0).Error
	}

	return err
}

// add default admin user, normal users and a guest to the DB
func DBAddAdminAndUserAndGuest() error {
	database.DBpool.AutoMigrate(&database.User{})

	//create a copy of global test data
	user0 := User0
	userA := UserA
	userB := UserB
	userC := UserC

	// add admin user to DB
	err := database.DBpool.Create(&user0).Error
	// add normal users to DB
	err = database.DBpool.Create(&userA).Error
	err = database.DBpool.Create(&userB).Error
	// add guest user to DB
	err = database.DBpool.Create(&userC).Error
	return err
}
