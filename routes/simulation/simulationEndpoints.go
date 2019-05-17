package simulation

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
)

func SimulationsRegister(r *gin.RouterGroup) {
	r.GET("/", simulationsReadEp)
	r.POST("/", simulationRegistrationEp)
	r.PUT("/:SimulationID", simulationUpdateEp)
	r.GET("/:SimulationID", simulationReadEp)
	r.DELETE("/:SimulationID", simulationDeleteEp)
	r.POST ("/:SimulationID/models/:SimulationModelID/file", fileRegistrationEp) // NEW
}

func simulationsReadEp(c *gin.Context) {
	allSimulations, _, _ := FindAllSimulations()
	serializer := SimulationsSerializerNoAssoc{c, allSimulations}
	c.JSON(http.StatusOK, gin.H{
		"simulations": serializer.Response(),
	})
}

func simulationRegistrationEp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}

func simulationUpdateEp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}

func simulationReadEp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}

func simulationDeleteEp(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})
}

func fileRegistrationEp(c *gin.Context) {

	// TODO Check if file upload is ok for this user or simulation (user or simulation exists)
	var widgetID_s = c.Param("WidgetID")
	var widgetID_i int
	var simulationmodelID_s = c.Param("SimulationModelID")
	var simulationmodelID_i int

	if widgetID_s != "" {
		widgetID_i, _ = strconv.Atoi(widgetID_s)
	} else {
		widgetID_i = -1
	}

	if simulationmodelID_s != "" {
		simulationmodelID_i, _ = strconv.Atoi(simulationmodelID_s)
	} else {
		simulationmodelID_i = -1
	}

	if simulationmodelID_i == -1 && widgetID_i == -1 {
		errormsg := fmt.Sprintf("Bad request. Did not provide simulation model ID or widget ID for file")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errormsg,
		})
		return;
	}

	// Extract file from POST request form
	file, err := c.FormFile("file")
	if err != nil {
		errormsg := fmt.Sprintf("Bad request. Get form error: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errormsg,
		})
		return;
	}

	// Obtain properties of file
	filetype := file.Header.Get("Content-Type") // TODO make sure this is properly set in file header
	filename := filepath.Base(file.Filename)
	foldername := "files/testfolder" //TODO replace this placeholder with systematic foldername (e.g. simulation ID)
	size := file.Size

	// Save file to local disc (NOT DB!)
	err = SaveFile(file, filename, foldername, uint(size))
	if err != nil {
		errormsg := fmt.Sprintf("Internal Server Error. Error saving file: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errormsg,
		})
		return
	}

	// Add File object with parameters to DB
	err = AddFile(filename, foldername, filetype, uint(size), widgetID_i, simulationmodelID_i )
	if common.ProvideErrorResponse(c, err) == false {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK.",
		})
	}

}
