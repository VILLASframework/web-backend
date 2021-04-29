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
	"log"

	"git.rwth-aachen.de/acs/public/villas/web-backend-go/configuration"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes"
	infrastructure_component "git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/infrastructure-component"
	"github.com/gin-gonic/gin"
	"github.com/zpatrick/go-config"
)

func addData(router *gin.Engine, cfg *config.Config) error {

	testDataPath, err := cfg.String("test.datapath")
	if err != nil {
		return err
	}

	err = routes.ReadTestDataFromJson(testDataPath)
	if err != nil {
		log.Println("testdata could not be read from json file")
		return err
	}

	resp, err := routes.AddTestData(cfg, router)
	if err != nil {
		fmt.Println("error: testdata could not be added to DB:", err.Error(), "Response body: ", resp)
		return err
	}

	return nil
}

// @title VILLASweb Backend API
// @version 2.0
// @description This is the [VILLASweb Backend](https://git.rwth-aachen.de/acs/public/villas/web-backend-go) API v2.0.
// @description This documentation is auto-generated based on the API documentation in the code. The tool [swag](https://github.com/swaggo/swag) is used to auto-generate API docs for the [gin-gonic](https://github.com/gin-gonic/gin) framework.
// @description Authentication: Use the authenticate endpoint below to obtain a token for your user account, copy the token into to the value field of the dialog showing up for the green Authorize button below and confirm with Done.
// @contact.name Sonja Happ
// @contact.email sonja.happ@eonerc.rwth-aachen.de
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

	gPath, err := configuration.GlobalConfig.String("groups.path")
	if err != nil {
		log.Fatalf("Error reading path to groups YAML file: %s, aborting.", err)
	}

	if gPath != "" {
		err = configuration.ReadGroupsFile(gPath)
		if err != nil {
			log.Fatalf("Error reading groups YAML file: %s, aborting.", err)
		}
	}

	// Init database
	err = database.InitDB(configuration.GlobalConfig, dbClear)
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
	AMQPhost, _ := configuration.GlobalConfig.String("amqp.host")
	AMQPuser, _ := configuration.GlobalConfig.String("amqp.user")
	AMQPpass, _ := configuration.GlobalConfig.String("amqp.pass")

	if AMQPhost != "" {
		// create amqp URL based on username, password and host
		amqpurl := "amqp://" + AMQPuser + ":" + AMQPpass + "@" + AMQPhost
		err = infrastructure_component.StartAMQP(amqpurl, api)
		if err != nil {
			log.Fatal(err)
		}
	}

	// Make sure that at least one admin user exists in DB
	err = database.DBAddAdminUser(configuration.GlobalConfig)
	if err != nil {
		fmt.Println("error: adding admin user failed:", err.Error())
		log.Fatal(err)
	}

	// Add test/demo data to DB (if any)
	err = addData(r, configuration.GlobalConfig)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Running...")
	// Server at port 4000 to match frontend's redirect path
	r.Run(":" + port)
}
