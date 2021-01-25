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
	"strings"
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
// @Success 200 {object} api.ResponseAuthenticate "JSON web token, success status, message and authenticated user object"
// @Failure 401 {object} api.ResponseError "Unauthorized"
// @Failure 500 {object} api.ResponseError "Internal server error."
// @Router /authenticate [post]
func authenticate(c *gin.Context) {
	var user *User

	externalAuth, err := configuration.GlobalConfig.Bool("external-auth")
	if err != nil {
		helper.UnauthorizedError(c, "Backend configuration error")
		return
	}

	if err != nil || !externalAuth {
		user = authenticateStandard(c)
	} else {
		user = authenticateExternal(c)
	}

	expiresStr, err := configuration.GlobalConfig.String("jwt.expires-after")
	if err != nil {
		helper.UnauthorizedError(c, "Backend configuration error")
		return
	}

	expiresDuration, err := time.ParseDuration(expiresStr)
	if err != nil {
		helper.UnauthorizedError(c, "Backend configuration error")
		return
	}

	secret, err := configuration.GlobalConfig.String("jwt.secret")
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

func authenticateStandard(c *gin.Context) *User {
	// Bind the response (context) with the loginRequest struct
	var credentials loginRequest
	if err := c.ShouldBindJSON(&credentials); err != nil {
		helper.UnauthorizedError(c, "Wrong username or password")
		return nil
	}

	// Validate the login request
	if errs := credentials.validate(); errs != nil {
		helper.UnauthorizedError(c, "Failed to validate request")
		return nil
	}

	// Find the username in the database
	var user User
	err := user.ByUsername(credentials.Username)
	if err != nil {
		helper.UnauthorizedError(c, "Unknown username")
		return nil
	}

	// Check if this is an active user
	if !user.Active {
		helper.UnauthorizedError(c, "User is not active")
		return nil
	}

	// Validate the password
	err = user.validatePassword(credentials.Password)
	if err != nil {
		helper.UnauthorizedError(c, "Invalid password")
		return nil
	}

	return &user
}

func authenticateExternal(c *gin.Context) *User {
	username := c.Request.Header.Get("X-Forwarded-User")
	if username == "" {
		helper.UnauthorizedAbort(c, "Authentication failed (X-Forwarded-User headers)")
		return nil
	}

	email := c.Request.Header.Get("X-Forwarded-Email")
	if email == "" {
		helper.UnauthorizedAbort(c, "Authentication failed (X-Forwarded-Email headers)")
		return nil
	}

	groups := strings.Split(c.Request.Header.Get("X-Forwarded-Groups"), ",")
	// preferred_username := c.Request.Header.Get("X-Forwarded-Preferred-Username")

	var user User
	if err := user.ByUsername(username); err == nil {
		// There is already a user by this name
		return &user
	} else {
		role := "User"
		if _, found := helper.Find(groups, "admin"); found {
			role = "Admin"
		}

		newUser, err := NewUser(username, "", email, role, true)
		if err != nil {
			helper.UnauthorizedAbort(c, "Authentication failed (failed to create new user: "+err.Error()+")")
			return nil
		}

		return newUser
	}
}
