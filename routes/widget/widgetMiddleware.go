package widget

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/visualization"
)

func CheckPermissions(c *gin.Context, operation common.CRUD) (bool, Widget) {

	var w Widget

	err := common.ValidateRole(c, common.ModelWidget, operation)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Access denied (role validation failed).")
		return false, w
	}

	widgetID, err := strconv.Atoi(c.Param("widgetID"))
	if err != nil {
		errormsg := fmt.Sprintf("Bad request. No or incorrect format of widgetID path parameter")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errormsg,
		})
		return false, w
	}

	err = w.ByID(uint(widgetID))
	if common.ProvideErrorResponse(c, err) {
		return false, w
	}

	ok, _ := visualization.CheckPermissions(c, operation, "body", int(w.VisualizationID))
	if !ok {
		return false, w
	}

	return true, w
}
