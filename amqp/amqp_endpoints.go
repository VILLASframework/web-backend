/** AMQP package, endpoints.
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
package amqp

import (
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/helper"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/infrastructure-component"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func RegisterAMQPEndpoint(r *gin.RouterGroup) {
	r.POST("/ICID/action", sendActionToIC)
}

// sendActionToIC godoc
// @Summary Send an action to IC (only available if backend server is started with -amqp parameter)
// @ID sendActionToIC
// @Tags AMQP
// @Produce json
// @Param inputAction query string true "Action for IC"
// @Success 200 {object} docs.ResponseError "Action sent successfully"
// @Failure 400 {object} docs.ResponseError "Bad request"
// @Failure 404 {object} docs.ResponseError "Not found"
// @Failure 422 {object} docs.ResponseError "Unprocessable entity"
// @Failure 500 {object} docs.ResponseError "Internal server error"
// @Param ICID path int true "InfrastructureComponent ID"
// @Router /ic/{ICID}/action [post]
func sendActionToIC(c *gin.Context) {

	ok, s := infrastructure_component.CheckPermissions(c, database.ModelInfrastructureComponentAction, database.Update, true)
	if !ok {
		return
	}

	var actions []Action
	err := c.BindJSON(&actions)
	if err != nil {
		helper.BadRequestError(c, "Error binding form data to JSON: "+err.Error())
		return
	}

	now := time.Now()

	for _, action := range actions {
		if action.When == 0 {
			action.When = float32(now.Unix())
		}

		err = SendActionAMQP(action, s.UUID)
		if err != nil {
			helper.InternalServerError(c, "Unable to send actions to IC: "+err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "OK.",
	})
}
