package healthz

import (
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/amqp"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/configuration"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/helper"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
	"net/http"

	"testing"
)

var router *gin.Engine
var db *gorm.DB

func TestHealthz(t *testing.T) {
	err := configuration.InitConfig()
	assert.NoError(t, err)

	// connect DB
	db, err = database.InitDB(configuration.GolbalConfig)
	assert.NoError(t, err)
	defer db.Close()

	router = gin.Default()

	RegisterHealthzEndpoint(router.Group("/healthz"))

	// close db connection
	err = db.Close()
	assert.NoError(t, err)

	// test healthz endpoint for unconnected DB and AMQP client
	code, resp, err := helper.TestEndpoint(router, "", "healthz", http.MethodGet, nil)
	assert.NoError(t, err)
	assert.Equalf(t, 500, code, "Response body: \n%v\n", resp)

	// reconnect DB
	db, err = database.InitDB(configuration.GolbalConfig)
	assert.NoError(t, err)
	defer db.Close()

	// test healthz endpoint for connected DB and unconnected AMQP client
	code, resp, err = helper.TestEndpoint(router, "", "healthz", http.MethodGet, nil)
	assert.NoError(t, err)
	assert.Equalf(t, 500, code, "Response body: \n%v\n", resp)

	// connect AMQP client (make sure that AMQP_URL is set via command line parameter -amqp)
	url, err := configuration.GolbalConfig.String("amqp.url")
	assert.NoError(t, err)
	err = amqp.ConnectAMQP(url)
	assert.NoError(t, err)

	// test healthz endpoint for connected DB and AMQP client
	code, resp, err = helper.TestEndpoint(router, "", "healthz", http.MethodGet, nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)
}
