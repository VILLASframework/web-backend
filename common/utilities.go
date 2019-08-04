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
				"error": errormsg,
			})
		} else {
			errormsg := "Error on DB Query or transaction: " + err.Error()
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": errormsg,
			})
		}
		return true // Error
	}
	return false // No error
}

func LengthOfResponse(router *gin.Engine, token string, url string,
	method string, body []byte) int {

	w := httptest.NewRecorder()
	responseLength := 0

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

	// Get the response
	var body_data map[string][]interface{}

	err := json.Unmarshal([]byte(w.Body.String()), &body_data)
	if err != nil {
		return responseLength
	}

	// Get an arbitrary key from tha map. The only key (entry) of course
	// is the model's name. With that trick we do not have to pass the
	// higher level key as argument.
	for arbitrary_tag := range body_data {
		responseLength = len(body_data[arbitrary_tag])
		break
	}

	return responseLength
}

func NewTestEndpoint(router *gin.Engine, token string, url string,
	method string, body []byte, expected_code int,
	expected_response []byte) error {

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

	// Check the return HTTP Code
	if w.Code != expected_code {
		return fmt.Errorf("HTTP Code: Expected \"%v\". Got \"%v\".",
			expected_code, w.Code)
	}

	// Check the response
	opts := jsondiff.DefaultConsoleOptions()
	diff, _ := jsondiff.Compare(w.Body.Bytes(), expected_response, &opts)
	if diff.String() != "FullMatch" {
		return fmt.Errorf("Response: Expected \"%v\". Got \"%v\".",
			"FullMatch", diff.String())
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
	method string, body []byte, expected_code int) (string, error) {

	w := httptest.NewRecorder()

	req, _ := http.NewRequest(method, url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	// Check the return HTTP Code
	if w.Code != expected_code {
		return "", fmt.Errorf("HTTP Code: Expected \"%v\". Got \"%v\".",
			expected_code, w.Code)
	}

	var body_data map[string]interface{}

	// Get the response
	err := json.Unmarshal([]byte(w.Body.String()), &body_data)
	if err != nil {
		return "", err
	}

	// Check the response
	success := body_data["success"].(bool)
	if !success {
		fmt.Println("Authentication not successful: ", body_data["message"])
		return "", fmt.Errorf("Authentication unsuccessful!")
	}

	// Return the token and nil error
	return body_data["token"].(string), nil
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
