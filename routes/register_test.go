/** Routes package, testing
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

package routes

import (
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/configuration"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
	infrastructure_component "git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/infrastructure-component"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var router *gin.Engine

func TestMain(m *testing.M) {
	err := configuration.InitConfig()
	if err != nil {
		panic(m)
	}

	err = database.InitDB(configuration.GolbalConfig)
	if err != nil {
		panic(m)
	}
	defer database.DBpool.Close()

	router = gin.Default()

	// connect AMQP client (make sure that AMQP_HOST, AMQP_USER, AMQP_PASS are set via command line parameters)
	host, err := configuration.GolbalConfig.String("amqp.host")
	user, err := configuration.GolbalConfig.String("amqp.user")
	pass, err := configuration.GolbalConfig.String("amqp.pass")

	amqpURI := "amqp://" + user + ":" + pass + "@" + host

	err = infrastructure_component.ConnectAMQP(amqpURI)

	os.Exit(m.Run())
}

func TestRegisterEndpoints(t *testing.T) {
	database.DropTables()
	database.MigrateModels()

	api := router.Group("/api")
	RegisterEndpoints(router, api)
}

func TestAddTestData(t *testing.T) {
	resp, err := AddTestData("/api", router)
	assert.NoError(t, err, "Response body: %v", resp)
}
