package main

import (
	"fmt"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/simulator"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
)

type Simulator simulator.Simulator

const (
	DB_NAME = "villasdb"
)

func main() {
	// Init connection's information
	dbinfo := fmt.Sprintf("host=/tmp sslmode=disable dbname=%s", DB_NAME)
	db, err := gorm.Open("postgres", dbinfo)
	checkErr(err)
	defer db.Close()

	// Check that db is reachable
	err = db.DB().Ping()
	checkErr(err)

	// Migrate one model
	db.AutoMigrate(&Simulator{})

	// Create
	db.Create(&Simulator{UUID: "12"})

	// Read
	var dummy Simulator
	db.First(&dummy, 1)
	fmt.Printf("%s\n", dummy.UUID)

	// Update
	db.Model(&dummy).Update("UUID", "100")
	db.First(&dummy, 1)
	fmt.Printf("%s\n", dummy.UUID)

	// Delete
	db.Unscoped().Delete(&dummy)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
