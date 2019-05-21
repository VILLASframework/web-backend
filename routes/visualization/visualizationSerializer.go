package visualization

import (
	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/widget"
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
	Widgets	   []widget.WidgetResponse
}

func (self *VisualizationSerializer) Response() VisualizationResponse {

	w, _, _:= widget.FindVisualizationWidgets(&self.Visualization)
	widgetsSerializer := widget.WidgetsSerializer{self.Ctx, w}

	response := VisualizationResponse{
		Name:    	self.Name,
		Grid:	 	self.Grid,
		SimulationID: self.SimulationID,
		Widgets:    widgetsSerializer.Response(),
	}
	return response
}

