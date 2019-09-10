package simulator

import (
	"fmt"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/database"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/helper"
	"github.com/gin-gonic/gin"
	"strconv"
)

func checkPermissions(c *gin.Context, modeltype database.ModelName, operation database.CRUD, hasID bool) (bool, Simulator) {

	var s Simulator

	err := database.ValidateRole(c, modeltype, operation)
	if err != nil {
		helper.UnprocessableEntityError(c, err.Error())
		return false, s
	}

	if hasID {
		// Get the ID of the simulator from the context
		simulatorID, err := strconv.Atoi(c.Param("simulatorID"))
		if err != nil {
			helper.BadRequestError(c, fmt.Sprintf("Could not get simulator's ID from context"))
			return false, s
		}

		err = s.ByID(uint(simulatorID))
		if helper.DBError(c, err) {
			return false, s
		}

	}

	return true, s
}
