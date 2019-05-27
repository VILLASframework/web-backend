package file

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
)

func RegisterFileEndpoints(r *gin.RouterGroup){
	r.GET("/", GetFiles)
	r.POST ("/", AddFile)
	r.GET("/:fileID", GetFile)
	r.PUT("/:fileID", UpdateFile)
	r.DELETE("/:fileID", DeleteFile)
	//r.GET("/:simulationID/visualizations/:visualizationID/widgets/:widgetID/files", GetFilesOfWidget)
	//r.POST ("/:simulationID/visualizations/:visualizationID/widgets/:widgetID/file", AddFileToWidget)
	//r.GET("/:simulationID/visualizations/:visualizationID/widgets/:widgetID/file", GetFileOfWidget)
	//r.PUT("/:simulationID/visualizations/:visualizationID/widgets/:widgetID/file", UpdateFileOfWidget)
	//r.DELETE("/:simulationID/visualizations/:visualizationID/widgets/:widgetID/file", DeleteFileOfWidget)
}



// GetFiles godoc
// @Summary Get all files of a specific model or widget
// @ID GetFiles
// @Tags files
// @Produce json
// @Success 200 {array} common.FileResponse "File parameters requested by user"
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param originType query string true "Set to model for files of model, set to widget for files of widget"
// @Param originID query int true "ID of either model or widget of which files are requested"
// @Router /files [get]
func GetFiles(c *gin.Context) {

	// TODO if originType == "model" --> GetFilesOfModel, if originType == "vis" --> GetFilesOfWidget

}

// AddFile godoc
// @Summary Add a file to a specific model or widget
// @ID AddFile
// @Tags files
// @Produce json
// @Accept text/plain
// @Accept png
// @Accept jpeg
// @Accept gif
// @Accept model/x-cim
// @Accept model/x-cim.zip
// @Success 200 "OK"
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param inputFile formData file true "File to be uploaded"
// @Param originType query string true "Set to model for files of model, set to widget for files of widget"
// @Param originID query int true "ID of either model or widget of which files are requested"
// @Router /files [post]
func AddFile(c *gin.Context){
	// TODO if originType == "model" --> AddFileToModel, if originType == "vis" --> AddFileToWidget
}

// GetFile godoc
// @Summary Download a file
// @ID GetFile
// @Tags files
// @Produce text/plain
// @Produce png
// @Produce jpeg
// @Produce gif
// @Produce model/x-cim
// @Produce model/x-cim.zip
// @Success 200 "OK"
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param fileID path int true "ID of the file to download"
// @Router /files/{fileID} [get]
func GetFile(c *gin.Context){
	// TODO
}

// UpdateFile godoc
// @Summary Update a file
// @ID UpdateFile
// @Tags files
// @Produce json
// @Accept text/plain
// @Accept png
// @Accept jpeg
// @Accept gif
// @Accept model/x-cim
// @Accept model/x-cim.zip
// @Success 200 "OK"
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param fileID path int true "ID of the file to update"
// @Router /files/{fileID} [put]
func UpdateFile(c *gin.Context){

	//TODO parse this info based on fileID parameter
	simulationID := 1
	modelID := 1
	widgetID := 1


	// Extract file from PUT request form
	err := c.Request.ParseForm()
	if err != nil {
		errormsg := fmt.Sprintf("Bad request. Get form error: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errormsg,
		})
		return;
	}

	file_header, err := c.FormFile("file")
	if err != nil {
		errormsg := fmt.Sprintf("Bad request. Get form error: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errormsg,
		})
		return;
	}

	filename := filepath.Base(file_header.Filename)
	filetype := file_header.Header.Get("Content-Type") // TODO make sure this is properly set in file header
	size := file_header.Size
	foldername := getFolderName(simulationID, modelID, widgetID)

	err = modifyFileOnDisc(file_header, filename, foldername, uint(size), false)
	if err != nil {
		errormsg := fmt.Sprintf("Internal Server Error. Error saving file: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errormsg,
		})
		return
	}

	saveFileInDB(c, filename, foldername, filetype, uint(size), widgetID, modelID, false)
}


// DeleteFile godoc
// @Summary Delete a file
// @ID DeleteFile
// @Tags files
// @Produce json
// @Success 200 "OK"
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param fileID path int true "ID of the file to update"
// @Router /files/{fileID} [delete]
func DeleteFile(c *gin.Context){
	// TODO
}


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

	//simulationID, modelID, err := getRequestParams(c)
	//if err != nil{
	//	return
	//}

	// Update file locally and update file entry in DB, HTTP response is set by this method
	//UpdateFile(c,-1, modelID, simulationID)
}

func DeleteFileOfModel(c *gin.Context) {

	//simulationID, modelID, err := getRequestParams(c)
	//if err != nil{
	//	return
	//}

	// Delete file from disk and remove entry from DB, HTTP response is set by this method
	//DeleteFile(c, -1, modelID, simulationID)


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

	//simulationID, widgetID, err := getRequestParams(c)
	//if err != nil{
	//	return
	//}
	//
	//// Update file locally and update file entry in DB, HTTP response is set by this method
	//UpdateFile(c,widgetID, -1, simulationID)
}


func DeleteFileOfWidget(c *gin.Context) {

	//simulationID, widgetID, err := getRequestParams(c)
	//if err != nil{
	//	return
	//}
	//
	//// Delete file from disk and remove entry from DB, HTTP response is set by this method
	//DeleteFile(c, widgetID, -1, simulationID)


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