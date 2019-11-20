/** Helper package, utilities.
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
package helper

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetIDOfElement(c *gin.Context, elementName string, source string, providedID int) (int, error) {

	if source == "path" {
		id, err := strconv.Atoi(c.Param(elementName))
		if err != nil {
			BadRequestError(c, fmt.Sprintf("No or incorrect format of path parameter"))
			return -1, err
		}
		return id, nil
	} else if source == "query" {
		id, err := strconv.Atoi(c.Request.URL.Query().Get(elementName))
		if err != nil {
			BadRequestError(c, fmt.Sprintf("No or incorrect format of query parameter"))
			return -1, err
		}
		return id, nil
	} else if source == "body" {
		id := providedID
		return id, nil
	} else {
		return -1, fmt.Errorf("invalid source of element ID")
	}
}
