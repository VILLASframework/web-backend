package simulationmodel

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/simulation"
)

func RegisterSimulationModelEndpoints(r *gin.RouterGroup) {
	r.GET("", getSimulationModels)
	r.POST("", addSimulationModel)
	//r.POST("/:modelID", cloneSimulationModel)
	r.PUT("/:modelID", updateSimulationModel)
	r.GET("/:modelID", getSimulationModel)
	r.DELETE("/:modelID", deleteSimulationModel)
	r.GET("/:modelID/signals", getSignals)
	r.PUT("/:modelID/signals", addSignal)
	r.DELETE("/:modelID/signals", deleteSignals)
}

// getSimulationModels godoc
// @Summary Get all simulation models of simulation
// @ID getSimulationModels
// @Produce  json
// @Tags models
// @Success 200 {array} common.SimulationModelResponse "Array of models to which belong to simulation"
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param simulationID query int true "Simulation ID"
// @Router /models [get]
func getSimulationModels(c *gin.Context) {

	ok, sim := simulation.CheckPermissions(c, common.Read, "query", -1)
	if !ok {
		return
	}

	db := common.GetDB()
	var models []common.SimulationModel
	err := db.Order("ID asc").Model(sim).Related(&models, "Models").Error
	if common.ProvideErrorResponse(c, err) {
		return
	}

	serializer := common.SimulationModelsSerializer{c, models}
	c.JSON(http.StatusOK, gin.H{
		"models": serializer.Response(),
	})
}

// addSimulationModel godoc
// @Summary Add a simulation model to a simulation
// @ID addSimulationModel
// @Accept json
// @Produce json
// @Tags models
// @Param inputSimulationModel body common.SimulationModelResponse true "Simulation model to be added incl. IDs of simulation and simulator"
// @Success 200 "OK."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Router /models [post]
func addSimulationModel(c *gin.Context) {

	var newModel SimulationModel
	err := c.BindJSON(&newModel)
	if err != nil {
		errormsg := "Bad request. Error binding form data to JSON: " + err.Error()
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errormsg,
		})
		return
	}

	ok, _ := simulation.CheckPermissions(c, common.Create, "body", int(newModel.SimulationID))
	if !ok {
		return
	}

	err = newModel.addToSimulation()
	if common.ProvideErrorResponse(c, err) == false {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK.",
		})
	}
}

func cloneSimulationModel(c *gin.Context) {

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

// updateSimulationModel godoc
// @Summary Update a simulation model
// @ID updateSimulationModel
// @Tags models
// @Accept json
// @Produce json
// @Param inputSimulationModel body common.SimulationModelResponse true "Simulation model to be updated"
// @Success 200 "OK."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param modelID path int true "Model ID"
// @Router /models/{modelID} [put]
func updateSimulationModel(c *gin.Context) {

	ok, m := checkPermissions(c, common.Update)
	if !ok {
		return
	}

	var modifiedModel SimulationModel
	err := c.BindJSON(&modifiedModel)
	if err != nil {
		errormsg := "Bad request. Error binding form data to JSON: " + err.Error()
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errormsg,
		})
		return
	}

	err = m.update(modifiedModel)
	if common.ProvideErrorResponse(c, err) == false {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK.",
		})
	}
}

// getSimulationModel godoc
// @Summary Get a simulation model
// @ID getSimulationModel
// @Tags models
// @Produce json
// @Success 200 {object} common.SimulationModelResponse "Requested simulation model."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param modelID path int true "Model ID"
// @Router /models/{modelID} [get]
func getSimulationModel(c *gin.Context) {

	ok, m := checkPermissions(c, common.Read)
	if !ok {
		return
	}

	serializer := common.SimulationModelSerializer{c, m.SimulationModel}
	c.JSON(http.StatusOK, gin.H{
		"model": serializer.Response(),
	})
}

// deleteSimulationModel godoc
// @Summary Delete a simulation model
// @ID deleteSimulationModel
// @Tags models
// @Produce json
// @Success 200 "OK."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param modelID path int true "Model ID"
// @Router /models/{modelID} [delete]
func deleteSimulationModel(c *gin.Context) {

	ok, m := checkPermissions(c, common.Delete)
	if !ok {
		return
	}

	err := m.delete()
	if common.ProvideErrorResponse(c, err) {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OK.",
	})
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

	ok, m := checkPermissions(c, common.Read)
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
// @Tags models
// @Param inputSignal body common.Signal true "A signal to be added to the model incl. direction"
// @Success 200 "OK."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Router /models/{modelID}/signals [put]
func addSignal(c *gin.Context) {

	ok, m := checkPermissions(c, common.Update)
	if !ok {
		return
	}

	var sig common.Signal
	err := c.BindJSON(&sig)
	if err != nil {
		errormsg := "Bad request. Error binding form data to JSON: " + err.Error()
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errormsg,
		})
		return
	}

	// Add signal to model
	err = m.addSignal(sig)
	if common.ProvideErrorResponse(c, err) == false {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK.",
		})
	}
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

	ok, m := checkPermissions(c, common.Update)
	if !ok {
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

	err := m.deleteSignals(direction)
	if common.ProvideErrorResponse(c, err) == false {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK.",
		})
	}

}
