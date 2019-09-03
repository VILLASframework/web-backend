package common

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var db *gorm.DB

// find model string lambda
func fM(s string, id uint) string {
	return fmt.Sprintf("Find %s with ID=%d", s, id)
}

func TestMain(m *testing.M) {
	db = DummyInitDB()
	defer db.Close()

	os.Exit(m.Run())
}

// Verify that you can connect to the database
func TestDBConnection(t *testing.T) {
	db := InitDB()
	defer db.Close()

	assert.NoError(t, VerifyConnection(db), "DB must ping")
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
	assert.NoError(t, db.Find(&usr1, "ID = ?", 2).Error, fM("User", 2))
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
	assert.NoError(t, db.Find(&scenario1, 1).Error, fM("Scenario", 1))
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

func TestSimulatorAssociations(t *testing.T) {

	DropTables(db)
	MigrateModels(db)

	// create copies of global test data
	simulatorA := SimulatorA
	simulatorB := SimulatorB
	modelA := SimulationModelA
	modelB := SimulationModelB

	// add simulators to DB
	assert.NoError(t, db.Create(&simulatorA).Error)
	assert.NoError(t, db.Create(&simulatorB).Error)

	// add simulation models to DB
	assert.NoError(t, db.Create(&modelA).Error)
	assert.NoError(t, db.Create(&modelB).Error)

	// add simulator has many simulation models association to DB
	assert.NoError(t, db.Model(&simulatorA).Association("SimulationModels").Append(&modelA).Error)
	assert.NoError(t, db.Model(&simulatorA).Association("SimulationModels").Append(&modelB).Error)

	var simulator1 Simulator
	assert.NoError(t, db.Find(&simulator1, 1).Error, fM("Simulator", 1))
	assert.EqualValues(t, "Host_A", simulator1.Host)

	// Get simulation models of simulator1
	var models []SimulationModel
	assert.NoError(t, db.Model(&simulator1).Association("SimulationModels").Find(&models).Error)
	if len(models) != 2 {
		assert.Fail(t, "Simulator Associations",
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
	simulatorA := SimulatorA
	simulatorB := SimulatorB

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

	// add simulators to DB
	assert.NoError(t, db.Create(&simulatorA).Error)
	assert.NoError(t, db.Create(&simulatorB).Error)

	// add simulation model has many signals associations
	assert.NoError(t, db.Model(&modelA).Association("InputMapping").Append(&inSignalA).Error)
	assert.NoError(t, db.Model(&modelA).Association("InputMapping").Append(&inSignalB).Error)
	assert.NoError(t, db.Model(&modelA).Association("OutputMapping").Append(&outSignalA).Error)
	assert.NoError(t, db.Model(&modelA).Association("OutputMapping").Append(&outSignalB).Error)

	// add simulation model has many files associations
	assert.NoError(t, db.Model(&modelA).Association("Files").Append(&fileC).Error)
	assert.NoError(t, db.Model(&modelA).Association("Files").Append(&fileD).Error)

	// associate simulation models with simulators
	assert.NoError(t, db.Model(&simulatorA).Association("SimulationModels").Append(&modelA).Error)
	assert.NoError(t, db.Model(&simulatorA).Association("SimulationModels").Append(&modelB).Error)

	var model1 SimulationModel
	assert.NoError(t, db.Find(&model1, 1).Error, fM("SimulationModel", 1))
	assert.EqualValues(t, "SimulationModel_A", model1.Name)

	// Check simulator ID
	if model1.SimulatorID != 1 {
		assert.Fail(t, "Simulation Model expected to have Simulator ID 1, but is %v", model1.SimulatorID)
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
	assert.NoError(t, db.Find(&dashboard1, 1).Error, fM("Dashboard", 1))
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
	assert.NoError(t, db.Find(&widget1, 1).Error, fM("Widget", 1))
	assert.EqualValues(t, "Widget_A", widget1.Name)

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
	assert.NoError(t, db.Find(&file1, 1).Error, fM("File", 1))
	assert.EqualValues(t, "File_A", file1.Name)
}
