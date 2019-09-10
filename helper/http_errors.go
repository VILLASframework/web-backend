package helper

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

func DBError(c *gin.Context, err error) bool {
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			NotFoundError(c, "Record not Found in DB: "+err.Error())
		} else {
			InternalServerError(c, "Error on DB Query or transaction: "+err.Error())
		}
		return true // Error
	}
	return false // No error
}

func BadRequestError(c *gin.Context, err string) {
	c.JSON(http.StatusBadRequest, gin.H{
		"success": false,
		"message": fmt.Sprintf("%v", err),
	})
}

func UnprocessableEntityError(c *gin.Context, err string) {
	c.JSON(http.StatusUnprocessableEntity, gin.H{
		"success": false,
		"message": fmt.Sprintf("%v", err),
	})
}

func InternalServerError(c *gin.Context, err string) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"success": false,
		"message": fmt.Sprintf("%v", err),
	})
}

func UnauthorizedError(c *gin.Context, err string) {
	c.JSON(http.StatusUnauthorized, gin.H{
		"success": false,
		"message": fmt.Sprintf("%v", err),
	})
}

func UnauthorizedAbort(c *gin.Context, err string) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"succes":  false,
		"message": fmt.Sprintf("%v", err),
	})
}

func NotFoundError(c *gin.Context, err string) {
	c.JSON(http.StatusNotFound, gin.H{
		"success": false,
		"message": fmt.Sprintf("%v", err),
	})
}

func ForbiddenError(c *gin.Context, err string) {
	c.JSON(http.StatusForbidden, gin.H{
		"success": false,
		"message": fmt.Sprintf("%v", err),
	})
}
