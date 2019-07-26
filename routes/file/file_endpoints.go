package file

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/simulationmodel"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/widget"
)

func RegisterFileEndpoints(r *gin.RouterGroup) {
	r.GET("", getFiles)
	r.POST("", addFile)
	r.GET("/:fileID", getFile)
	r.PUT("/:fileID", updateFile)
	r.DELETE("/:fileID", deleteFile)
}

// getFiles godoc
// @Summary Get all files of a specific model or widget
// @ID getFiles
// @Tags files
// @Produce json
// @Success 200 {array} common.FileResponse "File parameters requested by user"
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param objectType query string true "Set to model for files of model, set to widget for files of widget"
// @Param objectID query int true "ID of either model or widget of which files are requested"
// @Router /files [get]
func getFiles(c *gin.Context) {

	objectType := c.Request.URL.Query().Get("objectType")
	if objectType != "model" && objectType != "widget" {
		errormsg := fmt.Sprintf("Bad request. Object type not supported for files: %s", objectType)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errormsg,
		})
		return
	}
	objectID_s := c.Request.URL.Query().Get("objectID")
	objectID, err := strconv.Atoi(objectID_s)
	if err != nil {
		errormsg := fmt.Sprintf("Bad request. Error on ID conversion: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errormsg,
		})
		return
	}

	//Check access
	var ok bool
	var m simulationmodel.SimulationModel
	var w widget.Widget
	if objectType == "model" {
		ok, m = simulationmodel.CheckPermissions(c, common.Read, "body", objectID)
		if !ok {
			return
		}
	} else {
		ok, w = widget.CheckPermissions(c, common.Read, objectID)
		if !ok {
			return
		}
	}

	// get meta data of files
	db := common.GetDB()

	var files []common.File
	if objectType == "model" {
		err = db.Order("ID asc").Model(&m).Related(&files, "Files").Error
		if common.ProvideErrorResponse(c, err) {
			return
		}
	} else {
		err = db.Order("ID asc").Model(&w).Related(&files, "Files").Error
		if common.ProvideErrorResponse(c, err) {
			return
		}
	}

	serializer := common.FilesSerializerNoAssoc{c, files}
	c.JSON(http.StatusOK, gin.H{
		"files": serializer.Response(),
	})

}

// addFile godoc
// @Summary Add a file to a specific model or widget
// @ID addFile
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
// @Param objectType query string true "Set to model for files of model, set to widget for files of widget"
// @Param objectID query int true "ID of either model or widget of which files are requested"
// @Router /files [post]
func addFile(c *gin.Context) {

	objectType := c.Request.URL.Query().Get("objectType")
	if objectType != "model" && objectType != "widget" {
		errormsg := fmt.Sprintf("Bad request. Object type not supported for files: %s", objectType)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errormsg,
		})
		return
	}
	objectID_s := c.Request.URL.Query().Get("objectID")
	objectID, err := strconv.Atoi(objectID_s)
	if err != nil {
		errormsg := fmt.Sprintf("Bad request. Error on ID conversion: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errormsg,
		})
		return
	}

	// Check access
	var ok bool
	if objectType == "model" {
		ok, _ = simulationmodel.CheckPermissions(c, common.Create, "body", objectID)
		if !ok {
			return
		}
	} else {
		ok, _ = widget.CheckPermissions(c, common.Create, objectID)
		if !ok {
			return
		}
	}

	// Extract file from POST request form
	file_header, err := c.FormFile("file")
	if err != nil {
		errormsg := fmt.Sprintf("Bad request. Get form error: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errormsg,
		})
		return
	}

	var newFile File
	err = newFile.register(file_header, objectType, uint(objectID))
	if common.ProvideErrorResponse(c, err) == false {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK.",
		})
	}
}

// getFile godoc
// @Summary Download a file
// @ID getFile
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
func getFile(c *gin.Context) {

	// check access
	ok, f := checkPermissions(c, common.Read)
	if !ok {
		return
	}

	err := f.download(c)
	if common.ProvideErrorResponse(c, err) {
		return
	}
}

// updateFile godoc
// @Summary Update a file
// @ID updateFile
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
// @Param fileID path int true "ID of the file to update"
// @Router /files/{fileID} [put]
func updateFile(c *gin.Context) {

	// check access
	ok, f := checkPermissions(c, common.Update)
	if !ok {
		return
	}

	// Extract file from PUT request form
	err := c.Request.ParseForm()
	if err != nil {
		errormsg := fmt.Sprintf("Bad request. Get form error: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errormsg,
		})
		return
	}

	file_header, err := c.FormFile("file")
	if err != nil {
		errormsg := fmt.Sprintf("Bad request. Get form error: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errormsg,
		})
		return
	}

	err = f.update(file_header)
	if common.ProvideErrorResponse(c, err) == false {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK.",
		})
	}
}

// deleteFile godoc
// @Summary Delete a file
// @ID deleteFile
// @Tags files
// @Produce json
// @Success 200 "OK"
// @Failure 401 "Unauthorized Access"
// @Failure 403 "Access forbidden."
// @Failure 404 "Not found"
// @Failure 500 "Internal server error"
// @Param fileID path int true "ID of the file to update"
// @Router /files/{fileID} [delete]
func deleteFile(c *gin.Context) {

	// check access
	ok, f := checkPermissions(c, common.Delete)
	if !ok {
		return
	}

	err := f.delete()
	if common.ProvideErrorResponse(c, err) == false {
		c.JSON(http.StatusOK, gin.H{
			"message": "OK.",
		})
	}
}
