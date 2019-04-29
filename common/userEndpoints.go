package common

import (
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
)

func UsersRegister(r *gin.RouterGroup) {
	r.GET("/", usersReadEp)
	r.POST("/", userRegistrationEp)
	r.PUT("/:userID", userUpdateEp)
	r.GET("/:userID", userReadEp)
	r.DELETE("/:userID", userDeleteEp)
	//r.GET("/me", userSelfEp) // TODO: this conflicts with GET /:userID
}

func usersReadEp(c *gin.Context) {
	allUsers, _, _ := FindAllUsers()
	serializer := UsersSerializer{c, allUsers}
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

func userReadEp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}

func userDeleteEp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}

func userSelfEp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}
