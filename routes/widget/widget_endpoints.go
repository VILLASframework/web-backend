/** Widget package, endpoints.
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
package widget

import (
	"net/http"

	"git.rwth-aachen.de/acs/public/villas/web-backend-go/helper"

	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
)

func RegisterWidgetEndpoints(r *gin.RouterGroup) {
	r.GET("", getWidgets)
	r.POST("", addWidget)
	r.PUT("/:widgetID", updateWidget)
	r.GET("/:widgetID", getWidget)
	r.DELETE("/:widgetID", deleteWidget)
}

// getWidgets godoc
// @Summary Get all widgets of dashboard
// @ID getWidgets
// @Produce  json
// @Tags widgets
// @Success 200 {object} api.ResponseWidgets "Widgets to which belong to dashboard"
// @Failure 404 {object} api.ResponseError "Not found"
// @Failure 422 {object} api.ResponseError "Unprocessable entity"
// @Failure 500 {object} api.ResponseError "Internal server error"
// @Param dashboardID query int true "Dashboard ID"
// @Router /widgets [get]
// @Security Bearer
func getWidgets(c *gin.Context) {

	ok, dab := database.CheckDashboardPermissions(c, database.Read, "query", -1)
	if !ok {
		return
	}

	db := database.GetDB()
	var widgets []database.Widget
	err := db.Order("ID asc").Model(dab).Related(&widgets, "Widgets").Error
	if !helper.DBError(c, err) {
		c.JSON(http.StatusOK, gin.H{"widgets": widgets})
	}

}

// addWidget godoc
// @Summary Add a widget to a dashboard
// @ID addWidget
// @Accept json
// @Produce json
// @Tags widgets
// @Success 200 {object} api.ResponseWidget "Widget that was added"
// @Failure 400 {object} api.ResponseError "Bad request"
// @Failure 404 {object} api.ResponseError "Not found"
// @Failure 422 {object} api.ResponseError "Unprocessable entity"
// @Failure 500 {object} api.ResponseError "Internal server error"
// @Param inputWidget body widget.addWidgetRequest true "Widget to be added incl. ID of dashboard"
// @Router /widgets [post]
// @Security Bearer
func addWidget(c *gin.Context) {

	var req addWidgetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.BadRequestError(c, err.Error())
		return
	}

	// Validate the request
	if err := req.validate(); err != nil {
		helper.UnprocessableEntityError(c, err.Error())
		return
	}

	// Create the new widget from the request
	newWidget := req.createWidget()

	// Check if user is allowed to modify selected dashboard (scenario)
	ok, _ := database.CheckDashboardPermissions(c, database.Update, "body", int(newWidget.DashboardID))
	if !ok {
		return
	}

	err := newWidget.addToDashboard()
	if !helper.DBError(c, err) {
		c.JSON(http.StatusOK, gin.H{"widget": newWidget.Widget})
	}

}

// updateWidget godoc
// @Summary Update a widget
// @ID updateWidget
// @Tags widgets
// @Accept json
// @Produce json
// @Success 200 {object} api.ResponseWidget "Widget that was updated"
// @Failure 400 {object} api.ResponseError "Bad request"
// @Failure 404 {object} api.ResponseError "Not found"
// @Failure 422 {object} api.ResponseError "Unprocessable entity"
// @Failure 500 {object} api.ResponseError "Internal server error"
// @Param inputWidget body widget.updateWidgetRequest true "Widget to be updated"
// @Param widgetID path int true "Widget ID"
// @Router /widgets/{widgetID} [put]
// @Security Bearer
func updateWidget(c *gin.Context) {

	ok, oldWidget_r := database.CheckWidgetPermissions(c, database.Update, -1)
	if !ok {
		return
	}

	var oldWidget Widget
	oldWidget.Widget = oldWidget_r

	var req updateWidgetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.BadRequestError(c, err.Error())
		return
	}

	// Validate the request
	if err := req.Widget.validate(); err != nil {
		helper.BadRequestError(c, err.Error())
		return
	}

	// Create the updatedScenario from oldScenario
	updatedWidget := req.updatedWidget(oldWidget)

	// Update the widget in the DB
	err := oldWidget.update(updatedWidget)
	if !helper.DBError(c, err) {
		c.JSON(http.StatusOK, gin.H{"widget": updatedWidget.Widget})
	}

}

// getWidget godoc
// @Summary Get a widget
// @ID getWidget
// @Tags widgets
// @Produce json
// @Success 200 {object} api.ResponseWidget "Widget that was requested"
// @Failure 400 {object} api.ResponseError "Bad request"
// @Failure 404 {object} api.ResponseError "Not found"
// @Failure 422 {object} api.ResponseError "Unprocessable entity"
// @Failure 500 {object} api.ResponseError "Internal server error"
// @Param widgetID path int true "Widget ID"
// @Router /widgets/{widgetID} [get]
// @Security Bearer
func getWidget(c *gin.Context) {

	ok, w := database.CheckWidgetPermissions(c, database.Read, -1)
	if !ok {
		return
	}

	c.JSON(http.StatusOK, gin.H{"widget": w})
}

// deleteWidget godoc
// @Summary Delete a widget
// @ID deleteWidget
// @Tags widgets
// @Produce json
// @Success 200 {object} api.ResponseWidget "Widget that was deleted"
// @Failure 400 {object} api.ResponseError "Bad request"
// @Failure 404 {object} api.ResponseError "Not found"
// @Failure 422 {object} api.ResponseError "Unprocessable entity"
// @Failure 500 {object} api.ResponseError "Internal server error"
// @Param widgetID path int true "Widget ID"
// @Router /widgets/{widgetID} [delete]
// @Security Bearer
func deleteWidget(c *gin.Context) {

	ok, w_r := database.CheckWidgetPermissions(c, database.Delete, -1)
	if !ok {
		return
	}

	var w Widget
	w.Widget = w_r

	err := w.delete()
	if !helper.DBError(c, err) {
		c.JSON(http.StatusOK, gin.H{"widget": w.Widget})
	}

}
