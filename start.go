package main

import (
	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/simulation"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/simulator"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/user"
)

func main() {
	// Testing
	db := common.InitDB()
	common.MigrateModels(db)
	defer db.Close()

	r := gin.Default()

	api := r.Group("/api")
	user.UsersRegister(api.Group("/users"))
	//file.FilesRegister(api.Group("/files"))
	//project.ProjectsRegister(api.Group("/projects"))
	simulation.SimulationsRegister(api.Group("/simulations"))
	//model.ModelsRegister(api.Group("/models"))
	simulator.SimulatorsRegister(api.Group("/simulators"))
	//visualization.VisualizationsRegister(api.Group("/visualizations"))




	r.Run()
}
