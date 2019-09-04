package user

import (
	"fmt"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
)

var router *gin.Engine
var db *gorm.DB

func TestMain(m *testing.M) {

	db = common.DummyInitDB()
	defer db.Close()

	router = gin.Default()
	api := router.Group("/api")

	RegisterAuthenticate(api.Group("/authenticate"))
	api.Use(Authentication(true))
	RegisterUserEndpoints(api.Group("/users"))

	os.Exit(m.Run())
}

func TestAddGetUser(t *testing.T) {

	common.DropTables(db)
	common.MigrateModels(db)
	common.DummyAddOnlyUserTableWithAdminDB(db)

	// authenticate as admin
	token, err := common.NewAuthenticateForTest(router,
		"/api/authenticate", "POST", common.AdminCredentials)
	assert.NoError(t, err)

	// test POST user/ $newUser
	newUser := common.Request{
		Username: "Alice483",
		Password: "th1s_I5_@lice#",
		Mail:     "mail@domain.com",
		Role:     "User",
	}
	code, resp, err := common.NewTestEndpoint(router, token,
		"/api/users", "POST", common.KeyModels{"user": newUser})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Turn password member of newUser to empty string so it is omitted
	// in marshaling. The password will never be included in the
	// response and if is non empty in request we will not be able to do
	// request-response comparison. Reserve it also for later usage.
	newUser.Password = ""

	// Compare POST's response with the newUser (Password omitted)
	err = common.CompareResponse(resp, common.KeyModels{"user": newUser})
	assert.NoError(t, err)

	// Read newUser's ID from the response
	newUserID, err := common.GetResponseID(resp)
	assert.NoError(t, err)

	// Get the newUser
	code, resp, err = common.NewTestEndpoint(router, token,
		fmt.Sprintf("/api/users/%v", newUserID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare GET's response with the newUser (Password omitted)
	err = common.CompareResponse(resp, common.KeyModels{"user": newUser})
	assert.NoError(t, err)
}

func TestUsersNotAllowedActions(t *testing.T) {

	common.DropTables(db)
	common.MigrateModels(db)
	common.DummyAddOnlyUserTableWithAdminDB(db)

	// authenticate as admin
	token, err := common.NewAuthenticateForTest(router,
		"/api/authenticate", "POST", common.AdminCredentials)
	assert.NoError(t, err)

	// Add a user
	newUser := common.Request{
		Username: "NotAllowedActions",
		Password: "N0t_@LLow3d_act10n5",
		Mail:     "not@allowed.ac",
		Role:     "User",
	}
	code, resp, err := common.NewTestEndpoint(router, token,
		"/api/users", "POST", common.KeyModels{"user": newUser})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	newUserID, err := common.GetResponseID(resp)
	assert.NoError(t, err)

	// Authenticate as the added user
	token, err = common.NewAuthenticateForTest(router,
		"/api/authenticate", "POST", common.Request{
			Username: newUser.Username,
			Password: newUser.Password,
		})
	assert.NoError(t, err)

	// Try to get all the users (NOT ALLOWED)
	code, resp, err = common.NewTestEndpoint(router, token,
		"/api/users", "POST", common.KeyModels{"user": newUser})
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

	// Try to read another user than self (eg. Admin) (NOT ALLOWED)
	code, resp, err = common.NewTestEndpoint(router, token,
		"/api/users/0", "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 403, code, "Response body: \n%v\n", resp)

	// Try to delete another user (eg. Admin) (NOT ALLOWED)
	code, resp, err = common.NewTestEndpoint(router, token,
		"/api/users/0", "DELETE", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

	// Try to delete self (NOT ALLOWED)
	code, resp, err = common.NewTestEndpoint(router, token,
		fmt.Sprintf("/api/users/%v", newUserID), "DELETE", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

}

func TestGetAllUsers(t *testing.T) {

	common.DropTables(db)
	common.MigrateModels(db)
	common.DummyAddOnlyUserTableWithAdminDB(db)

	// authenticate as admin
	token, err := common.NewAuthenticateForTest(router,
		"/api/authenticate", "POST", common.AdminCredentials)
	assert.NoError(t, err)

	// get the length of the GET all users response
	initialNumber, err := common.LengthOfResponse(router, token,
		"/api/users", "GET", nil)
	assert.NoError(t, err)

	// Add a user
	newUser := common.Request{
		Username: "UserGetAllUsers",
		Password: "get@ll_User5",
		Mail:     "get@all.users",
		Role:     "User",
	}
	code, resp, err := common.NewTestEndpoint(router, token,
		"/api/users", "POST", common.KeyModels{"user": newUser})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// get the length of the GET all users response again
	finalNumber, err := common.LengthOfResponse(router, token,
		"/api/users", "GET", nil)
	assert.NoError(t, err)

	assert.Equal(t, finalNumber, initialNumber+1)
}

func TestModifyAddedUserAsUser(t *testing.T) {

	common.DropTables(db)
	common.MigrateModels(db)
	common.DummyAddOnlyUserTableWithAdminDB(db)

	// authenticate as admin
	token, err := common.NewAuthenticateForTest(router,
		"/api/authenticate", "POST", common.AdminCredentials)
	assert.NoError(t, err)

	// Add a user that will modify itself
	newUser := common.Request{
		Username: "modMyself",
		Password: "mod_My5elf",
		Mail:     "mod@my.self",
		Role:     "User",
	}
	code, resp, err := common.NewTestEndpoint(router, token,
		"/api/users", "POST", common.KeyModels{"user": newUser})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	newUserID, err := common.GetResponseID(resp)
	assert.NoError(t, err)

	// authenticate as the new user
	token, err = common.NewAuthenticateForTest(router,
		"/api/authenticate", "POST", common.Request{
			Username: newUser.Username,
			Password: newUser.Password,
		})
	assert.NoError(t, err)

	// Turn password member of newUser to empty string so it is omitted
	// in marshaling. The password will never be included in the
	// response and if is non empty in request we will not be able to do
	// request-response comparison
	newUser.Password = ""

	// modify newUser's own name
	modRequest := common.Request{Username: "myNewName"}
	newUser.Username = modRequest.Username
	code, resp, err = common.NewTestEndpoint(router, token,
		fmt.Sprintf("/api/users/%v", newUserID), "PUT",
		common.KeyModels{"user": modRequest})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)
	err = common.CompareResponse(resp, common.KeyModels{"user": newUser})
	assert.NoError(t, err)

	// modify Admin's name (ILLEGAL)
	modRequest = common.Request{Username: "myNewName"}
	newUser.Username = modRequest.Username
	code, resp, err = common.NewTestEndpoint(router, token,
		"/api/users/1", "PUT", common.KeyModels{"user": modRequest})
	assert.NoError(t, err)
	assert.Equalf(t, 403, code, "Response body: \n%v\n", resp)

	// modify newUser's own email
	modRequest = common.Request{Mail: "my@new.email"}
	newUser.Mail = modRequest.Mail
	code, resp, err = common.NewTestEndpoint(router, token,
		fmt.Sprintf("/api/users/%v", newUserID), "PUT",
		common.KeyModels{"user": modRequest})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)
	err = common.CompareResponse(resp, common.KeyModels{"user": newUser})
	assert.NoError(t, err)

	// modify Admin's own email (ILLEGAL)
	modRequest = common.Request{Mail: "my@new.email"}
	newUser.Mail = modRequest.Mail
	code, resp, err = common.NewTestEndpoint(router, token,
		"/api/users/1", "PUT", common.KeyModels{"user": modRequest})
	assert.NoError(t, err)
	assert.Equalf(t, 403, code, "Response body: \n%v\n", resp)

	// modify newUser's role (ILLEGAL)
	modRequest = common.Request{Role: "Admin"}
	newUser.Role = modRequest.Role
	code, resp, err = common.NewTestEndpoint(router, token,
		fmt.Sprintf("/api/users/%v", newUserID), "PUT",
		common.KeyModels{"user": modRequest})
	assert.NoError(t, err)
	assert.Equalf(t, 403, code, "Response body: \n%v\n", resp)

	// modify newUser's password
	modRequest = common.Request{Password: "5tr0ng_pw!"}
	code, resp, err = common.NewTestEndpoint(router, token,
		fmt.Sprintf("/api/users/%v", newUserID), "PUT",
		common.KeyModels{"user": modRequest})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// try to login as newUser with the modified username and password
	token, err = common.NewAuthenticateForTest(router,
		"/api/authenticate", "POST", common.Request{
			Username: newUser.Username,
			Password: modRequest.Password,
		})
	assert.NoError(t, err)

	// modify Admin's password (ILLEGAL)
	modRequest = common.Request{Password: "4dm1ns_pw!"}
	code, resp, err = common.NewTestEndpoint(router, token,
		"/api/users/1", "PUT", common.KeyModels{"user": modRequest})
	assert.NoError(t, err)
	assert.Equalf(t, 403, code, "Response body: \n%v\n", resp)
}

func TestInvalidUserUpdate(t *testing.T) {

	common.DropTables(db)
	common.MigrateModels(db)
	common.DummyAddOnlyUserTableWithAdminDB(db)

	// authenticate as admin
	token, err := common.NewAuthenticateForTest(router,
		"/api/authenticate", "POST", common.AdminCredentials)
	assert.NoError(t, err)

	// Add a user
	newUser := common.Request{
		Username: "invalidUpdatedUser",
		Password: "wr0ng_Upd@te!",
		Mail:     "inv@user.upd",
		Role:     "User",
	}
	code, resp, err := common.NewTestEndpoint(router, token,
		"/api/users", "POST", common.KeyModels{"user": newUser})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	newUserID, err := common.GetResponseID(resp)
	assert.NoError(t, err)

	// modify newUser's password with INVALID password
	modRequest := common.Request{Password: "short"}
	code, resp, err = common.NewTestEndpoint(router, token,
		fmt.Sprintf("/api/users/%v", newUserID), "PUT",
		common.KeyModels{"user": modRequest})
	assert.NoError(t, err)
	assert.Equalf(t, 400, code, "Response body: \n%v\n", resp)

	// modify newUser's email with INVALID email
	modRequest = common.Request{Mail: "notEmail"}
	code, resp, err = common.NewTestEndpoint(router, token,
		fmt.Sprintf("/api/users/%v", newUserID), "PUT",
		common.KeyModels{"user": modRequest})
	assert.NoError(t, err)
	assert.Equalf(t, 400, code, "Response body: \n%v\n", resp)

	// modify newUser's role with INVALID role
	modRequest = common.Request{Role: "noRole"}
	code, resp, err = common.NewTestEndpoint(router, token,
		fmt.Sprintf("/api/users/%v", newUserID), "PUT",
		common.KeyModels{"user": modRequest})
	assert.NoError(t, err)
	assert.Equalf(t, 400, code, "Response body: \n%v\n", resp)
}

func TestModifyAddedUserAsAdmin(t *testing.T) {

	common.DropTables(db)
	common.MigrateModels(db)
	common.DummyAddOnlyUserTableWithAdminDB(db)

	// authenticate as admin
	token, err := common.NewAuthenticateForTest(router,
		"/api/authenticate", "POST", common.AdminCredentials)
	assert.NoError(t, err)

	// Add a user
	newUser := common.Request{
		Username: "modAddedUser",
		Password: "mod_4d^2ed_0ser",
		Mail:     "mod@added.user",
		Role:     "User",
	}
	code, resp, err := common.NewTestEndpoint(router, token,
		"/api/users", "POST", common.KeyModels{"user": newUser})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	newUserID, err := common.GetResponseID(resp)
	assert.NoError(t, err)

	// Turn password member of newUser to empty string so it is omitted
	// in marshaling. The password will never be included in the
	// response and if is non empty in request we will not be able to do
	// request-response comparison
	newUser.Password = ""

	// modify newUser's name
	modRequest := common.Request{Username: "NewUsername"}
	newUser.Username = modRequest.Username
	code, resp, err = common.NewTestEndpoint(router, token,
		fmt.Sprintf("/api/users/%v", newUserID), "PUT",
		common.KeyModels{"user": modRequest})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)
	err = common.CompareResponse(resp, common.KeyModels{"user": newUser})
	assert.NoError(t, err)

	// modify newUser's email
	modRequest = common.Request{Mail: "new@e.mail"}
	newUser.Mail = modRequest.Mail
	code, resp, err = common.NewTestEndpoint(router, token,
		fmt.Sprintf("/api/users/%v", newUserID), "PUT",
		common.KeyModels{"user": modRequest})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)
	err = common.CompareResponse(resp, common.KeyModels{"user": newUser})
	assert.NoError(t, err)

	// modify newUser's role
	modRequest = common.Request{Role: "Admin"}
	newUser.Role = modRequest.Role
	code, resp, err = common.NewTestEndpoint(router, token,
		fmt.Sprintf("/api/users/%v", newUserID), "PUT",
		common.KeyModels{"user": modRequest})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)
	err = common.CompareResponse(resp, common.KeyModels{"user": newUser})
	assert.NoError(t, err)

	// modify newUser's password
	modRequest = common.Request{Password: "4_g00d_pw!"}
	code, resp, err = common.NewTestEndpoint(router, token,
		fmt.Sprintf("/api/users/%v", newUserID), "PUT",
		common.KeyModels{"user": modRequest})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// try to login as newUser with the modified username and password
	_, err = common.NewAuthenticateForTest(router,
		"/api/authenticate", "POST", common.Request{
			Username: newUser.Username,
			Password: modRequest.Password,
		})
	assert.NoError(t, err)
}

