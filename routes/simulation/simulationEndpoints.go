package simulation

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/user"
)

func RegisterSimulationEndpoints(r *gin.RouterGroup) {
	r.GET("/", getSimulations)
	r.POST("/", addSimulation)
	r.PUT("/:simulationID", updateSimulation)
	r.GET("/:simulationID", getSimulation)
	r.DELETE("/:simulationID", deleteSimulation)
	r.GET("/:simulationID/users", getUsersOfSimulation)
	r.PUT("/:simulationID/user", addUserToSimulation)
	r.DELETE("/:simulationID/user", deleteUserFromSimulation)
}

// getSimulations godoc
// @Summary Get all simulations
// @ID getSimulations
// @Produce  json
// @Tags simulations
// @Success 200 {array} common.SimulationResponse "Array of simulations to which user has access"
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Router /simulations [get]
func getSimulations(c *gin.Context) {

	err := common.ValidateRole(c, common.ModelSimulation, common.Read)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, fmt.Sprintf("%v", err))
		return
	}

	// ATTENTION: do not use c.GetInt (common.UserIDCtx) since user_id is of type uint and not int
	userID, _ := c.Get(common.UserIDCtx)
	userRole, _ := c.Get(common.UserRoleCtx)

	var u user.User
	err = u.ByID(userID.(uint))
	if common.ProvideErrorResponse(c, err) {
		return
	}

	// get all simulations for the user who issues the request
	db := common.GetDB()
	var simulations []common.Simulation
	if userRole == "Admin" { // Admin can see all simulations
		err = db.Order("ID asc").Find(&simulations).Error
		if common.ProvideErrorResponse(c, err) {
			return
		}

	} else { // User or Guest roles see only their simulations
		err = db.Order("ID asc").Model(&u).Related(&simulations, "Simulations").Error
		if common.ProvideErrorResponse(c, err) {
			return
		}
	}

	serializer := common.SimulationsSerializer{c, simulations}
	c.JSON(http.StatusOK, gin.H{
		"simulations": serializer.Response(),
	})
}

// addSimulation godoc
// @Summary Add a simulation
// @ID addSimulation
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
func addSimulation(c *gin.Context) {

	userID, _ := c.Get("user_id")

	var u user.User
	err := u.ByID(userID.(uint))
	if common.ProvideErrorResponse(c, err) {
		return
	}

	var sim Simulation
	err = c.BindJSON(&sim)
	if err != nil {
		errormsg := "Bad request. Error binding form data to JSON: " + err.Error()
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errormsg,
		})
		return
	}

	// save new simulation to DB
	err = sim.save()
	if common.ProvideErrorResponse(c, err) {
		return
	}

	// add user to new simulation
	err = sim.addUser(&(u.User))
	if common.ProvideErrorResponse(c, err) == false {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK.",
		})
	}
}

// updateSimulation godoc
// @Summary Update a simulation
// @ID updateSimulation
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
func updateSimulation(c *gin.Context) {

	// TODO check if user has access to this simulation

	simID, err := common.GetSimulationID(c)
	if err != nil {
		return
	}

	var modifiedSim Simulation
	err = c.BindJSON(&modifiedSim)
	if err != nil {
		errormsg := "Bad request. Error binding form data to JSON: " + err.Error()
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errormsg,
		})
		return
	}

	var sim Simulation
	err = sim.ByID(uint(simID))
	if common.ProvideErrorResponse(c, err) {
		return
	}

	err = sim.update(modifiedSim)
	if common.ProvideErrorResponse(c, err) == false {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK.",
		})
	}
}

// getSimulation godoc
// @Summary Get simulation
// @ID getSimulation
// @Produce  json
// @Tags simulations
// @Success 200 {object} common.SimulationResponse "Simulation requested by user"
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param simulationID path int true "Simulation ID"
// @Router /simulations/{simulationID} [get]
func getSimulation(c *gin.Context) {

	// TODO check if user has access to this simulation

	simID, err := common.GetSimulationID(c)
	if err != nil {
		return
	}

	var sim Simulation
	err = sim.ByID(uint(simID))
	if common.ProvideErrorResponse(c, err) {
		return
	}

	serializer := common.SimulationSerializer{c, sim.Simulation}
	c.JSON(http.StatusOK, gin.H{
		"simulation": serializer.Response(),
	})
}

// deleteSimulation godoc
// @Summary Delete a simulation
// @ID deleteSimulation
// @Tags simulations
// @Produce json
// @Success 200 "OK."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param simulationID path int true "Simulation ID"
// @Router /simulations/{simulationID} [delete]
func deleteSimulation(c *gin.Context) {

	// TODO check if user has access to this simulation

	simID, err := common.GetSimulationID(c)
	if err != nil {
		return
	}

	var sim Simulation
	err = sim.ByID(uint(simID))
	if common.ProvideErrorResponse(c, err) {
		return
	}

	err = sim.delete()
	if common.ProvideErrorResponse(c, err) == false {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK.",
		})
	}
}

// getUsersOfSimulation godoc
// @Summary Get users of simulation
// @ID getUsersOfSimulation
// @Produce  json
// @Tags simulations
// @Success 200 {array} common.UserResponse "Array of users that have access to the simulation"
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param simulationID path int true "Simulation ID"
// @Router /simulations/{simulationID}/users/ [get]
func getUsersOfSimulation(c *gin.Context) {

	// TODO check if user has access to this simulation

	simID, err := common.GetSimulationID(c)
	if err != nil {
		return
	}

	var sim Simulation
	err = sim.ByID(uint(simID))
	if common.ProvideErrorResponse(c, err) {
		return
	}

	// Find all users of simulation
	allUsers, _, err := sim.getUsers()
	if common.ProvideErrorResponse(c, err) {
		return
	}

	serializer := common.UsersSerializer{c, allUsers}
	c.JSON(http.StatusOK, gin.H{
		"users": serializer.Response(false),
	})
}

// addUserToSimulation godoc
// @Summary Add a user to a a simulation
// @ID addUserToSimulation
// @Tags simulations
// @Produce json
// @Success 200 "OK."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param simulationID path int true "Simulation ID"
// @Param username query string true "User name"
// @Router /simulations/{simulationID}/user [put]
func addUserToSimulation(c *gin.Context) {

	// TODO check if user has access to this simulation

	simID, err := common.GetSimulationID(c)
	if err != nil {
		return
	}

	var sim Simulation
	err = sim.ByID(uint(simID))
	if common.ProvideErrorResponse(c, err) {
		return
	}

	username := c.Request.URL.Query().Get("username")

	var u user.User
	err = u.ByUsername(username)
	if common.ProvideErrorResponse(c, err) {
		return
	}

	err = sim.addUser(&(u.User))
	if common.ProvideErrorResponse(c, err) {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OK.",
	})
}

// deleteUserFromSimulation godoc
// @Summary Delete a user from a simulation
// @ID deleteUserFromSimulation
// @Tags simulations
// @Produce json
// @Success 200 "OK."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param simulationID path int true "Simulation ID"
// @Param username query string true "User name"
// @Router /simulations/{simulationID}/user [delete]
func deleteUserFromSimulation(c *gin.Context) {

	// TODO check if user has access to this simulation

	simID, err := common.GetSimulationID(c)
	if err != nil {
		return
	}

	var sim Simulation
	err = sim.ByID(uint(simID))
	if common.ProvideErrorResponse(c, err) {
		return
	}

	username := c.Request.URL.Query().Get("username")

	err = sim.deleteUser(username)
	if common.ProvideErrorResponse(c, err) {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OK.",
	})
}
