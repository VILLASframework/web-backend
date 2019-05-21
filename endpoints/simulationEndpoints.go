package endpoints

import (
	"fmt"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/queries"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/serializers"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)


func simulationReadAllEp(c *gin.Context) {

	//TODO Identify user who is issuing the request and return only those simulations that are known to the user

	allSimulations, _, _ := queries.FindAllSimulations()
	serializer := serializers.SimulationsSerializer{c, allSimulations}
	c.JSON(http.StatusOK, gin.H{
		"simulations": serializer.Response(),
	})
}

func simulationRegistrationEp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}

func simulationCloneEp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}

func simulationUpdateEp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}

func simulationReadEp(c *gin.Context) {

	simID, err := GetSimulationID(c)
	if err != nil {
		return
	}

	sim, err := queries.FindSimulation(simID)
	if common.ProvideErrorResponse(c, err) {
		return
	}

	serializer := serializers.SimulationSerializer{c, sim}
	c.JSON(http.StatusOK, gin.H{
		"simulation": serializer.Response(),
	})
}

func simulationDeleteEp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}


func GetSimulationID(c *gin.Context) (int, error) {

	simID, err := strconv.Atoi(c.Param("SimulationID"))

	if err != nil {
		errormsg := fmt.Sprintf("Bad request. No or incorrect format of simulation ID")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errormsg,
		})
		return -1, err
	} else {
		return simID, err

	}
}