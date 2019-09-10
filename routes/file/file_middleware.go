package file

import (
	"fmt"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/database"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/helper"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/simulationmodel"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/widget"
	"github.com/gin-gonic/gin"
)

func checkPermissions(c *gin.Context, operation database.CRUD) (bool, File) {

	var f File

	err := database.ValidateRole(c, database.ModelFile, operation)
	if err != nil {
		helper.UnprocessableEntityError(c, fmt.Sprintf("Access denied (role validation failed): %v", err.Error()))
		return false, f
	}

	fileID, err := helper.GetIDOfElement(c, "fileID", "path", -1)
	if err != nil {
		return false, f
	}

	err = f.byID(uint(fileID))
	if helper.DBError(c, err) {
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
