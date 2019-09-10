package scenario

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
)

func CheckPermissions(c *gin.Context, operation common.CRUD, screnarioIDSource string, scenarioIDbody int) (bool, Scenario) {

	var so Scenario

	err := common.ValidateRole(c, common.ModelScenario, operation)
	if err != nil {
		common.UnprocessableEntityError(c, fmt.Sprintf("Access denied (role validation failed): %v", err))
		return false, so
	}

	if operation == common.Create || (operation == common.Read && screnarioIDSource == "none") {
		return true, so
	}

	var scenarioID int
	if screnarioIDSource == "path" {
		scenarioID, err = strconv.Atoi(c.Param("scenarioID"))
		if err != nil {
			common.BadRequestError(c, fmt.Sprintf("No or incorrect format of scenarioID path parameter"))
			return false, so
		}
	} else if screnarioIDSource == "query" {
		scenarioID, err = strconv.Atoi(c.Request.URL.Query().Get("scenarioID"))
		if err != nil {
			common.BadRequestError(c, fmt.Sprintf("No or incorrect format of scenarioID query parameter"))
			return false, so
		}
	} else if screnarioIDSource == "body" {
		scenarioID = scenarioIDbody

	} else {
		common.BadRequestError(c, fmt.Sprintf("The following source of scenario ID is not valid: %s", screnarioIDSource))
		return false, so
	}

	userID, _ := c.Get(common.UserIDCtx)
	userRole, _ := c.Get(common.UserRoleCtx)

	err = so.ByID(uint(scenarioID))
	if common.DBError(c, err) {
		return false, so
	}

	if so.checkAccess(userID.(uint), userRole.(string)) == false {
		common.UnprocessableEntityError(c, "Access denied (for scenario ID).")
		return false, so
	}

	return true, so
}
