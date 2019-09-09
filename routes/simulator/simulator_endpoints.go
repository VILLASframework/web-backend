package simulator

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
)

func RegisterSimulatorEndpoints(r *gin.RouterGroup) {
	r.GET("", getSimulators)
	r.POST("", addSimulator)
	r.PUT("/:simulatorID", updateSimulator)
	r.GET("/:simulatorID", getSimulator)
	r.DELETE("/:simulatorID", deleteSimulator)
	r.GET("/:simulatorID/models", getModelsOfSimulator)
	// register action endpoint only if AMQP client is used
	if common.WITH_AMQP == true {
		r.POST("/:simulatorID/action", sendActionToSimulator)
	}
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
// @Router /simulators [get]
func getSimulators(c *gin.Context) {

	ok, _ := checkPermissions(c, common.ModelSimulator, common.Read, false)
	if !ok {
		return
	}

	db := common.GetDB()
	var simulators []common.Simulator
	err := db.Order("ID asc").Find(&simulators).Error
	if common.DBError(c, err) {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"simulators": simulators,
	})
}

// addSimulator godoc
// @Summary Add a simulator
// @ID addSimulator
// @Accept json
// @Produce json
// @Tags simulators
// @Param inputSimulator body simulator.validNewSimulator true "Simulator to be added"
// @Success 200 {object} docs.ResponseSimulator "Simulator that was added"
// @Failure 400 {object} docs.ResponseError "Bad request"
// @Failure 404 {object} docs.ResponseError "Not found"
// @Failure 422 {object} docs.ResponseError "Unprocessable entity"
// @Failure 500 {object} docs.ResponseError "Internal server error"
// @Router /simulators [post]
func addSimulator(c *gin.Context) {

	ok, _ := checkPermissions(c, common.ModelSimulator, common.Create, false)
	if !ok {
		return
	}

	var req addSimulatorRequest
	err := c.BindJSON(&req)
	if err != nil {
		errormsg := "Bad request. Error binding form data to JSON: " + err.Error()
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": errormsg,
		})
		return
	}

	// Validate the request
	if err = req.validate(); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"success": false,
			"message": fmt.Sprintf("%v", err),
		})
		return
	}

	// Create the new simulator from the request
	newSimulator := req.createSimulator()

	// Save new simulator to DB
	err = newSimulator.save()
	if err != nil {
		common.DBError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"simulator": newSimulator.Simulator,
	})
}

// updateSimulator godoc
// @Summary Update a simulator
// @ID updateSimulator
// @Tags simulators
// @Accept json
// @Produce json
// @Param inputSimulator body simulator.validUpdatedSimulator true "Simulator to be updated"
// @Success 200 {object} docs.ResponseSimulator "Simulator that was updated"
// @Failure 400 {object} docs.ResponseError "Bad request"
// @Failure 404 {object} docs.ResponseError "Not found"
// @Failure 422 {object} docs.ResponseError "Unprocessable entity"
// @Failure 500 {object} docs.ResponseError "Internal server error"
// @Param simulatorID path int true "Simulator ID"
// @Router /simulators/{simulatorID} [put]
func updateSimulator(c *gin.Context) {

	ok, oldSimulator := checkPermissions(c, common.ModelSimulator, common.Update, true)
	if !ok {
		return
	}

	var req updateSimulatorRequest
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Bad request. Error binding form data to JSON: " + err.Error(),
		})
		return
	}

	// Validate the request
	if err = req.validate(); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"success": false,
			"message": fmt.Sprintf("%v", err),
		})
		return
	}

	// Create the updatedSimulator from oldSimulator
	updatedSimulator, err := req.updatedSimulator(oldSimulator)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": fmt.Sprintf("%v", err),
		})
		return
	}

	// Finally update the simulator in the DB
	err = oldSimulator.update(updatedSimulator)
	if err != nil {
		common.DBError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"simulator": updatedSimulator.Simulator,
	})

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
// @Param simulatorID path int true "Simulator ID"
// @Router /simulators/{simulatorID} [get]
func getSimulator(c *gin.Context) {

	ok, s := checkPermissions(c, common.ModelSimulator, common.Read, true)
	if !ok {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"simulator": s.Simulator,
	})
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
// @Param simulatorID path int true "Simulator ID"
// @Router /simulators/{simulatorID} [delete]
func deleteSimulator(c *gin.Context) {

	ok, s := checkPermissions(c, common.ModelSimulator, common.Delete, true)
	if !ok {
		return
	}

	// Delete the simulator
	err := s.delete()
	if common.DBError(c, err) {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"simulator": s.Simulator,
	})
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
// @Param simulatorID path int true "Simulator ID"
// @Router /simulators/{simulatorID}/models [get]
func getModelsOfSimulator(c *gin.Context) {

	ok, s := checkPermissions(c, common.ModelSimulator, common.Read, true)
	if !ok {
		return
	}

	// get all associated simulation models
	allModels, _, err := s.getModels()
	if common.DBError(c, err) {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"models": allModels,
	})
}

// sendActionToSimulator godoc
// @Summary Send an action to simulator (only available if backend server is started with -amqp parameter)
// @ID sendActionToSimulator
// @Tags simulators
// @Produce json
// @Param inputAction query string true "Action for simulator"
// @Success 200 {object} docs.ResponseError "Action sent successfully"
// @Failure 400 {object} docs.ResponseError "Bad request"
// @Failure 404 {object} docs.ResponseError "Not found"
// @Failure 422 {object} docs.ResponseError "Unprocessable entity"
// @Failure 500 {object} docs.ResponseError "Internal server error"
// @Param simulatorID path int true "Simulator ID"
// @Router /simulators/{simulatorID}/action [post]
func sendActionToSimulator(c *gin.Context) {

	ok, s := checkPermissions(c, common.ModelSimulatorAction, common.Update, true)
	if !ok {
		return
	}

	var actions []common.Action
	err := c.BindJSON(&actions)
	if err != nil {
		errormsg := "Bad request. Error binding form data to JSON: " + err.Error()
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errormsg,
		})
		return
	}

	now := time.Now()

	for _, action := range actions {
		if action.When == 0 {
			action.When = float32(now.Unix())
		}

		err = common.SendActionAMQP(action, s.UUID)
		if err != nil {
			errormsg := "Internal Server Error. Unable to send actions to simulator: " + err.Error()
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": errormsg,
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "OK.",
	})
}
