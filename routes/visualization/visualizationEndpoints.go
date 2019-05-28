package visualization

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/simulation"
)

func RegisterVisualizationEndpoints(r *gin.RouterGroup){

	r.GET("/", GetVisualizations)
	r.POST("/", AddVisualization)
	//r.POST("/:visualizationID", CloneVisualization)
	r.PUT("/:visualizationID", UpdateVisualization)
	r.GET("/:visualizationID", GetVisualization)
	r.DELETE("/:visualizationID", DeleteVisualization)

}

// GetVisualizations godoc
// @Summary Get all visualizations of simulation
// @ID GetVisualizations
// @Produce  json
// @Tags visualizations
// @Success 200 {array} common.VisualizationResponse "Array of visualizations to which belong to simulation"
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param simulationID query int true "Simulation ID"
// @Router /visualizations [get]
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

// AddVisualization godoc
// @Summary Add a visualization to a simulation
// @ID AddVisualization
// @Accept json
// @Produce json
// @Tags visualizations
// @Param inputVis body common.VisualizationResponse true "Visualization to be added incl. ID of simulation"
// @Success 200 "OK."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Router /visualizations [post]
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

// UpdateVisualization godoc
// @Summary Update a visualization
// @ID UpdateVisualization
// @Tags visualizations
// @Accept json
// @Produce json
// @Param inputVis body common.VisualizationResponse true "Visualization to be updated"
// @Success 200 "OK."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param visualizationID path int true "Visualization ID"
// @Router /visualizations/{visualizationID} [put]
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

// GetVisualization godoc
// @Summary Get a visualization
// @ID GetVisualization
// @Tags visualizations
// @Produce json
// @Success 200 {object} common.VisualizationResponse "Requested visualization."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param visualizationID path int true "Visualization ID"
// @Router /visualizations/{visualizationID} [get]
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

// DeleteVisualization godoc
// @Summary Delete a visualization
// @ID DeleteVisualization
// @Tags visualizations
// @Produce json
// @Success 200 "OK."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param visualizationID path int true "Visualization ID"
// @Router /visualizations/{visualizationID} [delete]
func DeleteVisualization(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}


