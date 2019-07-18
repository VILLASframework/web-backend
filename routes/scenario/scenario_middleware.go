package scenario

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
)

func CheckPermissions(c *gin.Context, operation common.CRUD, simIDSource string, simIDBody int) (bool, Scenario) {

	var so Scenario

	err := common.ValidateRole(c, common.ModelScenario, operation)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Access denied (role validation failed).")
		return false, so
	}

	if operation == common.Create || (operation == common.Read && simIDSource == "none") {
		return true, so
	}

	var simID int
	if simIDSource == "path" {
		simID, err = strconv.Atoi(c.Param("scenarioID"))
		if err != nil {
			errormsg := fmt.Sprintf("Bad request. No or incorrect format of scenarioID path parameter")
			c.JSON(http.StatusBadRequest, gin.H{
				"error": errormsg,
			})
			return false, so
		}
	} else if simIDSource == "query" {
		simID, err = strconv.Atoi(c.Request.URL.Query().Get("scenarioID"))
		if err != nil {
			errormsg := fmt.Sprintf("Bad request. No or incorrect format of scenarioID query parameter")
			c.JSON(http.StatusBadRequest, gin.H{
				"error": errormsg,
			})
			return false, so
		}
	} else if simIDSource == "body" {
		simID = simIDBody

	} else {
		errormsg := fmt.Sprintf("Bad request. The following source of your scenario ID is not valid: %s", simIDSource)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errormsg,
		})
		return false, so
	}

	userID, _ := c.Get(common.UserIDCtx)
	userRole, _ := c.Get(common.UserRoleCtx)

	err = so.ByID(uint(simID))
	if common.ProvideErrorResponse(c, err) {
		return false, so
	}

	if so.checkAccess(userID.(uint), userRole.(string)) == false {
		c.JSON(http.StatusUnprocessableEntity, "Access denied (for scenario ID).")
		return false, so
	}

	return true, so
}
