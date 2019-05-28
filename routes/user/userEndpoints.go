package user

import (
	//"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
)

// TODO: the signing secret must be environmental variable
const jwtSigningSecret = "This should NOT be here!!@33$8&"
const weekHours = time.Hour * 24 * 7

type tokenClaims struct {
	UserID string `json:"id"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

type AuthResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Token   string `json:"token"`
}

// `/authenticate` endpoint does not require Authentication
func VisitorAuthenticate(r *gin.RouterGroup) {
	r.POST("", authenticationEp)
}

func RegisterUserEndpoints(r *gin.RouterGroup) {
	r.POST("", addUser)
	r.PUT("/:UserID", updateUser)
	r.GET("", getUsers)
	r.GET("/:UserID", getUser)
	r.DELETE("/:UserID", deleteUser)
}

// authenticationEp godoc
// @Summary Authentication for user
// @ID authenticationEp
// @Accept json
// @Produce json
// @Tags users
// @Param inputUser body user.Credentials true "Credentials of user"
// @Success 200 {object} user.AuthResponse "JSON web token and message"
// @Failure 401 "Unauthorized Access"
// @Failure 404 "Not found"
// @Failure 422 "Unprocessable entity."
// @Router /authenticate [post]
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
	var user User
	err := user.ByUsername(loginRequest.Username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "User not found",
		})
		return
	}

	// Validate the password
	err = user.validatePassword(loginRequest.Password)
	if err != nil {
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
		"success": true,
		"message": "Authenticated",
		"token":   tokenString,
	})
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
func getUsers(c *gin.Context) {
	//// dummy TODO: check in the middleware if the user is authorized
	//authorized := false
	//// TODO: move this redirect in the authentication middleware
	//if !authorized {
	//c.Redirect(http.StatusSeeOther, "/authenticate")
	//return
	//}

	db := common.GetDB()
	var users []common.User
	err := db.Order("ID asc").Find(&users).Error
	if common.ProvideErrorResponse(c, err) {
		return
	}
	serializer := common.UsersSerializer{c, users}
	c.JSON(http.StatusOK, gin.H{
		"users": serializer.Response(true),
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
func addUser(c *gin.Context) {

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
	err := newUser.ByUsername(newUser.Username)
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
func updateUser(c *gin.Context) {
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
func getUser(c *gin.Context) {

	var user User
	id, _ := strconv.ParseInt(c.Param("UserID"), 10, 64)

	err := user.byID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, fmt.Sprintf("%v", err))
		return
	}

	serializer := common.UserSerializer{c, user.User}
	c.JSON(http.StatusOK, gin.H{
		"user": serializer.Response(false),
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
func deleteUser(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}
