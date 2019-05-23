package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/queries"
)

func signalRegistrationEp(c *gin.Context) {

	modelID, err := GetModelID(c)
	if err != nil {
		return
	}

	model, err := queries.FindModel(modelID)
	if common.ProvideErrorResponse(c, err) {
		return
	}

	direction := c.Param("direction")
	if !(direction == "out") && !(direction == "in") {
		errormsg := "Bad request. Direction has to be in or out"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errormsg,
		})
		return
	}

	var sigs []common.Signal
	err = c.BindJSON(&sigs)
	if err != nil {
		errormsg := "Bad request. Error binding form data to JSON: " + err.Error()
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errormsg,
		})
		return
	}

	// Add signals to model and remove all existing Signals of the requested direction (if any)
	err = queries.ReplaceSignals(&model, sigs, direction)
	if common.ProvideErrorResponse(c, err) == false {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK.",
		})
	}
}

func signalReadAllEp(c *gin.Context) {

	modelID, err := GetModelID(c)
	if err != nil {
		return
	}

	model, err := queries.FindModel(modelID)
	if common.ProvideErrorResponse(c, err) {
		return
	}

	direction := c.Param("direction")
	if !(direction == "out") && !(direction == "in") {
		errormsg := "Bad request. Direction has to be in or out"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errormsg,
		})
		return
	}

	var signals []common.Signal
	if direction == "in" {
		signals = model.InputMapping
	} else {
		signals = model.OutputMapping
	}

	c.JSON(http.StatusOK, gin.H{
		"signals": signals,
	})
}

