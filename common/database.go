package common

import (
	"flag"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
)

var DB_HOST string
var DB_NAME string
var DB_DUMMY string
var DB_SSLMODE string
var WITH_AMQP bool

var DBpool *gorm.DB

// Initialize input command line flags
func init() {
	flag.StringVar(&DB_HOST, "dbhost", "/tmp", "Host of the PostgreSQL database (default is /tmp)")
	flag.StringVar(&DB_NAME, "dbname", "villasdb", "Name of the database to use (default is villasdb)")
	flag.StringVar(&DB_DUMMY, "dbdummy", "testvillasdb", "Name of the test database to use (default is testvillasdb)")
	flag.StringVar(&DB_SSLMODE, "dbsslmode", "disable", "SSL mode of DB (default is disable)") // TODO: change default for production
	flag.BoolVar(&WITH_AMQP, "amqp", false, "If AMQP client for simulators shall be enabled, set this option to true (default is false)")
	flag.Parse()
	fmt.Println("DB_HOST has value ", DB_HOST)
	fmt.Println("DB_NAME has value ", DB_NAME)
	fmt.Println("DB_DUMMY has value ", DB_DUMMY)
	fmt.Println("DB_SSLMODE has value ", DB_SSLMODE)
	fmt.Println("WITH_AMQP has value ", WITH_AMQP)
}

// Initialize connection to the database
func InitDB() *gorm.DB {
	dbinfo := fmt.Sprintf("host=%s sslmode=%s dbname=%s",
		DB_HOST, DB_SSLMODE, DB_NAME)
	db, err := gorm.Open("postgres", dbinfo)
	checkErr(err)
	DBpool = db
	return db
}

// Connection pool to already opened DB
func GetDB() *gorm.DB {
	return DBpool
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
	db.DropTableIfExists(&Scenario{})
	db.DropTableIfExists(&User{})
	db.DropTableIfExists(&Dashboard{})
	db.DropTableIfExists(&Widget{})
}

// AutoMigrate the models
func MigrateModels(db *gorm.DB) {
	db.AutoMigrate(&Simulator{})
	db.AutoMigrate(&Signal{})
	db.AutoMigrate(&SimulationModel{})
	db.AutoMigrate(&File{})
	db.AutoMigrate(&Scenario{})
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Dashboard{})
	db.AutoMigrate(&Widget{})
}

// Start a dummy database for testing
func DummyInitDB() *gorm.DB {

	dbinfo := fmt.Sprintf("host=%s sslmode=%s dbname=%s",
		DB_HOST, DB_SSLMODE, DB_DUMMY)
	test_db, err := gorm.Open("postgres", dbinfo)
	checkErr(err)
	DBpool = test_db
	// drop tables from previous tests
	DropTables(test_db)

	return test_db
}

// Migrates models and populates them with data
func DummyPopulateDB(test_db *gorm.DB) {

	MigrateModels(test_db)

	// Create entries of each model (data defined in testdata.go)

	// Users
	checkErr(test_db.Create(&User0).Error) // Admin
	checkErr(test_db.Create(&UserA).Error) // Normal User
	checkErr(test_db.Create(&UserB).Error) // Normal User

	// Simulators
	checkErr(test_db.Create(&SimulatorA).Error)
	checkErr(test_db.Create(&SimulatorB).Error)

	// Scenarios
	checkErr(test_db.Create(&ScenarioA).Error)
	checkErr(test_db.Create(&ScenarioB).Error)

	// Signals
	checkErr(test_db.Create(&OutSignalA).Error)
	checkErr(test_db.Create(&OutSignalB).Error)
	checkErr(test_db.Create(&InSignalA).Error)
	checkErr(test_db.Create(&InSignalB).Error)

	// Simulation Models
	checkErr(test_db.Create(&SimulationModelA).Error)
	checkErr(test_db.Create(&SimulationModelB).Error)

	// Dashboards
	checkErr(test_db.Create(&DashboardA).Error)
	checkErr(test_db.Create(&DashboardB).Error)

	// Files
	checkErr(test_db.Create(&FileA).Error)
	checkErr(test_db.Create(&FileB).Error)
	checkErr(test_db.Create(&FileC).Error)
	checkErr(test_db.Create(&FileD).Error)

	widg_A := Widget{Name: "Widget_A"}
	widg_B := Widget{Name: "Widget_B"}
	checkErr(test_db.Create(&widg_A).Error)
	checkErr(test_db.Create(&widg_B).Error)

	// Associations between models
	// For `belongs to` use the model with id=1
	// For `has many` use the models with id=1 and id=2

	// User HM Scenarios, Scenario HM Users (Many-to-Many)
	checkErr(test_db.Model(&ScenarioA).Association("Users").Append(&UserA).Error)
	checkErr(test_db.Model(&ScenarioA).Association("Users").Append(&UserB).Error)
	checkErr(test_db.Model(&ScenarioB).Association("Users").Append(&UserA).Error)
	checkErr(test_db.Model(&ScenarioB).Association("Users").Append(&UserB).Error)

	// Scenario HM SimulationModels
	checkErr(test_db.Model(&ScenarioA).Association("SimulationModels").Append(&SimulationModelA).Error)
	checkErr(test_db.Model(&ScenarioA).Association("SimulationModels").Append(&SimulationModelB).Error)

	// Scenario HM Dashboards
	checkErr(test_db.Model(&ScenarioA).Association("Dashboards").Append(&DashboardA).Error)
	checkErr(test_db.Model(&ScenarioA).Association("Dashboards").Append(&DashboardB).Error)

	// Dashboard HM Widget
	checkErr(test_db.Model(&DashboardA).Association("Widgets").Append(&widg_A).Error)
	checkErr(test_db.Model(&DashboardA).Association("Widgets").Append(&widg_B).Error)

	// SimulationModel HM Signals
	checkErr(test_db.Model(&SimulationModelA).Association("InputMapping").Append(&InSignalA).Error)
	checkErr(test_db.Model(&SimulationModelA).Association("InputMapping").Append(&InSignalB).Error)
	checkErr(test_db.Model(&SimulationModelA).Association("OutputMapping").Append(&OutSignalA).Error)
	checkErr(test_db.Model(&SimulationModelA).Association("OutputMapping").Append(&OutSignalB).Error)

	// SimulationModel HM Files
	checkErr(test_db.Model(&SimulationModelA).Association("Files").Append(&FileC).Error)
	checkErr(test_db.Model(&SimulationModelA).Association("Files").Append(&FileD).Error)

	// Simulator HM SimulationModels
	checkErr(test_db.Model(&SimulatorA).Association("SimulationModels").Append(&SimulationModelA).Error)
	checkErr(test_db.Model(&SimulatorA).Association("SimulationModels").Append(&SimulationModelB).Error)

	// Widget HM Files
	checkErr(test_db.Model(&widg_A).Association("Files").Append(&FileA).Error)
	checkErr(test_db.Model(&widg_A).Association("Files").Append(&FileB).Error)
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
