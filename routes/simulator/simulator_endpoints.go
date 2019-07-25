package simulator

import (
	"fmt"
	"net/http"
	"strconv"
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
// @Success 200 {array} common.ResponseMsgSimulators "Simulator parameters requested by user"
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Router /simulators [get]
func getSimulators(c *gin.Context) {

	err := common.ValidateRole(c, common.ModelSimulator, common.Read)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Access denied (role validation failed).")
		return
	}

	db := common.GetDB()
	var simulators []common.Simulator
	err = db.Order("ID asc").Find(&simulators).Error
	if common.ProvideErrorResponse(c, err) {
		return
	}
	serializer := common.SimulatorsSerializer{c, simulators}
	c.JSON(http.StatusOK, gin.H{
		"simulators": serializer.Response(),
	})
}

// addSimulator godoc
// @Summary Add a simulator
// @ID addSimulator
// @Accept json
// @Produce json
// @Tags simulators
// @Param inputSimulator body common.ResponseMsgSimulator true "Simulator to be added"
// @Success 200 "OK."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Router /simulators [post]
func addSimulator(c *gin.Context) {

	err := common.ValidateRole(c, common.ModelSimulator, common.Create)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Access denied (role validation failed).")
		return
	}

	var newSimulatorData common.ResponseMsgSimulator
	err = c.BindJSON(&newSimulatorData)
	if err != nil {
		errormsg := "Bad request. Error binding form data to JSON: " + err.Error()
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errormsg,
		})
		return
	}

	var newSimulator Simulator
	newSimulator.ID = newSimulatorData.Simulator.ID
	newSimulator.State = newSimulatorData.Simulator.State
	newSimulator.StateUpdateAt = newSimulatorData.Simulator.StateUpdateAt
	newSimulator.Modeltype = newSimulatorData.Simulator.Modeltype
	newSimulator.UUID = newSimulatorData.Simulator.UUID
	newSimulator.Uptime = newSimulatorData.Simulator.Uptime
	newSimulator.Host = newSimulatorData.Simulator.Host
	newSimulator.RawProperties = newSimulatorData.Simulator.RawProperties
	newSimulator.Properties = newSimulatorData.Simulator.Properties

	err = newSimulator.save()
	if common.ProvideErrorResponse(c, err) == false {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK.",
		})
	}
}

// updateSimulator godoc
// @Summary Update a simulator
// @ID updateSimulator
// @Tags simulators
// @Accept json
// @Produce json
// @Param inputSimulator body common.ResponseMsgSimulator true "Simulator to be updated"
// @Success 200 "OK."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param simulatorID path int true "Simulator ID"
// @Router /simulators/{simulatorID} [put]
func updateSimulator(c *gin.Context) {
	err := common.ValidateRole(c, common.ModelSimulator, common.Update)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Access denied (role validation failed).")
		return
	}

	var modifiedSimulator common.ResponseMsgSimulator

	err = c.BindJSON(&modifiedSimulator)
	if err != nil {
		errormsg := "Bad request. Error unmarshalling data to JSON: " + err.Error()
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errormsg,
		})
		return
	}

	simulatorID, err := strconv.Atoi(c.Param("simulatorID"))
	if err != nil {
		errormsg := fmt.Sprintf("Bad request. No or incorrect format of simulatorID path parameter")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errormsg,
		})
		return
	}

	var s Simulator
	err = s.ByID(uint(simulatorID))
	if common.ProvideErrorResponse(c, err) {
		return
	}

	err = s.update(modifiedSimulator.Simulator)
	if common.ProvideErrorResponse(c, err) == false {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK.",
		})
	}
}