func TestDeleteUser(t *testing.T) {

	common.DropTables(db)
	common.MigrateModels(db)
	common.DummyAddOnlyUserTableWithAdminDB(db)

	// authenticate as admin
	token, err := common.NewAuthenticateForTest(router,
		"/api/authenticate", "POST", common.AdminCredentials)
	assert.NoError(t, err)

	// Add a user
	newUser := common.Request{
		Username: "toBeDeletedUser",
		Password: "f0r_deletIOn_0ser",
		Mail:     "to@be.deleted",
		Role:     "User",
	}
	code, resp, err := common.NewTestEndpoint(router, token,
		"/api/users", "POST", common.KeyModels{"user": newUser})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	newUserID, err := common.GetResponseID(resp)
	assert.NoError(t, err)

	// Count the number of all the users returned
	initialNumber, err := common.LengthOfResponse(router, token,
		"/api/users", "GET", nil)
	assert.NoError(t, err)

	// Delete the added newUser
	code, resp, err = common.NewTestEndpoint(router, token,
		fmt.Sprintf("/api/users/%v", newUserID), "DELETE", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Again count the number of all the users returned
	finalNumber, err := common.LengthOfResponse(router, token,
		"/api/users", "GET", nil)
	assert.NoError(t, err)

	assert.Equal(t, finalNumber, initialNumber-1)
}
