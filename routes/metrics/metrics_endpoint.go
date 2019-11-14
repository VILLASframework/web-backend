package metrics

import (
	"github.com/chenjiandongx/ginprom"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	SimulatorCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "simulators",
			Help: "A counter for the total number of simulators",
		},
	)

	SimulationModelCounter = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "simulation_models",
			Help: "A counter for the total number of simulation models",
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
		SimulatorCounter,
		SimulationModelCounter,
		FileCounter,
		ScenarioCounter,
		UserCounter,
		DashboardCounter,
	)
}

func InitCounters(db *gorm.DB) {
	var simulators, simulation_models, files, scenarios, users, dashboards float64

	db.Table("simulators").Count(&simulators)
	db.Table("simulation_models").Count(&simulation_models)
	db.Table("files").Count(&files)
	db.Table("scenarios").Count(&scenarios)
	db.Table("users").Count(&users)
	db.Table("dashboards").Count(&dashboards)

	SimulatorCounter.Add(simulators)
	SimulationModelCounter.Add(simulation_models)
	FileCounter.Add(files)
	ScenarioCounter.Add(scenarios)
	UserCounter.Add(users)
	DashboardCounter.Add(dashboards)
}
