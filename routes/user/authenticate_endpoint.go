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
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"git.rwth-aachen.de/acs/public/villas/web-backend-go/configuration"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
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
	r.GET("", authenticated)
	r.POST("/:mechanism", authenticate)
}

// authenticated godoc
// @Summary Check if user is authenticated and provide details on how the user can authenticate
// @ID authenticated
// @Accept json
// @Produce json
// @Tags authentication
// @Success 200 {object} api.ResponseAuthenticate "JSON web token, success status, message and authenticated user object"
// @Failure 401 {object} api.ResponseError "Unauthorized"
// @Failure 500 {object} api.ResponseError "Internal server error."
// @Router /authenticate [get]
func authenticated(c *gin.Context) {
	ok, err := isAuthenticated(c)
	if err != nil {
		helper.InternalServerError(c, err.Error())
		return
	}

	if ok {
		// ATTENTION: do not use c.GetInt (common.UserIDCtx) since userID is of type uint and not int
		userID, _ := c.Get(database.UserIDCtx)

		var user User
		err := user.ByID(userID.(uint))
		if helper.DBError(c, err) {
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success":       true,
			"authenticated": true,
			"user":          user.User,
		})
	} else {
		authExternal, err := configuration.GlobalConfig.Bool("auth.external.enabled")
		if err != nil {
			helper.UnauthorizedError(c, "Backend configuration error")
			return
		}

		if authExternal {
			c.JSON(http.StatusOK, gin.H{
				"success":       true,
				"authenticated": false,
				"external":      "/oauth/start",
			})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"success":       true,
				"authenticated": false,
			})
		}
	}
}

// authenticate godoc
// @Summary Authentication for user
// @ID authenticate
// @Accept json
// @Produce json
// @Tags authentication
// @Param inputUser body user.loginRequest true "loginRequest of user"
// @Param mechanism path string true "Login mechanism" Enums(internal, external)
// @Success 200 {object} api.ResponseAuthenticate "JSON web token, success status, message and authenticated user object"
// @Failure 401 {object} api.ResponseError "Unauthorized"
// @Failure 500 {object} api.ResponseError "Internal server error."
// @Router /authenticate/{mechanism} [post]
func authenticate(c *gin.Context) {
	var myUser User
	var err error

	switch c.Param("mechanism") {
	case "internal":
		myUser, err = authenticateInternal(c)
		if err != nil {
			helper.BadRequestError(c, err.Error())
		}
	case "external":
		var authExternal bool
		authExternal, err = configuration.GlobalConfig.Bool("auth.external.enabled")
		if err == nil && authExternal {
			myUser, err = authenticateExternal(c)
			if err != nil {
				helper.BadRequestError(c, err.Error())
			}
		} else {
			helper.BadRequestError(c, "External authentication is not activated")
		}
	default:
		helper.BadRequestError(c, "Invalid authentication mechanism")
	}

	// Check if this is an active user
	if !myUser.Active {
		helper.UnauthorizedError(c, "User is not active")
		return
	}

	expiresStr, err := configuration.GlobalConfig.String("jwt.expires-after")
	if err != nil {
		helper.InternalServerError(c, "Invalid backend configuration: jwt.expires-after")
		return
	}

	expiresDuration, err := time.ParseDuration(expiresStr)
	if err != nil {
		helper.InternalServerError(c, "Invalid backend configuration: jwt.expires-after")
		return
	}

	secret, err := configuration.GlobalConfig.String("jwt.secret")
	if err != nil {
		helper.InternalServerError(c, "Invalid backend configuration: jwt.secret")
		return
	}

	// Create authentication token
	claims := tokenClaims{
		myUser.ID,
		myUser.Role,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expiresDuration).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "http://web.villas.fein-aachen.org/",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		helper.InternalServerError(c, "Invalid backend configuration: jwt.secret")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Authenticated",
		"token":   tokenString,
		"user":    myUser.User,
	})
}

func authenticateInternal(c *gin.Context) (User, error) {
	// Bind the response (context) with the loginRequest struct
	var myUser User
	var credentials loginRequest
	if err := c.ShouldBindJSON(&credentials); err != nil {
		helper.UnauthorizedError(c, "Wrong username or password")
		return myUser, err
	}

	// Validate the login request
	if errs := credentials.validate(); errs != nil {
		helper.UnauthorizedError(c, "Failed to validate request")
		return myUser, errs
	}

	// Find the username in the database
	err := myUser.ByUsername(credentials.Username)
	if err != nil {
		helper.UnauthorizedError(c, "Unknown username")
		return myUser, err
	}

	// Validate the password
	err = myUser.validatePassword(credentials.Password)
	if err != nil {
		helper.UnauthorizedError(c, "Invalid password")
		return myUser, err
	}

	return myUser, nil
}

func authenticateExternal(c *gin.Context) (User, error) {
	var myUser User
	username := c.Request.Header.Get("X-Forwarded-User")
	if username == "" {
		helper.UnauthorizedAbort(c, "Authentication failed (X-Forwarded-User headers)")
		return myUser, fmt.Errorf("no username")
	}

	email := c.Request.Header.Get("X-Forwarded-Email")
	if email == "" {
		helper.UnauthorizedAbort(c, "Authentication failed (X-Forwarded-Email headers)")
		return myUser, fmt.Errorf("no email")
	}

	groups := strings.Split(c.Request.Header.Get("X-Forwarded-Groups"), ",")
	// preferred_username := c.Request.Header.Get("X-Forwarded-Preferred-Username")

	if err := myUser.ByUsername(username); err != nil {
		role := "User"
		if _, found := helper.Find(groups, "admin"); found {
			role = "Admin"
		}

		myUser, err = NewUser(username, "", email, role, true)
		if err != nil {
			helper.UnauthorizedAbort(c, "Authentication failed (failed to create new user: "+err.Error()+")")
			return myUser, fmt.Errorf("failed to create new user")
		}

		log.Printf("Created new external user %s (id=%d)", myUser.Username, myUser.ID)
	}

	// Add users to scenarios based on static map
	db := database.GetDB()
	for _, group := range groups {
		if soIDs, ok := configuration.ScenarioGroupMap[group]; ok {
			for _, soID := range soIDs {
				var so database.Scenario
				err := db.Find(&so, soID).Error
				if err != nil {
					log.Printf("Failed to add user %s (id=%d) to scenario %s (id=%d): %s\n", myUser.Username, myUser.ID, so.Name, so.ID, err)
					continue
				}

				err = db.Model(&so).Association("Users").Append(&myUser).Error
				if err != nil {
					log.Printf("Failed to add user %s (id=%d) to scenario %s (id=%d): %s\n", myUser.Username, myUser.ID, so.Name, so.ID, err)
					continue
				}

				log.Printf("Added user %s (id=%d) to scenario %s (id=%d)", myUser.Username, myUser.ID, so.Name, so.ID)
			}
		}
	}

	return myUser, nil
}
