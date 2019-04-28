package main

import (
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"gopkg.in/gin-gonic/gin.v1"
)

func main() {
	// Testing
	db := common.InitDB()
	common.MigrateModels(db)
	defer db.Close()

	r := gin.Default()

	api := r.Group("/api")
	common.UsersRegister(api.Group("/users"))

	r.Run()
}
