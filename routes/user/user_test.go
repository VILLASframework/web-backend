package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/helper"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/database"
)

var router *gin.Engine
var db *gorm.DB

type UserRequest struct {
	Username    string `json:"username,omitempty"`
	Password    string `json:"password,omitempty"`
	OldPassword string `json:"oldPassword,omitempty"`
	Mail        string `json:"mail,omitempty"`
	Role        string `json:"role,omitempty"`
	Active      string `json:"active,omitempty"`
}

func TestMain(m *testing.M) {

	db = database.InitDB(database.DB_TEST)
	defer db.Close()

	router = gin.Default()
	api := router.Group("/api")

	RegisterAuthenticate(api.Group("/authenticate"))
	api.Use(Authentication(true))
	RegisterUserEndpoints(api.Group("/users"))

	os.Exit(m.Run())
}

func TestAuthenticate(t *testing.T) {
	database.DropTables(db)
	database.MigrateModels(db)
	assert.NoError(t, database.DBAddAdminAndUserAndGuest(db))

	// try to authenticate with non JSON body
	// should result in unauthorized
	w1 := httptest.NewRecorder()
	body, _ := json.Marshal("This is no JSON")
	req, err := http.NewRequest("POST", "/api/authenticate", bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w1, req)
	assert.Equalf(t, 401, w1.Code, "Response body: \n%v\n", w1.Body)

	malformedCredentials := helper.Credentials{
		Username: "TEST1",
	}
	// try to authenticate with malformed credentials
	// should result in unauthorized
	w2 := httptest.NewRecorder()
	body, _ = json.Marshal(malformedCredentials)
	req, err = http.NewRequest("POST", "/api/authenticate", bytes.NewBuffer(body))
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
	req, err = http.NewRequest("POST", "/api/authenticate", bytes.NewBuffer(body))
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
	req, err = http.NewRequest("POST", "/api/authenticate", bytes.NewBuffer(body))
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w4, req)
	assert.Equal(t, 401, w4.Code, w4.Body)

	// authenticate as admin
	_, err = helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.AdminCredentials)
	assert.NoError(t, err)

}

func TestAddGetUser(t *testing.T) {

	database.DropTables(db)
	database.MigrateModels(db)
	assert.NoError(t, database.DBAddAdminAndUserAndGuest(db))

	// authenticate as admin
	token, err := helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.AdminCredentials)
	assert.NoError(t, err)

	// try to POST with non JSON body
	// should result in bad request
	code, resp, err := helper.TestEndpoint(router, token,
		"/api/users", "POST", "This is not JSON")
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
		"/api/users", "POST", helper.KeyModels{"user": wrongUser})
	assert.NoError(t, err)
	assert.Equalf(t, 400, code, "Response body: \n%v\n", resp)

	// try POST with too short password
	// should result in bad request
	wrongUser.Username = "Longenoughusername"
	wrongUser.Password = "short"
	code, resp, err = helper.TestEndpoint(router, token,
		"/api/users", "POST", helper.KeyModels{"user": wrongUser})
	assert.NoError(t, err)
	assert.Equalf(t, 400, code, "Response body: \n%v\n", resp)

	// try POST with wrong role
	// should result in bad request
	wrongUser.Password = "longenough"
	wrongUser.Role = "ThisIsNotARole"
	code, resp, err = helper.TestEndpoint(router, token,
		"/api/users", "POST", helper.KeyModels{"user": wrongUser})
	assert.NoError(t, err)
	assert.Equalf(t, 400, code, "Response body: \n%v\n", resp)

	//try POST with wrong email
	// should result in bad request
	wrongUser.Mail = "noemailaddress"
	wrongUser.Role = "Guest"
	code, resp, err = helper.TestEndpoint(router, token,
		"/api/users", "POST", helper.KeyModels{"user": wrongUser})
	assert.NoError(t, err)
	assert.Equalf(t, 400, code, "Response body: \n%v\n", resp)

	// try POST with missing required fields
	// should result in bad request
	wrongUser.Mail = ""
	wrongUser.Role = ""
	wrongUser.Password = ""
	code, resp, err = helper.TestEndpoint(router, token,
		"/api/users", "POST", helper.KeyModels{"user": wrongUser})
	assert.NoError(t, err)
	assert.Equalf(t, 400, code, "Response body: \n%v\n", resp)

	// try POST with username that is already taken
	// should result in unprocessable entity
	wrongUser.Mail = "test@test.com"
	wrongUser.Role = "Guest"
	wrongUser.Password = "blablabla"
	wrongUser.Username = "User_A"
	code, resp, err = helper.TestEndpoint(router, token,
		"/api/users", "POST", helper.KeyModels{"user": wrongUser})
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
		"/api/users", "POST", helper.KeyModels{"user": newUser})
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
		fmt.Sprintf("/api/users/%v", newUserID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare GET's response with the newUser (Password omitted)
	err = helper.CompareResponse(resp, helper.KeyModels{"user": newUser})
	assert.NoError(t, err)

	// try to GET user with invalid user ID
	// should result in bad request
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/users/bla"), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 400, code, "Response body: \n%v\n", resp)
}

