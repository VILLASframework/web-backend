package user

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
)

// TODO: the signing secret must be environmental variable
const jwtSigningSecret = "This should NOT be here!!@33$8&"
const weekHours = time.Hour * 24 * 7

type tokenClaims struct {
	UserID uint   `json:"id"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

type AuthResponse struct {
	Success bool                `json:"success"`
	Message string              `json:"message"`
	Token   string              `json:"token"`
	User    common.UserResponse `json:"user"`
}

func VisitorAuthenticate(r *gin.RouterGroup) {
	r.POST("", authenticate)
}

func RegisterUserEndpoints(r *gin.RouterGroup) {
	r.POST("", addUser)
	r.PUT("/:userID", updateUser)
	r.GET("", getUsers)
	r.GET("/:userID", getUser)
	r.DELETE("/:userID", deleteUser)
}

// authenticate godoc
// @Summary Authentication for user
// @ID authenticate
// @Accept json
// @Produce json
// @Tags users
// @Param inputUser body user.loginRequest true "loginRequest of user"
// @Success 200 {object} user.AuthResponse "JSON web token and message"
// @Failure 401 "Unauthorized Access"
// @Failure 404 "Not found"
// @Failure 422 "Unprocessable entity."
// @Router /authenticate [post]
func authenticate(c *gin.Context) {

	// Bind the response (context) with the loginRequest struct
	var credentials loginRequest
	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"success": false,
			"message": fmt.Sprintf("%v", err),
		})
		return
	}

	// Validate the login request
	if errs := credentials.validate(); errs != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": fmt.Sprintf("%v", errs),
		})
		return
	}

	// Check if the Username or Password are empty
	if credentials.Username == "" || credentials.Password == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Invalid credentials",
		})
		return
	}

	// Find the username in the database
	var user User
	err := user.ByUsername(credentials.Username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": "User not found",
		})
		return
	}

	// Validate the password
	err = user.validatePassword(credentials.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "Invalid password",
		})
		return
	}

	// create authentication token
	claims := tokenClaims{
		user.ID,
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

	serializer := common.UserSerializer{c, user.User}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Authenticated",
		"token":   tokenString,
		"user":    serializer.Response(false),
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

	err := common.ValidateRole(c, common.ModelUsers, common.Read)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, fmt.Sprintf("%v", err))
		return
	}

	db := common.GetDB()
	var users []common.User
	err = db.Order("ID asc").Find(&users).Error
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

	err := common.ValidateRole(c, common.ModelUser, common.Create)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, fmt.Sprintf("%v", err))
		return
	}

	// Bind the request
	var req addUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"success": false,
			"message": fmt.Sprintf("%v", err),
		})
		return
	}

	// Validate the request
	if err = req.validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": fmt.Sprintf("%v", err),
		})
		return
	}

	// Create the new user from the request
	newUser := req.createUser()

	// Check that the username is NOT taken
	err = newUser.ByUsername(newUser.Username)
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
// @Param inputUser body common.User true "User to be updated (anything except for ID can be changed, role can only be change by admin)"
// @Success 200 "OK."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param userID path int true "User ID"
// @Router /users/{userID} [put]
func updateUser(c *gin.Context) {

	err := common.ValidateRole(c, common.ModelUser, common.Update)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, fmt.Sprintf("%v", err))
		return
	}

	// Find the user
	var user User
	toBeUpdatedID, _ := common.UintParamFromCtx(c, "userID")
	err = user.ByID(toBeUpdatedID)
	if err != nil {
		c.JSON(http.StatusNotFound, fmt.Sprintf("%v", err))
		return
	}

	// Cases
	// 1: If the logged in user has NOT the same id as the user that is
	// going to be updated AND the role is NOT admin (is already saved
	// in the context from the Authentication middleware) the operation
	// is elegal
	// 2: If the udpate is done by the Admin every field can be updated
	// 3: If the update is done by a User everything except Role
	userID, _ := c.Get(common.UserIDCtx)
	userRole, _ := c.Get(common.UserRoleCtx)

	if toBeUpdatedID != userID && userRole != "Admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "Invalid authorization",
		})
		return
	}

	// Bind the (context) with the User struct
	var req updateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"success": false,
			"message": fmt.Sprintf("%v", err),
		})
		return
	}

	// Validate the request (taking into acount the role)
	if err = req.validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": fmt.Sprintf("%v", err),
		})
		return
	}

	updatedUser := req.createUser(userRole)

	// Check that the username is NOT taken
	err = updatedUser.ByUsername(updatedUser.Username)
	if err == nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Username is already taken",
		})
		return
	}

	// Hash the password before updating it to the DB
	err = updatedUser.setPassword(updatedUser.Password)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Unable to encrypt new password",
		})
		return
	}

	// To change the role of a user admin role is required
	if (updatedUser.Role != user.Role) && (userRole != "Admin") {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "Invalid authorization. User role can only be changed by Admin",
		})
		return
	}

	// Finaly update the user
	err = user.update(updatedUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to update user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": fmt.Sprintf(updatedUser.Username),
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

	err := common.ValidateRole(c, common.ModelUser, common.Read)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, fmt.Sprintf("%v", err))
		return
	}

	var user User
	id, _ := common.UintParamFromCtx(c, "userID")

	reqUserID, _ := c.Get(common.UserIDCtx)
	reqUserRole, _ := c.Get(common.UserRoleCtx)

	if id != reqUserID && reqUserRole != "Admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "Invalid authorization",
		})
		return
	}

	err = user.ByID(id)
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

	err := common.ValidateRole(c, common.ModelUser, common.Delete)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, fmt.Sprintf("%v", err))
		return
	}

	var user User
	id, _ := common.UintParamFromCtx(c, "userID")

	// Check that the user exist
	err = user.ByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, fmt.Sprintf("%v", err))
		return
	}

	// Try to remove user
	err = user.remove()
	if err != nil {
		c.JSON(http.StatusInternalServerError, fmt.Sprintf("%v", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{})
}
