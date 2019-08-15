package user

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

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

func RegisterAuthenticate(r *gin.RouterGroup) {
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

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Authenticated",
		"token":   tokenString,
		"user":    user.User,
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

	err := common.ValidateRole(c, common.ModelUser, common.Read)
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

	c.JSON(http.StatusOK, gin.H{"users": users})
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
		"id": newUser.ID,
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
	var oldUser User
	toBeUpdatedID, _ := common.UintParamFromCtx(c, "userID")
	err = oldUser.ByID(toBeUpdatedID)
	if err != nil {
		c.JSON(http.StatusNotFound, fmt.Sprintf("%v", err))
		return
	}

	// Cases
	// 1: If the logged in user has NOT the same id as the user that is
	// going to be updated AND the role is NOT admin (is already saved
	// in the context from the Authentication middleware) the operation
	// is illegal
	// 2: If the udpate is done by the Admin every field can be updated
	// 3: If the update is done by a User everything can be updated
	// except Role
	callerID, _ := c.Get(common.UserIDCtx)
	callerRole, _ := c.Get(common.UserRoleCtx)

	if toBeUpdatedID != callerID && callerRole != "Admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "Invalid authorization",
		})
		return
	}

	// Bind the (context) with the updateUserRequest struct
	var req updateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"success": false,
			"message": fmt.Sprintf("%v", err),
		})
		return
	}

	// Validate the request based on struct updateUserRequest json tags
	if err = req.validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": fmt.Sprintf("%v", err),
		})
		return
	}

	// Create the updatedUser from oldUser considering callerRole (in
	// case that the request updates the role of the old user)
	updatedUser, err := req.updatedUser(callerRole, oldUser)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": fmt.Sprintf("%v", err),
		})
		return
	}

	// Finaly update the user
	err = oldUser.update(updatedUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Unable to update user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id": updatedUser.ID,
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

	err = user.ByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, fmt.Sprintf("%v", err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user.User})
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

	c.JSON(http.StatusOK, gin.H{
		"id": user.ID,
	})
}
