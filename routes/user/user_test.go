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

	// authenticate
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
}
