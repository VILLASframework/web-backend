package serializers

import (
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/queries"
	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
)

type VisualizationsSerializer struct {
	Ctx         *gin.Context
	Visualizations []common.Visualization
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
	common.Visualization
}

type VisualizationResponse struct {
	Name    string `json:"Name"`
	Grid 		int   `json:"Grid"`
	SimulationID uint  `json:"SimulationID"`
	Widgets	   []WidgetResponse
}

func (self *VisualizationSerializer) Response() VisualizationResponse {

	w, _, _:= queries.FindVisualizationWidgets(&self.Visualization)
	widgetsSerializer := WidgetsSerializer{self.Ctx, w}

	response := VisualizationResponse{
		Name:    	self.Name,
		Grid:	 	self.Grid,
		SimulationID: self.SimulationID,
		Widgets:    widgetsSerializer.Response(),
	}
	return response
}

