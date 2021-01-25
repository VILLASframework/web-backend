/** User package, authentication endpoint.
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
	"net/http"
	"time"

	"git.rwth-aachen.de/acs/public/villas/web-backend-go/configuration"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/helper"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

type tokenClaims struct {
	UserID uint   `json:"id"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

func RegisterAuthenticate(r *gin.RouterGroup) {
	r.POST("", authenticate)
}

// authenticate godoc
// @Summary Authentication for user
// @ID authenticate
// @Accept json
// @Produce json
// @Tags authentication
// @Param inputUser body user.loginRequest true "loginRequest of user"
// @Success 200 {object} docs.ResponseAuthenticate "JSON web token, success status, message and authenticated user object"
// @Failure 401 {object} docs.ResponseError "Unauthorized"
// @Failure 500 {object} docs.ResponseError "Internal server error."
// @Router /authenticate [post]
func authenticate(c *gin.Context) {

	// Bind the response (context) with the loginRequest struct
	var credentials loginRequest
	if err := c.ShouldBindJSON(&credentials); err != nil {
		helper.UnauthorizedError(c, "Wrong username or password")
		return
	}

	// Validate the login request
	if errs := credentials.validate(); errs != nil {
		helper.UnauthorizedError(c, "Wrong username or password")
		return
	}

	// Find the username in the database
	var user User
	err := user.ByUsername(credentials.Username)
	if err != nil {
		helper.UnauthorizedError(c, "Wrong username or password")
		return
	}

	// Check if this is an active user
	if !user.Active {
		helper.UnauthorizedError(c, "Wrong username or password")
		return
	}

	// Validate the password
	err = user.validatePassword(credentials.Password)
	if err != nil {
		helper.UnauthorizedError(c, "Wrong username or password")
		return
	}

	expiresStr, err := configuration.GolbalConfig.String("jwt.expires-after")
	if err != nil {
		helper.UnauthorizedError(c, "Backend configuration error")
		return
	}

	expiresDuration, err := time.ParseDuration(expiresStr)
	if err != nil {
		helper.UnauthorizedError(c, "Backend configuration error")
		return
	}

	secret, err := configuration.GolbalConfig.String("jwt.secret")
	if err != nil {
		helper.UnauthorizedError(c, "Backend configuration error")
		return
	}

	// create authentication token
	claims := tokenClaims{
		user.ID,
		user.Role,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expiresDuration).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "http://web.villas.fein-aachen.org/",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		helper.InternalServerError(c, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Authenticated",
		"token":   tokenString,
		"user":    user.User,
	})
}
