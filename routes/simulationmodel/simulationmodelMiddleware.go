package simulationmodel

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/simulation"
)

func checkPermissions(c *gin.Context, operation common.CRUD) (bool, SimulationModel) {

	var m SimulationModel

	modelID, err := strconv.Atoi(c.Param("modelID"))

	if err != nil {
		errormsg := fmt.Sprintf("Bad request. No or incorrect format of model ID in path")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errormsg,
		})
		return false, m
	}

	err = m.ByID(uint(modelID))
	if common.ProvideErrorResponse(c, err) {
		return false, m
	}

	ok, _ := simulation.CheckPermissions(c, common.ModelSimulationModel, operation, "body", int(m.SimulationID))
	if !ok {
		return false, m
	}

	return true, m
}
