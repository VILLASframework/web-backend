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
	var simn Simulation
	var usr User
	var usrs []User
	var vis Visualization
	var widg Widget

	var sigs []Signal
	var mos []SimulationModel
	var files []File
	var files_sm []File
	var simns []Simulation
	var viss []Visualization
	var widgs []Widget

	// User

	a.NoError(db.Find(&usr, 2).Error, fM("User"))
	a.EqualValues("User_A", usr.Username)

	// User Associations

	a.NoError(db.Model(&usr).Related(&simns, "Simulations").Error)
	if len(simns) != 2 {
		a.Fail("User Associations",
			"Expected to have %v Simulations. Has %v.", 2, len(simns))
	}

	// Simulation

	a.NoError(db.Find(&simn, 1).Error, fM("Simulation"))
	a.EqualValues("Simulation_A", simn.Name)

	// Simulation Associations

	a.NoError(db.Model(&simn).Association("Users").Find(&usrs).Error)
	if len(usrs) != 2 {
		a.Fail("Simulations Associations",
			"Expected to have %v Users. Has %v.", 2, len(usrs))
	}

	a.NoError(db.Model(&simn).Related(&mos, "SimulationModels").Error)
	if len(mos) != 2 {
		a.Fail("Simulation Associations",
			"Expected to have %v simulation models. Has %v.", 2, len(mos))
	}

	a.NoError(db.Model(&simn).Related(&viss, "Visualizations").Error)
	if len(viss) != 2 {
		a.Fail("Simulation Associations",
			"Expected to have %v Visualizations. Has %v.", 2, len(viss))
	}

	// SimulationModel

	a.NoError(db.Find(&mo, 1).Error, fM("SimulationModel"))
	a.EqualValues("SimulationModel_A", mo.Name)

	// SimulationModel Associations

	a.NoError(db.Model(&mo).Association("Simulator").Find(&simr).Error)
	a.EqualValues("Host_A", simr.Host, "Expected Host_A")

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

	// Visualization

	a.NoError(db.Find(&vis, 1).Error, fM("Visualization"))
	a.EqualValues("Visualization_A", vis.Name)

	// Visualization Associations

	a.NoError(db.Model(&vis).Related(&widgs, "Widgets").Error)
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
