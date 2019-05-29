package common

import (
	"flag"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/crypto/bcrypt"
	"log"
)

var DB_HOST string
var DB_NAME string
var DB_DUMMY string
var DB_SSLMODE string

var DBpool *gorm.DB

// Initialize input command line flags
func init() {
	flag.StringVar(&DB_HOST, "dbhost", "/tmp", "Host of the PostgreSQL database (default is /tmp)")
	flag.StringVar(&DB_NAME, "dbname", "villasdb", "Name of the database to use (default is villasdb)")
	flag.StringVar(&DB_DUMMY, "dbdummy", "testvillasdb", "Name of the test database to use (default is testvillasdb)")
	flag.StringVar(&DB_SSLMODE, "dbsslmode", "disable", "SSL mode of DB (default is disable)") // TODO: change default for production
	flag.Parse()
	fmt.Println("DB_HOST has value ", DB_HOST)
	fmt.Println("DB_NAME has value ", DB_NAME)
	fmt.Println("DB_DUMMY has value ", DB_DUMMY)
	fmt.Println("DB_SSLMODE has value ", DB_SSLMODE)
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
	//db.DropTableIfExists(&Signal{})
	db.DropTableIfExists(&Model{})
	db.DropTableIfExists(&File{})
	db.DropTableIfExists(&Simulation{})
	db.DropTableIfExists(&User{})
	db.DropTableIfExists(&Visualization{})
	db.DropTableIfExists(&Widget{})
}

// AutoMigrate the models
func MigrateModels(db *gorm.DB) {
	db.AutoMigrate(&Simulator{})
	//db.AutoMigrate(&Signal{})
	db.AutoMigrate(&Model{})
	db.AutoMigrate(&File{})
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
	DBpool = test_db
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

	//outSig_A := Signal{Name: "outSignal_A", Direction: "out"}
	//outSig_B := Signal{Name: "outSignal_B", Direction: "out"}
	//inSig_A := Signal{Name: "inSignal_A", Direction: "in"}
	//inSig_B := Signal{Name: "inSignal_B", Direction: "in"}
	//checkErr(test_db.Create(&outSig_A).Error)
	//checkErr(test_db.Create(&outSig_B).Error)
	//checkErr(test_db.Create(&inSig_A).Error)
	//checkErr(test_db.Create(&inSig_B).Error)

	mo_A := Model{Name: "Model_A"}
	mo_B := Model{Name: "Model_B"}
	checkErr(test_db.Create(&mo_A).Error)
	checkErr(test_db.Create(&mo_B).Error)

	file_A := File{Name: "File_A"}
	file_B := File{Name: "File_B"}
	checkErr(test_db.Create(&file_A).Error)
	checkErr(test_db.Create(&file_B).Error)

	simn_A := Simulation{Name: "Simulation_A"}
	simn_B := Simulation{Name: "Simulation_B"}
	checkErr(test_db.Create(&simn_A).Error)
	checkErr(test_db.Create(&simn_B).Error)

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

	// User HM Simulations, Simulation HM Users (Many-to-Many)
	checkErr(test_db.Model(&simn_A).Association("Users").Append(&usr_A).Error)
	checkErr(test_db.Model(&simn_A).Association("Users").Append(&usr_B).Error)
	checkErr(test_db.Model(&simn_B).Association("Users").Append(&usr_A).Error)
	checkErr(test_db.Model(&simn_B).Association("Users").Append(&usr_B).Error)

	// Simulation HM Model
	checkErr(test_db.Model(&simn_A).Association("Models").Append(&mo_A).Error)
	checkErr(test_db.Model(&simn_A).Association("Models").Append(&mo_B).Error)

	// Simulation HM Visualizations
	checkErr(test_db.Model(&simn_A).Association("Visualizations").Append(&vis_A).Error)
	checkErr(test_db.Model(&simn_A).Association("Visualizations").Append(&vis_B).Error)

	// Visualization HM Widget
	checkErr(test_db.Model(&vis_A).Association("Widgets").Append(&widg_A).Error)
	checkErr(test_db.Model(&vis_A).Association("Widgets").Append(&widg_B).Error)

	// Model HM Signal
	//checkErr(test_db.Model(&mo_A).Association("InputMapping").Append(&inSig_A).Error)
	//checkErr(test_db.Model(&mo_A).Association("InputMapping").Append(&inSig_B).Error)
	//checkErr(test_db.Model(&mo_A).Association("OutputMapping").Append(&outSig_A).Error)
	//checkErr(test_db.Model(&mo_A).Association("OutputMapping").Append(&outSig_B).Error)

	// Model HM Files
	checkErr(test_db.Model(&mo_A).Association("Files").Append(&file_A).Error)
	checkErr(test_db.Model(&mo_A).Association("Files").Append(&file_B).Error)

	// Simulator BT Model
	checkErr(test_db.Model(&mo_A).Association("Simulator").Append(&simr_A).Error)

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
