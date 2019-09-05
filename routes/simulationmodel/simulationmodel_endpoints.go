package simulationmodel

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/scenario"
)

func RegisterSimulationModelEndpoints(r *gin.RouterGroup) {
	r.GET("", getSimulationModels)
	r.POST("", addSimulationModel)
	r.PUT("/:modelID", updateSimulationModel)
	r.GET("/:modelID", getSimulationModel)
	r.DELETE("/:modelID", deleteSimulationModel)
}

// getSimulationModels godoc
// @Summary Get all simulation models of scenario
// @ID getSimulationModels
// @Produce  json
// @Tags models
// @Success 200 {object} docs.ResponseSimulationModels "Simulation models which belong to scenario"
// @Failure 404 {object} docs.ResponseError "Not found"
// @Failure 422 {object} docs.ResponseError "Unprocessable entity"
// @Failure 500 {object} docs.ResponseError "Internal server error"
// @Param scenarioID query int true "Scenario ID"
// @Router /models [get]
func getSimulationModels(c *gin.Context) {

	ok, so := scenario.CheckPermissions(c, common.Read, "query", -1)
	if !ok {
		return
	}

	db := common.GetDB()
	var models []common.SimulationModel
	err := db.Order("ID asc").Model(so).Related(&models, "Models").Error
	if common.ProvideErrorResponse(c, err) {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"models": models,
	})
}

// addSimulationModel godoc
// @Summary Add a simulation model to a scenario
// @ID addSimulationModel
// @Accept json
// @Produce json
// @Tags models
// @Param inputSimulationModel body simulationmodel.validNewSimulationModel true "Simulation model to be added incl. IDs of scenario and simulator"
// @Success 200 {object} docs.ResponseSimulationModel "simulation model that was added"
// @Failure 400 {object} docs.ResponseError "Bad request"
// @Failure 404 {object} docs.ResponseError "Not found"
// @Failure 422 {object} docs.ResponseError "Unprocessable entity"
// @Failure 500 {object} docs.ResponseError "Internal server error"
// @Router /models [post]
func addSimulationModel(c *gin.Context) {

	// Bind the request to JSON
	var req addSimulationModelRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Bad request. Error binding form data to JSON: " + err.Error(),
		})
		return
	}

	// validate the request
	if err = req.validate(); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"success": false,
			"message": fmt.Sprintf("%v", err),
		})
		return
	}

	// Create the new simulation model from the request
	newSimulationModel := req.createSimulationModel()

	// check access to the scenario
	ok, _ := scenario.CheckPermissions(c, common.Create, "body", int(newSimulationModel.ScenarioID))
	if !ok {
		return
	}

	// add the new simulation model to the scenario
	err = newSimulationModel.addToScenario()
	if err != nil {
		common.ProvideErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"model": newSimulationModel.SimulationModel,
	})
}

// updateSimulationModel godoc
// @Summary Update a simulation model
// @ID updateSimulationModel
// @Tags models
// @Accept json
// @Produce json
// @Param inputSimulationModel body simulationmodel.validUpdatedSimulationModel true "Simulation model to be updated"
// @Success 200 {object} docs.ResponseSimulationModel "simulation model that was added"
// @Failure 400 {object} docs.ResponseError "Bad request"
// @Failure 404 {object} docs.ResponseError "Not found"
// @Failure 422 {object} docs.ResponseError "Unprocessable entity"
// @Failure 500 {object} docs.ResponseError "Internal server error"
// @Param modelID path int true "Model ID"
// @Router /models/{modelID} [put]
func updateSimulationModel(c *gin.Context) {

	ok, oldSimulationModel := CheckPermissions(c, common.Update, "path", -1)
	if !ok {
		return
	}

	var req updateSimulationModelRequest
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Bad request. Error binding form data to JSON: " + err.Error(),
		})
		return
	}

	// Validate the request
	if err := req.validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": fmt.Sprintf("%v", err),
		})
		return
	}

	// Create the updatedSimulationModel from oldSimulationModel
	updatedSimulationModel, err := req.updatedSimulationModel(oldSimulationModel)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": fmt.Sprintf("%v", err),
		})
		return
	}

	// Finally, update the simulation model
	err = oldSimulationModel.Update(updatedSimulationModel)
	if err != nil {
		common.ProvideErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"model": updatedSimulationModel.SimulationModel,
	})

}

// getSimulationModel godoc
// @Summary Get a simulation model
// @ID getSimulationModel
// @Tags models
// @Produce json
// @Success 200 {object} docs.ResponseSimulationModel "simulation model that was requested"
// @Failure 400 {object} docs.ResponseError "Bad request"
// @Failure 404 {object} docs.ResponseError "Not found"
// @Failure 422 {object} docs.ResponseError "Unprocessable entity"
// @Failure 500 {object} docs.ResponseError "Internal server error"
// @Param modelID path int true "Model ID"
// @Router /models/{modelID} [get]
func getSimulationModel(c *gin.Context) {

	ok, m := CheckPermissions(c, common.Read, "path", -1)
	if !ok {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"model": m.SimulationModel,
	})
}

// deleteSimulationModel godoc
// @Summary Delete a simulation model
// @ID deleteSimulationModel
// @Tags models
// @Produce json
// @Success 200 {object} docs.ResponseSimulationModel "simulation model that was deleted"
// @Failure 400 {object} docs.ResponseError "Bad request"
// @Failure 404 {object} docs.ResponseError "Not found"
// @Failure 422 {object} docs.ResponseError "Unprocessable entity"
// @Failure 500 {object} docs.ResponseError "Internal server error"
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
		"model": m.SimulationModel,
	})
}
