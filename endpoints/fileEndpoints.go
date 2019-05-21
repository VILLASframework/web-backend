package endpoints

import (
	"fmt"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/file"
	"strconv"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Endpoint functions

func fileMReadAllEp(c *gin.Context) {

	simulationID, modelID, err := getRequestParams(c)
	if err != nil{
		return
	}

	// Find files' properties in DB and return in HTTP response, no change to DB
	allFiles, _, err := file.FindFiles(c, -1, modelID, simulationID)

	if common.ProvideErrorResponse(c, err) == false {
		serializer := file.FilesSerializerNoAssoc{c, allFiles}
		c.JSON(http.StatusOK, gin.H{
			"files": serializer.Response(),
		})
	}

}

func fileMRegistrationEp(c *gin.Context) {

	simulationID, modelID, err := getRequestParams(c)
	if err != nil{
		return
	}

	// Save file locally and register file in DB, HTTP response is set by this method
	file.RegisterFile(c,-1, modelID, simulationID)

}

func fileMReadEp(c *gin.Context) {

	simulationID, modelID, err := getRequestParams(c)
	if err != nil{
		return
	}

	// Read file from disk and return in HTTP response, no change to DB
	file.ReadFile(c, -1, modelID, simulationID)
}

func fileMUpdateEp(c *gin.Context) {

	simulationID, modelID, err := getRequestParams(c)
	if err != nil{
		return
	}

	// Update file locally and update file entry in DB, HTTP response is set by this method
	file.UpdateFile(c,-1, modelID, simulationID)
}

func fileMDeleteEp(c *gin.Context) {

	simulationID, modelID, err := getRequestParams(c)
	if err != nil{
		return
	}

	// Delete file from disk and remove entry from DB, HTTP response is set by this method
	file.DeleteFile(c, -1, modelID, simulationID)


}

func fileWReadAllEp(c *gin.Context) {

	simulationID, widgetID, err := getRequestParams(c)
	if err != nil{
		return
	}

	// Find files' properties in DB and return in HTTP response, no change to DB
	allFiles, _, err := file.FindFiles(c, widgetID, -1, simulationID)

	if common.ProvideErrorResponse(c, err) == false {
		serializer := file.FilesSerializerNoAssoc{c, allFiles}
		c.JSON(http.StatusOK, gin.H{
			"files": serializer.Response(),
		})
	}

}

func fileWRegistrationEp(c *gin.Context) {

	simulationID, widgetID, err := getRequestParams(c)
	if err != nil{
		return
	}

	// Save file locally and register file in DB, HTTP response is set by this method
	file.RegisterFile(c,widgetID, -1, simulationID)

}

func fileWReadEp(c *gin.Context) {

	simulationID, widgetID, err := getRequestParams(c)
	if err != nil{
		return
	}

	// Read file from disk and return in HTTP response, no change to DB
	file.ReadFile(c, widgetID, -1, simulationID)
}

func fileWUpdateEp(c *gin.Context) {

	simulationID, widgetID, err := getRequestParams(c)
	if err != nil{
		return
	}

	// Update file locally and update file entry in DB, HTTP response is set by this method
	file.UpdateFile(c,widgetID, -1, simulationID)
}

func fileWDeleteEp(c *gin.Context) {

	simulationID, widgetID, err := getRequestParams(c)
	if err != nil{
		return
	}

	// Delete file from disk and remove entry from DB, HTTP response is set by this method
	file.DeleteFile(c, widgetID, -1, simulationID)


}


// local functions

//func filesReadEp(c *gin.Context)  {
//	// Database query
//	allFiles, _, err := FindAllFiles()
//
//	if common.ProvideErrorResponse(c, err) == false {
//		serializer := FilesSerializerNoAssoc{c, allFiles}
//		c.JSON(http.StatusOK, gin.H{
//			"files": serializer.Response(),
//		})
//	}
//
//}
//
//
//
//func fileUpdateEp(c *gin.Context) {
//	c.JSON(http.StatusOK, gin.H{
//		"message": "NOT implemented",
//	})
//}
//
//func fileReadEp(c *gin.Context) {
//	var err error
//	var file common.File
//	fileID := c.Param("FileID")
//	desc := c.GetHeader("X-Request-FileDesc")
//	desc_b, _ := strconv.ParseBool(desc)
//
//	userID := 1 // TODO obtain ID of user making the request
//
//	//check if description of file or file itself shall be returned
//	if desc_b {
//		file, err = FindFile(userID, fileID)
//		if common.ProvideErrorResponse(c, err) == false {
//			serializer := FileSerializerNoAssoc{c, file}
//			c.JSON(http.StatusOK, gin.H{
//				"file": serializer.Response(),
//			})
//		}
//
//
//	} else {
//		//TODO: return file itself
//	}
//}
//
//func fileDeleteEp(c *gin.Context) {
//	c.JSON(http.StatusOK, gin.H{
//		"message": "NOT implemented",
//	})
//}


func getRequestParams(c *gin.Context) (int, int, error){
	simulationID, err := strconv.Atoi(c.Param("SimulationID"))

	if err != nil {
		errormsg := fmt.Sprintf("Bad request. No or incorrect format of simulation ID")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errormsg,
		})
		return -1, -1, err
	}

	var subID int
	subID, err = GetModelID(c)
	if err != nil{
		subID, err = GetWidgetID(c)
		if err != nil {
			return -1, -1, err
		}
	}

	return simulationID, subID, err
}