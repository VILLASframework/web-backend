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
	"bytes"
	"encoding/json"
	"fmt"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/helper"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/crypto/bcrypt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
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

type UserRequest struct {
	Username    string `json:"username,omitempty"`
	Password    string `json:"password,omitempty"`
	OldPassword string `json:"oldPassword,omitempty"`
	Mail        string `json:"mail,omitempty"`
	Role        string `json:"role,omitempty"`
	Active      string `json:"active,omitempty"`
}

var newUserA = UserRequest{
	Username: UserA.Username,
	Password: StrPasswordA,
	Role:     UserA.Role,
	Mail:     UserA.Mail,
}

var newUserB = UserRequest{
	Username: UserB.Username,
	Password: StrPasswordB,
	Role:     UserB.Role,
	Mail:     UserB.Mail,
}

var newUserC = UserRequest{
	Username: UserC.Username,
	Password: StrPasswordC,
	Role:     UserC.Role,
	Mail:     UserC.Mail,
}

// Infrastructure components

var propertiesA = json.RawMessage(`{"name" : "DPsim simulator", "category" : "Simulator", "location" : "ACSlab", "type": "DPsim"}`)
var propertiesB = json.RawMessage(`{"name" : "VILLASnode gateway", "category" : "Gateway", "location" : "ACSlab", "type": "VILLASnode"}`)

var ICA = InfrastructureComponent{
	UUID:          "4854af30-325f-44a5-ad59-b67b2597de68",
	Host:          "Host_A",
	Modeltype:     "ModelTypeA",
	Uptime:        0,
	State:         "running",
	StateUpdateAt: time.Now().String(),
	Properties:    postgres.Jsonb{propertiesA},
	RawProperties: postgres.Jsonb{propertiesA},
}

