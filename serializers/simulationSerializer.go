package serializers

import (
	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
)

type SimulationsSerializerNoAssoc struct {
	Ctx         *gin.Context
	Simulations []common.Simulation
}

func (self *SimulationsSerializerNoAssoc) Response() []SimulationResponseNoAssoc {
	response := []SimulationResponseNoAssoc{}
	for _, simulation := range self.Simulations {
		serializer := SimulationSerializerNoAssoc{self.Ctx, simulation}
		response = append(response, serializer.Response())
	}
	return response
}

type SimulationSerializerNoAssoc struct {
	Ctx *gin.Context
	common.Simulation
}

type SimulationResponseNoAssoc struct {
	Name    string `json:"Name"`
	ID      uint   `json:"SimulationID"`
	Running bool   `json:"Running"`
	//StartParams postgres.Jsonb `json:"Starting Parameters"`
}

func (self *SimulationSerializerNoAssoc) Response() SimulationResponseNoAssoc {
	response := SimulationResponseNoAssoc{
		Name:    self.Name,
		ID:      self.ID,
		Running: self.Running,
		//StartParams: self.StartParameters,
	}
	return response
}
