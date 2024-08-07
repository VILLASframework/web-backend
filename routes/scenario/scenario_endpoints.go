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

package scenario

import (
	"net/http"
	"strconv"

	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/helper"
	"github.com/gin-gonic/gin"
)

func RegisterScenarioEndpoints(r *gin.RouterGroup) {
	r.GET("", getScenarios)
	r.POST("", addScenario)
	r.PUT("/:scenarioID", updateScenario)
	r.GET("/:scenarioID", getScenario)
	r.DELETE("/:scenarioID", deleteScenario)
	r.GET("/:scenarioID/users", getUsersOfScenario)
	r.PUT("/:scenarioID/user", addUserToScenario)
	r.DELETE("/:scenarioID/user", deleteUserFromScenario)
}

// getScenarios godoc
// @Summary Get all scenarios of requesting user
// @ID getScenarios
// @Produce  json
// @Tags scenarios
// @Success 200 {object} api.ResponseScenarios "Scenarios to which user has access"
// @Failure 404 {object} api.ResponseError "Not found"
// @Failure 422 {object} api.ResponseError "Unprocessable entity"
// @Failure 500 {object} api.ResponseError "Internal server error"
// @Router /scenarios [get]
// @Security Bearer
func getScenarios(c *gin.Context) {

	// Checking permissions is not required here as read access is independent of user's role

	// ATTENTION: do not use c.GetInt (common.UserIDCtx) since userID is of type uint and not int
	userID, _ := c.Get(database.UserIDCtx)
	db := database.GetDB()
	var u database.User
	err := db.Find(&u, userID.(uint)).Error
	if helper.DBNotFoundError(c, err, strconv.FormatUint(uint64(userID.(uint)), 10), "User") {
		return
	}

	// get all scenarios for the user who issues the request

	var scenarios []database.Scenario
	if u.Role == "Admin" { // Admin can see all scenarios
		err = db.Order("ID asc").Find(&scenarios).Error
		if helper.DBError(c, err) {
			return
		}

	} else { // User or Guest roles see only their scenarios
		err = db.Order("ID asc").Model(&u).Related(&scenarios, "Scenarios").Error
		if helper.DBError(c, err) {
			return
		}
	}
	// TODO return list of configIDs, dashboardIDs and userIDs per scenario
	c.JSON(http.StatusOK, gin.H{"scenarios": scenarios})
}

// addScenario godoc
// @Summary Add a scenario
// @ID addScenario
// @Accept json
// @Produce json
// @Tags scenarios
// @Success 200 {object} api.ResponseScenario "scenario that was added"
// @Failure 400 {object} api.ResponseError "Bad request"
// @Failure 404 {object} api.ResponseError "Not found"
// @Failure 422 {object} api.ResponseError "Unprocessable entity"
// @Failure 500 {object} api.ResponseError "Internal server error"
// @Param inputScenario body scenario.addScenarioRequest true "Scenario to be added"
// @Router /scenarios [post]
// @Security Bearer
func addScenario(c *gin.Context) {

	ok, _ := database.CheckScenarioPermissions(c, database.Create, "none", -1)
	if !ok {
		return
	}

	// ATTENTION: do not use c.GetInt (common.UserIDCtx) since userID is of type uint and not int
	userID, _ := c.Get(database.UserIDCtx)
	db := database.GetDB()
	var u database.User
	err := db.Find(&u, userID.(uint)).Error
	if helper.DBNotFoundError(c, err, strconv.FormatUint(uint64(userID.(uint)), 10), "User") {
		return
	}

	var req addScenarioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.BadRequestError(c, err.Error())
		return
	}

	// Validate the request
	if err = req.validate(); err != nil {
		helper.UnprocessableEntityError(c, err.Error())
		return
	}

	// Create the new scenario from the request
	newScenario := req.createScenario()

	// Save the new scenario in the DB
	err = newScenario.save()
	if helper.DBError(c, err) {
		return
	}

	// add user to new scenario
	err = newScenario.addUser(&(u))
	if helper.DBError(c, err) {
		return
	}

	c.JSON(http.StatusOK, gin.H{"scenario": newScenario.Scenario})
}

// updateScenario godoc
// @Summary Update a scenario
// @ID updateScenario
// @Tags scenarios
// @Accept json
// @Produce json
// @Success 200 {object} api.ResponseScenario "Updated scenario."
// @Failure 400 {object} api.ResponseError "Bad request"
// @Failure 404 {object} api.ResponseError "Not found"
// @Failure 422 {object} api.ResponseError "Unprocessable entity"
// @Failure 500 {object} api.ResponseError "Internal server error"
// @Param inputScenario body scenario.updateScenarioRequest true "Scenario to be updated"
// @Param scenarioID path int true "Scenario ID"
// @Router /scenarios/{scenarioID} [put]
// @Security Bearer
func updateScenario(c *gin.Context) {

	ok, oldScenario_r := database.CheckScenarioPermissions(c, database.Update, "path", -1)
	if !ok {
		return
	}

	var oldScenario Scenario
	oldScenario.Scenario = oldScenario_r

	// Bind the (context) with the updateScenarioRequest struct
	var req updateScenarioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.BadRequestError(c, err.Error())
		return
	}

	// Validate the request based on struct updateScenarioRequest json tags
	if err := req.Scenario.validate(); err != nil {
		helper.BadRequestError(c, err.Error())
		return
	}

	// Create the updatedScenario from oldScenario
	userRole, _ := c.Get(database.UserRoleCtx)
	updatedScenario := req.updatedScenario(oldScenario, userRole.(string))

	// Finally update the scenario
	err := oldScenario.update(updatedScenario)
	if helper.DBError(c, err) {
		return
	}

	c.JSON(http.StatusOK, gin.H{"scenario": updatedScenario.Scenario})
}

