package common

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
)

const (
	DB_NAME    = "villasdb"
	DB_HOST    = "/tmp"
	DB_SSLMODE = "disable" // TODO: change that for production
)

func StartDB() {
	// Init connection's information
	dbinfo := fmt.Sprintf("host=%s sslmode=%s dbname=%s",
		DB_HOST, DB_SSLMODE, DB_NAME)
	db, err := gorm.Open("postgres", dbinfo)
	checkErr(err)
	defer db.Close()

	// Check that db is reachable
	err = db.DB().Ping()
	checkErr(err)

	// Migrate one model
	db.AutoMigrate(&Simulator{})
	db.AutoMigrate(&Signal{})
	db.AutoMigrate(&SimulationModel{})
	db.AutoMigrate(&File{})
	db.AutoMigrate(&Project{})
	db.AutoMigrate(&Simulation{})
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Visualization{})
	db.AutoMigrate(&Signal{})
	db.AutoMigrate(&Widget{})

	// Create
	db.Create(&Simulator{UUID: "12"})
	db.Create(&Signal{Name: "Some", Unit: "314"})
	fooSimMod := SimulationModel{Name: "buz",
		InputMapping: []Signal{
			{Name: "foo", Unit: "42"},
			{Name: "buz", Unit: "511"},
		},
	}
	db.Create(&fooSimMod)

	// get number of associations from SimulationModel table InputMapping column
	fmt.Println("Number of associations of InputMapping: ",
		db.Model(&fooSimMod).Association("InputMapping").Count())

	// get the associations from SimulationModel table InputMapping column
	var inSignals []Signal
	db.Model(&fooSimMod).Association("InputMapping").Find(&inSignals)
	fmt.Println(inSignals)

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
