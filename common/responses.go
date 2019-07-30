package common

import "github.com/jinzhu/gorm/dialects/postgres"

type UserResponse struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	Mail     string `json:"mail"`
	ID       uint   `json:"id"`
}

type ScenarioResponse struct {
	Name            string         `json:"name"`
	ID              uint           `json:"id"`
	Running         bool           `json:"running"`
	StartParameters postgres.Jsonb `json:"startParameters"`
}

type SimulationModelResponse struct {
	ID              uint           `json:"id"`
	Name            string         `json:"name"`
	OutputLength    int            `json:"outputLength"`
	InputLength     int            `json:"inputLength"`
	ScenarioID      uint           `json:"scenarioID"`
	SimulatorID     uint           `json:"simulatorID"`
	StartParameters postgres.Jsonb `json:"startParameters"`
}

type SimulatorResponse struct {
	ID            uint           `json:"id"`
	UUID          string         `json:"uuid"`
	Host          string         `json:"host"`
	Modeltype     string         `json:"modelType"`
	Uptime        int            `json:"uptime"`
	State         string         `json:"state"`
	StateUpdateAt string         `json:"stateUpdateAt"`
	Properties    postgres.Jsonb `json:"properties"`
	RawProperties postgres.Jsonb `json:"rawProperties"`
}

type DashboardResponse struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	Grid       int    `json:"grid"`
	ScenarioID uint   `json:"scenarioID"`
}

type WidgetResponse struct {
	ID               uint           `json:"id"`
	Name             string         `json:"name"`
	Type             string         `json:"type"`
	Width            uint           `json:"width"`
	Height           uint           `json:"height"`
	MinWidth         uint           `json:"minWidth"`
	MinHeight        uint           `json:"minHeight"`
	X                int            `json:"x"`
	Y                int            `json:"y"`
	Z                int            `json:"z"`
	DashboardID      uint           `json:"dashboardID"`
	IsLocked         bool           `json:"isLocked"`
	CustomProperties postgres.Jsonb `json:"customProperties"`
}

type FileResponse struct {
	Name              string `json:"name"`
	ID                uint   `json:"id"`
	Type              string `json:"type"`
	Size              uint   `json:"size"`
	ImageWidth        uint   `json:"imageHeight"`
	ImageHeight       uint   `json:"imageWidth"`
	Date              string `json:"date"`
	WidgetID          uint   `json:"widgetID"`
	SimulationModelID uint   `json:"simulationModelID"`
}

type SignalResponse struct {
	Name              string `json:"name"`
	Unit              string `json:"unit"`
	Index             uint   `json:"index"`
	Direction         string `json:"direction"`
	SimulationModelID uint   `json:"simulationModelID"`
}

// Response messages

type ResponseMsg struct {
	Message string `json:"message"`
}

type ResponseMsgUsers struct {
	Users []User `json:"users"`
}

type ResponseMsgUser struct {
	User UserResponse `json:"user"`
}

type ResponseMsgScenarios struct {
	Scenarios []ScenarioResponse `json:"scenarios"`
}

type ResponseMsgScenario struct {
	Scenario ScenarioResponse `json:"scenario"`
}

type ResponseMsgSimulationModels struct {
	SimulationModels []SimulationModelResponse `json:"models"`
}

type ResponseMsgSimulationModel struct {
	SimulationModel SimulationModelResponse `json:"model"`
}

type ResponseMsgSignals struct {
	Signals []SignalResponse `json:"signals"`
}

type ResponseMsgSignal struct {
	Signal SignalResponse `json:"signal"`
}

type ResponseMsgDashboards struct {
	Dashboards []DashboardResponse `json:"dashboards"`
}

type ResponseMsgDashboard struct {
	Dashboard DashboardResponse `json:"dashboard"`
}

type ResponseMsgWidgets struct {
	Widgets []WidgetResponse `json:"widgets"`
}

type ResponseMsgWidget struct {
	Widget WidgetResponse `json:"widget"`
}

type ResponseMsgSimulators struct {
	Simulators []SimulatorResponse `json:"simulators"`
}

type ResponseMsgSimulator struct {
	Simulator SimulatorResponse `json:"simulator"`
}

type ResponseMsgFiles struct {
	Files []FileResponse `json:"files"`
}

type ResponseMsgFile struct {
	File FileResponse `json:"file"`
}
