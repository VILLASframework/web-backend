package main

import (
	"fmt"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/amqp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/database"
	_ "git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/doc/api" // doc/api folder is used by Swag CLI, you have to import it
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/dashboard"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/file"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/scenario"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/signal"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/simulationmodel"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/simulator"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/user"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/widget"
)

// @title VILLASweb Backend API
// @version 2.0
// @description This is the API of the VILLASweb Backend
// @description WORK IN PROGRESS! PLEASE BE PATIENT!

// @description This documentation is auto-generated based on the API documentation in the code.
// @description The tool https://github.com/swaggo/swag is used to auto-generate API docs for gin.

// @contact.name Sonja Happ
// @contact.email sonja.happ@eonerc.rwth-aachen.de

// @license.name GNU GPL 3.0
// @license.url http://www.gnu.de/documents/gpl-3.0.en.html

// @host localhost:4000
// @BasePath /api/v2
func main() {
	// TODO DB_TEST is used for testing, should be DB_NAME in production
	db := database.InitDB(database.DB_TEST)
	database.MigrateModels(db)
	defer db.Close()

	// TODO the following line should be removed in production, it adds test data to the DB
	database.DBAddTestData(db)

	r := gin.Default()

	api := r.Group("/api/v2")

	// All endpoints require authentication except when someone wants to
	// login (POST /authenticate)
	user.RegisterAuthenticate(api.Group("/authenticate"))

	api.Use(user.Authentication(true))

	scenario.RegisterScenarioEndpoints(api.Group("/scenarios"))
	simulationmodel.RegisterSimulationModelEndpoints(api.Group("/models"))
	signal.RegisterSignalEndpoints(api.Group("/signals"))
	dashboard.RegisterDashboardEndpoints(api.Group("/dashboards"))
	widget.RegisterWidgetEndpoints(api.Group("/widgets"))
	file.RegisterFileEndpoints(api.Group("/files"))
	user.RegisterUserEndpoints(api.Group("/users"))
	simulator.RegisterSimulatorEndpoints(api.Group("/simulators"))

	r.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	if database.WITH_AMQP == true {
		fmt.Println("Starting AMQP client")
		err := amqp.ConnectAMQP("amqp://localhost")
		if err != nil {
			panic(err)
		}

		// Periodically call the Ping function to check which simulators are still there
		ticker := time.NewTicker(10 * time.Second)
		go func() {

			for {
				select {
				case <-ticker.C:
					err = amqp.PingAMQP()
					if err != nil {
						fmt.Println("AMQP Error: ", err.Error())
					}
				}
			}

		}()
	}
	// server at port 4000 to match frontend's redirect path
	r.Run(":4000")
}
