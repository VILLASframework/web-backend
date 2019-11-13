package main

import (
	"log"
	"time"

	"git.rwth-aachen.de/acs/public/villas/web-backend-go/amqp"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/healthz"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	c "git.rwth-aachen.de/acs/public/villas/web-backend-go/config"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
	docs "git.rwth-aachen.de/acs/public/villas/web-backend-go/doc/api" // doc/api folder is used by Swag CLI, you have to import it
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/dashboard"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/file"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/scenario"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/signal"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/simulationmodel"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/simulator"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/user"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/widget"
)

// @title VILLASweb Backend API
// @version 2.0
// @description This is the VILLASweb Backend API v2.0.
// @description Parts of this API are still in development. Please check the [VILLASweb-backend-go repository](https://git.rwth-aachen.de/acs/public/villas/web-backend-go) for more information.
// @description This documentation is auto-generated based on the API documentation in the code. The tool [swag](https://github.com/swaggo/swag) is used to auto-generate API docs for the [gin-gonic](https://github.com/gin-gonic/gin) framework.
// @contact.name Sonja Happ
// @contact.email sonja.happ@eonerc.rwth-aachen.de
// @license.name GNU GPL 3.0
// @license.url http://www.gnu.de/documents/gpl-3.0.en.html
// @BasePath /api/v2
func main() {
	log.Println("Starting VILLASweb-backend-go")

	c := c.InitConfig()
	db := database.InitDB(c)
	defer db.Close()

	if m, _ := c.String("mode"); m == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	baseHost, _ := c.String("base.host")
	basePath, _ := c.String("base.path")
	docs.SwaggerInfo.Host = baseHost
	docs.SwaggerInfo.BasePath = basePath

	r := gin.Default()

	api := r.Group(basePath)

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
	healthz.RegisterHealthzEndpoint(api.Group("/healthz"))

	r.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	amqpurl, _ := c.String("amqp.url")
	if amqpurl != "" {
		log.Println("Starting AMQP client")

		err := amqp.ConnectAMQP(amqpurl)
		if err != nil {
			log.Panic(err)
		}

		// register simulator action endpoint only if AMQP client is used
		amqp.RegisterAMQPEndpoint(api.Group("/simulators"))

		// Periodically call the Ping function to check which simulators are still there
		ticker := time.NewTicker(10 * time.Second)
		go func() {

			for {
				select {
				case <-ticker.C:
					err = amqp.PingAMQP()
					if err != nil {
						log.Println("AMQP Error: ", err.Error())
					}
				}
			}

		}()

		log.Printf("Connected AMQP client to %s", amqpurl)
	}
	// server at port 4000 to match frontend's redirect path
	r.Run(":4000")
}
