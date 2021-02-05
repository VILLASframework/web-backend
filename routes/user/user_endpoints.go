/** User package, endpoints.
*
* @author Sonja Happ <sonja.happ@eonerc.rwth-aachen.de>
* @copyright 2014-2019, Institute for Automation of Complex Power Systems, EONERC
* @license GNU General Public License (version 3)
*
* VILLASweb-backend-go
*
* This program is free software: you can redistribute it and/or modify
* it under the terms of the GNU General Public License as published by
* the Free Software Foundation, either version 3 of the License, or
* any later version.
*
* This program is distributed in the hope that it will be useful,
* but WITHOUT ANY WARRANTY; without even the implied warranty of
* MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
* GNU General Public License for more details.
*
* You should have received a copy of the GNU General Public License
* along with this program.  If not, see <http://www.gnu.org/licenses/>.
*********************************************************************************/
package user

import (
	"fmt"
	"net/http"
	"strings"

	"git.rwth-aachen.de/acs/public/villas/web-backend-go/helper"

	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
)

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
// @Success 200 {object} api.ResponseUsers "Array of users"
// @Failure 404 {object} api.ResponseError "Not found"
// @Failure 422 {object} api.ResponseError "Unprocessable entity"
// @Failure 500 {object} api.ResponseError "Internal server error"
// @Router /users [get]
// @Security Bearer
func getUsers(c *gin.Context) {

	err := database.ValidateRole(c, database.ModelUsers, database.Read)
	if err != nil {
		helper.UnprocessableEntityError(c, err.Error())
		return
	}

	db := database.GetDB()
	var users []database.User
	err = db.Order("ID asc").Find(&users).Error
	if !helper.DBError(c, err) {
		c.JSON(http.StatusOK, gin.H{"users": users})
	}

}

// AddUser godoc
// @Summary Add a user
// @ID AddUser
// @Accept json
// @Produce json
// @Tags users
// @Param inputUser body user.addUserRequest true "User to be added"
// @Success 200 {object} api.ResponseUser "Contains added user object"
// @Failure 400 {object} api.ResponseError "Bad request"
// @Failure 422 {object} api.ResponseError "Unprocessable entity"
// @Failure 500 {object} api.ResponseError "Internal server error"
// @Router /users [post]
// @Security Bearer
func addUser(c *gin.Context) {

	err := database.ValidateRole(c, database.ModelUser, database.Create)
	if err != nil {
		helper.UnprocessableEntityError(c, err.Error())
		return
	}

	// Bind the request
	var req addUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.BadRequestError(c, err.Error())
		return
	}

	// Validate the request
	if err = req.validate(); err != nil {
		helper.BadRequestError(c, err.Error())
		return
	}

	// Create the new user from the request
	newUser := req.createUser()

	// Check that the username is NOT taken
	err = newUser.ByUsername(newUser.Username)
	if err == nil {
		helper.UnprocessableEntityError(c, "Username is already taken")
		return
	}

	// Hash the password before saving it to the DB
	err = newUser.setPassword(newUser.Password)
	if err != nil {
		helper.InternalServerError(c, "Unable to encrypt the password")
		return
	}

	// Save the user in the DB
	err = newUser.save()
	if !helper.DBError(c, err) {
		c.JSON(http.StatusOK, gin.H{"user": newUser.User})
	}

}

