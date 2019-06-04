package common

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

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

func GetSimulationID(c *gin.Context) (int, error) {

	simID, err := strconv.Atoi(c.Param("simulationID"))

	if err != nil {
		errormsg := fmt.Sprintf("Bad request. No or incorrect format of simulation ID")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errormsg,
		})
		return -1, err
	} else {
		return simID, err

	}
}

func GetModelID(c *gin.Context) (int, error) {

	modelID, err := strconv.Atoi(c.Param("modelID"))

	if err != nil {
		errormsg := fmt.Sprintf("Bad request. No or incorrect format of model ID")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errormsg,
		})
		return -1, err
	} else {
		return modelID, err

	}
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

func ValidateRole(c *gin.Context, model ModelName, action CRUD) error {
	// Extracts and validates the role which is saved in the context for
	// executing a specific CRUD operation on a specific model. In case
	// of invalid role return an error.

	// Get user's role from context
	role, exists := c.Get("user_role")
	if !exists {
		return fmt.Errorf("Request does not contain user's role")
	}

	// Check if the role can execute the action on the model
	if !Roles[role.(string)][model][action] {
		return fmt.Errorf("Action not allowed for role %v", role)
	}

	return nil
}
