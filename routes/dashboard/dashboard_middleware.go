package dashboard

import (
	"fmt"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/helper"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/scenario"
	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
)

func CheckPermissions(c *gin.Context, operation database.CRUD, dabIDSource string, dabIDBody int) (bool, Dashboard) {

	var dab Dashboard

	err := database.ValidateRole(c, database.ModelDashboard, operation)
	if err != nil {
		helper.UnprocessableEntityError(c, fmt.Sprintf("Access denied (role validation failed): %v", err.Error()))
		return false, dab
	}

	dabID, err := helper.GetIDOfElement(c, "dashboardID", dabIDSource, dabIDBody)
	if err != nil {
		return false, dab
	}

	err = dab.ByID(uint(dabID))
	if helper.DBError(c, err) {
		return false, dab
	}

	ok, _ := scenario.CheckPermissions(c, operation, "body", int(dab.ScenarioID))
	if !ok {
		return false, dab
	}

	return true, dab
}
