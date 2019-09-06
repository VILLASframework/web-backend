package common

import (
	"github.com/gin-gonic/gin"
)

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
		ImageHeight:       self.ImageHeight,
		ImageWidth:        self.ImageWidth,
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
