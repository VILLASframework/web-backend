package endpoints

import (
	"fmt"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/queries"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/serializers"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

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

