package config

import (
	"flag"
	"log"
	"os"

	"github.com/zpatrick/go-config"
)

// Global configuration
var Config *config.Config

func InitConfig() *config.Config {
	var (
		dbHost     = flag.String("dbhost", "/var/run/postgresql", "Host of the PostgreSQL database (default is /var/run/postgresql for localhost DB on Ubuntu systems)")
		dbName     = flag.String("dbname", "villas", "Name of the database to use (default is villasdb)")
		dbUser     = flag.String("dbuser", "", "Username of database connection (default is <empty>)")
		dbPass     = flag.String("dbpass", "", "Password of database connection (default is <empty>)")
		dbInit     = flag.Bool("dbinit", false, "Initialize database with test data (default is off)")
		dbSSLMode  = flag.String("dbsslmode", "disable", "SSL mode of DB (default is disable)") // TODO: change default for production
		amqpURL    = flag.String("amqp", "", "If set, use this url to connect to an AMQP broker (default is disabled)")
		configFile = flag.String("configFile", "", "Path to YAML configuration file")
		mode       = flag.String("mode", "release", "Select debug/release/test mode (default is release)")
		baseHost   = flag.String("base-host", "localhost:4000", "The host:port at which the backend is hosted (default: localhost:4000)")
		basePath   = flag.String("base-path", "/api/v2", "The path at which the API routes are located (default /api/v2)")
	)
	flag.Parse()

	var dbInitStr string
	if *dbInit {
		dbInitStr = "true"
	} else {
		dbInitStr = "false"
	}

	static := map[string]string{
		"db.host":   *dbHost,
		"db.name":   *dbName,
		"db.user":   *dbUser,
		"db.pass":   *dbPass,
		"db.init":   dbInitStr,
		"db.ssl":    *dbSSLMode,
		"amqp.url":  *amqpURL,
		"mode":      *mode,
		"base.host": *baseHost,
		"base.path": *basePath,
	}

	mappings := map[string]string{
		"DB_HOST":   "db.host",
		"DB_NAME":   "db.name",
		"DB_USER":   "db.user",
		"DB_PASS":   "db.pass",
		"DB_SSLMOE": "db.ssl",
		"DB_INIT":   "db.init",
		"AMQP_URL":  "amqp.url",
		"BASE_HOST": "base.host",
		"BASE_PATH": "base.path",
	}

	defaults := config.NewStatic(static)
	env := config.NewEnvironment(mappings)

	var c *config.Config
	if _, err := os.Stat(*configFile); os.IsExist(err) {
		yamlFile := config.NewYAMLFile(*configFile)
		c = config.NewConfig([]config.Provider{defaults, yamlFile, env})
	} else {
		c = config.NewConfig([]config.Provider{defaults, env})
	}

	err := c.Load()
	if err != nil {
		log.Fatal("failed to parse config")
	}

	settings, _ := c.Settings()

	log.Print("All settings:")
	for key, val := range settings {
		log.Printf("   %s = %s \n", key, val)
	}

	// Save pointer to global variable
	Config = c

	return c
}
