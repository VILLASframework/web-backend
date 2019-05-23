package simulation

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
)

func RegisterSimulationEndpoints(r *gin.RouterGroup){
	r.GET("/", GetSimulations)
	r.POST("/", AddSimulation)
	r.POST("/:simulationID", CloneSimulation)
	r.PUT("/:simulationID", UpdateSimulation)
	r.GET("/:simulationID", GetSimulation)
	r.DELETE("/:simulationID", DeleteSimulation)
}

// GetSimulations godoc
// @Summary Get all simulations
// @ID GetSimulations
// @Produce  json
// @Tags simulation
// @Success 200 {array} common.SimulationResponse "Array of simulations to which user has access"
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Router /simulations [get]
func GetSimulations(c *gin.Context) {

	//TODO Identify user who is issuing the request and return only those simulations that are known to the user

	allSimulations, _, _ := FindAllSimulations()
	serializer := common.SimulationsSerializer{c, allSimulations}
	c.JSON(http.StatusOK, gin.H{
		"simulations": serializer.Response(),
	})
}

func AddSimulation(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}

func CloneSimulation(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}

func UpdateSimulation(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}

// GetSimulation godoc
// @Summary Get simulation
// @ID GetSimulation
// @Produce  json
// @Tags simulation
// @Success 200 {object} common.SimulationResponse "Simulation requested by user"
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param simulationID path int true "Simulation ID"
// @Router /simulations/{simulationID} [get]
func GetSimulation(c *gin.Context) {

	simID, err := common.GetSimulationID(c)
	if err != nil {
		return
	}

	sim, err := FindSimulation(simID)
	if common.ProvideErrorResponse(c, err) {
		return
	}

	serializer := common.SimulationSerializer{c, sim}
	c.JSON(http.StatusOK, gin.H{
		"simulation": serializer.Response(),
	})
}

func DeleteSimulation(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}


