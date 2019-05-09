package simulator

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func SimulatorsRegister(r *gin.RouterGroup) {
	r.GET("/", simulatorsReadEp)
	r.POST("/", simulatorRegistrationEp)
	r.PUT("/:SimulatorID", simulatorUpdateEp)
	r.GET("/:SimulatorID", simulatorReadEp)
	r.DELETE("/:SimulatorID", simulatorDeleteEp)
	r.POST("/:SimulatorID", simulatorSendActionEp)
}

func simulatorsReadEp(c *gin.Context) {
	allSimulators, _, _ := FindAllSimulators()
	serializer := SimulatorsSerializer{c, allSimulators}
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

func simulatorReadEp(c *gin.Context) {
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