package widget

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/dashboard"
)

func CheckPermissions(c *gin.Context, operation common.CRUD, widgetIDBody int) (bool, Widget) {

	var w Widget

	err := common.ValidateRole(c, common.ModelWidget, operation)
	if err != nil {
		common.UnprocessableEntityError(c, fmt.Sprintf("Access denied (role validation failed): %v", err.Error()))
		return false, w
	}

	var widgetID int
	if widgetIDBody < 0 {
		widgetID, err = strconv.Atoi(c.Param("widgetID"))
		if err != nil {
			common.BadRequestError(c, fmt.Sprintf("No or incorrect format of widgetID path parameter"))
			return false, w
		}
	} else {
		widgetID = widgetIDBody
	}

	err = w.ByID(uint(widgetID))
	if common.DBError(c, err) {
		return false, w
	}

	ok, _ := dashboard.CheckPermissions(c, operation, "body", int(w.DashboardID))
	if !ok {
		return false, w
	}

	return true, w
}
