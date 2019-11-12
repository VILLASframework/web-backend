package healthz

import (
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/amqp"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/helper"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/user"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

var router *gin.Engine
var db *gorm.DB

const test_amqp_url = "amqp://villas:villas@rabbitmq:5672"

//const test_amqp_url = "amqp://rabbit@goofy:5672"

func TestHealthz(t *testing.T) {
	database.AMQP_URL = test_amqp_url
	// connect DB
	db = database.InitDB(database.DB_NAME, true)
	defer db.Close()

	assert.NoError(t, database.DBAddAdminAndUserAndGuest(db))

	router = gin.Default()
	api := router.Group("/api")

	user.RegisterAuthenticate(api.Group("/authenticate"))
	api.Use(user.Authentication(true))
	RegisterHealthzEndpoint(api.Group("/healthz"))

	// authenticate as normal user
	token, err := helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserACredentials)
	assert.NoError(t, err)

	// connect AMQP client

	err = amqp.ConnectAMQP(database.AMQP_URL)
	assert.NoError(t, err)

	// test healthz endpoint for connected DB and AMQP client
	code, resp, err := helper.TestEndpoint(router, token, "api/healthz", http.MethodGet, nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

}
