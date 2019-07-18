package common

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Verify that you can connect to the database
func TestDBConnection(t *testing.T) {
	db := InitDB()
	defer db.Close()

	assert.NoError(t, VerifyConnection(db), "DB must ping")
}

// Verify that the associations between each model are done properly
func TestDummyDBAssociations(t *testing.T) {
	a := assert.New(t)

	// find model string lambda
	fM := func(s string) string { return fmt.Sprintf("Find %s with ID=1", s) }

	db := DummyInitDB()
	defer db.Close()

	DummyPopulateDB(db)

	// Variables for tests
	var simr Simulator
	var mo SimulationModel
	var file File
	var so Scenario
	var usr User
	var usrs []User
	var dab Dashboard
	var widg Widget

	var sigs []Signal
	var mos []SimulationModel
	var files []File
	var files_sm []File
	var sos []Scenario
	var dabs []Dashboard
	var widgs []Widget

	// User

	a.NoError(db.Find(&usr, 2).Error, fM("User"))
	a.EqualValues("User_A", usr.Username)

	// User Associations

	a.NoError(db.Model(&usr).Related(&sos, "Scenarios").Error)
	if len(sos) != 2 {
		a.Fail("User Associations",
			"Expected to have %v Scenarios. Has %v.", 2, len(sos))
	}

	// Scenario

	a.NoError(db.Find(&so, 1).Error, fM("Scenario"))
	a.EqualValues("Scenario_A", so.Name)

	// Scenario Associations

	a.NoError(db.Model(&so).Association("Users").Find(&usrs).Error)
	if len(usrs) != 2 {
		a.Fail("Scenario Associations",
			"Expected to have %v Users. Has %v.", 2, len(usrs))
	}

	a.NoError(db.Model(&so).Related(&mos, "SimulationModels").Error)
	if len(mos) != 2 {
		a.Fail("Scenario Associations",
			"Expected to have %v simulation models. Has %v.", 2, len(mos))
	}

	a.NoError(db.Model(&so).Related(&dabs, "Dashboards").Error)
	if len(dabs) != 2 {
		a.Fail("Scenario Associations",
			"Expected to have %v Dashboards. Has %v.", 2, len(dabs))
	}

	// Simulator
	a.NoError(db.Find(&simr, 1).Error, fM("Simulator"))
	a.EqualValues("Host_A", simr.Host)

	// Simulator Associations
	a.NoError(db.Model(&simr).Association("SimulationModels").Find(&mos).Error)
	if len(mos) != 2 {
		a.Fail("Simulator Associations",
			"Expected to have %v SimulationModels. Has %v.", 2, len(mos))
	}

	// SimulationModel

	a.NoError(db.Find(&mo, 1).Error, fM("SimulationModel"))
	a.EqualValues("SimulationModel_A", mo.Name)

	// SimulationModel Associations

	a.NoError(db.Model(&mo).Where("Direction = ?", "out").Related(&sigs, "OutputMapping").Error)
	if len(sigs) != 2 {
		a.Fail("SimulationModel Associations",
			"Expected to have %v Output Signals. Has %v.", 2, len(sigs))
	}

	a.NoError(db.Model(&mo).Related(&files_sm, "Files").Error)
	if len(files_sm) != 2 {
		a.Fail("SimulationModel Associations",
			"Expected to have %v Files. Has %v.", 2, len(files_sm))
	}

	fmt.Println("SimulatorID: ", mo.SimulatorID)

	// Dashboard

	a.NoError(db.Find(&dab, 1).Error, fM("Dashboard"))
	a.EqualValues("Dashboard_A", dab.Name)

	// Dashboard Associations

	a.NoError(db.Model(&dab).Related(&widgs, "Widgets").Error)
	if len(widgs) != 2 {
		a.Fail("Widget Associations",
			"Expected to have %v Widget. Has %v.", 2, len(widgs))
	}

	// Widget
	a.NoError(db.Find(&widg, 1).Error, fM("Widget"))
	a.EqualValues("Widget_A", widg.Name)

	// Widget Association
	a.NoError(db.Model(&widg).Related(&files, "Files").Error)
	if len(files) != 2 {
		a.Fail("Widget Associations",
			"Expected to have %v Files. Has %v.", 2, len(files))
	}

	// File

	a.NoError(db.Find(&file, 1).Error, fM("File"))
	a.EqualValues("File_A", file.Name)

}
