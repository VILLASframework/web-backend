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
	"fmt"
	component_configuration "git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/component-configuration"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/dashboard"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/file"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/healthz"
	infrastructure_component "git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/infrastructure-component"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/scenario"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/signal"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/user"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/widget"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"log"
	"time"

	"git.rwth-aachen.de/acs/public/villas/web-backend-go/amqp"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/configuration"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
	apidocs "git.rwth-aachen.de/acs/public/villas/web-backend-go/doc/api" // doc/api folder is used by Swag CLI, you have to import it
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/metrics"
	"github.com/gin-gonic/gin"
)

func configureBackend() (string, string, string, string, string, string, string, error) {

	err := configuration.InitConfig()
	if err != nil {
		log.Printf("Error during initialization of global configuration: %v, aborting.", err.Error())
		return "", "", "", "", "", "", "", err
	}

	err = database.InitDB(configuration.GolbalConfig)
	if err != nil {
		log.Printf("Error during initialization of database: %v, aborting.", err.Error())
		return "", "", "", "", "", "", "", err
	}

	mode, err := configuration.GolbalConfig.String("mode")
	if err != nil {
		log.Printf("Error reading mode from global configuration: %v, aborting.", err.Error())
		return "", "", "", "", "", "", "", err
	}

	if mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	baseHost, err := configuration.GolbalConfig.String("base.host")
	if err != nil {
		log.Printf("Error reading base.host from global configuration: %v, aborting.", err.Error())
		return "", "", "", "", "", "", "", err
	}
	basePath, err := configuration.GolbalConfig.String("base.path")
	if err != nil {
		log.Printf("Error reading base.path from global configuration: %v, aborting.", err.Error())
		return "", "", "", "", "", "", "", err
	}
	port, err := configuration.GolbalConfig.String("port")
	if err != nil {
		log.Printf("Error reading port from global configuration: %v, aborting.", err.Error())
		return "", "", "", "", "", "", "", err
	}

	apidocs.SwaggerInfo.Host = baseHost
	apidocs.SwaggerInfo.BasePath = basePath

	metrics.InitCounters()

	AMQPhost, _ := configuration.GolbalConfig.String("amqp.host")
	AMQPuser, _ := configuration.GolbalConfig.String("amqp.user")
	AMQPpass, _ := configuration.GolbalConfig.String("amqp.pass")

	return mode, baseHost, basePath, port, AMQPhost, AMQPuser, AMQPpass, nil

}

func registerEndpoints(router *gin.Engine, api *gin.RouterGroup) {

	healthz.RegisterHealthzEndpoint(api.Group("/healthz"))
	metrics.RegisterMetricsEndpoint(api.Group("/metrics"))
	// All endpoints (except for /healthz and /metrics) require authentication except when someone wants to
	// login (POST /authenticate)
	user.RegisterAuthenticate(api.Group("/authenticate"))

	api.Use(user.Authentication(true))

	scenario.RegisterScenarioEndpoints(api.Group("/scenarios"))
	component_configuration.RegisterComponentConfigurationEndpoints(api.Group("/configs"))
	signal.RegisterSignalEndpoints(api.Group("/signals"))
	dashboard.RegisterDashboardEndpoints(api.Group("/dashboards"))
	widget.RegisterWidgetEndpoints(api.Group("/widgets"))
	file.RegisterFileEndpoints(api.Group("/files"))
	user.RegisterUserEndpoints(api.Group("/users"))
	infrastructure_component.RegisterICEndpoints(api.Group("/ic"))

	router.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

}

func addData(router *gin.Engine, mode string, basePath string) error {

	if mode == "test" {
		// test mode: drop all tables and add test data to DB
		database.DropTables()
		log.Println("Database tables dropped, adding test data to DB")
		err := database.DBAddTestData(basePath, router)
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("error: testdata could not be added to DB, aborting")
			return err
		}
		log.Println("Database initialized with test data")
	} else {
		// release mode: make sure that at least one admin user exists in DB
		err := database.DBAddAdminUser()
		if err != nil {
			fmt.Println(err.Error())
			fmt.Println("error: adding admin user failed, aborting")
			return err
		}
	}
	return nil
}

func connectAMQP(AMQPurl string, api *gin.RouterGroup) error {
	if AMQPurl != "" {
		log.Println("Starting AMQP client")

		err := amqp.ConnectAMQP(AMQPurl)
		if err != nil {
			return err
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

		log.Printf("Connected AMQP client to %s", AMQPurl)
	}

	return nil
}

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

	mode, _, basePath, port, amqphost, amqpuser, amqppass, err := configureBackend()
	if err != nil {
		panic(err)
	}
	defer database.DBpool.Close()

	r := gin.Default()
	api := r.Group(basePath)
	registerEndpoints(r, api)

	err = addData(r, mode, basePath)
	if err != nil {
		panic(err)
	}

	// create amqp URL based on username, password and host
	amqpurl := amqpuser + ":" + amqppass + "@" + amqphost
	err = connectAMQP(amqpurl, api)
	if err != nil {
		panic(err)
	}

	// server at port 4000 to match frontend's redirect path
	r.Run(":" + port)
}
