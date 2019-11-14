package main

import (
	"log"
	"time"

	"git.rwth-aachen.de/acs/public/villas/web-backend-go/amqp"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/healthz"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	"git.rwth-aachen.de/acs/public/villas/web-backend-go/configuration"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
	apidocs "git.rwth-aachen.de/acs/public/villas/web-backend-go/doc/api" // doc/api folder is used by Swag CLI, you have to import it
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/dashboard"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/file"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/metrics"
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

	err := configuration.InitConfig()
	if err != nil {
		log.Printf("Error during initialization of global configuration: %v, aborting.", err.Error())
		return
	}
	db := database.InitDB(configuration.GolbalConfig)
	defer db.Close()

	m, err := configuration.GolbalConfig.String("mode")
	if err != nil {
		log.Printf("Error reading mode from global configuration: %v, aborting.", err.Error())
		return
	}

	if m == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	baseHost, err := configuration.GolbalConfig.String("base.host")
	if err != nil {
		log.Printf("Error reading base.host from global configuration: %v, aborting.", err.Error())
		return
	}
	basePath, err := configuration.GolbalConfig.String("base.path")
	if err != nil {
		log.Printf("Error reading base.path from global configuration: %v, aborting.", err.Error())
		return
	}

	apidocs.SwaggerInfo.Host = baseHost
	apidocs.SwaggerInfo.BasePath = basePath

	metrics.InitCounters(db)

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
	healthz.RegisterHealthzEndpoint(r.Group("/healthz"))
	metrics.RegisterMetricsEndpoint(r.Group("/metrics"))

	r.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	amqpurl, _ := configuration.GolbalConfig.String("amqp.url")
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
