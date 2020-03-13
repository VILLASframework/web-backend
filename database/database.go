/** Database package.
*
* @author Sonja Happ <sonja.happ@eonerc.rwth-aachen.de>
* @copyright 2014-2019, Institute for Automation of Complex Power Systems, EONERC
* @license GNU General Public License (version 3)
*
* VILLASweb-backend-go
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
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/zpatrick/go-config"
)

var DBpool *gorm.DB // database used by backend

// Initialize connection to the database
func InitDB(cfg *config.Config) (*gorm.DB, error) {
	name, err := cfg.String("db.name")
	if err != nil {
		return nil, err
	}
	host, err := cfg.String("db.host")
	if err != nil {
		return nil, err
	}
	user, err := cfg.String("db.user")
	if err != nil && !strings.Contains(err.Error(), "Required setting 'db.user' not set") {
		return nil, err
	}
	pass := ""
	if user != "" {
		pass, err = cfg.String("db.pass")
		if err != nil {
			return nil, err
		}
	}
	sslmode, err := cfg.String("db.ssl")
	if err != nil {
		return nil, err
	}

	dbinfo := fmt.Sprintf("host=%s sslmode=%s dbname=%s", host, sslmode, name)
	if user != "" && pass != "" {
		dbinfo += fmt.Sprintf(" user=%s password=%s", user, pass)
	}

	db, err := gorm.Open("postgres", dbinfo)
	if err != nil {
		return nil, err
	}

	DBpool = db
	MigrateModels(db)
	log.Println("Database connection established")

	return db, nil
}

// Connection pool to already opened DB
func GetDB() *gorm.DB {
	return DBpool
}

// Drop all the tables of the database
// TODO: Remove that function from the codebase and substitute the body
// to the Dummy*() where it is called
func DropTables(db *gorm.DB) {
	db.DropTableIfExists(&InfrastructureComponent{})
	db.DropTableIfExists(&Signal{})
	db.DropTableIfExists(&ComponentConfiguration{})
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
	db.AutoMigrate(&InfrastructureComponent{})
	db.AutoMigrate(&Signal{})
	db.AutoMigrate(&ComponentConfiguration{})
	db.AutoMigrate(&File{})
	db.AutoMigrate(&Scenario{})
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Dashboard{})
	db.AutoMigrate(&Widget{})
}
