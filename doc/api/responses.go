package docs

import "git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"

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
	user    common.User
}

type ResponseUsers struct {
	users []common.User
}

type ResponseUser struct {
	user common.User
}

type ResponseSimulators struct {
	simulators []common.Simulator
}

type ResponseSimulator struct {
	simulator common.Simulator
}

type ResponseScenarios struct {
	scenarios []common.Scenario
}

type ResponseScenario struct {
	scenario common.Scenario
}

type ResponseSimulationModels struct {
	models []common.SimulationModel
}

type ResponseSimulationModel struct {
	model common.SimulationModel
}

type ResponseDashboards struct {
	dashboards []common.Dashboard
}

type ResponseDashboard struct {
	dashboard common.Dashboard
}

type ResponseWidgets struct {
	widgets []common.Widget
}

type ResponseWidget struct {
	widget common.Widget
}

type ResponseSignals struct {
	signals []common.Signal
}

type ResponseSignal struct {
	signal common.Signal
}
