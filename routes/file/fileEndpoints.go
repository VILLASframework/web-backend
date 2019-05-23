package file

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
)

func RegisterFileEndpoints(r *gin.RouterGroup){
	r.GET("/:simulationID/models/:modelID/files", GetFilesOfModel)
	r.POST ("/:simulationID/models/:modelID/file", AddFileToModel)
	//r.POST ("/:simulationID/models/:modelID/file", CloneFileOfModel)
	r.GET("/:simulationID/models/:modelID/file", GetFileOfModel)
	r.PUT("/:simulationID/models/:modelID/file", UpdateFileOfModel)
	r.DELETE("/:simulationID/models/:modelID/file", DeleteFileOfModel)

	r.GET("/:simulationID/visualizations/:visualizationID/widgets/:widgetID/files", GetFilesOfWidget)
	r.POST ("/:simulationID/visualizations/:visualizationID/widgets/:widgetID/file", AddFileToWidget)
	//r.POST ("/:simulationID/visualizations/:visualizationID/widgets/:widgetID/file", CloneFileOfWidget)
	r.GET("/:simulationID/visualizations/:visualizationID/widgets/:widgetID/file", GetFileOfWidget)
	r.PUT("/:simulationID/visualizations/:visualizationID/widgets/:widgetID/file", UpdateFileOfWidget)
	r.DELETE("/:simulationID/visualizations/:visualizationID/widgets/:widgetID/file", DeleteFileOfWidget)
}

// GetFilesOfModel godoc
// @Summary Get all parameters of files of model
// @ID GetFilesOfModel
// @Tags file
// @Success 200 {array} common.File "File parameters requested by user"
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param simulationID path int true "Simulation ID"
// @Param modelID path int true "Model ID"
// @Router simulations/{simulationID}/models/{modelID}/files [get]
func GetFilesOfModel(c *gin.Context) {

	simulationID, modelID, err := getRequestParams(c)
	if err != nil{
		return
	}

	// Find files' properties in DB and return in HTTP response, no change to DB
	allFiles, _, err := FindFiles(c, -1, modelID, simulationID)

	if common.ProvideErrorResponse(c, err) == false {
		serializer := common.FilesSerializerNoAssoc{c, allFiles}
		c.JSON(http.StatusOK, gin.H{
			"files": serializer.Response(),
		})
	}

}

// AddFileToModel godoc
// @Summary Get all parameters of files of model
// @ID AddFileToModel
// @Tags file
// @Success 200 "OK."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param simulationID path int true "Simulation ID"
// @Param modelID path int true "Model ID"
// @Router simulations/{simulationID}/models/{modelID}/file [post]
func AddFileToModel(c *gin.Context) {

	simulationID, modelID, err := getRequestParams(c)
	if err != nil{
		return
	}

	// Save file locally and register file in DB, HTTP response is set by this method
	RegisterFile(c,-1, modelID, simulationID)

}

func CloneFileOfModel(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})

}

func GetFileOfModel(c *gin.Context) {

	simulationID, modelID, err := getRequestParams(c)
	if err != nil{
		return
	}

	// Read file from disk and return in HTTP response, no change to DB
	ReadFile(c, -1, modelID, simulationID)
}

func UpdateFileOfModel(c *gin.Context) {

	simulationID, modelID, err := getRequestParams(c)
	if err != nil{
		return
	}

	// Update file locally and update file entry in DB, HTTP response is set by this method
	UpdateFile(c,-1, modelID, simulationID)
}

func DeleteFileOfModel(c *gin.Context) {

	simulationID, modelID, err := getRequestParams(c)
	if err != nil{
		return
	}

	// Delete file from disk and remove entry from DB, HTTP response is set by this method
	DeleteFile(c, -1, modelID, simulationID)


}

func GetFilesOfWidget(c *gin.Context) {

	simulationID, widgetID, err := getRequestParams(c)
	if err != nil{
		return
	}

	// Find files' properties in DB and return in HTTP response, no change to DB
	allFiles, _, err := FindFiles(c, widgetID, -1, simulationID)

	if common.ProvideErrorResponse(c, err) == false {
		serializer := common.FilesSerializerNoAssoc{c, allFiles}
		c.JSON(http.StatusOK, gin.H{
			"files": serializer.Response(),
		})
	}

}

func AddFileToWidget(c *gin.Context) {

	simulationID, widgetID, err := getRequestParams(c)
	if err != nil{
		return
	}

	// Save file locally and register file in DB, HTTP response is set by this method
	RegisterFile(c,widgetID, -1, simulationID)

}

func CloneFileOfWidget(c *gin.Context) {

	c.JSON(http.StatusOK, gin.H{
		"message": "NOT implemented",
	})

}

func GetFileOfWidget(c *gin.Context) {

	simulationID, widgetID, err := getRequestParams(c)
	if err != nil{
		return
	}

	// Read file from disk and return in HTTP response, no change to DB
	ReadFile(c, widgetID, -1, simulationID)
}

func UpdateFileOfWidget(c *gin.Context) {

	simulationID, widgetID, err := getRequestParams(c)
	if err != nil{
		return
	}

	// Update file locally and update file entry in DB, HTTP response is set by this method
	UpdateFile(c,widgetID, -1, simulationID)
}

func DeleteFileOfWidget(c *gin.Context) {

	simulationID, widgetID, err := getRequestParams(c)
	if err != nil{
		return
	}

	// Delete file from disk and remove entry from DB, HTTP response is set by this method
	DeleteFile(c, widgetID, -1, simulationID)


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
	subID, err = common.GetModelID(c)
	if err != nil{
		subID, err = common.GetWidgetID(c)
		if err != nil {
			return -1, -1, err
		}
	}

	return simulationID, subID, err
}