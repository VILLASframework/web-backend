/**
* This file is part of VILLASweb-backend-go
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
	"log"
	"time"

	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/user"

	"git.rwth-aachen.de/acs/public/villas/web-backend-go/configuration"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/helper"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/healthz"
	infrastructure_component "git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/infrastructure-component"
	"github.com/gin-gonic/gin"
	"github.com/zpatrick/go-config"
)

func addData(router *gin.Engine, cfg *config.Config) error {

	testDataPath, err := cfg.String("test.datapath")
	if err != nil {
		// if param is missing, no test data will be added
		return nil
	}

	err = routes.ReadTestDataFromJson(testDataPath)
	if err != nil {
		log.Println("WARNING: test data cannot not be read from file, continue without it: ", err)
		return nil
	}

	resp, err := routes.AddTestData(cfg, router)
	if err != nil {
		fmt.Println("ERROR: test data could not be added to DB:", err.Error(), "Response body: ", resp)
		return err
	}

	return nil
}

// @title VILLASweb Backend API
// @version 2.0
// @description This is the [VILLASweb Backend](https://git.rwth-aachen.de/acs/public/villas/web-backend-go) API v2.0.
// @description This documentation is auto-generated based on the API documentation in the code. The tool [swag](https://github.com/swaggo/swag) is used to auto-generate API docs for the [gin-gonic](https://github.com/gin-gonic/gin) framework.
// @description Authentication: Use the authenticate endpoint below to obtain a token for your user account, copy the token into to the value field of the dialog showing up for the green Authorize button below and confirm with Done.
// @contact.name Institute for Automation of Complex Power Systems, RWTH Aachen University
// @license.name GNU GPL 3.0
// @license.url http://www.gnu.de/documents/gpl-3.0.en.html
// @BasePath /api/v2
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
func main() {
	log.Println("Starting VILLASweb-backend-go")

	err := configuration.InitConfig()
	if err != nil {
		log.Fatalf("Error during initialization of global configuration: %s, aborting.", err)
	}

	dbClear, err := configuration.GlobalConfig.String("db.clear")
	if err != nil {
		log.Fatalf("Error reading db.clear parameter from global configuration: %s, aborting.", err)
	}

	port, err := configuration.GlobalConfig.String("port")
	if err != nil {
		log.Fatalf("Error reading port from global configuration: %s, aborting.", err)
	}

	// Init database
	err = database.InitDB(configuration.GlobalConfig, dbClear == "true")
	if err != nil {
		log.Fatalf("Error during initialization of database: %s, aborting.", err)
	}
	defer database.DBpool.Close()

	// Init endpoints
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()
	api := r.Group("/api/v2")
	routes.RegisterEndpoints(r, api)

	// Start AMQP client
	amqpHost, _ := configuration.GlobalConfig.String("amqp.host")
	amqpUser, _ := configuration.GlobalConfig.String("amqp.user")
	amqpPass, _ := configuration.GlobalConfig.String("amqp.pass")

	if amqpHost != "" {
		// create amqp URL based on username, password and host
		amqpurl := "amqp://" + amqpUser + ":" + amqpPass + "@" + amqpHost
		session := helper.NewAMQPSession("villas-amqp-session", amqpurl, "villas", infrastructure_component.ProcessMessage)
		healthz.SetAMQPSession(session)                  // healthz needs to know the amqp session to check the health of the backend
		infrastructure_component.SetAMQPSession(session) // IC needs to know the session to send amqp messages
		user.SetAMQPSession(session)                     // User needs to know the session to duplicate ICs upon login

		// send Ping to all externally managed ICs
		for {
			if session.IsReady {
				err = infrastructure_component.SendPing("")
				if err != nil {
					log.Println("error sending ping action via AMQP: ", err.Error())
				}
				break
			}
		}
	}

	// Make sure that at least one admin user exists in DB
	_, err = database.AddAdminUser(configuration.GlobalConfig)
	if err != nil {
		fmt.Println("error: adding admin user failed:", err.Error())
		log.Fatal(err)
	}

	// Add test/demo data to DB (if any)
	err = addData(r, configuration.GlobalConfig)
	if err != nil {
		log.Fatal(err)
	}

	// Update via external APIs of ICs (if not managed via AMQP)
	intervalStr, _ := configuration.GlobalConfig.String("apiupdateinterval")
	interval, _ := time.ParseDuration(intervalStr)
	infrastructure_component.QueryICAPIs(interval)

	log.Println("Running...")
	// Server at port 4000 to match frontend's redirect path
	r.Run(":" + port)
}
