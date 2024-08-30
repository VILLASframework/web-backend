/**
* This file is part of VILLASweb-backend-go
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

package usergroup

import (
	"net/http"

	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/helper"
	"github.com/gin-gonic/gin"
)

func RegisterUserGroupEndpoints(r *gin.RouterGroup) {
	r.POST("", addUserGroup)
	r.PUT("/:userGroupID", updateUserGroup)
	r.GET("", getUserGroups)
	r.GET("/:userGroupID", getUserGroup)
	r.DELETE("/:userGroupID", deleteUserGroup)
}

// addUserGroup godoc
// @Summary Add a user group
// @ID addUserGroup
// @Accept json
// @Produce json
// @Tags usergroups
// @Success 200 {object} api.ResponseUserGroup "user group that was added"
// @Failure 400 {object} api.ResponseError "Bad request"
// @Failure 404 {object} api.ResponseError "Not found"
// @Failure 422 {object} api.ResponseError "Unprocessable entity"
// @Failure 500 {object} api.ResponseError "Internal server error"
// @Param inputUserGroup body usergroup.addUserGroupRequest true "User group to be added"
// @Router /usergroups [post]
// @Security Bearer
func addUserGroup(c *gin.Context) {
	ok, _ := database.CheckUserGroupPermissions(c, database.Create, "none", -1)
	if !ok {
		return
	}

	var req addUserGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.BadRequestError(c, err.Error())
		return
	}

	// Validate the request
	if err := req.validate(); err != nil {
		helper.UnprocessableEntityError(c, err.Error())
		return
	}

	// Create the new user group from the request
	newUserGroup := req.createUserGroup()

	// Save the new user group to the database
	err := newUserGroup.save()
	if !helper.DBError(c, err) {
		c.JSON(http.StatusOK, gin.H{"usergroup": newUserGroup.UserGroup})
	}
}

// updateUserGroup godoc
// @Summary Update a user group
// @ID updateUserGroup
// @Tags usergroups
// @Accept json
// @Produce json
// @Success 200 {object} api.ResponseUserGroup "User group that was updated"
// @Failure 400 {object} api.ResponseError "Bad request"
// @Failure 404 {object} api.ResponseError "Not found"
// @Failure 422 {object} api.ResponseError "Unprocessable entity"
// @Failure 500 {object} api.ResponseError "Internal server error"
// @Param inputUserGroup body usergroup.updateUserGroupRequest true "User group to be updated"
// @Param usergroupID path int true "User group ID"
// @Router /usergroups/{usergroupID} [put]
// @Security Bearer
func updateUserGroup(c *gin.Context) {
	ok, oldUserGroup_r := database.CheckUserGroupPermissions(c, database.Update, "path", -1)
	if !ok {
		return
	}

	var req updateUserGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.BadRequestError(c, err.Error())
		return
	}

	if err := req.UserGroup.validate(); err != nil {
		helper.BadRequestError(c, err.Error())
		return
	}

	var oldUserGroup UserGroup
	oldUserGroup.UserGroup = oldUserGroup_r
	updatedUserGroup := req.updatedUserGroup(oldUserGroup)

	// update the user group in the database
	err := oldUserGroup.update(updatedUserGroup, req.UserGroup.ScenarioMappings)
	if !helper.DBError(c, err) {
		c.JSON(http.StatusOK, gin.H{"usergroup": updatedUserGroup.UserGroup})
	}
}

// getUserGroups godoc
// @Summary Get all user groups
// @ID getUserGroups
// @Produce json
// @Tags usergroups
// @Success 200 {object} api.ResponseUserGroups "List of user groups"
// @Failure 404 {object} api.ResponseError "Not found"
// @Failure 422 {object} api.ResponseError "Unprocessable entity"
// @Failure 500 {object} api.ResponseError "Internal server error"
// @Router /usergroups [get]
// @Security Bearer
func getUserGroups(c *gin.Context) {

	err := database.ValidateRole(c, database.ModelUserGroup, database.Read)
	if err != nil {
		helper.UnprocessableEntityError(c, err.Error())
		return
	}

	db := database.GetDB()
	var usergroups []database.UserGroup
	err = db.Preload("ScenarioMappings.Scenario").Order("ID asc").Find(&usergroups).Error
	if !helper.DBError(c, err) {
		c.JSON(http.StatusOK, gin.H{"usergroups": usergroups})
	}
}

// getUserGroup godoc
// @Summary Get user group by ID
// @ID getUserGroup
// @Produce  json
// @Tags usergroups
// @Success 200 {object} api.ResponseUserGroup "requested user group"
// @Failure 403 {object} api.ResponseError "Access forbidden."
// @Failure 404 {object} api.ResponseError "Not found"
// @Failure 422 {object} api.ResponseError "Unprocessable entity"
// @Failure 500 {object} api.ResponseError "Internal server error"
// @Param usergroupID path int true "User group ID"
// @Router /usergroups/{usergroupID} [get]
// @Security Bearer
func getUserGroup(c *gin.Context) {
	ok, ug := database.CheckUserGroupPermissions(c, database.Read, "path", -1)
	if !ok {
		return
	}

	c.JSON(http.StatusOK, gin.H{"usergroup": ug})
}

// deleteUserGroup godoc
// @Summary Delete a user group
// @ID deleteUserGroup
// @Tags usergroups
// @Produce json
// @Success 200 {object} api.ResponseUserGroup "deleted user group"
// @Failure 404 {object} api.ResponseError "Not found"
// @Failure 422 {object} api.ResponseError "Unprocessable entity"
// @Failure 500 {object} api.ResponseError "Internal server error"
// @Param usergroupID path int true "User group ID"
// @Router /usergroups/{usergroupID} [delete]
// @Security Bearer
func deleteUserGroup(c *gin.Context) {

	ok, ug_r := database.CheckUserGroupPermissions(c, database.Delete, "path", -1)
	if !ok {
		return
	}

	var ug UserGroup
	ug.UserGroup = ug_r

	// Try to remove user group
	err := ug.remove()
	if !helper.DBError(c, err) {
		c.JSON(http.StatusOK, gin.H{"usergroup": ug})
	}

}
