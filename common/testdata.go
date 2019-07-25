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

// Hash passwords with bcrypt algorithm
var bcryptCost = 10
var pw0, _ = bcrypt.GenerateFromPassword([]byte("xyz789"), bcryptCost)
var pwA, _ = bcrypt.GenerateFromPassword([]byte("abc123"), bcryptCost)
var pwB, _ = bcrypt.GenerateFromPassword([]byte("bcd234"), bcryptCost)
var User0 = User{ID: 1, Username: "User_0", Password: string(pw0), Role: "Admin", Mail: "User_0@example.com"}
var User0_response = UserResponse{Username: User0.Username, Role: User0.Role, ID: User0.ID, Mail: User0.Mail}
var UserA = User{ID: 2, Username: "User_A", Password: string(pwA), Role: "User", Mail: "User_A@example.com"}
var UserA_response = UserResponse{Username: UserA.Username, Role: UserA.Role, ID: UserA.ID, Mail: UserA.Mail}
var UserB = User{ID: 3, Username: "User_B", Password: string(pwB), Role: "User", Mail: "User_B@example.com"}
var UserB_response = UserResponse{Username: UserB.Username, Role: UserB.Role, ID: UserB.ID, Mail: UserB.Mail}

// Credentials

var CredAdmin = credentials{
	Username: User0.Username,
	Password: "xyz789",
}

var CredUser = credentials{
	Username: UserA.Username,
	Password: "abc123",
}

// Simulators

var propertiesA = json.RawMessage(`{"name" : "TestNameA", "category" : "CategoryA", "location" : "anywhere on earth", "type": "dummy"}`)
var propertiesB = json.RawMessage(`{"name" : "TestNameB", "category" : "CategoryB", "location" : "where ever you want", "type": "generic"}`)
var propertiesC = json.RawMessage(`{"name" : "TestNameC", "category" : "CategoryC", "location" : "my desk", "type": "blubb"}`)
var propertiesCupdated = json.RawMessage(`{"name" : "TestNameCUpdate", "category" : "CategoryC", "location" : "my desk", "type": "blubb"}`)

