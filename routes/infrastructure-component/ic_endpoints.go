/** InfrastructureComponent package, endpoints.
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
package infrastructure_component

import (
	"log"
	"net/http"

	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/helper"
	"github.com/gin-gonic/gin"
)

func RegisterICEndpoints(r *gin.RouterGroup) {
	r.GET("", getICs)
	r.POST("", addIC)
	r.PUT("/:ICID", updateIC)
	r.GET("/:ICID", getIC)
	r.DELETE("/:ICID", deleteIC)
	r.GET("/:ICID/configs", getConfigsOfIC)
}

func RegisterAMQPEndpoint(r *gin.RouterGroup) {
	r.POST("/:ICID/action", sendActionToIC)
}

// getICs godoc
// @Summary Get all infrastructure components
// @ID getICs
// @Tags infrastructure-components
// @Produce json
// @Success 200 {object} api.ResponseICs "ICs requested"
// @Failure 404 {object} api.ResponseError "Not found"
// @Failure 422 {object} api.ResponseError "Unprocessable entity"
// @Failure 500 {object} api.ResponseError "Internal server error"
// @Router /ic [get]
// @Security Bearer
func getICs(c *gin.Context) {

	// Checking permission is not required here since READ access is independent of user's role

	db := database.GetDB()
	var ics []database.InfrastructureComponent
	err := db.Order("ID asc").Find(&ics).Error
	if !helper.DBError(c, err) {
		c.JSON(http.StatusOK, gin.H{"ics": ics})
	}

}

// addIC godoc
// @Summary Add an infrastructure component
// @ID addIC
// @Accept json
// @Produce json
// @Tags infrastructure-components
// @Success 200 {object} api.ResponseIC "Infrastructure Component that was added"
// @Failure 400 {object} api.ResponseError "Bad request"
// @Failure 404 {object} api.ResponseError "Not found"
// @Failure 422 {object} api.ResponseError "Unprocessable entity"
// @Failure 500 {object} api.ResponseError "Internal server error"
// @Param inputIC body infrastructure_component.AddICRequest true "Infrastructure Component to be added"
// @Router /ic [post]
// @Security Bearer
func addIC(c *gin.Context) {

	ok, _ := CheckPermissions(c, database.ModelInfrastructureComponent, database.Create, false)
	if !ok {
		return
	}

	var req AddICRequest
	err := c.BindJSON(&req)
	if err != nil {
		helper.BadRequestError(c, "Error binding form data to JSON: "+err.Error())
		return
	}

	// Validate the request
	if err = req.validate(); err != nil {
		helper.UnprocessableEntityError(c, err.Error())
		return
	}

	// Check if IC to be created is managed externally
	if *req.InfrastructureComponent.ManagedExternally == true {
		// if so: refuse creation
		helper.BadRequestError(c, "create for externally managed IC not possible with this endpoint - use /ic/{ICID}/action endpoint instead to request creation of the component")
		return
	}

	// Create the new IC from the request
	newIC, err := req.createIC()
	if err != nil {
		helper.InternalServerError(c, "Unable to send create action: "+err.Error())
		return
	}

	if !(newIC.ManagedExternally) {
		// Save new IC to DB if not managed externally
		err = newIC.save()

		if helper.DBError(c, err) {
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"ic": newIC.InfrastructureComponent})
}

// updateIC godoc
// @Summary Update an infrastructure component
// @ID updateIC
// @Tags infrastructure-components
// @Accept json
// @Produce json
// @Success 200 {object} api.ResponseIC "Infrastructure Component that was updated"
// @Failure 400 {object} api.ResponseError "Bad request"
// @Failure 404 {object} api.ResponseError "Not found"
// @Failure 422 {object} api.ResponseError "Unprocessable entity"
// @Failure 500 {object} api.ResponseError "Internal server error"
// @Param inputIC body infrastructure_component.UpdateICRequest true "InfrastructureComponent to be updated"
// @Param ICID path int true "InfrastructureComponent ID"
// @Router /ic/{ICID} [put]
// @Security Bearer
func updateIC(c *gin.Context) {

	ok, oldIC := CheckPermissions(c, database.ModelInfrastructureComponent, database.Update, true)
	if !ok {
		return
	}

	if oldIC.ManagedExternally {
		helper.BadRequestError(c, "Cannot update externally managed component via API")
		return
	}

	var req UpdateICRequest
	err := c.BindJSON(&req)
	if err != nil {
		helper.BadRequestError(c, "Error binding form data to JSON: "+err.Error())
		return
	}

	// Validate the request
	if err = req.validate(); err != nil {
		helper.UnprocessableEntityError(c, err.Error())
		return
	}

	// Create the updatedIC from oldIC
	updatedIC := req.updatedIC(oldIC)

	// Finally update the IC in the DB
	err = oldIC.update(updatedIC)
	if !helper.DBError(c, err) {
		c.JSON(http.StatusOK, gin.H{"ic": updatedIC.InfrastructureComponent})
	}

}

// getIC godoc
// @Summary Get infrastructure component
// @ID getIC
// @Produce  json
// @Tags infrastructure-components
// @Success 200 {object} api.ResponseIC "Infrastructure Component that was requested"
// @Failure 400 {object} api.ResponseError "Bad request"
// @Failure 404 {object} api.ResponseError "Not found"
// @Failure 422 {object} api.ResponseError "Unprocessable entity"
// @Failure 500 {object} api.ResponseError "Internal server error"
// @Param ICID path int true "Infrastructure Component ID"
// @Router /ic/{ICID} [get]
// @Security Bearer
func getIC(c *gin.Context) {

	ok, s := CheckPermissions(c, database.ModelInfrastructureComponent, database.Read, true)
	if !ok {
		return
	}

	c.JSON(http.StatusOK, gin.H{"ic": s.InfrastructureComponent})
}

// deleteIC godoc
// @Summary Delete an infrastructure component
// @ID deleteIC
// @Tags infrastructure-components
// @Produce json
// @Success 200 {object} api.ResponseIC "Infrastructure Component that was deleted"
// @Failure 400 {object} api.ResponseError "Bad request"
// @Failure 404 {object} api.ResponseError "Not found"
// @Failure 422 {object} api.ResponseError "Unprocessable entity"
// @Failure 500 {object} api.ResponseError "Internal server error"
// @Param ICID path int true "Infrastructure Component ID"
// @Router /ic/{ICID} [delete]
// @Security Bearer
func deleteIC(c *gin.Context) {

	ok, s := CheckPermissions(c, database.ModelInfrastructureComponent, database.Delete, true)
	if !ok {
		return
	}

	// Check if IC is managed externally
	if s.ManagedExternally {
		// if so: refuse deletion
		helper.BadRequestError(c, "delete for externally managed IC not possible with this endpoint - use /ic/{ICID}/action endpoint instead to request deletion of the component")
		return
	}

	// Delete the IC
	err := s.delete()
	if helper.DBError(c, err) {
		return
	} else if err != nil {
		helper.InternalServerError(c, "Unable to send delete action: "+err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{"ic": s.InfrastructureComponent})
}

// getConfigsOfIC godoc
// @Summary Get all configurations of the infrastructure component
// @ID getConfigsOfIC
// @Tags infrastructure-components
// @Produce json
// @Success 200 {object} api.ResponseConfigs "Configs requested by user"
// @Failure 400 {object} api.ResponseError "Bad request"
// @Failure 404 {object} api.ResponseError "Not found"
// @Failure 422 {object} api.ResponseError "Unprocessable entity"
// @Failure 500 {object} api.ResponseError "Internal server error"
// @Param ICID path int true "Infrastructure Component ID"
// @Router /ic/{ICID}/configs [get]
// @Security Bearer
func getConfigsOfIC(c *gin.Context) {

	ok, s := CheckPermissions(c, database.ModelInfrastructureComponent, database.Read, true)
	if !ok {
		return
	}

	// get all associated configurations
	allConfigs, _, err := s.getConfigs()
	if !helper.DBError(c, err) {
		c.JSON(http.StatusOK, gin.H{"configs": allConfigs})
	}

}

// sendActionToIC godoc
// @Summary Send an action to IC (only available if backend server is started with -amqp parameter)
// @ID sendActionToIC
// @Tags infrastructure-components
// @Produce json
// @Param inputAction query string true "Action for IC"
// @Success 200 {object} api.ResponseError "Action sent successfully"
// @Failure 400 {object} api.ResponseError "Bad request"
// @Failure 404 {object} api.ResponseError "Not found"
// @Failure 422 {object} api.ResponseError "Unprocessable entity"
// @Failure 500 {object} api.ResponseError "Internal server error"
// @Param ICID path int true "InfrastructureComponent ID"
// @Router /ic/{ICID}/action [post]
// @Security Bearer
func sendActionToIC(c *gin.Context) {

	ok, s := CheckPermissions(c, database.ModelInfrastructureComponentAction, database.Update, true)
	if !ok {
		return
	}

	var actions []helper.Action
	err := c.BindJSON(&actions)
	if err != nil {
		helper.BadRequestError(c, "Error binding form data to JSON: "+err.Error())
		return
	}

	for _, action := range actions {
		if (action.Act == "delete" || action.Act == "create") && s.Category != "manager" {
			helper.BadRequestError(c, "cannot send a delete or create action to an IC of category "+s.Category)
		}
		err = helper.SendActionAMQP(action, s.UUID)
		if err != nil {
			helper.InternalServerError(c, "Unable to send actions to IC: "+err.Error())
			return
		}
	}
	log.Println("AMQP: Sending actions:", actions)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "OK.",
	})
}
