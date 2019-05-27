package widget

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/simulation"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/visualization"
)

func RegisterWidgetEndpoints(r *gin.RouterGroup){
	r.GET("/", GetWidgets)
	r.POST("/", AddWidget)
	//r.POST("/:widgetID", CloneWidget)
	r.PUT("/:widgetID", UpdateWidget)
	r.GET("/:widgetID", GetWidget)
	r.DELETE("/:widgetID", DeleteWidget)
}

// GetWidgets godoc
// @Summary Get all widgets of visualization
// @ID GetWidgets
// @Produce  json
// @Tags widgets
// @Success 200 {array} common.WidgetResponse "Array of widgets to which belong to visualization"
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param visualizationID query int true "Visualization ID"
// @Router /widgets [get]
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

// AddWidget godoc
// @Summary Add a widget to a visualization
// @ID AddWidget
// @Accept json
// @Produce json
// @Tags widgets
// @Param inputWidget body common.WidgetResponse true "Widget to be added incl. ID of visualization"
// @Success 200 "OK."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Router /widgets [post]
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

// UpdateWidget godoc
// @Summary Update a widget
// @ID UpdateWidget
// @Tags widgets
// @Accept json
// @Produce json
// @Param inputWidget body common.WidgetResponse true "Widget to be updated"
// @Success 200 "OK."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param widgetID path int true "Widget ID"
// @Router /widgets/{widgetID} [put]
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

// GetWidget godoc
// @Summary Get a widget
// @ID GetWidget
// @Tags widgets
// @Produce json
// @Success 200 {object} common.WidgetResponse "Requested widget."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param widgetID path int true "Widget ID"
// @Router /widgets/{widgetID} [get]
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


// DeleteWidget godoc
// @Summary Delete a widget
// @ID DeleteWidget
// @Tags widgets
// @Produce json
// @Success 200 "OK."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param widgetID path int true "Widget ID"
// @Router /widgets/{widgetID} [delete]
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


