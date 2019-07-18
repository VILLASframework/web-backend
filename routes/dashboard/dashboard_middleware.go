package dashboard

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/simulation"
)

func CheckPermissions(c *gin.Context, operation common.CRUD, dabIDSource string, dabIDBody int) (bool, Dashboard) {

	var dab Dashboard

	err := common.ValidateRole(c, common.ModelDashboard, operation)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Access denied (role validation failed).")
		return false, dab
	}

	var dabID int
	if dabIDSource == "path" {
		dabID, err = strconv.Atoi(c.Param("dashboardID"))
		if err != nil {
			errormsg := fmt.Sprintf("Bad request. No or incorrect format of dashboardID path parameter")
			c.JSON(http.StatusBadRequest, gin.H{
				"error": errormsg,
			})
			return false, dab
		}
	} else if dabIDSource == "query" {
		dabID, err = strconv.Atoi(c.Request.URL.Query().Get("dashboardID"))
		if err != nil {
			errormsg := fmt.Sprintf("Bad request. No or incorrect format of dashboardID query parameter")
			c.JSON(http.StatusBadRequest, gin.H{
				"error": errormsg,
			})
			return false, dab
		}
	} else if dabIDSource == "body" {
		dabID = dabIDBody
	}

	err = dab.ByID(uint(dabID))
	if common.ProvideErrorResponse(c, err) {
		return false, dab
	}

	ok, _ := simulation.CheckPermissions(c, operation, "body", int(dab.SimulationID))
	if !ok {
		return false, dab
	}

	return true, dab
}
