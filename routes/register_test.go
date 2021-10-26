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
	"os"
	"testing"

	"git.rwth-aachen.de/acs/public/villas/web-backend-go/configuration"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var router *gin.Engine
var api *gin.RouterGroup

func TestMain(m *testing.M) {
	err := configuration.InitConfig()
	if err != nil {
		panic(m)
	}

	err = database.InitDB(configuration.GlobalConfig, true)
	if err != nil {
		panic(m)
	}
	defer database.DBpool.Close()

	router = gin.Default()

	api = router.Group("/api/v2")
	os.Exit(m.Run())
}

func TestRegisterEndpoints(t *testing.T) {
	database.DropTables()
	database.MigrateModels()

	RegisterEndpoints(router, api)
}

func TestAddTestData(t *testing.T) {

	err := ReadTestDataFromJson("../database/testdata.json")
	assert.NoError(t, err)

	resp, err := AddTestData(configuration.GlobalConfig, router)
	assert.NoError(t, err, "Response body: %v", resp)
}
