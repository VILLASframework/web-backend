package signal

import (
	"fmt"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/helper"
	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/database"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/simulationmodel"
)

func checkPermissions(c *gin.Context, operation database.CRUD) (bool, Signal) {

	var sig Signal

	err := database.ValidateRole(c, database.ModelSignal, operation)
	if err != nil {
		helper.UnprocessableEntityError(c, fmt.Sprintf("Access denied (role validation failed): %v", err.Error()))
		return false, sig
	}

	signalID, err := helper.GetIDOfElement(c, "signalID", "path", -1)
	if err != nil {
		return false, sig
	}

	err = sig.byID(uint(signalID))
	if helper.DBError(c, err) {
		return false, sig
	}

	ok, _ := simulationmodel.CheckPermissions(c, operation, "body", int(sig.SimulationModelID))
	if !ok {
		return false, sig
	}

	return true, sig
}
