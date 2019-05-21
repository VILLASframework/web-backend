package endpoints

import (
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/simulation"
	"github.com/gin-gonic/gin"
	"net/http"
)


func simulationReadAllEp(c *gin.Context) {
	allSimulations, _, _ := simulation.FindAllSimulations()
	serializer := simulation.SimulationsSerializerNoAssoc{c, allSimulations}
	c.JSON(http.StatusOK, gin.H{
		"simulations": serializer.Response(),
	})
}

func simulationRegistrationEp(c *gin.Context) {
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
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}

func simulationDeleteEp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}
