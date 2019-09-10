package dashboard

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/scenario"
)

func RegisterDashboardEndpoints(r *gin.RouterGroup) {

	r.GET("", getDashboards)
	r.POST("", addDashboard)
	r.PUT("/:dashboardID", updateDashboard)
	r.GET("/:dashboardID", getDashboard)
	r.DELETE("/:dashboardID", deleteDashboard)
}

// getDashboards godoc
// @Summary Get all dashboards of scenario
// @ID getDashboards
// @Produce  json
// @Tags dashboards
// @Success 200 {object} docs.ResponseDashboards "Dashboards which belong to scenario"
// @Failure 404 {object} docs.ResponseError "Not found"
// @Failure 422 {object} docs.ResponseError "Unprocessable entity"
// @Failure 500 {object} docs.ResponseError "Internal server error"
// @Param scenarioID query int true "Scenario ID"
// @Router /dashboards [get]
func getDashboards(c *gin.Context) {

	ok, sim := scenario.CheckPermissions(c, common.Read, "query", -1)
	if !ok {
		return
	}

	db := common.GetDB()
	var dab []common.Dashboard
	err := db.Order("ID asc").Model(sim).Related(&dab, "Dashboards").Error
	if common.DBError(c, err) {
		return
	}

	c.JSON(http.StatusOK, gin.H{"dashboards": dab})
}

// addDashboard godoc
// @Summary Add a dashboard to a scenario
// @ID addDashboard
// @Accept json
// @Produce json
// @Tags dashboards
// @Param inputDab body dashboard.validNewDashboard true "Dashboard to be added incl. ID of Scenario"
// @Success 200 {object} docs.ResponseDashboard "Dashboard that was added"
// @Failure 400 {object} docs.ResponseError "Bad request"
// @Failure 404 {object} docs.ResponseError "Not found"
// @Failure 422 {object} docs.ResponseError "Unprocessable entity"
// @Failure 500 {object} docs.ResponseError "Internal server error"
// @Router /dashboards [post]
func addDashboard(c *gin.Context) {

	// bind request to JSON
	var req addDashboardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.BadRequestError(c, err.Error())
		return
	}

	// Validate the request
	if err := req.validate(); err != nil {
		common.UnprocessableEntityError(c, err.Error())
		return
	}

	// Create the new dashboard from the request
	newDashboard := req.createDashboard()

	// Check if user is allowed to modify scenario specified in request
	ok, _ := scenario.CheckPermissions(c, common.Update, "body", int(newDashboard.ScenarioID))
	if !ok {
		return
	}

	// add dashboard to DB and add association to scenario
	err := newDashboard.addToScenario()
	if common.DBError(c, err) {
		return
	}

	c.JSON(http.StatusOK, gin.H{"dashboard": newDashboard.Dashboard})
}

// updateDashboard godoc
// @Summary Update a dashboard
// @ID updateDashboard
// @Tags dashboards
// @Accept json
// @Produce json
// @Param inputDab body dashboard.validUpdatedDashboard true "Dashboard to be updated"
// @Success 200 {object} docs.ResponseDashboard "Dashboard that was updated"
// @Failure 400 {object} docs.ResponseError "Bad request"
// @Failure 404 {object} docs.ResponseError "Not found"
// @Failure 422 {object} docs.ResponseError "Unprocessable entity"
// @Failure 500 {object} docs.ResponseError "Internal server error"
// @Param dashboardID path int true "Dashboard ID"
// @Router /dashboards/{dashboardID} [put]
func updateDashboard(c *gin.Context) {

	ok, oldDashboard := CheckPermissions(c, common.Update, "path", -1)
	if !ok {
		return
	}

	var req updateDashboardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.BadRequestError(c, err.Error())
		return
	}

	// Validate the request
	if err := req.validate(); err != nil {
		common.BadRequestError(c, err.Error())
		return
	}
	// Create the updatedDashboard from oldDashboard
	updatedDashboard := req.updatedDashboard(oldDashboard)

	// update the dashboard in the DB
	err := oldDashboard.update(updatedDashboard)
	if common.DBError(c, err) {
		return
	}

	c.JSON(http.StatusOK, gin.H{"dashboard": updatedDashboard.Dashboard})
}

// getDashboard godoc
// @Summary Get a dashboard
// @ID getDashboard
// @Tags dashboards
// @Produce json
// @Success 200 {object} docs.ResponseDashboard "Dashboard that was requested"
// @Failure 400 {object} docs.ResponseError "Bad request"
// @Failure 404 {object} docs.ResponseError "Not found"
// @Failure 422 {object} docs.ResponseError "Unprocessable entity"
// @Failure 500 {object} docs.ResponseError "Internal server error"
// @Param dashboardID path int true "Dashboard ID"
// @Router /dashboards/{dashboardID} [get]
func getDashboard(c *gin.Context) {

	ok, dab := CheckPermissions(c, common.Read, "path", -1)
	if !ok {
		return
	}

	c.JSON(http.StatusOK, gin.H{"dashboard": dab.Dashboard})
}

// deleteDashboard godoc
// @Summary Delete a dashboard
// @ID deleteDashboard
// @Tags dashboards
// @Produce json
// @Success 200 {object} docs.ResponseDashboard "Dashboard that was deleted"
// @Failure 400 {object} docs.ResponseError "Bad request"
// @Failure 404 {object} docs.ResponseError "Not found"
// @Failure 422 {object} docs.ResponseError "Unprocessable entity"
// @Failure 500 {object} docs.ResponseError "Internal server error"
// @Param dashboardID path int true "Dashboard ID"
// @Router /dashboards/{dashboardID} [delete]
func deleteDashboard(c *gin.Context) {
	ok, dab := CheckPermissions(c, common.Delete, "path", -1)
	if !ok {
		return
	}

	err := dab.delete()
	if common.DBError(c, err) {
		return
	}

	c.JSON(http.StatusOK, gin.H{"dashboard": dab.Dashboard})
}
