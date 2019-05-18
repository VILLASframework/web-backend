package user

import (
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Credentials struct {
	Username string `form:"Username"`
	Password string `form:"Password"`
	Role     string `form:"Role"`
	Mail     string `form:"Mail"`
}

func UsersRegister(r *gin.RouterGroup) {
	r.POST("/authenticate", authenticationEp)
	r.GET("/", usersReadEp)
	r.POST("/", userRegistrationEp)
	r.PUT("/:UserID", userUpdateEp)
	r.GET("/:UserID", userReadEp)
	r.DELETE("/:UserID", userDeleteEp)
	//r.GET("/me", userSelfEp) // TODO: this conflicts with GET /:userID
}

func authenticationEp(c *gin.Context) {

	// Bind the response (context) with the Credentials struct
	var userLogin Credentials
	err := c.BindJSON(&userLogin)
	if err != nil {
		panic(err)
	}

	// Check if the Username or Password are empty
	if userLogin.Username == "" || userLogin.Password == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Invalid credentials",
		})
		return
	}

	// Find the username in the database
	db := common.GetDB()
	var user common.User
	err = db.Find(&user, "Username = ?", userLogin.Username).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "User not found",
		})
		return
	}

	// TODO: Validate password

	// TODO: generate jwt

	c.JSON(http.StatusOK, gin.H{
		"success":          true,
		"message":          "Authenticated",
		"token":            "NOT yet implemented",
		"Original request": userLogin, // TODO: remove that
	})
}

func usersReadEp(c *gin.Context) {
	allUsers, _, _ := FindAllUsers()
	serializer := UsersSerializer{c, allUsers}
	c.JSON(http.StatusOK, gin.H{
		"users": serializer.Response(),
	})
}

func userRegistrationEp(c *gin.Context) {
	//// dummy TODO: check in the middleware if the user is authorized
	//authorized := false
	//// TODO: move this redirect in the authentication middleware
	//if !authorized {
	//c.Redirect(http.StatusSeeOther, "/authenticate")
	//return
	//}
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
