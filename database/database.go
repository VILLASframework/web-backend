package database

import (
	"flag"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"log"
)

var DB_HOST string    // host of the database system
var DB_NAME string    // name of the production database
var DB_USER string    // name of the database user
var DB_PASS string    // database password
var DB_SSLMODE string // set to enable if database uses SSL
var WITH_AMQP bool    // set to true if backend shall be used with AMQP client

var DBpool *gorm.DB // database used by backend

// Initialize input command line flags
func init() {
	flag.StringVar(&DB_HOST, "dbhost", "/var/run/postgresql", "Host of the PostgreSQL database (default is /var/run/postgresql for localhost DB on Ubuntu systems)")
	flag.StringVar(&DB_NAME, "dbname", "villasdb", "Name of the database to use (default is villasdb)")
	flag.StringVar(&DB_USER, "dbuser", "villas", "Username of database connection (default is villas)")
	flag.StringVar(&DB_PASS, "dbpass", "villas", "Password of database connection (default is villas)")
	flag.StringVar(&DB_SSLMODE, "dbsslmode", "disable", "SSL mode of DB (default is disable)") // TODO: change default for production
	flag.BoolVar(&WITH_AMQP, "amqp", false, "If AMQP client for simulators shall be enabled, set this option to true (default is false)")
	flag.Parse()
	fmt.Println("DB_HOST has value ", DB_HOST)
	fmt.Println("DB_USER has value ", DB_USER)
	fmt.Println("DB_NAME has value ", DB_NAME)
	fmt.Println("DB_SSLMODE has value ", DB_SSLMODE)
	fmt.Println("WITH_AMQP has value ", WITH_AMQP)
}

// Initialize connection to the database
func InitDB(dbname string, isTest bool) *gorm.DB {
	dbinfo := fmt.Sprintf("host=%s sslmode=%s user=%s password=%s dbname=%s",
		DB_HOST, DB_SSLMODE, DB_USER, DB_PASS, dbname)
	db, err := gorm.Open("postgres", dbinfo)
	if err != nil {
		log.Fatal(err)
	}
	DBpool = db

	if isTest {
		// drop tables for testing case
		DropTables(db)
	}

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