func TestUsersNotAllowedActions(t *testing.T) {

	database.DropTables(db)
	database.MigrateModels(db)
	assert.NoError(t, database.DBAddAdminAndUserAndGuest(db))

	// authenticate as admin
	token, err := helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.AdminCredentials)
	assert.NoError(t, err)

	// Add a user
	newUser := UserRequest{
		Username: "NotAllowedActions",
		Password: "N0t_@LLow3d_act10n5",
		Mail:     "not@allowed.ac",
		Role:     "User",
	}
	code, resp, err := helper.TestEndpoint(router, token,
		"/api/users", "POST", helper.KeyModels{"user": newUser})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	newUserID, err := helper.GetResponseID(resp)
	assert.NoError(t, err)

	// Authenticate as the added user
	token, err = helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", UserRequest{
			Username: newUser.Username,
			Password: newUser.Password,
		})
	assert.NoError(t, err)

	// Try to get all the users (NOT ALLOWED)
	code, resp, err = helper.TestEndpoint(router, token,
		"/api/users", "POST", helper.KeyModels{"user": newUser})
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

	// Try to read another user than self (eg. Admin) (NOT ALLOWED)
	code, resp, err = helper.TestEndpoint(router, token,
		"/api/users/0", "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 403, code, "Response body: \n%v\n", resp)

	// Try to delete another user (eg. Admin) (NOT ALLOWED)
	code, resp, err = helper.TestEndpoint(router, token,
		"/api/users/0", "DELETE", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

	// Try to delete self (NOT ALLOWED)
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/users/%v", newUserID), "DELETE", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

}

func TestGetAllUsers(t *testing.T) {

	database.DropTables(db)
	database.MigrateModels(db)
	assert.NoError(t, database.DBAddAdminAndUserAndGuest(db))

	// authenticate as admin
	token, err := helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.AdminCredentials)
	assert.NoError(t, err)

	// get the length of the GET all users response
	initialNumber, err := helper.LengthOfResponse(router, token,
		"/api/users", "GET", nil)
	assert.NoError(t, err)

	// Add a user
	newUser := UserRequest{
		Username: "UserGetAllUsers",
		Password: "get@ll_User5",
		Mail:     "get@all.users",
		Role:     "User",
	}
	code, resp, err := helper.TestEndpoint(router, token,
		"/api/users", "POST", helper.KeyModels{"user": newUser})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// get the length of the GET all users response again
	finalNumber, err := helper.LengthOfResponse(router, token,
		"/api/users", "GET", nil)
	assert.NoError(t, err)

	assert.Equal(t, finalNumber, initialNumber+1)

	// authenticate as normal user
	token, err = helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserACredentials)
	assert.NoError(t, err)

	// try to get all users as normal user
	// should result in unprocessable entity eror
	code, resp, err = helper.TestEndpoint(router, token,
		"/api/users", "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

}

