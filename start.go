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
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/configuration"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
	apidocs "git.rwth-aachen.de/acs/public/villas/web-backend-go/doc/api" // doc/api folder is used by Swag CLI, you have to import it
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/helper"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/infrastructure-component"
	"github.com/gin-gonic/gin"
	"log"
)

func addData(router *gin.Engine, mode string, basePath string, amqphost string) error {

	if mode == "test" {
		// test mode: drop all tables and add test data to DB
		database.DropTables()
		log.Println("Database tables dropped, using API to add test data")
		resp, err := routes.AddTestData(basePath, router, amqphost)
		if err != nil {
			fmt.Println("error: testdata could not be added to DB:", err.Error(), "Response body: ", resp)
			return err
		}
	} else {
		// release mode: make sure that at least one admin user exists in DB
		err := helper.DBAddAdminUser()
		if err != nil {
			fmt.Println("error: adding admin user failed:", err.Error())
			return err
		}
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

	mode, baseHost, basePath, port, amqphost, amqpuser, amqppass, err := configuration.ConfigureBackend()
	if err != nil {
		panic(err)
	}

	//init database
	err = database.InitDB(configuration.GolbalConfig)
	if err != nil {
		log.Printf("Error during initialization of database: %v, aborting.", err.Error())
		panic(err)
	}
	defer database.DBpool.Close()

	// init endpoints
	if mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.Default()
	api := r.Group(basePath)
	routes.RegisterEndpoints(r, api)
	apidocs.SwaggerInfo.Host = baseHost
	apidocs.SwaggerInfo.BasePath = basePath

	// add data to DB (if any)
	err = addData(r, mode, basePath, amqphost)
	if err != nil {
		panic(err)
	}

	//Start AMQP client
	if amqphost != "" {
		// create amqp URL based on username, password and host
		amqpurl := "amqp://" + amqpuser + ":" + amqppass + "@" + amqphost
		err = infrastructure_component.StartAMQP(amqpurl, api)
		if err != nil {
			panic(err)
		}
	}

	// server at port 4000 to match frontend's redirect path
	r.Run(":" + port)
}
