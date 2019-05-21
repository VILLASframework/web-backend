package endpoints

import (
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/visualization"
	"github.com/gin-gonic/gin"
	"net/http"
)

func visualizationReadAllEp(c *gin.Context) {
	allVisualizations, _, _ := visualization.FindAllVisualizations()
	serializer := visualization.VisualizationsSerializer{c, allVisualizations}
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
