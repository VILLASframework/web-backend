package endpoints

import (
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/queries"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/serializers"
	"github.com/gin-gonic/gin"
	"net/http"
)



func userReadAllEp(c *gin.Context) {
	allUsers, _, _ := queries.FindAllUsers()
	serializer := serializers.UsersSerializer{c, allUsers}
	c.JSON(http.StatusOK, gin.H{
		"users": serializer.Response(),
	})
}

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