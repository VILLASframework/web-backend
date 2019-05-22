package endpoints

import (
	"net/http"
	"testing"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
)

// Test /simulation endpoints
func TestSimulationEndpoints(t *testing.T) {

	db := common.DummyInitDB()
	defer db.Close()
	common.DummyPopulateDB(db)


	router := gin.Default()
	api := router.Group("/api")
	SimulationsRegister(api.Group("/simulations"))

	w := httptest.NewRecorder()

	// test GET simulations/
	req, _ := http.NewRequest("GET", "/api/simulations/", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	var simulations_response = "{\"simulations\":[{\"Name\":\"Simulation_A\",\"SimulationID\":1,\"Running\":false,\"Starting Parameters\":null},{\"Name\":\"Simulation_B\",\"SimulationID\":2,\"Running\":false,\"Starting Parameters\":null}]}"
	assert.Equal(t, simulations_response, w.Body.String())

	// test get simulations/:SimulationID
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/simulations/1", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	simulations_response = "{\"simulation\":{\"Name\":\"Simulation_A\",\"SimulationID\":1,\"Running\":false,\"Starting Parameters\":null}}"
	assert.Equal(t, simulations_response, w.Body.String())

	// test get simulations/:SimulationID/users
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/simulations/1/users", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	simulations_response = "{\"users\":[{\"Username\":\"User_A\",\"Role\":\"user\",\"Mail\":\"\"},{\"Username\":\"User_B\",\"Role\":\"user\",\"Mail\":\"\"}]}"
	assert.Equal(t, simulations_response, w.Body.String())

	// TODO add more tests

}
