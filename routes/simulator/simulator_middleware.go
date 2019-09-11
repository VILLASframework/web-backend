package simulator

import (
	"fmt"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/database"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/helper"
	"github.com/gin-gonic/gin"
)

func CheckPermissions(c *gin.Context, modeltype database.ModelName, operation database.CRUD, hasID bool) (bool, Simulator) {

	var s Simulator

	err := database.ValidateRole(c, modeltype, operation)
	if err != nil {
		helper.UnprocessableEntityError(c, fmt.Sprintf("Access denied (role validation of simulator failed): %v", err.Error()))
		return false, s
	}

	if hasID {
		// Get the ID of the simulator from the context
		simulatorID, err := helper.GetIDOfElement(c, "simulatorID", "path", -1)
		if err != nil {
			return false, s
		}

		err = s.ByID(uint(simulatorID))
		if helper.DBError(c, err) {
			return false, s
		}
	}

	return true, s
}
