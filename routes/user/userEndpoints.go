package user

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
)

func RegisterUserEndpoints(r *gin.RouterGroup){
	r.GET("/", GetUsers)
	r.POST("/", AddUser)
	r.PUT("/:userID", UpdateUser)
	r.GET("/:userID", GetUser)
	r.DELETE("/:userID", DeleteUser)
	//r.GET("/me", userSelfEp) // TODO redirect to users/:userID
}

// GetUsers godoc
// @Summary Get all users
// @ID GetUsers
// @Produce  json
// @Tags users
// @Success 200 {array} common.UserResponse "Array of users"
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Router /users [get]
func GetUsers(c *gin.Context) {
	allUsers, _, _ := FindAllUsers()
	serializer := common.UsersSerializer{c, allUsers}
	c.JSON(http.StatusOK, gin.H{
		"users": serializer.Response(),
	})
}

// AddUser godoc
// @Summary Add a user
// @ID AddUser
// @Accept json
// @Produce json
// @Tags users
// @Param inputUser body common.UserResponse true "User to be added"
// @Success 200 "OK."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Router /users [post]
func AddUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}

// UpdateUser godoc
// @Summary Update a user
// @ID UpdateUser
// @Tags users
// @Accept json
// @Produce json
// @Param inputUser body common.UserResponse true "User to be updated"
// @Success 200 "OK."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param userID path int true "User ID"
// @Router /users/{userID} [put]
func UpdateUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}

// GetUser godoc
// @Summary Get user
// @ID GetUser
// @Produce  json
// @Tags users
// @Success 200 {object} common.UserResponse "User requested by user"
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param userID path int true "User ID"
// @Router /users/{userID} [get]
func GetUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}

// DeleteUser godoc
// @Summary Delete a user
// @ID DeleteUser
// @Tags users
// @Produce json
// @Success 200 "OK."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param userID path int true "User ID"
// @Router /users/{userID} [delete]
func DeleteUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}
