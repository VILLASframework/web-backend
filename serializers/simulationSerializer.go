package serializers

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm/dialects/postgres"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
)

type SimulationsSerializer struct {
	Ctx         *gin.Context
	Simulations []common.Simulation
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
	common.Simulation
}

type SimulationResponse struct {
	Name    string `json:"Name"`
	ID      uint   `json:"SimulationID"`
	Running bool   `json:"Running"`
	StartParams postgres.Jsonb `json:"Starting Parameters"`
}

func (self *SimulationSerializer) Response() SimulationResponse {
	response := SimulationResponse{
		Name:    self.Name,
		ID:      self.ID,
		Running: self.Running,
		StartParams: self.StartParameters,
	}
	return response
}
