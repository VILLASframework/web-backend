package healthz

import (
	"net/http"
	"os"

	"git.rwth-aachen.de/acs/public/villas/web-backend-go/amqp"
	c "git.rwth-aachen.de/acs/public/villas/web-backend-go/config"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/helper"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"

	"testing"
)

var router *gin.Engine
var db *gorm.DB

func TestMain(m *testing.M) {
	c.InitConfig()

	os.Exit(m.Run())
}

func TestHealthz(t *testing.T) {
	// connect DB
	db = database.InitDB(c.Config)

	router = gin.Default()

	RegisterHealthzEndpoint(router.Group("/healthz"))

	// close db connection
	err := db.Close()
	assert.NoError(t, err)

	// test healthz endpoint for unconnected DB and AMQP client
	code, resp, err := helper.TestEndpoint(router, "", "healthz", http.MethodGet, nil)
	assert.NoError(t, err)
	assert.Equalf(t, 500, code, "Response body: \n%v\n", resp)

	// reconnect DB
	db = database.InitDB(c.Config)
	defer db.Close()

	// test healthz endpoint for connected DB and unconnected AMQP client
	code, resp, err = helper.TestEndpoint(router, "", "healthz", http.MethodGet, nil)
	assert.NoError(t, err)
	assert.Equalf(t, 500, code, "Response body: \n%v\n", resp)

	// connect AMQP client (make sure that AMQP_URL is set via command line parameter -amqp)
	url, _ := c.Config.String("amqp.url")
	err = amqp.ConnectAMQP(url)
	assert.NoError(t, err)

	// test healthz endpoint for connected DB and AMQP client
	code, resp, err = helper.TestEndpoint(router, "", "healthz", http.MethodGet, nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)
}
