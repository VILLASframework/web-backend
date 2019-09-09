package simulator

import (
	"fmt"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"github.com/gin-gonic/gin"
	"strconv"
)

func checkPermissions(c *gin.Context, modeltype common.ModelName, operation common.CRUD, hasID bool) (bool, Simulator) {

	var s Simulator

	err := common.ValidateRole(c, modeltype, operation)
	if err != nil {
		common.UnprocessableEntityError(c, err.Error())
		return false, s
	}

	if hasID {
		// Get the ID of the simulator from the context
		simulatorID, err := strconv.Atoi(c.Param("simulatorID"))
		if err != nil {
			common.BadRequestError(c, fmt.Sprintf("Could not get simulator's ID from context"))
			return false, s
		}

		err = s.ByID(uint(simulatorID))
		if common.DBError(c, err) {
			return false, s
		}

	}

	return true, s
}
