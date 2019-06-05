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

func GetVisualizationID(c *gin.Context) (int, error) {

	simID, err := strconv.Atoi(c.Param("visualizationID"))

	if err != nil {
		errormsg := fmt.Sprintf("Bad request. No or incorrect format of visualization ID")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errormsg,
		})
		return -1, err
	} else {
		return simID, err

	}
}

func GetWidgetID(c *gin.Context) (int, error) {

	widgetID, err := strconv.Atoi(c.Param("widgetID"))

	if err != nil {
		errormsg := fmt.Sprintf("Bad request. No or incorrect format of widget ID")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errormsg,
		})
		return -1, err
	} else {
		return widgetID, err
	}
}

func GetFileID(c *gin.Context) (int, error) {

	fileID, err := strconv.Atoi(c.Param("fileID"))

	if err != nil {
		errormsg := fmt.Sprintf("Bad request. No or incorrect format of file ID")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errormsg,
		})
		return -1, err
	} else {
		return fileID, err
	}
}

func TestEndpoint(t *testing.T, router *gin.Engine, token string, url string, method string, body []byte, expected_code int, expected_response string) {
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
	fmt.Println(w.Body.String())
	assert.Equal(t, expected_response, w.Body.String())
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
		panic(-1)
	}

	fmt.Println(w.Body.String())

	return body_data["token"].(string)
}
