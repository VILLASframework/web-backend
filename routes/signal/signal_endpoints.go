package signal

import (
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/helper"
	"net/http"

	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/database"
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

	ok, m := simulationmodel.CheckPermissions(c, database.Read, "query", -1)
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
		helper.BadRequestError(c, "Bad request. Direction has to be in or out")
		return
	}

	db := database.GetDB()
	var sigs []database.Signal
	err := db.Order("ID asc").Model(m).Where("Direction = ?", direction).Related(&sigs, mapping).Error
	if helper.DBError(c, err) {
		return
	}

	c.JSON(http.StatusOK, gin.H{"signals": sigs})
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
		helper.BadRequestError(c, err.Error())
		return
	}

	// Validate the request
	if err := req.validate(); err != nil {
		helper.UnprocessableEntityError(c, err.Error())
		return
	}

	// Create the new signal from the request
	newSignal := req.createSignal()

	ok, _ := simulationmodel.CheckPermissions(c, database.Update, "body", int(newSignal.SimulationModelID))
	if !ok {
		return
	}

	// Add signal to model
	err := newSignal.addToSimulationModel()
	if err != nil {
		helper.DBError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"signal": newSignal.Signal})
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
	ok, oldSignal := checkPermissions(c, database.Delete)
	if !ok {
		return
	}

	var req updateSignalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.BadRequestError(c, err.Error())
		return
	}

	// Validate the request
	if err := req.validate(); err != nil {
		helper.BadRequestError(c, err.Error())
		return
	}

	// Create the updatedSignal from oldDashboard
	updatedSignal, err := req.updatedSignal(oldSignal)
	if err != nil {
		helper.BadRequestError(c, err.Error())
		return
	}

	// Update the signal in the DB
	err = oldSignal.update(updatedSignal)
	if err != nil {
		helper.DBError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"signal": updatedSignal.Signal})
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
	ok, sig := checkPermissions(c, database.Delete)
	if !ok {
		return
	}

	c.JSON(http.StatusOK, gin.H{"signal": sig.Signal})
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

	ok, sig := checkPermissions(c, database.Delete)
	if !ok {
		return
	}

	err := sig.delete()
	if err != nil {
		helper.DBError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"signal": sig.Signal})
}
