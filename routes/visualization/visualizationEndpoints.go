package visualization

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/simulation"
)

func RegisterVisualizationEndpoints(r *gin.RouterGroup) {

	r.GET("/", getVisualizations)
	r.POST("/", addVisualization)
	//r.POST("/:visualizationID", cloneVisualization)
	r.PUT("/:visualizationID", updateVisualization)
	r.GET("/:visualizationID", getVisualization)
	r.DELETE("/:visualizationID", deleteVisualization)

}

// getVisualizations godoc
// @Summary Get all visualizations of simulation
// @ID getVisualizations
// @Produce  json
// @Tags visualizations
// @Success 200 {array} common.VisualizationResponse "Array of visualizations to which belong to simulation"
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param simulationID query int true "Simulation ID"
// @Router /visualizations [get]
func getVisualizations(c *gin.Context) {

	simID, err := strconv.Atoi(c.Request.URL.Query().Get("simulationID"))
	if err != nil {
		errormsg := fmt.Sprintf("Bad request. No or incorrect format of simulationID query parameter")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errormsg,
		})
		return
	}

	var sim simulation.Simulation
	err = sim.ByID(uint(simID))
	if common.ProvideErrorResponse(c, err) {
		return
	}

	db := common.GetDB()
	var visualizations []common.Visualization
	err = db.Order("ID asc").Model(sim).Related(&visualizations, "Visualizations").Error

	serializer := common.VisualizationsSerializer{c, visualizations}
	c.JSON(http.StatusOK, gin.H{
		"visualizations": serializer.Response(),
	})
}

// addVisualization godoc
// @Summary Add a visualization to a simulation
// @ID addVisualization
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
func addVisualization(c *gin.Context) {

	var vis Visualization
	err := c.BindJSON(&vis)
	if err != nil {
		errormsg := "Bad request. Error binding form data to JSON: " + err.Error()
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errormsg,
		})
		return
	}

	// add visualization to DB and add association to simulation
	err = vis.addToSimulation(vis.SimulationID)
	if common.ProvideErrorResponse(c, err) == false {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK",
		})
	}
}

func cloneVisualization(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}

// updateVisualization godoc
// @Summary Update a visualization
// @ID updateVisualization
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
func updateVisualization(c *gin.Context) {

	visID, err := common.GetVisualizationID(c)
	if err != nil {
		return
	}

	var modifiedVis Visualization
	err = c.BindJSON(&modifiedVis)
	if err != nil {
		errormsg := "Bad request. Error binding form data to JSON: " + err.Error()
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errormsg,
		})
		return
	}

	var vis Visualization
	err = vis.ByID(uint(visID))
	if common.ProvideErrorResponse(c, err) {
		return
	}

	err = vis.update(modifiedVis)
	if common.ProvideErrorResponse(c, err) == false {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK",
		})
	}
}

// getVisualization godoc
// @Summary Get a visualization
// @ID getVisualization
// @Tags visualizations
// @Produce json
// @Success 200 {object} common.VisualizationResponse "Requested visualization."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param visualizationID path int true "Visualization ID"
// @Router /visualizations/{visualizationID} [get]
func getVisualization(c *gin.Context) {

	visID, err := common.GetVisualizationID(c)
	if err != nil {
		return
	}

	var vis Visualization
	err = vis.ByID(uint(visID))
	if common.ProvideErrorResponse(c, err) {
		return
	}

	serializer := common.VisualizationSerializer{c, vis.Visualization}
	c.JSON(http.StatusOK, gin.H{
		"visualization": serializer.Response(),
	})
}

// deleteVisualization godoc
// @Summary Delete a visualization
// @ID deleteVisualization
// @Tags visualizations
// @Produce json
// @Success 200 "OK."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param visualizationID path int true "Visualization ID"
// @Router /visualizations/{visualizationID} [delete]
func deleteVisualization(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}
