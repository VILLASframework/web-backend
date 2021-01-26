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
	"io/ioutil"
	"log"
	"os"
)

// #######################################################################
// #################### Data used for testing ############################
// #######################################################################

type jsonUser struct {
	Username string
	Password string
	Mail     string
	Role     string
}

var GlobalTestData struct {
	Users      []jsonUser
	ICs        []database.InfrastructureComponent
	Scenarios  []database.Scenario
	Results    []database.Result
	Configs    []database.ComponentConfiguration
	Dashboards []database.Dashboard
	Widgets    []database.Widget
	Signals    []database.Signal
}

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
	Role: "Admin", Mail: "User_0@example.com"}
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

// Infrastructure components

var propertiesA = json.RawMessage(`{"prop1" : "a nice prop"}`)
var propertiesB = json.RawMessage(`{"prop1" : "not so nice"}`)

var ICA = database.InfrastructureComponent{
	UUID:                 "7be0322d-354e-431e-84bd-ae4c9633138b",
	WebsocketURL:         "https://villas.k8s.eonerc.rwth-aachen.de/ws/ws_sig",
	APIURL:               "https://villas.k8s.eonerc.rwth-aachen.de/ws/api/v2",
	Type:                 "villas-node",
	Category:             "gateway",
	Name:                 "ACS Demo Signals",
	Uptime:               -1.0,
	State:                "idle",
	Location:             "k8s",
	Description:          "A signal generator for testing purposes",
	StartParameterScheme: postgres.Jsonb{propertiesA},
	ManagedExternally:    false,
}

var ICB = database.InfrastructureComponent{
	UUID:                 "4854af30-325f-44a5-ad59-b67b2597de68",
	WebsocketURL:         "xxx.yyy.zzz.aaa",
	Type:                 "dpsim",
	Category:             "simulator",
	Name:                 "Test DPsim Simulator",
	Uptime:               -1.0,
	State:                "running",
	Location:             "ACS Laboratory",
	Description:          "This is a test description",
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

// Widgets
var customPropertiesSlider = json.RawMessage(`{"default_value" : 0, "orientation" : 0, "rangeUseMinMax": false, "rangeMin" : 0, "rangeMax": 200, "rangeUseMinMax" : true, "showUnit": true, "continous_update": false, "value": "", "resizeLeftRightLock": false, "resizeTopBottomLock": true, "step": 0.1 }`)
var customPropertiesLabel = json.RawMessage(`{"textSize" : "20", "fontColor" : "#4287f5", "fontColor_opacity": 1}`)

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

func ReadTestDataFromJson(path string) error {

	jsonFile, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("error opening json file: %v", err)
	}
	log.Println("Successfully opened json data file", path)

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	err = json.Unmarshal(byteValue, &GlobalTestData)
	if err != nil {
		return fmt.Errorf("error unmarshalling json: %v", err)
	}

	log.Println(len(GlobalTestData.Users))

	return nil
}

// add test users defined above
func AddTestUsers() error {

	testUsers := []database.User{User0, UserA, UserB, UserC}
	database.DBpool.AutoMigrate(&database.User{})

	for _, user := range testUsers {
		err := database.DBpool.Create(&user).Error
		if err != nil {
			return err
		}

	}

	return nil
}
