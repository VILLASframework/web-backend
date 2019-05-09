package simulationmodel

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func SimulationModelsRegister(r *gin.RouterGroup) {
	r.GET("/", simulationmodelsReadEp)
	r.POST("/", simulationmodelRegistrationEp)
	r.PUT("/:SimulationModelID", simulationmodelUpdateEp)
	r.GET("/:SimulationModelID", simulationmodelReadEp)
	r.DELETE("/:SimulationModelID", simulationmodelDeleteEp)
}

func simulationmodelsReadEp(c *gin.Context) {
	allSimulationModels, _, _ := FindAllSimulationModels()
	serializer := SimulationModelsSerializerNoAssoc{c, allSimulationModels}
	c.JSON(http.StatusOK, gin.H{
		"simulationmodels": serializer.Response(),
	})
}

func simulationmodelRegistrationEp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}

func simulationmodelUpdateEp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}

func simulationmodelReadEp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}

func simulationmodelDeleteEp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}
