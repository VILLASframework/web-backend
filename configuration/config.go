/** Configuration package.
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
package configuration

import (
	"flag"
	"log"
	"os"

	"github.com/zpatrick/go-config"
)

// Global configuration
var GolbalConfig *config.Config = nil

func InitConfig() error {
	if GolbalConfig != nil {
		return nil
	}

	var (
		dbHost     = flag.String("db-host", "/var/run/postgresql", "Host of the PostgreSQL database (default is /var/run/postgresql for localhost DB on Ubuntu systems)")
		dbName     = flag.String("db-name", "villasdb", "Name of the database to use (default is villasdb)")
		dbUser     = flag.String("db-user", "", "Username of database connection (default is <empty>)")
		dbPass     = flag.String("db-pass", "", "Password of database connection (default is <empty>)")
		dbSSLMode  = flag.String("db-ssl-mode", "disable", "SSL mode of DB (default is disable)") // TODO: change default for production
		amqpHost   = flag.String("amqp-host", "", "If set, use this as host for AMQP broker (default is disabled)")
		amqpUser   = flag.String("amqp-user", "", "Username for AMQP broker")
		amqpPass   = flag.String("amqp-pass", "", "Password for AMQP broker")
		configFile = flag.String("config", "", "Path to YAML configuration file")
		mode       = flag.String("mode", "release", "Select debug/release/test mode (default is release)")
		port       = flag.String("port", "4000", "Port of the backend (default is 4000)")
		baseHost   = flag.String("base-host", "localhost:4000", "The host at which the backend is hosted (default: localhost)")
		basePath   = flag.String("base-path", "/api/v2", "The path at which the API routes are located (default /api/v2)")
		adminUser  = flag.String("admin-user", "", "Initial admin username")
		adminPass  = flag.String("admin-pass", "", "Initial admin password")
		adminMail  = flag.String("admin-mail", "", "Initial admin mail address")
	)
	flag.Parse()

	static := map[string]string{
		"db.host":    *dbHost,
		"db.name":    *dbName,
		"db.user":    *dbUser,
		"db.pass":    *dbPass,
		"db.ssl":     *dbSSLMode,
		"amqp.host":  *amqpHost,
		"amqp.user":  *amqpUser,
		"amqp.pass":  *amqpPass,
		"mode":       *mode,
		"port":       *port,
		"base.host":  *baseHost,
		"base.path":  *basePath,
		"admin.user": *adminUser,
		"admin.pass": *adminPass,
		"admin.mail": *adminMail,
	}

	mappings := map[string]string{
		"DB_HOST":    "db.host",
		"DB_NAME":    "db.name",
		"DB_USER":    "db.user",
		"DB_PASS":    "db.pass",
		"DB_SSLMODE": "db.ssl",
		"AMQP_HOST":  "amqp.host",
		"AMQP_USER":  "amqp.user",
		"AMQP_PASS":  "amqp.pass",
		"BASE_HOST":  "base.host",
		"BASE_PATH":  "base.path",
		"MODE":       "mode",
		"PORT":       "port",
		"ADMIN_USER": "admin.user",
		"ADMIN_PASS": "admin.pass",
		"ADMIN_MAIL": "admin.mail",
	}

	defaults := config.NewStatic(static)
	env := config.NewEnvironment(mappings)

	if _, err := os.Stat(*configFile); os.IsExist(err) {
		yamlFile := config.NewYAMLFile(*configFile)
		GolbalConfig = config.NewConfig([]config.Provider{defaults, yamlFile, env})
	} else {
		GolbalConfig = config.NewConfig([]config.Provider{defaults, env})
	}

	err := GolbalConfig.Load()
	if err != nil {
		log.Fatal("failed to parse config")
		return err
	}

	m, err := GolbalConfig.String("mode")
	if err != nil {
		return err
	}

	if m != "test" {
		settings, _ := GolbalConfig.Settings()
		log.Print("All settings:")
		for key, val := range settings {
			// TODO password settings should be excluded!
			log.Printf("   %s = %s \n", key, val)
		}
	}

	return nil
}

func ConfigureBackend() (string, string, string, string, string, string, string, error) {

	err := InitConfig()
	if err != nil {
		log.Printf("Error during initialization of global configuration: %v, aborting.", err.Error())
		return "", "", "", "", "", "", "", err
	}

	mode, err := GolbalConfig.String("mode")
	if err != nil {
		log.Printf("Error reading mode from global configuration: %v, aborting.", err.Error())
		return "", "", "", "", "", "", "", err
	}

	baseHost, err := GolbalConfig.String("base.host")
	if err != nil {
		log.Printf("Error reading base.host from global configuration: %v, aborting.", err.Error())
		return "", "", "", "", "", "", "", err
	}
	basePath, err := GolbalConfig.String("base.path")
	if err != nil {
		log.Printf("Error reading base.path from global configuration: %v, aborting.", err.Error())
		return "", "", "", "", "", "", "", err
	}
	port, err := GolbalConfig.String("port")
	if err != nil {
		log.Printf("Error reading port from global configuration: %v, aborting.", err.Error())
		return "", "", "", "", "", "", "", err
	}

	AMQPhost, _ := GolbalConfig.String("amqp.host")
	AMQPuser, _ := GolbalConfig.String("amqp.user")
	AMQPpass, _ := GolbalConfig.String("amqp.pass")

	return mode, baseHost, basePath, port, AMQPhost, AMQPuser, AMQPpass, nil
}
