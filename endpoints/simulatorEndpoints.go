package endpoints

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/queries"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/serializers"
)

func simulatorReadAllEp(c *gin.Context) {
	allSimulators, _, _ := queries.FindAllSimulators()
	serializer := serializers.SimulatorsSerializer{c, allSimulators}
	c.JSON(http.StatusOK, gin.H{
		"simulators": serializer.Response(),
	})
}

func simulatorRegistrationEp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}

func simulatorUpdateEp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}

func simulatorUpdateModelEp(c *gin.Context) {

	// simulator ID as parameter of Query, e.g. simulations/:SimulationID/models/:ModelID/simulator?simulatorID=42
	simulatorID, err := strconv.Atoi(c.Query("simulatorID"))
	if err != nil {
		errormsg := fmt.Sprintf("Bad request. No or incorrect simulator ID")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errormsg,
		})
		return
	}

	modelID, err := GetModelID(c)
	if err != nil {
		return
	}

	simulator, err := queries.FindSimulator(simulatorID)
	if common.ProvideErrorResponse(c, err) {
		return
	}

	model, err := queries.FindModel(modelID)
	if common.ProvideErrorResponse(c, err) {
		return
	}

	err = queries.UpdateSimulatorOfModel(&model, &simulator)
	if common.ProvideErrorResponse(c, err) == false {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK",
		})
	}

}

func simulatorReadEp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}

func simulatorReadModelEp(c *gin.Context) {

	modelID, err := GetModelID(c)
	if err != nil {
		return
	}

	model, err := queries.FindModel(modelID)
	if common.ProvideErrorResponse(c, err) {
		return
	}

	simulator, err := queries.FindSimulator(int(model.SimulatorID))
	if common.ProvideErrorResponse(c, err) {
		return
	}

	serializer := serializers.SimulatorSerializer{c, simulator}
	c.JSON(http.StatusOK, gin.H{
		"simulator": serializer.Response(),
	})
}

func simulatorDeleteEp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}

func simulatorSendActionEp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}


