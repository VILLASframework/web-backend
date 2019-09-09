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
// @Success 200 {object} docs.ResponseFiles "Files which belong to simulation model or widget"
// @Failure 404 {object} docs.ResponseError "Not found"
// @Failure 422 {object} docs.ResponseError "Unprocessable entity"
// @Failure 500 {object} docs.ResponseError "Internal server error"
// @Param objectType query string true "Set to model for files of model, set to widget for files of widget"
// @Param objectID query int true "ID of either model or widget of which files are requested"
// @Router /files [get]
func getFiles(c *gin.Context) {

	objectType := c.Request.URL.Query().Get("objectType")
	if objectType != "model" && objectType != "widget" {
		common.BadRequestError(c, fmt.Sprintf("Object type not supported for files: %s", objectType))
		return
	}
	objectID_s := c.Request.URL.Query().Get("objectID")
	objectID, err := strconv.Atoi(objectID_s)
	if err != nil {
		common.BadRequestError(c, fmt.Sprintf("Error on ID conversion: %s", err.Error()))
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
		if common.DBError(c, err) {
			return
		}
	} else {
		err = db.Order("ID asc").Model(&w).Related(&files, "Files").Error
		if common.DBError(c, err) {
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"files": files})
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
// @Success 200 {object} docs.ResponseFile "File that was added"
// @Failure 400 {object} docs.ResponseError "Bad request"
// @Failure 404 {object} docs.ResponseError "Not found"
// @Failure 422 {object} docs.ResponseError "Unprocessable entity"
// @Failure 500 {object} docs.ResponseError "Internal server error"
// @Param inputFile formData file true "File to be uploaded"
// @Param objectType query string true "Set to model for files of model, set to widget for files of widget"
// @Param objectID query int true "ID of either model or widget of which files are requested"
// @Router /files [post]
func addFile(c *gin.Context) {

	objectType := c.Request.URL.Query().Get("objectType")
	if objectType != "model" && objectType != "widget" {
		common.BadRequestError(c, fmt.Sprintf("Object type not supported for files: %s", objectType))
		return
	}
	objectID_s := c.Request.URL.Query().Get("objectID")
	objectID, err := strconv.Atoi(objectID_s)
	if err != nil {
		common.BadRequestError(c, fmt.Sprintf("Error on ID conversion: %s", err.Error()))
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
		common.BadRequestError(c, fmt.Sprintf("Get form error: %s", err.Error()))
		return
	}

	var newFile File
	err = newFile.register(file_header, objectType, uint(objectID))
	if err != nil {
		common.DBError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"file": newFile.File})
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
// @Success 200 {object} docs.ResponseFile "File that was requested"
// @Failure 400 {object} docs.ResponseError "Bad request"
// @Failure 404 {object} docs.ResponseError "Not found"
// @Failure 422 {object} docs.ResponseError "Unprocessable entity"
// @Failure 500 {object} docs.ResponseError "Internal server error"
// @Param fileID path int true "ID of the file to download"
// @Router /files/{fileID} [get]
func getFile(c *gin.Context) {

	// check access
	ok, f := checkPermissions(c, common.Read)
	if !ok {
		return
	}

	err := f.download(c)
	if err != nil {
		common.DBError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"file": f.File})
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
// @Success 200 {object} docs.ResponseFile "File that was updated"
// @Failure 400 {object} docs.ResponseError "Bad request"
// @Failure 404 {object} docs.ResponseError "Not found"
// @Failure 422 {object} docs.ResponseError "Unprocessable entity"
// @Failure 500 {object} docs.ResponseError "Internal server error"
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
		common.BadRequestError(c, fmt.Sprintf("Get form error: %s", err.Error()))
		return
	}

	file_header, err := c.FormFile("file")
	if err != nil {
		common.BadRequestError(c, fmt.Sprintf("Get form error: %s", err.Error()))
		return
	}

	err = f.update(file_header)
	if err != nil {
		common.DBError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"file": f.File})
}

// deleteFile godoc
// @Summary Delete a file
// @ID deleteFile
// @Tags files
// @Produce json
// @Success 200 {object} docs.ResponseFile "File that was deleted"
// @Failure 400 {object} docs.ResponseError "Bad request"
// @Failure 404 {object} docs.ResponseError "Not found"
// @Failure 422 {object} docs.ResponseError "Unprocessable entity"
// @Failure 500 {object} docs.ResponseError "Internal server error"
// @Param fileID path int true "ID of the file to update"
// @Router /files/{fileID} [delete]
func deleteFile(c *gin.Context) {

	// check access
	ok, f := checkPermissions(c, common.Delete)
	if !ok {
		return
	}

	err := f.delete()
	if err != nil {
		common.DBError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"file": f.File})
}
