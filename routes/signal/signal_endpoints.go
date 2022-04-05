/**
* This file is part of VILLASweb-backend-go
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

package signal

import (
	"net/http"

	"git.rwth-aachen.de/acs/public/villas/web-backend-go/helper"

	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
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
// @Param configID query string true "Config ID of signals to be obtained"
// @Success 200 {object} api.ResponseSignals "Signals which belong to component configuration"
// @Failure 404 {object} api.ResponseError "Not found"
// @Failure 422 {object} api.ResponseError "Unprocessable entity"
// @Failure 500 {object} api.ResponseError "Internal server error"
// @Router /signals [get]
// @Security Bearer
func getSignals(c *gin.Context) {

	ok, m := database.CheckComponentConfigPermissions(c, database.Read, "query", -1)
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
	if !helper.DBError(c, err) {
		c.JSON(http.StatusOK, gin.H{"signals": sigs})
	}

}

// AddSignal godoc
// @Summary Add a signal to a signal mapping of a component configuration
// @ID AddSignal
// @Accept json
// @Produce json
// @Tags signals
// @Success 200 {object} api.ResponseSignal "Signal that was added"
// @Failure 400 {object} api.ResponseError "Bad request"
// @Failure 404 {object} api.ResponseError "Not found"
// @Failure 422 {object} api.ResponseError "Unprocessable entity"
// @Failure 500 {object} api.ResponseError "Internal server error"
// @Param inputSignal body signal.addSignalRequest true "A signal to be added to the component configuration incl. direction and config ID to which signal shall be added"
// @Router /signals [post]
// @Security Bearer
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

	ok, _ := database.CheckComponentConfigPermissions(c, database.Update, "body", int(newSignal.ConfigID))
	if !ok {
		return
	}

	// Add signal to component configuration
	err := newSignal.AddToConfig()
	if !helper.DBError(c, err) {
		c.JSON(http.StatusOK, gin.H{"signal": newSignal.Signal})
	}

}

// updateSignal godoc
// @Summary Update a signal
// @ID updateSignal
// @Tags signals
// @Produce json
// @Success 200 {object} api.ResponseSignal "Signal that was updated"
// @Failure 400 {object} api.ResponseError "Bad request"
// @Failure 404 {object} api.ResponseError "Not found"
// @Failure 422 {object} api.ResponseError "Unprocessable entity"
// @Failure 500 {object} api.ResponseError "Internal server error"
// @Param inputSignal body signal.updateSignalRequest true "A signal to be updated"
// @Param signalID path int true "ID of signal to be updated"
// @Router /signals/{signalID} [put]
// @Security Bearer
func updateSignal(c *gin.Context) {
	ok, oldSignal_r := database.CheckSignalPermissions(c, database.Delete)
	if !ok {
		return
	}

	var oldSignal Signal
	oldSignal.Signal = oldSignal_r

	var req updateSignalRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.BadRequestError(c, err.Error())
		return
	}

	// Validate the request
	if err := req.Signal.validate(); err != nil {
		helper.BadRequestError(c, err.Error())
		return
	}

	// Create the updatedSignal from oldDashboard
	updatedSignal := req.updatedSignal(oldSignal)

	// Update the signal in the DB
	err := oldSignal.update(updatedSignal)
	if !helper.DBError(c, err) {
		c.JSON(http.StatusOK, gin.H{"signal": updatedSignal.Signal})
	}

}

// getSignal godoc
// @Summary Get a signal
// @ID getSignal
// @Tags signals
// @Produce json
// @Success 200 {object} api.ResponseSignal "Signal that was requested"
// @Failure 400 {object} api.ResponseError "Bad request"
// @Failure 404 {object} api.ResponseError "Not found"
// @Failure 422 {object} api.ResponseError "Unprocessable entity"
// @Failure 500 {object} api.ResponseError "Internal server error"
// @Param signalID path int true "ID of signal to be obtained"
// @Router /signals/{signalID} [get]
// @Security Bearer
func getSignal(c *gin.Context) {
	ok, sig := database.CheckSignalPermissions(c, database.Delete)
	if !ok {
		return
	}

	c.JSON(http.StatusOK, gin.H{"signal": sig})
}

// deleteSignal godoc
// @Summary Delete a signal
// @ID deleteSignal
// @Tags signals
// @Produce json
// @Success 200 {object} api.ResponseSignal "Signal that was deleted"
// @Failure 400 {object} api.ResponseError "Bad request"
// @Failure 404 {object} api.ResponseError "Not found"
// @Failure 422 {object} api.ResponseError "Unprocessable entity"
// @Failure 500 {object} api.ResponseError "Internal server error"
// @Param signalID path int true "ID of signal to be deleted"
// @Router /signals/{signalID} [delete]
// @Security Bearer
func deleteSignal(c *gin.Context) {

	ok, sig_r := database.CheckSignalPermissions(c, database.Delete)
	if !ok {
		return
	}

	var sig Signal
	sig.Signal = sig_r

	err := sig.delete()
	if !helper.DBError(c, err) {
		c.JSON(http.StatusOK, gin.H{"signal": sig.Signal})
	}

}
