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

package configuration

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/zpatrick/go-config"
)

// Global configuration
var GlobalConfig *config.Config = nil

type GroupedScenario struct {
	Scenario  int  `yaml:"scenario"`
	Duplicate bool `default:"false" yaml:"duplicate"`
}

var ScenarioGroupMap = map[string][]GroupedScenario{}

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
		dbClear                  = flag.Bool("db-clear", false, "Set to true if you want to clear all DB tables upon startup. This parameter has to be used with great care, its effects cannot be reverted.")
		amqpHost                 = flag.String("amqp-host", "", "If set, use this as host for AMQP broker (default is disabled)")
		amqpUser                 = flag.String("amqp-user", "", "Username for AMQP broker")
		amqpPass                 = flag.String("amqp-pass", "", "Password for AMQP broker")
		configFile               = flag.String("config", "", "Path to YAML configuration file")
		port                     = flag.String("port", "4000", "Port of the backend (default is 4000)")
		adminUser                = flag.String("admin-user", "", "Initial admin username")
		adminPass                = flag.String("admin-pass", "", "Initial admin password")
		adminMail                = flag.String("admin-mail", "", "Initial admin mail address")
		s3Bucket                 = flag.String("s3-bucket", "", "S3 Bucket for uploading files")
		s3Endpoint               = flag.String("s3-endpoint", "", "Endpoint of S3 API for file uploads")
		s3EndpointPublic         = flag.String("s3-endpoint-public", "", "Public endpoint address of S3 API for file uploads")
		s3Region                 = flag.String("s3-region", "default", "S3 Region for file uploads")
		s3NoSSL                  = flag.Bool("s3-nossl", false, "Use encrypted connections to the S3 API")
		s3PathStyle              = flag.Bool("s3-pathstyle", false, "Use path-style S3 API")
		jwtSecret                = flag.String("jwt-secret", "This should NOT be here!!@33$8&", "The JSON Web Token secret")
		jwtExpiresAfter          = flag.String("jwt-expires-after", "168h" /* 1 week */, "The time after which the JSON Web Token expires")
		authExternal             = flag.Bool("auth-external", false, "Use external authentication via X-Forwarded-User header (e.g. OAuth2 Proxy)")
		authExternalLoginURL     = flag.String("auth-external-login-url", "/oauth2/start", "A URL to initiate external login procedure")
		authExternalProviderName = flag.String("auth-external-provider-name", "JupyterHub", "A name of the external authentication provider")
		authLogoutURL            = flag.String("auth-logout-url", "/oauth2/sign_out?rd=https%3A%2F%2Fjupyter.k8s.eonerc.rwth-aachen.de%2Fhub%2Flogout", "The URL to redirect the user to log out")
		title                    = flag.String("title", "VILLASweb", "Title shown in the frontend")
		subTitle                 = flag.String("sub-title", "", "Sub-title shown in the frontend")
		contactName              = flag.String("contact-name", "Steffen Vogel", "Name of the administrative contact")
		contactMail              = flag.String("contact-mail", "svogel2@eonerc.rwth-aachen.de", "EMail of the administrative contact")
		testDataPath             = flag.String("test-data-path", "", "The path to a test data json file")
		groupsPath               = flag.String("groups-path", "", "The path to a YAML file that maps user groups to scenario IDs")
		apiUpdateInterval        = flag.String("api-update-interval", "10s" /* 10 sec */, "Interval in which API URL is queried for status updates of ICs")
		k8sRancherURL            = flag.String("k8s-rancher-url", "https://rancher.k8s.eonerc.rwth-aachen.de", "URL of Rancher instance that is used to deploy the backend")
		k8sClusterName           = flag.String("k8s-cluster-name", "local", "Name of the Kubernetes cluster where the backend is deployed")
		staleICTime              = flag.String("stale-ic-time", "1h" /* 1 hour */, "Time after which an IC is considered stale")
		webRTCiceUrls            = flag.String("webrtc-ice-urls",
			"stun:stun.l.google.com:19302,villas:villas@stun:stun.0l.de,villas:villas@turn:turn.0l.de?transport=udp,villas:villas@turn:turn.0l.de?transport=tcp",
			"WebRTC ICE URLs (comma-separated list, use username:password@url style for non-anonymous URLs)")
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
		"port":                        *port,
		"admin.user":                  *adminUser,
		"admin.pass":                  *adminPass,
		"admin.mail":                  *adminMail,
		"s3.bucket":                   *s3Bucket,
		"s3.endpoint":                 *s3Endpoint,
		"s3.endpoint-public":          *s3EndpointPublic,
		"s3.region":                   *s3Region,
		"jwt.secret":                  *jwtSecret,
		"jwt.expires-after":           *jwtExpiresAfter,
		"auth.external.login-url":     *authExternalLoginURL,
		"auth.external.provider-name": *authExternalProviderName,
		"auth.logout-url":             *authLogoutURL,
		"title":                       *title,
		"sub-title":                   *subTitle,
		"contact.name":                *contactName,
		"contact.mail":                *contactMail,
		"test.datapath":               *testDataPath,
		"groups.path":                 *groupsPath,
		"config.file":                 *configFile,
		"apiupdateinterval":           *apiUpdateInterval,
		"k8s.rancher-url":             *k8sRancherURL,
		"k8s.cluster-name":            *k8sClusterName,
		"staleictime":                 *staleICTime,
		"webrtc.ice-urls":             *webRTCiceUrls,
	}

	if *dbClear {
		static["db.clear"] = "true"
	} else {
		static["db.clear"] = "false"
	}

	if *s3NoSSL {
		static["s3.nossl"] = "true"
	} else {
		static["s3.nossl"] = "false"
	}

	if *s3PathStyle {
		static["s3.pathstyle"] = "true"
	} else {
		static["s3.pathstyle"] = "false"
	}

	if *authExternal {
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

	settings, _ := GlobalConfig.Settings()

	keys := make([]string, 0, len(settings))
	for k := range settings {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	log.Print("All settings (except for PW and secrets):")
	for _, k := range keys {
		if !strings.Contains(k, "pass") && !strings.Contains(k, "secret") {
			log.Printf("   %s = %s \n", k, settings[k])
		}
	}
	return nil
}

func remove(arr []GroupedScenario, index int) []GroupedScenario {
	arr[index] = arr[len(arr)-1]
	return arr[:len(arr)-1]
}

func ReadGroupsFile(path string) error {

	_, err := os.Stat(path)
	if err != nil {
		return err
	}

	yamlFile, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("error opening yaml file for groups: %v", err)
	}
	log.Println("Successfully opened yaml groups file", path)

	defer yamlFile.Close()

	byteValue, _ := ioutil.ReadAll(yamlFile)

	err = yaml.Unmarshal(byteValue, &ScenarioGroupMap)
	if err != nil {
		return fmt.Errorf("error unmarshalling yaml into ScenarioGroupMap: %v", err)
	}

	for _, group := range ScenarioGroupMap {
		for i, scenario := range group {
			// remove invalid values that might have been introduced by typos
			// (Unmarshal sets default values when it doesn't find a field)
			if scenario.Scenario == 0 {
				log.Println("Removing entry from ScenarioGroupMap, check for typos in the yaml!")
				remove(group, i)
			}
		}
	}

	log.Println("ScenarioGroupMap", ScenarioGroupMap)

	return nil
}
