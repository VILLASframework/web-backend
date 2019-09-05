package common

import "github.com/jinzhu/gorm/dialects/postgres"

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

type ResponseMsgFiles struct {
	Files []FileResponse `json:"files"`
}

type ResponseMsgFile struct {
	File FileResponse `json:"file"`
}
