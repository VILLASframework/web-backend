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

	common.DummyAddOnlyUserTableWithAdminDB(db)

	router = gin.Default()
	api := router.Group("/api")

	RegisterAuthenticate(api.Group("/authenticate"))
	api.Use(Authentication(true))
	RegisterUserEndpoints(api.Group("/users"))

	os.Exit(m.Run())
}

func TestAddUser(t *testing.T) {

	// authenticate as admin
	token, err := common.NewAuthenticateForTest(router,
		"/api/authenticate", "POST", common.AdminCredentials, 200)
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
	err = common.CompareResponse(resp,
		common.KeyModels{"user": common.Request{
			Username: newUser.Username,
			Mail:     newUser.Mail,
			Role:     newUser.Role,
		}})
	assert.NoError(t, err)
}

func TestGetAllUsers(t *testing.T) {

	// authenticate as admin
	token, err := common.NewAuthenticateForTest(router,
		"/api/authenticate", "POST", common.AdminCredentials, 200)
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

func TestModifyAddedUserAsAdmin(t *testing.T) {

	// authenticate as admin
	token, err := common.NewAuthenticateForTest(router,
		"/api/authenticate", "POST", common.AdminCredentials, 200)
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
	modRequest1 := common.Request{Username: "NewUsername"}
	newUser.Username = modRequest1.Username
	code, resp, err = common.NewTestEndpoint(router, token,
		fmt.Sprintf("/api/users/%v", newUserID), "PUT",
		common.KeyModels{"user": modRequest1})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)
	err = common.CompareResponse(resp, common.KeyModels{"user": newUser})
	assert.NoError(t, err)

	// modify newUser's email
	modRequest2 := common.Request{Mail: "new@e.mail"}
	newUser.Mail = modRequest2.Mail
	code, resp, err = common.NewTestEndpoint(router, token,
		fmt.Sprintf("/api/users/%v", newUserID), "PUT",
		common.KeyModels{"user": modRequest2})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)
	err = common.CompareResponse(resp, common.KeyModels{"user": newUser})
	assert.NoError(t, err)

	// modify newUser's role
	modRequest3 := common.Request{Role: "Admin"}
	newUser.Role = modRequest3.Role
	code, resp, err = common.NewTestEndpoint(router, token,
		fmt.Sprintf("/api/users/%v", newUserID), "PUT",
		common.KeyModels{"user": modRequest3})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)
	err = common.CompareResponse(resp, common.KeyModels{"user": newUser})
	assert.NoError(t, err)

	// modify newUser's password with INVALID password
	modRequest4 := common.Request{Password: "short"}
	code, resp, err = common.NewTestEndpoint(router, token,
		fmt.Sprintf("/api/users/%v", newUserID), "PUT",
		common.KeyModels{"user": modRequest4})
	assert.NoError(t, err)
	assert.Equalf(t, 400, code, "Response body: \n%v\n", resp) // HTTP 400

	// modify newUser's password with VALID password
	modRequest5 := common.Request{Password: "4_g00d_pw!"}
	code, resp, err = common.NewTestEndpoint(router, token,
		fmt.Sprintf("/api/users/%v", newUserID), "PUT",
		common.KeyModels{"user": modRequest5})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// try to login as newUser with the modified username and password
	_, err = common.NewAuthenticateForTest(router,
		"/api/authenticate", "POST", common.Request{
			Username: modRequest1.Username,
			Password: modRequest5.Password,
		}, 200)
	assert.NoError(t, err)
}
