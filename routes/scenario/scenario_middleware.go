package scenario

import (
	"fmt"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/helper"
	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/database"
)

func CheckPermissions(c *gin.Context, operation database.CRUD, screnarioIDSource string, scenarioIDbody int) (bool, Scenario) {

	var so Scenario

	err := database.ValidateRole(c, database.ModelScenario, operation)
	if err != nil {
		helper.UnprocessableEntityError(c, fmt.Sprintf("Access denied (role validation of scenario failed): %v", err))
		return false, so
	}

	if operation == database.Create || (operation == database.Read && screnarioIDSource == "none") {
		return true, so
	}

	scenarioID, err := helper.GetIDOfElement(c, "scenarioID", screnarioIDSource, scenarioIDbody)
	if err != nil {
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
