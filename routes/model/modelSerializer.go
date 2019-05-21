package model

import (
	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
)

type ModelsSerializerNoAssoc struct {
	Ctx         *gin.Context
	Models []common.Model
}

func (self *ModelsSerializerNoAssoc) Response() []ModelResponseNoAssoc {
	response := []ModelResponseNoAssoc{}
	for _, model := range self.Models {
		serializer := ModelSerializerNoAssoc{self.Ctx, model}
		response = append(response, serializer.Response())
	}
	return response
}

type ModelSerializerNoAssoc struct {
	Ctx *gin.Context
	common.Model
}

type ModelResponseNoAssoc struct {
	Name    		string `json:"Name"`
	OutputLength    int   `json:"OutputLength"`
	InputLength 	int   `json:"InputLength"`
	SimulationID uint `json:"SimulationID"`
	SimulatorID uint `json:"SimulatiorID"`
	//StartParams postgres.Jsonb `json:"Starting Parameters"`
	//Output Mapping
	//Input Mapping
}

func (self *ModelSerializerNoAssoc) Response() ModelResponseNoAssoc {
	response := ModelResponseNoAssoc{
		Name:    		self.Name,
		OutputLength:   self.OutputLength,
		InputLength: 	self.InputLength,
		SimulationID: self.SimulationID,
		SimulatorID: self.SimulatorID,
		//StartParams: self.StartParameters,
		//InputMapping
		//OutputMapping
	}
	return response
}

