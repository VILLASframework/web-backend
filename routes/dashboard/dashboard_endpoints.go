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

package dashboard

import (
	"net/http"

	"git.rwth-aachen.de/acs/public/villas/web-backend-go/helper"

	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
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
// @Success 200 {object} api.ResponseDashboards "Dashboards which belong to scenario"
// @Failure 404 {object} api.ResponseError "Not found"
// @Failure 422 {object} api.ResponseError "Unprocessable entity"
// @Failure 500 {object} api.ResponseError "Internal server error"
// @Param scenarioID query int true "Scenario ID"
// @Router /dashboards [get]
// @Security Bearer
func getDashboards(c *gin.Context) {

	ok, sim := database.CheckScenarioPermissions(c, database.Read, "query", -1)
	if !ok {
		return
	}

	db := database.GetDB()
	var dab []database.Dashboard
	err := db.Order("ID asc").Model(sim).Related(&dab, "Dashboards").Error
	if !helper.DBError(c, err) {
		c.JSON(http.StatusOK, gin.H{"dashboards": dab})
	}

}

// addDashboard godoc
// @Summary Add a dashboard to a scenario
// @ID addDashboard
// @Accept json
// @Produce json
// @Tags dashboards
// @Success 200 {object} api.ResponseDashboard "Dashboard that was added"
// @Failure 400 {object} api.ResponseError "Bad request"
// @Failure 404 {object} api.ResponseError "Not found"
// @Failure 422 {object} api.ResponseError "Unprocessable entity"
// @Failure 500 {object} api.ResponseError "Internal server error"
// @Param inputDab body dashboard.addDashboardRequest true "Dashboard to be added incl. ID of Scenario"
// @Router /dashboards [post]
// @Security Bearer
func addDashboard(c *gin.Context) {

	// bind request to JSON
	var req addDashboardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.BadRequestError(c, err.Error())
		return
	}

	// Validate the request
	if err := req.validate(); err != nil {
		helper.UnprocessableEntityError(c, err.Error())
		return
	}

	// Create the new dashboard from the request
	newDashboard := req.createDashboard()

	// Check if user is allowed to modify scenario specified in request
	ok, _ := database.CheckScenarioPermissions(c, database.Update, "body", int(newDashboard.ScenarioID))
	if !ok {
		return
	}

	// add dashboard to DB and add association to scenario
	err := newDashboard.addToScenario()
	if !helper.DBError(c, err) {
		c.JSON(http.StatusOK, gin.H{"dashboard": newDashboard.Dashboard})
	}

}

// updateDashboard godoc
// @Summary Update a dashboard
// @ID updateDashboard
// @Tags dashboards
// @Accept json
// @Produce json
// @Success 200 {object} api.ResponseDashboard "Dashboard that was updated"
// @Failure 400 {object} api.ResponseError "Bad request"
// @Failure 404 {object} api.ResponseError "Not found"
// @Failure 422 {object} api.ResponseError "Unprocessable entity"
// @Failure 500 {object} api.ResponseError "Internal server error"
// @Param inputDab body dashboard.updateDashboardRequest true "Dashboard to be updated"
// @Param dashboardID path int true "Dashboard ID"
// @Router /dashboards/{dashboardID} [put]
// @Security Bearer
func updateDashboard(c *gin.Context) {

	ok, oldDashboard_r := database.CheckDashboardPermissions(c, database.Update, "path", -1)
	if !ok {
		return
	}

	var oldDashboard Dashboard
	oldDashboard.Dashboard = oldDashboard_r

	var req updateDashboardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.BadRequestError(c, err.Error())
		return
	}

	// Validate the request
	if err := req.Dashboard.validate(); err != nil {
		helper.BadRequestError(c, err.Error())
		return
	}
	// Create the updatedDashboard from oldDashboard
	updatedDashboard := req.updatedDashboard(oldDashboard)

	// update the dashboard in the DB
	err := oldDashboard.update(updatedDashboard)
	if !helper.DBError(c, err) {
		c.JSON(http.StatusOK, gin.H{"dashboard": updatedDashboard.Dashboard})
	}

}

// getDashboard godoc
// @Summary Get a dashboard
// @ID getDashboard
// @Tags dashboards
// @Produce json
// @Success 200 {object} api.ResponseDashboard "Dashboard that was requested"
// @Failure 400 {object} api.ResponseError "Bad request"
// @Failure 404 {object} api.ResponseError "Not found"
// @Failure 422 {object} api.ResponseError "Unprocessable entity"
// @Failure 500 {object} api.ResponseError "Internal server error"
// @Param dashboardID path int true "Dashboard ID"
// @Router /dashboards/{dashboardID} [get]
// @Security Bearer
func getDashboard(c *gin.Context) {

	ok, dab := database.CheckDashboardPermissions(c, database.Read, "path", -1)
	if !ok {
		return
	}

	c.JSON(http.StatusOK, gin.H{"dashboard": dab})
}

// deleteDashboard godoc
// @Summary Delete a dashboard
// @ID deleteDashboard
// @Tags dashboards
// @Produce json
// @Success 200 {object} api.ResponseDashboard "Dashboard that was deleted"
// @Failure 400 {object} api.ResponseError "Bad request"
// @Failure 404 {object} api.ResponseError "Not found"
// @Failure 422 {object} api.ResponseError "Unprocessable entity"
// @Failure 500 {object} api.ResponseError "Internal server error"
// @Param dashboardID path int true "Dashboard ID"
// @Router /dashboards/{dashboardID} [delete]
// @Security Bearer
func deleteDashboard(c *gin.Context) {
	ok, dab_r := database.CheckDashboardPermissions(c, database.Delete, "path", -1)
	if !ok {
		return
	}

	var dab Dashboard
	dab.Dashboard = dab_r

	err := dab.delete()
	if !helper.DBError(c, err) {
		c.JSON(http.StatusOK, gin.H{"dashboard": dab.Dashboard})
	}

}
