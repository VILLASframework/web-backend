/** Main package.
*
* @author Sonja Happ <sonja.happ@eonerc.rwth-aachen.de>
* @copyright 2014-2019, Institute for Automation of Complex Power Systems, EONERC
* @license GNU General Public License (version 3)
*
* VILLASweb-backend-go
*
* This program is free software: you can redistribute it and/or modify
* it under the terms of the GNU General Public License as published by
* the Free Software Foundation, either version 3 of the License, or
* any later version.
*
* This program is distributed in the hope that it will be useful,
* but WITHOUT ANY WARRANTY; without even the implied warranty of
* MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
* GNU General Public License for more details.
*
* You should have received a copy of the GNU General Public License
* along with this program.  If not, see <http://www.gnu.org/licenses/>.
*********************************************************************************/
package main

import (
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/healthz"
	"log"
	"time"

	"git.rwth-aachen.de/acs/public/villas/web-backend-go/amqp"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	"git.rwth-aachen.de/acs/public/villas/web-backend-go/configuration"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
	apidocs "git.rwth-aachen.de/acs/public/villas/web-backend-go/doc/api" // doc/api folder is used by Swag CLI, you have to import it
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/dashboard"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/file"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/infrastructure-component"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/metrics"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/scenario"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/signal"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/simulationmodel"
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
	db, err := database.InitDB(configuration.GolbalConfig)
	if err != nil {
		log.Printf("Error during initialization of database: %v, aborting.", err.Error())
		return
	}
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
	port, err := configuration.GolbalConfig.String("port")
	if err != nil {
		log.Printf("Error reading port from global configuration: %v, aborting.", err.Error())
		return
	}

	apidocs.SwaggerInfo.Host = baseHost
	apidocs.SwaggerInfo.BasePath = basePath

	metrics.InitCounters(db)

	r := gin.Default()

	api := r.Group(basePath)

	healthz.RegisterHealthzEndpoint(api.Group("/healthz"))
	metrics.RegisterMetricsEndpoint(api.Group("/metrics"))
	// All endpoints (except for /healthz and /metrics) require authentication except when someone wants to
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
	infrastructure_component.RegisterICEndpoints(api.Group("/ic"))

	r.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	amqpurl, _ := configuration.GolbalConfig.String("amqp.url")
	if amqpurl != "" {
		log.Println("Starting AMQP client")

		err := amqp.ConnectAMQP(amqpurl)
		if err != nil {
			log.Panic(err)
		}

		// register IC action endpoint only if AMQP client is used
		amqp.RegisterAMQPEndpoint(api.Group("/ic"))

		// Periodically call the Ping function to check which ICs are still there
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
	r.Run(":" + port)
}
