package user

import (
	//"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// TODO: the signing secret must be environmental variable
const jwtSigningSecret = "This should NOT be here!!@33$8&"
const weekHours = time.Hour * 24 * 7

type tokenClaims struct {
	UserID string `json:"id"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

// `/authenticate` endpoint does not require Authentication
func VisitorAuthenticate(r *gin.RouterGroup) {
	r.POST("", authenticationEp)
}

func UsersRegister(r *gin.RouterGroup) {
	r.POST("", userRegistrationEp)
	r.PUT("/:UserID", userUpdateEp)
	r.GET("", usersReadEp)
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
	user, err := findUserByUsername(loginRequest.Username)
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

	// create authentication token
	claims := tokenClaims{
		string(user.ID),
		user.Role,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(weekHours).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "http://web.villas.fein-aachen.org/",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(jwtSigningSecret))
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"success": false,
			"message": fmt.Sprintf("%v", err),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":          true,
		"message":          "Authenticated",
		"token":            tokenString,
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

	// Bind the response (context) with the User struct
	var newUser User
	if err := c.BindJSON(&newUser); err != nil {
		// TODO: do something other than panic ...
		panic(err)
	}

	// TODO: validate the User for:
	//       - username
	//       - email
	//       - role
	// and in case of error raise 422

	// Check that the username is NOT taken
	_, err := findUserByUsername(newUser.Username)
	if err == nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Username is already taken",
		})
		return
	}

	// Hash the password before saving it to the DB
	err = newUser.setPassword(newUser.Password)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Unable to encrypt the password",
		})
		return
	}

	// Save the user in the DB
	err = newUser.save()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to create new user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": fmt.Sprintf(newUser.Username),
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
