package user

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	_ "git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/doc/api"
)

// TODO: the signing secret must be environmental variable
const jwtSigningSecret = "This should NOT be here!!@33$8&"
const weekHours = time.Hour * 24 * 7

func RegisterUserEndpoints(r *gin.RouterGroup) {
	r.POST("", addUser)
	r.PUT("/:userID", updateUser)
	r.GET("", getUsers)
	r.GET("/:userID", getUser)
	r.DELETE("/:userID", deleteUser)
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
		common.UnprocessableEntityError(c, err.Error())
		return
	}

	db := common.GetDB()
	var users []common.User
	err = db.Order("ID asc").Find(&users).Error
	if common.DBError(c, err) {
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
// @Param inputUser body user.validNewUser true "User to be added"
// @Success 200 {object} docs.ResponseUser "Contains added user object"
// @Failure 400 {object} docs.ResponseError "Bad request"
// @Failure 422 {object} docs.ResponseError "Unprocessable entity"
// @Failure 500 {object} docs.ResponseError "Internal server error"
// @Router /users [post]
func addUser(c *gin.Context) {

	err := common.ValidateRole(c, common.ModelUser, common.Create)
	if err != nil {
		common.UnprocessableEntityError(c, err.Error())
		return
	}

	// Bind the request
	var req addUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.BadRequestError(c, err.Error())
		return
	}

	// Validate the request
	if err = req.validate(); err != nil {
		common.UnprocessableEntityError(c, err.Error())
		return
	}

	// Create the new user from the request
	newUser := req.createUser()

	// Check that the username is NOT taken
	err = newUser.ByUsername(newUser.Username)
	if err == nil {
		common.UnprocessableEntityError(c, "Username is already taken")
		return
	}

	// Hash the password before saving it to the DB
	err = newUser.setPassword(newUser.Password)
	if err != nil {
		common.UnprocessableEntityError(c, "Unable to encrypt the password")
		return
	}

	// Save the user in the DB
	err = newUser.save()
	if err != nil {
		common.DBError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": newUser.User})
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
		common.UnprocessableEntityError(c, err.Error())
		return
	}

	// Get the user's (to be updated) ID from the context
	var oldUser User
	toBeUpdatedID, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		common.NotFoundError(c, fmt.Sprintf("Could not get user's ID from context"))
		return
	}

	// Find the user
	err = oldUser.ByID(uint(toBeUpdatedID))
	if err != nil {
		common.DBError(c, err)
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
		common.NotFoundError(c, fmt.Sprintf("Could not get caller's ID from context"))
		return
	}

	// Get caller's Role from context
	callerRole, exists := c.Get(common.UserRoleCtx)
	if !exists {
		common.NotFoundError(c, fmt.Sprintf("Could not get caller's Role from context"))
		return
	}

	if uint(toBeUpdatedID) != callerID && callerRole != "Admin" {
		common.ForbiddenError(c, "Invalid authorization")
		return
	}

	// Bind the (context) with the updateUserRequest struct
	var req updateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.UnprocessableEntityError(c, fmt.Sprintf("%v", err))
		return
	}

	// Validate the request based on struct updateUserRequest json tags
	if err = req.validate(); err != nil {
		common.BadRequestError(c, err.Error())
		return
	}

	// Create the updatedUser from oldUser considering callerRole (in
	// case that the request updates the role of the old user)
	updatedUser, err := req.updatedUser(callerRole, oldUser)
	if err != nil {
		if strings.Contains(err.Error(), "Admin") {
			common.ForbiddenError(c, err.Error())

		} else if strings.Contains(err.Error(), "Username") || strings.Contains(err.Error(), "password") {
			common.BadRequestError(c, err.Error())
		}
		return
	}

	// Finally update the user
	err = oldUser.update(updatedUser)
	if err != nil {
		common.DBError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": updatedUser.User})
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
		common.UnprocessableEntityError(c, err.Error())
		return
	}

	id, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		common.NotFoundError(c, fmt.Sprintf("Could not get user's ID from context"))
		return
	}

	reqUserID, _ := c.Get(common.UserIDCtx)
	reqUserRole, _ := c.Get(common.UserRoleCtx)

	if id != reqUserID && reqUserRole != "Admin" {
		common.ForbiddenError(c, "Invalid authorization")
		return
	}

	var user User
	err = user.ByID(uint(id))
	if err != nil {
		common.DBError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user.User})
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
		common.UnprocessableEntityError(c, err.Error())
		return
	}

	var user User
	id, err := strconv.Atoi(c.Param("userID"))
	if err != nil {
		common.NotFoundError(c, fmt.Sprintf("Could not get user's ID from context"))
		return
	}

	// Check that the user exist
	err = user.ByID(uint(id))
	if err != nil {
		common.DBError(c, err)
		return
	}

	// Try to remove user
	err = user.remove()
	if err != nil {
		common.DBError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user.User})
}
