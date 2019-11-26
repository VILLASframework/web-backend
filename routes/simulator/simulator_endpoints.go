/** Simulator package, endpoints.
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
package simulator

import (
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/helper"
	"github.com/gin-gonic/gin"
	"net/http"

	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
)

func RegisterSimulatorEndpoints(r *gin.RouterGroup) {
	r.GET("", getSimulators)
	r.POST("", addSimulator)
	r.PUT("/:simulatorID", updateSimulator)
	r.GET("/:simulatorID", getSimulator)
	r.DELETE("/:simulatorID", deleteSimulator)
	r.GET("/:simulatorID/models", getModelsOfSimulator)
}

// getSimulators godoc
// @Summary Get all simulators
// @ID getSimulators
// @Tags simulators
// @Produce json
// @Success 200 {object} docs.ResponseSimulators "Simulators requested"
// @Failure 404 {object} docs.ResponseError "Not found"
// @Failure 422 {object} docs.ResponseError "Unprocessable entity"
// @Failure 500 {object} docs.ResponseError "Internal server error"
// @Param Authorization header string true "Authorization token"
// @Router /simulators [get]
func getSimulators(c *gin.Context) {

	// Checking permission is not required here since READ access is independent of user's role

	db := database.GetDB()
	var simulators []database.Simulator
	err := db.Order("ID asc").Find(&simulators).Error
	if !helper.DBError(c, err) {
		c.JSON(http.StatusOK, gin.H{"simulators": simulators})
	}

}

// addSimulator godoc
// @Summary Add a simulator
// @ID addSimulator
// @Accept json
// @Produce json
// @Tags simulators
// @Success 200 {object} docs.ResponseSimulator "Simulator that was added"
// @Failure 400 {object} docs.ResponseError "Bad request"
// @Failure 404 {object} docs.ResponseError "Not found"
// @Failure 422 {object} docs.ResponseError "Unprocessable entity"
// @Failure 500 {object} docs.ResponseError "Internal server error"
// @Param Authorization header string true "Authorization token"
// @Param inputSimulator body simulator.addSimulatorRequest true "Simulator to be added"
// @Router /simulators [post]
func addSimulator(c *gin.Context) {

	ok, _ := CheckPermissions(c, database.ModelSimulator, database.Create, false)
	if !ok {
		return
	}

	var req addSimulatorRequest
	err := c.BindJSON(&req)
	if err != nil {
		helper.BadRequestError(c, "Error binding form data to JSON: "+err.Error())
		return
	}

	// Validate the request
	if err = req.validate(); err != nil {
		helper.UnprocessableEntityError(c, err.Error())
		return
	}

	// Create the new simulator from the request
	newSimulator := req.createSimulator()

	// Save new simulator to DB
	err = newSimulator.save()
	if !helper.DBError(c, err) {
		c.JSON(http.StatusOK, gin.H{"simulator": newSimulator.Simulator})
	}

}

// updateSimulator godoc
// @Summary Update a simulator
// @ID updateSimulator
// @Tags simulators
// @Accept json
// @Produce json
// @Success 200 {object} docs.ResponseSimulator "Simulator that was updated"
// @Failure 400 {object} docs.ResponseError "Bad request"
// @Failure 404 {object} docs.ResponseError "Not found"
// @Failure 422 {object} docs.ResponseError "Unprocessable entity"
// @Failure 500 {object} docs.ResponseError "Internal server error"
// @Param Authorization header string true "Authorization token"
// @Param inputSimulator body simulator.updateSimulatorRequest true "Simulator to be updated"
// @Param simulatorID path int true "Simulator ID"
// @Router /simulators/{simulatorID} [put]
func updateSimulator(c *gin.Context) {

	ok, oldSimulator := CheckPermissions(c, database.ModelSimulator, database.Update, true)
	if !ok {
		return
	}

	var req updateSimulatorRequest
	err := c.BindJSON(&req)
	if err != nil {
		helper.BadRequestError(c, "Error binding form data to JSON: "+err.Error())
		return
	}

	// Validate the request
	if err = req.Simulator.validate(); err != nil {
		helper.UnprocessableEntityError(c, err.Error())
		return
	}

	// Create the updatedSimulator from oldSimulator
	updatedSimulator := req.updatedSimulator(oldSimulator)

	// Finally update the simulator in the DB
	err = oldSimulator.update(updatedSimulator)
	if !helper.DBError(c, err) {
		c.JSON(http.StatusOK, gin.H{"simulator": updatedSimulator.Simulator})
	}

}

// getSimulator godoc
// @Summary Get simulator
// @ID getSimulator
// @Produce  json
// @Tags simulators
// @Success 200 {object} docs.ResponseSimulator "Simulator that was requested"
// @Failure 400 {object} docs.ResponseError "Bad request"
// @Failure 404 {object} docs.ResponseError "Not found"
// @Failure 422 {object} docs.ResponseError "Unprocessable entity"
// @Failure 500 {object} docs.ResponseError "Internal server error"
// @Param Authorization header string true "Authorization token"
// @Param simulatorID path int true "Simulator ID"
// @Router /simulators/{simulatorID} [get]
func getSimulator(c *gin.Context) {

	ok, s := CheckPermissions(c, database.ModelSimulator, database.Read, true)
	if !ok {
		return
	}

	c.JSON(http.StatusOK, gin.H{"simulator": s.Simulator})
}

// deleteSimulator godoc
// @Summary Delete a simulator
// @ID deleteSimulator
// @Tags simulators
// @Produce json
// @Success 200 {object} docs.ResponseSimulator "Simulator that was deleted"
// @Failure 400 {object} docs.ResponseError "Bad request"
// @Failure 404 {object} docs.ResponseError "Not found"
// @Failure 422 {object} docs.ResponseError "Unprocessable entity"
// @Failure 500 {object} docs.ResponseError "Internal server error"
// @Param Authorization header string true "Authorization token"
// @Param simulatorID path int true "Simulator ID"
// @Router /simulators/{simulatorID} [delete]
func deleteSimulator(c *gin.Context) {

	ok, s := CheckPermissions(c, database.ModelSimulator, database.Delete, true)
	if !ok {
		return
	}

	// Delete the simulator
	err := s.delete()
	if !helper.DBError(c, err) {
		c.JSON(http.StatusOK, gin.H{"simulator": s.Simulator})
	}

}

// getModelsOfSimulator godoc
// @Summary Get all simulation models in which the simulator is used
// @ID getModelsOfSimulator
// @Tags simulators
// @Produce json
// @Success 200 {object} docs.ResponseSimulationModels "Simulation models requested by user"
// @Failure 400 {object} docs.ResponseError "Bad request"
// @Failure 404 {object} docs.ResponseError "Not found"
// @Failure 422 {object} docs.ResponseError "Unprocessable entity"
// @Failure 500 {object} docs.ResponseError "Internal server error"
// @Param Authorization header string true "Authorization token"
// @Param simulatorID path int true "Simulator ID"
// @Router /simulators/{simulatorID}/models [get]
func getModelsOfSimulator(c *gin.Context) {

	ok, s := CheckPermissions(c, database.ModelSimulator, database.Read, true)
	if !ok {
		return
	}

	// get all associated simulation models
	allModels, _, err := s.getModels()
	if !helper.DBError(c, err) {
		c.JSON(http.StatusOK, gin.H{"simulationModels": allModels})
	}

}
