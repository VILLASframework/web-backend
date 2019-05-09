package main

import (
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/user"

	"github.com/gin-gonic/gin"
)

func main() {
	// Testing
	db := common.InitDB()
	common.MigrateModels(db)
	defer db.Close()

	r := gin.Default()

	api := r.Group("/api")
	user.UsersRegister(api.Group("/users"))

	r.Run()
}
