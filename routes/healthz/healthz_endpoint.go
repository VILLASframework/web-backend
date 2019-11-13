package healthz

import (
	"fmt"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/amqp"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
	"github.com/gin-gonic/gin"
	"net/http"
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
	if err != nil {
		return
	}

	// check if connection to AMQP broker is alive if backend was started with AMQP client
	if len(database.AMQP_URL) != 0 {
		err = amqp.CheckConnection()
		if err != nil {
			fmt.Println(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"success:": false,
				"message":  err.Error(),
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{})
}