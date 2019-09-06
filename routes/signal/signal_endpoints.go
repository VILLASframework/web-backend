package signal

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/simulationmodel"
)

func RegisterSignalEndpoints(r *gin.RouterGroup) {
	r.GET("", getSignals)
	r.POST("", addSignal)
	r.PUT("/:signalID", updateSignal)
	r.GET("/:signalID", getSignal)
	r.DELETE("/:signalID", deleteSignal)
}

// getSignals godoc
// @Summary Get all signals of one direction
// @ID getSignals
// @Produce json
// @Tags signals
// @Param direction query string true "Direction of signal (in or out)"
// @Param modelID query string true "Model ID of signals to be obtained"
// @Success 200 {object} docs.ResponseSignals "Signals which belong to simulation model"
// @Failure 404 {object} docs.ResponseError "Not found"
// @Failure 422 {object} docs.ResponseError "Unprocessable entity"
// @Failure 500 {object} docs.ResponseError "Internal server error"
// @Router /signals [get]
func getSignals(c *gin.Context) {

	ok, m := simulationmodel.CheckPermissions(c, common.Read, "query", -1)
	if !ok {
		return
	}

	var mapping string
	direction := c.Request.URL.Query().Get("direction")
	if direction == "in" {
		mapping = "InputMapping"
	} else if direction == "out" {
		mapping = "OutputMapping"
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Bad request. Direction has to be in or out",
		})
		return
	}

	db := common.GetDB()
	var sigs []common.Signal
	err := db.Order("ID asc").Model(m).Where("Direction = ?", direction).Related(&sigs, mapping).Error
	if common.ProvideErrorResponse(c, err) {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"signals": sigs,
	})
}

// AddSignal godoc
// @Summary Add a signal to a signal mapping of a model
// @ID AddSignal
// @Accept json
// @Produce json
// @Tags signals
// @Param inputSignal body signal.validNewSignal true "A signal to be added to the model incl. direction and model ID to which signal shall be added"
// @Success 200 {object} docs.ResponseSignal "Signal that was added"
// @Failure 400 {object} docs.ResponseError "Bad request"
// @Failure 404 {object} docs.ResponseError "Not found"
// @Failure 422 {object} docs.ResponseError "Unprocessable entity"
// @Failure 500 {object} docs.ResponseError "Internal server error"
// @Router /signals [post]
func addSignal(c *gin.Context) {

	var req addSignalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": fmt.Sprintf("%v", err),
		})
		return
	}

	// Validate the request
	if err := req.validate(); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"success": false,
			"message": fmt.Sprintf("%v", err),
		})
		return
	}

	// Create the new signal from the request
	newSignal := req.createSignal()

	ok, _ := simulationmodel.CheckPermissions(c, common.Update, "body", int(newSignal.SimulationModelID))
	if !ok {
		return
	}

	// Add signal to model
	err := newSignal.addToSimulationModel()
	if err != nil {
		common.ProvideErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"signal": newSignal.Signal,
	})
}

// updateSignal godoc
// @Summary Update a signal
// @ID updateSignal
// @Tags signals
// @Produce json
// @Param inputSignal body signal.validUpdatedSignal true "A signal to be updated"
// @Success 200 {object} docs.ResponseSignal "Signal that was updated"
// @Failure 400 {object} docs.ResponseError "Bad request"
// @Failure 404 {object} docs.ResponseError "Not found"
// @Failure 422 {object} docs.ResponseError "Unprocessable entity"
// @Failure 500 {object} docs.ResponseError "Internal server error"
// @Param signalID path int true "ID of signal to be updated"
// @Router /signals/{signalID} [put]
func updateSignal(c *gin.Context) {
	ok, oldSignal := checkPermissions(c, common.Delete)
	if !ok {
		return
	}

	var req updateSignalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": fmt.Sprintf("%v", err),
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

	// Create the updatedSignal from oldDashboard
	updatedSignal, err := req.updatedSignal(oldSignal)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": fmt.Sprintf("%v", err),
		})
		return
	}

	// Update the signal in the DB
	err = oldSignal.update(updatedSignal)
	if err != nil {
		common.ProvideErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"signal": updatedSignal.Signal,
	})
}

// getSignal godoc
// @Summary Get a signal
// @ID getSignal
// @Tags signals
// @Produce json
// @Success 200 {object} docs.ResponseSignal "Signal that was requested"
// @Failure 400 {object} docs.ResponseError "Bad request"
// @Failure 404 {object} docs.ResponseError "Not found"
// @Failure 422 {object} docs.ResponseError "Unprocessable entity"
// @Failure 500 {object} docs.ResponseError "Internal server error"
// @Param signalID path int true "ID of signal to be obtained"
// @Router /signals/{signalID} [get]
func getSignal(c *gin.Context) {
	ok, sig := checkPermissions(c, common.Delete)
	if !ok {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"signal": sig.Signal,
	})
}

// deleteSignal godoc
// @Summary Delete a signal
// @ID deleteSignal
// @Tags signals
// @Produce json
// @Success 200 {object} docs.ResponseSignal "Signal that was deleted"
// @Failure 400 {object} docs.ResponseError "Bad request"
// @Failure 404 {object} docs.ResponseError "Not found"
// @Failure 422 {object} docs.ResponseError "Unprocessable entity"
// @Failure 500 {object} docs.ResponseError "Internal server error"
// @Param signalID path int true "ID of signal to be deleted"
// @Router /signals/{signalID} [delete]
func deleteSignal(c *gin.Context) {

	ok, sig := checkPermissions(c, common.Delete)
	if !ok {
		return
	}

	err := sig.delete()
	if err != nil {
		common.ProvideErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"signal": sig.Signal,
	})
}
