/** InfrastructureComponent package, testing.
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
package infrastructure_component

import (
	"encoding/json"
	"fmt"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/helper"
	component_configuration "git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/component-configuration"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/scenario"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/streadway/amqp"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/web-backend-go/configuration"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/user"
)

var router *gin.Engine
var api *gin.RouterGroup
var waitingTime time.Duration = 2

type ICRequest struct {
	UUID                 string         `json:"uuid,omitempty"`
	WebsocketURL         string         `json:"websocketurl,omitempty"`
	APIURL               string         `json:"apiurl,omitempty"`
	Type                 string         `json:"type,omitempty"`
	Name                 string         `json:"name,omitempty"`
	Category             string         `json:"category,omitempty"`
	State                string         `json:"state,omitempty"`
	Location             string         `json:"location,omitempty"`
	Description          string         `json:"description,omitempty"`
	StartParameterScheme postgres.Jsonb `json:"startparameterscheme,omitempty"`
	ManagedExternally    *bool          `json:"managedexternally"`
}

type ScenarioRequest struct {
	Name            string         `json:"name,omitempty"`
	Running         bool           `json:"running,omitempty"`
	StartParameters postgres.Jsonb `json:"startParameters,omitempty"`
}

type ConfigRequest struct {
	Name            string         `json:"name,omitempty"`
	ScenarioID      uint           `json:"scenarioID,omitempty"`
	ICID            uint           `json:"icID,omitempty"`
	StartParameters postgres.Jsonb `json:"startParameters,omitempty"`
	FileIDs         []int64        `json:"fileIDs,omitempty"`
}

type ICAction struct {
	Act        string `json:"action,omitempty"`
	When       int64  `json:"when,omitempty"`
	Properties struct {
		UUID        *string `json:"uuid,omitempty"`
		Name        *string `json:"name,omitempty"`
		Category    *string `json:"category,omitempty"`
		Type        *string `json:"type,omitempty"`
		Location    *string `json:"location,omitempty"`
		WS_url      *string `json:"ws_url,omitempty"`
		API_url     *string `json:"api_url,omitempty"`
		Description *string `json:"description,omitempty"`
	} `json:"properties,omitempty"`
}

func TestMain(m *testing.M) {
	err := configuration.InitConfig()
	if err != nil {
		panic(m)
	}

	err = database.InitDB(configuration.GolbalConfig)
	if err != nil {
		panic(m)
	}
	defer database.DBpool.Close()

	router = gin.Default()
	api = router.Group("/api")

	user.RegisterAuthenticate(api.Group("/authenticate"))
	api.Use(user.Authentication(true))
	RegisterICEndpoints(api.Group("/ic"))
	// component configuration endpoints required to associate an IC with a component config
	component_configuration.RegisterComponentConfigurationEndpoints(api.Group("/configs"))
	// scenario endpoints required here to first add a scenario to the DB
	// that can be associated with a new component configuration
	scenario.RegisterScenarioEndpoints(api.Group("/scenarios"))

	os.Exit(m.Run())
}

func TestAddICAsAdmin(t *testing.T) {
	database.DropTables()
	database.MigrateModels()
	assert.NoError(t, helper.DBAddAdminAndUserAndGuest())

	// check AMQP connection
	err := CheckConnection()
	assert.Errorf(t, err, "connection is nil")

	// connect AMQP client
	// Make sure that AMQP_HOST, AMQP_USER, AMQP_PASS are set
	host, err := configuration.GolbalConfig.String("amqp.host")
	user, err := configuration.GolbalConfig.String("amqp.user")
	pass, err := configuration.GolbalConfig.String("amqp.pass")
	amqpURI := "amqp://" + user + ":" + pass + "@" + host

	// AMQP Connection startup is tested here
	// Not repeated in other tests because it is only needed once
	err = StartAMQP(amqpURI, api)
	assert.NoError(t, err)

	// authenticate as admin
	token, err := helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.AdminCredentials)
	assert.NoError(t, err)

	// try to POST with non JSON body
	// should result in bad request
	code, resp, err := helper.TestEndpoint(router, token,
		"/api/ic", "POST", "This is no JSON")
	assert.NoError(t, err)
	assert.Equalf(t, 400, code, "Response body: \n%v\n", resp)

	// try to POST malformed IC (required fields missing, validation should fail)
	// should result in an unprocessable entity
	newMalformedIC := ICRequest{
		UUID: helper.ICB.UUID,
	}
	code, resp, err = helper.TestEndpoint(router, token,
		"/api/ic", "POST", helper.KeyModels{"ic": newMalformedIC})
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

	// test POST ic/ $newIC
	newIC := ICRequest{
		UUID:                 helper.ICA.UUID,
		WebsocketURL:         helper.ICA.WebsocketURL,
		APIURL:               helper.ICB.APIURL,
		Type:                 helper.ICA.Type,
		Name:                 helper.ICA.Name,
		Category:             helper.ICA.Category,
		State:                helper.ICA.State,
		Location:             helper.ICA.Location,
		Description:          helper.ICA.Description,
		StartParameterScheme: helper.ICA.StartParameterScheme,
		ManagedExternally:    newFalse(),
	}
	code, resp, err = helper.TestEndpoint(router, token,
		"/api/ic", "POST", helper.KeyModels{"ic": newIC})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare POST's response with the newIC
	err = helper.CompareResponse(resp, helper.KeyModels{"ic": newIC})
	assert.NoError(t, err)

	// Read newIC's ID from the response
	newICID, err := helper.GetResponseID(resp)
	assert.NoError(t, err)

	// Get the newIC
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/ic/%v", newICID), "GET", nil)
	assert.NoError(t, err)

	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare GET's response with the newIC
	err = helper.CompareResponse(resp, helper.KeyModels{"ic": newIC})
	assert.NoError(t, err)

	// Try to GET a IC that does not exist
	// should result in not found
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/ic/%v", newICID+1), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 404, code, "Response body: \n%v\n", resp)

	newExternalIC := ICRequest{
		UUID:                 helper.ICB.UUID,
		WebsocketURL:         helper.ICB.WebsocketURL,
		APIURL:               helper.ICB.APIURL,
		Type:                 helper.ICB.Type,
		Name:                 helper.ICB.Name,
		Category:             helper.ICB.Category,
		State:                helper.ICB.State,
		Location:             helper.ICB.Location,
		Description:          helper.ICB.Description,
		StartParameterScheme: helper.ICB.StartParameterScheme,
		ManagedExternally:    newTrue(),
	}

	// test creation of external IC (should lead to emission of AMQP message to VILLAS)
	code, resp, err = helper.TestEndpoint(router, token,
		"/api/ic", "POST", helper.KeyModels{"ic": newExternalIC})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare POST's response with the newExternalIC
	err = helper.CompareResponse(resp, helper.KeyModels{"ic": newExternalIC})
	assert.NoError(t, err)
}

func TestAddICAsUser(t *testing.T) {
	database.DropTables()
	database.MigrateModels()
	assert.NoError(t, helper.DBAddAdminAndUserAndGuest())

	// authenticate as user
	token, err := helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserACredentials)
	assert.NoError(t, err)

	// test POST ic/ $newIC
	newIC := ICRequest{
		UUID:                 helper.ICA.UUID,
		WebsocketURL:         helper.ICA.WebsocketURL,
		Type:                 helper.ICA.Type,
		Name:                 helper.ICA.Name,
		Category:             helper.ICA.Category,
		State:                helper.ICA.State,
		Location:             helper.ICA.Location,
		Description:          helper.ICA.Description,
		StartParameterScheme: helper.ICA.StartParameterScheme,
		ManagedExternally:    newFalse(),
	}

	// This should fail with unprocessable entity 422 error code
	// Normal users are not allowed to add ICs
	code, resp, err := helper.TestEndpoint(router, token,
		"/api/ic", "POST", helper.KeyModels{"ic": newIC})
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)
}

func TestUpdateICAsAdmin(t *testing.T) {
	database.DropTables()
	database.MigrateModels()
	assert.NoError(t, helper.DBAddAdminAndUserAndGuest())

	// authenticate as admin
	token, err := helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.AdminCredentials)
	assert.NoError(t, err)

	// test POST ic/ $newIC
	newIC := ICRequest{
		UUID:                 helper.ICA.UUID,
		WebsocketURL:         helper.ICA.WebsocketURL,
		Type:                 helper.ICA.Type,
		Name:                 helper.ICA.Name,
		Category:             helper.ICA.Category,
		State:                helper.ICA.State,
		Location:             helper.ICA.Location,
		Description:          helper.ICA.Description,
		StartParameterScheme: helper.ICA.StartParameterScheme,
		ManagedExternally:    newFalse(),
	}
	code, resp, err := helper.TestEndpoint(router, token,
		"/api/ic", "POST", helper.KeyModels{"ic": newIC})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare POST's response with the newIC
	err = helper.CompareResponse(resp, helper.KeyModels{"ic": newIC})
	assert.NoError(t, err)

	// Read newIC's ID from the response
	newICID, err := helper.GetResponseID(resp)
	assert.NoError(t, err)

	// try to PUT with non JSON body
	// should result in bad request
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/ic/%v", newICID), "PUT", "This is no JSON")
	assert.NoError(t, err)
	assert.Equalf(t, 400, code, "Response body: \n%v\n", resp)

	// Test PUT IC
	newIC.WebsocketURL = "ThisIsMyNewURL"
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/ic/%v", newICID), "PUT", helper.KeyModels{"ic": newIC})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare PUT's response with the updated newIC
	err = helper.CompareResponse(resp, helper.KeyModels{"ic": newIC})
	assert.NoError(t, err)

	// Get the updated newIC
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/ic/%v", newICID), "GET", nil)
	assert.NoError(t, err)

	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare GET's response with the updated newIC
	err = helper.CompareResponse(resp, helper.KeyModels{"ic": newIC})
	assert.NoError(t, err)

	// fake an IC update (create) message
	var update ICUpdate
	update.Status = new(ICStatus)
	update.Status.State = new(string)
	*update.Status.State = "idle"
	update.Status.Name = new(string)
	*update.Status.Name = helper.ICB.Name
	update.Status.Category = new(string)
	*update.Status.Category = helper.ICB.Category
	update.Status.Type = new(string)
	*update.Status.Type = helper.ICB.Type

	payload, err := json.Marshal(update)
	assert.NoError(t, err)

	var headers map[string]interface{}
	headers = make(map[string]interface{}) // empty map
	headers["uuid"] = helper.ICB.UUID      // set uuid

	msg := amqp.Publishing{
		DeliveryMode:    2,
		Timestamp:       time.Now(),
		ContentType:     "application/json",
		ContentEncoding: "utf-8",
		Priority:        0,
		Body:            payload,
		Headers:         headers,
	}

	err = CheckConnection()
	assert.NoError(t, err)

	err = client.channel.Publish(VILLAS_EXCHANGE,
		"",
		false,
		false,
		msg)
	assert.NoError(t, err)

	// Wait until externally managed IC is created (happens async)
	time.Sleep(waitingTime * time.Second)

	// try to update this IC
	var updatedIC ICRequest
	updatedIC.Name = "a new name"

	// Should result in forbidden return code 403
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/ic/%v", 2), "PUT", helper.KeyModels{"ic": updatedIC})
	assert.NoError(t, err)
	assert.Equalf(t, 403, code, "Response body: \n%v\n", resp)
}

func TestUpdateICAsUser(t *testing.T) {
	database.DropTables()
	database.MigrateModels()
	assert.NoError(t, helper.DBAddAdminAndUserAndGuest())

	// authenticate as admin
	token, err := helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.AdminCredentials)
	assert.NoError(t, err)

	// test POST ic/ $newIC
	newIC := ICRequest{
		UUID:                 helper.ICA.UUID,
		WebsocketURL:         helper.ICA.WebsocketURL,
		Type:                 helper.ICA.Type,
		Name:                 helper.ICA.Name,
		Category:             helper.ICA.Category,
		State:                helper.ICA.State,
		Location:             helper.ICA.Location,
		Description:          helper.ICA.Description,
		StartParameterScheme: helper.ICA.StartParameterScheme,
		ManagedExternally:    newFalse(),
	}
	code, resp, err := helper.TestEndpoint(router, token,
		"/api/ic", "POST", helper.KeyModels{"ic": newIC})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Read newIC's ID from the response
	newICID, err := helper.GetResponseID(resp)
	assert.NoError(t, err)

	// authenticate as user
	token, err = helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserACredentials)
	assert.NoError(t, err)

	// Test PUT IC
	// This should fail with unprocessable entity status code 422
	newIC.WebsocketURL = "ThisIsMyNewURL"
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/ic/%v", newICID), "PUT", helper.KeyModels{"ic": newIC})
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)

}

func TestDeleteICAsAdmin(t *testing.T) {
	database.DropTables()
	database.MigrateModels()
	assert.NoError(t, helper.DBAddAdminAndUserAndGuest())

	// authenticate as admin
	token, err := helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.AdminCredentials)
	assert.NoError(t, err)

	// test POST ic/ $newIC
	newIC := ICRequest{
		UUID:                 helper.ICA.UUID,
		WebsocketURL:         helper.ICA.WebsocketURL,
		Type:                 helper.ICA.Type,
		Name:                 helper.ICA.Name,
		Category:             helper.ICA.Category,
		State:                helper.ICA.State,
		Location:             helper.ICA.Location,
		Description:          helper.ICA.Description,
		StartParameterScheme: helper.ICA.StartParameterScheme,
		ManagedExternally:    newFalse(),
	}
	code, resp, err := helper.TestEndpoint(router, token,
		"/api/ic", "POST", helper.KeyModels{"ic": newIC})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Read newIC's ID from the response
	newICID, err := helper.GetResponseID(resp)
	assert.NoError(t, err)

	// Count the number of all the ICs returned for admin
	initialNumber, err := helper.LengthOfResponse(router, token,
		"/api/ic", "GET", nil)
	assert.NoError(t, err)

	// Delete the added newIC
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/ic/%v", newICID), "DELETE", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare DELETE's response with the newIC
	err = helper.CompareResponse(resp, helper.KeyModels{"ic": newIC})
	assert.NoError(t, err)

	// Again count the number of all the ICs returned
	finalNumber, err := helper.LengthOfResponse(router, token,
		"/api/ic", "GET", nil)
	assert.NoError(t, err)

	assert.Equal(t, finalNumber, initialNumber-1)

	// fake an IC update (create) message
	var update ICUpdate
	update.Status = new(ICStatus)
	update.Status.State = new(string)
	*update.Status.State = "idle"
	update.Status.Name = new(string)
	*update.Status.Name = helper.ICB.Name
	update.Status.Category = new(string)
	*update.Status.Category = helper.ICB.Category
	update.Status.Type = new(string)
	*update.Status.Type = helper.ICB.Type

	payload, err := json.Marshal(update)
	assert.NoError(t, err)

	var headers map[string]interface{}
	headers = make(map[string]interface{}) // empty map
	headers["uuid"] = helper.ICB.UUID      // set uuid

	msg := amqp.Publishing{
		DeliveryMode:    2,
		Timestamp:       time.Now(),
		ContentType:     "application/json",
		ContentEncoding: "utf-8",
		Priority:        0,
		Body:            payload,
		Headers:         headers,
	}

	err = CheckConnection()
	assert.NoError(t, err)

	err = client.channel.Publish(VILLAS_EXCHANGE,
		"",
		false,
		false,
		msg)
	assert.NoError(t, err)

	// Wait until externally managed IC is created (happens async)
	time.Sleep(waitingTime * time.Second)

	// Delete the added external IC (triggers an AMQP message, but should not remove the IC from the DB)
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/ic/%v", 2), "DELETE", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Again count the number of all the ICs returned
	finalNumberAfterExtneralDelete, err := helper.LengthOfResponse(router, token,
		"/api/ic", "GET", nil)
	assert.NoError(t, err)

	assert.Equal(t, finalNumber+1, finalNumberAfterExtneralDelete)

}

func TestDeleteICAsUser(t *testing.T) {
	database.DropTables()
	database.MigrateModels()
	assert.NoError(t, helper.DBAddAdminAndUserAndGuest())

	// authenticate as admin
	token, err := helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.AdminCredentials)
	assert.NoError(t, err)

	// test POST ic/ $newIC
	newIC := ICRequest{
		UUID:                 helper.ICA.UUID,
		WebsocketURL:         helper.ICA.WebsocketURL,
		Type:                 helper.ICA.Type,
		Name:                 helper.ICA.Name,
		Category:             helper.ICA.Category,
		State:                helper.ICA.State,
		Location:             helper.ICA.Location,
		Description:          helper.ICA.Description,
		StartParameterScheme: helper.ICA.StartParameterScheme,
		ManagedExternally:    newFalse(),
	}
	code, resp, err := helper.TestEndpoint(router, token,
		"/api/ic", "POST", helper.KeyModels{"ic": newIC})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Read newIC's ID from the response
	newICID, err := helper.GetResponseID(resp)
	assert.NoError(t, err)

	// authenticate as user
	token, err = helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserACredentials)
	assert.NoError(t, err)

	// Test DELETE ICs
	// This should fail with unprocessable entity status code 422
	newIC.WebsocketURL = "ThisIsMyNewURL"
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/ic/%v", newICID), "DELETE", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 422, code, "Response body: \n%v\n", resp)
}

func TestGetAllICs(t *testing.T) {
	database.DropTables()
	database.MigrateModels()
	assert.NoError(t, helper.DBAddAdminAndUserAndGuest())

	// authenticate as admin
	token, err := helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.AdminCredentials)
	assert.NoError(t, err)

	// get the length of the GET all ICs response for user
	initialNumber, err := helper.LengthOfResponse(router, token,
		"/api/ic", "GET", nil)
	assert.NoError(t, err)

	// test POST ic/ $newICA
	newICA := ICRequest{
		UUID:                 helper.ICA.UUID,
		WebsocketURL:         helper.ICA.WebsocketURL,
		Type:                 helper.ICA.Type,
		Name:                 helper.ICA.Name,
		Category:             helper.ICA.Category,
		State:                helper.ICA.State,
		Location:             helper.ICA.Location,
		Description:          helper.ICA.Description,
		StartParameterScheme: helper.ICA.StartParameterScheme,
		ManagedExternally:    newFalse(),
	}
	code, resp, err := helper.TestEndpoint(router, token,
		"/api/ic", "POST", helper.KeyModels{"ic": newICA})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// test POST ic/ $newICB
	newICB := ICRequest{
		UUID:                 helper.ICB.UUID,
		WebsocketURL:         helper.ICB.WebsocketURL,
		Type:                 helper.ICB.Type,
		Name:                 helper.ICB.Name,
		Category:             helper.ICB.Category,
		State:                helper.ICB.State,
		Location:             helper.ICB.Location,
		Description:          helper.ICB.Description,
		StartParameterScheme: helper.ICB.StartParameterScheme,
		ManagedExternally:    newFalse(),
	}

	code, resp, err = helper.TestEndpoint(router, token,
		"/api/ic", "POST", helper.KeyModels{"ic": newICB})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// get the length of the GET all ICs response again
	finalNumber, err := helper.LengthOfResponse(router, token,
		"/api/ic", "GET", nil)
	assert.NoError(t, err)

	assert.Equal(t, finalNumber, initialNumber+2)

	// authenticate as normal user
	token, err = helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserACredentials)
	assert.NoError(t, err)

	// get the length of the GET all ICs response again
	finalNumber2, err := helper.LengthOfResponse(router, token,
		"/api/ic", "GET", nil)
	assert.NoError(t, err)

	assert.Equal(t, finalNumber2, initialNumber+2)
}

func TestGetConfigsOfIC(t *testing.T) {
	database.DropTables()
	database.MigrateModels()
	assert.NoError(t, helper.DBAddAdminAndUserAndGuest())

	// authenticate as admin
	token, err := helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.AdminCredentials)
	assert.NoError(t, err)

	// test POST ic/ $newICA
	newICA := ICRequest{
		UUID:                 helper.ICA.UUID,
		WebsocketURL:         helper.ICA.WebsocketURL,
		Type:                 helper.ICA.Type,
		Name:                 helper.ICA.Name,
		Category:             helper.ICA.Category,
		State:                helper.ICA.State,
		Location:             helper.ICA.Location,
		Description:          helper.ICA.Description,
		StartParameterScheme: helper.ICA.StartParameterScheme,
		ManagedExternally:    newFalse(),
	}
	code, resp, err := helper.TestEndpoint(router, token,
		"/api/ic", "POST", helper.KeyModels{"ic": newICA})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Read newIC's ID from the response
	newICID, err := helper.GetResponseID(resp)
	assert.NoError(t, err)

	// test GET ic/ID/confis
	numberOfConfigs, err := helper.LengthOfResponse(router, token,
		fmt.Sprintf("/api/ic/%v/configs", newICID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	assert.Equal(t, 0, numberOfConfigs)

	// authenticate as normal user
	token, err = helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.UserACredentials)
	assert.NoError(t, err)

	// test GET ic/ID/configs
	numberOfConfigs, err = helper.LengthOfResponse(router, token,
		fmt.Sprintf("/api/ic/%v/configs", newICID), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	assert.Equal(t, 0, numberOfConfigs)

	// Try to get configs of IC that does not exist
	// should result in not found
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/ic/%v/configs", newICID+1), "GET", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 404, code, "Response body: \n%v\n", resp)
}

func TestSendActionToIC(t *testing.T) {
	database.DropTables()
	database.MigrateModels()
	assert.NoError(t, helper.DBAddAdminAndUserAndGuest())

	// authenticate as admin
	token, err := helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.AdminCredentials)
	assert.NoError(t, err)

	// test POST ic/ $newICA
	newICA := ICRequest{
		UUID:                 helper.ICA.UUID,
		WebsocketURL:         helper.ICA.WebsocketURL,
		Type:                 helper.ICA.Type,
		Name:                 helper.ICA.Name,
		Category:             helper.ICA.Category,
		State:                helper.ICA.State,
		Location:             helper.ICA.Location,
		Description:          helper.ICA.Description,
		StartParameterScheme: helper.ICA.StartParameterScheme,
		ManagedExternally:    newFalse(),
	}
	code, resp, err := helper.TestEndpoint(router, token,
		"/api/ic", "POST", helper.KeyModels{"ic": newICA})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Read newIC's ID from the response
	newICID, err := helper.GetResponseID(resp)
	assert.NoError(t, err)

	// create action to be sent to IC
	action1 := ICAction{
		Act:  "start",
		When: time.Now().Unix(),
	}
	action1.Properties.UUID = new(string)
	*action1.Properties.UUID = newICA.UUID
	actions := [1]ICAction{action1}

	// Send action to IC
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/ic/%v/action", newICID), "POST", actions)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Send malformed actions array to IC (should yield bad request)
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/ic/%v/action", newICID), "POST", action1)
	assert.NoError(t, err)
	assert.Equalf(t, 400, code, "Response body: \n%v\n", resp)
}

func TestCreateUpdateViaAMQPRecv(t *testing.T) {

	database.DropTables()
	database.MigrateModels()
	assert.NoError(t, helper.DBAddAdminAndUserAndGuest())

	// authenticate as admin
	token, err := helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.AdminCredentials)
	assert.NoError(t, err)

	// fake an IC update message
	var update ICUpdate
	update.Status = new(ICStatus)
	update.Status.State = new(string)
	*update.Status.State = "idle"

	payload, err := json.Marshal(update)
	assert.NoError(t, err)

	var headers map[string]interface{}
	headers = make(map[string]interface{}) // empty map
	headers["uuid"] = helper.ICB.UUID      // set uuid

	msg := amqp.Publishing{
		DeliveryMode:    2,
		Timestamp:       time.Now(),
		ContentType:     "application/json",
		ContentEncoding: "utf-8",
		Priority:        0,
		Body:            payload,
		Headers:         headers,
	}

	err = CheckConnection()
	assert.NoError(t, err)

	err = client.channel.Publish(VILLAS_EXCHANGE,
		"",
		false,
		false,
		msg)
	assert.NoError(t, err)

	time.Sleep(waitingTime * time.Second)

	// get the length of the GET all ICs response for user
	number, err := helper.LengthOfResponse(router, token,
		"/api/ic", "GET", nil)
	assert.NoError(t, err)
	assert.Equal(t, 0, number)

	// complete the (required) data of an IC
	update.Status.Name = new(string)
	*update.Status.Name = helper.ICA.Name
	update.Status.Category = new(string)
	*update.Status.Category = helper.ICA.Category
	update.Status.Type = new(string)
	*update.Status.Type = helper.ICA.Type
	update.Status.Uptime = new(float64)
	*update.Status.Uptime = -1.0
	update.Status.WS_url = new(string)
	*update.Status.WS_url = helper.ICA.WebsocketURL
	update.Status.API_url = new(string)
	*update.Status.API_url = helper.ICA.APIURL
	update.Status.Description = new(string)
	*update.Status.Description = helper.ICA.Description
	update.Status.Location = new(string)
	*update.Status.Location = helper.ICA.Location

	payload, err = json.Marshal(update)
	assert.NoError(t, err)

	var headersA map[string]interface{}
	headersA = make(map[string]interface{}) // empty map
	headersA["uuid"] = helper.ICA.UUID

	msg = amqp.Publishing{
		DeliveryMode:    2,
		Timestamp:       time.Now(),
		ContentType:     "application/json",
		ContentEncoding: "utf-8",
		Priority:        0,
		Body:            payload,
		Headers:         headersA,
	}

	err = client.channel.Publish(VILLAS_EXCHANGE,
		"",
		false,
		false,
		msg)
	assert.NoError(t, err)

	time.Sleep(waitingTime * time.Second)

	// get the length of the GET all ICs response for user
	number, err = helper.LengthOfResponse(router, token,
		"/api/ic", "GET", nil)
	assert.NoError(t, err)
	assert.Equal(t, 1, number)

	// modify status update
	*update.Status.Name = "This is the new name"
	payload, err = json.Marshal(update)
	assert.NoError(t, err)

	msg = amqp.Publishing{
		DeliveryMode:    2,
		Timestamp:       time.Now(),
		ContentType:     "application/json",
		ContentEncoding: "utf-8",
		Priority:        0,
		Body:            payload,
		Headers:         headersA,
	}

	err = client.channel.Publish(VILLAS_EXCHANGE,
		"",
		false,
		false,
		msg)
	assert.NoError(t, err)

	time.Sleep(waitingTime * time.Second)
	// get the length of the GET all ICs response for user
	number, err = helper.LengthOfResponse(router, token,
		"/api/ic", "GET", nil)
	assert.NoError(t, err)
	assert.Equal(t, 1, number)

}

func TestDeleteICViaAMQPRecv(t *testing.T) {

	database.DropTables()
	database.MigrateModels()
	assert.NoError(t, helper.DBAddAdminAndUserAndGuest())

	// authenticate as admin
	token, err := helper.AuthenticateForTest(router,
		"/api/authenticate", "POST", helper.AdminCredentials)
	assert.NoError(t, err)

	// fake an IC update message
	var update ICUpdate
	update.Status = new(ICStatus)

	update.Status.State = new(string)
	*update.Status.State = "idle"
	// complete the (required) data of an IC
	update.Status.Name = new(string)
	*update.Status.Name = helper.ICA.Name
	update.Status.Category = new(string)
	*update.Status.Category = helper.ICA.Category
	update.Status.Type = new(string)
	*update.Status.Type = helper.ICA.Type
	update.Status.Uptime = new(float64)
	*update.Status.Uptime = -1.0
	update.Status.WS_url = new(string)
	*update.Status.WS_url = helper.ICA.WebsocketURL
	update.Status.API_url = new(string)
	*update.Status.API_url = helper.ICA.APIURL
	update.Status.Description = new(string)
	*update.Status.Description = helper.ICA.Description
	update.Status.Location = new(string)
	*update.Status.Location = helper.ICA.Location

	payload, err := json.Marshal(update)
	assert.NoError(t, err)

	var headers map[string]interface{}
	headers = make(map[string]interface{}) // empty map
	headers["uuid"] = helper.ICB.UUID      // set uuid

	msg := amqp.Publishing{
		DeliveryMode:    2,
		Timestamp:       time.Now(),
		ContentType:     "application/json",
		ContentEncoding: "utf-8",
		Priority:        0,
		Body:            payload,
		Headers:         headers,
	}

	err = CheckConnection()
	assert.NoError(t, err)
	err = client.channel.Publish(VILLAS_EXCHANGE,
		"",
		false,
		false,
		msg)
	assert.NoError(t, err)

	time.Sleep(waitingTime * time.Second)

	// get the length of the GET all ICs response for user
	number, err := helper.LengthOfResponse(router, token,
		"/api/ic", "GET", nil)
	assert.NoError(t, err)
	assert.Equal(t, 1, number)

	// add scenario
	newScenario := ScenarioRequest{
		Name:            helper.ScenarioA.Name,
		Running:         helper.ScenarioA.Running,
		StartParameters: helper.ScenarioA.StartParameters,
	}

	code, resp, err := helper.TestEndpoint(router, token,
		"/api/scenarios", "POST", helper.KeyModels{"scenario": newScenario})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare POST's response with the newScenario
	err = helper.CompareResponse(resp, helper.KeyModels{"scenario": newScenario})
	assert.NoError(t, err)

	// Read newScenario's ID from the response
	newScenarioID, err := helper.GetResponseID(resp)
	assert.NoError(t, err)

	// Add component config and associate with IC and scenario
	newConfig := ConfigRequest{
		Name:            helper.ConfigA.Name,
		ScenarioID:      uint(newScenarioID),
		ICID:            1,
		StartParameters: helper.ConfigA.StartParameters,
		FileIDs:         helper.ConfigA.FileIDs,
	}

	code, resp, err = helper.TestEndpoint(router, token,
		"/api/configs", "POST", helper.KeyModels{"config": newConfig})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare POST's response with the newConfig
	err = helper.CompareResponse(resp, helper.KeyModels{"config": newConfig})
	assert.NoError(t, err)

	// Read newConfig's ID from the response
	newConfigID, err := helper.GetResponseID(resp)
	assert.NoError(t, err)

	// modify status update to state "gone"
	*update.Status.State = "gone"
	payload, err = json.Marshal(update)
	assert.NoError(t, err)

	msg = amqp.Publishing{
		DeliveryMode:    2,
		Timestamp:       time.Now(),
		ContentType:     "application/json",
		ContentEncoding: "utf-8",
		Priority:        0,
		Body:            payload,
		Headers:         headers,
	}

	// attempt to delete IC (should not work immediately because IC is still associated with component config)
	err = client.channel.Publish(VILLAS_EXCHANGE,
		"",
		false,
		false,
		msg)
	assert.NoError(t, err)

	time.Sleep(waitingTime * time.Second)

	// get the length of the GET all ICs response for user
	number, err = helper.LengthOfResponse(router, token,
		"/api/ic", "GET", nil)
	assert.NoError(t, err)
	assert.Equal(t, 1, number)

	// Delete component config from earlier
	code, resp, err = helper.TestEndpoint(router, token,
		fmt.Sprintf("/api/configs/%v", newConfigID), "DELETE", nil)
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Compare DELETE's response with the newConfig
	err = helper.CompareResponse(resp, helper.KeyModels{"config": newConfig})
	assert.NoError(t, err)

	// get the length of the GET all ICs response for user
	number, err = helper.LengthOfResponse(router, token,
		"/api/ic", "GET", nil)
	assert.NoError(t, err)
	assert.Equal(t, 0, number)
}
