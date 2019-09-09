package file

import (
	"fmt"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/simulationmodel"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/widget"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func checkPermissions(c *gin.Context, operation common.CRUD) (bool, File) {

	var f File

	err := common.ValidateRole(c, common.ModelFile, operation)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"success": false,
			"message": fmt.Sprintf("Access denied (role validation failed): %v", err),
		})
		return false, f
	}

	fileID, err := strconv.Atoi(c.Param("fileID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   fmt.Sprintf("Bad request. No or incorrect format of fileID path parameter"),
		})
		return false, f
	}

	err = f.byID(uint(fileID))
	if common.DBError(c, err) {
		return false, f
	}

	if f.SimulationModelID > 0 {
		ok, _ := simulationmodel.CheckPermissions(c, operation, "body", int(f.SimulationModelID))
		if !ok {
			return false, f
		}
	} else {
		ok, _ := widget.CheckPermissions(c, operation, int(f.WidgetID))
		if !ok {
			return false, f
		}
	}

	return true, f
}
