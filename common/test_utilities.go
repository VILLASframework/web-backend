package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/nsf/jsondiff"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/http/httptest"
	"time"
)

// data type used in testing
type KeyModels map[string]interface{}

// #######################################################################
// #################### Data used for testing ############################
// #######################################################################

// Users
var StrPassword0 = "xyz789"
var StrPasswordA = "abc123"
var StrPasswordB = "bcd234"

// Hash passwords with bcrypt algorithm
var bcryptCost = 10
var pw0, _ = bcrypt.GenerateFromPassword([]byte(StrPassword0), bcryptCost)
var pwA, _ = bcrypt.GenerateFromPassword([]byte(StrPasswordA), bcryptCost)
var pwB, _ = bcrypt.GenerateFromPassword([]byte(StrPasswordB), bcryptCost)

var User0 = User{Username: "User_0", Password: string(pw0),
	Role: "Admin", Mail: "User_0@example.com"}
var UserA = User{Username: "User_A", Password: string(pwA),
	Role: "User", Mail: "User_A@example.com"}
var UserB = User{Username: "User_B", Password: string(pwB),
	Role: "User", Mail: "User_B@example.com"}

// Credentials

type Credentials struct {
	Username string
	Password string
}

var AdminCredentials = Credentials{
	Username: User0.Username,
	Password: StrPassword0,
}

var UserACredentials = Credentials{
	Username: UserA.Username,
	Password: StrPasswordA,
}

var UserBCredentials = Credentials{
	Username: UserB.Username,
	Password: StrPasswordB,
}

// Simulators

var propertiesA = json.RawMessage(`{"name" : "TestNameA", "category" : "CategoryA", "location" : "anywhere on earth", "type": "dummy"}`)
var propertiesB = json.RawMessage(`{"name" : "TestNameB", "category" : "CategoryB", "location" : "where ever you want", "type": "generic"}`)

var SimulatorA = Simulator{
	UUID:          "4854af30-325f-44a5-ad59-b67b2597de68",
	Host:          "Host_A",
	Modeltype:     "ModelTypeA",
	Uptime:        0,
	State:         "running",
	StateUpdateAt: time.Now().String(),
	Properties:    postgres.Jsonb{propertiesA},
	RawProperties: postgres.Jsonb{propertiesA},
}

var SimulatorB = Simulator{
	UUID:          "7be0322d-354e-431e-84bd-ae4c9633138b",
	Host:          "Host_B",
	Modeltype:     "ModelTypeB",
	Uptime:        0,
	State:         "idle",
	StateUpdateAt: time.Now().String(),
	Properties:    postgres.Jsonb{propertiesB},
	RawProperties: postgres.Jsonb{propertiesB},
}

// Scenarios

var startParametersA = json.RawMessage(`{"parameter1" : "testValue1A", "parameter2" : "testValue2A", "parameter3" : 42}`)
var startParametersB = json.RawMessage(`{"parameter1" : "testValue1B", "parameter2" : "testValue2B", "parameter3" : 43}`)

var ScenarioA = Scenario{
	Name:            "Scenario_A",
	Running:         true,
	StartParameters: postgres.Jsonb{startParametersA},
}
var ScenarioB = Scenario{
	Name:            "Scenario_B",
	Running:         false,
	StartParameters: postgres.Jsonb{startParametersB},
}

// Simulation Models

var SimulationModelA = SimulationModel{
	Name:            "SimulationModel_A",
	StartParameters: postgres.Jsonb{startParametersA},
}

var SimulationModelB = SimulationModel{
	Name:            "SimulationModel_B",
	StartParameters: postgres.Jsonb{startParametersB},
}

// Signals

var OutSignalA = Signal{
	Name:      "outSignal_A",
	Direction: "out",
	Index:     0,
	Unit:      "V",
}

var OutSignalB = Signal{
	Name:      "outSignal_B",
	Direction: "out",
	Index:     1,
	Unit:      "V",
}

var InSignalA = Signal{
	Name:      "inSignal_A",
	Direction: "in",
	Index:     0,
	Unit:      "A",
}

var InSignalB = Signal{
	Name:      "inSignal_B",
	Direction: "in",
	Index:     1,
	Unit:      "A",
}

// Dashboards

var DashboardA = Dashboard{
	Name: "Dashboard_A",
	Grid: 15,
}
var DashboardB = Dashboard{
	Name: "Dashboard_B",
	Grid: 10,
}

// Files

var FileA = File{
	Name:        "File_A",
	Type:        "text/plain",
	Size:        42,
	ImageHeight: 333,
	ImageWidth:  111,
	Date:        time.Now().String(),
}

var FileB = File{
	Name:        "File_B",
	Type:        "text/plain",
	Size:        1234,
	ImageHeight: 55,
	ImageWidth:  22,
	Date:        time.Now().String(),
}

var FileC = File{
	Name:        "File_C",
	Type:        "text/plain",
	Size:        32,
	ImageHeight: 10,
	ImageWidth:  10,
	Date:        time.Now().String(),
}
var FileD = File{
	Name:        "File_D",
	Type:        "text/plain",
	Size:        5000,
	ImageHeight: 400,
	ImageWidth:  800,
	Date:        time.Now().String(),
}

// Widgets
var customPropertiesA = json.RawMessage(`{"property1" : "testValue1A", "property2" : "testValue2A", "property3" : 42}`)
var customPropertiesB = json.RawMessage(`{"property1" : "testValue1B", "property2" : "testValue2B", "property3" : 43}`)

var WidgetA = Widget{
	Name:             "Widget_A",
	Type:             "graph",
	Width:            100,
	Height:           50,
	MinHeight:        40,
	MinWidth:         80,
	X:                10,
	Y:                10,
	Z:                10,
	IsLocked:         false,
	CustomProperties: postgres.Jsonb{customPropertiesA},
}

var WidgetB = Widget{
	Name:             "Widget_B",
	Type:             "slider",
	Width:            200,
	Height:           20,
	MinHeight:        10,
	MinWidth:         50,
	X:                100,
	Y:                -40,
	Z:                -1,
	IsLocked:         false,
	CustomProperties: postgres.Jsonb{customPropertiesB},
}

// ############################################################################
// #################### Functions used for testing ############################
// ############################################################################

// Return the ID of an element contained in a response
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

// Return the length of an response in case it is an array
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

// Make a request to an endpoint
func TestEndpoint(router *gin.Engine, token string, url string,
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

// Compare the response of a query with a JSON
func CompareResponse(resp *bytes.Buffer, expected interface{}) error {
	// Serialize expected response
	expectedBytes, err := json.Marshal(expected)
	if err != nil {
		return fmt.Errorf("Failed to marshal expected response: %v", err)
	}
	// Compare
	opts := jsondiff.DefaultConsoleOptions()
	diff, text := jsondiff.Compare(resp.Bytes(), expectedBytes, &opts)
	if diff.String() != "FullMatch" && diff.String() != "SupersetMatch" {
		fmt.Println(text)
		return fmt.Errorf("Response: Expected \"%v\". Got \"%v\".",
			"(FullMatch OR SupersetMatch)", diff.String())
	}

	return nil
}

// Authenticate a user for testing purposes
func AuthenticateForTest(router *gin.Engine, url string,
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
