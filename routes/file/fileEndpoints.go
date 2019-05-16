package file

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func FilesRegister(r *gin.RouterGroup) {
	r.GET("/", filesReadEp)
	r.POST("/", fileRegistrationEp) // NEW in API
	r.PUT("/:FileID", fileUpdateEp) // NEW in API
	r.GET("/:FileID", fileReadEp)
	r.DELETE("/:FileID", fileDeleteEp)
}

func filesReadEp(c *gin.Context) {
	allFiles, _, _ := FindAllFiles()
	serializer := FilesSerializerNoAssoc{c, allFiles}
	c.JSON(http.StatusOK, gin.H{
		"files": serializer.Response(),
	})
}

func fileRegistrationEp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}

func fileUpdateEp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}

func fileReadEp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}

func fileDeleteEp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}
