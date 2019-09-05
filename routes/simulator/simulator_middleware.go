package simulator

import (
	"fmt"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"github.com/gin-gonic/gin"
	"net/http"
)

func checkPermissions(c *gin.Context, modeltype common.ModelName, operation common.CRUD, hasID bool) (bool, Simulator) {

	var s Simulator

	err := common.ValidateRole(c, modeltype, operation)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"success": false,
			"message": fmt.Sprintf("%v", err),
		})
		return false, s
	}

	if hasID {
		// Get the ID of the simulator from the context
		simulatorID, err := common.UintParamFromCtx(c, "simulatorID")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": fmt.Sprintf("Could not get simulator's ID from context"),
			})
			return false, s
		}

		err = s.ByID(uint(simulatorID))
		if common.ProvideErrorResponse(c, err) {
			return false, s
		}

	}

	return true, s
}
