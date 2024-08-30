/**
* This file is part of VILLASweb-backend-go
*
* This program is free software: you can redistribute it and/or modify
* it under the terms of the GNU General Public License as published by
* the Free Software Foundation, either version 3 of the License, or
* any later version.
*
* This program is distributed in the hope that it will be useful,
* but WITHOUT ANY WARRANTY; without even the implied warranty of
* MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
* GNU General Public License for more details.
*
* You should have received a copy of the GNU General Public License
* along with this program.  If not, see <http://www.gnu.org/licenses/>.
*********************************************************************************/

package database

import (
	"fmt"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/zpatrick/go-config"
)

var DBpool *gorm.DB // database used by backend

// InitDB Initialize connection to the database
func InitDB(cfg *config.Config, clear bool) error {
	name, err := cfg.String("db.name")
	if err != nil {
		return err
	}

	host, err := cfg.String("db.host")
	if err != nil {
		return err
	}

	port, err := cfg.IntOr("db.port", -1)
	if err != nil {
		return err
	}

	user, err := cfg.StringOr("db.user", "")
	if err != nil {
		return err
	}

	pass := ""
	if user != "" {
		pass, err = cfg.String("db.pass")
		if err != nil {
			return err
		}
	}

	sslmode, err := cfg.String("db.ssl")
	if err != nil {
		return err
	}

	dbinfo := fmt.Sprintf("host=%s sslmode=%s dbname=%s", host, sslmode, name)

	if user != "" && pass != "" {
		dbinfo += fmt.Sprintf(" user=%s password=%s", user, pass)
	}

	if port > 0 {
		dbinfo += fmt.Sprintf(" port=%d", port)
	}

	db, err := gorm.Open("postgres", dbinfo)
	if err != nil {
		return err
	}

	DBpool = db

	// drop tables if parameter set
	if clear {
		DropTables()
		log.Println("Database tables dropped")
	}

	MigrateModels()
	log.Println("Database connection established")

	return nil
}

// GetDB Connection pool to already opened DB
func GetDB() *gorm.DB {
	return DBpool
}

// DropTables drops all the tables of the database (use with care!)
func DropTables() {
	DBpool.DropTableIfExists(&InfrastructureComponent{})
	DBpool.DropTableIfExists(&Signal{})
	DBpool.DropTableIfExists(&ComponentConfiguration{})
	DBpool.DropTableIfExists(&File{})
	DBpool.DropTableIfExists(&Scenario{})
	DBpool.DropTableIfExists(&User{})
	DBpool.DropTableIfExists(&UserGroup{})
	DBpool.DropTableIfExists(&ScenarioMapping{})
	DBpool.DropTableIfExists(&Dashboard{})
	DBpool.DropTableIfExists(&Widget{})
	DBpool.DropTableIfExists(&Result{})
	// The following statement deletes the many to many relationship between users and scenarios
	DBpool.DropTableIfExists("user_scenarios")
}

// MigrateModels AutoMigrate the models
func MigrateModels() {
	DBpool.AutoMigrate(&InfrastructureComponent{})
	DBpool.AutoMigrate(&Signal{})
	DBpool.AutoMigrate(&ComponentConfiguration{})
	DBpool.AutoMigrate(&File{})
	DBpool.AutoMigrate(&Scenario{})
	DBpool.AutoMigrate(&User{})
	DBpool.AutoMigrate(&UserGroup{})
	DBpool.AutoMigrate(&ScenarioMapping{})
	DBpool.AutoMigrate(&Dashboard{})
	DBpool.AutoMigrate(&Widget{})
	DBpool.AutoMigrate(&Result{})
}
