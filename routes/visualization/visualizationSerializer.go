package visualization

import (
	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/project"
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
	UserID      uint   `json:"UserID"`
	Grid 		int   `json:"Grid"`
	ProjectID	uint   `json:"ProjectID"`
	Project    project.ProjectResponseNoAssoc
	Widgets	   []widget.WidgetResponse
}

func (self *VisualizationSerializer) Response() VisualizationResponse {

	// TODO: maybe all those should be made in one transaction
	p, _, _ := project.FindVisualizationProject(&self.Visualization)
	projectSerializer := project.ProjectSerializerNoAssoc{self.Ctx, p}

	w, _, _:= widget.FindVisualizationWidgets(&self.Visualization)
	widgetsSerializer := widget.WidgetsSerializer{self.Ctx, w}


	response := VisualizationResponse{
		Name:    	self.Name,
		UserID:  	self.UserID,
		Grid:	 	self.Grid,
		ProjectID: 	self.ProjectID,
		Project: 	projectSerializer.Response(),
		Widgets:    widgetsSerializer.Response(),
	}
	return response
}

