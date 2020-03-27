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
	"fmt"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/helper"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"

	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/web-backend-go/configuration"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/user"
)

var router *gin.Engine

type ICRequest struct {
	UUID       string         `json:"uuid,omitempty"`
	Host       string         `json:"host,omitempty"`
	Type       string         `json:"type,omitempty"`
	Name       string         `json:"name,omitempty"`
	Category   string         `json:"category,omitempty"`
	State      string         `json:"state,omitempty"`
	Properties postgres.Jsonb `json:"properties,omitempty"`
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
	api := router.Group("/api")

	user.RegisterAuthenticate(api.Group("/authenticate"))
	api.Use(user.Authentication(true))
	RegisterICEndpoints(api.Group("/ic"))

	os.Exit(m.Run())
}

func TestAddICAsAdmin(t *testing.T) {
	database.DropTables()
	database.MigrateModels()
	assert.NoError(t, helper.DBAddAdminAndUserAndGuest())

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
		UUID:       helper.ICA.UUID,
		Host:       helper.ICA.Host,
		Type:       helper.ICA.Type,
		Name:       helper.ICA.Name,
		Category:   helper.ICA.Category,
		State:      helper.ICA.State,
		Properties: helper.ICA.Properties,
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
		UUID:       helper.ICA.UUID,
		Host:       helper.ICA.Host,
		Type:       helper.ICA.Type,
		Name:       helper.ICA.Name,
		Category:   helper.ICA.Category,
		State:      helper.ICA.State,
		Properties: helper.ICA.Properties,
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
		UUID:       helper.ICA.UUID,
		Host:       helper.ICA.Host,
		Type:       helper.ICA.Type,
		Name:       helper.ICA.Name,
		Category:   helper.ICA.Category,
		State:      helper.ICA.State,
		Properties: helper.ICA.Properties,
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
	newIC.Host = "ThisIsMyNewHost"
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
		UUID:       helper.ICA.UUID,
		Host:       helper.ICA.Host,
		Type:       helper.ICA.Type,
		Name:       helper.ICA.Name,
		Category:   helper.ICA.Category,
		State:      helper.ICA.State,
		Properties: helper.ICA.Properties,
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
	newIC.Host = "ThisIsMyNewHost"
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
		UUID:       helper.ICA.UUID,
		Host:       helper.ICA.Host,
		Type:       helper.ICA.Type,
		Name:       helper.ICA.Name,
		Category:   helper.ICA.Category,
		State:      helper.ICA.State,
		Properties: helper.ICA.Properties,
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
		UUID:       helper.ICA.UUID,
		Host:       helper.ICA.Host,
		Type:       helper.ICA.Type,
		Name:       helper.ICA.Name,
		Category:   helper.ICA.Category,
		State:      helper.ICA.State,
		Properties: helper.ICA.Properties,
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
	newIC.Host = "ThisIsMyNewHost"
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
		UUID:       helper.ICA.UUID,
		Host:       helper.ICA.Host,
		Type:       helper.ICA.Type,
		Name:       helper.ICA.Name,
		Category:   helper.ICA.Category,
		State:      helper.ICA.State,
		Properties: helper.ICA.Properties,
	}
	code, resp, err := helper.TestEndpoint(router, token,
		"/api/ic", "POST", helper.KeyModels{"ic": newICA})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// test POST ic/ $newICB
	newICB := ICRequest{
		UUID:       helper.ICB.UUID,
		Host:       helper.ICB.Host,
		Type:       helper.ICB.Type,
		Name:       helper.ICB.Name,
		Category:   helper.ICB.Category,
		State:      helper.ICB.State,
		Properties: helper.ICB.Properties,
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
		UUID:       helper.ICA.UUID,
		Host:       helper.ICA.Host,
		Type:       helper.ICA.Type,
		Name:       helper.ICA.Name,
		Category:   helper.ICA.Category,
		State:      helper.ICA.State,
		Properties: helper.ICA.Properties,
	}
	code, resp, err := helper.TestEndpoint(router, token,
		"/api/ic", "POST", helper.KeyModels{"ic": newICA})
	assert.NoError(t, err)
	assert.Equalf(t, 200, code, "Response body: \n%v\n", resp)

	// Read newIC's ID from the response
	newICID, err := helper.GetResponseID(resp)
	assert.NoError(t, err)

	// test GET ic/ID/confis
	// TODO how to properly test this without using component configuration endpoints?
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
	// TODO how to properly test this without using component configuration endpoints?
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