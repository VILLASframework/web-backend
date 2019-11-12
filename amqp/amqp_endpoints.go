package amqp

import (
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/helper"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/simulator"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func RegisterAMQPEndpoint(r *gin.RouterGroup) {
	r.POST("/:simulatorID/action", sendActionToSimulator)
}

// sendActionToSimulator godoc
// @Summary Send an action to simulator (only available if backend server is started with -amqp parameter)
// @ID sendActionToSimulator
// @Tags AMQP
// @Produce json
// @Param inputAction query string true "Action for simulator"
// @Success 200 {object} docs.ResponseError "Action sent successfully"
// @Failure 400 {object} docs.ResponseError "Bad request"
// @Failure 404 {object} docs.ResponseError "Not found"
// @Failure 422 {object} docs.ResponseError "Unprocessable entity"
// @Failure 500 {object} docs.ResponseError "Internal server error"
// @Param simulatorID path int true "Simulator ID"
// @Router /simulators/{simulatorID}/action [post]
func sendActionToSimulator(c *gin.Context) {

	ok, s := simulator.CheckPermissions(c, database.ModelSimulatorAction, database.Update, true)
	if !ok {
		return
	}

	var actions []Action
	err := c.BindJSON(&actions)
	if err != nil {
		helper.BadRequestError(c, "Error binding form data to JSON: "+err.Error())
		return
	}

	now := time.Now()

	for _, action := range actions {
		if action.When == 0 {
			action.When = float32(now.Unix())
		}

		err = SendActionAMQP(action, s.UUID)
		if err != nil {
			helper.InternalServerError(c, "Unable to send actions to simulator: "+err.Error())
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "OK.",
	})
}
