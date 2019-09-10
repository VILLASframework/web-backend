package simulationmodel

import (
	"fmt"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/helper"
	"strconv"

	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/database"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/scenario"
)

func CheckPermissions(c *gin.Context, operation database.CRUD, modelIDSource string, modelIDBody int) (bool, SimulationModel) {

	var m SimulationModel

	err := database.ValidateRole(c, database.ModelSimulationModel, operation)
	if err != nil {
		helper.UnprocessableEntityError(c, fmt.Sprintf("Access denied (role validation failed): %v", err.Error()))
		return false, m
	}

	var modelID int
	if modelIDSource == "path" {
		modelID, err = strconv.Atoi(c.Param("modelID"))
		if err != nil {
			helper.BadRequestError(c, fmt.Sprintf("No or incorrect format of modelID path parameter"))
			return false, m
		}
	} else if modelIDSource == "query" {
		modelID, err = strconv.Atoi(c.Request.URL.Query().Get("modelID"))
		if err != nil {
			helper.BadRequestError(c, fmt.Sprintf("No or incorrect format of modelID query parameter"))
			return false, m
		}
	} else if modelIDSource == "body" {
		modelID = modelIDBody
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
