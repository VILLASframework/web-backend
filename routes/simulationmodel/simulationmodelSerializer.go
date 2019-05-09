package simulationmodel

import (
	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
)

type SimulationModelsSerializerNoAssoc struct {
	Ctx         *gin.Context
	SimulationModels []common.SimulationModel
}

func (self *SimulationModelsSerializerNoAssoc) Response() []SimulationModelResponseNoAssoc {
	response := []SimulationModelResponseNoAssoc{}
	for _, simulationmodel := range self.SimulationModels {
		serializer := SimulationModelSerializerNoAssoc{self.Ctx, simulationmodel}
		response = append(response, serializer.Response())
	}
	return response
}

type SimulationModelSerializerNoAssoc struct {
	Ctx *gin.Context
	common.SimulationModel
}

type SimulationModelResponseNoAssoc struct {
	Name    		string `json:"Name"`
	OutputLength    int   `json:"OutputLength"`
	InputLength 	int   `json:"InputLength"`
	BelongsToSimulationID uint `json:"BelongsToSimulationID"`
	BelongsToSimulatorID uint `json:"BelongsToSimulatiorID"`
	//StartParams postgres.Jsonb `json:"Starting Parameters"`
	//Output Mapping
	//Input Mapping
}

func (self *SimulationModelSerializerNoAssoc) Response() SimulationModelResponseNoAssoc {
	response := SimulationModelResponseNoAssoc{
		Name:    		self.Name,
		OutputLength:   self.OutputLength,
		InputLength: 	self.InputLength,
		BelongsToSimulationID: self.BelongsToSimulationID,
		BelongsToSimulatorID: self.BelongsToSimulatorID,
		//StartParams: self.StartParameters,
		//InputMapping
		//OutputMapping
	}
	return response
}

