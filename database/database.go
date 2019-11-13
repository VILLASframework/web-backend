package database

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/zpatrick/go-config"
)

var DBpool *gorm.DB // database used by backend

// Initialize connection to the database
func InitDB(cfg *config.Config) *gorm.DB {
	name, err := cfg.String("db.name")
	host, err := cfg.String("db.host")
	user, err := cfg.String("db.user")
	pass, err := cfg.String("db.pass")
	sslmode, err := cfg.String("db.ssl")
	init, err := cfg.Bool("db.init")
	mode, err := cfg.String("mode")

	dbinfo := fmt.Sprintf("host=%s sslmode=%s dbname=%s", host, sslmode, name)
	if user != "" && pass != "" {
		dbinfo += fmt.Sprintf(" user=%s password=%s", user, pass)
	}

	db, err := gorm.Open("postgres", dbinfo)
	if err != nil {
		log.Fatal(err)
	}
	DBpool = db

	MigrateModels(db)

	if mode == "test" || init {
		DropTables(db)
		log.Println("Database tables dropped")

		DBAddTestData(db)
		log.Println("Database initialized with test data")
	}

	log.Println("Database connection established")

	return db
}

// Connection pool to already opened DB
func GetDB() *gorm.DB {
	return DBpool
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
	// The following statement deletes the many to many relationship between users and scenarios
	db.DropTableIfExists("user_scenarios")
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
