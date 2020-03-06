/** Metrics package, endpoints.
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
package metrics

import (
	"github.com/chenjiandongx/ginprom"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	ICCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "infrastructure_components",
			Help: "A counter for the total number of infrastructure_components",
		},
	)

	ComponentConfigurationCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "component_configurations",
			Help: "A counter for the total number of component configurations",
		},
	)

	FileCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "files",
			Help: "A counter for the total number of files",
		},
	)

	ScenarioCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "scenarios",
			Help: "A counter for the total number of scenarios",
		},
	)

	UserCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "users",
			Help: "A counter for the total number of users",
		},
	)

	DashboardCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "dashboards",
			Help: "A counter for the total number of dashboards",
		},
	)
)

// RegisterMetricsEndpoint godoc
// @Summary Prometheus metrics endpoint
// @ID getMetrics
// @Produce  json
// @Tags metrics
// @Success 200 "Returns Prometheus metrics"
// @Router /metrics [get]
func RegisterMetricsEndpoint(rg *gin.RouterGroup) {
	// use prometheus metrics exporter middleware.
	//
	// ginprom.PromMiddleware() expects a ginprom.PromOpts{} poniter.
	// It was used for filtering labels with regex. `nil` will pass every requests.
	//
	// ginprom promethues-labels:
	//   `status`, `endpoint`, `method`
	//
	// for example:
	// 1). I want not to record the 404 status request. That's easy for it.
	// ginprom.PromMiddleware(&ginprom.PromOpts{ExcludeRegexStatus: "404"})
	//
	// 2). And I wish ignore endpoint start with `/prefix`.
	// ginprom.PromMiddleware(&ginprom.PromOpts{ExcludeRegexEndpoint: "^/prefix"})
	r := gin.Default()
	r.Use(ginprom.PromMiddleware(nil))

	rg.GET("", ginprom.PromHandler(promhttp.Handler()))

	// Register metrics
	prometheus.MustRegister(
		ICCounter,
		ComponentConfigurationCounter,
		FileCounter,
		ScenarioCounter,
		UserCounter,
		DashboardCounter,
	)
}

func InitCounters(db *gorm.DB) {
	var infrastructure_components, component_configurations, files, scenarios, users, dashboards float64

	db.Table("infrastructure_components").Count(&infrastructure_components)
	db.Table("component_configurations").Count(&component_configurations)
	db.Table("files").Count(&files)
	db.Table("scenarios").Count(&scenarios)
	db.Table("users").Count(&users)
	db.Table("dashboards").Count(&dashboards)

	ICCounter.Add(infrastructure_components)
	ComponentConfigurationCounter.Add(component_configurations)
	FileCounter.Add(files)
	ScenarioCounter.Add(scenarios)
	UserCounter.Add(users)
	DashboardCounter.Add(dashboards)
}
