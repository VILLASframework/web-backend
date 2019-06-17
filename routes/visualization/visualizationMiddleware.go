package visualization

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/simulation"
)

func CheckPermissions(c *gin.Context, operation common.CRUD, visIDSource string, visIDBody int) (bool, Visualization) {

	var vis Visualization

	err := common.ValidateRole(c, common.ModelVisualization, operation)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Access denied (role validation failed).")
		return false, vis
	}

	var visID int
	if visIDSource == "path" {
		visID, err = strconv.Atoi(c.Param("visualizationID"))
		if err != nil {
			errormsg := fmt.Sprintf("Bad request. No or incorrect format of simulationID path parameter")
			c.JSON(http.StatusBadRequest, gin.H{
				"error": errormsg,
			})
			return false, vis
		}
	} else if visIDSource == "query" {
		visID, err = strconv.Atoi(c.Request.URL.Query().Get("visualizationID"))
		if err != nil {
			errormsg := fmt.Sprintf("Bad request. No or incorrect format of visualizationID query parameter")
			c.JSON(http.StatusBadRequest, gin.H{
				"error": errormsg,
			})
			return false, vis
		}
	} else if visIDSource == "body" {
		visID = visIDBody
	}

	err = vis.ByID(uint(visID))
	if common.ProvideErrorResponse(c, err) {
		return false, vis
	}

	ok, _ := simulation.CheckPermissions(c, operation, "body", int(vis.SimulationID))
	if !ok {
		return false, vis
	}

	return true, vis
}
