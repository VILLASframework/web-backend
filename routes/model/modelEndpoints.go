package model

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/simulation"
)

func RegisterModelEndpoints(r *gin.RouterGroup) {
	r.GET("/", getModels)
	r.POST("/", addModel)
	//r.POST("/:modelID", cloneModel)
	r.PUT("/:modelID", updateModel)
	r.GET("/:modelID", getModel)
	r.DELETE("/:modelID", deleteModel)
	r.GET("/:modelID/signals", getSignals)
	r.POST("/:modelID/signals", addSignal)
	r.DELETE("/:modelID/signals", deleteSignals)
}

// getModels godoc
// @Summary Get all models of simulation
// @ID getModels
// @Produce  json
// @Tags models
// @Success 200 {array} common.ModelResponse "Array of models to which belong to simulation"
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param simulationID query int true "Simulation ID"
// @Router /models [get]
func getModels(c *gin.Context) {

	simID, err := common.GetSimulationID(c)
	if err != nil {
		return
	}

	db := common.GetDB()
	var models []common.Model

	var sim simulation.Simulation
	err = sim.ByID(uint(simID))
	if common.ProvideErrorResponse(c, err) {
		return
	}

	err = db.Order("ID asc").Model(sim).Related(&models, "Models").Error
	if common.ProvideErrorResponse(c, err) {
		return
	}

	serializer := common.ModelsSerializer{c, models}
	c.JSON(http.StatusOK, gin.H{
		"models": serializer.Response(),
	})
}

// addModel godoc
// @Summary Add a model to a simulation
// @ID addModel
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
func addModel(c *gin.Context) {

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

func cloneModel(c *gin.Context) {

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

// updateModel godoc
// @Summary Update a model
// @ID updateModel
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
func updateModel(c *gin.Context) {

	modelID, err := common.GetModelID(c)
	if err != nil {
		return
	}

	var modifiedModel Model
	err = c.BindJSON(&modifiedModel)
	if err != nil {
		errormsg := "Bad request. Error binding form data to JSON: " + err.Error()
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errormsg,
		})
		return
	}

	var m Model
	err = m.ByID(uint(modelID))
	if common.ProvideErrorResponse(c, err) {
		return
	}

	err = m.update(modifiedModel)
	if common.ProvideErrorResponse(c, err) == false {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK.",
		})
	}

}

// getModel godoc
// @Summary Get a model
// @ID getModel
// @Tags models
// @Produce json
// @Success 200 {object} common.ModelResponse "Requested model."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param modelID path int true "Model ID"
// @Router /models/{modelID} [get]
func getModel(c *gin.Context) {

	modelID, err := common.GetModelID(c)
	if err != nil {
		return
	}

	var m Model
	err = m.ByID(uint(modelID))
	if common.ProvideErrorResponse(c, err) {
		return
	}

	serializer := common.ModelSerializer{c, m.Model}
	c.JSON(http.StatusOK, gin.H{
		"model": serializer.Response(),
	})
}

// deleteModel godoc
// @Summary Delete a model
// @ID deleteModel
// @Tags models
// @Produce json
// @Success 200 "OK."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param modelID path int true "Model ID"
// @Router /models/{modelID} [delete]
func deleteModel(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"message": "Not implemented.",
	})
}

// AddSignal godoc
// @Summary Add a signal to a signal mapping of a model
// @ID AddSignal
// @Accept json
// @Produce json
// @Tags models
// @Param inputSignal body common.Signal true "A signal to be added to the model"
// @Param direction query string true "Direction of signal (in or out)"
// @Success 200 "OK."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Router /models/{modelID}/signals [post]
func addSignal(c *gin.Context) {

	modelID, err := common.GetModelID(c)
	if err != nil {
		return
	}

	var m Model
	err = m.ByID(uint(modelID))
	if common.ProvideErrorResponse(c, err) {
		return
	}

	direction := c.Request.URL.Query().Get("direction")
	if !(direction == "out") && !(direction == "in") {
		errormsg := "Bad request. Direction has to be in or out"
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errormsg,
		})
		return
	}

	var sig common.Signal
	err = c.BindJSON(&sig)
	if err != nil {
		errormsg := "Bad request. Error binding form data to JSON: " + err.Error()
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errormsg,
		})
		return
	}

	// Add signal to model
	err = m.addSignal(sig, direction)
	if common.ProvideErrorResponse(c, err) == false {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK.",
		})
	}
}

// getSignals godoc
// @Summary Get all signals of one direction
// @ID getSignals
// @Produce json
// @Tags models
// @Param direction query string true "Direction of signal (in or out)"
// @Success 200 {array} common.Signal "Requested signals."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Router /models/{modelID}/signals [get]
func getSignals(c *gin.Context) {

	modelID, err := common.GetModelID(c)
	if err != nil {
		return
	}

	var m Model
	err = m.ByID(uint(modelID))
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

// deleteSignals godoc
// @Summary Delete all signals of a direction
// @ID deleteSignals
// @Tags models
// @Produce json
// @Success 200 "OK."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param modelID path int true "Model ID"
// @Param direction query string true "Direction of signals to delete (in or out)"
// @Router /models/{modelID}/signals [delete]
func deleteSignals(c *gin.Context) {
	modelID, err := common.GetModelID(c)
	if err != nil {
		return
	}

	var m Model
	err = m.ByID(uint(modelID))
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

	err = m.deleteSignals(direction)
	if common.ProvideErrorResponse(c, err) == false {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK.",
		})
	}

}
