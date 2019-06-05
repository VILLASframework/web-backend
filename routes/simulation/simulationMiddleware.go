package simulation

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
)

func CheckPermissions(c *gin.Context, modelname common.ModelName, operation common.CRUD, simIDSource string) (bool, Simulation) {

	var sim Simulation

	err := common.ValidateRole(c, modelname, operation)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Access denied (role validation failed).")
		return false, sim
	}

	if operation == common.Create || (operation == common.Read && simIDSource == "none") {
		return true, sim
	}

	var simID int
	if simIDSource == "path" {
		simID, err = strconv.Atoi(c.Param("simulationID"))
		if err != nil {
			errormsg := fmt.Sprintf("Bad request. No or incorrect format of simulationID path parameter")
			c.JSON(http.StatusBadRequest, gin.H{
				"error": errormsg,
			})
			return false, sim
		}
	} else if simIDSource == "query" {
		simID, err = strconv.Atoi(c.Request.URL.Query().Get("simulationID"))
		if err != nil {
			errormsg := fmt.Sprintf("Bad request. No or incorrect format of simulationID query parameter")
			c.JSON(http.StatusBadRequest, gin.H{
				"error": errormsg,
			})
			return false, sim
		}
	} else {
		errormsg := fmt.Sprintf("Bad request. The following source of your simulation ID is not valid: %s", simIDSource)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errormsg,
		})
		return false, sim
	}

	userID, _ := c.Get(common.UserIDCtx)
	userRole, _ := c.Get(common.UserRoleCtx)

	err = sim.ByID(uint(simID))
	if common.ProvideErrorResponse(c, err) {
		return false, sim
	}

	if sim.checkAccess(userID.(uint), userRole.(string)) == false {
		c.JSON(http.StatusUnprocessableEntity, "Access denied (for simulation ID).")
		return false, sim
	}

	return true, sim
}
