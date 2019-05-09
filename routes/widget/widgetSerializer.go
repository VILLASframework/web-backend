package widget

import (
	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
)

type WidgetsSerializer struct {
	Ctx         *gin.Context
	Widgets []common.Widget
}

func (self *WidgetsSerializer) Response() []WidgetResponse {
	response := []WidgetResponse{}
	for _, widget := range self.Widgets {
		serializer := WidgetSerializer{self.Ctx, widget}
		response = append(response, serializer.Response())
	}
	return response
}

type WidgetSerializer struct {
	Ctx *gin.Context
	common.Widget
}

type WidgetResponse struct {
	Name    	string `json:"Name"`
	Type      	string   `json:"Type"`
	Width 		uint   `json:"Width"`
	Height		uint   `json:"Height"`
	MinWidth    uint `json:"MinWidth"`
	MinHeight	uint `json:"MinHeight"`
	X			int `json:"X"`
	Y			int `json:"Y"`
	Z			int `json:"Z"`
	VisualizationID	uint `json:"VisualizationID"`
	IsLocked	bool `json:"IsLocked"`
	//CustomProperties
}

func (self *WidgetSerializer) Response() WidgetResponse {

	response := WidgetResponse{
		Name:    			self.Name,
		Type:	 			self.Type,
		Width: 				self.Width,
		Height: 			self.Height,
		MinWidth:			self.MinWidth,
		MinHeight:			self.MinHeight,
		X: 					self.X,
		Y: 					self.Y,
		Z: 					self.Z,
		VisualizationID:  	self.VisualizationID,
		IsLocked: 			self.IsLocked,
		//CustomProperties
	}
	return response
}
