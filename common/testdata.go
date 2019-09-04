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

var AdminCredentials = Request{
	Username: User0.Username,
	Password: StrPassword0,
}

var UserACredentials = Request{
	Username: UserA.Username,
	Password: StrPasswordA,
}

var UserBCredentials = Request{
	Username: UserB.Username,
	Password: StrPasswordB,
}

// Simulators

var propertiesA = json.RawMessage(`{"name" : "TestNameA", "category" : "CategoryA", "location" : "anywhere on earth", "type": "dummy"}`)
var propertiesB = json.RawMessage(`{"name" : "TestNameB", "category" : "CategoryB", "location" : "where ever you want", "type": "generic"}`)
var propertiesC = json.RawMessage(`{"name" : "TestNameC", "category" : "CategoryC", "location" : "my desk", "type": "blubb"}`)
var propertiesCupdated = json.RawMessage(`{"name" : "TestNameCUpdate", "category" : "CategoryC", "location" : "my desk", "type": "blubb"}`)

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

var SimulatorA_response = SimulatorResponse{
	ID:            1,
	UUID:          SimulatorA.UUID,
	Host:          SimulatorA.Host,
	Modeltype:     SimulatorA.Modeltype,
	Uptime:        SimulatorA.Uptime,
	State:         SimulatorA.State,
	StateUpdateAt: SimulatorA.StateUpdateAt,
	Properties:    SimulatorA.Properties,
	RawProperties: SimulatorA.RawProperties,
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

var SimulatorB_response = SimulatorResponse{
	ID:            2,
	UUID:          SimulatorB.UUID,
	Host:          SimulatorB.Host,
	Modeltype:     SimulatorB.Modeltype,
	Uptime:        SimulatorB.Uptime,
	State:         SimulatorB.State,
	StateUpdateAt: SimulatorB.StateUpdateAt,
	Properties:    SimulatorB.Properties,
	RawProperties: SimulatorB.RawProperties,
}

var SimulatorC = Simulator{
	UUID:          "6d9776bf-b693-45e8-97b6-4c13d151043f",
	Host:          "Host_C",
	Modeltype:     "ModelTypeC",
	Uptime:        0,
	State:         "idle",
	StateUpdateAt: time.Now().String(),
	Properties:    postgres.Jsonb{propertiesC},
	RawProperties: postgres.Jsonb{propertiesC},
}

var SimulatorC_response = SimulatorResponse{
	ID:            3,
	UUID:          SimulatorC.UUID,
	Host:          SimulatorC.Host,
	Modeltype:     SimulatorC.Modeltype,
	Uptime:        SimulatorC.Uptime,
	State:         SimulatorC.State,
	StateUpdateAt: SimulatorC.StateUpdateAt,
	Properties:    SimulatorC.Properties,
	RawProperties: SimulatorC.RawProperties,
}

var SimulatorCUpdated = Simulator{
	UUID:          SimulatorC.UUID,
	Host:          "Host_Cupdated",
	Modeltype:     "ModelTypeCUpdated",
	Uptime:        SimulatorC.Uptime,
	State:         "running",
	StateUpdateAt: time.Now().String(),
	Properties:    postgres.Jsonb{propertiesCupdated},
	RawProperties: postgres.Jsonb{propertiesCupdated},
}

var SimulatorCUpdated_response = SimulatorResponse{
	ID:            3,
	UUID:          SimulatorCUpdated.UUID,
	Host:          SimulatorCUpdated.Host,
	Modeltype:     SimulatorCUpdated.Modeltype,
	Uptime:        SimulatorCUpdated.Uptime,
	State:         SimulatorCUpdated.State,
	StateUpdateAt: SimulatorCUpdated.StateUpdateAt,
	Properties:    SimulatorCUpdated.Properties,
	RawProperties: SimulatorCUpdated.RawProperties,
}

// Scenarios

var startParametersA = json.RawMessage(`{"parameter1" : "testValue1A", "parameter2" : "testValue2A", "parameter3" : 42}`)
var startParametersB = json.RawMessage(`{"parameter1" : "testValue1B", "parameter2" : "testValue2B", "parameter3" : 43}`)
var startParametersC = json.RawMessage(`{"parameter1" : "testValue1C", "parameter2" : "testValue2C", "parameter3" : 44}`)

var ScenarioA = Scenario{Name: "Scenario_A", Running: true, StartParameters: postgres.Jsonb{startParametersA}}
var ScenarioA_response = ScenarioResponse{ID: 1, Name: ScenarioA.Name, Running: ScenarioA.Running, StartParameters: ScenarioA.StartParameters}
var ScenarioB = Scenario{Name: "Scenario_B", Running: false, StartParameters: postgres.Jsonb{startParametersB}}
var ScenarioB_response = ScenarioResponse{ID: 2, Name: ScenarioB.Name, Running: ScenarioB.Running, StartParameters: ScenarioB.StartParameters}
var ScenarioC = Scenario{Name: "Scenario_C", Running: false, StartParameters: postgres.Jsonb{startParametersC}}
var ScenarioC_response = ScenarioResponse{ID: 3, Name: ScenarioC.Name, Running: ScenarioC.Running, StartParameters: ScenarioC.StartParameters}
var ScenarioCUpdated = Scenario{Name: "Scenario_Cupdated", Running: true, StartParameters: postgres.Jsonb{startParametersC}}
var ScenarioCUpdated_response = ScenarioResponse{ID: 3, Name: ScenarioCUpdated.Name, Running: ScenarioCUpdated.Running, StartParameters: ScenarioCUpdated.StartParameters}

// Simulation Models

var SimulationModelA = SimulationModel{
	Name:            "SimulationModel_A",
	OutputLength:    1,
	InputLength:     1,
	StartParameters: postgres.Jsonb{startParametersA},
}

var SimulationModelA_response = SimulationModelResponse{
	ID:              1,
	Name:            SimulationModelA.Name,
	InputLength:     SimulationModelA.InputLength,
	OutputLength:    SimulationModelA.OutputLength,
	StartParameters: SimulationModelA.StartParameters,
}

var SimulationModelB = SimulationModel{
	Name:            "SimulationModel_B",
	OutputLength:    1,
	InputLength:     1,
	StartParameters: postgres.Jsonb{startParametersB},
}

var SimulationModelB_response = SimulationModelResponse{
	ID:              2,
	Name:            SimulationModelB.Name,
	InputLength:     SimulationModelB.InputLength,
	OutputLength:    SimulationModelB.OutputLength,
	StartParameters: SimulationModelB.StartParameters,
}

var SimulationModelC = SimulationModel{
	Name:            "SimulationModel_C",
	OutputLength:    1,
	InputLength:     1,
	StartParameters: postgres.Jsonb{startParametersC},
}

var SimulationModelC_response = SimulationModelResponse{
	ID:              3,
	Name:            SimulationModelC.Name,
	InputLength:     SimulationModelC.InputLength,
	OutputLength:    SimulationModelC.OutputLength,
	ScenarioID:      SimulationModelC.ScenarioID,
	SimulatorID:     SimulationModelC.SimulatorID,
	StartParameters: SimulationModelC.StartParameters,
}

var SimulationModelCUpdated = SimulationModel{
	Name:            "SimulationModel_CUpdated",
	OutputLength:    SimulationModelC.OutputLength,
	InputLength:     SimulationModelC.InputLength,
	ScenarioID:      SimulationModelC.ScenarioID,
	SimulatorID:     2,
	StartParameters: SimulationModelC.StartParameters,
	InputMapping:    SimulationModelC.InputMapping,
	OutputMapping:   SimulationModelC.OutputMapping,
}

var SimulationModelCUpdated_response = SimulationModelResponse{
	ID:              3,
	Name:            SimulationModelCUpdated.Name,
	InputLength:     SimulationModelCUpdated.InputLength,
	OutputLength:    SimulationModelCUpdated.OutputLength,
	ScenarioID:      SimulationModelCUpdated.ScenarioID,
	SimulatorID:     SimulationModelCUpdated.SimulatorID,
	StartParameters: SimulationModelCUpdated.StartParameters,
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

var DashboardA = Dashboard{Name: "Dashboard_A", Grid: 15}
var DashboardA_response = DashboardResponse{ID: 1, Name: DashboardA.Name, Grid: DashboardA.Grid, ScenarioID: DashboardA.ScenarioID}
var DashboardB = Dashboard{Name: "Dashboard_B", Grid: 10}
var DashboardB_response = DashboardResponse{ID: 2, Name: DashboardB.Name, Grid: DashboardB.Grid, ScenarioID: DashboardB.ScenarioID}
var DashboardC = Dashboard{Name: "Dashboard_C", Grid: 25}
var DashboardC_response = DashboardResponse{ID: 3, Name: DashboardC.Name, Grid: DashboardC.Grid, ScenarioID: DashboardC.ScenarioID}
var DashboardCUpdated = Dashboard{Name: "Dashboard_Cupdated", Grid: 24}
var DashboardCUpdated_response = DashboardResponse{ID: 3, Name: DashboardCUpdated.Name, Grid: DashboardCUpdated.Grid, ScenarioID: DashboardCUpdated.ScenarioID}

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
var customPropertiesC = json.RawMessage(`{"property1" : "testValue1C", "property2" : "testValue2C", "property3" : 44}`)

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

var WidgetA_response = WidgetResponse{
	ID:               1,
	Name:             WidgetA.Name,
	Type:             WidgetA.Type,
	Width:            WidgetA.Width,
	Height:           WidgetA.Height,
	MinWidth:         WidgetA.MinWidth,
	MinHeight:        WidgetA.MinHeight,
	X:                WidgetA.X,
	Y:                WidgetA.Y,
	Z:                WidgetA.Z,
	IsLocked:         WidgetA.IsLocked,
	CustomProperties: WidgetA.CustomProperties,
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
	Z:                0,
	IsLocked:         false,
	CustomProperties: postgres.Jsonb{customPropertiesB},
}

var WidgetB_response = WidgetResponse{
	ID:               2,
	Name:             WidgetB.Name,
	Type:             WidgetB.Type,
	Width:            WidgetB.Width,
	Height:           WidgetB.Height,
	MinWidth:         WidgetB.MinWidth,
	MinHeight:        WidgetB.MinHeight,
	X:                WidgetB.X,
	Y:                WidgetB.Y,
	Z:                WidgetB.Z,
	IsLocked:         WidgetB.IsLocked,
	CustomProperties: WidgetB.CustomProperties,
}

var WidgetC = Widget{
	Name:             "Widget_C",
	Type:             "bargraph",
	Height:           30,
	Width:            100,
	MinHeight:        20,
	MinWidth:         50,
	X:                11,
	Y:                12,
	Z:                13,
	IsLocked:         false,
	CustomProperties: postgres.Jsonb{customPropertiesC},
}

var WidgetC_response = WidgetResponse{
	ID:               3,
	Name:             WidgetC.Name,
	Type:             WidgetC.Type,
	Width:            WidgetC.Width,
	Height:           WidgetC.Height,
	MinWidth:         WidgetC.MinWidth,
	MinHeight:        WidgetC.MinHeight,
	X:                WidgetC.X,
	Y:                WidgetC.Y,
	Z:                WidgetC.Z,
	IsLocked:         WidgetC.IsLocked,
	CustomProperties: WidgetC.CustomProperties,
}

var WidgetCUpdated_response = WidgetResponse{
	ID:               3,
	Name:             "Widget_CUpdated",
	Type:             WidgetC.Type,
	Height:           35,
	Width:            110,
	MinHeight:        WidgetC.MinHeight,
	MinWidth:         WidgetC.MinWidth,
	X:                WidgetC.X,
	Y:                WidgetC.Y,
	Z:                WidgetC.Z,
	IsLocked:         WidgetC.IsLocked,
	CustomProperties: WidgetC.CustomProperties,
}
