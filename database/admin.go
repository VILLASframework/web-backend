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
		user := User{
			Username: adminName,
			Password: string(pwEnc),
			Role:     "Admin",
			Mail:     mail,
			Active:   true,
		}

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