func TestModifyAddedUserAsUser(t *testing.T) {

	database.DropTables(db)
	database.MigrateModels(db)
	assert.NoError(t, database.DBAddAdminAndUserAndGuest(db))

	// authenticate as admin
	token, err := helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.AdminCredentials)
	assert.NoError(t, err)

	// Add a user that will modify itself
	newUser := UserRequest{
		Username: "modMyself",
		Password: "mod_My5elf",
		Mail:     "mod@my.self",
		Role:     "User",
	}
	code, resp, err := helper.TestEndpoint(router, token,
		"/api/users", "POST", helper.KeyModels{"user": newUser})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	newUserID, err := helper.GetResponseID(resp)
	assert.NoError(t, err)

	// authenticate as the new user
	token, err = helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", UserRequest{
			Username: newUser.Username,
			Password: newUser.Password,
		})
	assert.NoError(t, err)

	// Try PUT with invalid user ID in path
	// Should return a bad request
	code, resp, err = helper.TestEndpoint(router, token, fmt.Sprintf("/api/users/blabla"), "PUT",
		helper.KeyModels{"user": newUser})
	assert.NoError(t, err)
	assert.Equalf(t, 400, code, "Response body: \n%v\n", resp)

	// Try to PUT a non JSON body
	// Should return a bad request
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/users/%v", newUserID), "PUT", "This is no JSON")
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
		fmt.Sprintf("/api/users/%v", newUserID), "PUT", helper.KeyModels{"user": modActiveState})
	assert.NoError(t, err)
	assert.Equalf(t, 403, code, "Response body: \n%v\n", resp)

	// modify newUser's own name
	modRequest := UserRequest{Username: "myNewName"}
	newUser.Username = modRequest.Username
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/users/%v", newUserID), "PUT",
		helper.KeyModels{"user": modRequest})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)
	err = helper.CompareResponse(resp, helper.KeyModels{"user": newUser})
	assert.NoError(t, err)

	// modify Admin's name (ILLEGAL)
	modRequest = UserRequest{Username: "myNewName"}
	newUser.Username = modRequest.Username
	code, resp, err = helper.TestEndpoint(router, token,
		"/api/users/1", "PUT", helper.KeyModels{"user": modRequest})
	assert.NoError(t, err)
	assert.Equalf(t, 403, code, "Response body: \n%v\n", resp)

	// modify newUser's own email
	modRequest = UserRequest{Mail: "my@new.email"}
	newUser.Mail = modRequest.Mail
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/users/%v", newUserID), "PUT",
		helper.KeyModels{"user": modRequest})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)
	err = helper.CompareResponse(resp, helper.KeyModels{"user": newUser})
	assert.NoError(t, err)

	// modify Admin's own email (ILLEGAL)
	modRequest = UserRequest{Mail: "my@new.email"}
	newUser.Mail = modRequest.Mail
	code, resp, err = helper.TestEndpoint(router, token,
		"/api/users/1", "PUT", helper.KeyModels{"user": modRequest})
	assert.NoError(t, err)
	assert.Equalf(t, 403, code, "Response body: \n%v\n", resp)

	// modify newUser's role (ILLEGAL)
	modRequest = UserRequest{Role: "Admin"}
	newUser.Role = modRequest.Role
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/users/%v", newUserID), "PUT",
		helper.KeyModels{"user": modRequest})
	assert.NoError(t, err)
	assert.Equalf(t, 403, code, "Response body: \n%v\n", resp)

	// modify newUser's password without providing old password
	modRequest = UserRequest{
		Password: "5tr0ng_pw!",
	}
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/users/%v", newUserID), "PUT",
		helper.KeyModels{"user": modRequest})
	assert.NoError(t, err)
	assert.Equalf(t, 400, code, "Response body: \n%v\n", resp)

	// modify newUser's password
	modRequest = UserRequest{
		Password:    "5tr0ng_pw!",
		OldPassword: "mod_My5elf",
	}
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/users/%v", newUserID), "PUT",
		helper.KeyModels{"user": modRequest})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// try to login as newUser with the modified username and password
	token, err = helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", UserRequest{
			Username: newUser.Username,
			Password: modRequest.Password,
		})
	assert.NoError(t, err)

	// modify Admin's password (ILLEGAL)
	modRequest = UserRequest{Password: "4dm1ns_pw!"}
	code, resp, err = helper.TestEndpoint(router, token,
		"/api/users/1", "PUT", helper.KeyModels{"user": modRequest})
	assert.NoError(t, err)
	assert.Equalf(t, 403, code, "Response body: \n%v\n", resp)
}

func TestInvalidUserUpdate(t *testing.T) {

	database.DropTables(db)
	database.MigrateModels(db)
	assert.NoError(t, database.DBAddAdminAndUserAndGuest(db))

	// authenticate as admin
	token, err := helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.AdminCredentials)
	assert.NoError(t, err)

	// Add a user
	newUser := UserRequest{
		Username: "invalidUpdatedUser",
		Password: "wr0ng_Upd@te!",
		Mail:     "inv@user.upd",
		Role:     "User",
	}
	code, resp, err := helper.TestEndpoint(router, token,
		"/api/users", "POST", helper.KeyModels{"user": newUser})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	newUserID, err := helper.GetResponseID(resp)
	assert.NoError(t, err)

	// try PUT with userID that does not exist
	// should result in not found
	modRequest := UserRequest{
		Password:    "longenough",
		OldPassword: "wr0ng_Upd@te!",
	}
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/users/%v", newUserID+1), "PUT",
		helper.KeyModels{"user": modRequest})
	assert.NoError(t, err)
	assert.Equalf(t, 404, code, "Response body: \n%v\n", resp)

	// try to PUT with username that is already taken
	// should result in bad request
	modRequest.Password = ""
	modRequest.Username = "User_A"
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/users/%v", newUserID), "PUT",
		helper.KeyModels{"user": modRequest})
	assert.NoError(t, err)
	assert.Equalf(t, 400, code, "Response body: \n%v\n", resp)

	// modify newUser's password with INVALID password
	modRequest.Password = "short"
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/users/%v", newUserID), "PUT",
		helper.KeyModels{"user": modRequest})
	assert.NoError(t, err)
	assert.Equalf(t, 400, code, "Response body: \n%v\n", resp)

	// modify newUser's email with INVALID email
	modRequest = UserRequest{Mail: "notEmail"}
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/users/%v", newUserID), "PUT",
		helper.KeyModels{"user": modRequest})
	assert.NoError(t, err)
	assert.Equalf(t, 400, code, "Response body: \n%v\n", resp)

	// modify newUser's role with INVALID role
	modRequest = UserRequest{Role: "noRole"}
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/users/%v", newUserID), "PUT",
		helper.KeyModels{"user": modRequest})
	assert.NoError(t, err)
	assert.Equalf(t, 400, code, "Response body: \n%v\n", resp)

}

