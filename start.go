package main

import (
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/file"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	_ "git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/doc/api" // doc/api folder is used by Swag CLI, you have to import it
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/simulation"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/simulationmodel"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/simulator"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/user"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/visualization"
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
	// Testing
	db := common.DummyInitDB()
	common.MigrateModels(db)
	defer db.Close()

	common.DummyPopulateDB(db)

	r := gin.Default()

	api := r.Group("/api/v2")

	// All endpoints require authentication except when someone wants to
	// login (POST /authenticate)
	user.VisitorAuthenticate(api.Group("/authenticate"))

	api.Use(user.Authentication(true))

	simulation.RegisterSimulationEndpoints(api.Group("/simulations"))
	simulationmodel.RegisterSimulationModelEndpoints(api.Group("/models"))
	visualization.RegisterVisualizationEndpoints(api.Group("/visualizations"))
	widget.RegisterWidgetEndpoints(api.Group("/widgets"))
	file.RegisterFileEndpoints(api.Group("/files"))
	user.RegisterUserEndpoints(api.Group("/users"))
	simulator.RegisterSimulatorEndpoints(api.Group("/simulators"))

	r.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// server at port 4000 to match frontend's redirect path
	r.Run(":4000")
}
