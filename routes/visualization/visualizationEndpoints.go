package visualization

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/simulation"
)

func RegisterVisualizationEndpoints(r *gin.RouterGroup){

	r.GET("/:simulationID/visualizations", GetVisualizations)
	r.POST("/:simulationID/visualization", AddVisualization)
	r.POST("/:simulationID/visualization/:visualizationID", CloneVisualization)
	r.PUT("/:simulationID/visualization/:visualizationID", UpdateVisualization)
	r.GET("/:simulationID/visualization/:visualizationID", GetVisualization)
	r.DELETE("/:simulationID/visualization/:visualizationID", DeleteVisualization)

}

func GetVisualizations(c *gin.Context) {

	simID, err := common.GetSimulationID(c)
	if err != nil {
		return
	}

	sim, err := simulation.FindSimulation(simID)
	if common.ProvideErrorResponse(c, err) {
		return
	}

	allVisualizations, _, _ := FindAllVisualizationsOfSim(&sim)
	serializer := common.VisualizationsSerializer{c, allVisualizations}
	c.JSON(http.StatusOK, gin.H{
		"visualizations": serializer.Response(),
	})
}

func AddVisualization(c *gin.Context) {

	simID, err := common.GetSimulationID(c)
	if err != nil {
		return
	}

	sim, err := simulation.FindSimulation(simID)
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
	err = AddVisualizationToSim(&sim, &vis)
	if common.ProvideErrorResponse(c, err) == false {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK",
		})
	}
}

func CloneVisualization(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}

func UpdateVisualization(c *gin.Context) {

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

	var vis common.Visualization
	err = c.BindJSON(&vis)
	if err != nil {
		errormsg := "Bad request. Error binding form data to JSON: " + err.Error()
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errormsg,
		})
		return
	}

	err = UpdateVisualizationOfSim(&sim, vis, visID)
	if common.ProvideErrorResponse(c, err) == false {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK",
		})
	}
}

func GetVisualization(c *gin.Context) {

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

	visualization, err := FindVisualizationOfSim(&sim, visID)
	if common.ProvideErrorResponse(c, err) {
		return
	}

	serializer := common.VisualizationSerializer{c, visualization}
	c.JSON(http.StatusOK, gin.H{
		"visualization": serializer.Response(),
	})
}

func DeleteVisualization(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}


