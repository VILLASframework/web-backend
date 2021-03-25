/** Scenario package, middleware.
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
package scenario

import (
	"fmt"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/helper"
	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
)

func CheckPermissions(c *gin.Context, operation database.CRUD, scenarioIDsource string, scenarioIDbody int) (bool, Scenario) {

	var so Scenario

	err := database.ValidateRole(c, database.ModelScenario, operation)
	if err != nil {
		helper.UnprocessableEntityError(c, fmt.Sprintf("Access denied (role validation of scenario failed): %v", err))
		return false, so
	}

	if operation == database.Create || (operation == database.Read && scenarioIDsource == "none") {
		return true, so
	}

	scenarioID, err := helper.GetIDOfElement(c, "scenarioID", scenarioIDsource, scenarioIDbody)
	if err != nil {
		return false, so
	}

	userID, _ := c.Get(database.UserIDCtx)
	userRole, _ := c.Get(database.UserRoleCtx)

	err = so.ByID(uint(scenarioID))
	if helper.DBError(c, err) {
		return false, so
	}

	if so.checkAccess(userID.(uint), userRole.(string), operation) == false {
		helper.UnprocessableEntityError(c, "Access denied (user has no access or scenario is locked).")
		return false, so
	}

	return true, so
}
