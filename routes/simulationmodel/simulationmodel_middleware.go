package simulationmodel

import (
	"fmt"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/helper"
	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/database"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/scenario"
)

func CheckPermissions(c *gin.Context, operation database.CRUD, modelIDSource string, modelIDBody int) (bool, SimulationModel) {

	var m SimulationModel

	err := database.ValidateRole(c, database.ModelSimulationModel, operation)
	if err != nil {
		helper.UnprocessableEntityError(c, fmt.Sprintf("Access denied (role validation of simulation model failed): %v", err.Error()))
		return false, m
	}

	modelID, err := helper.GetIDOfElement(c, "modelID", modelIDSource, modelIDBody)
	if err != nil {
		return false, m
	}

	err = m.ByID(uint(modelID))
	if helper.DBError(c, err) {
		return false, m
	}

	ok, _ := scenario.CheckPermissions(c, operation, "body", int(m.ScenarioID))
	if !ok {
		return false, m
	}

	return true, m
}
