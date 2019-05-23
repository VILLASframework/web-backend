package endpoints

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/queries"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/serializers"
)



func userReadAllEp(c *gin.Context) {
	allUsers, _, _ := queries.FindAllUsers()
	serializer := serializers.UsersSerializer{c, allUsers}
	c.JSON(http.StatusOK, gin.H{
		"users": serializer.Response(),
	})
}

// userReadAllSimEp godoc
// @Summary Get users of simulation
// @ID GetAllUsersOfSimulation
// @Produce  json
// @Tags user
// @Success 200 {array} common.User "Array of users that have access to the simulation"
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param simulationID path int true "Simulation ID"
// @Router /simulations/{simulationID}/users [get]
func userReadAllSimEp(c *gin.Context) {

	simID, err := GetSimulationID(c)
	if err != nil {
		return
	}

	sim, err := queries.FindSimulation(simID)
	if common.ProvideErrorResponse(c, err) {
		return
	}

	// Find all users of simulation
	allUsers, _, err := queries.FindAllUsersSim(&sim)
	if common.ProvideErrorResponse(c, err) {
		return
	}

	serializer := serializers.UsersSerializer{c, allUsers}
	c.JSON(http.StatusOK, gin.H{
		"users": serializer.Response(),
	})
}

func userRegistrationEp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}

func userUpdateEp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}

// userUpdateSimEp godoc
// @Summary Add user to simulation
// @ID AddUserToSimulation
// @Tags user
// @Success 200 "OK."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param simulationID path int true "Simulation ID"
// @Param username path int true "Username of user to be added"
// @Router /simulations/{simulationID}/users/{username} [put]
func userUpdateSimEp(c *gin.Context) {


	simID, err := GetSimulationID(c)
	if err != nil {
		return
	}

	sim, err := queries.FindSimulation(simID)
	if common.ProvideErrorResponse(c, err) {
		return
	}

	username := c.Param("username")

	user, err := queries.FindUserByName(username)
	if common.ProvideErrorResponse(c, err) {
		return
	}

	err = queries.AddUserToSim(&sim, &user)
	if common.ProvideErrorResponse(c, err){
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OK.",
	})
}

func userReadEp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}

func userReadSimEp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}

func userDeleteEp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}

// userDeleteSimEp godoc
// @Summary Delete user from simulation
// @ID DeleteUserFromSimulation
// @Tags user
// @Success 200 "OK."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param simulationID path int true "Simulation ID"
// @Param username path int true "Username of user"
// @Router /simulations/{simulationID}/users/{username} [delete]
func userDeleteSimEp(c *gin.Context) {

	simID, err := GetSimulationID(c)
	if err != nil {
		return
	}

	sim, err := queries.FindSimulation(simID)
	if common.ProvideErrorResponse(c, err) {
		return
	}

	username := c.Param("username")

	err = queries.RemoveUserFromSim(&sim, username)
	if common.ProvideErrorResponse(c, err) {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OK.",
	})
}

func userSelfEp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}