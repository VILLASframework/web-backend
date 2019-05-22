package endpoints

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/queries"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/serializers"
)


func widgetReadAllEp(c *gin.Context) {

	simID, err := GetSimulationID(c)
	if err != nil {
		return
	}

	sim, err := queries.FindSimulation(simID)
	if common.ProvideErrorResponse(c, err) {
		return
	}

	visID, err := GetVisualizationID(c)
	if err != nil {
		return
	}

	vis, err := queries.FindVisualizationOfSim(&sim, visID)
	if common.ProvideErrorResponse(c, err) {
		return
	}

	widgets,_, err := queries.FindWidgetsOfVisualization(&vis)
	if common.ProvideErrorResponse(c, err) {
		return
	}

	serializer := serializers.WidgetsSerializer{c, widgets}
	c.JSON(http.StatusOK, gin.H{
		"widgets": serializer.Response(),
	})
}

func widgetRegistrationEp(c *gin.Context) {

	simID, err := GetSimulationID(c)
	if err != nil {
		return
	}

	sim, err := queries.FindSimulation(simID)
	if common.ProvideErrorResponse(c, err) {
		return
	}

	visID, err := GetVisualizationID(c)
	if err != nil {
		return
	}

	vis, err := queries.FindVisualizationOfSim(&sim, visID)
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

	err = queries.AddWidgetToVisualization(&vis, &widget_input)
	if common.ProvideErrorResponse(c, err) == false {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK.",
		})
	}


}

func widgetCloneEp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}

func widgetUpdateEp(c *gin.Context) {
	simID, err := GetSimulationID(c)
	if err != nil {
		return
	}

	sim, err := queries.FindSimulation(simID)
	if common.ProvideErrorResponse(c, err) {
		return
	}

	visID, err := GetVisualizationID(c)
	if err != nil {
		return
	}

	vis, err := queries.FindVisualizationOfSim(&sim, visID)
	if common.ProvideErrorResponse(c, err) {
		return
	}

	widgetID, err := GetWidgetID(c)
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

	err = queries.UpdateWidgetOfVisualization(&vis, widget_input, widgetID)
	if common.ProvideErrorResponse(c, err) == false {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK",
		})
	}
}

func widgetReadEp(c *gin.Context) {

	simID, err := GetSimulationID(c)
	if err != nil {
		return
	}

	sim, err := queries.FindSimulation(simID)
	if common.ProvideErrorResponse(c, err) {
		return
	}

	visID, err := GetVisualizationID(c)
	if err != nil {
		return
	}

	visualization, err := queries.FindVisualizationOfSim(&sim, visID)
	if common.ProvideErrorResponse(c, err) {
		return
	}

	widgetID, err := GetWidgetID(c)
	if err != nil {
		return
	}

	widget, err := queries.FindWidgetOfVisualization(&visualization, widgetID)
	serializer := serializers.WidgetSerializer{c, widget}
	c.JSON(http.StatusOK, gin.H{
		"widget": serializer.Response(),
	})
}

func widgetDeleteEp(c *gin.Context) {

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


func GetWidgetID(c *gin.Context) (int, error) {

	widgetID, err := strconv.Atoi(c.Param("WidgetID"))

	if err != nil {
		errormsg := fmt.Sprintf("Bad request. No or incorrect format of widget ID")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errormsg,
		})
		return -1, err
	} else {
		return widgetID, err

	}
}