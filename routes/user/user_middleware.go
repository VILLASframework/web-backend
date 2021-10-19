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

func claimsToContext(c *gin.Context, claims jwt.MapClaims) error {
	userID, ok := claims["id"].(float64)
	if !ok {
		return fmt.Errorf("Authentication failed (claims casting)")
	}

	var user User

	err := user.ByID(uint(userID))
	if err != nil {
		return err
	}

	c.Set(database.UserRoleCtx, user.Role)
	c.Set(database.UserIDCtx, uint(userID))

	return nil
}

func isAuthenticated(c *gin.Context) (bool, error) {
	// Authentication's access token extraction
	token, err := request.ParseFromRequest(c.Request,
		request.MultiExtractor{
			request.AuthorizationHeaderExtractor,
			request.ArgumentExtractor{"token"},
		},
		func(token *jwt.Token) (interface{}, error) {
			// Validate alg for signing the jwt
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing alg: %v",
					token.Header["alg"])
			}

			// Return secret in byte format
			secret, _ := configuration.GlobalConfig.String("jwt.secret")
			return []byte(secret), nil
		})

	// If the authentication extraction fails return HTTP code 401
	if err != nil {
		return false, fmt.Errorf("Authentication failed (claims extraction: %s)", err)
	}

	// If the token is ok, pass user id to context
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		err = claimsToContext(c, claims)
		if err != nil {
			return false, fmt.Errorf("Authentication failed (claims casting: %s)", err)
		}
	}

	return true, nil
}

func Authentication() gin.HandlerFunc {

	return func(c *gin.Context) {
		ok, err := isAuthenticated(c)
		if !ok || err != nil {
			helper.UnauthorizedAbort(c, err.Error())
		}
	}
}
