package common

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
)

const (
	DB_NAME    = "villasdb"
	DB_DUMMY   = "testvillasdb"
	DB_HOST    = "/tmp"
	DB_SSLMODE = "disable" // TODO: change that for production
)

// Initialize connection to the database
func InitDB() *gorm.DB {
	dbinfo := fmt.Sprintf("host=%s sslmode=%s dbname=%s",
		DB_HOST, DB_SSLMODE, DB_NAME)
	db, err := gorm.Open("postgres", dbinfo)
	checkErr(err)
	return db
}

// Verify that the database connection is alive
func VerifyConnection(db *gorm.DB) error {
	return db.DB().Ping()
}

// Drop all the tables of the database
// TODO: Remove that function from the codebase and substitute the body
// to the Dummy*() where it is called
func DropTables(db *gorm.DB) {
	db.DropTableIfExists(&Simulator{})
	db.DropTableIfExists(&Signal{})
	db.DropTableIfExists(&SimulationModel{})
	db.DropTableIfExists(&File{})
	db.DropTableIfExists(&Project{})
	db.DropTableIfExists(&Simulation{})
	db.DropTableIfExists(&User{})
	db.DropTableIfExists(&Visualization{})
	db.DropTableIfExists(&Widget{})
}

// AutoMigrate the models
func MigrateModels(db *gorm.DB) {
	db.AutoMigrate(&Simulator{})
	db.AutoMigrate(&Signal{})
	db.AutoMigrate(&SimulationModel{})
	db.AutoMigrate(&File{})
	db.AutoMigrate(&Project{})
	db.AutoMigrate(&Simulation{})
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Visualization{})
	db.AutoMigrate(&Widget{})
}

// Start a dummy database for testing
func DummyInitDB() *gorm.DB {

	dbinfo := fmt.Sprintf("host=%s sslmode=%s dbname=%s",
		DB_HOST, DB_SSLMODE, DB_DUMMY)
	test_db, err := gorm.Open("postgres", dbinfo)
	checkErr(err)

	// drop tables from previous tests
	DropTables(test_db)

	return test_db
}

// Migrates models and populates them with data
func DummyPopulateDB(test_db *gorm.DB) {

	MigrateModels(test_db)

	// Create two entries of each model

	simr_A := Simulator{UUID: "1", Host: "Host_A"}
	simr_B := Simulator{UUID: "2", Host: "Host_B"}
	checkErr(test_db.Create(&simr_A).Error)
	checkErr(test_db.Create(&simr_B).Error)

	sig_A := Signal{Name: "Signal_A"}
	sig_B := Signal{Name: "Signal_B"}
	checkErr(test_db.Create(&sig_A).Error)
	checkErr(test_db.Create(&sig_B).Error)

	smo_A := SimulationModel{Name: "SimModel_A"}
	smo_B := SimulationModel{Name: "SimModel_B"}
	checkErr(test_db.Create(&smo_A).Error)
	checkErr(test_db.Create(&smo_B).Error)

	file_A := File{Name: "File_A"}
	file_B := File{Name: "File_B"}
	checkErr(test_db.Create(&file_A).Error)
	checkErr(test_db.Create(&file_B).Error)

	proj_A := Project{Name: "Proj_A"}
	proj_B := Project{Name: "Proj_B"}
	checkErr(test_db.Create(&proj_A).Error)
	checkErr(test_db.Create(&proj_B).Error)

	simn_A := Simulation{Name: "Simulation_A"}
	simn_B := Simulation{Name: "Simulation_B"}
	checkErr(test_db.Create(&simn_A).Error)
	checkErr(test_db.Create(&simn_B).Error)

	usr_A := User{Username: "User_A"}
	usr_B := User{Username: "User_B"}
	checkErr(test_db.Create(&usr_A).Error)
	checkErr(test_db.Create(&usr_B).Error)

	vis_A := Visualization{Name: "Visualization_A"}
	vis_B := Visualization{Name: "Visualization_B"}
	checkErr(test_db.Create(&vis_A).Error)
	checkErr(test_db.Create(&vis_B).Error)

	widg_A := Widget{Name: "Widget_A"}
	widg_B := Widget{Name: "Widget_B"}
	checkErr(test_db.Create(&widg_A).Error)
	checkErr(test_db.Create(&widg_B).Error)

	// Associations betweend models
	// For `belongs to` use the model with id=1
	// For `has many` use the models with id=1 and id=2

	checkErr(test_db.Model(&smo_A).Association("BelongsToSimulation").Append(&simn_A).Error)
	checkErr(test_db.Model(&smo_A).Association("BelongsToSimulator").Append(&simr_A).Error)
	checkErr(test_db.Model(&smo_A).Association("OutputMapping").Append(&sig_A).Error)
	checkErr(test_db.Model(&smo_A).Association("OutputMapping").Append(&sig_B).Error)
	checkErr(test_db.Model(&smo_A).Association("InputMapping").Append(&sig_B).Error)
	checkErr(test_db.Model(&smo_A).Association("InputMapping").Append(&sig_A).Error)

	checkErr(test_db.Model(&simn_A).Association("User").Append(&usr_A).Error)
	checkErr(test_db.Model(&simn_A).Association("Models").Append(&smo_A).Error)
	checkErr(test_db.Model(&simn_A).Association("Models").Append(&smo_B).Error)
	checkErr(test_db.Model(&simn_A).Association("Projects").Append(&proj_A).Error)
	checkErr(test_db.Model(&simn_A).Association("Projects").Append(&proj_B).Error)

	checkErr(test_db.Model(&proj_A).Association("Simulation").Append(&simn_A).Error)
	checkErr(test_db.Model(&proj_A).Association("User").Append(&usr_A).Error)
	checkErr(test_db.Model(&proj_A).Association("Visualizations").Append(&vis_A).Error)
	checkErr(test_db.Model(&proj_A).Association("Visualizations").Append(&vis_B).Error)

	checkErr(test_db.Model(&usr_A).Association("Projects").Append(&proj_A).Error)
	checkErr(test_db.Model(&usr_A).Association("Projects").Append(&proj_B).Error)
	checkErr(test_db.Model(&usr_A).Association("Simulations").Append(&simn_A).Error)
	checkErr(test_db.Model(&usr_A).Association("Simulations").Append(&simn_B).Error)
	checkErr(test_db.Model(&usr_A).Association("Files").Append(&file_A).Error)
	checkErr(test_db.Model(&usr_A).Association("Files").Append(&file_B).Error)

	checkErr(test_db.Model(&vis_A).Association("Project").Append(&proj_A).Error)
	checkErr(test_db.Model(&vis_A).Association("User").Append(&usr_A).Error)
	checkErr(test_db.Model(&vis_A).Association("Widgets").Append(&widg_A).Error)
	checkErr(test_db.Model(&vis_A).Association("Widgets").Append(&widg_B).Error)

	checkErr(test_db.Model(&file_A).Association("User").Append(&usr_A).Error)
}

// Erase tables and glose the testdb
func DummyCloseDB(test_db *gorm.DB) {
	test_db.Close()
}

// Quick error check
// NOTE: maybe this is not a good idea
func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
