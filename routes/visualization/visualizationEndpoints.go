package visualization

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func VisualizationsRegister(r *gin.RouterGroup) {
	r.GET("/", visualizationsReadEp)
	r.POST("/", visualizationRegistrationEp)
	r.PUT("/:VisualizationID", visualizationUpdateEp)
	r.GET("/:VisualizationID", visualizationReadEp)
	r.DELETE("/:VisualizationID", visualizationDeleteEp)
}

func visualizationsReadEp(c *gin.Context) {
	allVisualizations, _, _ := FindAllVisualizations()
	serializer := VisualizationsSerializer{c, allVisualizations}
	c.JSON(http.StatusOK, gin.H{
		"visualizations": serializer.Response(),
	})
}

func visualizationRegistrationEp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}

func visualizationUpdateEp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}

func visualizationReadEp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}

func visualizationDeleteEp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}
