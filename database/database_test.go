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
	"fmt"
	"github.com/jinzhu/gorm"
	"log"
	"os"
	"testing"

	"git.rwth-aachen.de/acs/public/villas/web-backend-go/configuration"
	"github.com/stretchr/testify/assert"
)

var db *gorm.DB

func TestMain(m *testing.M) {
	err := configuration.InitConfig()
	if err != nil {
		panic(m)
	}

	db, err = InitDB(configuration.GolbalConfig)
	if err != nil {
		panic(m)
	}

	// Verify that you can connect to the database
	err = db.DB().Ping()
	if err != nil {
		log.Panic("Error: DB must ping to run tests")
	}

	defer db.Close()
	os.Exit(m.Run())
}

func TestUserAssociations(t *testing.T) {

	DropTables(db)
	MigrateModels(db)

	// create copies of global test data
	scenarioA := ScenarioA
	scenarioB := ScenarioB
	user0 := User0
	userA := UserA
	userB := UserB

	// add three users to DB
	assert.NoError(t, db.Create(&user0).Error) // Admin
	assert.NoError(t, db.Create(&userA).Error) // Normal User
	assert.NoError(t, db.Create(&userB).Error) // Normal User

	// add two scenarios to DB
	assert.NoError(t, db.Create(&scenarioA).Error)
	assert.NoError(t, db.Create(&scenarioB).Error)

	// add many-to-many associations between users and scenarios
	// User HM Scenarios, Scenario HM Users (Many-to-Many)
	assert.NoError(t, db.Model(&userA).Association("Scenarios").Append(&scenarioA).Error)
	assert.NoError(t, db.Model(&userA).Association("Scenarios").Append(&scenarioB).Error)
	assert.NoError(t, db.Model(&userB).Association("Scenarios").Append(&scenarioA).Error)
	assert.NoError(t, db.Model(&userB).Association("Scenarios").Append(&scenarioB).Error)

	var usr1 User
	assert.NoError(t, db.Find(&usr1, "ID = ?", 2).Error, fmt.Sprintf("Find User with ID=2"))
	assert.EqualValues(t, "User_A", usr1.Username)

	// Get scenarios of usr1
	var scenarios []Scenario
	assert.NoError(t, db.Model(&usr1).Related(&scenarios, "Scenarios").Error)
	if len(scenarios) != 2 {
		assert.Fail(t, "User Associations",
			"Expected to have %v Scenarios. Has %v.", 2, len(scenarios))
	}
}

func TestScenarioAssociations(t *testing.T) {

	DropTables(db)
	MigrateModels(db)

	// create copies of global test data
	scenarioA := ScenarioA
	scenarioB := ScenarioB
	user0 := User0
	userA := UserA
	userB := UserB
	modelA := SimulationModelA
	modelB := SimulationModelB
	dashboardA := DashboardA
	dashboardB := DashboardB

	// add scenarios to DB
	assert.NoError(t, db.Create(&scenarioA).Error)
	assert.NoError(t, db.Create(&scenarioB).Error)

	// add users to DB
	assert.NoError(t, db.Create(&user0).Error) // Admin
	assert.NoError(t, db.Create(&userA).Error) // Normal User
	assert.NoError(t, db.Create(&userB).Error) // Normal User

	// add simulation models to DB
	assert.NoError(t, db.Create(&modelA).Error)
	assert.NoError(t, db.Create(&modelB).Error)

	// add dashboards to DB
	assert.NoError(t, db.Create(&dashboardA).Error)
	assert.NoError(t, db.Create(&dashboardB).Error)

	// add many-to-many associations between users and scenarios
	// User HM Scenarios, Scenario HM Users (Many-to-Many)
	assert.NoError(t, db.Model(&scenarioA).Association("Users").Append(&userA).Error)
	assert.NoError(t, db.Model(&scenarioA).Association("Users").Append(&userB).Error)
	assert.NoError(t, db.Model(&scenarioB).Association("Users").Append(&userA).Error)
	assert.NoError(t, db.Model(&scenarioB).Association("Users").Append(&userB).Error)

	// add scenario has many simulation models associations
	assert.NoError(t, db.Model(&scenarioA).Association("SimulationModels").Append(&modelA).Error)
	assert.NoError(t, db.Model(&scenarioA).Association("SimulationModels").Append(&modelB).Error)

	// Scenario HM Dashboards
	assert.NoError(t, db.Model(&scenarioA).Association("Dashboards").Append(&dashboardA).Error)
	assert.NoError(t, db.Model(&scenarioA).Association("Dashboards").Append(&dashboardB).Error)

	var scenario1 Scenario
	assert.NoError(t, db.Find(&scenario1, 1).Error, fmt.Sprintf("Find Scenario with ID=1"))
	assert.EqualValues(t, "Scenario_A", scenario1.Name)

	// Get users of scenario1
	var users []User
	assert.NoError(t, db.Model(&scenario1).Association("Users").Find(&users).Error)
	if len(users) != 2 {
		assert.Fail(t, "Scenario Associations",
			"Expected to have %v Users. Has %v.", 2, len(users))
	}

	// Get simulation models of scenario1
	var models []SimulationModel
	assert.NoError(t, db.Model(&scenario1).Related(&models, "SimulationModels").Error)
	if len(models) != 2 {
		assert.Fail(t, "Scenario Associations",
			"Expected to have %v simulation models. Has %v.", 2, len(models))
	}

	// Get dashboards of scenario1
	var dashboards []Dashboard
	assert.NoError(t, db.Model(&scenario1).Related(&dashboards, "Dashboards").Error)
	if len(dashboards) != 2 {
		assert.Fail(t, "Scenario Associations",
			"Expected to have %v Dashboards. Has %v.", 2, len(dashboards))
	}
}

