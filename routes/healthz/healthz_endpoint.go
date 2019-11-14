package healthz

import (
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/amqp"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/configuration"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/helper"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

func RegisterHealthzEndpoint(r *gin.RouterGroup) {

	r.GET("", getHealth)
}

// getHealth godoc
// @Summary Get health status of backend
// @ID getHealth
// @Produce  json
// @Tags healthz
// @Success 200 "Backend is healthy, database and AMQP broker connections are alive"
// @Failure 500 {object} docs.ResponseError "Backend is NOT healthy"
// @Param Authorization header string true "Authorization token"
// @Router /healthz [get]
func getHealth(c *gin.Context) {

	// check if DB connection is active
	db := database.GetDB()
	err := db.DB().Ping()
	if helper.DBError(c, err) {
		return
	}

	// check if connection to AMQP broker is alive if backend was started with AMQP client
	url, err := configuration.GolbalConfig.String("amqp.url")
	if err != nil && strings.Contains(err.Error(), "Required setting 'amqp.url' not set") {
		c.JSON(http.StatusOK, gin.H{})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success:": false,
			"message":  err.Error(),
		})
		return
	}

	if len(url) != 0 {
		err = amqp.CheckConnection()
		if err != nil {
			log.Println(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"success:": false,
				"message":  err.Error(),
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{})
}
