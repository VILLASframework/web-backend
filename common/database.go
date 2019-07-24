package common

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/crypto/bcrypt"
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

	// Create two entries of each model

	propertiesA := json.RawMessage(`{"name" : "TestNameA", "category" : "CategoryA", "location" : "anywhere on earth", "type": "dummy"}`)
	propertiesB := json.RawMessage(`{"name" : "TestNameB", "category" : "CategoryB", "location" : "where ever you want", "type": "generic"}`)
	simr_A := Simulator{UUID: "4854af30-325f-44a5-ad59-b67b2597de68", Host: "Host_A", State: "running", Modeltype: "ModelTypeA", StateUpdateAt: "placeholder", Properties: postgres.Jsonb{propertiesA}, RawProperties: postgres.Jsonb{propertiesA}}
	simr_B := Simulator{UUID: "7be0322d-354e-431e-84bd-ae4c9633138b", Host: "Host_B", State: "idle", Modeltype: "ModelTypeB", StateUpdateAt: "placeholder", Properties: postgres.Jsonb{propertiesB}, RawProperties: postgres.Jsonb{propertiesB}}
	checkErr(test_db.Create(&simr_A).Error)
	checkErr(test_db.Create(&simr_B).Error)

	outSig_A := Signal{Name: "outSignal_A", Direction: "out", Index: 0, Unit: "V"}
	outSig_B := Signal{Name: "outSignal_B", Direction: "out", Index: 1, Unit: "V"}
	inSig_A := Signal{Name: "inSignal_A", Direction: "in", Index: 0, Unit: "A"}
	inSig_B := Signal{Name: "inSignal_B", Direction: "in", Index: 1, Unit: "A"}
	checkErr(test_db.Create(&outSig_A).Error)
	checkErr(test_db.Create(&outSig_B).Error)
	checkErr(test_db.Create(&inSig_A).Error)
	checkErr(test_db.Create(&inSig_B).Error)

	mo_A := SimulationModel{Name: "SimulationModel_A"}
	mo_B := SimulationModel{Name: "SimulationModel_B"}
	checkErr(test_db.Create(&mo_A).Error)
	checkErr(test_db.Create(&mo_B).Error)

	file_A := File{Name: "File_A"}
	file_B := File{Name: "File_B"}
	file_C := File{Name: "File_C"}
	file_D := File{Name: "File_D"}
	checkErr(test_db.Create(&file_A).Error)
	checkErr(test_db.Create(&file_B).Error)
	checkErr(test_db.Create(&file_C).Error)
	checkErr(test_db.Create(&file_D).Error)

	so_A := Scenario{Name: "Scenario_A"}
	so_B := Scenario{Name: "Scenario_B"}
	checkErr(test_db.Create(&so_A).Error)
	checkErr(test_db.Create(&so_B).Error)

	// Hash passwords with bcrypt algorithm
	var bcryptCost = 10

	pw_0, err :=
		bcrypt.GenerateFromPassword([]byte("xyz789"), bcryptCost)
	checkErr(err)

	pw_A, err :=
		bcrypt.GenerateFromPassword([]byte("abc123"), bcryptCost)
	checkErr(err)

	pw_B, err :=
		bcrypt.GenerateFromPassword([]byte("bcd234"), bcryptCost)
	checkErr(err)

	usr_0 := User{Username: "User_0", Password: string(pw_0), Role: "Admin"}
	usr_A := User{Username: "User_A", Password: string(pw_A), Role: "User"}
	usr_B := User{Username: "User_B", Password: string(pw_B), Role: "User"}
	checkErr(test_db.Create(&usr_0).Error)
	checkErr(test_db.Create(&usr_A).Error)
	checkErr(test_db.Create(&usr_B).Error)

	dab_A := Dashboard{Name: "Dashboard_A"}
	dab_B := Dashboard{Name: "Dashboard_B"}
	checkErr(test_db.Create(&dab_A).Error)
	checkErr(test_db.Create(&dab_B).Error)

	widg_A := Widget{Name: "Widget_A"}
	widg_B := Widget{Name: "Widget_B"}
	checkErr(test_db.Create(&widg_A).Error)
	checkErr(test_db.Create(&widg_B).Error)

	// Associations between models
	// For `belongs to` use the model with id=1
	// For `has many` use the models with id=1 and id=2

	// User HM Scenarios, Scenario HM Users (Many-to-Many)
	checkErr(test_db.Model(&so_A).Association("Users").Append(&usr_A).Error)
	checkErr(test_db.Model(&so_A).Association("Users").Append(&usr_B).Error)
	checkErr(test_db.Model(&so_B).Association("Users").Append(&usr_A).Error)
	checkErr(test_db.Model(&so_B).Association("Users").Append(&usr_B).Error)

	// Scenario HM SimulationModels
	checkErr(test_db.Model(&so_A).Association("SimulationModels").Append(&mo_A).Error)
	checkErr(test_db.Model(&so_A).Association("SimulationModels").Append(&mo_B).Error)

	// Scenario HM Dashboards
	checkErr(test_db.Model(&so_A).Association("Dashboards").Append(&dab_A).Error)
	checkErr(test_db.Model(&so_A).Association("Dashboards").Append(&dab_B).Error)

	// Dashboard HM Widget
	checkErr(test_db.Model(&dab_A).Association("Widgets").Append(&widg_A).Error)
	checkErr(test_db.Model(&dab_A).Association("Widgets").Append(&widg_B).Error)

	// SimulationModel HM Signals
	checkErr(test_db.Model(&mo_A).Association("InputMapping").Append(&inSig_A).Error)
	checkErr(test_db.Model(&mo_A).Association("InputMapping").Append(&inSig_B).Error)
	checkErr(test_db.Model(&mo_A).Association("OutputMapping").Append(&outSig_A).Error)
	checkErr(test_db.Model(&mo_A).Association("OutputMapping").Append(&outSig_B).Error)

	// SimulationModel HM Files
	checkErr(test_db.Model(&mo_A).Association("Files").Append(&file_C).Error)
	checkErr(test_db.Model(&mo_A).Association("Files").Append(&file_D).Error)

	// Simulator HM SimulationModels
	checkErr(test_db.Model(&simr_A).Association("SimulationModels").Append(&mo_A).Error)
	checkErr(test_db.Model(&simr_A).Association("SimulationModels").Append(&mo_B).Error)

	// Widget HM Files
	checkErr(test_db.Model(&widg_A).Association("Files").Append(&file_A).Error)
	checkErr(test_db.Model(&widg_A).Association("Files").Append(&file_B).Error)

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
