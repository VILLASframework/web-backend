package docs

import "git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/database"

// This file defines the responses to any endpoint in the backend
// The defined structures are only used for documentation purposes with swaggo and are NOT used in the code

type ResponseError struct {
	success bool
	message string
}

type ResponseAuthenticate struct {
	success bool
	token   string
	message string
	user    database.User
}

type ResponseUsers struct {
	users []database.User
}

type ResponseUser struct {
	user database.User
}

type ResponseSimulators struct {
	simulators []database.Simulator
}

type ResponseSimulator struct {
	simulator database.Simulator
}

type ResponseScenarios struct {
	scenarios []database.Scenario
}

type ResponseScenario struct {
	scenario database.Scenario
}

type ResponseSimulationModels struct {
	models []database.SimulationModel
}

type ResponseSimulationModel struct {
	model database.SimulationModel
}

type ResponseDashboards struct {
	dashboards []database.Dashboard
}

type ResponseDashboard struct {
	dashboard database.Dashboard
}

type ResponseWidgets struct {
	widgets []database.Widget
}

type ResponseWidget struct {
	widget database.Widget
}

type ResponseSignals struct {
	signals []database.Signal
}

type ResponseSignal struct {
	signal database.Signal
}

type ResponseFiles struct {
	files []database.File
}

type ResponseFile struct {
	file database.File
}
