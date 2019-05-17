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
	a := assert.New(t)

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
	var widg Widget

	var sigs []Signal
	var smos []SimulationModel
	var files []File
	var files_sm []File
	var projs []Project
	var simns []Simulation
	var viss []Visualization
	var widgs []Widget

	// Simulation Model

	a.NoError(db.Find(&smo, 1).Error, fM("SimulationModel"))
	a.EqualValues("SimModel_A", smo.Name)

	// Simulation Model Associations

	a.NoError(db.Model(&smo).Association("BelongsToSimulation").Find(&simn).Error)
	a.EqualValues("Simulation_A", simn.Name, "Expected Simulation_A")

	a.NoError(db.Model(&smo).Association("BelongsToSimulator").Find(&simr).Error)
	a.EqualValues("Host_A", simr.Host, "Expected Host_A")

	a.NoError(db.Model(&smo).Related(&sigs, "OutputMapping").Error)
	if len(sigs) != 4 {
		a.Fail("Simulation Model Associations",
			"Expected to have %v Output AND Input Signals. Has %v.", 4, len(sigs))
	}

	a.NoError(db.Model(&smo).Related(&files_sm, "Files").Error)
	if len(files_sm) != 2 {
		a.Fail("Simulation Model Associations",
			"Expected to have %v Files. Has %v.", 2, len(files_sm))
	}

	// Simulation

	a.NoError(db.Find(&simn, 1).Error, fM("Simulation"))
	a.EqualValues("Simulation_A", simn.Name)

	// Simulation Associations

	a.NoError(db.Model(&simn).Association("User").Find(&usr).Error)
	a.EqualValues("User_A", usr.Username)

	a.NoError(db.Model(&simn).Related(&smos, "Models").Error)
	if len(smos) != 2 {
		a.Fail("Simulation Associations",
			"Expected to have %v Simulation Models. Has %v.", 2, len(smos))
	}

	a.NoError(db.Model(&simn).Related(&projs, "Projects").Error)
	if len(projs) != 2 {
		a.Fail("Simulation Associations",
			"Expected to have %v Projects. Has %v.", 2, len(projs))
	}

	// Project

	a.NoError(db.Find(&proj, 1).Error, fM("Project"))
	a.EqualValues("Project_A", proj.Name)

	// Project Associations

	a.NoError(db.Model(&proj).Association("Simulation").Find(&simn).Error)
	a.EqualValues("Simulation_A", simn.Name)

	a.NoError(db.Model(&proj).Association("User").Find(&usr).Error)
	a.EqualValues("User_A", usr.Username)

	a.NoError(db.Model(&proj).Related(&viss, "Visualizations").Error)
	if len(viss) != 2 {
		a.Fail("Project Associations",
			"Expected to have %v Visualizations. Has %v.", 2, len(viss))
	}

	// User

	a.NoError(db.Find(&usr, 1).Error, fM("User"))
	a.EqualValues("User_A", usr.Username)

	// User Associations

	a.NoError(db.Model(&usr).Related(&projs, "Projects").Error)
	if len(projs) != 2 {
		a.Fail("User Associations",
			"Expected to have %v Projects. Has %v.", 2, len(projs))
	}

	a.NoError(db.Model(&usr).Related(&simns, "Simulations").Error)
	if len(simns) != 2 {
		a.Fail("User Associations",
			"Expected to have %v Simulations. Has %v.", 2, len(simns))
	}



	// Visualization

	a.NoError(db.Find(&vis, 1).Error, fM("Visualization"))
	a.EqualValues("Visualization_A", vis.Name)

	// Visualization Associations

	a.NoError(db.Model(&vis).Association("Project").Find(&proj).Error)
	a.EqualValues("Project_A", proj.Name)

	a.NoError(db.Model(&vis).Association("User").Find(&usr).Error)
	a.EqualValues("User_A", usr.Username)

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

	// File Associations

	//a.NoError(db.Model(&file).Association("User").Find(&usr).Error)
	//a.EqualValues("User_A", usr.Username)

}