// getScenario godoc
// @Summary Get scenario
// @ID getScenario
// @Produce  json
// @Tags scenarios
// @Success 200 {object} api.ResponseScenario "Scenario requested by user"
// @Failure 404 {object} api.ResponseError "Not found"
// @Failure 422 {object} api.ResponseError "Unprocessable entity"
// @Failure 500 {object} api.ResponseError "Internal server error"
// @Param scenarioID path int true "Scenario ID"
// @Router /scenarios/{scenarioID} [get]
// @Security Bearer
func getScenario(c *gin.Context) {

	ok, so := database.CheckScenarioPermissions(c, database.Read, "path", -1)
	if !ok {
		return
	}

	// TODO return list of configIDs, dashboardIDs and userIDs per scenario
	c.JSON(http.StatusOK, gin.H{"scenario": so})
}

// deleteScenario godoc
// @Summary Delete a scenario
// @ID deleteScenario
// @Tags scenarios
// @Produce json
// @Success 200 {object} api.ResponseScenario  "Deleted scenario"
// @Failure 404 {object} api.ResponseError "Not found"
// @Failure 422 {object} api.ResponseError "Unprocessable entity"
// @Failure 500 {object} api.ResponseError "Internal server error"
// @Param scenarioID path int true "Scenario ID"
// @Router /scenarios/{scenarioID} [delete]
// @Security Bearer
func deleteScenario(c *gin.Context) {

	ok, so_r := database.CheckScenarioPermissions(c, database.Delete, "path", -1)
	if !ok {
		return
	}

	var so Scenario
	so.Scenario = so_r

	errs := so.delete()
	if len(errs) > 0 {
		var errorString = "DB errors:"
		for _, e := range errs {
			if e != nil {
				errorString += ", " + e.Error()
			}
		}
		helper.InternalServerError(c, errorString)
		return
	}

	c.JSON(http.StatusOK, gin.H{"scenario": so.Scenario})
}

// getUsersOfScenario godoc
// @Summary Get users of a scenario
// @ID getUsersOfScenario
// @Produce  json
// @Tags scenarios
// @Success 200 {object} api.ResponseUsers "Array of users that have access to the scenario"
// @Failure 404 {object} api.ResponseError "Not found"
// @Failure 422 {object} api.ResponseError "Unprocessable entity"
// @Failure 500 {object} api.ResponseError "Internal server error"
// @Param scenarioID path int true "Scenario ID"
// @Router /scenarios/{scenarioID}/users/ [get]
// @Security Bearer
func getUsersOfScenario(c *gin.Context) {

	ok, so_r := database.CheckScenarioPermissions(c, database.Read, "path", -1)
	if !ok {
		return
	}

	var so Scenario
	so.Scenario = so_r

	// Find all users of scenario
	allUsers, _, err := so.getUsers()
	if helper.DBError(c, err) {
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": allUsers})
}

// addUserToScenario godoc
// @Summary Add a user to a a scenario
// @ID addUserToScenario
// @Tags scenarios
// @Produce json
// @Success 200 {object} api.ResponseUser "User that was added to scenario"
// @Failure 404 {object} api.ResponseError "Not found"
// @Failure 422 {object} api.ResponseError "Unprocessable entity"
// @Failure 500 {object} api.ResponseError "Internal server error"
// @Param scenarioID path int true "Scenario ID"
// @Param username query string true "User name"
// @Router /scenarios/{scenarioID}/user [put]
// @Security Bearer
func addUserToScenario(c *gin.Context) {

	ok, so_r := database.CheckScenarioPermissions(c, database.Update, "path", -1)
	if !ok {
		return
	}

	var so Scenario
	so.Scenario = so_r

	username := c.Request.URL.Query().Get("username")
	var u database.User
	db := database.GetDB()
	err := db.Find(&u, "Username = ?", username).Error
	if helper.DBNotFoundError(c, err, username, "User") {
		return
	}

	if !u.Active {
		helper.BadRequestError(c, "bad user")
		return
	}

	err = so.addUser(&(u))
	if helper.DBError(c, err) {
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": u})
}

// deleteUserFromScenario godoc
// @Summary Delete a user from a scenario
// @ID deleteUserFromScenario
// @Tags scenarios
// @Produce json
// @Success 200 {object} api.ResponseUser "User that was deleted from scenario"
// @Failure 404 {object} api.ResponseError "Not found"
// @Failure 422 {object} api.ResponseError "Unprocessable entity"
// @Failure 500 {object} api.ResponseError "Internal server error"
// @Param scenarioID path int true "Scenario ID"
// @Param username query string true "User name"
// @Router /scenarios/{scenarioID}/user [delete]
// @Security Bearer
func deleteUserFromScenario(c *gin.Context) {

	ok, so_r := database.CheckScenarioPermissions(c, database.Update, "path", -1)
	if !ok {
		return
	}

	var so Scenario
	so.Scenario = so_r

	username := c.Request.URL.Query().Get("username")
	var u database.User
	db := database.GetDB()
	err := db.Find(&u, "Username = ?", username).Error
	if helper.DBNotFoundError(c, err, username, "User") {
		return
	}

	err = so.deleteUser(username)
	if helper.DBError(c, err) {
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": u})
}
