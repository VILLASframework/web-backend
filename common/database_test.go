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

	DummyPopulateDB(db)

	os.Exit(m.Run())
}

// Verify that you can connect to the database
func TestDBConnection(t *testing.T) {
	db := InitDB()
	defer db.Close()

	assert.NoError(t, VerifyConnection(db), "DB must ping")
}

func TestUserAssociations(t *testing.T) {
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
	var file1 File
	assert.NoError(t, db.Find(&file1, 1).Error, fM("File", 1))
	assert.EqualValues(t, "File_A", file1.Name)
}
