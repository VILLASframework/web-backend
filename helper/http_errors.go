/**
* This file is part of VILLASweb-backend-go
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
	"net/http"

	"github.com/gin-gonic/gin"
)

func BadRequestError(c *gin.Context, err string) {
	c.JSON(http.StatusBadRequest, gin.H{
		"success": false,
		"message": fmt.Sprintf("%v", err),
	})
}

func UnprocessableEntityError(c *gin.Context, err string) {
	c.JSON(http.StatusUnprocessableEntity, gin.H{
		"success": false,
		"message": fmt.Sprintf("%v", err),
	})
}

func InternalServerError(c *gin.Context, err string) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"success": false,
		"message": fmt.Sprintf("%v", err),
	})
}

func UnauthorizedError(c *gin.Context, err string) {
	c.JSON(http.StatusUnauthorized, gin.H{
		"success": false,
		"message": fmt.Sprintf("%v", err),
	})
}

func UnauthorizedAbort(c *gin.Context, err string) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"success": false,
		"message": fmt.Sprintf("%v", err),
	})
}

func NotFoundError(c *gin.Context, err string) {
	c.JSON(http.StatusNotFound, gin.H{
		"success": false,
		"message": fmt.Sprintf("%v", err),
	})
}

func ForbiddenError(c *gin.Context, err string) {
	c.JSON(http.StatusForbidden, gin.H{
		"success": false,
		"message": fmt.Sprintf("%v", err),
	})
}
