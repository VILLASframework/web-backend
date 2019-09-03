package user

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	_ "git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/doc/api"
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
// @Success 200 {object} docs.ResponseAuthenticate "JSON web token, success status, message and authenticated user object"
// @Failure 400 {object} docs.ResponseError "Bad request"
// @Failure 401 {object} docs.ResponseError "Unauthorized"
// @Failure 404 {object} docs.ResponseError "Not found"
// @Failure 422 {object} docs.ResponseError "Unprocessable entity."
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
// @Success 200 {object} docs.ResponseUsers "Array of users"
// @Failure 404 {object} docs.ResponseError "Not found"
// @Failure 422 {object} docs.ResponseError "Unprocessable entity"
// @Failure 500 {object} docs.ResponseError "Internal server error"
// @Router /users [get]
func getUsers(c *gin.Context) {

	err := common.ValidateRole(c, common.ModelUsers, common.Read)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"success": false,
			"message": fmt.Sprintf("%v", err),
		})
		return
	}

	db := common.GetDB()
	var users []common.User
	err = db.Order("ID asc").Find(&users).Error
	if common.ProvideErrorResponse(c, err) {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"users":   users,
	})
}

// AddUser godoc
// @Summary Add a user
// @ID AddUser
// @Accept json
// @Produce json
// @Tags users
// @Param inputUser body user.validNewUser true "User to be added"
// @Success 200 {object} docs.ResponseUser "Contains added user object"
// @Failure 400 {object} docs.ResponseError "Bad request"
// @Failure 422 {object} docs.ResponseError "Unprocessable entity"
// @Failure 500 {object} docs.ResponseError "Internal server error"
// @Router /users [post]
func addUser(c *gin.Context) {

	err := common.ValidateRole(c, common.ModelUser, common.Create)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"success": false,
			"message": fmt.Sprintf("%v", err),
		})
		return
	}

	// Bind the request
	var req addUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": fmt.Sprintf("%v", err),
		})
		return
	}

	// Validate the request
	if err = req.validate(); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
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
			"success": false,
			"message": "Username is already taken",
		})
		return
	}

	// Hash the password before saving it to the DB
	err = newUser.setPassword(newUser.Password)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"success": false,
			"message": "Unable to encrypt the password",
		})
		return
	}

	// Save the user in the DB
	err = newUser.save()
	if err != nil {
		common.ProvideErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"user":    newUser.User,
	})
}

// UpdateUser godoc
// @Summary Update a user
// @ID UpdateUser
// @Tags users
// @Accept json
// @Produce json
// @Param inputUser body user.validUpdatedRequest true "User to be updated (anything except for ID can be changed, role can only be change by admin)"
// @Success 200 {object} docs.ResponseUser "Contains updated user"
// @Failure 400 {object} docs.ResponseError "Bad request."
// @Failure 403 {object} docs.ResponseError "Access forbidden."
// @Failure 404 {object} docs.ResponseError "Not found"
// @Failure 422 {object} docs.ResponseError "Unprocessable entity"
// @Failure 500 {object} docs.ResponseError "Internal server error"
// @Param userID path int true "User ID"
// @Router /users/{userID} [put]
func updateUser(c *gin.Context) {

	err := common.ValidateRole(c, common.ModelUser, common.Update)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"success": false,
			"message": fmt.Sprintf("%v", err),
		})
		return
	}

	// Get the user's (to be updated) ID from the context
	var oldUser User
	toBeUpdatedID, err := common.UintParamFromCtx(c, "userID")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": fmt.Sprintf("Could not get user's ID from context"),
		})
		return
	}

	// Find the user
	err = oldUser.ByID(toBeUpdatedID)
	if err != nil {
		common.ProvideErrorResponse(c, err)
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

	// Get caller's ID from context
	callerID, exists := c.Get(common.UserIDCtx)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": fmt.Sprintf("Could not get caller's ID from context"),
		})
		return
	}

	// Get caller's Role from context
	callerRole, exists := c.Get(common.UserRoleCtx)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": fmt.Sprintf("Could not get caller's Role from context"),
		})
		return
	}

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
		if strings.Contains(err.Error(), "Admin") {
			c.JSON(http.StatusForbidden, gin.H{
				"success": false,
				"message": fmt.Sprintf("%v", err),
			})

		} else if strings.Contains(err.Error(), "Username") || strings.Contains(err.Error(), "password") {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": fmt.Sprintf("%v", err),
			})
		}
		return
	}

	// Finally update the user
	err = oldUser.update(updatedUser)
	if err != nil {
		common.ProvideErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"user":    updatedUser.User,
	})
}

// GetUser godoc
// @Summary Get user
// @ID GetUser
// @Produce  json
// @Tags users
// @Success 200 {object} docs.ResponseUser "requested user"
// @Failure 403 {object} docs.ResponseError "Access forbidden."
// @Failure 404 {object} docs.ResponseError "Not found"
// @Failure 422 {object} docs.ResponseError "Unprocessable entity"
// @Failure 500 {object} docs.ResponseError "Internal server error"
// @Param userID path int true "User ID"
// @Router /users/{userID} [get]
func getUser(c *gin.Context) {

	err := common.ValidateRole(c, common.ModelUser, common.Read)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"success": false,
			"message": fmt.Sprintf("%v", err),
		})
		return
	}

	id, err := common.UintParamFromCtx(c, "userID")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": fmt.Sprintf("Could not get user's ID from context"),
		})
		return
	}

	reqUserID, _ := c.Get(common.UserIDCtx)
	reqUserRole, _ := c.Get(common.UserRoleCtx)

	if id != reqUserID && reqUserRole != "Admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"success": false,
			"message": "Invalid authorization",
		})
		return
	}

	var user User
	err = user.ByID(id)
	if err != nil {
		common.ProvideErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"user":    user.User,
	})
}

// DeleteUser godoc
// @Summary Delete a user
// @ID DeleteUser
// @Tags users
// @Produce json
// @Success 200 {object} docs.ResponseUser "deleted user"
// @Failure 404 {object} docs.ResponseError "Not found"
// @Failure 422 {object} docs.ResponseError "Unprocessable entity"
// @Failure 500 {object} docs.ResponseError "Internal server error"
// @Param userID path int true "User ID"
// @Router /users/{userID} [delete]
func deleteUser(c *gin.Context) {

	err := common.ValidateRole(c, common.ModelUser, common.Delete)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"success": false,
			"message": fmt.Sprintf("%v", err),
		})
		return
	}

	var user User
	id, err := common.UintParamFromCtx(c, "userID")
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": fmt.Sprintf("Could not get user's ID from context"),
		})
		return
	}

	// Check that the user exist
	err = user.ByID(uint(id))
	if err != nil {
		common.ProvideErrorResponse(c, err)
		return
	}

	// Try to remove user
	err = user.remove()
	if err != nil {
		common.ProvideErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"user":    user.User,
	})
}
