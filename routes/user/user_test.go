package user

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
)

func TestUserEndpoints(t *testing.T) {

	db := common.DummyInitDB()
	defer db.Close()
	common.DummyOnlyAdminDB(db)

	router := gin.Default()
	api := router.Group("/api")

	VisitorAuthenticate(api.Group("/authenticate"))
	api.Use(Authentication(true))
	RegisterUserEndpoints(api.Group("/users"))

	// authenticate as admin
	token, err := common.NewAuthenticateForTest(router,
		"/api/authenticate", "POST", common.AdminCredentials, 200)
	assert.NoError(t, err)

	// test GET user/
	err = common.NewTestEndpoint(router, token,
		"/api/users", "GET", nil,
		200, common.KeyModels{"users": []common.User{common.User0}})
	assert.NoError(t, err)

	// test GET user/1 (the admin)
	err = common.NewTestEndpoint(router, token,
		"/api/users/1", "GET", nil,
		200, common.KeyModels{"user": common.User0})
	assert.NoError(t, err)

	// test POST user/ $newUser
	newUser := common.Request{
		Username: common.UserA.Username,
		Password: common.StrPasswordA,
		Mail:     common.UserA.Mail,
		Role:     common.UserA.Role,
	}
	// TODO: For now the response from this endpoint has the form
	// {"user":$Username}. Make sure that this should be the usual
	// {"user":$User{}} response.
	err = common.NewTestEndpoint(router, token,
		"/api/users", "POST", common.KeyModels{"user": newUser},
		200, common.KeyModels{"user": newUser.Username})
	assert.NoError(t, err)

	// test PUT user/1 $modifiedUser
	modifiedUser := common.Request{Role: "Admin"}
	err = common.NewTestEndpoint(router, token,
		"/api/users/1", "PUT", common.KeyModels{"user": modifiedUser},
		200, common.KeyModels{"user": modifiedUser.Username})
	assert.NoError(t, err)
}
