package widget

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/simulation"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/visualization"
)

func RegisterWidgetEndpoints(r *gin.RouterGroup){
	r.GET("/:simulationID/visualization/:visualizationID/widgets", GetWidgets)
	r.POST("/:simulationID/visualization/:visualizationID/widget", AddWidget)
	r.POST("/:simulationID/visualization/:visualizationID/widget:widgetID", CloneWidget)
	r.PUT("/:simulationID/visualization/:visualizationID/widget/:widgetID", UpdateWidget)
	r.GET("/:simulationID/visualization/:visualizationID/widget/:widgetID", GetWidget)
	r.DELETE("/:simulationID/visualization/:visualizationID/widget/:widgetID", DeleteWidget)
}

func GetWidgets(c *gin.Context) {

	simID, err := common.GetSimulationID(c)
	if err != nil {
		return
	}

	sim, err := simulation.FindSimulation(simID)
	if common.ProvideErrorResponse(c, err) {
		return
	}

	visID, err := common.GetVisualizationID(c)
	if err != nil {
		return
	}

	vis, err := visualization.FindVisualizationOfSim(&sim, visID)
	if common.ProvideErrorResponse(c, err) {
		return
	}

	widgets,_, err := FindWidgetsOfVisualization(&vis)
	if common.ProvideErrorResponse(c, err) {
		return
	}

	serializer := common.WidgetsSerializer{c, widgets}
	c.JSON(http.StatusOK, gin.H{
		"widgets": serializer.Response(),
	})
}

func AddWidget(c *gin.Context) {

	simID, err := common.GetSimulationID(c)
	if err != nil {
		return
	}

	sim, err := simulation.FindSimulation(simID)
	if common.ProvideErrorResponse(c, err) {
		return
	}

	visID, err := common.GetVisualizationID(c)
	if err != nil {
		return
	}

	vis, err := visualization.FindVisualizationOfSim(&sim, visID)
	if common.ProvideErrorResponse(c, err) {
		return
	}

	var widget_input common.Widget
	err = c.BindJSON(&widget_input)
	if err != nil {
		errormsg := "Bad request. Error binding form data to JSON: " + err.Error()
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errormsg,
		})
		return
	}

	err = AddWidgetToVisualization(&vis, &widget_input)
	if common.ProvideErrorResponse(c, err) == false {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK.",
		})
	}


}

func CloneWidget(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}

func UpdateWidget(c *gin.Context) {
	simID, err := common.GetSimulationID(c)
	if err != nil {
		return
	}

	sim, err := simulation.FindSimulation(simID)
	if common.ProvideErrorResponse(c, err) {
		return
	}

	visID, err := common.GetVisualizationID(c)
	if err != nil {
		return
	}

	vis, err := visualization.FindVisualizationOfSim(&sim, visID)
	if common.ProvideErrorResponse(c, err) {
		return
	}

	widgetID, err := common.GetWidgetID(c)
	if err != nil {
		return
	}

	var widget_input common.Widget
	err = c.BindJSON(&widget_input)
	if err != nil {
		errormsg := "Bad request. Error binding form data to JSON: " + err.Error()
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errormsg,
		})
		return
	}

	err = UpdateWidgetOfVisualization(&vis, widget_input, widgetID)
	if common.ProvideErrorResponse(c, err) == false {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK",
		})
	}
}

func GetWidget(c *gin.Context) {

	simID, err := common.GetSimulationID(c)
	if err != nil {
		return
	}

	sim, err := simulation.FindSimulation(simID)
	if common.ProvideErrorResponse(c, err) {
		return
	}

	visID, err := common.GetVisualizationID(c)
	if err != nil {
		return
	}

	vis, err := visualization.FindVisualizationOfSim(&sim, visID)
	if common.ProvideErrorResponse(c, err) {
		return
	}

	widgetID, err := common.GetWidgetID(c)
	if err != nil {
		return
	}

	widget, err := FindWidgetOfVisualization(&vis, widgetID)
	serializer := common.WidgetSerializer{c, widget}
	c.JSON(http.StatusOK, gin.H{
		"widget": serializer.Response(),
	})
}

func DeleteWidget(c *gin.Context) {

	// simID, err := GetSimulationID(c)
	// if err != nil {
	// 	return
	// }
	//
	// sim, err := queries.FindSimulation(simID)
	// if common.ProvideErrorResponse(c, err) {
	// 	return
	// }
	//
	// visID, err := GetVisualizationID(c)
	// if err != nil {
	// 	return
	// }
	//
	// visualization, err := queries.FindVisualizationOfSim(&sim, visID)
	// if common.ProvideErrorResponse(c, err) {
	// 	return
	// }
	//
	// widgetID, err := GetWidgetID(c)
	// if err != nil {
	// 	return
	// }
	//
	// widget, err := queries.FindWidgetOfVisualization(&visualization, widgetID)
	// if common.ProvideErrorResponse(c, err) {
	// 	return
	// }

	// TODO delete files of widget in DB and on disk


	// TODO Delete widget itself + association with visualization

	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}


