package serializers

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm/dialects/postgres"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
)

type ModelsSerializer struct {
	Ctx         *gin.Context
	Models []common.Model
}

func (self *ModelsSerializer) Response() []ModelResponse {
	response := []ModelResponse{}
	for _, model := range self.Models {
		serializer := ModelSerializer{self.Ctx, model}
		response = append(response, serializer.Response())
	}
	return response
}

type ModelSerializer struct {
	Ctx *gin.Context
	common.Model
}

type ModelResponse struct {
	Name    		string `json:"Name"`
	OutputLength    int   `json:"OutputLength"`
	InputLength 	int   `json:"InputLength"`
	SimulationID uint `json:"SimulationID"`
	SimulatorID uint `json:"SimulatorID"`
	StartParams postgres.Jsonb `json:"StartParams"`
	//StartParams postgres.Jsonb `json:"Starting Parameters"`
	//Output Mapping
	//Input Mapping
}

func (self *ModelSerializer) Response() ModelResponse {
	response := ModelResponse{
		Name:    		self.Name,
		OutputLength:   self.OutputLength,
		InputLength: 	self.InputLength,
		SimulationID: self.SimulationID,
		SimulatorID: self.SimulatorID,
		StartParams: self.StartParameters,
		//InputMapping
		//OutputMapping
	}
	return response
}

