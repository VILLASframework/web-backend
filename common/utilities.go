package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/nsf/jsondiff"
	"github.com/stretchr/testify/assert"
)

const UserIDCtx = "user_id"
const UserRoleCtx = "user_role"

func ProvideErrorResponse(c *gin.Context, err error) bool {
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			errormsg := "Record not Found in DB: " + err.Error()
			c.JSON(http.StatusNotFound, gin.H{
				"success": false,
				"message": errormsg,
			})
		} else {
			errormsg := "Error on DB Query or transaction: " + err.Error()
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": errormsg,
			})
		}
		return true // Error
	}
	return false // No error
}

func GetResponseID(resp *bytes.Buffer) (int, error) {

	// Transform bytes buffer into byte slice
	respBytes := []byte(resp.String())

	// Map JSON response to a map[string]map[string]interface{}
	var respRemapped map[string]map[string]interface{}
	err := json.Unmarshal(respBytes, &respRemapped)
	if err != nil {
		return 0, fmt.Errorf("Unmarshal failed for respRemapped %v", err)
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
			return 0, fmt.Errorf("Cannot type assert respRemapped")
		}
		return int(id), nil
	}
	return 0, fmt.Errorf("GetResponse reached exit")
}

func LengthOfResponse(router *gin.Engine, token string, url string,
	method string, body []byte) (int, error) {

	w := httptest.NewRecorder()

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return 0, fmt.Errorf("Failed to create new request: %v", err)
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
	responseBytes := []byte(w.Body.String())

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
	return 0, fmt.Errorf("Length of response cannot be detected")
}

func NewTestEndpoint(router *gin.Engine, token string, url string,
	method string, requestBody interface{}) (int, *bytes.Buffer, error) {

	w := httptest.NewRecorder()

	// Marshal the HTTP request body
	body, err := json.Marshal(requestBody)
	if err != nil {
		return 0, nil, fmt.Errorf("Failed to marshal request body: %v", err)
	}

	// Create the request
	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return 0, nil, fmt.Errorf("Failed to create new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)
	router.ServeHTTP(w, req)

	return w.Code, w.Body, nil
}

func CompareResponse(resp *bytes.Buffer, expected interface{}) error {
	// Serialize expected response
	expectedBytes, err := json.Marshal(expected)
	if err != nil {
		return fmt.Errorf("Failed to marshal expected response: %v", err)
	}
	// Compare
	opts := jsondiff.DefaultConsoleOptions()
	diff, _ := jsondiff.Compare(resp.Bytes(), expectedBytes, &opts)
	if diff.String() != "FullMatch" && diff.String() != "SupersetMatch" {
		return fmt.Errorf("Response: Expected \"%v\". Got \"%v\".",
			"(FullMatch OR SupersetMatch)", diff.String())
	}

	return nil
}

func TestEndpoint(t *testing.T, router *gin.Engine, token string, url string, method string, body []byte, expected_code int, expected_response []byte) {
	w := httptest.NewRecorder()

	if body != nil {
		req, _ := http.NewRequest(method, url, bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Add("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)
	} else {
		req, _ := http.NewRequest(method, url, nil)
		req.Header.Add("Authorization", "Bearer "+token)
		router.ServeHTTP(w, req)
	}

	assert.Equal(t, expected_code, w.Code)
	//fmt.Println("Actual:", w.Body.String())
	//fmt.Println("Expected: ", string(expected_response))
	opts := jsondiff.DefaultConsoleOptions()
	diff, _ := jsondiff.Compare(w.Body.Bytes(), expected_response, &opts)
	assert.Equal(t, "FullMatch", diff.String())

}

func NewAuthenticateForTest(router *gin.Engine, url string,
	method string, credentials interface{}) (string, error) {

	w := httptest.NewRecorder()

	// Marshal credentials
	body, err := json.Marshal(credentials)
	if err != nil {
		return "", fmt.Errorf("Failed to marshal credentials: %v", err)
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("Faile to create new request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Check that return HTTP Code is 200 (OK)
	if w.Code != http.StatusOK {
		return "", fmt.Errorf("HTTP Code: Expected \"%v\". Got \"%v\".",
			http.StatusOK, w.Code)
	}

	// Get the response
	var body_data map[string]interface{}
	err = json.Unmarshal([]byte(w.Body.String()), &body_data)
	if err != nil {
		return "", err
	}

	// Check the response
	success, ok := body_data["success"].(bool)
	if !ok {
		return "", fmt.Errorf("Type asssertion of response[\"success\"] failed")
	}
	if !success {
		return "", fmt.Errorf("Authentication failed: %v", body_data["message"])
	}

	// Extract the token
	token, ok := body_data["token"].(string)
	if !ok {
		return "", fmt.Errorf("Type assertion of response[\"token\"] failed")
	}

	// Return the token and nil error
	return token, nil
}

func AuthenticateForTest(t *testing.T, router *gin.Engine, url string, method string, body []byte, expected_code int) string {
	w := httptest.NewRecorder()

	req, _ := http.NewRequest(method, url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	assert.Equal(t, expected_code, w.Code)

	var body_data map[string]interface{}

	err := json.Unmarshal([]byte(w.Body.String()), &body_data)
	if err != nil {
		panic(err)
	}

	success := body_data["success"].(bool)
	if !success {
		fmt.Println("Authentication not successful: ", body_data["message"])
		panic(-1)
	}

	fmt.Println(w.Body.String())

	return body_data["token"].(string)
}

// Read the parameter with name paramName from the gin Context and
// return it as uint variable
func UintParamFromCtx(c *gin.Context, paramName string) (uint, error) {

	param, err := strconv.Atoi(c.Param(paramName))

	return uint(param), err
}
