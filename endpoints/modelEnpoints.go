package endpoints

import (
	"fmt"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func modelReadAllEp(c *gin.Context) {
	allModels, _, _ := model.FindAllModels()
	serializer := model.ModelsSerializerNoAssoc{c, allModels}
	c.JSON(http.StatusOK, gin.H{
		"models": serializer.Response(),
	})
}

func modelRegistrationEp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}

func modelUpdateEp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}

func modelReadEp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}

func modelDeleteEp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
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

