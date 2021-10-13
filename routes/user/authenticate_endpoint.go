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

func duplicateFiles(originalSo *database.Scenario, duplicateSo *database.Scenario) error {
	db := database.GetDB()
	var files []database.File
	err := db.Order("ID asc").Model(originalSo).Related(&files, "Files").Error
	if err != nil {
		log.Printf("error getting files for scenario %d", originalSo.ID)
	}

	for _, file := range files {
		var duplicateF database.File
		duplicateF.Name = file.Name
		duplicateF.Key = file.Key
		duplicateF.Type = file.Type
		duplicateF.Size = file.Size
		duplicateF.Date = file.Date
		duplicateF.ScenarioID = duplicateSo.ID
		duplicateF.FileData = file.FileData
		duplicateF.ImageHeight = file.ImageHeight
		duplicateF.ImageWidth = file.ImageWidth
		err = db.Create(&duplicateF).Error
		if err != nil {
			log.Print("error creating duplicate file")
			return err
		}
	}
	return nil
}

func duplicateDashboards(originalSo *database.Scenario, duplicateSo *database.Scenario,
	signalMap map[uint]uint, appendix string) error {

	db := database.GetDB()
	var dabs []database.Dashboard
	err := db.Order("ID asc").Model(originalSo).Related(&dabs, "Dashboards").Error
	if err != nil {
		log.Printf("error getting dashboards for scenario %d", originalSo.ID)
	}

	for _, dab := range dabs {
		var duplicateD database.Dashboard
		duplicateD.Grid = dab.Grid
		duplicateD.Name = dab.Name + appendix
		duplicateD.ScenarioID = duplicateSo.ID
		duplicateD.Height = dab.Height
		err = db.Create(&duplicateD).Error
		if err != nil {
			log.Print("error creating duplicate dashboard")
			continue
		}

		// add widgets to duplicated dashboards
		var widgets []database.Widget
		err = db.Order("ID asc").Model(&dab).Related(&widgets, "Widgets").Error
		if err != nil {
			log.Printf("error getting widgets for dashboard %d", dab.ID)
		}
		for _, widget := range widgets {
			var duplicateW database.Widget
			duplicateW.DashboardID = duplicateD.ID
			duplicateW.CustomProperties = widget.CustomProperties
			duplicateW.Height = widget.Height
			duplicateW.Width = widget.Width
			duplicateW.MinHeight = widget.MinHeight
			duplicateW.MinWidth = widget.MinWidth
			duplicateW.Name = widget.Name
			duplicateW.Type = widget.Type
			duplicateW.X = widget.X
			duplicateW.Y = widget.Y

			duplicateW.SignalIDs = []int64{}
			for _, id := range widget.SignalIDs {
				duplicateW.SignalIDs = append(duplicateW.SignalIDs, int64(signalMap[uint(id)]))
			}

			err = db.Create(&duplicateW).Error
			if err != nil {
				log.Print("error creating duplicate widget")
				continue
			}
			// associate dashboard with simulation
			err = db.Model(&duplicateD).Association("Widgets").Append(&duplicateW).Error
			if err != nil {
				log.Print("error associating duplicate widget and dashboard")
			}
		}

	}
	return nil
}

func duplicateComponentConfig(config *database.ComponentConfiguration,
	duplicateSo *database.Scenario, icIds map[uint]string, appendix string, signalMap *map[uint]uint) error {
	var configDpl database.ComponentConfiguration
	configDpl.Name = config.Name
	configDpl.StartParameters = config.StartParameters
	configDpl.ScenarioID = duplicateSo.ID
	configDpl.OutputMapping = config.OutputMapping
	configDpl.InputMapping = config.InputMapping

	db := database.GetDB()
	if icIds[config.ICID] == "" {
		configDpl.ICID = config.ICID
	} else {
		var duplicatedIC database.InfrastructureComponent
		err := db.Find(&duplicatedIC, "UUID = ?", icIds[config.ICID]).Error
		if err != nil {
			log.Print(err)
			return err
		}
		configDpl.ICID = duplicatedIC.ID
	}
	err := db.Create(&configDpl).Error
	if err != nil {
		log.Print(err)
		return err
	}

	// get all signals corresponding to component config
	var sigs []database.Signal
	err = db.Order("ID asc").Model(&config).Related(&sigs, "OutputMapping").Error
	smap := *signalMap
	for _, signal := range sigs {
		var sig database.Signal
		sig.Direction = signal.Direction
		sig.Index = signal.Index
		sig.Name = signal.Name + appendix
		sig.ScalingFactor = signal.ScalingFactor
		sig.Unit = signal.Unit
		sig.ConfigID = configDpl.ID
		err = db.Create(&sig).Error
		if err == nil {
			smap[signal.ID] = sig.ID
		}
	}

	return err
}

