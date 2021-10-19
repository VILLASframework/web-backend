/** Healthz package, testing.
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
	"testing"

	"git.rwth-aachen.de/acs/public/villas/web-backend-go/configuration"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/helper"
	infrastructure_component "git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/infrastructure-component"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var router *gin.Engine

func TestHealthz(t *testing.T) {
	err := configuration.InitConfig()
	assert.NoError(t, err)

	// connect DB
	err = database.InitDB(configuration.GlobalConfig, true)
	assert.NoError(t, err)
	defer database.DBpool.Close()

	router = gin.Default()
	api := router.Group("/api/v2")

	RegisterHealthzEndpoint(api.Group("/healthz"))

	// close db connection
	err = database.DBpool.Close()
	assert.NoError(t, err)

	// test healthz endpoint for unconnected DB and AMQP client
	code, resp, err := helper.TestEndpoint(router, "", "/api/v2/healthz", http.MethodGet, nil)
	assert.NoError(t, err)
	assert.Equalf(t, 500, code, "Response body: \n%v\n", resp)

	// reconnect DB
	err = database.InitDB(configuration.GlobalConfig, true)
	assert.NoError(t, err)
	defer database.DBpool.Close()

	// test healthz endpoint for connected DB and unconnected AMQP client
	code, resp, err = helper.TestEndpoint(router, "", "/api/v2/healthz", http.MethodGet, nil)
	assert.NoError(t, err)
	assert.Equalf(t, 500, code, "Response body: \n%v\n", resp)

	// connect AMQP client (make sure that AMQP_HOST, AMQP_USER, AMQP_PASS are set via command line parameters)
	host, err := configuration.GlobalConfig.String("amqp.host")
	assert.NoError(t, err)
	user, err := configuration.GlobalConfig.String("amqp.user")
	assert.NoError(t, err)
	pass, err := configuration.GlobalConfig.String("amqp.pass")
	assert.NoError(t, err)

	amqpURI := "amqp://" + user + ":" + pass + "@" + host
	log.Println("AMQP URI is", amqpURI)

	err = helper.ConnectAMQP(amqpURI, infrastructure_component.ProcessMessage)
	assert.NoError(t, err)

	// test healthz endpoint for connected DB and AMQP client
	code, resp, err = helper.TestEndpoint(router, "", "/api/v2/healthz", http.MethodGet, nil)
	assert.NoError(t, err)
	assert.Equalf(t, 204, code, "Response body: \n%v\n", resp)
}
