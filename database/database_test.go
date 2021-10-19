/** Database package, testing.
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
	"os"
	"testing"

	"github.com/zpatrick/go-config"

	"git.rwth-aachen.de/acs/public/villas/web-backend-go/configuration"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {

	defer DBpool.Close()
	os.Exit(m.Run())
}

func TestInitDB(t *testing.T) {

	err := configuration.InitConfig()
	assert.NoError(t, err)

	static := map[string]string{}
	mappings := map[string]string{}

	defaults := config.NewStatic(static)
	env := config.NewEnvironment(mappings)

	ownconfig := config.NewConfig([]config.Provider{defaults, env})

	err = InitDB(ownconfig, true)
	assert.Error(t, err)
	dbname, err := configuration.GlobalConfig.String("db.name")
	assert.NoError(t, err)
	static["db.name"] = dbname
	ownconfig = config.NewConfig([]config.Provider{defaults, env})
	err = InitDB(ownconfig, true)
	assert.Error(t, err)

	dbhost, err := configuration.GlobalConfig.String("db.host")
	assert.NoError(t, err)
	static["db.host"] = dbhost
	ownconfig = config.NewConfig([]config.Provider{defaults, env})
	err = InitDB(ownconfig, true)
	assert.Error(t, err)

	dbuser, err := configuration.GlobalConfig.String("db.user")
	assert.NoError(t, err)
	static["db.user"] = dbuser
	ownconfig = config.NewConfig([]config.Provider{defaults, env})
	err = InitDB(ownconfig, true)
	assert.Error(t, err)

	dbpass, err := configuration.GlobalConfig.String("db.pass")
	assert.NoError(t, err)
	static["db.pass"] = dbpass
	ownconfig = config.NewConfig([]config.Provider{defaults, env})
	err = InitDB(ownconfig, true)
	assert.Error(t, err)

	dbssl, err := configuration.GlobalConfig.String("db.ssl")
	assert.NoError(t, err)
	static["db.ssl"] = dbssl
	ownconfig = config.NewConfig([]config.Provider{defaults, env})
	err = InitDB(ownconfig, true)
	assert.NoError(t, err)

	// Verify that you can connect to the database
	db := GetDB()
	err = db.DB().Ping()
	assert.NoError(t, err)
}

func TestUserAssociations(t *testing.T) {

	DropTables()
	MigrateModels()

	// create copies of global test data
	scenarioA := Scenario{}
	scenarioB := Scenario{}
	userA := User{Username: "abcdef"}
	userB := User{Username: "ghijkl"}

	// add users to DB
	assert.NoError(t, DBpool.Create(&userA).Error) // Normal User
	assert.NoError(t, DBpool.Create(&userB).Error) // Normal User

	// add two scenarios to DB
	assert.NoError(t, DBpool.Create(&scenarioA).Error)
	assert.NoError(t, DBpool.Create(&scenarioB).Error)

	// add many-to-many associations between users and scenarios
	// User HM Scenarios, Scenario HM Users (Many-to-Many)
	assert.NoError(t, DBpool.Model(&userA).Association("Scenarios").Append(&scenarioA).Error)
	assert.NoError(t, DBpool.Model(&userA).Association("Scenarios").Append(&scenarioB).Error)
	assert.NoError(t, DBpool.Model(&userB).Association("Scenarios").Append(&scenarioA).Error)

	var usr1 User
	assert.NoError(t, DBpool.Find(&usr1, "ID = ?", 1).Error, "Find User with ID=1")

	// Get scenarios of usr1
	var scenarios []Scenario
	assert.NoError(t, DBpool.Model(&usr1).Related(&scenarios, "Scenarios").Error)
	if len(scenarios) != 2 {
		assert.Fail(t, "User Associations",
			"Expected to have %v Scenarios. Has %v.", 2, len(scenarios))
	}
}

func TestScenarioAssociations(t *testing.T) {

	DropTables()
	MigrateModels()

	// create copies of global test data
	scenarioA := Scenario{}
	scenarioB := Scenario{}
	userA := User{Username: "abcdef"}
	userB := User{Username: "ghijkl"}
	configA := ComponentConfiguration{}
	configB := ComponentConfiguration{}
	dashboardA := Dashboard{}
	dashboardB := Dashboard{}
	fileA := File{}
	fileB := File{}
	resultA := Result{}
	resultB := Result{}

	// add scenarios to DB
	assert.NoError(t, DBpool.Create(&scenarioA).Error)
	assert.NoError(t, DBpool.Create(&scenarioB).Error)

	// add users to DB
	assert.NoError(t, DBpool.Create(&userA).Error) // Normal User
	assert.NoError(t, DBpool.Create(&userB).Error) // Normal User

	// add component configurations to DB
	assert.NoError(t, DBpool.Create(&configA).Error)
	assert.NoError(t, DBpool.Create(&configB).Error)

	// add dashboards to DB
	assert.NoError(t, DBpool.Create(&dashboardA).Error)
	assert.NoError(t, DBpool.Create(&dashboardB).Error)

	// add files to DB
	assert.NoError(t, DBpool.Create(&fileA).Error)
	assert.NoError(t, DBpool.Create(&fileB).Error)

	// add results to DB
	assert.NoError(t, DBpool.Create(&resultA).Error)
	assert.NoError(t, DBpool.Create(&resultB).Error)

	// add many-to-many associations between users and scenarios
	// User HM Scenarios, Scenario HM Users (Many-to-Many)
	assert.NoError(t, DBpool.Model(&scenarioA).Association("Users").Append(&userA).Error)
	assert.NoError(t, DBpool.Model(&scenarioA).Association("Users").Append(&userB).Error)
	assert.NoError(t, DBpool.Model(&scenarioB).Association("Users").Append(&userA).Error)
	assert.NoError(t, DBpool.Model(&scenarioB).Association("Users").Append(&userB).Error)

	// add Component Configuration has many files associations
	assert.NoError(t, DBpool.Model(&scenarioA).Association("Files").Append(&fileA).Error)
	assert.NoError(t, DBpool.Model(&scenarioA).Association("Files").Append(&fileB).Error)

	// add scenario has many component configurations associations
	assert.NoError(t, DBpool.Model(&scenarioA).Association("ComponentConfigurations").Append(&configA).Error)
	assert.NoError(t, DBpool.Model(&scenarioA).Association("ComponentConfigurations").Append(&configB).Error)

	// Scenario HM Dashboards
	assert.NoError(t, DBpool.Model(&scenarioA).Association("Dashboards").Append(&dashboardA).Error)
	assert.NoError(t, DBpool.Model(&scenarioA).Association("Dashboards").Append(&dashboardB).Error)

	// Scenario HM Results
	assert.NoError(t, DBpool.Model(&scenarioA).Association("Results").Append(&resultA).Error)
	assert.NoError(t, DBpool.Model(&scenarioA).Association("Results").Append(&resultB).Error)

	var scenario1 Scenario
	assert.NoError(t, DBpool.Find(&scenario1, 1).Error, "Find Scenario with ID=1")

	// Get users of scenario1
	var users []User
	assert.NoError(t, DBpool.Model(&scenario1).Association("Users").Find(&users).Error)
	if len(users) != 2 {
		assert.Fail(t, "Scenario Associations",
			"Expected to have %v Users. Has %v.", 2, len(users))
	}

	// Get component configurations of scenario1
	var configs []ComponentConfiguration
	assert.NoError(t, DBpool.Model(&scenario1).Related(&configs, "ComponentConfigurations").Error)
	if len(configs) != 2 {
		assert.Fail(t, "Scenario Associations",
			"Expected to have %v component configs. Has %v.", 2, len(configs))
	}

	// Get dashboards of scenario1
	var dashboards []Dashboard
	assert.NoError(t, DBpool.Model(&scenario1).Related(&dashboards, "Dashboards").Error)
	if len(dashboards) != 2 {
		assert.Fail(t, "Scenario Associations",
			"Expected to have %v Dashboards. Has %v.", 2, len(dashboards))
	}

	// Get files of scenario1
	var files []File
	assert.NoError(t, DBpool.Model(&scenario1).Related(&files, "Files").Error)
	if len(files) != 2 {
		assert.Fail(t, "Scenario Associations",
			"Expected to have %v Files. Has %v.", 2, len(files))
	}

	// Get results of scenario1
	var results []File
	assert.NoError(t, DBpool.Model(&scenario1).Related(&results, "Results").Error)
	if len(files) != 2 {
		assert.Fail(t, "Scenario Associations",
			"Expected to have %v Results. Has %v.", 2, len(results))
	}
}

func TestICAssociations(t *testing.T) {

	DropTables()
	MigrateModels()

	// create copies of global test data
	icA := InfrastructureComponent{}
	icB := InfrastructureComponent{}
	configA := ComponentConfiguration{}
	configB := ComponentConfiguration{}

	// add ICs to DB
	assert.NoError(t, DBpool.Create(&icA).Error)
	assert.NoError(t, DBpool.Create(&icB).Error)

	// add component configurations to DB
	assert.NoError(t, DBpool.Create(&configA).Error)
	assert.NoError(t, DBpool.Create(&configB).Error)

	// add IC has many component configurations association to DB
	assert.NoError(t, DBpool.Model(&icA).Association("ComponentConfigurations").Append(&configA).Error)
	assert.NoError(t, DBpool.Model(&icA).Association("ComponentConfigurations").Append(&configB).Error)

	var ic1 InfrastructureComponent
	assert.NoError(t, DBpool.Find(&ic1, 1).Error, "Find InfrastructureComponent with ID=1")

	// Get Component Configurations of ic1
	var configs []ComponentConfiguration
	assert.NoError(t, DBpool.Model(&ic1).Association("ComponentConfigurations").Find(&configs).Error)
	if len(configs) != 2 {
		assert.Fail(t, "InfrastructureComponent Associations",
			"Expected to have %v Component Configurations. Has %v.", 2, len(configs))
	}
}

func TestComponentConfigurationAssociations(t *testing.T) {

	DropTables()
	MigrateModels()

	// create copies of global test data
	configA := ComponentConfiguration{}
	configB := ComponentConfiguration{}
	outSignalA := Signal{Direction: "out"}
	outSignalB := Signal{Direction: "out"}
	inSignalA := Signal{Direction: "in"}
	inSignalB := Signal{Direction: "in"}
	icA := InfrastructureComponent{}
	icB := InfrastructureComponent{}

	// add Component Configurations to DB
	assert.NoError(t, DBpool.Create(&configA).Error)
	assert.NoError(t, DBpool.Create(&configB).Error)

	// add signals to DB
	assert.NoError(t, DBpool.Create(&outSignalA).Error)
	assert.NoError(t, DBpool.Create(&outSignalB).Error)
	assert.NoError(t, DBpool.Create(&inSignalA).Error)
	assert.NoError(t, DBpool.Create(&inSignalB).Error)

	// add ICs to DB
	assert.NoError(t, DBpool.Create(&icA).Error)
	assert.NoError(t, DBpool.Create(&icB).Error)

	// add Component Configuration has many signals associations
	assert.NoError(t, DBpool.Model(&configA).Association("InputMapping").Append(&inSignalA).Error)
	assert.NoError(t, DBpool.Model(&configA).Association("InputMapping").Append(&inSignalB).Error)
	assert.NoError(t, DBpool.Model(&configA).Association("OutputMapping").Append(&outSignalA).Error)
	assert.NoError(t, DBpool.Model(&configA).Association("OutputMapping").Append(&outSignalB).Error)

	// associate Component Configurations with IC
	assert.NoError(t, DBpool.Model(&icA).Association("ComponentConfigurations").Append(&configA).Error)
	assert.NoError(t, DBpool.Model(&icA).Association("ComponentConfigurations").Append(&configB).Error)

	var config1 ComponentConfiguration
	assert.NoError(t, DBpool.Find(&config1, 1).Error, "Find ComponentConfiguration with ID=1")

	// Check IC ID
	if config1.ICID != 1 {
		assert.Fail(t, "Component Configurations expected to have Infrastructure Component ID 1, but is %v", config1.ICID)
	}

	// Get OutputMapping signals of config1
	var signals []Signal
	assert.NoError(t, DBpool.Model(&config1).Where("Direction = ?", "out").Related(&signals, "OutputMapping").Error)
	if len(signals) != 2 {
		assert.Fail(t, "ComponentConfiguration Associations",
			"Expected to have %v Output Signals. Has %v.", 2, len(signals))
	}

}

func TestDashboardAssociations(t *testing.T) {

	DropTables()
	MigrateModels()

	// create copies of global test data
	dashboardA := Dashboard{}
	dashboardB := Dashboard{}
	widgetA := Widget{}
	widgetB := Widget{}

	// add dashboards to DB
	assert.NoError(t, DBpool.Create(&dashboardA).Error)
	assert.NoError(t, DBpool.Create(&dashboardB).Error)

	// add widgets to DB
	assert.NoError(t, DBpool.Create(&widgetA).Error)
	assert.NoError(t, DBpool.Create(&widgetB).Error)

	// add dashboard has many widgets associations to DB
	assert.NoError(t, DBpool.Model(&dashboardA).Association("Widgets").Append(&widgetA).Error)
	assert.NoError(t, DBpool.Model(&dashboardA).Association("Widgets").Append(&widgetB).Error)

	var dashboard1 Dashboard
	assert.NoError(t, DBpool.Find(&dashboard1, 1).Error, "Find Dashboard with ID=1")

	//Get widgets of dashboard1
	var widgets []Widget
	assert.NoError(t, DBpool.Model(&dashboard1).Related(&widgets, "Widgets").Error)
	if len(widgets) != 2 {
		assert.Fail(t, "Dashboard Associations",
			"Expected to have %v Widget. Has %v.", 2, len(widgets))
	}
}

func TestWidgetAssociations(t *testing.T) {

	DropTables()
	MigrateModels()

	// create copies of global test data
	widgetA := Widget{}
	widgetB := Widget{}

	// add widgets to DB
	assert.NoError(t, DBpool.Create(&widgetA).Error)
	assert.NoError(t, DBpool.Create(&widgetB).Error)

	var widget1 Widget
	assert.NoError(t, DBpool.Find(&widget1, 1).Error, "Find Widget with ID=1")
}

func TestFileAssociations(t *testing.T) {

	DropTables()
	MigrateModels()

	// create copies of global test data
	fileA := File{}
	fileB := File{}
	fileC := File{}
	fileD := File{}

	// add files to DB
	assert.NoError(t, DBpool.Create(&fileA).Error)
	assert.NoError(t, DBpool.Create(&fileB).Error)
	assert.NoError(t, DBpool.Create(&fileC).Error)
	assert.NoError(t, DBpool.Create(&fileD).Error)

	var file1 File
	assert.NoError(t, DBpool.Find(&file1, 1).Error, "Find File with ID=1")
}
