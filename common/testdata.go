package common

import (
	"encoding/json"
	"github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// Generic

var MsgOK = ResponseMsg{
	Message: "OK.",
}

// Users
var StrPassword0 = "xyz789"
var StrPasswordA = "abc123"
var StrPasswordB = "bcd234"

// Hash passwords with bcrypt algorithm
var bcryptCost = 10
var pw0, _ = bcrypt.GenerateFromPassword([]byte(StrPassword0), bcryptCost)
var pwA, _ = bcrypt.GenerateFromPassword([]byte(StrPasswordA), bcryptCost)
var pwB, _ = bcrypt.GenerateFromPassword([]byte(StrPasswordB), bcryptCost)

var User0 = User{Username: "User_0", Password: string(pw0),
	Role: "Admin", Mail: "User_0@example.com"}
var UserA = User{Username: "User_A", Password: string(pwA),
	Role: "User", Mail: "User_A@example.com"}
var UserB = User{Username: "User_B", Password: string(pwB),
	Role: "User", Mail: "User_B@example.com"}

// Credentials

type Credentials struct {
	Username string
	Password string
}

var AdminCredentials = Credentials{
	Username: User0.Username,
	Password: StrPassword0,
}

var UserACredentials = Credentials{
	Username: UserA.Username,
	Password: StrPasswordA,
}

var UserBCredentials = Credentials{
	Username: UserB.Username,
	Password: StrPasswordB,
}

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
	Name:            "SimulationModel_A",
	StartParameters: postgres.Jsonb{startParametersA},
}

var SimulationModelB = SimulationModel{
	Name:            "SimulationModel_B",
	StartParameters: postgres.Jsonb{startParametersB},
}

// Signals

var OutSignalA = Signal{
	Name:      "outSignal_A",
	Direction: "out",
	Index:     0,
	Unit:      "V",
}

var OutSignalA_response = SignalResponse{
	Name:      OutSignalA.Name,
	Direction: OutSignalA.Direction,
	Index:     OutSignalA.Index,
	Unit:      OutSignalA.Unit,
}

var OutSignalB = Signal{
	Name:      "outSignal_B",
	Direction: "out",
	Index:     1,
	Unit:      "V",
}

var OutSignalB_response = SignalResponse{
	Name:      OutSignalB.Name,
	Direction: OutSignalB.Direction,
	Index:     OutSignalB.Index,
	Unit:      OutSignalB.Unit,
}

var InSignalA = Signal{
	Name:      "inSignal_A",
	Direction: "in",
	Index:     0,
	Unit:      "A",
}

var InSignalA_response = SignalResponse{
	Name:      InSignalA.Name,
	Direction: InSignalA.Direction,
	Index:     InSignalA.Index,
	Unit:      InSignalA.Unit,
}

var InSignalB = Signal{
	Name:      "inSignal_B",
	Direction: "in",
	Index:     1,
	Unit:      "A",
}

var InSignalB_response = SignalResponse{
	Name:      InSignalB.Name,
	Direction: InSignalB.Direction,
	Index:     InSignalB.Index,
	Unit:      InSignalB.Unit,
}

var InSignalC = Signal{
	Name:      "inSignal_C",
	Direction: "in",
	Index:     2,
	Unit:      "A",
}

var InSignalC_response = SignalResponse{
	Name:      InSignalC.Name,
	Direction: InSignalC.Direction,
	Index:     InSignalC.Index,
	Unit:      InSignalC.Unit,
}

var InSignalCUpdated = Signal{
	Name:      "inSignalupdated_C",
	Direction: InSignalC.Direction,
	Index:     InSignalC.Index,
	Unit:      "Ohm",
}

var InSignalCUpdated_response = SignalResponse{
	Name:      InSignalCUpdated.Name,
	Direction: InSignalCUpdated.Direction,
	Index:     InSignalCUpdated.Index,
	Unit:      InSignalCUpdated.Unit,
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

var FileA_response = FileResponse{
	ID:          1,
	Name:        FileA.Name,
	Type:        FileA.Type,
	Size:        FileA.Size,
	ImageWidth:  FileA.ImageWidth,
	ImageHeight: FileA.ImageHeight,
	Date:        FileA.Date,
}

var FileB = File{
	Name:        "File_B",
	Type:        "text/plain",
	Size:        1234,
	ImageHeight: 55,
	ImageWidth:  22,
	Date:        time.Now().String(),
}

var FileB_response = FileResponse{
	ID:          2,
	Name:        FileB.Name,
	Type:        FileB.Type,
	Size:        FileB.Size,
	ImageWidth:  FileB.ImageWidth,
	ImageHeight: FileB.ImageHeight,
	Date:        FileB.Date,
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
var customPropertiesB = json.RawMessage(`{"property1" : "testValue1B", "property2" : "testValue2B", "property3" : 43}`)

var WidgetA = Widget{
	Name:             "Widget_A",
	Type:             "graph",
	Width:            100,
	Height:           50,
	MinHeight:        40,
	MinWidth:         80,
	X:                10,
	Y:                10,
	Z:                10,
	IsLocked:         false,
	CustomProperties: postgres.Jsonb{customPropertiesA},
}

var WidgetB = Widget{
	Name:             "Widget_B",
	Type:             "slider",
	Width:            200,
	Height:           20,
	MinHeight:        10,
	MinWidth:         50,
	X:                100,
	Y:                -40,
	Z:                -1,
	IsLocked:         false,
	CustomProperties: postgres.Jsonb{customPropertiesB},
}