func TestICAssociations(t *testing.T) {

	DropTables(db)
	MigrateModels(db)

	// create copies of global test data
	icA := ICA
	icB := ICB
	modelA := SimulationModelA
	modelB := SimulationModelB

	// add ICs to DB
	assert.NoError(t, db.Create(&icA).Error)
	assert.NoError(t, db.Create(&icB).Error)

	// add simulation models to DB
	assert.NoError(t, db.Create(&modelA).Error)
	assert.NoError(t, db.Create(&modelB).Error)

	// add IC has many simulation models association to DB
	assert.NoError(t, db.Model(&icA).Association("SimulationModels").Append(&modelA).Error)
	assert.NoError(t, db.Model(&icA).Association("SimulationModels").Append(&modelB).Error)

	var ic1 InfrastructureComponent
	assert.NoError(t, db.Find(&ic1, 1).Error, fmt.Sprintf("Find InfrastructureComponent with ID=1"))
	assert.EqualValues(t, "Host_A", ic1.Host)

	// Get simulation models of ic1
	var models []SimulationModel
	assert.NoError(t, db.Model(&ic1).Association("SimulationModels").Find(&models).Error)
	if len(models) != 2 {
		assert.Fail(t, "InfrastructureComponent Associations",
			"Expected to have %v SimulationModels. Has %v.", 2, len(models))
	}
}

func TestSimulationModelAssociations(t *testing.T) {

	DropTables(db)
	MigrateModels(db)

	// create copies of global test data
	modelA := SimulationModelA
	modelB := SimulationModelB
	outSignalA := OutSignalA
	outSignalB := OutSignalB
	inSignalA := InSignalA
	inSignalB := InSignalB
	fileA := FileA
	fileB := FileB
	fileC := FileC
	fileD := FileD
	icA := ICA
	icB := ICB

	// add simulation models to DB
	assert.NoError(t, db.Create(&modelA).Error)
	assert.NoError(t, db.Create(&modelB).Error)

	// add signals to DB
	assert.NoError(t, db.Create(&outSignalA).Error)
	assert.NoError(t, db.Create(&outSignalB).Error)
	assert.NoError(t, db.Create(&inSignalA).Error)
	assert.NoError(t, db.Create(&inSignalB).Error)

	// add files to DB
	assert.NoError(t, db.Create(&fileA).Error)
	assert.NoError(t, db.Create(&fileB).Error)
	assert.NoError(t, db.Create(&fileC).Error)
	assert.NoError(t, db.Create(&fileD).Error)

	// add ICs to DB
	assert.NoError(t, db.Create(&icA).Error)
	assert.NoError(t, db.Create(&icB).Error)

	// add simulation model has many signals associations
	assert.NoError(t, db.Model(&modelA).Association("InputMapping").Append(&inSignalA).Error)
	assert.NoError(t, db.Model(&modelA).Association("InputMapping").Append(&inSignalB).Error)
	assert.NoError(t, db.Model(&modelA).Association("OutputMapping").Append(&outSignalA).Error)
	assert.NoError(t, db.Model(&modelA).Association("OutputMapping").Append(&outSignalB).Error)

	// add simulation model has many files associations
	assert.NoError(t, db.Model(&modelA).Association("Files").Append(&fileC).Error)
	assert.NoError(t, db.Model(&modelA).Association("Files").Append(&fileD).Error)

	// associate simulation models with IC
	assert.NoError(t, db.Model(&icA).Association("SimulationModels").Append(&modelA).Error)
	assert.NoError(t, db.Model(&icA).Association("SimulationModels").Append(&modelB).Error)

	var model1 SimulationModel
	assert.NoError(t, db.Find(&model1, 1).Error, fmt.Sprintf("Find SimulationModel with ID=1"))
	assert.EqualValues(t, "SimulationModel_A", model1.Name)

	// Check IC ID
	if model1.ICID != 1 {
		assert.Fail(t, "Simulation Model expected to have InfrastructureComponent ID 1, but is %v", model1.ICID)
	}

	// Get OutputMapping signals of model1
	var signals []Signal
	assert.NoError(t, db.Model(&model1).Where("Direction = ?", "out").Related(&signals, "OutputMapping").Error)
	if len(signals) != 2 {
		assert.Fail(t, "SimulationModel Associations",
			"Expected to have %v Output Signals. Has %v.", 2, len(signals))
	}

	// Get files of model1
	var files []File
	assert.NoError(t, db.Model(&model1).Related(&files, "Files").Error)
	if len(files) != 2 {
		assert.Fail(t, "SimulationModel Associations",
			"Expected to have %v Files. Has %v.", 2, len(files))
	}
}

