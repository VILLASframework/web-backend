package result

import (
	"fmt"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/helper"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/scenario"
	"github.com/gin-gonic/gin"
)

func CheckPermissions(c *gin.Context, operation database.CRUD, resultIDSource string, resultIDBody int) (bool, Result) {

	var result Result

	err := database.ValidateRole(c, database.ModelResult, operation)
	if err != nil {
		helper.UnprocessableEntityError(c, fmt.Sprintf("Access denied (role validation failed): %v", err.Error()))
		return false, result
	}

	resultID, err := helper.GetIDOfElement(c, "resultID", resultIDSource, resultIDBody)
	if err != nil {
		return false, result
	}

	err = result.ByID(uint(resultID))
	if helper.DBError(c, err) {
		return false, result
	}

	ok, _ := scenario.CheckPermissions(c, operation, "body", int(result.ScenarioID))
	if !ok {
		return false, result
	}

	return true, result
}
