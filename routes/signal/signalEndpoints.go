package signal

import (
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
// @Param direction modelID string true "Model ID of signals to be obtained"
// @Success 200 {array} common.Signal "Requested signals."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
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
		errormsg := "Bad request. Direction has to be in or out"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errormsg,
		})
		return
	}

	db := common.GetDB()
	var sigs []common.Signal
	err := db.Order("ID asc").Model(m).Where("Direction = ?", direction).Related(&sigs, mapping).Error
	if common.ProvideErrorResponse(c, err) {
		return
	}

	serializer := common.SignalsSerializer{c, sigs}
	c.JSON(http.StatusOK, gin.H{
		"signals": serializer.Response(),
	})
}

// AddSignal godoc
// @Summary Add a signal to a signal mapping of a model
// @ID AddSignal
// @Accept json
// @Produce json
// @Tags signals
// @Param inputSignal body common.Signal true "A signal to be added to the model incl. direction and model ID to which signal shall be added"
// @Success 200 "OK."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Router /signals [post]
func addSignal(c *gin.Context) {

	var newSignal Signal
	err := c.BindJSON(&newSignal)
	if err != nil {
		errormsg := "Bad request. Error binding form data to JSON: " + err.Error()
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errormsg,
		})
		return
	}

	ok, _ := simulationmodel.CheckPermissions(c, common.Update, "body", int(newSignal.SimulationModelID))
	if !ok {
		return
	}

	// Add signal to model
	err = newSignal.addToSimulationModel()
	if common.ProvideErrorResponse(c, err) == false {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK.",
		})
	}
}

// updateSignal godoc
// @Summary Update a signal
// @ID updateSignal
// @Tags signals
// @Produce json
// @Success 200 "OK."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param signalID path int true "ID of signal to be updated"
// @Router /signals/{signalID} [put]
func updateSignal(c *gin.Context) {
	ok, sig := checkPermissions(c, common.Delete)
	if !ok {
		return
	}

	var modifiedSignal Signal
	err := c.BindJSON(&modifiedSignal)
	if err != nil {
		errormsg := "Bad request. Error binding form data to JSON: " + err.Error()
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errormsg,
		})
		return
	}

	err = sig.update(modifiedSignal)
	if common.ProvideErrorResponse(c, err) == false {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK.",
		})
	}

}

// getSignal godoc
// @Summary Get a signal
// @ID getSignal
// @Tags signals
// @Produce json
// @Success 200 "OK."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param signalID path int true "ID of signal to be obtained"
// @Router /signals/{signalID} [get]
func getSignal(c *gin.Context) {
	ok, sig := checkPermissions(c, common.Delete)
	if !ok {
		return
	}

	serializer := common.SignalSerializer{c, sig.Signal}
	c.JSON(http.StatusOK, gin.H{
		"signal": serializer.Response(),
	})
}

// deleteSignal godoc
// @Summary Delete a signal
// @ID deleteSignal
// @Tags signals
// @Produce json
// @Success 200 "OK."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param signalID path int true "ID of signal to be deleted"
// @Router /signals/{signalID} [delete]
func deleteSignal(c *gin.Context) {

	ok, sig := checkPermissions(c, common.Delete)
	if !ok {
		return
	}

	err := sig.delete()
	if common.ProvideErrorResponse(c, err) == false {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK.",
		})
	}

}
