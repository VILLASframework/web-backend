package user

import (
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
)

func TestUserEndpoints(t *testing.T) {

	myUsers := []common.User{common.User0}
	msgUsers := common.ResponseMsgUsers{Users: myUsers}

	db := common.DummyInitDB()
	defer db.Close()
	common.DummyOnlyAdminDB(db)

	router := gin.Default()
	api := router.Group("/api")

	VisitorAuthenticate(api.Group("/authenticate"))
	api.Use(Authentication(true))
	RegisterUserEndpoints(api.Group("/users"))

	token, err := common.NewAuthenticateForTest(router,
		"/api/authenticate", "POST", common.CredAdmin, 200)
	assert.NoError(t, err)

	// test GET user/
	err = common.NewTestEndpoint(router, token, "/api/users", "GET",
		nil, 200, msgUsers)
	assert.NoError(t, err)
}
