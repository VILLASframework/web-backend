/** component_configuration package, endpoints.
*
* @author Sonja Happ <sonja.happ@eonerc.rwth-aachen.de>
* @copyright 2014-2019, Institute for Automation of Complex Power Systems, EONERC
* @license GNU General Public License (version 3)
*
* VILLASweb-backend-go
*
* This program is free software: you can redistribute it and/or modify
* it under the terms of the GNU General Public License as published by
* the Free Software Foundation, either version 3 of the License, or
* any later version.
*
* This program is distributed in the hope that it will be useful,
* but WITHOUT ANY WARRANTY; without even the implied warranty of
* MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
* GNU General Public License for more details.
*
* You should have received a copy of the GNU General Public License
* along with this program.  If not, see <http://www.gnu.org/licenses/>.
*********************************************************************************/
package component_configuration

import (
	"net/http"

	"git.rwth-aachen.de/acs/public/villas/web-backend-go/helper"

	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/scenario"
)

func RegisterComponentConfigurationEndpoints(r *gin.RouterGroup) {
	r.GET("", getConfigs)
	r.POST("", addConfig)
	r.PUT("/:configID", updateConfig)
	r.GET("/:configID", getConfig)
	r.DELETE("/:configID", deleteConfig)
}

// getConfigs godoc
// @Summary Get all component configurations of scenario
// @ID getConfigs
// @Produce  json
// @Tags component-configurations
// @Success 200 {object} api.ResponseConfigs "Component configurations which belong to scenario"
// @Failure 404 {object} api.ResponseError "Not found"
// @Failure 422 {object} api.ResponseError "Unprocessable entity"
// @Failure 500 {object} api.ResponseError "Internal server error"
// @Param scenarioID query int true "Scenario ID"
// @Router /configs [get]
// @Security Bearer
func getConfigs(c *gin.Context) {

	ok, so := scenario.CheckPermissions(c, database.Read, "query", -1)
	if !ok {
		return
	}

	db := database.GetDB()
	var configs []database.ComponentConfiguration
	err := db.Order("ID asc").Model(so).Related(&configs, "ComponentConfigurations").Error
	if !helper.DBError(c, err) {
		c.JSON(http.StatusOK, gin.H{"configs": configs})
	}

}

// addConfig godoc
// @Summary Add a component configuration to a scenario
// @ID addConfig
// @Accept json
// @Produce json
// @Tags component-configurations
// @Success 200 {object} api.ResponseConfig "Component configuration that was added"
// @Failure 400 {object} api.ResponseError "Bad request"
// @Failure 404 {object} api.ResponseError "Not found"
// @Failure 422 {object} api.ResponseError "Unprocessable entity"
// @Failure 500 {object} api.ResponseError "Internal server error"
// @Param inputConfig body component_configuration.addConfigRequest true "component configuration to be added incl. IDs of scenario and IC"
// @Router /configs [post]
// @Security Bearer
func addConfig(c *gin.Context) {

	// Bind the request to JSON
	var req addConfigRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		helper.BadRequestError(c, "Bad request. Error binding form data to JSON: "+err.Error())
		return
	}

	// validate the request
	if err = req.validate(); err != nil {
		helper.UnprocessableEntityError(c, err.Error())
		return
	}

	// Create the new Component Configuration from the request
	newConfig := req.createConfig()

	// check access to the scenario
	ok, _ := scenario.CheckPermissions(c, database.Update, "body", int(newConfig.ScenarioID))
	if !ok {
		return
	}

	// add the new Component Configuration to the scenario
	err = newConfig.addToScenario()
	if !helper.DBError(c, err) {
		c.JSON(http.StatusOK, gin.H{"config": newConfig.ComponentConfiguration})
	}

}

// updateConfig godoc
// @Summary Update a component configuration
// @ID updateConfig
// @Tags component-configurations
// @Accept json
// @Produce json
// @Success 200 {object} api.ResponseConfig "Component configuration that was added"
// @Failure 400 {object} api.ResponseError "Bad request"
// @Failure 404 {object} api.ResponseError "Not found"
// @Failure 422 {object} api.ResponseError "Unprocessable entity"
// @Failure 500 {object} api.ResponseError "Internal server error"
// @Param inputConfig body component_configuration.updateConfigRequest true "component configuration to be updated"
// @Param configID path int true "Config ID"
// @Router /configs/{configID} [put]
// @Security Bearer
func updateConfig(c *gin.Context) {

	ok, oldConfig := CheckPermissions(c, database.Update, "path", -1)
	if !ok {
		return
	}

	var req updateConfigRequest
	err := c.BindJSON(&req)
	if err != nil {
		helper.BadRequestError(c, "Error binding form data to JSON: "+err.Error())
		return
	}

	// Validate the request
	if err := req.Config.validate(); err != nil {
		helper.BadRequestError(c, err.Error())
		return
	}

	// Create the updateConfig from oldConfig
	updatedConfig := req.updateConfig(oldConfig)

	// Finally, update the Component Configuration
	err = oldConfig.Update(updatedConfig)
	if !helper.DBError(c, err) {
		c.JSON(http.StatusOK, gin.H{"config": updatedConfig.ComponentConfiguration})
	}

}

// getConfig godoc
// @Summary Get a component configuration
// @ID getConfig
// @Tags component-configurations
// @Produce json
// @Success 200 {object} api.ResponseConfig "component configuration that was requested"
// @Failure 400 {object} api.ResponseError "Bad request"
// @Failure 404 {object} api.ResponseError "Not found"
// @Failure 422 {object} api.ResponseError "Unprocessable entity"
// @Failure 500 {object} api.ResponseError "Internal server error"
// @Param configID path int true "Config ID"
// @Router /configs/{configID} [get]
// @Security Bearer
func getConfig(c *gin.Context) {

	ok, m := CheckPermissions(c, database.Read, "path", -1)
	if !ok {
		return
	}

	c.JSON(http.StatusOK, gin.H{"config": m.ComponentConfiguration})
}

// deleteConfig godoc
// @Summary Delete a component configuration
// @ID deleteConfig
// @Tags component-configurations
// @Produce json
// @Success 200 {object} api.ResponseConfig "component configuration that was deleted"
// @Failure 400 {object} api.ResponseError "Bad request"
// @Failure 404 {object} api.ResponseError "Not found"
// @Failure 422 {object} api.ResponseError "Unprocessable entity"
// @Failure 500 {object} api.ResponseError "Internal server error"
// @Param configID path int true "Config ID"
// @Router /configs/{configID} [delete]
// @Security Bearer
func deleteConfig(c *gin.Context) {

	ok, m := CheckPermissions(c, database.Delete, "path", -1)
	if !ok {
		return
	}

	err := m.delete()
	if !helper.DBError(c, err) {
		c.JSON(http.StatusOK, gin.H{"config": m.ComponentConfiguration})
	}
}
