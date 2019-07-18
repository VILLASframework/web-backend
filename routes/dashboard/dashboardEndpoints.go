package dashboard

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/simulation"
)

func RegisterDashboardEndpoints(r *gin.RouterGroup) {

	r.GET("", getDashboards)
	r.POST("", addDashboard)
	r.PUT("/:dashboardID", updateDashboard)
	r.GET("/:dashboardID", getDashboard)
	r.DELETE("/:dashboardID", deleteDashboard)
}

// getDashboards godoc
// @Summary Get all dashboards of simulation
// @ID getDashboards
// @Produce  json
// @Tags dashboards
// @Success 200 {array} common.DashboardResponse "Array of dashboards to which belong to simulation"
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param simulationID query int true "Simulation ID"
// @Router /dashboards [get]
func getDashboards(c *gin.Context) {

	ok, sim := simulation.CheckPermissions(c, common.Read, "query", -1)
	if !ok {
		return
	}

	db := common.GetDB()
	var dab []common.Dashboard
	err := db.Order("ID asc").Model(sim).Related(&dab, "Dashboards").Error
	if common.ProvideErrorResponse(c, err) {
		return
	}

	serializer := common.DashboardsSerializer{c, dab}
	c.JSON(http.StatusOK, gin.H{
		"dashboards": serializer.Response(),
	})
}

// addDashboard godoc
// @Summary Add a dashboard to a simulation
// @ID addDashboard
// @Accept json
// @Produce json
// @Tags dashboards
// @Param inputDab body common.DashboardResponse true "Dashboard to be added incl. ID of simulation"
// @Success 200 "OK."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Router /dashboards [post]
func addDashboard(c *gin.Context) {

	var newDab Dashboard
	err := c.BindJSON(&newDab)
	if err != nil {
		errormsg := "Bad request. Error binding form data to JSON: " + err.Error()
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errormsg,
		})
		return
	}

	ok, _ := simulation.CheckPermissions(c, common.Create, "body", int(newDab.SimulationID))
	if !ok {
		return
	}

	// add dashboard to DB and add association to simulation
	err = newDab.addToSimulation()
	if common.ProvideErrorResponse(c, err) == false {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK.",
		})
	}

}

// updateDashboard godoc
// @Summary Update a dashboard
// @ID updateDashboard
// @Tags dashboards
// @Accept json
// @Produce json
// @Param inputDab body common.DashboardResponse true "Dashboard to be updated"
// @Success 200 "OK."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param dashboardID path int true "Dashboard ID"
// @Router /dashboards/{dashboardID} [put]
func updateDashboard(c *gin.Context) {

	ok, d := CheckPermissions(c, common.Update, "path", -1)
	if !ok {
		return
	}

	var modifiedDab Dashboard
	err := c.BindJSON(&modifiedDab)
	if err != nil {
		errormsg := "Bad request. Error binding form data to JSON: " + err.Error()
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errormsg,
		})
		return
	}

	err = d.update(modifiedDab)
	if common.ProvideErrorResponse(c, err) == false {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK.",
		})
	}
}

// getDashboard godoc
// @Summary Get a dashboard
// @ID getDashboard
// @Tags dashboards
// @Produce json
// @Success 200 {object} common.DashboardResponse "Requested dashboard."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param dashboardID path int true "Dashboard ID"
// @Router /dashboards/{dashboardID} [get]
func getDashboard(c *gin.Context) {

	ok, dab := CheckPermissions(c, common.Read, "path", -1)
	if !ok {
		return
	}

	serializer := common.DashboardSerializer{c, dab.Dashboard}
	c.JSON(http.StatusOK, gin.H{
		"dashboard": serializer.Response(),
	})
}

// deleteDashboard godoc
// @Summary Delete a dashboard
// @ID deleteDashboard
// @Tags dashboards
// @Produce json
// @Success 200 "OK."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param dashboardID path int true "Dashboard ID"
// @Router /dashboards/{dashboardID} [delete]
func deleteDashboard(c *gin.Context) {
	ok, dab := CheckPermissions(c, common.Delete, "path", -1)
	if !ok {
		return
	}

	err := dab.delete()
	if common.ProvideErrorResponse(c, err) {
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "OK.",
	})
}
