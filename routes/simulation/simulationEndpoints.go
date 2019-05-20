package simulation

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func SimulationsRegister(r *gin.RouterGroup) {

	//simulations
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

func GetSimulationID(c *gin.Context) (int, error) {

	simulationID, err := strconv.Atoi(c.Param("SimulationID"))

	if err != nil {
		errormsg := fmt.Sprintf("Bad request. No or incorrect format of simulation ID")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errormsg,
		})
		return -1, err
	} else {
		return simulationID, err

	}
}