// getSimulator godoc
// @Summary Get simulator
// @ID getSimulator
// @Produce  json
// @Tags simulators
// @Success 200 {object} common.ResponseMsgSimulator "Simulator requested by user"
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param simulatorID path int true "Simulator ID"
// @Router /simulators/{simulatorID} [get]
func getSimulator(c *gin.Context) {

	err := common.ValidateRole(c, common.ModelSimulator, common.Read)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Access denied (role validation failed).")
		return
	}

	simulatorID, err := strconv.Atoi(c.Param("simulatorID"))
	if err != nil {
		errormsg := fmt.Sprintf("Bad request. No or incorrect format of simulatorID path parameter")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errormsg,
		})
		return
	}

	var s Simulator
	err = s.ByID(uint(simulatorID))
	if common.ProvideErrorResponse(c, err) {
		return
	}

	serializer := common.SimulatorSerializer{c, s.Simulator}
	c.JSON(http.StatusOK, gin.H{
		"simulator": serializer.Response(),
	})
}

// deleteSimulator godoc
// @Summary Delete a simulator
// @ID deleteSimulator
// @Tags simulators
// @Produce json
// @Success 200 "OK."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param simulatorID path int true "Simulator ID"
// @Router /simulators/{simulatorID} [delete]
func deleteSimulator(c *gin.Context) {

	err := common.ValidateRole(c, common.ModelSimulator, common.Delete)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Access denied (role validation failed).")
		return
	}

	simulatorID, err := strconv.Atoi(c.Param("simulatorID"))
	if err != nil {
		errormsg := fmt.Sprintf("Bad request. No or incorrect format of simulatorID path parameter")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errormsg,
		})
		return
	}

	var s Simulator
	err = s.ByID(uint(simulatorID))
	if common.ProvideErrorResponse(c, err) {
		return
	}

	err = s.delete()
	if common.ProvideErrorResponse(c, err) {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OK.",
	})
}

// getModelsOfSimulator godoc
// @Summary Get all simulation models in which the simulator is used
// @ID getModelsOfSimulator
// @Tags simulators
// @Produce json
// @Success 200 {object} common.ResponseMsgSimulationModels "Simulation models requested by user"
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param simulatorID path int true "Simulator ID"
// @Router /simulators/{simulatorID}/models [get]
func getModelsOfSimulator(c *gin.Context) {

	err := common.ValidateRole(c, common.ModelSimulator, common.Read)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Access denied (role validation failed).")
		return
	}

	simulatorID, err := strconv.Atoi(c.Param("simulatorID"))
	if err != nil {
		errormsg := fmt.Sprintf("Bad request. No or incorrect format of simulatorID path parameter")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errormsg,
		})
		return
	}

	var s Simulator
	err = s.ByID(uint(simulatorID))
	if common.ProvideErrorResponse(c, err) {
		return
	}

	// get all associated simulation models
	allModels, _, err := s.getModels()
	if common.ProvideErrorResponse(c, err) {
		return
	}

	serializer := common.SimulationModelsSerializer{c, allModels}
	c.JSON(http.StatusOK, gin.H{
		"models": serializer.Response(),
	})

}

// sendActionToSimulator godoc
// @Summary Send an action to simulator (only available if backend server is started with -amqp parameter)
// @ID sendActionToSimulator
// @Tags simulators
// @Produce json
// @Param inputAction query string true "Action for simulator"
// @Success 200 "OK."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param simulatorID path int true "Simulator ID"
// @Router /simulators/{simulatorID}/action [post]
func sendActionToSimulator(c *gin.Context) {

	err := common.ValidateRole(c, common.ModelSimulatorAction, common.Update)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Access denied (role validation failed).")
		return
	}

	simulatorID, err := strconv.Atoi(c.Param("simulatorID"))
	if err != nil {
		errormsg := fmt.Sprintf("Bad request. No or incorrect format of simulatorID path parameter")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errormsg,
		})
		return
	}

	var actions []common.Action
	err = c.BindJSON(&actions)
	if err != nil {
		errormsg := "Bad request. Error binding form data to JSON: " + err.Error()
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errormsg,
		})
		return
	}

	var s Simulator
	err = s.ByID(uint(simulatorID))
	if common.ProvideErrorResponse(c, err) {
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
		"message": "OK.",
	})
}
