package user

import (
	//"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"github.com/gin-gonic/gin"
	"net/http"
)

// `/authenticate` endpoint does not require Authentication
func VisitorAuthenticate(r *gin.RouterGroup) {
	r.POST("", authenticationEp)
}

func UsersRegister(r *gin.RouterGroup) {
	r.POST("/users", userRegistrationEp)
	r.PUT("/:UserID", userUpdateEp)
	r.GET("/", usersReadEp)
	r.GET("/:UserID", userReadEp)
	//r.GET("/me", userSelfEp) // TODO: this conflicts with GET /:userID
	r.DELETE("/:UserID", userDeleteEp)
}

func authenticationEp(c *gin.Context) {

	// Bind the response (context) with the Credentials struct
	var loginRequest Credentials
	if err := c.BindJSON(&loginRequest); err != nil {
		// TODO: do something other than panic ...
		panic(err)
	}

	// Check if the Username or Password are empty
	if loginRequest.Username == "" || loginRequest.Password == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Invalid credentials",
		})
		return
	}

	// Find the username in the database
	user, err := FindUserByUsername(loginRequest.Username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "User not found",
		})
		return
	}

	// Validate the password
	if user.validatePassword(loginRequest.Password) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Invalid password",
		})
		return
	}

	// TODO: generate jwt

	c.JSON(http.StatusOK, gin.H{
		"success":          true,
		"message":          "Authenticated",
		"token":            "NOT yet implemented",
		"Original request": loginRequest, // TODO: remove that
	})
}

func usersReadEp(c *gin.Context) {
	//// dummy TODO: check in the middleware if the user is authorized
	//authorized := false
	//// TODO: move this redirect in the authentication middleware
	//if !authorized {
	//c.Redirect(http.StatusSeeOther, "/authenticate")
	//return
	//}
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
