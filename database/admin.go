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
	"math/rand"
	"strings"
	"time"

	"github.com/zpatrick/go-config"
	"golang.org/x/crypto/bcrypt"
)

// AddAdminUser adds a default admin user to the DB
func AddAdminUser(cfg *config.Config) (string, error) {
	DBpool.AutoMigrate(User{})

	updatedPW := false
	generatedPW := false

	adminName, _ := cfg.StringOr("admin.user", "admin")
	adminMail, _ := cfg.StringOr("admin.mail", "admin@example.com")
	adminPW, err := cfg.String("admin.pass")
	if err == nil && adminPW != "" {
		updatedPW = true
	} else if err != nil || adminPW == "" {
		adminPW = generatePassword(16)
		generatedPW = true
	}

	adminPWEnc, err := bcrypt.GenerateFromPassword([]byte(adminPW), 10)
	if err != nil {
		return "", err
	}

	// Check if admin user exists in DB
	var users []User
	err = DBpool.Where("Username = ?", adminName).Find(&users).Error
	if err != nil {
		return "", err
	}

	if len(users) == 0 {
		fmt.Println("No admin user found in DB, adding default admin user.")
		if generatedPW {
			fmt.Printf("  Generated admin password: %s for admin user %s\n", adminPW, adminName)
		}

		user := User{
			Username: adminName,
			Password: string(adminPWEnc),
			Role:     "Admin",
			Mail:     adminMail,
			Active:   true,
		}

		// add admin user to DB
		err = DBpool.Create(&user).Error
		if err != nil {
			return "", err
		}
	} else if updatedPW {
		fmt.Println("Found existing admin user in DB, updating user from CLI parameters.")

		user := users[0]

		user.Password = string(adminPWEnc)
		user.Role = "Admin"
		user.Mail = adminMail
		user.Active = true

		err = DBpool.Model(user).Update(&user).Error
		if err != nil {
			return "", err
		}
	}

	return adminPW, err
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
