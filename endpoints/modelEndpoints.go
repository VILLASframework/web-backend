package endpoints

import (
	"fmt"
	"net/http"
	"strconv"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/queries"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/serializers"

	"github.com/gin-gonic/gin"
)

// modelReadAllEp godoc
// @Summary Get all models of simulation
// @ID GetAllModelsOfSimulation
// @Produce  json
// @Tags model
// @Success 200 {array} common.Model "Array of models to which belong to simulation"
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Router /simulations/{simulationID}/models [get]
func modelReadAllEp(c *gin.Context) {

	simID, err := GetSimulationID(c)
	if err != nil {
		return
	}

	allModels, _, err := queries.FindAllModels(simID)
	if common.ProvideErrorResponse(c, err) {
		return
	}

	serializer := serializers.ModelsSerializer{c, allModels}
	c.JSON(http.StatusOK, gin.H{
		"models": serializer.Response(),
	})
}

// modelRegistrationEp godoc
// @Summary Add a model to a simulation
// @ID AddModelToSimulation
// @Tags model
// @Param inputModel body common.Model true "Model to be added"
// @Success 200 "OK."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Router /simulations/{simulationID}/models [post]
func modelRegistrationEp(c *gin.Context) {

	simID, err := GetSimulationID(c)
	if err != nil {
		return
	}

	var m common.Model
	err = c.BindJSON(&m)
	if err != nil {
		errormsg := "Bad request. Error binding form data to JSON: " + err.Error()
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errormsg,
		})
		return
	}

	err = queries.AddModel(simID, &m)
	if common.ProvideErrorResponse(c, err) == false {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK.",
		})
	}
}

func modelCloneEp(c *gin.Context) {

	modelID, err := GetModelID(c)
	if err != nil {
		return
	}

	targetSimID, err := strconv.Atoi(c.PostForm("TargetSim"))
	if err != nil {
		errormsg := fmt.Sprintf("Bad request. No or incorrect format of target sim ID")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errormsg,
		})
		return
	}

	err = queries.CloneModel(targetSimID, modelID)
	if common.ProvideErrorResponse(c, err) == false {
		c.JSON(http.StatusOK, gin.H{
			"message": "Not implemented.",
		})
	}

}

func modelUpdateEp(c *gin.Context) {

	modelID, err := GetModelID(c)
	if err != nil {
		return
	}

	var m common.Model
	err = c.BindJSON(&m)
	if err != nil {
		errormsg := "Bad request. Error binding form data to JSON: " + err.Error()
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errormsg,
		})
		return
	}

	err = queries.UpdateModel(modelID, &m)
	if common.ProvideErrorResponse(c, err) == false {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK.",
		})
	}

}

func modelReadEp(c *gin.Context) {

	modelID, err := GetModelID(c)
	if err != nil {
		return
	}

	m, err := queries.FindModel(modelID)
	if common.ProvideErrorResponse(c, err) {
		return
	}

	serializer := serializers.ModelSerializer{c, m}
	c.JSON(http.StatusOK, gin.H{
		"model": serializer.Response(),
	})
}

func modelDeleteEp(c *gin.Context) {

	simID, err := GetSimulationID(c)
	if err != nil {
		return
	}

	modelID, err := GetModelID(c)
	if err != nil {
		return
	}

	err = queries.DeleteModel(simID, modelID)
	if common.ProvideErrorResponse(c, err) == false {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK.",
		})
	}
}

func GetModelID(c *gin.Context) (int, error) {

	modelID, err := strconv.Atoi(c.Param("ModelID"))

	if err != nil {
		errormsg := fmt.Sprintf("Bad request. No or incorrect format of model ID")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errormsg,
		})
		return -1, err
	} else {
		return modelID, err

	}
}