func TestModifyAddedUserAsAdmin(t *testing.T) {

	database.DropTables(db)
	database.MigrateModels(db)
	assert.NoError(t, database.DBAddAdminAndUserAndGuest(db))

	// authenticate as admin
	token, err := helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.AdminCredentials)
	assert.NoError(t, err)

	// Add a user
	newUser := UserRequest{
		Username: "modAddedUser",
		Password: "mod_4d^2ed_0ser",
		Mail:     "mod@added.user",
		Role:     "User",
	}
	code, resp, err := helper.TestEndpoint(router, token,
		"/api/users", "POST", helper.KeyModels{"user": newUser})
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
		fmt.Sprintf("/api/users/%v", newUserID), "PUT",
		helper.KeyModels{"user": modRequest})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)
	err = helper.CompareResponse(resp, helper.KeyModels{"user": newUser})
	assert.NoError(t, err)

	// modify newUser's email
	modRequest = UserRequest{Mail: "new@e.mail"}
	newUser.Mail = modRequest.Mail
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/users/%v", newUserID), "PUT",
		helper.KeyModels{"user": modRequest})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)
	err = helper.CompareResponse(resp, helper.KeyModels{"user": newUser})
	assert.NoError(t, err)

	// modify newUser's role
	modRequest = UserRequest{Role: "Admin"}
	newUser.Role = modRequest.Role
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/users/%v", newUserID), "PUT",
		helper.KeyModels{"user": modRequest})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)
	err = helper.CompareResponse(resp, helper.KeyModels{"user": newUser})
	assert.NoError(t, err)

	// modify newUser's password
	modRequest = UserRequest{
		Password:    "4_g00d_pw!",
		OldPassword: "mod_4d^2ed_0ser",
	}
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/users/%v", newUserID), "PUT",
		helper.KeyModels{"user": modRequest})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// try to login as newUser with the modified username and password
	_, err = helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", UserRequest{
			Username: newUser.Username,
			Password: modRequest.Password,
		})
	assert.NoError(t, err)

	// authenticate as admin
	token, err = helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.AdminCredentials)
	assert.NoError(t, err)

	// modify newUser's Active status
	modRequest = UserRequest{Active: "no"}
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/users/%v", newUserID), "PUT",
		helper.KeyModels{"user": modRequest})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// try to login as newUser with the modified active status
	// should NOT work anymore!
	code, resp, err = helper.TestEndpoint(router, token,
		"/api/authenticate", "POST",
		UserRequest{
			Username: newUser.Username,
			Password: "4_g00d_pw!",
		})
	assert.NoError(t, err)
	assert.Equalf(t, 401, code, "Response body: \n%v\n", resp)
}

func TestDeleteUser(t *testing.T) {

	database.DropTables(db)
	database.MigrateModels(db)
	assert.NoError(t, database.DBAddAdminAndUserAndGuest(db))

	// authenticate as admin
	token, err := helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.AdminCredentials)
	assert.NoError(t, err)

	// Add a user
	newUser := UserRequest{
		Username: "toBeDeletedUser",
		Password: "f0r_deletIOn_0ser",
		Mail:     "to@be.deleted",
		Role:     "User",
	}
	code, resp, err := helper.TestEndpoint(router, token,
		"/api/users", "POST", helper.KeyModels{"user": newUser})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	newUserID, err := helper.GetResponseID(resp)
	assert.NoError(t, err)

	// try to DELETE with invalid ID
	// should result in bad request
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/users/bla"), "DELETE", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 400, code, "Response body: \n%v\n", resp)

	// try to DELETE with ID that does not exist
	// should result in not found
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/users/%v", newUserID+1), "DELETE", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 404, code, "Response body: \n%v\n", resp)

	// Count the number of all the users returned
	initialNumber, err := helper.LengthOfResponse(router, token,
		"/api/users", "GET", nil)
	assert.NoError(t, err)

	// Delete the added newUser
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/users/%v", newUserID), "DELETE", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Again count the number of all the users returned
	finalNumber, err := helper.LengthOfResponse(router, token,
		"/api/users", "GET", nil)
	assert.NoError(t, err)

	assert.Equal(t, finalNumber, initialNumber-1)
}
