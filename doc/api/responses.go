/** Docs package, responses.
*
* @author Sonja Happ <sonja.happ@eonerc.rwth-aachen.de>
* @copyright 2014-2019, Institute for Automation of Complex Power Systems, EONERC
* @license GNU General Public License (version 3)
*
* VILLASweb-backend-go
*
* This program is free software: you can redistribute it and/or modify
* it under the terms of the GNU General Public License as published by
* the Free Software Foundation, either version 3 of the License, or
* any later version.
*
* This program is distributed in the hope that it will be useful,
* but WITHOUT ANY WARRANTY; without even the implied warranty of
* MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
* GNU General Public License for more details.
*
* You should have received a copy of the GNU General Public License
* along with this program.  If not, see <http://www.gnu.org/licenses/>.
*********************************************************************************/
package docs

import "git.rwth-aachen.de/acs/public/villas/web-backend-go/database"

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

type ResponseICs struct {
	ics []database.InfrastructureComponent
}

type ResponseIC struct {
	ic database.InfrastructureComponent
}

type ResponseScenarios struct {
	scenarios []database.Scenario
}

type ResponseScenario struct {
	scenario database.Scenario
}

type ResponseConfigs struct {
	configs []database.ComponentConfiguration
}

type ResponseConfig struct {
	config database.ComponentConfiguration
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

type ResponseResults struct {
	results []database.Result
}

type ResponseResult struct {
	result database.Result
}
