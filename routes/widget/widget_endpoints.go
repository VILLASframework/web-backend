package widget

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/dashboard"
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
// @Success 200 {object} docs.ResponseWidgets "Widgets to which belong to dashboard"
// @Failure 404 {object} docs.ResponseError "Not found"
// @Failure 422 {object} docs.ResponseError "Unprocessable entity"
// @Failure 500 {object} docs.ResponseError "Internal server error"
// @Param dashboardID query int true "Dashboard ID"
// @Router /widgets [get]
func getWidgets(c *gin.Context) {

	ok, dab := dashboard.CheckPermissions(c, common.Read, "query", -1)
	if !ok {
		return
	}

	db := common.GetDB()
	var widgets []common.Widget
	err := db.Order("ID asc").Model(dab).Related(&widgets, "Widgets").Error
	if common.DBError(c, err) {
		return
	}

	c.JSON(http.StatusOK, gin.H{"widgets": widgets})
}

// addWidget godoc
// @Summary Add a widget to a dashboard
// @ID addWidget
// @Accept json
// @Produce json
// @Tags widgets
// @Param inputWidget body widget.validNewWidget true "Widget to be added incl. ID of dashboard"
// @Success 200 {object} docs.ResponseWidget "Widget that was added"
// @Failure 400 {object} docs.ResponseError "Bad request"
// @Failure 404 {object} docs.ResponseError "Not found"
// @Failure 422 {object} docs.ResponseError "Unprocessable entity"
// @Failure 500 {object} docs.ResponseError "Internal server error"
// @Router /widgets [post]
func addWidget(c *gin.Context) {

	var req addWidgetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.BadRequestError(c, err.Error())
		return
	}

	// Validate the request
	if err := req.validate(); err != nil {
		common.UnprocessableEntityError(c, err.Error())
		return
	}

	// Create the new widget from the request
	newWidget := req.createWidget()

	// Check if user is allowed to modify selected dashboard (scenario)
	ok, _ := dashboard.CheckPermissions(c, common.Create, "body", int(newWidget.DashboardID))
	if !ok {
		return
	}

	err := newWidget.addToDashboard()
	if err != nil {
		common.DBError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"widget": newWidget.Widget})
}

// updateWidget godoc
// @Summary Update a widget
// @ID updateWidget
// @Tags widgets
// @Accept json
// @Produce json
// @Param inputWidget body widget.validUpdatedWidget true "Widget to be updated"
// @Success 200 {object} docs.ResponseWidget "Widget that was updated"
// @Failure 400 {object} docs.ResponseError "Bad request"
// @Failure 404 {object} docs.ResponseError "Not found"
// @Failure 422 {object} docs.ResponseError "Unprocessable entity"
// @Failure 500 {object} docs.ResponseError "Internal server error"
// @Param widgetID path int true "Widget ID"
// @Router /widgets/{widgetID} [put]
func updateWidget(c *gin.Context) {

	ok, oldWidget := CheckPermissions(c, common.Update, -1)
	if !ok {
		return
	}

	var req updateWidgetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.BadRequestError(c, err.Error())
		return
	}

	// Validate the request
	if err := req.validate(); err != nil {
		common.BadRequestError(c, err.Error())
		return
	}

	// Create the updatedScenario from oldScenario
	updatedWidget, err := req.updatedWidget(oldWidget)
	if err != nil {
		common.BadRequestError(c, err.Error())
		return
	}

	// Update the widget in the DB
	err = oldWidget.update(updatedWidget)
	if err != nil {
		common.DBError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"widget": updatedWidget.Widget})

}

// getWidget godoc
// @Summary Get a widget
// @ID getWidget
// @Tags widgets
// @Produce json
// @Success 200 {object} docs.ResponseWidget "Widget that was requested"
// @Failure 400 {object} docs.ResponseError "Bad request"
// @Failure 404 {object} docs.ResponseError "Not found"
// @Failure 422 {object} docs.ResponseError "Unprocessable entity"
// @Failure 500 {object} docs.ResponseError "Internal server error"
// @Param widgetID path int true "Widget ID"
// @Router /widgets/{widgetID} [get]
func getWidget(c *gin.Context) {

	ok, w := CheckPermissions(c, common.Read, -1)
	if !ok {
		return
	}

	c.JSON(http.StatusOK, gin.H{"widget": w.Widget})
}

// deleteWidget godoc
// @Summary Delete a widget
// @ID deleteWidget
// @Tags widgets
// @Produce json
// @Success 200 {object} docs.ResponseWidget "Widget that was deleted"
// @Failure 400 {object} docs.ResponseError "Bad request"
// @Failure 404 {object} docs.ResponseError "Not found"
// @Failure 422 {object} docs.ResponseError "Unprocessable entity"
// @Failure 500 {object} docs.ResponseError "Internal server error"
// @Param widgetID path int true "Widget ID"
// @Router /widgets/{widgetID} [delete]
func deleteWidget(c *gin.Context) {

	ok, w := CheckPermissions(c, common.Delete, -1)
	if !ok {
		return
	}

	err := w.delete()
	if common.DBError(c, err) {
		return
	}

	c.JSON(http.StatusOK, gin.H{"widget": w.Widget})
}
