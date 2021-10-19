/** Package database
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
	"math/rand"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/zpatrick/go-config"
)

var DBpool *gorm.DB // database used by backend

// InitDB Initialize connection to the database
func InitDB(cfg *config.Config, dbClear string) error {
	name, err := cfg.String("db.name")
	if err != nil {
		return err
	}
	host, err := cfg.String("db.host")
	if err != nil {
		return err
	}
	user, err := cfg.String("db.user")
	if err != nil && !strings.Contains(err.Error(), "Required setting 'db.user' not set") {
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

	db, err := gorm.Open("postgres", dbinfo)
	if err != nil {
		return err
	}

	DBpool = db

	// drop tables if parameter set
	if dbClear == "true" {
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
	DBpool.AutoMigrate(&Dashboard{})
	DBpool.AutoMigrate(&Widget{})
	DBpool.AutoMigrate(&Result{})
}

// DBAddAdminUser adds a default admin user to the DB
func DBAddAdminUser(cfg *config.Config) (string, error) {
	DBpool.AutoMigrate(User{})

	// Check if admin user exists in DB
	var users []User
	DBpool.Where("Role = ?", "Admin").Find(&users)

	adminPW := ""

	if len(users) == 0 {
		fmt.Println("No admin user found in DB, adding default admin user.")

		adminName, err := cfg.String("admin.user")
		if err != nil || adminName == "" {
			adminName = "admin"
		}

		adminPW, err = cfg.String("admin.pass")
		if err != nil || adminPW == "" {
			adminPW = generatePassword(16)
			fmt.Printf("  Generated admin password: %s for admin user %s\n", adminPW, adminName)
		}

		mail, err := cfg.String("admin.mail")
		if err == nil || mail == "" {
			mail = "admin@example.com"
		}

		pwEnc, _ := bcrypt.GenerateFromPassword([]byte(adminPW), 10)

		// create a copy of global test data
		user := User{Username: adminName, Password: string(pwEnc),
			Role: "Admin", Mail: mail, Active: true}

		// add admin user to DB
		err = DBpool.Create(&user).Error
		if err != nil {
			return "", err
		}
	}
	return adminPW, nil
}

func generatePassword(Len int) string {
	rand.Seed(time.Now().UnixNano())
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")

	var b strings.Builder
	for i := 0; i < Len; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}

	return b.String()
}
