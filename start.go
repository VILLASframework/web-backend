package main

import (
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/endpoints"
	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
)

func main() {
	// Testing
	db := common.InitDB()
	common.MigrateModels(db)
	defer db.Close()

	r := gin.Default()

	api := r.Group("/api")
	endpoints.UsersRegister(api.Group("/users"))
	//file.FilesRegister(api.Group("/files"))
	//project.ProjectsRegister(api.Group("/projects"))
	endpoints.SimulationsRegister(api.Group("/simulations"))
	//model.ModelsRegister(api.Group("/simulations"))
	endpoints.SimulatorsRegister(api.Group("/simulators"))
	//visualization.VisualizationsRegister(api.Group("/visualizations"))




	r.Run()
}
