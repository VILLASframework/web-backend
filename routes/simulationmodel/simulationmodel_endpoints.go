/** Simulationmodel package, endpoints.
*
* @author Sonja Happ <sonja.happ@eonerc.rwth-aachen.de>
* @copyright 2014-2019, Institute for Automation of Complex Power Systems, EONERC
* @license GNU General Public License (version 3)
*
* VILLASweb-backend-go
*
* This program is free software: you can redistribute it and/or modify
* it under the terms of the GNU General Public License as published by
* the Free Software Foundation, either version 3 of the License, or
* any later version.
*
* This program is distributed in the hope that it will be useful,
* but WITHOUT ANY WARRANTY; without even the implied warranty of
* MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
* GNU General Public License for more details.
*
* You should have received a copy of the GNU General Public License
* along with this program.  If not, see <http://www.gnu.org/licenses/>.
*********************************************************************************/
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
// @Tags simulationModels
// @Success 200 {object} docs.ResponseSimulationModels "Simulation models which belong to scenario"
// @Failure 404 {object} docs.ResponseError "Not found"
// @Failure 422 {object} docs.ResponseError "Unprocessable entity"
// @Failure 500 {object} docs.ResponseError "Internal server error"
// @Param Authorization header string true "Authorization token"
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
		c.JSON(http.StatusOK, gin.H{"simulationModels": models})
	}

}

// addSimulationModel godoc
// @Summary Add a simulation model to a scenario
// @ID addSimulationModel
// @Accept json
// @Produce json
// @Tags simulationModels
// @Success 200 {object} docs.ResponseSimulationModel "simulation model that was added"
// @Failure 400 {object} docs.ResponseError "Bad request"
// @Failure 404 {object} docs.ResponseError "Not found"
// @Failure 422 {object} docs.ResponseError "Unprocessable entity"
// @Failure 500 {object} docs.ResponseError "Internal server error"
// @Param Authorization header string true "Authorization token"
// @Param inputSimulationModel body simulationmodel.addSimulationModelRequest true "Simulation model to be added incl. IDs of scenario and IC"
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
		c.JSON(http.StatusOK, gin.H{"simulationModel": newSimulationModel.SimulationModel})
	}

}

// updateSimulationModel godoc
// @Summary Update a simulation model
// @ID updateSimulationModel
// @Tags simulationModels
// @Accept json
// @Produce json
// @Success 200 {object} docs.ResponseSimulationModel "simulation model that was added"
// @Failure 400 {object} docs.ResponseError "Bad request"
// @Failure 404 {object} docs.ResponseError "Not found"
// @Failure 422 {object} docs.ResponseError "Unprocessable entity"
// @Failure 500 {object} docs.ResponseError "Internal server error"
// @Param Authorization header string true "Authorization token"
// @Param inputSimulationModel body simulationmodel.updateSimulationModelRequest true "Simulation model to be updated"
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
	if err := req.SimulationModel.validate(); err != nil {
		helper.BadRequestError(c, err.Error())
		return
	}

	// Create the updatedSimulationModel from oldSimulationModel
	updatedSimulationModel := req.updatedSimulationModel(oldSimulationModel)

	// Finally, update the simulation model
	err = oldSimulationModel.Update(updatedSimulationModel)
	if !helper.DBError(c, err) {
		c.JSON(http.StatusOK, gin.H{"simulationModel": updatedSimulationModel.SimulationModel})
	}

}

// getSimulationModel godoc
// @Summary Get a simulation model
// @ID getSimulationModel
// @Tags simulationModels
// @Produce json
// @Success 200 {object} docs.ResponseSimulationModel "simulation model that was requested"
// @Failure 400 {object} docs.ResponseError "Bad request"
// @Failure 404 {object} docs.ResponseError "Not found"
// @Failure 422 {object} docs.ResponseError "Unprocessable entity"
// @Failure 500 {object} docs.ResponseError "Internal server error"
// @Param Authorization header string true "Authorization token"
// @Param modelID path int true "Model ID"
// @Router /models/{modelID} [get]
func getSimulationModel(c *gin.Context) {

	ok, m := CheckPermissions(c, database.Read, "path", -1)
	if !ok {
		return
	}

	c.JSON(http.StatusOK, gin.H{"simulationModel": m.SimulationModel})
}

// deleteSimulationModel godoc
// @Summary Delete a simulation model
// @ID deleteSimulationModel
// @Tags simulationModels
// @Produce json
// @Success 200 {object} docs.ResponseSimulationModel "simulation model that was deleted"
// @Failure 400 {object} docs.ResponseError "Bad request"
// @Failure 404 {object} docs.ResponseError "Not found"
// @Failure 422 {object} docs.ResponseError "Unprocessable entity"
// @Failure 500 {object} docs.ResponseError "Internal server error"
// @Param Authorization header string true "Authorization token"
// @Param modelID path int true "Model ID"
// @Router /models/{modelID} [delete]
func deleteSimulationModel(c *gin.Context) {

	ok, m := CheckPermissions(c, database.Delete, "path", -1)
	if !ok {
		return
	}

	err := m.delete()
	if !helper.DBError(c, err) {
		c.JSON(http.StatusOK, gin.H{"simulationModel": m.SimulationModel})
	}
}
