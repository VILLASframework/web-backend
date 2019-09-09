package file

import (
	"fmt"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/simulationmodel"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/widget"
	"github.com/gin-gonic/gin"
	"strconv"
)

func checkPermissions(c *gin.Context, operation common.CRUD) (bool, File) {

	var f File

	err := common.ValidateRole(c, common.ModelFile, operation)
	if err != nil {
		common.UnprocessableEntityError(c, fmt.Sprintf("Access denied (role validation failed): %v", err.Error()))
		return false, f
	}

	fileID, err := strconv.Atoi(c.Param("fileID"))
	if err != nil {
		common.BadRequestError(c, fmt.Sprintf("No or incorrect format of fileID path parameter"))
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
