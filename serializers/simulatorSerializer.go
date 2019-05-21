package serializers

import (
	"time"

	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
)

type SimulatorsSerializer struct {
	Ctx   *gin.Context
	Simulators []common.Simulator
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
	common.Simulator
}

type SimulatorResponse struct {
	UUID	    	string `json:"UUID"`
	Host	    	string `json:"Host"`
	ModelType  		string `json:"ModelType"`
	Uptime     		int `json:"Uptime"`
	State    		string 	`json:"State"`
	StateUpdateAt 	time.Time `json:"StateUpdateAt"`
	// Properties
	// Raw Properties
}

func (self *SimulatorSerializer) Response() SimulatorResponse {

	response := SimulatorResponse{
		UUID:    		self.UUID,
		Host:    		self.Host,
		ModelType:      self.Modeltype,
		Uptime:        	self.Uptime,
		State:    		self.State,
		StateUpdateAt: 	self.StateUpdateAt,
	}
	return response
}
