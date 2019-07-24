package common

import (
	"github.com/gin-gonic/gin"
)

// User/s Serializers

type UsersSerializer struct {
	Ctx   *gin.Context
	Users []User
}

func (self *UsersSerializer) Response(assoc bool) []UserResponse {
	response := []UserResponse{}
	for _, user := range self.Users {
		serializer := UserSerializer{self.Ctx, user}
		response = append(response, serializer.Response(assoc))
	}
	return response
}

type UserSerializer struct {
	Ctx *gin.Context
	User
}

func (self *UserSerializer) Response(assoc bool) UserResponse {

	response := UserResponse{
		Username: self.Username,
		Role:     self.Role,
		Mail:     self.Mail,
		ID:       self.ID,
	}

	// Associated models MUST NOT called with assoc=true otherwise we
	// will have an infinite loop due to the circular dependencies
	if assoc {

		// TODO: maybe all those should be made in one transaction

		//scenarios, _, _ := scenario.FindUserScenarios(&self.User)
		//scenariosSerializer :=
		//	ScenariosSerializer{self.Ctx, scenarios}

		// Add the associated models to the response
		//response.Scenarios = scenariosSerializer.Response()
	}

	return response
}

// Scenario/s Serializers

type ScenariosSerializer struct {
	Ctx       *gin.Context
	Scenarios []Scenario
}

func (self *ScenariosSerializer) Response() []ScenarioResponse {
	response := []ScenarioResponse{}
	for _, so := range self.Scenarios {
		serializer := ScenarioSerializer{self.Ctx, so}
		response = append(response, serializer.Response())
	}
	return response
}

type ScenarioSerializer struct {
	Ctx *gin.Context
	Scenario
}

func (self *ScenarioSerializer) Response() ScenarioResponse {
	response := ScenarioResponse{
		Name:        self.Name,
		ID:          self.ID,
		Running:     self.Running,
		StartParams: self.StartParameters,
	}
	return response
}

// Model/s Serializers

type SimulationModelsSerializer struct {
	Ctx              *gin.Context
	SimulationModels []SimulationModel
}

func (self *SimulationModelsSerializer) Response() []SimulationModelResponse {
	response := []SimulationModelResponse{}
	for _, simulationmodel := range self.SimulationModels {
		serializer := SimulationModelSerializer{self.Ctx, simulationmodel}
		response = append(response, serializer.Response())
	}
	return response
}

type SimulationModelSerializer struct {
	Ctx *gin.Context
	SimulationModel
}

func (self *SimulationModelSerializer) Response() SimulationModelResponse {
	response := SimulationModelResponse{
		ID:           self.ID,
		Name:         self.Name,
		OutputLength: self.OutputLength,
		InputLength:  self.InputLength,
		ScenarioID:   self.ScenarioID,
		SimulatorID:  self.SimulatorID,
		StartParams:  self.StartParameters,
	}
	return response
}

// Simulator/s Serializers

type SimulatorsSerializer struct {
	Ctx        *gin.Context
	Simulators []Simulator
}

func (self *SimulatorsSerializer) Response() []SimulatorResponse {
	response := []SimulatorResponse{}
	for _, simulator := range self.Simulators {
		serializer := SimulatorSerializer{self.Ctx, simulator}
		response = append(response, serializer.Response())
	}
	return response
}

type SimulatorSerializer struct {
	Ctx *gin.Context
	Simulator
}

func (self *SimulatorSerializer) Response() SimulatorResponse {

	response := SimulatorResponse{
		ID:            self.ID,
		UUID:          self.UUID,
		Host:          self.Host,
		Modeltype:     self.Modeltype,
		Uptime:        self.Uptime,
		State:         self.State,
		StateUpdateAt: self.StateUpdateAt,
		Properties:    self.Properties,
		RawProperties: self.RawProperties,
	}
	return response
}

// Dashboard/s Serializers

type DashboardsSerializer struct {
	Ctx        *gin.Context
	Dashboards []Dashboard
}

func (self *DashboardsSerializer) Response() []DashboardResponse {
	response := []DashboardResponse{}
	for _, dashboard := range self.Dashboards {
		serializer := DashboardSerializer{self.Ctx, dashboard}
		response = append(response, serializer.Response())
	}
	return response
}

type DashboardSerializer struct {
	Ctx *gin.Context
	Dashboard
}

func (self *DashboardSerializer) Response() DashboardResponse {

	response := DashboardResponse{
		Name:       self.Name,
		Grid:       self.Grid,
		ScenarioID: self.ScenarioID,
		ID:         self.ID,
	}
	return response
}

// Widget/s Serializers

type WidgetsSerializer struct {
	Ctx     *gin.Context
	Widgets []Widget
}

func (self *WidgetsSerializer) Response() []WidgetResponse {
	response := []WidgetResponse{}
	for _, widget := range self.Widgets {
		serializer := WidgetSerializer{self.Ctx, widget}
		response = append(response, serializer.Response())
	}
	return response
}

type WidgetSerializer struct {
	Ctx *gin.Context
	Widget
}

func (self *WidgetSerializer) Response() WidgetResponse {

	response := WidgetResponse{
		ID:          self.ID,
		Name:        self.Name,
		Type:        self.Type,
		Width:       self.Width,
		Height:      self.Height,
		MinWidth:    self.MinWidth,
		MinHeight:   self.MinHeight,
		X:           self.X,
		Y:           self.Y,
		Z:           self.Z,
		DashboardID: self.DashboardID,
		IsLocked:    self.IsLocked,
		//CustomProperties
	}
	return response
}

// File/s Serializers

type FilesSerializerNoAssoc struct {
	Ctx   *gin.Context
	Files []File
}

func (self *FilesSerializerNoAssoc) Response() []FileResponse {
	response := []FileResponse{}
	for _, files := range self.Files {
		serializer := FileSerializerNoAssoc{self.Ctx, files}
		response = append(response, serializer.Response())
	}
	return response
}

type FileSerializerNoAssoc struct {
	Ctx *gin.Context
	File
}

func (self *FileSerializerNoAssoc) Response() FileResponse {
	response := FileResponse{
		Name: self.Name,
		ID:   self.ID,
		//Path: self.Path,
		Type:              self.Type,
		Size:              self.Size,
		H:                 self.ImageHeight,
		W:                 self.ImageWidth,
		Date:              self.Date,
		WidgetID:          self.WidgetID,
		SimulationModelID: self.SimulationModelID,
	}
	return response
}

// Signal/s Serializers
type SignalsSerializer struct {
	Ctx     *gin.Context
	Signals []Signal
}

func (self *SignalsSerializer) Response() []SignalResponse {
	response := []SignalResponse{}
	for _, s := range self.Signals {
		serializer := SignalSerializer{self.Ctx, s}
		response = append(response, serializer.Response())
	}
	return response

}

type SignalSerializer struct {
	Ctx *gin.Context
	Signal
}

func (self *SignalSerializer) Response() SignalResponse {
	response := SignalResponse{
		Name:              self.Name,
		Unit:              self.Unit,
		Direction:         self.Direction,
		SimulationModelID: self.SimulationModelID,
		Index:             self.Index,
	}
	return response
}
