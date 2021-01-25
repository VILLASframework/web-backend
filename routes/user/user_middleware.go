/** User package, middleware.
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
package user

import (
	"fmt"

	"git.rwth-aachen.de/acs/public/villas/web-backend-go/configuration"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/helper"

	"github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gin-gonic/gin"
)

func userToContext(ctx *gin.Context, user_id uint) {

	var user User

	err := user.ByID(user_id)
	if helper.DBError(ctx, err) {
		return
	}

	ctx.Set(database.UserRoleCtx, user.Role)
	ctx.Set(database.UserIDCtx, user_id)
}

func Authentication() gin.HandlerFunc {

	return func(ctx *gin.Context) {

		// Authentication's access token extraction
		// XXX: if we have a multi-header for Authorization (e.g. in
		// case of OAuth2 use the request.OAuth2Extractor and make sure
		// that the argument is 'access-token' or provide a custom one
		token, err := request.ParseFromRequest(ctx.Request,
			request.MultiExtractor{
				request.AuthorizationHeaderExtractor,
				request.ArgumentExtractor{"token"},
			},
			func(token *jwt.Token) (interface{}, error) {

				// validate alg for signing the jwt
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing alg: %v",
						token.Header["alg"])
				}

				// return secret in byte format
				secret, _ := configuration.GlobalConfig.String("jwt.secret")
				return []byte(secret), nil
			})

		// If the authentication extraction fails return HTTP CODE 401
		if err != nil {
			helper.UnauthorizedAbort(ctx, "Authentication failed (claims extraction)")
			return
		}

		// If the token is ok, pass user_id to context
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

			user_id, ok := claims["id"].(float64)

			if !ok {
				helper.UnauthorizedAbort(ctx, "Authentication failed (claims casting)")
				return
			}

			userToContext(ctx, uint(user_id))
		}

	}
}
