package simulationmodel

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/scenario"
)

func CheckPermissions(c *gin.Context, operation common.CRUD, modelIDSource string, modelIDBody int) (bool, SimulationModel) {

	var m SimulationModel

	err := common.ValidateRole(c, common.ModelSimulationModel, operation)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"success": false,
			"message": fmt.Sprintf("Access denied (role validation failed): %v", err),
		})
		return false, m
	}

	var modelID int
	if modelIDSource == "path" {
		modelID, err = strconv.Atoi(c.Param("modelID"))
		if err != nil {
			errormsg := fmt.Sprintf("Bad request. No or incorrect format of modelID path parameter")
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": errormsg,
			})
			return false, m
		}
	} else if modelIDSource == "query" {
		modelID, err = strconv.Atoi(c.Request.URL.Query().Get("modelID"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": fmt.Sprintf("Bad request. No or incorrect format of modelID query parameter"),
			})
			return false, m
		}
	} else if modelIDSource == "body" {
		modelID = modelIDBody
	}

	err = m.ByID(uint(modelID))
	if common.DBError(c, err) {
		return false, m
	}

	ok, _ := scenario.CheckPermissions(c, operation, "body", int(m.ScenarioID))
	if !ok {
		return false, m
	}

	return true, m
}
