package scenario

import (
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
// @Success 200 {array} common.ScenarioResponse "Array of scenarios to which user has access"
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

	serializer := common.ScenariosSerializer{c, scenarios}
	c.JSON(http.StatusOK, gin.H{
		"scenarios": serializer.Response(),
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

	var newScenarioData common.ResponseMsgScenario
	err = c.BindJSON(&newScenarioData)
	if err != nil {
		errormsg := "Bad request. Error binding form data to JSON: " + err.Error()
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errormsg,
		})
		return
	}

	var newScenario Scenario
	newScenario.ID = newScenarioData.Scenario.ID
	newScenario.StartParameters = newScenarioData.Scenario.StartParameters
	newScenario.Running = newScenarioData.Scenario.Running
	newScenario.Name = newScenarioData.Scenario.Name

	// save new scenario to DB
	err = newScenario.save()
	if common.ProvideErrorResponse(c, err) {
		return
	}

	// add user to new scenario
	err = newScenario.addUser(&(u.User))
	if common.ProvideErrorResponse(c, err) == false {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK.",
		})
	}
}

// updateScenario godoc
// @Summary Update a scenario
// @ID updateScenario
// @Tags scenarios
// @Accept json
// @Produce json
// @Param inputScenario body common.ResponseMsgScenario true "Scenario to be updated"
// @Success 200 "OK."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param scenarioID path int true "Scenario ID"
// @Router /scenarios/{scenarioID} [put]
func updateScenario(c *gin.Context) {

	ok, so := CheckPermissions(c, common.Update, "path", -1)
	if !ok {
		return
	}

	var modifiedScenarioData common.ResponseMsgScenario
	err := c.BindJSON(&modifiedScenarioData)
	if err != nil {
		errormsg := "Bad request. Error binding form data to JSON: " + err.Error()
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errormsg,
		})
		return
	}

	err = so.update(modifiedScenarioData.Scenario)
	if common.ProvideErrorResponse(c, err) == false {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK.",
		})
	}
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

	serializer := common.ScenarioSerializer{c, so.Scenario}
	c.JSON(http.StatusOK, gin.H{
		"scenario": serializer.Response(),
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
	if common.ProvideErrorResponse(c, err) == false {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK.",
		})
	}
}

// getUsersOfScenario godoc
// @Summary Get users of a scenario
// @ID getUsersOfScenario
// @Produce  json
// @Tags scenarios
// @Success 200 {array} common.UserResponse "Array of users that have access to the scenario"
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

	serializer := common.UsersSerializer{c, allUsers}
	c.JSON(http.StatusOK, gin.H{
		"users": serializer.Response(false),
	})
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
		"message": "OK.",
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

	err := so.deleteUser(username)
	if common.ProvideErrorResponse(c, err) {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OK.",
	})
}
