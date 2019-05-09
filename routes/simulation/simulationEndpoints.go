package simulation

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func SimulationsRegister(r *gin.RouterGroup) {
	r.GET("/", simulationsReadEp)
	r.POST("/", simulationRegistrationEp)
	r.PUT("/:SimulationID", simulationUpdateEp)
	r.GET("/:SimulationID", simulationReadEp)
	r.DELETE("/:SimulationID", simulationDeleteEp)
}

func simulationsReadEp(c *gin.Context) {
	allSimulations, _, _ := FindAllSimulations()
	serializer := SimulationsSerializerNoAssoc{c, allSimulations}
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

