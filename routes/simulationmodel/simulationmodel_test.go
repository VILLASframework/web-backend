/** Simulationmodel package, testing.
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
package simulationmodel

import (
	"fmt"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/configuration"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/helper"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/infrastructure-component"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/scenario"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/user"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var router *gin.Engine
var db *gorm.DB

type SimulationModelRequest struct {
	Name                string         `json:"name,omitempty"`
	ScenarioID          uint           `json:"scenarioID,omitempty"`
	ICID                uint           `json:"icID,omitempty"`
	StartParameters     postgres.Jsonb `json:"startParameters,omitempty"`
	SelectedModelFileID uint           `json:"selectedModelFileID,omitempty"`
}

type ICRequest struct {
	UUID       string         `json:"uuid,omitempty"`
	Host       string         `json:"host,omitempty"`
	Modeltype  string         `json:"modelType,omitempty"`
	State      string         `json:"state,omitempty"`
	Properties postgres.Jsonb `json:"properties,omitempty"`
}

type ScenarioRequest struct {
	Name            string         `json:"name,omitempty"`
	Running         bool           `json:"running,omitempty"`
	StartParameters postgres.Jsonb `json:"startParameters,omitempty"`
}

func addScenarioAndIC() (scenarioID uint, ICID uint) {

	// authenticate as admin
	token, _ := helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.AdminCredentials)

	// POST $newICA
	newICA := ICRequest{
		UUID:       database.ICA.UUID,
		Host:       database.ICA.Host,
		Modeltype:  database.ICA.Modeltype,
		State:      database.ICA.State,
		Properties: database.ICA.Properties,
	}
	_, resp, _ := helper.TestEndpoint(router, token,
		"/api/ic", "POST", helper.KeyModels{"ic": newICA})

	// Read newIC's ID from the response
	newICID, _ := helper.GetResponseID(resp)

	// POST a second IC to change to that IC during testing
	newICB := ICRequest{
		UUID:       database.ICB.UUID,
		Host:       database.ICB.Host,
		Modeltype:  database.ICB.Modeltype,
		State:      database.ICB.State,
		Properties: database.ICB.Properties,
	}
	_, resp, _ = helper.TestEndpoint(router, token,
		"/api/ic", "POST", helper.KeyModels{"ic": newICB})

	// authenticate as normal user
	token, _ = helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserACredentials)

	// POST $newScenario
	newScenario := ScenarioRequest{
		Name:            database.ScenarioA.Name,
		Running:         database.ScenarioA.Running,
		StartParameters: database.ScenarioA.StartParameters,
	}
	_, resp, _ = helper.TestEndpoint(router, token,
		"/api/scenarios", "POST", helper.KeyModels{"scenario": newScenario})

	// Read newScenario's ID from the response
	newScenarioID, _ := helper.GetResponseID(resp)

	// add the guest user to the new scenario
	_, resp, _ = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/scenarios/%v/user?username=User_C", newScenarioID), "PUT", nil)

	return uint(newScenarioID), uint(newICID)
}

func TestMain(m *testing.M) {
	err := configuration.InitConfig()
	if err != nil {
		panic(m)
	}

	db, err = database.InitDB(configuration.GolbalConfig)
	if err != nil {
		panic(m)
	}
	defer db.Close()

	router = gin.Default()
	api := router.Group("/api")

	user.RegisterAuthenticate(api.Group("/authenticate"))
	api.Use(user.Authentication(true))
	RegisterSimulationModelEndpoints(api.Group("/models"))
	// scenario endpoints required here to first add a scenario to the DB
	// that can be associated with a new simulation model
	scenario.RegisterScenarioEndpoints(api.Group("/scenarios"))
	// IC endpoints required here to first add a IC to the DB
	// that can be associated with a new simulation model
	infrastructure_component.RegisterICEndpoints(api.Group("/ic"))

	os.Exit(m.Run())
}

func TestAddSimulationModel(t *testing.T) {
	database.DropTables(db)
	database.MigrateModels(db)
	assert.NoError(t, database.DBAddAdminAndUserAndGuest(db))

	// prepare the content of the DB for testing
	// by adding a scenario and a IC to the DB
	// using the respective endpoints of the API
	scenarioID, ICID := addScenarioAndIC()

	newSimulationModel := SimulationModelRequest{
		Name:                database.SimulationModelA.Name,
		ScenarioID:          scenarioID,
		ICID:                ICID,
		StartParameters:     database.SimulationModelA.StartParameters,
		SelectedModelFileID: database.SimulationModelA.SelectedModelFileID,
	}

	// authenticate as normal userB who has no access to new scenario
	token, err := helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserBCredentials)
	assert.NoError(t, err)

	// try to POST with no access
	// should result in unprocessable entity
	code, resp, err := helper.TestEndpoint(router, token,
		"/api/models", "POST", helper.KeyModels{"simulationModel": newSimulationModel})
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

	// authenticate as normal user
	token, err = helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserACredentials)
	assert.NoError(t, err)

	// try to POST non JSON body
	code, resp, err = helper.TestEndpoint(router, token,
		"/api/models", "POST", "this is not JSON")
	assert.NoError(t, err)
	assert.Equalf(t, 400, code, "Response body: \n%v\n", resp)

	// authenticate as normal user
	token, err = helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserACredentials)
	assert.NoError(t, err)

	// test POST models/ $newSimulationModel
	code, resp, err = helper.TestEndpoint(router, token,
		"/api/models", "POST", helper.KeyModels{"simulationModel": newSimulationModel})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare POST's response with the newSimulationModel
	err = helper.CompareResponse(resp, helper.KeyModels{"simulationModel": newSimulationModel})
	assert.NoError(t, err)

	// Read newSimulationModel's ID from the response
	newSimulationModelID, err := helper.GetResponseID(resp)
	assert.NoError(t, err)

	// Get the newSimulationModel
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/models/%v", newSimulationModelID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare GET's response with the newSimulationModel
	err = helper.CompareResponse(resp, helper.KeyModels{"simulationModel": newSimulationModel})
	assert.NoError(t, err)

	// try to POST a malformed simulation model
	// Required fields are missing
	malformedNewSimulationModel := SimulationModelRequest{
		Name: "ThisIsAMalformedRequest",
	}
	// this should NOT work and return a unprocessable entity 442 status code
	code, resp, err = helper.TestEndpoint(router, token,
		"/api/models", "POST", helper.KeyModels{"simulationModel": malformedNewSimulationModel})
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

	// authenticate as normal userB who has no access to new scenario
	token, err = helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserBCredentials)
	assert.NoError(t, err)

	// Try to GET the newSimulationModel with no access
	// Should result in unprocessable entity
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/models/%v", newSimulationModelID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

}

func TestUpdateSimulationModel(t *testing.T) {

	database.DropTables(db)
	database.MigrateModels(db)
	assert.NoError(t, database.DBAddAdminAndUserAndGuest(db))

	// prepare the content of the DB for testing
	// by adding a scenario and a IC to the DB
	// using the respective endpoints of the API
	scenarioID, ICID := addScenarioAndIC()

	// authenticate as normal user
	token, err := helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserACredentials)
	assert.NoError(t, err)

	// test POST models/ $newSimulationModel
	newSimulationModel := SimulationModelRequest{
		Name:                database.SimulationModelA.Name,
		ScenarioID:          scenarioID,
		ICID:                ICID,
		StartParameters:     database.SimulationModelA.StartParameters,
		SelectedModelFileID: database.SimulationModelA.SelectedModelFileID,
	}
	code, resp, err := helper.TestEndpoint(router, token,
		"/api/models", "POST", helper.KeyModels{"simulationModel": newSimulationModel})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Read newSimulationModel's ID from the response
	newSimulationModelID, err := helper.GetResponseID(resp)
	assert.NoError(t, err)

	updatedSimulationModel := SimulationModelRequest{
		Name:            database.SimulationModelB.Name,
		StartParameters: database.SimulationModelB.StartParameters,
	}

	// authenticate as normal userB who has no access to new scenario
	token, err = helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserBCredentials)
	assert.NoError(t, err)

	// try to PUT with no access
	// should result in unprocessable entity
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/models/%v", newSimulationModelID), "PUT", helper.KeyModels{"simulationModel": updatedSimulationModel})
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

	// authenticate as guest user who has access to simulation model
	token, err = helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.GuestCredentials)
	assert.NoError(t, err)

	// try to PUT as guest
	// should NOT work and result in unprocessable entity
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/models/%v", newSimulationModelID), "PUT", helper.KeyModels{"simulationModel": updatedSimulationModel})
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

	// authenticate as normal user
	token, err = helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserACredentials)
	assert.NoError(t, err)

	// try to PUT a non JSON body
	// should result in a bad request
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/models/%v", newSimulationModelID), "PUT", "This is not JSON")
	assert.NoError(t, err)
	assert.Equalf(t, 400, code, "Response body: \n%v\n", resp)

	// test PUT
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/models/%v", newSimulationModelID), "PUT", helper.KeyModels{"simulationModel": updatedSimulationModel})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare PUT's response with the updatedSimulationModel
	err = helper.CompareResponse(resp, helper.KeyModels{"simulationModel": updatedSimulationModel})
	assert.NoError(t, err)

	//Change IC ID to use second IC available in DB
	updatedSimulationModel.ICID = ICID + 1
	// test PUT again
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/models/%v", newSimulationModelID), "PUT", helper.KeyModels{"simulationModel": updatedSimulationModel})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare PUT's response with the updatedSimulationModel
	err = helper.CompareResponse(resp, helper.KeyModels{"simulationModel": updatedSimulationModel})
	assert.NoError(t, err)

	// Get the updatedSimulationModel
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/models/%v", newSimulationModelID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare GET's response with the updatedSimulationModel
	err = helper.CompareResponse(resp, helper.KeyModels{"simulationModel": updatedSimulationModel})
	assert.NoError(t, err)

	// try to update a simulation model that does not exist (should return not found 404 status code)
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/models/%v", newSimulationModelID+1), "PUT", helper.KeyModels{"simulationModel": updatedSimulationModel})
	assert.NoError(t, err)
	assert.Equalf(t, 404, code, "Response body: \n%v\n", resp)
}

func TestDeleteSimulationModel(t *testing.T) {
	database.DropTables(db)
	database.MigrateModels(db)
	assert.NoError(t, database.DBAddAdminAndUserAndGuest(db))

	// prepare the content of the DB for testing
	// by adding a scenario and a IC to the DB
	// using the respective endpoints of the API
	scenarioID, ICID := addScenarioAndIC()

	newSimulationModel := SimulationModelRequest{
		Name:                database.SimulationModelA.Name,
		ScenarioID:          scenarioID,
		ICID:                ICID,
		StartParameters:     database.SimulationModelA.StartParameters,
		SelectedModelFileID: database.SimulationModelA.SelectedModelFileID,
	}

	// authenticate as normal user
	token, err := helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserACredentials)
	assert.NoError(t, err)

	// test POST models/ $newSimulationModel
	code, resp, err := helper.TestEndpoint(router, token,
		"/api/models", "POST", helper.KeyModels{"simulationModel": newSimulationModel})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Read newSimulationModel's ID from the response
	newSimulationModelID, err := helper.GetResponseID(resp)
	assert.NoError(t, err)

	// authenticate as normal userB who has no access to new scenario
	token, err = helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserBCredentials)
	assert.NoError(t, err)

	// try to DELETE with no access
	// should result in unprocessable entity
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/models/%v", newSimulationModelID), "DELETE", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

	// authenticate as normal user
	token, err = helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserACredentials)
	assert.NoError(t, err)

	// Count the number of all the simulation models returned for scenario
	initialNumber, err := helper.LengthOfResponse(router, token,
		fmt.Sprintf("/api/models?scenarioID=%v", scenarioID), "GET", nil)
	assert.NoError(t, err)

	// Delete the added newSimulationModel
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/models/%v", newSimulationModelID), "DELETE", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare DELETE's response with the newSimulationModel
	err = helper.CompareResponse(resp, helper.KeyModels{"simulationModel": newSimulationModel})
	assert.NoError(t, err)

	// Again count the number of all the simulation models returned
	finalNumber, err := helper.LengthOfResponse(router, token,
		fmt.Sprintf("/api/models?scenarioID=%v", scenarioID), "GET", nil)
	assert.NoError(t, err)

	assert.Equal(t, initialNumber-1, finalNumber)
}

func TestGetAllSimulationModelsOfScenario(t *testing.T) {
	database.DropTables(db)
	database.MigrateModels(db)
	assert.NoError(t, database.DBAddAdminAndUserAndGuest(db))

	// prepare the content of the DB for testing
	// by adding a scenario and a IC to the DB
	// using the respective endpoints of the API
	scenarioID, ICID := addScenarioAndIC()

	// authenticate as normal user
	token, err := helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserACredentials)
	assert.NoError(t, err)

	// test POST models/ $newSimulationModel
	newSimulationModel := SimulationModelRequest{
		Name:                database.SimulationModelA.Name,
		ScenarioID:          scenarioID,
		ICID:                ICID,
		StartParameters:     database.SimulationModelA.StartParameters,
		SelectedModelFileID: database.SimulationModelA.SelectedModelFileID,
	}
	code, resp, err := helper.TestEndpoint(router, token,
		"/api/models", "POST", helper.KeyModels{"simulationModel": newSimulationModel})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Count the number of all the simulation models returned for scenario
	NumberOfSimulationModels, err := helper.LengthOfResponse(router, token,
		fmt.Sprintf("/api/models?scenarioID=%v", scenarioID), "GET", nil)
	assert.NoError(t, err)

	assert.Equal(t, 1, NumberOfSimulationModels)

	// authenticate as normal userB who has no access to scenario
	token, err = helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserBCredentials)
	assert.NoError(t, err)

	// try to get models without access
	// should result in unprocessable entity
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/models?scenarioID=%v", scenarioID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

}
