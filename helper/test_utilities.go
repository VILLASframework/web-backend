/**
* This file is part of VILLASweb-backend-go
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

package helper

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/gin-gonic/gin"
	"github.com/nsf/jsondiff"
)

// data type used in testing
type KeyModels map[string]interface{}

// #######################################################################
// #################### User data used for testing #######################
// #######################################################################

type UserRequest struct {
	Username    string `json:"username,omitempty"`
	Password    string `json:"password,omitempty"`
	OldPassword string `json:"oldPassword,omitempty"`
	Mail        string `json:"mail,omitempty"`
	Role        string `json:"role,omitempty"`
	Active      string `json:"active,omitempty"`
}

// ############################################################################
// #################### Functions used for testing ############################
// ############################################################################

// Return the ID of an element contained in a response
func GetResponseID(resp *bytes.Buffer) (int, error) {

	// Transform bytes buffer into byte slice
	respBytes := resp.Bytes()

	// Map JSON response to a map[string]map[string]interface{}
	var respRemapped map[string]map[string]interface{}
	err := json.Unmarshal(respBytes, &respRemapped)
	if err != nil {
		return 0, fmt.Errorf("failed to unmarshal for respRemapped %v", err)
	}

	// Get an arbitrary key from tha map. The only key (entry) of
	// course is the model's name. With that trick we do not have to
	// pass the higher level key as argument.
	for arbitrary_key := range respRemapped {

		// The marshaler turns numerical values into float64 types so we
		// first have to make a type assertion to the interface and then
		// the conversion to integer before returning
		id, ok := respRemapped[arbitrary_key]["id"].(float64)
		if !ok {
			return 0, fmt.Errorf("cannot type assert respRemapped")
		}
		return int(id), nil
	}
	return 0, fmt.Errorf("getResponse reached exit")
}

// Return the length of an response in case it is an array
func LengthOfResponse(router *gin.Engine, token string, url string,
	method string, body []byte) (int, error) {

	w := httptest.NewRecorder()

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	// HTTP Code of response must be 200
	if w.Code != 200 {
		return 0, fmt.Errorf("HTTP Code: Expected \"200\". Got \"%v\""+
			".\nResponse message:\n%v", w.Code, w.Body.String())
	}

	// Convert the response in array of bytes
	responseBytes := w.Body.Bytes()

	// First we are trying to unmarshal the response into an array of
	// general type variables ([]interface{}). If this fails we will try
	// to unmarshal into a single general type variable (interface{}).
	// If that also fails we will return 0.

	// Response might be array of objects
	var arrayResponse map[string][]interface{}
	err = json.Unmarshal(responseBytes, &arrayResponse)
	if err == nil {

		// Get an arbitrary key from tha map. The only key (entry) of
		// course is the model's name. With that trick we do not have to
		// pass the higher level key as argument.
		for arbitrary_key := range arrayResponse {
			return len(arrayResponse[arbitrary_key]), nil
		}
	}

	// Response might be a single object
	var singleResponse map[string]interface{}
	err = json.Unmarshal(responseBytes, &singleResponse)
	if err == nil {
		return 1, nil
	}

	// Failed to identify response.
	return 0, fmt.Errorf("length of response cannot be detected")
}

// Make a request to an endpoint
func TestEndpoint(router *gin.Engine, token string, url string,
	method string, requestBody interface{}) (int, *bytes.Buffer, error) {

	w := httptest.NewRecorder()

	// Marshal the HTTP request body
	body, err := json.Marshal(requestBody)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	// Create the request
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return 0, nil, fmt.Errorf("failed to create new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	router.ServeHTTP(w, req)

	return handleRedirect(w, req)
}

func handleRedirect(w *httptest.ResponseRecorder, req *http.Request) (int, *bytes.Buffer, error) {
	if w.Code == http.StatusMovedPermanently ||
		w.Code == http.StatusFound ||
		w.Code == http.StatusTemporaryRedirect ||
		w.Code == http.StatusPermanentRedirect {

		// Follow external redirect
		redirURL, err := w.Result().Location()
		if err != nil {
			return 0, nil, fmt.Errorf("invalid location header")
		}

		log.Println("redirecting request to", redirURL.String())

		req, err := http.NewRequest(req.Method, redirURL.String(), req.Body)
		if err != nil {
			return 0, nil, fmt.Errorf("handle redirect: failed to create new request: %v", err)
		}

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return 0, nil, fmt.Errorf("handle redirect: failed to follow redirect: %v", err)
		}

		buf := new(bytes.Buffer)
		_, err = buf.ReadFrom(resp.Body)
		if err != nil {
			return 0, nil, fmt.Errorf("handle redirect: failed to follow redirect: %v", err)
		}

		return resp.StatusCode, buf, nil
	}

	// No redirect
	return w.Code, w.Body, nil
}

// Compare the response of a query with a JSON
func CompareResponse(resp *bytes.Buffer, expected interface{}) error {
	// Serialize expected response
	expectedBytes, err := json.Marshal(expected)
	if err != nil {
		return fmt.Errorf("failed to marshal expected response: %v", err)
	}
	// Compare
	opts := jsondiff.DefaultConsoleOptions()
	diff, text := jsondiff.Compare(resp.Bytes(), expectedBytes, &opts)
	if diff.String() != "FullMatch" && diff.String() != "SupersetMatch" {
		log.Println(text)
		return fmt.Errorf("response: Expected \"%v\". Got \"%v\"",
			"(FullMatch OR SupersetMatch)", diff.String())
	}

	return nil
}

// Authenticate a user for testing purposes
func AuthenticateForTest(router *gin.Engine, credentials interface{}) (string, error) {

	w := httptest.NewRecorder()

	// Marshal credentials
	body, err := json.Marshal(credentials)
	if err != nil {
		return "", fmt.Errorf("failed to marshal credentials: %v", err)
	}

	req, err := http.NewRequest("POST", "/api/v2/authenticate/internal", bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("failed to create new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Check that return HTTP Code is 200 (OK)
	if w.Code != http.StatusOK {
		return "", fmt.Errorf("http code: Expected \"%v\". Got \"%v\"",
			http.StatusOK, w.Code)
	}

	// Get the response
	var body_data map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &body_data)
	if err != nil {
		return "", err
	}

	// Check the response
	success, ok := body_data["success"].(bool)
	if !ok {
		return "", fmt.Errorf("type asssertion of response[\"success\"] failed")
	}
	if !success {
		return "", fmt.Errorf("authentication failed: %v", body_data["message"])
	}

	// Extract the token
	token, ok := body_data["token"].(string)
	if !ok {
		return "", fmt.Errorf("type assertion of response[\"token\"] failed")
	}

	// Return the token and nil error
	return token, nil
}
