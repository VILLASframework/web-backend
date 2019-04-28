package common

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Verify that you can connect to the database
func TestDBConnection(t *testing.T) {
	db := InitDB()
	defer db.Close()

	assert.NoError(t, VerifyConnection(db), "DB must ping")
}

// Verify that the associations between each model are done properly
func TestDummyDBAssociations(t *testing.T) {
	assert := assert.New(t)

	// find model string lambda
	fM := func(s string) string { return fmt.Sprintf("Find %s with ID=1", s) }

	db := DummyInitDB()
	defer db.Close()

	DummyPopulateDB(db)

	// Variables for tests
	var simr Simulator
	var smo SimulationModel
	var file File
	var proj Project
	var simn Simulation
	var usr User
	var vis Visualization

	var sigs []Signal
	var smos []SimulationModel
	var files []File
	var projs []Project
	var simns []Simulation
	var viss []Visualization
	var widgs []Widget

	// Simulation Model

	assert.NoError(db.Find(&smo, 1).Error, fM("SimulationModel"))
	assert.EqualValues("SimModel_A", smo.Name)

	// Simulation Model Associations

	assert.NoError(db.Model(&smo).Association("BelongsToSimulation").Find(&simn).Error)
	assert.EqualValues("Simulation_A", simn.Name, "Expected Simulation_A")

	assert.NoError(db.Model(&smo).Association("BelongsToSimulator").Find(&simr).Error)
	assert.EqualValues("Host_A", simr.Host, "Expected Host_A")

	assert.NoError(db.Model(&smo).Related(&sigs, "OutputMapping").Error)
	if len(sigs) != 2 {
		assert.Fail("Simulation Model Associations",
			"Expected to have %v Output Signals. Has %v.", 2, len(sigs))
	}

	assert.NoError(db.Model(&smo).Related(&sigs, "InputMapping").Error)
	if len(sigs) != 2 {
		assert.Fail("Simulation Model Associations",
			"Expected to have %v Input Signals. Has %v.", 2, len(sigs))
	}

	// Simulation

	assert.NoError(db.Find(&simn, 1).Error, fM("Simulation"))
	assert.EqualValues("Simulation_A", simn.Name)

	// Simulation Associations

	assert.NoError(db.Model(&simn).Association("User").Find(&usr).Error)
	assert.EqualValues("User_A", usr.Username)

	assert.NoError(db.Model(&simn).Related(&smos, "Models").Error)
	if len(smos) != 2 {
		assert.Fail("Simulation Associations",
			"Expected to have %v Simulation Models. Has %v.", 2, len(smos))
	}

	assert.NoError(db.Model(&simn).Related(&projs, "Projects").Error)
	if len(projs) != 2 {
		assert.Fail("Simulation Associations",
			"Expected to have %v Projects. Has %v.", 2, len(projs))
	}

	// Project

	assert.NoError(db.Find(&proj, 1).Error, fM("Project"))
	assert.EqualValues("Project_A", proj.Name)

	// Project Associations

	assert.NoError(db.Model(&proj).Association("Simulation").Find(&simn).Error)
	assert.EqualValues("Simulation_A", simn.Name)

	assert.NoError(db.Model(&proj).Association("User").Find(&usr).Error)
	assert.EqualValues("User_A", usr.Username)

	assert.NoError(db.Model(&proj).Related(&viss, "Visualizations").Error)
	if len(viss) != 2 {
		assert.Fail("Project Associations",
			"Expected to have %v Visualizations. Has %v.", 2, len(viss))
	}

	// User

	assert.NoError(db.Find(&usr, 1).Error, fM("User"))
	assert.EqualValues("User_A", usr.Username)

	// User Associations

	assert.NoError(db.Model(&usr).Related(&projs, "Projects").Error)
	if len(projs) != 2 {
		assert.Fail("User Associations",
			"Expected to have %v Projects. Has %v.", 2, len(projs))
	}

	assert.NoError(db.Model(&usr).Related(&simns, "Simulations").Error)
	if len(simns) != 2 {
		assert.Fail("User Associations",
			"Expected to have %v Simulations. Has %v.", 2, len(simns))
	}

	assert.NoError(db.Model(&usr).Related(&files, "Files").Error)
	if len(files) != 2 {
		assert.Fail("User Associations",
			"Expected to have %v Files. Has %v.", 2, len(files))
	}

	// Visualization

	assert.NoError(db.Find(&vis, 1).Error, fM("Visualization"))
	assert.EqualValues("Visualization_A", vis.Name)

	// Visualization Associations

	assert.NoError(db.Model(&vis).Association("Project").Find(&proj).Error)
	assert.EqualValues("Project_A", proj.Name)

	assert.NoError(db.Model(&vis).Association("User").Find(&usr).Error)
	assert.EqualValues("User_A", usr.Username)

	assert.NoError(db.Model(&vis).Related(&widgs, "Widgets").Error)
	if len(widgs) != 2 {
		assert.Fail("Widget Associations",
			"Expected to have %v Widget. Has %v.", 2, len(widgs))
	}

	// File

	assert.NoError(db.Find(&file, 1).Error, fM("File"))
	assert.EqualValues("File_A", file.Name)

	// File Associations

	assert.NoError(db.Model(&file).Association("User").Find(&usr).Error)
	assert.EqualValues("User_A", usr.Username)

}
