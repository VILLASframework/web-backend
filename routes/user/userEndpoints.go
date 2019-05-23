package user

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/simulation"
)


func RegisterUserEndpoints(r *gin.RouterGroup){
	r.GET("/", GetUsers)
	r.POST("/", AddUser)
	r.PUT("/:userID", UpdateUser)
	r.GET("/:userID", GetUser)
	r.GET("/:userID/simulations", GetSimulationsOfUser)
	r.DELETE("/:userID", DeleteUser)
	//r.GET("/me", userSelfEp) // TODO redirect to users/:userID
}

func RegisterUserEndpointsForSimulation(r *gin.RouterGroup){
	r.GET("/:simulationID/users", GetUsersOfSimulation)
	r.PUT("/:simulationID/user/:username", UpdateUserOfSimulation)
	r.DELETE("/:simulationID/user/:username", DeleteUserOfSimulation)

}

func GetUsers(c *gin.Context) {
	allUsers, _, _ := FindAllUsers()
	serializer := common.UsersSerializer{c, allUsers}
	c.JSON(http.StatusOK, gin.H{
		"users": serializer.Response(),
	})
}

// GetUsersOfSimulation godoc
// @Summary Get users of simulation
// @ID GetUsersOfSimulation
// @Produce  json
// @Tags user
// @Success 200 {array} common.UserResponse "Array of users that have access to the simulation"
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param simulationID path int true "Simulation ID"
// @Router /simulations/{simulationID}/users [get]
func GetUsersOfSimulation(c *gin.Context) {

	simID, err := common.GetSimulationID(c)
	if err != nil {
		return
	}

	sim, err := simulation.FindSimulation(simID)
	if common.ProvideErrorResponse(c, err) {
		return
	}

	// Find all users of simulation
	allUsers, _, err := FindAllUsersSim(&sim)
	if common.ProvideErrorResponse(c, err) {
		return
	}

	serializer := common.UsersSerializer{c, allUsers}
	c.JSON(http.StatusOK, gin.H{
		"users": serializer.Response(),
	})
}

func AddUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}

func UpdateUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}

// UpdateUserOfSimulation godoc
// @Summary Add user to simulation
// @ID UpdateUserOfSimulation
// @Tags user
// @Success 200 "OK."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param simulationID path int true "Simulation ID"
// @Param username path int true "Username of user to be added"
// @Router /simulations/{simulationID}/users/{username} [put]
func UpdateUserOfSimulation(c *gin.Context) {


	simID, err := common.GetSimulationID(c)
	if err != nil {
		return
	}

	sim, err := simulation.FindSimulation(simID)
	if common.ProvideErrorResponse(c, err) {
		return
	}

	username := c.Param("username")

	user, err := FindUserByName(username)
	if common.ProvideErrorResponse(c, err) {
		return
	}

	err = AddUserToSim(&sim, &user)
	if common.ProvideErrorResponse(c, err){
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OK.",
	})
}

func GetUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}

func GetSimulationsOfUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}

func DeleteUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}

// DeleteUserOfSimulation godoc
// @Summary Delete user from simulation
// @ID DeleteUserOfSimulation
// @Tags user
// @Success 200 "OK."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param simulationID path int true "Simulation ID"
// @Param username path int true "Username of user"
// @Router /simulations/{simulationID}/users/{username} [delete]
func DeleteUserOfSimulation(c *gin.Context) {

	simID, err := common.GetSimulationID(c)
	if err != nil {
		return
	}

	sim, err := simulation.FindSimulation(simID)
	if common.ProvideErrorResponse(c, err) {
		return
	}

	username := c.Param("username")

	err = RemoveUserFromSim(&sim, username)
	if common.ProvideErrorResponse(c, err) {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OK.",
	})
}