// UpdateUser godoc
// @Summary Update a user
// @ID UpdateUser
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} api.ResponseUser "Contains updated user"
// @Failure 400 {object} api.ResponseError "Bad request."
// @Failure 403 {object} api.ResponseError "Access forbidden."
// @Failure 404 {object} api.ResponseError "Not found"
// @Failure 422 {object} api.ResponseError "Unprocessable entity"
// @Failure 500 {object} api.ResponseError "Internal server error"
// @Param inputUser body user.updateUserRequest true "User to be updated (anything except for ID can be changed, role can only be change by admin)"
// @Param userID path int true "User ID"
// @Router /users/{userID} [put]
// @Security Bearer
func updateUser(c *gin.Context) {

	// no need to validate the role since updating a single user is role independent
	//err := database.ValidateRole(c, database.ModelUser, database.Update)
	//if err != nil {
	//	helper.UnprocessableEntityError(c, err.Error())
	//	return
	//}

	// Get the user's (to be updated) ID from the context
	toBeUpdatedID, err := helper.GetIDOfElement(c, "userID", "path", -1)
	if err != nil {
		return
	}

	// Cases
	// 1: If the logged in user has NOT the same id as the user that is
	// going to be updated AND the role is NOT admin (is already saved
	// in the context from the Authentication middleware) the operation
	// is illegal
	// 2: If the udpate is done by the Admin every field can be updated
	// 3: If the update is done by a User everything can be updated
	// except Role and Active state

	// Get caller's ID from context
	callerID, _ := c.Get(database.UserIDCtx)

	// Get caller's Role from context
	callerRole, _ := c.Get(database.UserRoleCtx)

	if uint(toBeUpdatedID) != callerID && callerRole != "Admin" {
		helper.ForbiddenError(c, "Invalid authorization")
		return
	}

	// Bind the (context) with the updateUserRequest struct
	var req updateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.BadRequestError(c, fmt.Sprintf("%v", err))
		return
	}

	// Validate the request based on struct updateUserRequest json tags
	if err = req.validate(); err != nil {
		helper.BadRequestError(c, err.Error())
		return
	}

	// Find the user
	var oldUser User
	err = oldUser.ByID(uint(toBeUpdatedID))
	if helper.DBError(c, err) {
		return
	}

	// Create the updatedUser from oldUser considering callerRole (in
	// case that the request updates the role of the old user)
	updatedUser, err := req.updatedUser(callerID, callerRole, oldUser)
	if err != nil {
		if strings.Contains(err.Error(), "Admin") || strings.Contains(err.Error(), "pw not changed") {
			helper.ForbiddenError(c, err.Error())
		} else if strings.Contains(err.Error(), "Username") || strings.Contains(err.Error(), "old or admin password") {
			helper.BadRequestError(c, err.Error())
		} else { // password encryption failed
			helper.InternalServerError(c, err.Error())
		}
		return
	}

	// Finally update the user
	err = oldUser.update(updatedUser)
	if !helper.DBError(c, err) {
		c.JSON(http.StatusOK, gin.H{"user": updatedUser.User})
	}

}

// GetUser godoc
// @Summary Get user
// @ID GetUser
// @Produce  json
// @Tags users
// @Success 200 {object} api.ResponseUser "requested user"
// @Failure 403 {object} api.ResponseError "Access forbidden."
// @Failure 404 {object} api.ResponseError "Not found"
// @Failure 422 {object} api.ResponseError "Unprocessable entity"
// @Failure 500 {object} api.ResponseError "Internal server error"
// @Param userID path int true "User ID"
// @Router /users/{userID} [get]
// @Security Bearer
func getUser(c *gin.Context) {

	// role validation not needed because updating a single user is role-independent
	//err := database.ValidateRole(c, database.ModelUser, database.Read)
	//if err != nil {
	//	helper.UnprocessableEntityError(c, err.Error())
	//	return
	//}

	id, err := helper.GetIDOfElement(c, "userID", "path", -1)
	if err != nil {
		return
	}

	reqUserID, _ := c.Get(database.UserIDCtx)
	reqUserRole, _ := c.Get(database.UserRoleCtx)

	if uint(id) != reqUserID && reqUserRole != "Admin" {
		helper.ForbiddenError(c, "Invalid authorization")
		return
	}

	var user User
	err = user.ByID(uint(id))
	if !helper.DBError(c, err) {
		c.JSON(http.StatusOK, gin.H{"user": user.User})
	}

}

// DeleteUser godoc
// @Summary Delete a user
// @ID DeleteUser
// @Tags users
// @Produce json
// @Success 200 {object} api.ResponseUser "deleted user"
// @Failure 404 {object} api.ResponseError "Not found"
// @Failure 422 {object} api.ResponseError "Unprocessable entity"
// @Failure 500 {object} api.ResponseError "Internal server error"
// @Param userID path int true "User ID"
// @Router /users/{userID} [delete]
// @Security Bearer
func deleteUser(c *gin.Context) {

	err := database.ValidateRole(c, database.ModelUser, database.Delete)
	if err != nil {
		helper.UnprocessableEntityError(c, err.Error())
		return
	}

	var user User
	id, err := helper.GetIDOfElement(c, "userID", "path", -1)
	if err != nil {
		return
	}

	// Check that the user exist
	err = user.ByID(uint(id))
	if helper.DBError(c, err) {
		return
	}

	// Try to remove user
	err = user.remove()
	if !helper.DBError(c, err) {
		c.JSON(http.StatusOK, gin.H{"user": user.User})
	}

}
