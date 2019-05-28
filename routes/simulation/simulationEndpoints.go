package simulation

import (
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/user"
	"net/http"

	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
)

func RegisterSimulationEndpoints(r *gin.RouterGroup){
	r.GET("/", GetSimulations)
	r.POST("/", AddSimulation)
	//r.POST("/:simulationID", CloneSimulation)
	r.PUT("/:simulationID", UpdateSimulation)
	r.GET("/:simulationID", GetSimulation)
	r.DELETE("/:simulationID", DeleteSimulation)

	r.GET("/:simulationID/users", GetUsersOfSimulation)
	r.PUT("/:simulationID/users/:username", AddUserToSimulation)
	r.DELETE("/:simulationID/users/:username", DeleteUserFromSimulation)
}

// GetSimulations godoc
// @Summary Get all simulations
// @ID GetSimulations
// @Produce  json
// @Tags simulations
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

// AddSimulation godoc
// @Summary Add a simulation
// @ID AddSimulation
// @Accept json
// @Produce json
// @Tags simulations
// @Param inputModel body common.ModelResponse true "Simulation to be added"
// @Success 200 "OK."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Router /simulations [post]
func AddSimulation(c *gin.Context) {


}

func CloneSimulation(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}

// UpdateSimulation godoc
// @Summary Update a simulation
// @ID UpdateSimulation
// @Tags simulations
// @Accept json
// @Produce json
// @Param inputSimulation body common.SimulationResponse true "Simulation to be updated"
// @Success 200 "OK."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param simulationID path int true "Simulation ID"
// @Router /simulations/{simulationID} [put]
func UpdateSimulation(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}

// GetSimulation godoc
// @Summary Get simulation
// @ID GetSimulation
// @Produce  json
// @Tags simulations
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

// DeleteSimulation godoc
// @Summary Delete a simulation
// @ID DeleteSimulation
// @Tags simulations
// @Produce json
// @Success 200 "OK."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param simulationID path int true "Simulation ID"
// @Router /simulations/{simulationID} [delete]
func DeleteSimulation(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}


// GetUsersOfSimulation godoc
// @Summary Get users of simulation
// @ID GetUsersOfSimulation
// @Produce  json
// @Tags simulations
// @Success 200 {array} common.UserResponse "Array of users that have access to the simulation"
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param simulationID path int true "Simulation ID"
// @Router /simulations/{simulationID}/users/ [get]
func GetUsersOfSimulation(c *gin.Context) {

	simID, err := common.GetSimulationID(c)
	if err != nil {
		return
	}

	sim, err := FindSimulation(simID)
	if common.ProvideErrorResponse(c, err) {
		return
	}

	// Find all users of simulation
	allUsers, _, err := user.FindAllUsersSim(&sim)
	if common.ProvideErrorResponse(c, err) {
		return
	}

	serializer := common.UsersSerializer{c, allUsers}
	c.JSON(http.StatusOK, gin.H{
		"users": serializer.Response(),
	})
}


// AddUserToSimulation godoc
// @Summary Add a user to a a simulation
// @ID AddUserToSimulation
// @Tags simulations
// @Produce json
// @Success 200 "OK."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param simulationID path int true "Simulation ID"
// @Param username path int true "User name"
// @Router /simulations/{simulationID}/users/{username} [put]
func AddUserToSimulation(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})

	simID, err := common.GetSimulationID(c)
	if err != nil {
		return
	}

	sim, err := FindSimulation(simID)
	if common.ProvideErrorResponse(c, err) {
		return
	}

	username := c.Param("username")

	u, err := user.FindUserByName(username)
	if common.ProvideErrorResponse(c, err) {
		return
	}

	err = user.AddUserToSim(&sim, &u)
	if common.ProvideErrorResponse(c, err){
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OK.",
	})
}

// DeleteUserFromSimulation godoc
// @Summary Delete a user from asimulation
// @ID DeleteUserFromSimulation
// @Tags simulations
// @Produce json
// @Success 200 "OK."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param simulationID path int true "Simulation ID"
// @Param username path int true "User ID"
// @Router /simulations/{simulationID}/users/{username} [delete]
func DeleteUserFromSimulation(c *gin.Context) {
	simID, err := common.GetSimulationID(c)
	if err != nil {
		return
	}

	sim, err := FindSimulation(simID)
	if common.ProvideErrorResponse(c, err) {
		return
	}

	username := c.Param("username")

	err = user.RemoveUserFromSim(&sim, username)
	if common.ProvideErrorResponse(c, err) {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OK.",
	})
}


