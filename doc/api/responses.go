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
	success bool
	users   []common.User
}

type ResponseUser struct {
	success bool
	user    common.User
}

type ResponseSimulators struct {
	success    bool
	simulators []common.Simulator
}

type ResponseSimulator struct {
	success   bool
	simulator common.Simulator
}

type ResponseScenarios struct {
	success   bool
	scenarios []common.Scenario
}

type ResponseScenario struct {
	success  bool
	scenario common.Scenario
}
