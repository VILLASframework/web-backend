package user

import (
	"encoding/json"
	"testing"

	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
)

func TestUserEndpoints(t *testing.T) {

	myUsers := []common.User{common.User0, common.UserA, common.UserB}
	msgUsers := common.ResponseMsgUsers{Users: myUsers}

	db := common.DummyInitDB()
	defer db.Close()
	common.DummyPopulateDB(db)

	router := gin.Default()
	api := router.Group("/api")

	VisitorAuthenticate(api.Group("/authenticate"))
	api.Use(Authentication(true))
	RegisterUserEndpoints(api.Group("/users"))

	credjson, _ := json.Marshal(common.CredAdmin)
	msgUsersjson, _ := json.Marshal(msgUsers)

	token := common.AuthenticateForTest(t, router, "/api/authenticate",
		"POST", credjson, 200)

	// test GET user/
	common.TestEndpoint(t, router, token, "/api/users", "GET", nil, 200,
		msgUsersjson)
}
