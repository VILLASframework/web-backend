package common

import (
	"time"
)

type UserResponse struct {
	Username string `json:"Username"`
	Role     string `json:"Role"`
	Mail     string `json:"Mail"`
}

type SimulationResponse struct {
	Name        string `json:"Name"`
	ID          uint   `json:"SimulationID"`
	Running     bool   `json:"Running"`
	StartParams string `json:"Starting Parameters"`
}

type SimulationModelResponse struct {
	ID           uint   `json:"ID"`
	Name         string `json:"Name"`
	OutputLength int    `json:"OutputLength"`
	InputLength  int    `json:"InputLength"`
	SimulationID uint   `json:"SimulationID"`
	SimulatorID  uint   `json:"SimulatorID"`
	StartParams  string `json:"StartParams"`
}

type SimulatorResponse struct {
	UUID          string    `json:"UUID"`
	Host          string    `json:"Host"`
	ModelType     string    `json:"ModelType"`
	Uptime        int       `json:"Uptime"`
	State         string    `json:"State"`
	StateUpdateAt time.Time `json:"StateUpdateAt"`
	Properties    string    `json:"Properties"`
	RawProperties string    `json:"RawProperties"`
}

type VisualizationResponse struct {
	Name         string `json:"Name"`
	Grid         int    `json:"Grid"`
	SimulationID uint   `json:"SimulationID"`
}

type WidgetResponse struct {
	Name             string `json:"Name"`
	Type             string `json:"Type"`
	Width            uint   `json:"Width"`
	Height           uint   `json:"Height"`
	MinWidth         uint   `json:"MinWidth"`
	MinHeight        uint   `json:"MinHeight"`
	X                int    `json:"X"`
	Y                int    `json:"Y"`
	Z                int    `json:"Z"`
	VisualizationID  uint   `json:"VisualizationID"`
	IsLocked         bool   `json:"IsLocked"`
	CustomProperties string `json:"CustomProperties"`
}

type FileResponse struct {
	Name string    `json:"Name"`
	ID   uint      `json:"FileID"`
	Path string    `json:"Path"`
	Type string    `json:"Type"`
	Size uint      `json:"Size"`
	H    uint      `json:"ImageHeight"`
	W    uint      `json:"ImageWidth"`
	Date time.Time `json:"Date"`
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

type ResponseMsgSimulations struct {
	Simulations []SimulationResponse `json:"simulations"`
}

type ResponseMsgSimulation struct {
	Simulation SimulationResponse `json:"simulation"`
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
