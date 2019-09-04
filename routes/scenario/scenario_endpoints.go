package scenario

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/user"
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
// @Summary Get all scenarios
// @ID getScenarios
// @Produce  json
// @Tags scenarios
// @Success 200 {array} docs.ResponseScenarios "Array of scenarios to which user has access"
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Router /scenarios [get]
func getScenarios(c *gin.Context) {

	ok, _ := CheckPermissions(c, common.Read, "none", -1)
	if !ok {
		return
	}

	// ATTENTION: do not use c.GetInt (common.UserIDCtx) since user_id is of type uint and not int
	userID, _ := c.Get(common.UserIDCtx)
	userRole, _ := c.Get(common.UserRoleCtx)

	var u user.User
	err := u.ByID(userID.(uint))
	if common.ProvideErrorResponse(c, err) {
		return
	}

	// get all scenarios for the user who issues the request
	db := common.GetDB()
	var scenarios []common.Scenario
	if userRole == "Admin" { // Admin can see all scenarios
		err = db.Order("ID asc").Find(&scenarios).Error
		if common.ProvideErrorResponse(c, err) {
			return
		}

	} else { // User or Guest roles see only their scenarios
		err = db.Order("ID asc").Model(&u).Related(&scenarios, "Scenarios").Error
		if common.ProvideErrorResponse(c, err) {
			return
		}
	}
	// TODO return list of simulationModelIDs, dashboardIDs and userIDs per scenario
	c.JSON(http.StatusOK, gin.H{
		"scenarios": scenarios,
	})
}

// addScenario godoc
// @Summary Add a scenario
// @ID addScenario
// @Accept json
// @Produce json
// @Tags scenarios
// @Param inputScenario body common.ResponseMsgScenario true "Scenario to be added"
// @Success 200 "OK."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Router /scenarios [post]
func addScenario(c *gin.Context) {

	ok, _ := CheckPermissions(c, common.Create, "none", -1)
	if !ok {
		return
	}

	userID, _ := c.Get(common.UserIDCtx)

	var u user.User
	err := u.ByID(userID.(uint))
	if common.ProvideErrorResponse(c, err) {
		return
	}

	var req addScenarioRequest
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

	// Create the new scenario from the request
	newScenario := req.createScenario()

	// Save the new scenario in the DB
	err = newScenario.save()
	if err != nil {
		common.ProvideErrorResponse(c, err)
		return
	}

	// add user to new scenario
	err = newScenario.addUser(&(u.User))
	if err != nil {
		common.ProvideErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"scenario": newScenario.Scenario,
	})
}

// updateScenario godoc
// @Summary Update a scenario
// @ID updateScenario
// @Tags scenarios
// @Accept json
// @Produce json
// @Param inputScenario body common.Scenario true "Scenario to be updated"
// @Success 200 "OK."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param scenarioID path int true "Scenario ID"
// @Router /scenarios/{scenarioID} [put]
func updateScenario(c *gin.Context) {

	ok, oldScenario := CheckPermissions(c, common.Update, "path", -1)
	if !ok {
		return
	}

	// Bind the (context) with the updateScenarioRequest struct
	var req updateScenarioRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": fmt.Sprintf("%v", err),
		})
		return
	}

	// Validate the request based on struct updateScenarioRequest json tags
	if err := req.validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": fmt.Sprintf("%v", err),
		})
		return
	}

	// Create the updatedScenario from oldScenario
	updatedScenario, err := req.updatedScenario(oldScenario)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": fmt.Sprintf("%v", err),
		})
		return
	}

	// Finally update the scenario
	err = oldScenario.update(updatedScenario)
	if err != nil {
		common.ProvideErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"scenario": updatedScenario.Scenario,
	})
}

// getScenario godoc
// @Summary Get scenario
// @ID getScenario
// @Produce  json
// @Tags scenarios
// @Success 200 {object} common.ScenarioResponse "Scenario requested by user"
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param scenarioID path int true "Scenario ID"
// @Router /scenarios/{scenarioID} [get]
func getScenario(c *gin.Context) {

	ok, so := CheckPermissions(c, common.Read, "path", -1)
	if !ok {
		return
	}

	// TODO return list of simulationModelIDs, dashboardIDs and userIDs per scenario
	c.JSON(http.StatusOK, gin.H{
		"scenario": so.Scenario,
	})
}

// deleteScenario godoc
// @Summary Delete a scenario
// @ID deleteScenario
// @Tags scenarios
// @Produce json
// @Success 200 "OK."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param scenarioID path int true "Scenario ID"
// @Router /scenarios/{scenarioID} [delete]
func deleteScenario(c *gin.Context) {

	ok, so := CheckPermissions(c, common.Delete, "path", -1)
	if !ok {
		return
	}

	err := so.delete()
	if err != nil {
		common.ProvideErrorResponse(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"scenario": so.Scenario,
	})
}

// getUsersOfScenario godoc
// @Summary Get users of a scenario
// @ID getUsersOfScenario
// @Produce  json
// @Tags scenarios
// @Success 200 {object} docs.ResponseUsers "Array of users that have access to the scenario"
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param scenarioID path int true "Scenario ID"
// @Router /scenarios/{scenarioID}/users/ [get]
func getUsersOfScenario(c *gin.Context) {

	ok, so := CheckPermissions(c, common.Read, "path", -1)
	if !ok {
		return
	}

	// Find all users of scenario
	allUsers, _, err := so.getUsers()
	if common.ProvideErrorResponse(c, err) {
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": allUsers})
}

// addUserToScenario godoc
// @Summary Add a user to a a scenario
// @ID addUserToScenario
// @Tags scenarios
// @Produce json
// @Success 200 "OK."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param scenarioID path int true "Scenario ID"
// @Param username query string true "User name"
// @Router /scenarios/{scenarioID}/user [put]
func addUserToScenario(c *gin.Context) {

	ok, so := CheckPermissions(c, common.Update, "path", -1)
	if !ok {
		return
	}

	username := c.Request.URL.Query().Get("username")

	var u user.User
	err := u.ByUsername(username)
	if common.ProvideErrorResponse(c, err) {
		return
	}

	err = so.addUser(&(u.User))
	if common.ProvideErrorResponse(c, err) {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": u.User,
	})
}

// deleteUserFromScenario godoc
// @Summary Delete a user from a scenario
// @ID deleteUserFromScenario
// @Tags scenarios
// @Produce json
// @Success 200 "OK."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param scenarioID path int true "Scenario ID"
// @Param username query string true "User name"
// @Router /scenarios/{scenarioID}/user [delete]
func deleteUserFromScenario(c *gin.Context) {

	ok, so := CheckPermissions(c, common.Update, "path", -1)
	if !ok {
		return
	}

	username := c.Request.URL.Query().Get("username")

	var u user.User
	err := u.ByUsername(username)
	if common.ProvideErrorResponse(c, err) {
		return
	}

	err = so.deleteUser(username)
	if common.ProvideErrorResponse(c, err) {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": u.User,
	})
}
