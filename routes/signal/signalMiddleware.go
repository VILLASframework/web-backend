package signal

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/simulationmodel"
)

func checkPermissions(c *gin.Context, operation common.CRUD) (bool, Signal) {

	var sig Signal

	err := common.ValidateRole(c, common.ModelSignal, operation)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Access denied (role validation failed).")
		return false, sig
	}

	signalID, err := strconv.Atoi(c.Param("signalID"))
	if err != nil {
		errormsg := fmt.Sprintf("Bad request. No or incorrect format of signalID path parameter")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errormsg,
		})
		return false, sig
	}

	err = sig.byID(uint(signalID))
	if common.ProvideErrorResponse(c, err) {
		return false, sig
	}

	ok, _ := simulationmodel.CheckPermissions(c, operation, "body", int(sig.SimulationModelID))
	if !ok {
		return false, sig
	}

	return true, sig
}
