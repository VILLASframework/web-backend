package signal

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/simulationmodel"
)

func checkPermissions(c *gin.Context, operation common.CRUD) (bool, Signal) {

	var sig Signal

	err := common.ValidateRole(c, common.ModelSignal, operation)
	if err != nil {
		common.UnprocessableEntityError(c, fmt.Sprintf("Access denied (role validation failed): %v", err.Error()))
		return false, sig
	}

	signalID, err := strconv.Atoi(c.Param("signalID"))
	if err != nil {
		common.BadRequestError(c, fmt.Sprintf("No or incorrect format of signalID path parameter"))
		return false, sig
	}

	err = sig.byID(uint(signalID))
	if common.DBError(c, err) {
		return false, sig
	}

	ok, _ := simulationmodel.CheckPermissions(c, operation, "body", int(sig.SimulationModelID))
	if !ok {
		return false, sig
	}

	return true, sig
}
