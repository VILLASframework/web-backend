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
	}

	// Associated models MUST NOT called with assoc=true otherwise we
	// will have an infinite loop due to the circular dependencies
	if assoc {

		// TODO: maybe all those should be made in one transaction

		//simulations, _, _ := simulation.FindUserSimulations(&self.User)
		//simulationsSerializer :=
		//	SimulationsSerializer{self.Ctx, simulations}

		// Add the associated models to the response
		//response.Simulations = simulationsSerializer.Response()
	}

	return response
}

// Simulation/s Serializers

type SimulationsSerializer struct {
	Ctx         *gin.Context
	Simulations []Simulation
}

func (self *SimulationsSerializer) Response() []SimulationResponse {
	response := []SimulationResponse{}
	for _, simulation := range self.Simulations {
		serializer := SimulationSerializer{self.Ctx, simulation}
		response = append(response, serializer.Response())
	}
	return response
}

type SimulationSerializer struct {
	Ctx *gin.Context
	Simulation
}

func (self *SimulationSerializer) Response() SimulationResponse {
	response := SimulationResponse{
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
		SimulationID: self.SimulationID,
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
		UUID:          self.UUID,
		Host:          self.Host,
		ModelType:     self.Modeltype,
		Uptime:        self.Uptime,
		State:         self.State,
		StateUpdateAt: self.StateUpdateAt,
	}
	return response
}

// Visualization/s Serializers

type VisualizationsSerializer struct {
	Ctx            *gin.Context
	Visualizations []Visualization
}

func (self *VisualizationsSerializer) Response() []VisualizationResponse {
	response := []VisualizationResponse{}
	for _, visualization := range self.Visualizations {
		serializer := VisualizationSerializer{self.Ctx, visualization}
		response = append(response, serializer.Response())
	}
	return response
}

type VisualizationSerializer struct {
	Ctx *gin.Context
	Visualization
}

func (self *VisualizationSerializer) Response() VisualizationResponse {

	response := VisualizationResponse{
		Name:         self.Name,
		Grid:         self.Grid,
		SimulationID: self.SimulationID,
		ID:           self.ID,
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
		ID:              self.ID,
		Name:            self.Name,
		Type:            self.Type,
		Width:           self.Width,
		Height:          self.Height,
		MinWidth:        self.MinWidth,
		MinHeight:       self.MinHeight,
		X:               self.X,
		Y:               self.Y,
		Z:               self.Z,
		VisualizationID: self.VisualizationID,
		IsLocked:        self.IsLocked,
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
		Path: self.Path,
		Type: self.Type,
		Size: self.Size,
		H:    self.ImageHeight,
		W:    self.ImageWidth,
		// Date
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