var ICB = InfrastructureComponent{
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

// Component Configuration

var ConfigA = ComponentConfiguration{
	Name:            "Example simulation",
	StartParameters: postgres.Jsonb{startParametersA},
	SelectedFileID:  3,
}

var ConfigB = ComponentConfiguration{
	Name:            "VILLASnode gateway X",
	StartParameters: postgres.Jsonb{startParametersB},
	SelectedFileID:  4,
}

// Signals

var OutSignalA = Signal{
	Name:      "outSignal_A",
	Direction: "out",
	Index:     1,
	Unit:      "V",
}

var OutSignalB = Signal{
	Name:      "outSignal_B",
	Direction: "out",
	Index:     2,
	Unit:      "V",
}

var InSignalA = Signal{
	Name:      "inSignal_A",
	Direction: "in",
	Index:     1,
	Unit:      "A",
}

var InSignalB = Signal{
	Name:      "inSignal_B",
	Direction: "in",
	Index:     2,
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
	Name: "File_A",
	Type: "text/plain",
	Size: 42,
	Date: time.Now().String(),
}

var FileB = File{
	Name: "File_B",
	Type: "text/plain",
	Size: 1234,
	Date: time.Now().String(),
}

var FileC = File{
	Name: "File_C",
	Type: "text/plain",
	Size: 32,
	Date: time.Now().String(),
}
var FileD = File{
	Name: "File_D",
	Type: "text/plain",
	Size: 5000,
	Date: time.Now().String(),
}

// Widgets
var customPropertiesA = json.RawMessage(`{"property1" : "testValue1A", "property2" : "testValue2A", "property3" : 42}`)
var customPropertiesBox = json.RawMessage(`{"border_color" : "18", "background_color" : "2", "background_color_opacity" : 0.3}`)
var customPropertiesSlider = json.RawMessage(`{"default_value" : 10, "orientation" : 0, "rangeMin" : 0, "rangeMax": 500, "step" : 1}`)
var customPropertiesLabel = json.RawMessage(`{"textSize" : "20", "fontColor" : 5}`)
var customPropertiesButton = json.RawMessage(`{"toggle" : "Value1", "on_value" : "Value2", "off_value" : Value3}`)
var customPropertiesCustomActions = json.RawMessage(`{"actions" : "Value1", "icon" : "Value2"}`)
var customPropertiesGauge = json.RawMessage(`{ "valueMin": 0, "valueMax": 1}`)
var customPropertiesLamp = json.RawMessage(`{"signal" : 0, "on_color" : 4, "off_color": 2 , "threshold" : 0.5}`)

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
	SignalIDs:        []int64{1},
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
	SignalIDs:        []int64{1},
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
	SignalIDs:        []int64{3},
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
	CustomProperties: postgres.Jsonb{customPropertiesCustomActions},
	SignalIDs:        []int64{2},
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
	SignalIDs:        []int64{4},
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
func DBAddTestData(db *gorm.DB, basePath string, router *gin.Engine) error {

	MigrateModels(db)
	// Create entries of each model (data defined in testdata.go)
	// add Admin user
	err := DBAddAdminUser(db)
	if err != nil {
		return err
	}

	// authenticate as admin
	token, err := helper.AuthenticateForTest(router, basePath+"/authenticate", "POST", helper.AdminCredentials)
	if err != nil {
		return err
	}

	// add 2 normal and 1 guest user
	code, _, err := helper.TestEndpoint(router, token, basePath+"/users", "POST", helper.KeyModels{"user": newUserA})
	if code != http.StatusOK {
		return fmt.Errorf("error adding User_A")
	}
	code, _, err = helper.TestEndpoint(router, token, basePath+"/users", "POST", helper.KeyModels{"user": newUserB})
	if code != http.StatusOK {
		return fmt.Errorf("error adding User_B")
	}
	code, _, err = helper.TestEndpoint(router, token, basePath+"/users", "POST", helper.KeyModels{"user": newUserC})
	if code != http.StatusOK {
		return fmt.Errorf("error adding User_C")
	}

	// add infrastructure components
	code, _, err = helper.TestEndpoint(router, token, basePath+"/ic", "POST", helper.KeyModels{"ic": ICA})
	if code != http.StatusOK {
		return fmt.Errorf("error adding IC A")
	}
	code, _, err = helper.TestEndpoint(router, token, basePath+"/ic", "POST", helper.KeyModels{"ic": ICB})
	if code != http.StatusOK {
		return fmt.Errorf("error adding IC B")
	}

	// add scenarios
	code, _, err = helper.TestEndpoint(router, token, basePath+"/scenarios", "POST", helper.KeyModels{"scenario": ScenarioA})
	if code != http.StatusOK {
		return fmt.Errorf("error adding Scenario A")
	}
	code, _, err = helper.TestEndpoint(router, token, basePath+"/scenarios", "POST", helper.KeyModels{"scenario": ScenarioB})
	if code != http.StatusOK {
		return fmt.Errorf("error adding Scenario B")
	}

	// add users to scenario
	code, _, err = helper.TestEndpoint(router, token, fmt.Sprintf("%v/scenarios/1/user?username=User_A", basePath), "PUT", nil)
	if code != http.StatusOK {
		return fmt.Errorf("error adding User_A to Scenario A")
	}
	code, _, err = helper.TestEndpoint(router, token, fmt.Sprintf("%v/scenarios/2/user?username=User_A", basePath), "PUT", nil)
	if code != http.StatusOK {
		return fmt.Errorf("error adding User_A to Scenario B")
	}
	code, _, err = helper.TestEndpoint(router, token, fmt.Sprintf("%v/scenarios/2/user?username=User_B", basePath), "PUT", nil)
	if code != http.StatusOK {
		return fmt.Errorf("error adding User_B to Scenario B")
	}
	code, _, err = helper.TestEndpoint(router, token, fmt.Sprintf("%v/scenarios/1/user?username=User_C", basePath), "PUT", nil)
	if code != http.StatusOK {
		return fmt.Errorf("error adding User_C to Scenario A")
	}

	// add component configurations
	configA := ConfigA
	configB := ConfigB
	configA.ScenarioID = 1
	configB.ScenarioID = 1
	configA.ICID = 1
	configB.ICID = 2
	code, _, err = helper.TestEndpoint(router, token, basePath+"/configs", "POST", helper.KeyModels{"config": configA})
	if code != http.StatusOK {
		return fmt.Errorf("error adding Config A")
	}
	code, _, err = helper.TestEndpoint(router, token, basePath+"/configs", "POST", helper.KeyModels{"config": configB})
	if code != http.StatusOK {
		return fmt.Errorf("error adding Config B")
	}

	// add dashboards
	dashboardA := DashboardA
	dashboardB := DashboardB
	dashboardA.ScenarioID = 1
	dashboardB.ScenarioID = 1
	code, _, err = helper.TestEndpoint(router, token, basePath+"/dashboards", "POST", helper.KeyModels{"dashboard": dashboardA})
	if code != http.StatusOK {
		return fmt.Errorf("error adding Dashboard B")
	}
	code, _, err = helper.TestEndpoint(router, token, basePath+"/dashboards", "POST", helper.KeyModels{"dashboard": dashboardB})
	if code != http.StatusOK {
		return fmt.Errorf("error adding Dashboard B")
	}

	// add widgets
	widgetA := WidgetA
	widgetB := WidgetB
	widgetC := WidgetC
	widgetD := WidgetD
	widgetE := WidgetE
	widgetA.DashboardID = 1
	widgetB.DashboardID = 1
	widgetC.DashboardID = 1
	widgetD.DashboardID = 1
	widgetE.DashboardID = 1
	code, _, err = helper.TestEndpoint(router, token, basePath+"/widgets", "POST", helper.KeyModels{"widget": widgetA})
	if code != http.StatusOK {
		return fmt.Errorf("error adding Widget A")
	}
	code, _, err = helper.TestEndpoint(router, token, basePath+"/widgets", "POST", helper.KeyModels{"widget": widgetB})
	if code != http.StatusOK {
		return fmt.Errorf("error adding Widget B")
	}
	code, _, err = helper.TestEndpoint(router, token, basePath+"/widgets", "POST", helper.KeyModels{"widget": widgetC})
	if code != http.StatusOK {
		return fmt.Errorf("error adding Widget C")
	}
	code, _, err = helper.TestEndpoint(router, token, basePath+"/widgets", "POST", helper.KeyModels{"widget": widgetD})
	if code != http.StatusOK {
		return fmt.Errorf("error adding Widget D")
	}
	code, _, err = helper.TestEndpoint(router, token, basePath+"/widgets", "POST", helper.KeyModels{"widget": widgetE})
	if code != http.StatusOK {
		return fmt.Errorf("error adding Widget E")
	}

	// add signals
	outSignalA := OutSignalA
	outSignalB := OutSignalB
	inSignalA := InSignalA
	inSignalB := InSignalB
	outSignalA.ConfigID = 1
	outSignalB.ConfigID = 1
	inSignalA.ConfigID = 1
	inSignalB.ConfigID = 1
	code, _, err = helper.TestEndpoint(router, token, basePath+"/signals", "POST", helper.KeyModels{"signal": outSignalB})
	if code != http.StatusOK {
		return fmt.Errorf("error adding outSignalB")
	}
	code, _, err = helper.TestEndpoint(router, token, basePath+"/signals", "POST", helper.KeyModels{"signal": outSignalA})
	if code != http.StatusOK {
		return fmt.Errorf("error adding outSignalA")
	}

	code, _, err = helper.TestEndpoint(router, token, basePath+"/signals", "POST", helper.KeyModels{"signal": inSignalA})
	if code != http.StatusOK {
		return fmt.Errorf("error adding inSignalA")
	}
	code, _, err = helper.TestEndpoint(router, token, basePath+"/signals", "POST", helper.KeyModels{"signal": inSignalB})
	if code != http.StatusOK {
		return fmt.Errorf("error adding inSignalB")
	}

	// upload files

	// upload readme file
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	fileWriter, _ := bodyWriter.CreateFormFile("file", "Readme.md")
	fh, _ := os.Open("README.md")
	defer fh.Close()

	// io copy
	_, err = io.Copy(fileWriter, fh)
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	// Create the request and add file to component config
	w1 := httptest.NewRecorder()
	req1, _ := http.NewRequest("POST", basePath+"/files?objectID=1&objectType=config", bodyBuf)
	req1.Header.Set("Content-Type", contentType)
	req1.Header.Add("Authorization", "Bearer "+token)
	router.ServeHTTP(w1, req1)

	// upload image file
	bodyBuf = &bytes.Buffer{}
	bodyWriter = multipart.NewWriter(bodyBuf)
	fileWriter, _ = bodyWriter.CreateFormFile("file", "logo.png")
	fh, _ = os.Open("doc/pictures/villas_web.png")
	defer fh.Close()

	// io copy
	_, err = io.Copy(fileWriter, fh)
	contentType = bodyWriter.FormDataContentType()
	bodyWriter.Close()

	// Create the request and add file to widget
	w2 := httptest.NewRecorder()
	req2, _ := http.NewRequest("POST", basePath+"/files?objectID=1&objectType=widget", bodyBuf)
	req2.Header.Set("Content-Type", contentType)
	req2.Header.Add("Authorization", "Bearer "+token)
	router.ServeHTTP(w2, req2)

	return nil
}
