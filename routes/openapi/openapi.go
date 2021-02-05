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
package openapi

import (
	_ "git.rwth-aachen.de/acs/public/villas/web-backend-go/doc/api" // doc/api folder is used by Swag CLI, you have to import it
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/helper"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/swag"
)

func RegisterOpenAPIEndpoint(r *gin.RouterGroup) {
	r.GET("", getOpenAPI)
}

// getOpenAPI godoc
// @Summary Get OpenAPI 2.0 spec of API
// @ID getOpenAPI
// @Produce json
// @Tags openapi
// @Success 200 string string "A OpenAPI 2.0 specification of the API"
// @Router /openapi [get]
func getOpenAPI(c *gin.Context) {
	doc, err := swag.ReadDoc()
	if err != nil {
		helper.InternalServerError(c, err.Error())
	}

	c.Header("Content-Type", "application/json")
	c.String(200, doc)
}
