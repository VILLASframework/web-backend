package main

import (
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
)

func main() {
	// Testing
	db := common.InitDB()
	defer db.Close()

	//Testing dependencies
	var testProject common.Project
	testProject.Name = "MyAweSomeProject"

	var testFile common.File
	testFile.Name = "MyAwesomeFilename"

	var testSimulation common.Simulation
	testSimulation.Name = "MyAwesomeSimulation"

	var testSimulationModel common.SimulationModel
	testSimulationModel.Name = "MyAwesomeSimulationModel"

	var testSimulator common.Simulator
	testSimulator.Host = "SimulatorHost"

	var testUser common.User
	testUser.Username = "MyUserName"

	var testVis common.Visualization
	testVis.Name = "MyAwesomeVisualization"

}
