/** Simulator package, middleware.
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
package simulator

import (
	"fmt"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/helper"
	"github.com/gin-gonic/gin"
)

func CheckPermissions(c *gin.Context, modeltype database.ModelName, operation database.CRUD, hasID bool) (bool, Simulator) {

	var s Simulator

	err := database.ValidateRole(c, modeltype, operation)
	if err != nil {
		helper.UnprocessableEntityError(c, fmt.Sprintf("Access denied (role validation of simulator failed): %v", err.Error()))
		return false, s
	}

	if hasID {
		// Get the ID of the simulator from the context
		simulatorID, err := helper.GetIDOfElement(c, "simulatorID", "path", -1)
		if err != nil {
			return false, s
		}

		err = s.ByID(uint(simulatorID))
		if helper.DBError(c, err) {
			return false, s
		}
	}

	return true, s
}
