package healthz

import (
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/amqp"
	//"git.rwth-aachen.de/acs/public/villas/web-backend-go/amqp"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/helper"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/user"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"net/http"

	//"net/http"
	"testing"
)

var router *gin.Engine
var db *gorm.DB

func TestHealthz(t *testing.T) {
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

	// close db connection
	err = db.Close()
	assert.NoError(t, err)

	// test healthz endpoint for unconnected DB and AMQP client
	code, resp, err := helper.TestEndpoint(router, token, "api/healthz", http.MethodGet, nil)
	assert.NoError(t, err)
	assert.Equalf(t, 500, code, "Response body: \n%v\n", resp)

	// reconnect DB
	db = database.InitDB(database.DB_NAME, false)
	defer db.Close()

	// test healthz endpoint for connected DB and unconnected AMQP client
	code, resp, err = helper.TestEndpoint(router, token, "api/healthz", http.MethodGet, nil)
	assert.NoError(t, err)
	assert.Equalf(t, 500, code, "Response body: \n%v\n", resp)

	// connect AMQP client (make sure that AMQP_URL is set via command line parameter -amqp)
	err = amqp.ConnectAMQP(database.AMQP_URL)
	assert.NoError(t, err)

	// test healthz endpoint for connected DB and AMQP client
	code, resp, err = helper.TestEndpoint(router, token, "api/healthz", http.MethodGet, nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)
}