func TestDashboardAssociations(t *testing.T) {

	DropTables(db)
	MigrateModels(db)

	// create copies of global test data
	dashboardA := DashboardA
	dashboardB := DashboardB
	widgetA := WidgetA
	widgetB := WidgetB

	// add dashboards to DB
	assert.NoError(t, db.Create(&dashboardA).Error)
	assert.NoError(t, db.Create(&dashboardB).Error)

	// add widgets to DB
	assert.NoError(t, db.Create(&widgetA).Error)
	assert.NoError(t, db.Create(&widgetB).Error)

	// add dashboard has many widgets associations to DB
	assert.NoError(t, db.Model(&dashboardA).Association("Widgets").Append(&widgetA).Error)
	assert.NoError(t, db.Model(&dashboardA).Association("Widgets").Append(&widgetB).Error)

	var dashboard1 Dashboard
	assert.NoError(t, db.Find(&dashboard1, 1).Error, fmt.Sprintf("Find Dashboard with ID=1"))
	assert.EqualValues(t, "Dashboard_A", dashboard1.Name)

	//Get widgets of dashboard1
	var widgets []Widget
	assert.NoError(t, db.Model(&dashboard1).Related(&widgets, "Widgets").Error)
	if len(widgets) != 2 {
		assert.Fail(t, "Dashboard Associations",
			"Expected to have %v Widget. Has %v.", 2, len(widgets))
	}
}

func TestWidgetAssociations(t *testing.T) {

	DropTables(db)
	MigrateModels(db)

	// create copies of global test data
	widgetA := WidgetA
	widgetB := WidgetB
	fileA := FileA
	fileB := FileB
	fileC := FileC
	fileD := FileD

	// add widgets to DB
	assert.NoError(t, db.Create(&widgetA).Error)
	assert.NoError(t, db.Create(&widgetB).Error)

	// add files to DB
	assert.NoError(t, db.Create(&fileA).Error)
	assert.NoError(t, db.Create(&fileB).Error)
	assert.NoError(t, db.Create(&fileC).Error)
	assert.NoError(t, db.Create(&fileD).Error)

	// add widget has many files associations to DB
	assert.NoError(t, db.Model(&widgetA).Association("Files").Append(&fileA).Error)
	assert.NoError(t, db.Model(&widgetA).Association("Files").Append(&fileB).Error)

	var widget1 Widget
	assert.NoError(t, db.Find(&widget1, 1).Error, fmt.Sprintf("Find Widget with ID=1"))
	assert.EqualValues(t, widgetA.Name, widget1.Name)

	// Get files of widget
	var files []File
	assert.NoError(t, db.Model(&widget1).Related(&files, "Files").Error)
	if len(files) != 2 {
		assert.Fail(t, "Widget Associations",
			"Expected to have %v Files. Has %v.", 2, len(files))
	}
}

func TestFileAssociations(t *testing.T) {

	DropTables(db)
	MigrateModels(db)

	// create copies of global test data
	fileA := FileA
	fileB := FileB
	fileC := FileC
	fileD := FileD

	// add files to DB
	assert.NoError(t, db.Create(&fileA).Error)
	assert.NoError(t, db.Create(&fileB).Error)
	assert.NoError(t, db.Create(&fileC).Error)
	assert.NoError(t, db.Create(&fileD).Error)

	var file1 File
	assert.NoError(t, db.Find(&file1, 1).Error, fmt.Sprintf("Find File with ID=1"))
	assert.EqualValues(t, "File_A", file1.Name)
}

func TestAddAdmin(t *testing.T) {
	DropTables(db)
	MigrateModels(db)

	assert.NoError(t, DBAddAdminUser(db))
}

func TestAddAdminAndUsers(t *testing.T) {
	DropTables(db)
	MigrateModels(db)

	assert.NoError(t, DBAddAdminAndUserAndGuest(db))
}

func TestAddData(t *testing.T) {
	DropTables(db)
	MigrateModels(db)

	assert.NoError(t, DBAddTestData(db))
}
