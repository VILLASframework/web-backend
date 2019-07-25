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
// @Success 200 {array} common.WidgetResponse "Array of widgets to which belong to dashboard"
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
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
	if common.ProvideErrorResponse(c, err) {
		return
	}

	serializer := common.WidgetsSerializer{c, widgets}
	c.JSON(http.StatusOK, gin.H{
		"widgets": serializer.Response(),
	})
}

// addWidget godoc
// @Summary Add a widget to a dashboard
// @ID addWidget
// @Accept json
// @Produce json
// @Tags widgets
// @Param inputWidget body common.ResponseMsgWidget true "Widget to be added incl. ID of dashboard"
// @Success 200 "OK."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Router /widgets [post]
func addWidget(c *gin.Context) {

	var newWidgetData common.ResponseMsgWidget
	err := c.BindJSON(&newWidgetData)
	if err != nil {
		errormsg := "Bad request. Error binding form data to JSON: " + err.Error()
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errormsg,
		})
		return
	}

	var newWidget Widget
	newWidget.Name = newWidgetData.Widget.Name
	newWidget.Type = newWidgetData.Widget.Type
	newWidget.Height = newWidgetData.Widget.Height
	newWidget.Width = newWidgetData.Widget.Width
	newWidget.MinHeight = newWidgetData.Widget.MinHeight
	newWidget.MinWidth = newWidgetData.Widget.MinWidth
	newWidget.X = newWidgetData.Widget.X
	newWidget.Y = newWidgetData.Widget.Y
	newWidget.Z = newWidgetData.Widget.Z
	newWidget.CustomProperties = newWidgetData.Widget.CustomProperties
	newWidget.IsLocked = newWidgetData.Widget.IsLocked
	newWidget.DashboardID = newWidgetData.Widget.DashboardID

	ok, _ := dashboard.CheckPermissions(c, common.Create, "body", int(newWidget.DashboardID))
	if !ok {
		return
	}

	err = newWidget.addToDashboard()

	if common.ProvideErrorResponse(c, err) == false {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK.",
		})
	}
}

// updateWidget godoc
// @Summary Update a widget
// @ID updateWidget
// @Tags widgets
// @Accept json
// @Produce json
// @Param inputWidget body common.ResponseMsgWidget true "Widget to be updated"
// @Success 200 "OK."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param widgetID path int true "Widget ID"
// @Router /widgets/{widgetID} [put]
func updateWidget(c *gin.Context) {

	ok, w := CheckPermissions(c, common.Update, -1)
	if !ok {
		return
	}

	var modifiedWidget common.ResponseMsgWidget
	err := c.BindJSON(&modifiedWidget)
	if err != nil {
		errormsg := "Bad request. Error binding form data to JSON: " + err.Error()
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errormsg,
		})
		return
	}

	err = w.update(modifiedWidget.Widget)
	if common.ProvideErrorResponse(c, err) == false {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK.",
		})
	}
}

// getWidget godoc
// @Summary Get a widget
// @ID getWidget
// @Tags widgets
// @Produce json
// @Success 200 {object} common.WidgetResponse "Requested widget."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param widgetID path int true "Widget ID"
// @Router /widgets/{widgetID} [get]
func getWidget(c *gin.Context) {

	ok, w := CheckPermissions(c, common.Read, -1)
	if !ok {
		return
	}

	serializer := common.WidgetSerializer{c, w.Widget}
	c.JSON(http.StatusOK, gin.H{
		"widget": serializer.Response(),
	})
}

// deleteWidget godoc
// @Summary Delete a widget
// @ID deleteWidget
// @Tags widgets
// @Produce json
// @Success 200 "OK."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param widgetID path int true "Widget ID"
// @Router /widgets/{widgetID} [delete]
func deleteWidget(c *gin.Context) {

	ok, w := CheckPermissions(c, common.Delete, -1)
	if !ok {
		return
	}

	err := w.delete()
	if common.ProvideErrorResponse(c, err) {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OK.",
	})
}
