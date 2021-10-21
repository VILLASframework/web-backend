/** User package, testing.
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
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	infrastructure_component "git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/infrastructure-component"

	"git.rwth-aachen.de/acs/public/villas/web-backend-go/helper"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/stretchr/testify/assert"

	"git.rwth-aachen.de/acs/public/villas/web-backend-go/configuration"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
	component_configuration "git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/component-configuration"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/dashboard"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/file"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/scenario"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/widget"
)

var router *gin.Engine

type UserRequest struct {
	Username    string `json:"username,omitempty"`
	Password    string `json:"password,omitempty"`
	OldPassword string `json:"oldPassword,omitempty"`
	Mail        string `json:"mail,omitempty"`
	Role        string `json:"role,omitempty"`
	Active      string `json:"active,omitempty"`
}

func TestMain(m *testing.M) {
	err := configuration.InitConfig()
	if err != nil {
		panic(err)
	}
	err = database.InitDB(configuration.GlobalConfig, true)
	if err != nil {
		panic(err)
	}
	defer database.DBpool.Close()

	router = gin.Default()
	api := router.Group("/api/v2")

	RegisterAuthenticate(api.Group("/authenticate"))
	api.Use(Authentication())
	RegisterUserEndpoints(api.Group("/users"))

	scenario.RegisterScenarioEndpoints(api.Group("/scenarios"))
	infrastructure_component.RegisterICEndpoints(api.Group("/ic"))
	component_configuration.RegisterComponentConfigurationEndpoints(api.Group("/configs"))
	dashboard.RegisterDashboardEndpoints(api.Group("/dashboards"))
	file.RegisterFileEndpoints(api.Group("/files"))
	widget.RegisterWidgetEndpoints(api.Group("/widgets"))

	os.Exit(m.Run())
}

func TestAuthenticate(t *testing.T) {
	database.DropTables()
	database.MigrateModels()
	adminpw, err := database.DBAddAdminUser(configuration.GlobalConfig)
	assert.NoError(t, err)

	// try to authenticate with non JSON body
	// should result in unauthorized
	w1 := httptest.NewRecorder()
	body, _ := json.Marshal("This is no JSON")
	req, err := http.NewRequest("POST", "/api/v2/authenticate/internal", bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w1, req)
	assert.Equalf(t, 401, w1.Code, "Response body: \n%v\n", w1.Body)

	malformedCredentials := database.Credentials{
		Username: "TEST1",
	}
	// try to authenticate with malformed credentials
	// should result in unauthorized
	w2 := httptest.NewRecorder()
	body, _ = json.Marshal(malformedCredentials)
	req, err = http.NewRequest("POST", "/api/v2/authenticate/internal", bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w2, req)
	assert.Equal(t, 401, w2.Code, w2.Body)

	// try to authenticate with a username that does not exist in the DB
	// should result in unauthorized
	malformedCredentials.Username = "NOTEXIST"
	malformedCredentials.Password = "blablabla"
	w3 := httptest.NewRecorder()
	body, _ = json.Marshal(malformedCredentials)
	req, err = http.NewRequest("POST", "/api/v2/authenticate/internal", bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w3, req)
	assert.Equal(t, 401, w3.Code, w3.Body)

	// try to authenticate with a correct user name and a wrong password
	// should result in unauthorized
	malformedCredentials.Username = "User_A"
	malformedCredentials.Password = "wrong password"
	w4 := httptest.NewRecorder()
	body, _ = json.Marshal(malformedCredentials)
	req, err = http.NewRequest("POST", "/api/v2/authenticate/internal", bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w4, req)
	assert.Equal(t, 401, w4.Code, w4.Body)

	// authenticate as admin
	_, err = helper.AuthenticateForTest(router, database.Credentials{Username: "admin", Password: adminpw})
	assert.NoError(t, err)

}

func TestUserGroups(t *testing.T) {
	// Create new user
	// (user, email and groups are read from request headers in real case)
	var myUser User
	username := "Fridolin"
	email := "Fridolin@rwth-aachen.de"
	role := "User"
	userGroups := strings.Split("testGroup1,testGroup2", ",")

	err := myUser.byUsername(username)
	assert.Error(t, err)
	myUser, err = NewUser(username, "", email, role, true)
	assert.NoError(t, err)

	// Read groups file
	err = configuration.ReadGroupsFile("notexisting.yaml")
	assert.Error(t, err)

	err = configuration.ReadGroupsFile("../../configuration/groups.yaml")
	assert.NoError(t, err)

	// Check whether duplicate flag is saved correctly in configuration
	for _, group := range userGroups {
		if gsarray, ok := configuration.ScenarioGroupMap[group]; ok {
			for _, groupedScenario := range gsarray {
				if group == "testGroup1" && groupedScenario.Scenario == 1 {
					assert.Equal(t, true, groupedScenario.Duplicate)
				} else if group == "testGroup2" && groupedScenario.Scenario == 4 {
					assert.Equal(t, true, groupedScenario.Duplicate)
				} else {
					assert.Equal(t, false, groupedScenario.Duplicate)
				}
			}
		}
	}
}

func TestAuthenticateQueryToken(t *testing.T) {

	database.DropTables()
	database.MigrateModels()
	adminpw, err := database.DBAddAdminUser(configuration.GlobalConfig)
	assert.NoError(t, err)

	// authenticate as admin
	token, err := helper.AuthenticateForTest(router, database.Credentials{Username: "admin", Password: adminpw})
	assert.NoError(t, err)

	w := httptest.NewRecorder()

	// Create the request
	req, err := http.NewRequest("GET", "/api/v2/users?token="+token, nil)
	assert.NoError(t, err)

	router.ServeHTTP(w, req)
	assert.Equal(t, w.Code, 200)
}

func TestAddGetUser(t *testing.T) {

	database.DropTables()
	database.MigrateModels()
	adminpw, err := database.DBAddAdminUser(configuration.GlobalConfig)
	assert.NoError(t, err)

	// authenticate as admin
	token, err := helper.AuthenticateForTest(router, database.Credentials{Username: "admin", Password: adminpw})
	assert.NoError(t, err)

	// try to POST with non JSON body
	// should result in bad request
	code, resp, err := helper.TestEndpoint(router, token,
		"/api/v2/users", "POST", "This is not JSON")
	assert.NoError(t, err)
	assert.Equalf(t, 400, code, "Response body: \n%v\n", resp)

	wrongUser := UserRequest{
		Username: "ab",
		Password: "123456",
		Mail:     "bla@test.com",
		Role:     "Guest",
	}
	// try POST with too short username
	// should result in bad request
	code, resp, err = helper.TestEndpoint(router, token,
		"/api/v2/users", "POST", helper.KeyModels{"user": wrongUser})
	assert.NoError(t, err)
	assert.Equalf(t, 400, code, "Response body: \n%v\n", resp)

	// try POST with too short password
	// should result in bad request
	wrongUser.Username = "Longenoughusername"
	wrongUser.Password = "short"
	code, resp, err = helper.TestEndpoint(router, token,
		"/api/v2/users", "POST", helper.KeyModels{"user": wrongUser})
	assert.NoError(t, err)
	assert.Equalf(t, 400, code, "Response body: \n%v\n", resp)

	// try POST with wrong role
	// should result in bad request
	wrongUser.Password = "longenough"
	wrongUser.Role = "ThisIsNotARole"
	code, resp, err = helper.TestEndpoint(router, token,
		"/api/v2/users", "POST", helper.KeyModels{"user": wrongUser})
	assert.NoError(t, err)
	assert.Equalf(t, 400, code, "Response body: \n%v\n", resp)

	//try POST with wrong email
	// should result in bad request
	wrongUser.Mail = "noemailaddress"
	wrongUser.Role = "Guest"
	code, resp, err = helper.TestEndpoint(router, token,
		"/api/v2/users", "POST", helper.KeyModels{"user": wrongUser})
	assert.NoError(t, err)
	assert.Equalf(t, 400, code, "Response body: \n%v\n", resp)

	// try POST with missing required fields
	// should result in bad request
	wrongUser.Mail = ""
	wrongUser.Role = ""
	wrongUser.Password = ""
	code, resp, err = helper.TestEndpoint(router, token,
		"/api/v2/users", "POST", helper.KeyModels{"user": wrongUser})
	assert.NoError(t, err)
	assert.Equalf(t, 400, code, "Response body: \n%v\n", resp)

	// try POST with username that is already taken
	// should result in unprocessable entity
	wrongUser.Mail = "test@test.com"
	wrongUser.Role = "Guest"
	wrongUser.Password = "blablabla"
	wrongUser.Username = "admin"
	code, resp, err = helper.TestEndpoint(router, token,
		"/api/v2/users", "POST", helper.KeyModels{"user": wrongUser})
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

	newUser := UserRequest{
		Username: "Alice483",
		Password: "th1s_I5_@lice#",
		Mail:     "mail@domain.com",
		Role:     "User",
	}

	// test POST user/ $newUser
	code, resp, err = helper.TestEndpoint(router, token,
		"/api/v2/users", "POST", helper.KeyModels{"user": newUser})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Turn password member of newUser to empty string so it is omitted
	// in marshaling. The password will never be included in the
	// response and if is non empty in request we will not be able to do
	// request-response comparison. Reserve it also for later usage.
	newUser.Password = ""

	// Compare POST's response with the newUser (Password omitted)
	err = helper.CompareResponse(resp, helper.KeyModels{"user": newUser})
	assert.NoError(t, err)

	// Read newUser's ID from the response
	newUserID, err := helper.GetResponseID(resp)
	assert.NoError(t, err)

	// Get the newUser
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/users/%v", newUserID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare GET's response with the newUser (Password omitted)
	err = helper.CompareResponse(resp, helper.KeyModels{"user": newUser})
	assert.NoError(t, err)

	// try to GET user with invalid user ID
	// should result in bad request
	code, resp, err = helper.TestEndpoint(router, token,
		"/api/v2/users/bla", "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 400, code, "Response body: \n%v\n", resp)
}

func TestUsersNotAllowedActions(t *testing.T) {

	database.DropTables()
	database.MigrateModels()
	adminpw, err := database.DBAddAdminUser(configuration.GlobalConfig)
	assert.NoError(t, err)

	// authenticate as admin
	token, err := helper.AuthenticateForTest(router, database.Credentials{Username: "admin", Password: adminpw})
	assert.NoError(t, err)

	// Add a user
	newUser := UserRequest{
		Username: "NotAllowedActions",
		Password: "N0t_@LLow3d_act10n5",
		Mail:     "not@allowed.ac",
		Role:     "User",
	}
	code, resp, err := helper.TestEndpoint(router, token,
		"/api/v2/users", "POST", helper.KeyModels{"user": newUser})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	newUserID, err := helper.GetResponseID(resp)
	assert.NoError(t, err)

	// Authenticate as the added user
	token, err = helper.AuthenticateForTest(router, UserRequest{
		Username: newUser.Username,
		Password: newUser.Password,
	})
	assert.NoError(t, err)

	// Try to get all the users (NOT ALLOWED)
	code, resp, err = helper.TestEndpoint(router, token,
		"/api/v2/users", "POST", helper.KeyModels{"user": newUser})
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

	// Try to read another user than self (eg. Admin) (NOT ALLOWED)
	code, resp, err = helper.TestEndpoint(router, token,
		"/api/v2/users/0", "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 403, code, "Response body: \n%v\n", resp)

	// Try to delete another user (eg. Admin) (NOT ALLOWED)
	code, resp, err = helper.TestEndpoint(router, token,
		"/api/v2/users/0", "DELETE", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

	// Try to delete self (NOT ALLOWED)
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/users/%v", newUserID), "DELETE", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

}

func TestGetAllUsers(t *testing.T) {

	database.DropTables()
	database.MigrateModels()
	adminpw, err := database.DBAddAdminUser(configuration.GlobalConfig)
	assert.NoError(t, err)

	// authenticate as admin
	token, err := helper.AuthenticateForTest(router, database.Credentials{Username: "admin", Password: adminpw})
	assert.NoError(t, err)

	// get the length of the GET all users response
	initialNumber, err := helper.LengthOfResponse(router, token,
		"/api/v2/users", "GET", nil)
	assert.NoError(t, err)

	// Add a user
	newUser := UserRequest{
		Username: "UserGetAllUsers",
		Password: "get@ll_User5",
		Mail:     "get@all.users",
		Role:     "User",
	}
	code, resp, err := helper.TestEndpoint(router, token,
		"/api/v2/users", "POST", helper.KeyModels{"user": newUser})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// get the length of the GET all users response again
	finalNumber, err := helper.LengthOfResponse(router, token,
		"/api/v2/users", "GET", nil)
	assert.NoError(t, err)

	assert.Equal(t, finalNumber, initialNumber+1)

	newUserCredentials := database.Credentials{
		Username: newUser.Username,
		Password: newUser.Password,
	}

	// authenticate as normal user
	token, err = helper.AuthenticateForTest(router, newUserCredentials)
	assert.NoError(t, err)

	// try to get all users as normal user
	// should result in unprocessable entity eror
	code, resp, err = helper.TestEndpoint(router, token,
		"/api/v2/users", "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

}

func TestModifyAddedUserAsUser(t *testing.T) {

	database.DropTables()
	database.MigrateModels()
	adminpw, err := database.DBAddAdminUser(configuration.GlobalConfig)
	assert.NoError(t, err)

	// authenticate as admin
	token, err := helper.AuthenticateForTest(router, database.Credentials{Username: "admin", Password: adminpw})
	assert.NoError(t, err)

	// Add a user that will modify itself
	newUser := UserRequest{
		Username: "modMyself",
		Password: "mod_My5elf",
		Mail:     "mod@my.self",
		Role:     "User",
	}
	code, resp, err := helper.TestEndpoint(router, token,
		"/api/v2/users", "POST", helper.KeyModels{"user": newUser})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	newUserID, err := helper.GetResponseID(resp)
	assert.NoError(t, err)

	// authenticate as the new user
	token, err = helper.AuthenticateForTest(router, UserRequest{
		Username: newUser.Username,
		Password: newUser.Password,
	})
	assert.NoError(t, err)

	// Try PUT with invalid user ID in path
	// Should return a bad request
	code, resp, err = helper.TestEndpoint(router, token, "/api/v2/users/blabla", "PUT",
		helper.KeyModels{"user": newUser})
	assert.NoError(t, err)
	assert.Equalf(t, 400, code, "Response body: \n%v\n", resp)

	// Try to PUT a non JSON body
	// Should return a bad request
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/users/%v", newUserID), "PUT", "This is no JSON")
	assert.NoError(t, err)
	assert.Equalf(t, 400, code, "Response body: \n%v\n", resp)

	// Turn password member of newUser to empty string so it is omitted
	// in marshaling. The password will never be included in the
	// response and if is non empty in request we will not be able to do
	// request-response comparison
	newUser.Password = ""

	// try to modify active state of user
	// should result in forbidden
	modActiveState := UserRequest{Active: "no"}
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/users/%v", newUserID), "PUT", helper.KeyModels{"user": modActiveState})
	assert.NoError(t, err)
	assert.Equalf(t, 403, code, "Response body: \n%v\n", resp)

	// modify newUser's own name
	modRequest := UserRequest{Username: "myNewName"}
	newUser.Username = modRequest.Username
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/users/%v", newUserID), "PUT",
		helper.KeyModels{"user": modRequest})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)
	err = helper.CompareResponse(resp, helper.KeyModels{"user": newUser})
	assert.NoError(t, err)

	// modify Admin's name (ILLEGAL)
	modRequest = UserRequest{Username: "myNewName"}
	newUser.Username = modRequest.Username
	code, resp, err = helper.TestEndpoint(router, token,
		"/api/v2/users/1", "PUT", helper.KeyModels{"user": modRequest})
	assert.NoError(t, err)
	assert.Equalf(t, 403, code, "Response body: \n%v\n", resp)

	// modify newUser's own email
	modRequest = UserRequest{Mail: "my@new.email"}
	newUser.Mail = modRequest.Mail
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/users/%v", newUserID), "PUT",
		helper.KeyModels{"user": modRequest})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)
	err = helper.CompareResponse(resp, helper.KeyModels{"user": newUser})
	assert.NoError(t, err)

	// modify Admin's own email (ILLEGAL)
	modRequest = UserRequest{Mail: "my@new.email"}
	newUser.Mail = modRequest.Mail
	code, resp, err = helper.TestEndpoint(router, token,
		"/api/v2/users/1", "PUT", helper.KeyModels{"user": modRequest})
	assert.NoError(t, err)
	assert.Equalf(t, 403, code, "Response body: \n%v\n", resp)

	// modify newUser's role (ILLEGAL)
	modRequest = UserRequest{Role: "Admin"}
	newUser.Role = modRequest.Role
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/users/%v", newUserID), "PUT",
		helper.KeyModels{"user": modRequest})
	assert.NoError(t, err)
	assert.Equalf(t, 403, code, "Response body: \n%v\n", resp)

	// modify newUser's password without providing old password
	modRequest = UserRequest{
		Password: "5tr0ng_pw!",
	}
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/users/%v", newUserID), "PUT",
		helper.KeyModels{"user": modRequest})
	assert.NoError(t, err)
	assert.Equalf(t, 403, code, "Response body: \n%v\n", resp)

	// modify newUser's password with wring old password
	modRequest = UserRequest{
		Password:    "5tr0ng_pw!",
		OldPassword: "wrongoldpassword",
	}
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/users/%v", newUserID), "PUT",
		helper.KeyModels{"user": modRequest})
	assert.NoError(t, err)
	assert.Equalf(t, 403, code, "Response body: \n%v\n", resp)

	// modify newUser's password
	modRequest = UserRequest{
		Password:    "5tr0ng_pw!",
		OldPassword: "mod_My5elf",
	}
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/users/%v", newUserID), "PUT",
		helper.KeyModels{"user": modRequest})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// try to login as newUser with the modified username and password
	token, err = helper.AuthenticateForTest(router, UserRequest{
		Username: newUser.Username,
		Password: modRequest.Password,
	})
	assert.NoError(t, err)

	// modify Admin's password (ILLEGAL)
	modRequest = UserRequest{Password: "4dm1ns_pw!"}
	code, resp, err = helper.TestEndpoint(router, token,
		"/api/v2/users/1", "PUT", helper.KeyModels{"user": modRequest})
	assert.NoError(t, err)
	assert.Equalf(t, 403, code, "Response body: \n%v\n", resp)
}

func TestInvalidUserUpdate(t *testing.T) {

	database.DropTables()
	database.MigrateModels()
	adminpw, err := database.DBAddAdminUser(configuration.GlobalConfig)
	assert.NoError(t, err)

	// authenticate as admin
	token, err := helper.AuthenticateForTest(router, database.Credentials{Username: "admin", Password: adminpw})
	assert.NoError(t, err)

	// Add a user
	newUser := UserRequest{
		Username: "invalidUpdatedUser",
		Password: "wr0ng_Upd@te!",
		Mail:     "inv@user.upd",
		Role:     "User",
	}
	code, resp, err := helper.TestEndpoint(router, token,
		"/api/v2/users", "POST", helper.KeyModels{"user": newUser})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	newUserID, err := helper.GetResponseID(resp)
	assert.NoError(t, err)

	// try PUT with userID that does not exist
	// should result in not found
	modRequest := UserRequest{
		Password: "longenough",
	}
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/users/%v", newUserID+1), "PUT",
		helper.KeyModels{"user": modRequest})
	assert.NoError(t, err)
	assert.Equalf(t, 404, code, "Response body: \n%v\n", resp)

	// try to PUT with username that is already taken
	// should result in bad request
	modRequest.Password = ""
	modRequest.Username = "admin"
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/users/%v", newUserID), "PUT",
		helper.KeyModels{"user": modRequest})
	assert.NoError(t, err)
	assert.Equalf(t, 400, code, "Response body: \n%v\n", resp)

	// modify newUser's password with INVALID password
	modRequest.Password = "short"
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/users/%v", newUserID), "PUT",
		helper.KeyModels{"user": modRequest})
	assert.NoError(t, err)
	assert.Equalf(t, 400, code, "Response body: \n%v\n", resp)

	// modify newUser's email with INVALID email
	modRequest = UserRequest{Mail: "notEmail"}
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/users/%v", newUserID), "PUT",
		helper.KeyModels{"user": modRequest})
	assert.NoError(t, err)
	assert.Equalf(t, 400, code, "Response body: \n%v\n", resp)

	// modify newUser's role with INVALID role
	modRequest = UserRequest{Role: "noRole"}
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/users/%v", newUserID), "PUT",
		helper.KeyModels{"user": modRequest})
	assert.NoError(t, err)
	assert.Equalf(t, 400, code, "Response body: \n%v\n", resp)
}

func TestModifyAddedUserAsAdmin(t *testing.T) {

	database.DropTables()
	database.MigrateModels()
	adminpw, err := database.DBAddAdminUser(configuration.GlobalConfig)
	assert.NoError(t, err)

	// authenticate as admin
	token, err := helper.AuthenticateForTest(router, database.Credentials{Username: "admin", Password: adminpw})
	assert.NoError(t, err)

	// Add a user
	newUser := UserRequest{
		Username: "modAddedUser",
		Password: "mod_4d^2ed_0ser",
		Mail:     "mod@added.user",
		Role:     "User",
	}
	code, resp, err := helper.TestEndpoint(router, token,
		"/api/v2/users", "POST", helper.KeyModels{"user": newUser})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	newUserID, err := helper.GetResponseID(resp)
	assert.NoError(t, err)

	// Turn password member of newUser to empty string so it is omitted
	// in marshaling. The password will never be included in the
	// response and if is non empty in request we will not be able to do
	// request-response comparison
	newUser.Password = ""

	// modify newUser's name
	modRequest := UserRequest{Username: "NewUsername"}
	newUser.Username = modRequest.Username
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/users/%v", newUserID), "PUT",
		helper.KeyModels{"user": modRequest})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)
	err = helper.CompareResponse(resp, helper.KeyModels{"user": newUser})
	assert.NoError(t, err)

	// modify newUser's email
	modRequest = UserRequest{Mail: "new@e.mail"}
	newUser.Mail = modRequest.Mail
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/users/%v", newUserID), "PUT",
		helper.KeyModels{"user": modRequest})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)
	err = helper.CompareResponse(resp, helper.KeyModels{"user": newUser})
	assert.NoError(t, err)

	// modify newUser's role
	modRequest = UserRequest{Role: "Admin"}
	newUser.Role = modRequest.Role
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/users/%v", newUserID), "PUT",
		helper.KeyModels{"user": modRequest})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)
	err = helper.CompareResponse(resp, helper.KeyModels{"user": newUser})
	assert.NoError(t, err)

	// modify newUser's password, should not work without admin password
	modRequest = UserRequest{
		Password: "4_g00d_pw!",
	}
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/users/%v", newUserID), "PUT",
		helper.KeyModels{"user": modRequest})
	assert.NoError(t, err)
	assert.Equalf(t, 403, code, "Response body: \n%v\n", resp)

	// modify newUser's password, requires admin password
	modRequest = UserRequest{
		Password:    "4_g00d_pw!",
		OldPassword: adminpw,
	}
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/users/%v", newUserID), "PUT",
		helper.KeyModels{"user": modRequest})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// try to login as newUser with the modified username and password
	_, err = helper.AuthenticateForTest(router, UserRequest{
		Username: newUser.Username,
		Password: modRequest.Password,
	})
	assert.NoError(t, err)

	// authenticate as admin
	token, err = helper.AuthenticateForTest(router, database.Credentials{Username: "admin", Password: adminpw})
	assert.NoError(t, err)

	// modify newUser's Active status
	modRequest = UserRequest{Active: "no"}
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/users/%v", newUserID), "PUT",
		helper.KeyModels{"user": modRequest})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// try to login as newUser with the modified active status
	// should NOT work anymore!
	code, resp, err = helper.TestEndpoint(router, token,
		"/api/v2/authenticate/internal", "POST",
		UserRequest{
			Username: newUser.Username,
			Password: "4_g00d_pw!",
		})
	assert.NoError(t, err)
	assert.Equalf(t, 401, code, "Response body: \n%v\n", resp)
}

func TestDeleteUser(t *testing.T) {

	database.DropTables()
	database.MigrateModels()
	adminpw, err := database.DBAddAdminUser(configuration.GlobalConfig)
	assert.NoError(t, err)

	// authenticate as admin
	token, err := helper.AuthenticateForTest(router, database.Credentials{Username: "admin", Password: adminpw})
	assert.NoError(t, err)

	// Add a user
	newUser := UserRequest{
		Username: "toBeDeletedUser",
		Password: "f0r_deletIOn_0ser",
		Mail:     "to@be.deleted",
		Role:     "User",
	}
	code, resp, err := helper.TestEndpoint(router, token,
		"/api/v2/users", "POST", helper.KeyModels{"user": newUser})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	newUserID, err := helper.GetResponseID(resp)
	assert.NoError(t, err)

	// try to DELETE with invalid ID
	// should result in bad request
	code, resp, err = helper.TestEndpoint(router, token,
		"/api/v2/users/bla", "DELETE", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 400, code, "Response body: \n%v\n", resp)

	// try to DELETE with ID that does not exist
	// should result in not found
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/users/%v", newUserID+1), "DELETE", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 404, code, "Response body: \n%v\n", resp)

	// Count the number of all the users returned
	initialNumber, err := helper.LengthOfResponse(router, token,
		"/api/v2/users", "GET", nil)
	assert.NoError(t, err)

	// Delete the added newUser
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/users/%v", newUserID), "DELETE", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Again count the number of all the users returned
	finalNumber, err := helper.LengthOfResponse(router, token,
		"/api/v2/users", "GET", nil)
	assert.NoError(t, err)

	assert.Equal(t, finalNumber, initialNumber-1)
}

func TestDuplicateScenarioForUser(t *testing.T) {

	database.DropTables()
	database.MigrateModels()
	assert.NoError(t, database.AddTestUsers())

	// connect AMQP client
	// Make sure that AMQP_HOST, AMQP_USER, AMQP_PASS are set
	host, _ := configuration.GlobalConfig.String("amqp.host")
	usr, _ := configuration.GlobalConfig.String("amqp.user")
	pass, _ := configuration.GlobalConfig.String("amqp.pass")
	amqpURI := "amqp://" + usr + ":" + pass + "@" + host

	// AMQP Connection startup is tested here
	// Not repeated in other tests because it is only needed once
	session = helper.NewAMQPSession("villas-test-session", amqpURI, "villas", infrastructure_component.ProcessMessage)

	// authenticate as admin (needed to create original IC)
	token, err := helper.AuthenticateForTest(router, database.AdminCredentials)
	assert.NoError(t, err)

	/*** Create original scenario and entities ***/
	scenarioID := addScenario(token)
	var originalSo database.Scenario
	db := database.GetDB()
	err = db.Find(&originalSo, scenarioID).Error
	assert.NoError(t, err)

	// add file to scenario
	fileID, err := addFile(scenarioID, token)
	assert.NoError(t, err)
	assert.NotEqual(t, 99, fileID)

	// add IC
	err = addIC(session, token)
	assert.NoError(t, err)

	// wait for IC to be created asynchronously
	time.Sleep(4 * time.Second)

	// add component config
	err = addComponentConfig(scenarioID, token)
	assert.NoError(t, err)

	// TODO: add signals

	// add dashboards to scenario
	err = addTwoDashboards(scenarioID, token)
	assert.NoError(t, err)
	var dashboards []database.Dashboard
	err = db.Find(&dashboards).Error
	assert.NoError(t, err)
	assert.Equal(t, 2, len(dashboards))

	// add widgets
	dashboardID_forAddingWidget := uint(0)
	err = addWidget(dashboardID_forAddingWidget, token)
	assert.NoError(t, err)

	/*** Duplicate scenario for new user ***/
	username := "Schnittlauch"
	myUser, err := NewUser(username, "", "schnitti@lauch.de", "User", true)
	assert.NoError(t, err)

	// create fake IC duplicate (would normally be created by villas-controller)
	err = addFakeIC(session, token)
	assert.NoError(t, err)

	SetAMQPSession(session)
	if err := <-duplicateScenarioForUser(originalSo, &myUser.User); err != nil {
		t.Fail()
	}

	/*** Check duplicated scenario for correctness ***/
	var dplScenarios []database.Scenario
	err = db.Find(&dplScenarios, "name = ?", originalSo.Name+" "+username).Error
	assert.NoError(t, err)
	assert.Equal(t, 1, len(dplScenarios))
	assert.Equal(t, originalSo.StartParameters, dplScenarios[0].StartParameters)

	// compare original and duplicated file
	var files []database.File
	err = db.Find(&files).Error
	assert.NoError(t, err)
	assert.Equal(t, 2, len(files))
	assert.Equal(t, files[0].Name, files[1].Name)
	assert.Equal(t, files[0].FileData, files[1].FileData)
	assert.Equal(t, files[0].ImageHeight, files[1].ImageHeight)
	assert.Equal(t, files[0].ImageWidth, files[1].ImageWidth)
	assert.Equal(t, files[0].Size, files[1].Size)
	assert.NotEqual(t, files[0].ScenarioID, files[1].ScenarioID)
	assert.NotEqual(t, files[0].ID, files[1].ID)

	// compare original and duplicated component config
	var configs []database.ComponentConfiguration
	err = db.Find(&configs).Error
	assert.NoError(t, err)
	assert.Equal(t, 2, len(configs))
	assert.Equal(t, configs[0].Name, configs[1].Name)
	assert.Equal(t, configs[0].FileIDs, configs[1].FileIDs)
	assert.Equal(t, configs[0].InputMapping, configs[1].InputMapping)
	assert.Equal(t, configs[0].OutputMapping, configs[1].OutputMapping)
	assert.Equal(t, configs[0].StartParameters, configs[1].StartParameters)
	assert.NotEqual(t, configs[0].ScenarioID, configs[1].ScenarioID)
	assert.NotEqual(t, configs[0].ICID, configs[1].ICID)
	assert.NotEqual(t, configs[0].ID, configs[1].ID)

	// compare original and duplicated infrastructure component

	// compare original and duplicated dashboards
	err = db.Order("created_at asc").Find(&dashboards).Error
	assert.NoError(t, err)
	assert.Equal(t, 4, len(dashboards))
	assert.Equal(t, dashboards[0].Name, dashboards[2].Name)
	assert.Equal(t, dashboards[0].Grid, dashboards[2].Grid)
	assert.Equal(t, dashboards[0].Height, dashboards[2].Height)
	assert.NotEqual(t, dashboards[0].ScenarioID, dashboards[2].ScenarioID)
	assert.NotEqual(t, dashboards[0].ID, dashboards[2].ID)

	assert.Equal(t, dashboards[1].Name, dashboards[3].Name)
	assert.Equal(t, dashboards[1].Grid, dashboards[3].Grid)
	assert.Equal(t, dashboards[1].Height, dashboards[3].Height)
	assert.NotEqual(t, dashboards[1].ScenarioID, dashboards[3].ScenarioID)
	assert.NotEqual(t, dashboards[1].ID, dashboards[3].ID)

	// compare original and duplicated widget
	var widgets []database.Widget
	err = db.Find(&widgets).Error
	assert.NoError(t, err)
	assert.Equal(t, 2, len(widgets))
	assert.Equal(t, widgets[0].Name, widgets[1].Name)
	assert.Equal(t, widgets[0].CustomProperties, widgets[1].CustomProperties)
	assert.Equal(t, widgets[0].SignalIDs, widgets[1].SignalIDs)
	assert.Equal(t, widgets[0].MinHeight, widgets[1].MinHeight)
	assert.Equal(t, widgets[0].MinWidth, widgets[1].MinWidth)
	assert.NotEqual(t, widgets[0].DashboardID, widgets[1].DashboardID)
	assert.NotEqual(t, widgets[0].ID, widgets[1].ID)

}

