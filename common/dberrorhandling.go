package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func ProvideErrorResponse(c *gin.Context, err error) bool {
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "No files found in DB",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Error on DB Query or transaction",
			})
		}
		return true // Error
	}
	return false // No error
}
