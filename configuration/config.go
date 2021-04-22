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
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/zpatrick/go-config"
)

// Global configuration
var GlobalConfig *config.Config = nil

var ScenarioGroupMap = map[string][]int{}

func InitConfig() error {
	if GlobalConfig != nil {
		return nil
	}

	var (
		dbHost                   = flag.String("db-host", "/var/run/postgresql", "Host of the PostgreSQL database (default is /var/run/postgresql for localhost DB on Ubuntu systems)")
		dbName                   = flag.String("db-name", "villasdb", "Name of the database to use (default is villasdb)")
		dbUser                   = flag.String("db-user", "", "Username of database connection (default is <empty>)")
		dbPass                   = flag.String("db-pass", "", "Password of database connection (default is <empty>)")
		dbSSLMode                = flag.String("db-ssl-mode", "disable", "SSL mode of DB (default is disable)") // TODO: change default for production
		amqpHost                 = flag.String("amqp-host", "", "If set, use this as host for AMQP broker (default is disabled)")
		amqpUser                 = flag.String("amqp-user", "", "Username for AMQP broker")
		amqpPass                 = flag.String("amqp-pass", "", "Password for AMQP broker")
		configFile               = flag.String("config", "", "Path to YAML configuration file")
		mode                     = flag.String("mode", "release", "Select debug/release/test mode (default is release)")
		port                     = flag.String("port", "4000", "Port of the backend (default is 4000)")
		adminUser                = flag.String("admin-user", "", "Initial admin username")
		adminPass                = flag.String("admin-pass", "", "Initial admin password")
		adminMail                = flag.String("admin-mail", "", "Initial admin mail address")
		s3Bucket                 = flag.String("s3-bucket", "", "S3 Bucket for uploading files")
		s3Endpoint               = flag.String("s3-endpoint", "", "Endpoint of S3 API for file uploads")
		s3Region                 = flag.String("s3-region", "default", "S3 Region for file uploads")
		s3NoSSL                  = flag.Bool("s3-nossl", false, "Use encrypted connections to the S3 API")
		s3PathStyle              = flag.Bool("s3-pathstyle", false, "Use path-style S3 API")
		jwtSecret                = flag.String("jwt-secret", "This should NOT be here!!@33$8&", "The JSON Web Token secret")
		jwtExpiresAfter          = flag.String("jwt-expires-after", "168h" /* 1 week */, "The time after which the JSON Web Token expires")
		authExternal             = flag.Bool("auth-external", false, "Use external authentication via X-Forwarded-User header (e.g. OAuth2 Proxy)")
		authExternalAuthorizeURL = flag.String("authexternal-authorize-url", "/oauth2/start", "A URL to initiate external login procedure")
		authExternalProviderName = flag.String("auth-external-provider-name", "JupyterHub", "A name of the external authentication provider")
		authLogoutURL            = flag.String("auth-logout-url", "/oauth2/sign_out?rd=https%3A%2F%2Fjupyter.k8s.eonerc.rwth-aachen.de%2Fhub%2Flogout", "The URL to redirect the user to log out")
		title                    = flag.String("title", "VILLASweb", "Title shown in the frontend")
		subTitle                 = flag.String("sub-title", "", "Sub-title shown in the frontend")
		contactName              = flag.String("contact-name", "Steffen Vogel", "Name of the administrative contact")
		contactMail              = flag.String("contact-mail", "svogel2@eonerc.rwth-aachen.de", "EMail of the administrative contact")
		testDataPath             = flag.String("test-data-path", "database/testdata.json", "The path to the test data json file (used in test mode)")
		groupsPath               = flag.String("groups-path", "configuration/groups.json", "The path to the JSON file that maps user groups to scenario IDs")
	)
	flag.Parse()

	static := map[string]string{
		"db.host":                     *dbHost,
		"db.name":                     *dbName,
		"db.user":                     *dbUser,
		"db.pass":                     *dbPass,
		"db.ssl":                      *dbSSLMode,
		"amqp.host":                   *amqpHost,
		"amqp.user":                   *amqpUser,
		"amqp.pass":                   *amqpPass,
		"mode":                        *mode,
		"port":                        *port,
		"admin.user":                  *adminUser,
		"admin.pass":                  *adminPass,
		"admin.mail":                  *adminMail,
		"s3.bucket":                   *s3Bucket,
		"s3.endpoint":                 *s3Endpoint,
		"s3.region":                   *s3Region,
		"jwt.secret":                  *jwtSecret,
		"jwt.expires-after":           *jwtExpiresAfter,
		"auth.external.authorize-url": *authExternalAuthorizeURL,
		"auth.external.provider-name": *authExternalProviderName,
		"auth.logout-url":             *authLogoutURL,
		"title":                       *title,
		"sub-title":                   *subTitle,
		"contact.name":                *contactName,
		"contact.mail":                *contactMail,
		"test.datapath":               *testDataPath,
		"groups.path":                 *groupsPath,
	}

	if *s3NoSSL == true {
		static["s3.nossl"] = "true"
	} else {
		static["s3.nossl"] = "false"
	}

	if *s3PathStyle == true {
		static["s3.pathstyle"] = "true"
	} else {
		static["s3.pathstyle"] = "false"
	}

	if *authExternal == true {
		static["auth.external.enabled"] = "true"
	} else {
		static["auth.external.enabled"] = "false"
	}

	mappings := map[string]string{}
	for name := range static {
		envName := strings.ReplaceAll(name, ".", "_")
		envName = strings.ReplaceAll(envName, "-", "_")
		envName = strings.ToUpper(envName)

		mappings[envName] = name
	}

	defaults := config.NewStatic(static)
	env := config.NewEnvironment(mappings)

	if _, err := os.Stat(*configFile); os.IsNotExist(err) {
		GlobalConfig = config.NewConfig([]config.Provider{defaults, env})
	} else {
		yamlFile := config.NewYAMLFile(*configFile)
		GlobalConfig = config.NewConfig([]config.Provider{defaults, yamlFile, env})
	}

	err := GlobalConfig.Load()
	if err != nil {
		log.Fatal("failed to parse config")
		return err
	}

	m, err := GlobalConfig.String("mode")
	if err != nil {
		return err
	}

	if m != "test" {
		settings, _ := GlobalConfig.Settings()
		log.Print("All settings:")
		for key, val := range settings {
			// TODO password settings should be excluded!
			log.Printf("   %s = %s \n", key, val)
		}
	}

	return nil
}

func ReadGroupsFile(path string) error {

	_, err := os.Stat(path)

	if err == nil {

		jsonFile, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("error opening json file for groups: %v", err)
		}
		log.Println("Successfully opened json groups file", path)

		defer jsonFile.Close()

		byteValue, _ := ioutil.ReadAll(jsonFile)

		err = json.Unmarshal(byteValue, &ScenarioGroupMap)
		if err != nil {
			return fmt.Errorf("error unmarshalling json into ScenarioGroupMap: %v", err)
		}

		log.Println("ScenarioGroupMap", ScenarioGroupMap)

		return nil
	} else if os.IsNotExist(err) {
		log.Println("File does not exist, no goups/scenarios mapping created:", path)
		return nil
	} else {
		log.Println("Something is wrong with this file path:", path)
		return nil
	}
}
