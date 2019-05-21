package endpoints

import (
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/simulator"
	"github.com/gin-gonic/gin"
	"net/http"
)

func simulatorReadAllEp(c *gin.Context) {
	allSimulators, _, _ := simulator.FindAllSimulators()
	serializer := simulator.SimulatorsSerializer{c, allSimulators}
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
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}

func simulatorReadEp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}

func simulatorReadModelEp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
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