type ScenarioRequest struct {
	Name            string         `json:"name,omitempty"`
	StartParameters postgres.Jsonb `json:"startParameters,omitempty"`
}

func addScenario(token string) (scenarioID uint) {
	newScenario := ScenarioRequest{
		Name:            "Scenario1",
		StartParameters: postgres.Jsonb{json.RawMessage(`{"parameter1" : "testValue1A", "parameter2" : "testValue2A", "parameter3" : 42}`)},
	}
	_, resp, err := helper.TestEndpoint(router, token,
		"/api/v2/scenarios", "POST", helper.KeyModels{"scenario": newScenario})
	if err != nil {
		log.Panic("The following error happend on POSTing a scenario: ", err.Error())
	}

	// Read newScenario's ID from the response
	newScenarioID, _ := helper.GetResponseID(resp)

	// add the guest user to the new scenario
	_, _, _ = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/v2/scenarios/%v/user?username=User_A", newScenarioID), "PUT", nil)

	return uint(newScenarioID)
}

func addFile(scenarioID uint, token string) (uint, error) {
	c1 := []byte(`<?xml version="1.0"?>
	<svg xmlns="http://www.w3.org/2000/svg"
			 width="400" height="400">
		<circle cx="100" cy="100" r="50" stroke="black"
			stroke-width="5" fill="red" />
	</svg>`)
	err := ioutil.WriteFile("circle.svg", c1, 0644)
	if err != nil {
		return 99, err
	}

	// test POST files
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	fileWriter, err := bodyWriter.CreateFormFile("file", "circle.svg")
	if err != nil {
		return 99, err
	}

	// open file handle
	fh, err := os.Open("circle.svg")
	if err != nil {
		return 99, err
	}
	defer fh.Close()

	// io copy
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		return 99, err
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	// Create the request
	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", fmt.Sprintf("/api/v2/files?scenarioID=%v", scenarioID), bodyBuf)
	if err != nil {
		return 99, err
	}
	req.Header.Set("Content-Type", contentType)
	req.Header.Add("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	newFileID, err := helper.GetResponseID(w.Body)
	if err != nil {
		return 99, err
	}
	return uint(newFileID), nil
}

type DashboardRequest struct {
	Name       string `json:"name,omitempty"`
	Grid       int    `json:"grid,omitempty"`
	Height     int    `json:"height,omitempty"`
	ScenarioID uint   `json:"scenarioID,omitempty"`
}

func addTwoDashboards(scenarioID uint, token string) error {
	newDashboardA := DashboardRequest{
		Name:       "Dashboard_A",
		Grid:       15,
		Height:     200,
		ScenarioID: scenarioID,
	}

	newDashboardB := DashboardRequest{
		Name:       "Dashboard_B",
		Grid:       35,
		Height:     555,
		ScenarioID: scenarioID,
	}

	_, _, err := helper.TestEndpoint(router, token,
		"/api/v2/dashboards", "POST", helper.KeyModels{"dashboard": newDashboardA})
	if err != nil {
		return err
	}

	_, _, err = helper.TestEndpoint(router, token,
		"/api/v2/dashboards", "POST", helper.KeyModels{"dashboard": newDashboardB})

	return err
}

type WidgetRequest struct {
	Name             string         `json:"name,omitempty"`
	Type             string         `json:"type,omitempty"`
	Width            uint           `json:"width,omitempty"`
	Height           uint           `json:"height,omitempty"`
	MinWidth         uint           `json:"minWidth,omitempty"`
	MinHeight        uint           `json:"minHeight,omitempty"`
	X                int            `json:"x,omitempty"`
	Y                int            `json:"y,omitempty"`
	Z                int            `json:"z,omitempty"`
	DashboardID      uint           `json:"dashboardID,omitempty"`
	IsLocked         bool           `json:"isLocked,omitempty"`
	CustomProperties postgres.Jsonb `json:"customProperties,omitempty"`
	SignalIDs        []int64        `json:"signalIDs,omitempty"`
}

func addWidget(dashboardID uint, token string) error {
	newWidget := WidgetRequest{
		Name:             "My label",
		Type:             "Label",
		Width:            100,
		Height:           50,
		MinWidth:         40,
		MinHeight:        80,
		X:                10,
		Y:                10,
		Z:                200,
		IsLocked:         false,
		CustomProperties: postgres.Jsonb{RawMessage: json.RawMessage(`{"textSize" : "20", "fontColor" : "#4287f5", "fontColor_opacity": 1}`)},
		SignalIDs:        []int64{},
	}

	newWidget.DashboardID = dashboardID

	_, _, err := helper.TestEndpoint(router, token,
		"/api/v2/widgets", "POST", helper.KeyModels{"widget": newWidget})

	return err
}

type ICRequest struct {
	UUID                  string         `json:"uuid,omitempty"`
	WebsocketURL          string         `json:"websocketurl,omitempty"`
	APIURL                string         `json:"apiurl,omitempty"`
	Type                  string         `json:"type,omitempty"`
	Name                  string         `json:"name,omitempty"`
	Category              string         `json:"category,omitempty"`
	State                 string         `json:"state,omitempty"`
	Location              string         `json:"location,omitempty"`
	Description           string         `json:"description,omitempty"`
	StartParameterSchema  postgres.Jsonb `json:"startparameterschema,omitempty"`
	CreateParameterSchema postgres.Jsonb `json:"createparameterschema,omitempty"`
	ManagedExternally     *bool          `json:"managedexternally"`
	Manager               string         `json:"manager,omitempty"`
}

func newTrue() *bool {
	b := true
	return &b
}

func addIC(session *helper.AMQPsession, token string) error {
	// create IC
	var newIC = ICRequest{
		UUID:                  "7be0322d-354e-431e-84bd-ae4c9633138b",
		WebsocketURL:          "https://villas.k8s.eonerc.rwth-aachen.de/ws/ws_sig",
		APIURL:                "https://villas.k8s.eonerc.rwth-aachen.de/ws/api/v2",
		Type:                  "kubernetes",
		Name:                  "Kubernetes Simulator",
		Category:              "simulator",
		State:                 "idle",
		Location:              "k8s",
		Description:           "A kubernetes simulator for testing purposes",
		StartParameterSchema:  postgres.Jsonb{json.RawMessage(`{"startprop1" : "a nice prop"}`)},
		CreateParameterSchema: postgres.Jsonb{json.RawMessage(`{"createprop1" : "a really nice prop"}`)},
		ManagedExternally:     newTrue(),
		Manager:               "7be0322d-354e-431e-84bd-ae4c9633beef",
	}

	// fake an IC update (create) message
	var update ICUpdateKubernetesJob
	update.Properties.Name = newIC.Name
	update.Properties.Category = newIC.Category
	update.Properties.Type = newIC.Type
	update.Properties.UUID = newIC.UUID
	update.Properties.Job.MetaData.JobName = "myJob"
	update.Properties.Job.Spec.Active = "70"
	update.Status.ManagedBy = newIC.Manager
	update.Status.State = newIC.State
	var container Container
	container.Name = "myContainer"
	container.Image = "python:latest"
	var Containers []Container
	update.Properties.Job.Spec.Template.Spec.Containers = append(Containers, container)

	payload, err := json.Marshal(update)
	if err != nil {
		return err
	}

	//session = helper.NewAMQPSession("villas-test-session", amqpURI, "villas", infrastructure_component.ProcessMessage)
	SetAMQPSession(session)

	//time.Sleep(3 * time.Second)

	err = session.CheckConnection()
	if err != nil {
		return err
	}

	err = session.Send(payload, newIC.Manager)
	time.Sleep(2 * time.Second)
	return err
}

func addComponentConfig(scenarioID uint, token string) error {
	type ConfigRequest struct {
		Name            string         `json:"name,omitempty"`
		ScenarioID      uint           `json:"scenarioID,omitempty"`
		ICID            uint           `json:"icID,omitempty"`
		StartParameters postgres.Jsonb `json:"startParameters,omitempty"`
		FileIDs         []int64        `json:"fileIDs,omitempty"`
	}

	var newConfig1 = ConfigRequest{
		Name:            "Example for Signal generator",
		ScenarioID:      scenarioID,
		ICID:            1,
		StartParameters: postgres.Jsonb{RawMessage: json.RawMessage(`{"parameter1" : "testValue1A", "parameter2" : "testValue2A", "parameter3" : 42}`)},
		FileIDs:         []int64{},
	}

	_, _, err := helper.TestEndpoint(router, token,
		"/api/v2/configs", "POST", helper.KeyModels{"config": newConfig1})

	return err
}

func addFakeIC(session *helper.AMQPsession, token string) error {

	db := database.GetDB()
	var originalIC database.InfrastructureComponent
	err := db.Find(&originalIC, 1).Error
	if err != nil {
		return err
	}

	//session = helper.NewAMQPSession("villas-test-session", amqpURI, "villas", infrastructure_component.ProcessMessage)
	log.Println(session)
	SetAMQPSession(session)

	//time.Sleep(3 * time.Second)

	err = session.CheckConnection()
	if err != nil {
		return err
	}

	var update ICUpdateKubernetesJob
	err = json.Unmarshal(originalIC.StatusUpdateRaw.RawMessage, &update)
	if err != nil {
		return err
	}
	update.Properties.UUID = "4854af30-325f-44a5-ad59-b67b2597de68"

	payload, err := json.Marshal(update)
	if err != nil {
		return err
	}

	log.Println(session)

	err = session.Send(payload, originalIC.Manager)
	time.Sleep(2 * time.Second)
	return err
}
