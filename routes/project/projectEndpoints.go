package project

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func ProjectsRegister(r *gin.RouterGroup) {
	r.GET("/", projectsReadEp)
	r.POST("/", projectRegistrationEp)
	r.PUT("/:ProjectID", projectUpdateEp)
	r.GET("/:ProjectID", projectReadEp)
	r.DELETE("/:ProjectID", projectDeleteEp)
}

func projectsReadEp(c *gin.Context) {
	allProjects, _, _ := FindAllProjects()
	serializer := ProjectsSerializerNoAssoc{c, allProjects}
	c.JSON(http.StatusOK, gin.H{
		"projects": serializer.Response(),
	})
}

func projectRegistrationEp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}

func projectUpdateEp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}

func projectReadEp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}

func projectDeleteEp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}
