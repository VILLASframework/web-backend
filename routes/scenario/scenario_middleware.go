package scenario

import (
	"fmt"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/helper"
	"strconv"

	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/database"
)

func CheckPermissions(c *gin.Context, operation database.CRUD, screnarioIDSource string, scenarioIDbody int) (bool, Scenario) {

	var so Scenario

	err := database.ValidateRole(c, database.ModelScenario, operation)
	if err != nil {
		helper.UnprocessableEntityError(c, fmt.Sprintf("Access denied (role validation failed): %v", err))
		return false, so
	}

	if operation == database.Create || (operation == database.Read && screnarioIDSource == "none") {
		return true, so
	}

	var scenarioID int
	if screnarioIDSource == "path" {
		scenarioID, err = strconv.Atoi(c.Param("scenarioID"))
		if err != nil {
			helper.BadRequestError(c, fmt.Sprintf("No or incorrect format of scenarioID path parameter"))
			return false, so
		}
	} else if screnarioIDSource == "query" {
		scenarioID, err = strconv.Atoi(c.Request.URL.Query().Get("scenarioID"))
		if err != nil {
			helper.BadRequestError(c, fmt.Sprintf("No or incorrect format of scenarioID query parameter"))
			return false, so
		}
	} else if screnarioIDSource == "body" {
		scenarioID = scenarioIDbody

	} else {
		helper.BadRequestError(c, fmt.Sprintf("The following source of scenario ID is not valid: %s", screnarioIDSource))
		return false, so
	}

	userID, _ := c.Get(database.UserIDCtx)
	userRole, _ := c.Get(database.UserRoleCtx)

	err = so.ByID(uint(scenarioID))
	if helper.DBError(c, err) {
		return false, so
	}

	if so.checkAccess(userID.(uint), userRole.(string)) == false {
		helper.UnprocessableEntityError(c, "Access denied (for scenario ID).")
		return false, so
	}

	return true, so
}
