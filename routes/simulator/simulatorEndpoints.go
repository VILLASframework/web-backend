package simulator

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
)

func RegisterSimulatorEndpoints(r *gin.RouterGroup){
	r.GET("/", GetSimulators)
	r.POST("/", AddSimulator)
	r.PUT("/:simulatorID", UpdateSimulator)
	r.GET("/:simulatorID", GetSimulator)
	r.DELETE("/:simulatorID", DeleteSimulator)
	r.POST("/:simulatorID", SendActionToSimulator)
}

// GetSimulators godoc
// @Summary Get all simulators
// @ID GetSimulators
// @Tags simulators
// @Produce json
// @Success 200 {array} common.SimulatorResponse "Simulator parameters requested by user"
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param simulationID path int true "Simulation ID"
// @Param modelID path int true "Model ID"
// @Router /simulators [get]
func GetSimulators(c *gin.Context) {
	allSimulators, _, _ := FindAllSimulators()
	serializer := common.SimulatorsSerializer{c, allSimulators}
	c.JSON(http.StatusOK, gin.H{
		"simulators": serializer.Response(),
	})
}

func AddSimulator(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}

func UpdateSimulator(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}

func GetSimulator(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}

func DeleteSimulator(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}

func SendActionToSimulator(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}


