package widget

import (
	"fmt"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/helper"
	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/database"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/dashboard"
)

func CheckPermissions(c *gin.Context, operation database.CRUD, widgetIDBody int) (bool, Widget) {

	var w Widget
	var err error
	err = database.ValidateRole(c, database.ModelWidget, operation)
	if err != nil {
		helper.UnprocessableEntityError(c, fmt.Sprintf("Access denied (role validation failed): %v", err.Error()))
		return false, w
	}

	var widgetID int
	if widgetIDBody < 0 {
		widgetID, err = helper.GetIDOfElement(c, "widgetID", "path", -1)
		if err != nil {
			return false, w
		}
	} else {
		widgetID = widgetIDBody
	}

	err = w.ByID(uint(widgetID))
	if helper.DBError(c, err) {
		return false, w
	}

	ok, _ := dashboard.CheckPermissions(c, operation, "body", int(w.DashboardID))
	if !ok {
		return false, w
	}

	return true, w
}
