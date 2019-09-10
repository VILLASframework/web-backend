package dashboard

import (
	"fmt"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/helper"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/scenario"
	"strconv"

	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/database"
)

func CheckPermissions(c *gin.Context, operation database.CRUD, dabIDSource string, dabIDBody int) (bool, Dashboard) {

	var dab Dashboard

	err := database.ValidateRole(c, database.ModelDashboard, operation)
	if err != nil {
		helper.UnprocessableEntityError(c, fmt.Sprintf("Access denied (role validation failed): %v", err.Error()))
		return false, dab
	}

	var dabID int
	if dabIDSource == "path" {
		dabID, err = strconv.Atoi(c.Param("dashboardID"))
		if err != nil {
			helper.BadRequestError(c, fmt.Sprintf("No or incorrect format of dashboardID path parameter"))
			return false, dab
		}
	} else if dabIDSource == "query" {
		dabID, err = strconv.Atoi(c.Request.URL.Query().Get("dashboardID"))
		if err != nil {
			helper.BadRequestError(c, fmt.Sprintf("No or incorrect format of dashboardID query parameter"))
			return false, dab
		}
	} else if dabIDSource == "body" {
		dabID = dabIDBody
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
