package main

import (
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/file"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/project"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/simulation"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/simulationmodel"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/simulator"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/user"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/visualization"

	"github.com/gin-gonic/gin"
)

func main() {
	// Testing
	db := common.DummyInitDB()
	common.MigrateModels(db)
	defer db.Close()

	common.DummyPopulateDB(db)

	r := gin.Default()

	api := r.Group("/api/v1")

	// All endpoints require authentication except when someone wants to
	// login (POST /authenticate)
	user.VisitorAuthenticate(api.Group("/authenticate"))

	api.Use(user.Authentication(true))

	user.UsersRegister(api.Group("/users"))
	file.FilesRegister(api.Group("/files"))
	project.ProjectsRegister(api.Group("/projects"))
	simulation.SimulationsRegister(api.Group("/simulations"))
	simulationmodel.SimulationModelsRegister(api.Group("/models"))
	simulator.SimulatorsRegister(api.Group("/simulators"))
	visualization.VisualizationsRegister(api.Group("/visualizations"))

	// server at port 4000 to match frontend's redirect path
	r.Run(":4000")
}