func duplicateScenario(so *database.Scenario, duplicateSo *database.Scenario, icIds map[uint]string, appendix string) error {
	duplicateSo.Name = so.Name + appendix
	duplicateSo.StartParameters.RawMessage = so.StartParameters.RawMessage
	db := database.GetDB()
	err := db.Create(&duplicateSo).Error
	if err != nil {
		log.Printf("Could not create duplicate of scenario %d", so.ID)
		return err
	}
	log.Print("created duplicate scenario")
	err = duplicateFiles(so, duplicateSo)
	if err != nil {
		return err
	}

	var configs []database.ComponentConfiguration
	// map existing signal IDs to duplicated signal IDs for widget duplication
	signalMap := make(map[uint]uint)
	err = db.Order("ID asc").Model(so).Related(&configs, "ComponentConfigurations").Error
	if err == nil {
		for _, config := range configs {
			err = duplicateComponentConfig(&config, duplicateSo, icIds, appendix, &signalMap)
			if err != nil {
				return err
			}
		}

	}

	err = duplicateDashboards(so, duplicateSo, signalMap, appendix)
	return err
}

func DuplicateScenarioForUser(so *database.Scenario, user *database.User) {
	go func() {

		// get all component configs of the scenario
		db := database.GetDB()
		var configs []database.ComponentConfiguration
		err := db.Order("ID asc").Model(so).Related(&configs, "ComponentConfigurations").Error
		if err != nil {
			log.Printf("Warning: scenario to duplicate (id=%d) has no component configurations", so.ID)
		}

		// iterate over component configs to check for ICs to duplicate
		duplicatedICuuids := make(map[uint]string) // key: icID; value: UUID of duplicate
		var externalUUIDs []string                 // external ICs to wait for
		for _, config := range configs {
			icID := config.ICID
			if duplicatedICuuids[icID] != "" { // this IC was already added
				continue
			}

			var ic database.InfrastructureComponent
			err = db.Find(&ic, icID).Error
			if err != nil {
				log.Printf("Cannot find IC with id %d in DB, will not duplicate for User %s", icID, user.Username)
				continue
			}

			if ic.Category == "simulator" && ic.Type == "kubernetes" {
				duplicateUUID, err := helper.RequestICcreateAMQP(&ic, ic.Manager)
				duplicatedICuuids[ic.ID] = duplicateUUID

				if err != nil { // TODO: should this function call be interrupted here?
					log.Printf("Duplication of IC (id=%d) unsuccessful", icID)
					continue
				}
				externalUUIDs = append(externalUUIDs, duplicateUUID)
			} else { // use existing IC
				duplicatedICuuids[ic.ID] = ""
				err = nil
			}
		}

		// copy scenario after all new external ICs are in DB
		icsToWaitFor := len(externalUUIDs)
		var duplicatedScenario database.Scenario
		var timeout = 5 // seconds

		for i := 0; i < timeout; i++ {
			log.Printf("i = %d", i)
			if icsToWaitFor == 0 {
				appendix := fmt.Sprintf("--%s-%d-%d", user.Username, user.ID, so.ID)
				duplicateScenario(so, &duplicatedScenario, duplicatedICuuids, appendix)

				// associate user to new scenario
				err = db.Model(&duplicatedScenario).Association("Users").Append(user).Error
				if err != nil {
					log.Printf("Could not associate User %s to scenario %d", user.Username, duplicatedScenario.ID)
				}
				log.Print("associated user to duplicated scenario")

				return
			} else {
				time.Sleep(1 * time.Second)
			}

			// check for new ICs with previously created UUIDs
			for _, uuid := range externalUUIDs {
				if uuid == "" {
					continue
				}
				log.Printf("looking for IC with UUID %s", uuid)
				var duplicatedIC database.InfrastructureComponent
				err = db.Find(&duplicatedIC, "UUID = ?", uuid).Error
				// TODO: check if not found or other error
				if err != nil {
					log.Print(err)
				} else {
					icsToWaitFor--
					uuid = ""
				}
			}
		}
	}()
}

func isAlreadyDuplicated(duplicatedName string) bool {
	db := database.GetDB()
	var scenarios []database.Scenario

	db.Find(&scenarios, "name = ?", duplicatedName)
	if len(scenarios) > 0 {
		return true
	}
	return false
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

	// check if user already exists
	err := myUser.ByUsername(username)

	if err != nil {
		// this is the first login, create new user
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
		if groupedArr, ok := configuration.ScenarioGroupMap[group]; ok {
			for _, groupedScenario := range groupedArr {
				var so database.Scenario
				err := db.Find(&so, groupedScenario.Scenario).Error
				if err != nil {
					log.Printf(`Cannot find scenario %s (id=%d) for adding/duplication.
					Affecting user %s (id=%d): %s\n`, so.Name, so.ID, myUser.Username, myUser.ID, err)
					continue
				}

				duplicateName := fmt.Sprintf("%s--%s-%d-%d", so.Name, myUser.Username, myUser.ID, so.ID)
				alreadyDuplicated := isAlreadyDuplicated(duplicateName)
				if alreadyDuplicated {
					log.Printf("Scenario %d already duplicated for user %s", so.ID, myUser.Username)
					return myUser, nil
				}

				if groupedScenario.Duplicate {
					DuplicateScenarioForUser(&so, &myUser.User)
				} else {
					err = db.Model(&so).Association("Users").Append(&(myUser.User)).Error
					if err != nil {
						log.Printf("Failed to add user %s (id=%d) to scenario %s (id=%d): %s\n", myUser.Username, myUser.ID, so.Name, so.ID, err)
						continue
					}
					log.Printf("Added user %s (id=%d) to scenario %s (id=%d)", myUser.Username, myUser.ID, so.Name, so.ID)
				}
			}
		}
	}

	return myUser, nil
}
