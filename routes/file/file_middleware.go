/** File package, middleware.
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
package file

import (
	"fmt"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/helper"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/scenario"
	"github.com/gin-gonic/gin"
)

func CheckPermissions(c *gin.Context, operation database.CRUD) (bool, File) {

	var f File

	err := database.ValidateRole(c, database.ModelFile, operation)
	if err != nil {
		helper.UnprocessableEntityError(c, fmt.Sprintf("Access denied (role validation of file failed): %v", err.Error()))
		return false, f
	}

	fileID, err := helper.GetIDOfElement(c, "fileID", "path", -1)
	if err != nil {
		return false, f
	}

	err = f.ByID(uint(fileID))
	if helper.DBError(c, err) {
		return false, f
	}

	if operation != database.Read {
		// check access to scenario only if operation is not Read (=download) of file
		ok, _ := scenario.CheckPermissions(c, operation, "body", int(f.ScenarioID))
		if !ok {
			return false, f
		}
	}

	return true, f
}
