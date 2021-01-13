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
	"math/rand"
	"strings"
	"time"

	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/zpatrick/go-config"
	"golang.org/x/crypto/bcrypt"
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

var propertiesA = json.RawMessage(`{"prop1" : "a nice prop"}`)
var propertiesB = json.RawMessage(`{"prop1" : "not so nice"}`)

var ICA = database.InfrastructureComponent{
	UUID:         "7be0322d-354e-431e-84bd-ae4c9633138b",
	WebsocketURL: "https://villas.k8s.eonerc.rwth-aachen.de/ws/ws_sig",
	APIURL:       "https://villas.k8s.eonerc.rwth-aachen.de/ws/api/v2",
	Type:         "villas-node",
	Category:     "gateway",
	Name:         "ACS Demo Signals",
	Uptime:       -1.0,
	State:        "idle",
	Location:     "k8s",
	Description:  "A signal generator for testing purposes",
	//StateUpdateAt: time.Now().Format(time.RFC1123),
	StartParameterScheme: postgres.Jsonb{propertiesA},
	ManagedExternally:    false,
}

var ICB = database.InfrastructureComponent{
	UUID:         "4854af30-325f-44a5-ad59-b67b2597de68",
	WebsocketURL: "xxx.yyy.zzz.aaa",
	Type:         "dpsim",
	Category:     "simulator",
	Name:         "Test DPsim Simulator",
	Uptime:       -1.0,
	State:        "running",
	Location:     "ACS Laboratory",
	Description:  "This is a test description",
	//StateUpdateAt: time.Now().Format(time.RFC1123),
	StartParameterScheme: postgres.Jsonb{propertiesB},
	ManagedExternally:    true,
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
	FileIDs:         []int64{},
}

var ConfigB = database.ComponentConfiguration{
	Name:            "VILLASnode gateway X",
	StartParameters: postgres.Jsonb{startParametersB},
	FileIDs:         []int64{},
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
var customPropertiesBox = json.RawMessage(`{"border_color" : "#4287f5", "border_color_opacity": 1, "background_color" : "#961520", "background_color_opacity" : 1}`)
var customPropertiesSlider = json.RawMessage(`{"default_value" : 0, "orientation" : 0, "rangeUseMinMax": false, "rangeMin" : 0, "rangeMax": 200, "rangeUseMinMax" : true, "showUnit": true, "continous_update": false, "value": "", "resizeLeftRightLock": false, "resizeTopBottomLock": true, "step": 0.1 }`)
var customPropertiesLabel = json.RawMessage(`{"textSize" : "20", "fontColor" : "#4287f5", "fontColor_opacity": 1}`)
var customPropertiesButton = json.RawMessage(`{"pressed": false, "toggle" : false, "on_value" : 1, "off_value" : 0, "background_color": "#961520", "font_color": "#4287f5", "border_color": "#4287f5", "background_color_opacity": 1}`)
var customPropertiesLamp = json.RawMessage(`{"signal" : 0, "on_color" : "#4287f5", "off_color": "#961520", "threshold" : 0.5, "on_color_opacity": 1, "off_color_opacity": 1}`)

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
	SignalIDs:        []int64{},
}

var WidgetB = database.Widget{
	Name:             "Slider",
	Type:             "Slider",
	Width:            400,
	Height:           50,
	MinHeight:        30,
	MinWidth:         380,
	X:                70,
	Y:                400,
	Z:                0,
	IsLocked:         false,
	CustomProperties: postgres.Jsonb{customPropertiesSlider},
	SignalIDs:        []int64{},
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
	SignalIDs:        []int64{},
}

var WidgetD = database.Widget{
	Name:             "Button",
	Type:             "Button",
	Width:            100,
	Height:           100,
	MinHeight:        50,
	MinWidth:         100,
	X:                10,
	Y:                50,
	Z:                0,
	IsLocked:         false,
	CustomProperties: postgres.Jsonb{customPropertiesButton},
	SignalIDs:        []int64{},
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
	SignalIDs:        []int64{},
}

func generatePassword(Len int) string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")

	var b strings.Builder
	for i := 0; i < Len; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}

	return b.String()
}

// DBAddAdminUser adds a default admin user to the DB
func DBAddAdminUser(cfg *config.Config) error {
	database.DBpool.AutoMigrate(&database.User{})

	// Check if admin user exists in DB
	var users []database.User
	err := database.DBpool.Where("Role = ?", "Admin").Find(&users).Error

	if len(users) == 0 {
		fmt.Println("No admin user found in DB, adding default admin user.")

		mode, err := cfg.String("mode")

		name, err := cfg.String("admin.user")
		if (err != nil || name == "") && mode != "test" {
			name = "admin"
		} else if mode == "test" {
			name = User0.Username
		}

		pw, err := cfg.String("admin.pass")
		if (err != nil || pw == "") && mode != "test" {
			pw = generatePassword(16)
			fmt.Printf("  Generated admin password: %s\n", pw)
		} else if mode == "test" {
			pw = StrPassword0
		}

		mail, err := cfg.String("admin.mail")
		if (err == nil || mail == "") && mode != "test" {
			mail = "admin@example.com"
		} else if mode == "test" {
			mail = User0.Mail
		}

		pwEnc, _ := bcrypt.GenerateFromPassword([]byte(pw), bcryptCost)

		// create a copy of global test data
		user := database.User{Username: name, Password: string(pwEnc),
			Role: "Admin", Mail: mail, Active: true}

		// add admin user to DB
		err = database.DBpool.Create(&user).Error
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
