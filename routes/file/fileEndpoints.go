package file

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func FilesRegister(r *gin.RouterGroup) {
	r.GET("/", filesReadEp)
	//r.POST("/", fileRegistrationEp) // TODO to be added to API
	//r.PUT("/:FileID", fileUpdateEp) // TODO to be added to API
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

// TODO to be added to API
//func fileRegistrationEp(c *gin.Context) {
//	c.JSON(http.StatusOK, gin.H{
//		"message": "NOT implemented",
//	})
//}

// TODO to be added to API
//func fileUpdateEp(c *gin.Context) {
//	c.JSON(http.StatusOK, gin.H{
//		"message": "NOT implemented",
//	})
//}

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
