package simulationmodel

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/simulation"
)

func RegisterSimulationModelEndpoints(r *gin.RouterGroup) {
	r.GET("", getSimulationModels)
	r.POST("", addSimulationModel)
	r.PUT("/:modelID", updateSimulationModel)
	r.GET("/:modelID", getSimulationModel)
	r.DELETE("/:modelID", deleteSimulationModel)
}

// getSimulationModels godoc
// @Summary Get all simulation models of simulation
// @ID getSimulationModels
// @Produce  json
// @Tags models
// @Success 200 {array} common.SimulationModelResponse "Array of models to which belong to simulation"
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param simulationID query int true "Simulation ID"
// @Router /models [get]
func getSimulationModels(c *gin.Context) {

	ok, sim := simulation.CheckPermissions(c, common.Read, "query", -1)
	if !ok {
		return
	}

	db := common.GetDB()
	var models []common.SimulationModel
	err := db.Order("ID asc").Model(sim).Related(&models, "Models").Error
	if common.ProvideErrorResponse(c, err) {
		return
	}

	serializer := common.SimulationModelsSerializer{c, models}
	c.JSON(http.StatusOK, gin.H{
		"models": serializer.Response(),
	})
}

// addSimulationModel godoc
// @Summary Add a simulation model to a simulation
// @ID addSimulationModel
// @Accept json
// @Produce json
// @Tags models
// @Param inputSimulationModel body common.SimulationModelResponse true "Simulation model to be added incl. IDs of simulation and simulator"
// @Success 200 "OK."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Router /models [post]
func addSimulationModel(c *gin.Context) {

	var newModel SimulationModel
	err := c.BindJSON(&newModel)
	if err != nil {
		errormsg := "Bad request. Error binding form data to JSON: " + err.Error()
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errormsg,
		})
		return
	}

	ok, _ := simulation.CheckPermissions(c, common.Create, "body", int(newModel.SimulationID))
	if !ok {
		return
	}

	err = newModel.addToSimulation()
	if common.ProvideErrorResponse(c, err) == false {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK.",
		})
	}
}

// updateSimulationModel godoc
// @Summary Update a simulation model
// @ID updateSimulationModel
// @Tags models
// @Accept json
// @Produce json
// @Param inputSimulationModel body common.SimulationModelResponse true "Simulation model to be updated"
// @Success 200 "OK."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param modelID path int true "Model ID"
// @Router /models/{modelID} [put]
func updateSimulationModel(c *gin.Context) {

	ok, m := CheckPermissions(c, common.Update, "path", -1)
	if !ok {
		return
	}

	var modifiedModel SimulationModel
	err := c.BindJSON(&modifiedModel)
	if err != nil {
		errormsg := "Bad request. Error binding form data to JSON: " + err.Error()
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errormsg,
		})
		return
	}

	err = m.Update(modifiedModel)
	if common.ProvideErrorResponse(c, err) == false {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK.",
		})
	}
}

// getSimulationModel godoc
// @Summary Get a simulation model
// @ID getSimulationModel
// @Tags models
// @Produce json
// @Success 200 {object} common.SimulationModelResponse "Requested simulation model."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param modelID path int true "Model ID"
// @Router /models/{modelID} [get]
func getSimulationModel(c *gin.Context) {

	ok, m := CheckPermissions(c, common.Read, "path", -1)
	if !ok {
		return
	}

	serializer := common.SimulationModelSerializer{c, m.SimulationModel}
	c.JSON(http.StatusOK, gin.H{
		"model": serializer.Response(),
	})
}

// deleteSimulationModel godoc
// @Summary Delete a simulation model
// @ID deleteSimulationModel
// @Tags models
// @Produce json
// @Success 200 "OK."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param modelID path int true "Model ID"
// @Router /models/{modelID} [delete]
func deleteSimulationModel(c *gin.Context) {

	ok, m := CheckPermissions(c, common.Delete, "path", -1)
	if !ok {
		return
	}

	err := m.delete()
	if common.ProvideErrorResponse(c, err) {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OK.",
	})
}
