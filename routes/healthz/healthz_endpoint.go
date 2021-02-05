/** Healthz package, endpoints.
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
package healthz

import (
	"log"
	"net/http"
	"strings"

	"git.rwth-aachen.de/acs/public/villas/web-backend-go/configuration"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/helper"
	infrastructure_component "git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/infrastructure-component"
	"github.com/gin-gonic/gin"
)

func RegisterHealthzEndpoint(r *gin.RouterGroup) {

	r.GET("", getHealth)
}

// getHealth godoc
// @Summary Get health status of backend
// @ID getHealth
// @Produce  json
// @Tags healthz
// @Success 204 string string "Backend is healthy, database and AMQP broker connections are alive"
// @Failure 500 {object} api.ResponseError "Backend is NOT healthy"
// @Router /healthz [get]
func getHealth(c *gin.Context) {

	// check if DB connection is active
	db := database.GetDB()
	err := db.DB().Ping()
	if helper.DBError(c, err) {
		return
	}

	// check if connection to AMQP broker is alive if backend was started with AMQP client
	url, err := configuration.GlobalConfig.String("amqp.host")
	if err != nil && strings.Contains(err.Error(), "Required setting 'amqp.host' not set") {
		c.JSON(http.StatusOK, gin.H{})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success:": false,
			"message":  err.Error(),
		})
		return
	}

	if len(url) != 0 {
		err = infrastructure_component.CheckConnection()
		if err != nil {
			log.Println(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"success:": false,
				"message":  err.Error(),
			})
			return
		}
	}

	// Send a 204 reponse
	c.Writer.WriteHeader(http.StatusNoContent)
}
