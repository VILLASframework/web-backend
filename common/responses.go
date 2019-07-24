package common

import "github.com/jinzhu/gorm/dialects/postgres"

type UserResponse struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	Mail     string `json:"mail"`
	ID       uint   `json:"id"`
}

type ScenarioResponse struct {
	Name        string `json:"Name"`
	ID          uint   `json:"ID"`
	Running     bool   `json:"Running"`
	StartParams string `json:"Starting Parameters"`
}

type SimulationModelResponse struct {
	ID           uint   `json:"ID"`
	Name         string `json:"Name"`
	OutputLength int    `json:"OutputLength"`
	InputLength  int    `json:"InputLength"`
	ScenarioID   uint   `json:"ScenarioID"`
	SimulatorID  uint   `json:"SimulatorID"`
	StartParams  string `json:"StartParams"`
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
	ID         uint   `json:"ID"`
	Name       string `json:"Name"`
	Grid       int    `json:"Grid"`
	ScenarioID uint   `json:"ScenarioID"`
}

type WidgetResponse struct {
	ID               uint   `json:"ID"`
	Name             string `json:"Name"`
	Type             string `json:"Type"`
	Width            uint   `json:"Width"`
	Height           uint   `json:"Height"`
	MinWidth         uint   `json:"MinWidth"`
	MinHeight        uint   `json:"MinHeight"`
	X                int    `json:"X"`
	Y                int    `json:"Y"`
	Z                int    `json:"Z"`
	DashboardID      uint   `json:"DashboardID"`
	IsLocked         bool   `json:"IsLocked"`
	CustomProperties string `json:"CustomProperties"`
}

type FileResponse struct {
	Name              string `json:"Name"`
	ID                uint   `json:"ID"`
	Type              string `json:"Type"`
	Size              uint   `json:"Size"`
	H                 uint   `json:"ImageHeight"`
	W                 uint   `json:"ImageWidth"`
	Date              string `json:"Date"`
	WidgetID          uint   `json:"WidgetID"`
	SimulationModelID uint   `json:"SimulationModelID"`
}

type SignalResponse struct {
	Name              string `json:"Name"`
	Unit              string `json:"Unit"`
	Index             uint   `json:"Index"`
	Direction         string `json:"Direction"`
	SimulationModelID uint   `json:"SimulationModelID"`
}

// Response messages

type ResponseMsg struct {
	Message string `json:"message"`
}

type ResponseMsgUsers struct {
	Users []UserResponse `json:"users"`
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
