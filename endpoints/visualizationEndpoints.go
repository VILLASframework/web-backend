package endpoints

import (
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/queries"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/serializers"
	"github.com/gin-gonic/gin"
	"fmt"
	"net/http"
	"strconv"
)

func visualizationReadAllEp(c *gin.Context) {

	simID, err := GetSimulationID(c)
	if err != nil {
		return
	}

	sim, err := queries.FindSimulation(simID)
	if common.ProvideErrorResponse(c, err) {
		return
	}

	allVisualizations, _, _ := queries.FindAllVisualizationsOfSim(&sim)
	serializer := serializers.VisualizationsSerializer{c, allVisualizations}
	c.JSON(http.StatusOK, gin.H{
		"visualizations": serializer.Response(),
	})
}

func visualizationRegistrationEp(c *gin.Context) {

	simID, err := GetSimulationID(c)
	if err != nil {
		return
	}

	sim, err := queries.FindSimulation(simID)
	if common.ProvideErrorResponse(c, err) {
		return
	}

	var vis common.Visualization
	err = c.BindJSON(&vis)
	if err != nil {
		errormsg := "Bad request. Error binding form data to JSON: " + err.Error()
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errormsg,
		})
		return
	}

	// add visualization to DB and add association to simulation
	err = queries.AddVisualizationToSim(&sim, &vis)
	if common.ProvideErrorResponse(c, err) == false {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK",
		})
	}
}

func visualizationCloneEp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}

func visualizationUpdateEp(c *gin.Context) {

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

	var vis common.Visualization
	err = c.BindJSON(&vis)
	if err != nil {
		errormsg := "Bad request. Error binding form data to JSON: " + err.Error()
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errormsg,
		})
		return
	}

	err = queries.UpdateVisualizationOfSim(&sim, vis, visID)
	if common.ProvideErrorResponse(c, err) == false {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK",
		})
	}
}

func visualizationReadEp(c *gin.Context) {

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

	serializer := serializers.VisualizationSerializer{c, visualization}
	c.JSON(http.StatusOK, gin.H{
		"visualization": serializer.Response(),
	})
}

func visualizationDeleteEp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}


func GetVisualizationID(c *gin.Context) (int, error) {

	simID, err := strconv.Atoi(c.Param("visualizationID"))

	if err != nil {
		errormsg := fmt.Sprintf("Bad request. No or incorrect format of visualization ID")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errormsg,
		})
		return -1, err
	} else {
		return simID, err

	}
}