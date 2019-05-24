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
	r.GET("/:simulationID/models/:modelID/file", GetFileOfModel)
	r.PUT("/:simulationID/models/:modelID/file", UpdateFileOfModel)
	r.DELETE("/:simulationID/models/:modelID/file", DeleteFileOfModel)
	r.GET("/:simulationID/visualizations/:visualizationID/widgets/:widgetID/files", GetFilesOfWidget)
	r.POST ("/:simulationID/visualizations/:visualizationID/widgets/:widgetID/file", AddFileToWidget)
	r.GET("/:simulationID/visualizations/:visualizationID/widgets/:widgetID/file", GetFileOfWidget)
	r.PUT("/:simulationID/visualizations/:visualizationID/widgets/:widgetID/file", UpdateFileOfWidget)
	r.DELETE("/:simulationID/visualizations/:visualizationID/widgets/:widgetID/file", DeleteFileOfWidget)
}

// GetFilesOfModel godoc
// @Summary Get all parameters of files of model
// @ID GetFilesOfModel
// @Tags files
// @Produce json
// @Success 200 {array} common.FileResponse "File parameters requested by user"
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param simulationID path int true "Simulation ID"
// @Param modelID path int true "Model ID"
// @Router /simulations/{simulationID}/models/{modelID}/files [get]
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
// @Tags files
// @Accept text/plain
// @Produce json
// @Success 200 "OK."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param simulationID path int true "Simulation ID"
// @Param modelID path int true "Model ID"
// @Router /simulations/{simulationID}/models/{modelID}/file [post]
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

// GetFileOfModel godoc
// @Summary Download a file that belongs to a model
// @ID GetFileOfModel
// @Tags files
// @Produce text/plain
// @Success 200 "OK, File included in response."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param simulationID path int true "Simulation ID"
// @Param modelID path int true "Model ID"
// @Router /simulations/{simulationID}/models/{modelID}/file [get]
func GetFileOfModel(c *gin.Context) {

	simulationID, modelID, err := getRequestParams(c)
	if err != nil{
		return
	}

	// Read file from disk and return in HTTP response, no change to DB
	ReadFile(c, -1, modelID, simulationID)
}

// UpdateFileOfModel godoc
// @Summary Update (overwrite) a file that belongs to a model
// @ID UpdateFileOfModel
// @Tags files
// @Accept text/plain
// @Produce json
// @Success 200 "OK."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param simulationID path int true "Simulation ID"
// @Param modelID path int true "Model ID"
// @Router /simulations/{simulationID}/models/{modelID}/file [put]
func UpdateFileOfModel(c *gin.Context) {

	simulationID, modelID, err := getRequestParams(c)
	if err != nil{
		return
	}

	// Update file locally and update file entry in DB, HTTP response is set by this method
	UpdateFile(c,-1, modelID, simulationID)
}

// DeleteFileOfModel godoc
// @Summary Delete a file that belongs to a model
// @ID DeleteFileOfModel
// @Tags files
// @Produce json
// @Success 200 "OK."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param simulationID path int true "Simulation ID"
// @Param modelID path int true "Model ID"
// @Router /simulations/{simulationID}/models/{modelID}/file [delete]
func DeleteFileOfModel(c *gin.Context) {

	simulationID, modelID, err := getRequestParams(c)
	if err != nil{
		return
	}

	// Delete file from disk and remove entry from DB, HTTP response is set by this method
	DeleteFile(c, -1, modelID, simulationID)


}

// GetFilesOfWidget godoc
// @Summary Get all parameters of files of widget
// @ID GetFilesOfWidget
// @Tags files
// @Produce json
// @Success 200 {array} common.WidgetResponse "File parameters requested by user"
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param simulationID path int true "Simulation ID"
// @Param visualizationID path int true "Visualization ID"
// @Param widgetID path int true "Widget ID"
// @Router /simulations/{simulationID}/visualizations/{visualizationID}/widgets/{widgetID}/files [get]
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

// AddFileToWidget godoc
// @Summary Get all parameters of files of widget
// @ID AddFileToWidget
// @Tags files
// @Accept text/plain
// @Produce json
// @Success 200 "OK."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param simulationID path int true "Simulation ID"
// @Param visualizationID path int true "Visualization ID"
// @Param widgetID path int true "Widget ID"
// @Router /simulations/{simulationID}/visualizations/{visualizationID}/widgets/{widgetID}/file [post]
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

// GetFileOfWidget godoc
// @Summary Download a file that belongs to a widget
// @ID GetFileOfWidget
// @Tags files
// @Produce text/plain
// @Success 200 "OK, File included in response."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param simulationID path int true "Simulation ID"
// @Param visualizationID path int true "Visualization ID"
// @Param widgetID path int true "Widget ID"
// @Router /simulations/{simulationID}/visualizations/{visualizationID}/widgets/{widgetID}/file [get]
func GetFileOfWidget(c *gin.Context) {

	simulationID, widgetID, err := getRequestParams(c)
	if err != nil{
		return
	}

	// Read file from disk and return in HTTP response, no change to DB
	ReadFile(c, widgetID, -1, simulationID)
}

// UpdateFileOfWidget godoc
// @Summary Update (overwrite) a file that belongs to a widget
// @ID UpdateFileOfWidget
// @Tags files
// @Accept text/plain
// @Produce json
// @Success 200 "OK."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param simulationID path int true "Simulation ID"
// @Param visualizationID path int true "Visualization ID"
// @Param widgetID path int true "Widget ID"
// @Router /simulations/{simulationID}/visualizations/{visualizationID}/widgets/{widgetID}/file [put]
func UpdateFileOfWidget(c *gin.Context) {

	simulationID, widgetID, err := getRequestParams(c)
	if err != nil{
		return
	}

	// Update file locally and update file entry in DB, HTTP response is set by this method
	UpdateFile(c,widgetID, -1, simulationID)
}

// DeleteFileOfWidget godoc
// @Summary Delete a file that belongs to a widget
// @ID DeleteFileOfWidget
// @Tags files
// @Produce json
// @Success 200 "OK."
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param simulationID path int true "Simulation ID"
// @Param visualizationID path int true "Visualization ID"
// @Param widgetID path int true "Widget ID"
// @Router /simulations/{simulationID}/visualizations/{visualizationID}/widgets/{widgetID}/file [delete]
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