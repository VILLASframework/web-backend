package simulationmodel

import (
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/helper"
	"net/http"

	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/scenario"
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

	ok, so := scenario.CheckPermissions(c, database.Read, "query", -1)
	if !ok {
		return
	}

	db := database.GetDB()
	var models []database.SimulationModel
	err := db.Order("ID asc").Model(so).Related(&models, "Models").Error
	if !helper.DBError(c, err) {
		c.JSON(http.StatusOK, gin.H{"models": models})
	}

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
		helper.BadRequestError(c, "Bad request. Error binding form data to JSON: "+err.Error())
		return
	}

	// validate the request
	if err = req.validate(); err != nil {
		helper.UnprocessableEntityError(c, err.Error())
		return
	}

	// Create the new simulation model from the request
	newSimulationModel := req.createSimulationModel()

	// check access to the scenario
	ok, _ := scenario.CheckPermissions(c, database.Update, "body", int(newSimulationModel.ScenarioID))
	if !ok {
		return
	}

	// add the new simulation model to the scenario
	err = newSimulationModel.addToScenario()
	if !helper.DBError(c, err) {
		c.JSON(http.StatusOK, gin.H{"model": newSimulationModel.SimulationModel})
	}

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

	ok, oldSimulationModel := CheckPermissions(c, database.Update, "path", -1)
	if !ok {
		return
	}

	var req updateSimulationModelRequest
	err := c.BindJSON(&req)
	if err != nil {
		helper.BadRequestError(c, "Error binding form data to JSON: "+err.Error())
		return
	}

	// Validate the request
	if err := req.validate(); err != nil {
		helper.BadRequestError(c, err.Error())
		return
	}

	// Create the updatedSimulationModel from oldSimulationModel
	updatedSimulationModel := req.updatedSimulationModel(oldSimulationModel)

	// Finally, update the simulation model
	err = oldSimulationModel.Update(updatedSimulationModel)
	if !helper.DBError(c, err) {
		c.JSON(http.StatusOK, gin.H{"model": updatedSimulationModel.SimulationModel})
	}

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

	ok, m := CheckPermissions(c, database.Read, "path", -1)
	if !ok {
		return
	}

	c.JSON(http.StatusOK, gin.H{"model": m.SimulationModel})
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

	ok, m := CheckPermissions(c, database.Delete, "path", -1)
	if !ok {
		return
	}

	err := m.delete()
	if !helper.DBError(c, err) {
		c.JSON(http.StatusOK, gin.H{"model": m.SimulationModel})
	}
}
