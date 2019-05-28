package simulator

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
)

func RegisterSimulatorEndpoints(r *gin.RouterGroup) {
	r.GET("/", GetSimulators)
	r.POST("/", AddSimulator)
	r.PUT("/:simulatorID", UpdateSimulator)
	r.GET("/:simulatorID", GetSimulator)
	r.DELETE("/:simulatorID", DeleteSimulator)
	r.POST("/:simulatorID/action", SendActionToSimulator)
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
// @Router /simulators [get]
func GetSimulators(c *gin.Context) {
	db := common.GetDB()
	var simulators []common.Simulator
	err := db.Order("ID asc").Find(&simulators).Error
	if common.ProvideErrorResponse(c, err) {
		return
	}
	serializer := common.SimulatorsSerializer{c, simulators}
	c.JSON(http.StatusOK, gin.H{
		"simulators": serializer.Response(),
	})
}

// AddSimulator godoc
// @Summary Add a simulator
// @ID AddSimulator
// @Accept json
// @Produce json
// @Tags simulators
// @Param inputSimulator body common.SimulatorResponse true "Simulator to be added"
// @Success 200 "OK."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Router /simulators [post]
func AddSimulator(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}

// UpdateSimulator godoc
// @Summary Update a simulator
// @ID UpdateSimulator
// @Tags simulators
// @Accept json
// @Produce json
// @Param inputSimulator body common.SimulatorResponse true "Simulator to be updated"
// @Success 200 "OK."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param simulatorID path int true "Simulator ID"
// @Router /simulators/{simulatorID} [put]
func UpdateSimulator(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}

// GetSimulator godoc
// @Summary Get simulator
// @ID GetSimulator
// @Produce  json
// @Tags simulators
// @Success 200 {object} common.SimulatorResponse "Simulator requested by user"
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param simulatorID path int true "Simulator ID"
// @Router /simulators/{simulatorID} [get]
func GetSimulator(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}

// DeleteSimulator godoc
// @Summary Delete a simulator
// @ID DeleteSimulator
// @Tags simulators
// @Produce json
// @Success 200 "OK."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param simulatorID path int true "Simulator ID"
// @Router /simulators/{simulatorID} [delete]
func DeleteSimulator(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}

// SendActionToSimulator godoc
// @Summary Send an action to simulator
// @ID SendActionToSimulator
// @Tags simulators
// @Produce json
// @Param inputAction query string true "Action for simulator"
// @Success 200 "OK."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param simulatorID path int true "Simulator ID"
// @Router /simulators/{simulatorID}/action [post]
func SendActionToSimulator(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}