var SimulatorA = Simulator{
	ID:            1,
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
	ID:            SimulatorA.ID,
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
	ID:            2,
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
	ID:            SimulatorB.ID,
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
	ID:            3,
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
	ID:            SimulatorC.ID,
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
	ID:            SimulatorC.ID,
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
	ID:            SimulatorCUpdated.ID,
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

var ScenarioA = Scenario{ID: 1, Name: "Scenario_A", Running: true, StartParameters: postgres.Jsonb{startParametersA}}
var ScenarioA_response = ScenarioResponse{ID: ScenarioA.ID, Name: ScenarioA.Name, Running: ScenarioA.Running, StartParameters: ScenarioA.StartParameters}
var ScenarioB = Scenario{ID: 2, Name: "Scenario_B", Running: false, StartParameters: postgres.Jsonb{startParametersB}}
var ScenarioB_response = ScenarioResponse{ID: ScenarioB.ID, Name: ScenarioB.Name, Running: ScenarioB.Running, StartParameters: ScenarioB.StartParameters}
var ScenarioC = Scenario{ID: 3, Name: "Scenario_C", Running: false, StartParameters: postgres.Jsonb{startParametersC}}
var ScenarioC_response = ScenarioResponse{ID: ScenarioC.ID, Name: ScenarioC.Name, Running: ScenarioC.Running, StartParameters: ScenarioC.StartParameters}
var ScenarioCUpdated = Scenario{ID: ScenarioC.ID, Name: "Scenario_Cupdated", Running: true, StartParameters: postgres.Jsonb{startParametersC}}
var ScenarioCUpdated_response = ScenarioResponse{ID: ScenarioCUpdated.ID, Name: ScenarioCUpdated.Name, Running: ScenarioCUpdated.Running, StartParameters: ScenarioCUpdated.StartParameters}

// Simulation Models

var SimulationModelA = SimulationModel{
	ID:              1,
	Name:            "SimulationModel_A",
	OutputLength:    1,
	InputLength:     1,
	ScenarioID:      1,
	SimulatorID:     1,
	StartParameters: postgres.Jsonb{startParametersA},
}

var SimulationModelA_response = SimulationModelResponse{
	ID:              SimulationModelA.ID,
	Name:            SimulationModelA.Name,
	InputLength:     SimulationModelA.InputLength,
	OutputLength:    SimulationModelA.OutputLength,
	ScenarioID:      SimulationModelA.ScenarioID,
	SimulatorID:     SimulationModelA.SimulatorID,
	StartParameters: SimulationModelA.StartParameters,
}

var SimulationModelB = SimulationModel{
	ID:              2,
	Name:            "SimulationModel_B",
	OutputLength:    1,
	InputLength:     1,
	ScenarioID:      1,
	SimulatorID:     1,
	StartParameters: postgres.Jsonb{startParametersB},
}

var SimulationModelB_response = SimulationModelResponse{
	ID:              SimulationModelB.ID,
	Name:            SimulationModelB.Name,
	InputLength:     SimulationModelB.InputLength,
	OutputLength:    SimulationModelB.OutputLength,
	ScenarioID:      SimulationModelB.ScenarioID,
	SimulatorID:     SimulationModelB.SimulatorID,
	StartParameters: SimulationModelB.StartParameters,
}

var SimulationModelC = SimulationModel{
	ID:              3,
	Name:            "SimulationModel_C",
	OutputLength:    1,
	InputLength:     1,
	ScenarioID:      1,
	SimulatorID:     1,
	StartParameters: postgres.Jsonb{startParametersC},
}

var SimulationModelC_response = SimulationModelResponse{
	ID:              SimulationModelC.ID,
	Name:            SimulationModelC.Name,
	InputLength:     SimulationModelC.InputLength,
	OutputLength:    SimulationModelC.OutputLength,
	ScenarioID:      SimulationModelC.ScenarioID,
	SimulatorID:     SimulationModelC.SimulatorID,
	StartParameters: SimulationModelC.StartParameters,
}

var SimulationModelCUpdated = SimulationModel{
	ID:              SimulationModelC.ID,
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
	ID:              SimulationModelCUpdated.ID,
	Name:            SimulationModelCUpdated.Name,
	InputLength:     SimulationModelCUpdated.InputLength,
	OutputLength:    SimulationModelCUpdated.OutputLength,
	ScenarioID:      SimulationModelCUpdated.ScenarioID,
	SimulatorID:     SimulationModelCUpdated.SimulatorID,
	StartParameters: SimulationModelCUpdated.StartParameters,
}

// Signals

var OutSignalA = Signal{
	Name:              "outSignal_A",
	Direction:         "out",
	Index:             0,
	Unit:              "V",
	SimulationModelID: 1,
}

var OutSignalA_response = SignalResponse{
	Name:              OutSignalA.Name,
	Direction:         OutSignalA.Direction,
	Index:             OutSignalA.Index,
	Unit:              OutSignalA.Unit,
	SimulationModelID: OutSignalA.SimulationModelID,
}

var OutSignalB = Signal{
	Name:              "outSignal_B",
	Direction:         "out",
	Index:             1,
	Unit:              "V",
	SimulationModelID: 1,
}

var OutSignalB_response = SignalResponse{
	Name:              OutSignalB.Name,
	Direction:         OutSignalB.Direction,
	Index:             OutSignalB.Index,
	Unit:              OutSignalB.Unit,
	SimulationModelID: OutSignalB.SimulationModelID,
}

var InSignalA = Signal{
	Name:              "inSignal_A",
	Direction:         "in",
	Index:             0,
	Unit:              "A",
	SimulationModelID: 1,
}

var InSignalA_response = SignalResponse{
	Name:              InSignalA.Name,
	Direction:         InSignalA.Direction,
	Index:             InSignalA.Index,
	Unit:              InSignalA.Unit,
	SimulationModelID: InSignalA.SimulationModelID,
}

var InSignalB = Signal{
	Name:              "inSignal_B",
	Direction:         "in",
	Index:             1,
	Unit:              "A",
	SimulationModelID: 1,
}

var InSignalB_response = SignalResponse{
	Name:              InSignalB.Name,
	Direction:         InSignalB.Direction,
	Index:             InSignalB.Index,
	Unit:              InSignalB.Unit,
	SimulationModelID: InSignalB.SimulationModelID,
}

var InSignalC = Signal{
	Name:              "inSignal_C",
	Direction:         "in",
	Index:             2,
	Unit:              "A",
	SimulationModelID: 1,
}

var InSignalC_response = SignalResponse{
	Name:              InSignalC.Name,
	Direction:         InSignalC.Direction,
	Index:             InSignalC.Index,
	Unit:              InSignalC.Unit,
	SimulationModelID: InSignalC.SimulationModelID,
}

var InSignalCUpdated = Signal{
	Name:              "inSignalupdated_C",
	Direction:         InSignalC.Direction,
	Index:             InSignalC.Index,
	Unit:              "Ohm",
	SimulationModelID: InSignalC.SimulationModelID,
}

var InSignalCUpdated_response = SignalResponse{
	Name:              InSignalCUpdated.Name,
	Direction:         InSignalCUpdated.Direction,
	Index:             InSignalCUpdated.Index,
	Unit:              InSignalCUpdated.Unit,
	SimulationModelID: InSignalCUpdated.SimulationModelID,
}

// Dashboards

var DashboardA = Dashboard{ID: 1, Name: "Dashboard_A", Grid: 15, ScenarioID: 1}
var DashboardA_response = DashboardResponse{ID: DashboardA.ID, Name: DashboardA.Name, Grid: DashboardA.Grid, ScenarioID: DashboardA.ScenarioID}
var DashboardB = Dashboard{ID: 2, Name: "Dashboard_B", Grid: 10, ScenarioID: 1}
var DashboardB_response = DashboardResponse{ID: DashboardB.ID, Name: DashboardB.Name, Grid: DashboardB.Grid, ScenarioID: DashboardB.ScenarioID}
var DashboardC = Dashboard{ID: 3, Name: "Dashboard_C", Grid: 25, ScenarioID: 1}
var DashboardC_response = DashboardResponse{ID: DashboardC.ID, Name: DashboardC.Name, Grid: DashboardC.Grid, ScenarioID: DashboardC.ScenarioID}
var DashboardCUpdated = Dashboard{ID: DashboardC.ID, Name: "Dashboard_Cupdated", Grid: 24, ScenarioID: DashboardC.ScenarioID}
var DashboardCUpdated_response = DashboardResponse{ID: DashboardCUpdated.ID, Name: DashboardCUpdated.Name, Grid: DashboardCUpdated.Grid, ScenarioID: DashboardCUpdated.ScenarioID}
