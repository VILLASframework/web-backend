package model

import (
	"fmt"
	"net/http"
	"strconv"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/simulator"

	"github.com/gin-gonic/gin"
)

func RegisterModelEndpoints(r *gin.RouterGroup){
	r.GET("/", GetModels)
	r.POST("/", AddModel)
	//r.POST("/:modelID", CloneModel)
	r.PUT("/:modelID", UpdateModel)
	r.GET("/:modelID", GetModel)
	r.DELETE("/:modelID", DeleteModel)
	//r.PUT("/:modelID/simulator", UpdateSimulator)
	//r.GET("/:modelID/simulator", GetSimulator)
	//r.POST("/:modelID/signals/:direction", UpdateSignals)
	//r.GET("/:modelID/signals/:direction", GetSignals)
}

// GetModels godoc
// @Summary Get all models of simulation
// @ID GetModels
// @Produce  json
// @Tags models
// @Success 200 {array} common.ModelResponse "Array of models to which belong to simulation"
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param simulationID query int true "Simulation ID"
// @Router /models [get]
func GetModels(c *gin.Context) {

	simID, err := common.GetSimulationID(c)
	if err != nil {
		return
	}

	allModels, _, err := FindAllModels(simID)
	if common.ProvideErrorResponse(c, err) {
		return
	}

	serializer := common.ModelsSerializer{c, allModels}
	c.JSON(http.StatusOK, gin.H{
		"models": serializer.Response(),
	})
}

// AddModel godoc
// @Summary Add a model to a simulation
// @ID AddModel
// @Accept json
// @Produce json
// @Tags models
// @Param inputModel body common.ModelResponse true "Model to be added incl. IDs of simulation and simulator"
// @Success 200 "OK."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Router /models [post]
func AddModel(c *gin.Context) {

	simID, err := common.GetSimulationID(c)
	if err != nil {
		return
	}

	var m Model
	err = c.BindJSON(&m)
	if err != nil {
		errormsg := "Bad request. Error binding form data to JSON: " + err.Error()
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errormsg,
		})
		return
	}

	err = m.addToSimulation(simID)
	if common.ProvideErrorResponse(c, err) == false {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK.",
		})
	}
}

func CloneModel(c *gin.Context) {

	// modelID, err := routes.GetModelID(c)
	// if err != nil {
	// 	return
	// }
	//
	// targetSimID, err := strconv.Atoi(c.PostForm("TargetSim"))
	// if err != nil {
	// 	errormsg := fmt.Sprintf("Bad request. No or incorrect format of target sim ID")
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"error": errormsg,
	// 	})
	// 	return
	// }

	// TODO TO BE IMPLEMENTED
	// Check if target sim exists
	// Check if model exists

	// Get all Signals of Model
	// Get Simulator of Model
	// Get Files of model

	// Add new model object to DB and associate with target sim
	// Add new signal objects to DB and associate with new model object (careful with directions)
	// Associate Simulator with new Model object


	c.JSON(http.StatusOK, gin.H{
		"message": "Not implemented.",
	})


}

// UpdateModel godoc
// @Summary Update a model
// @ID UpdateModel
// @Tags models
// @Accept json
// @Produce json
// @Param inputModel body common.ModelResponse true "Model to be updated"
// @Success 200 "OK."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param modelID path int true "Model ID"
// @Router /models/{modelID} [put]
func UpdateModel(c *gin.Context) {

	modelID, err := common.GetModelID(c)
	if err != nil {
		return
	}

	var m Model
	err = c.BindJSON(&m)
	if err != nil {
		errormsg := "Bad request. Error binding form data to JSON: " + err.Error()
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errormsg,
		})
		return
	}

	err = m.UpdateModel(modelID)
	if common.ProvideErrorResponse(c, err) == false {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK.",
		})
	}

}

// GetModel godoc
// @Summary Get a model
// @ID GetModel
// @Tags models
// @Produce json
// @Success 200 {object} common.ModelResponse "Requested model."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param modelID path int true "Model ID"
// @Router /models/{modelID} [get]
func GetModel(c *gin.Context) {

	modelID, err := common.GetModelID(c)
	if err != nil {
		return
	}

	m, err := FindModel(modelID)
	if common.ProvideErrorResponse(c, err) {
		return
	}

	serializer := common.ModelSerializer{c, m}
	c.JSON(http.StatusOK, gin.H{
		"model": serializer.Response(),
	})
}

// DeleteModel godoc
// @Summary Delete a model
// @ID DeleteModel
// @Tags models
// @Produce json
// @Success 200 "OK."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param modelID path int true "Model ID"
// @Router /models/{modelID} [delete]
func DeleteModel(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"message": "Not implemented.",
	})
}


func GetSimulator(c *gin.Context) {

	modelID, err := common.GetModelID(c)
	if err != nil {
		return
	}

	m, err := FindModel(modelID)
	if common.ProvideErrorResponse(c, err) {
		return
	}

	smtr, err := simulator.FindSimulator(int(m.SimulatorID))
	if common.ProvideErrorResponse(c, err) {
		return
	}

	serializer := common.SimulatorSerializer{c, smtr}
	c.JSON(http.StatusOK, gin.H{
		"simulator": serializer.Response(),
	})
}


func UpdateSimulator(c *gin.Context) {

	// simulator ID as parameter of Query, e.g. simulations/:SimulationID/models/:ModelID/simulator?simulatorID=42
	simulatorID, err := strconv.Atoi(c.Query("simulatorID"))
	if err != nil {
		errormsg := fmt.Sprintf("Bad request. No or incorrect simulator ID")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errormsg,
		})
		return
	}

	modelID, err := common.GetModelID(c)
	if err != nil {
		return
	}

	smtr, err := simulator.FindSimulator(simulatorID)
	if common.ProvideErrorResponse(c, err) {
		return
	}

	_m, err := FindModel(modelID)
	if common.ProvideErrorResponse(c, err) {
		return
	}

	var m = Model{_m}
	err = m.UpdateSimulator(&smtr)
	if common.ProvideErrorResponse(c, err) == false {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK",
		})
	}

}


func UpdateSignals(c *gin.Context) {

	modelID, err := common.GetModelID(c)
	if err != nil {
		return
	}

	_m, err := FindModel(modelID)
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
	var m = Model{_m}
	err = m.UpdateSignals(sigs, direction)
	if common.ProvideErrorResponse(c, err) == false {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK.",
		})
	}
}

func GetSignals(c *gin.Context) {

	modelID, err := common.GetModelID(c)
	if err != nil {
		return
	}

	m, err := FindModel(modelID)
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
		signals = m.InputMapping
	} else {
		signals = m.OutputMapping
	}

	c.JSON(http.StatusOK, gin.H{
		"signals": signals,
	})
}