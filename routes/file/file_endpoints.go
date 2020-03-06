/** File package, endpoints.
*
* @author Sonja Happ <sonja.happ@eonerc.rwth-aachen.de>
* @copyright 2014-2019, Institute for Automation of Complex Power Systems, EONERC
* @license GNU General Public License (version 3)
*
* VILLASweb-backend-go
*
* This program is free software: you can redistribute it and/or modify
* it under the terms of the GNU General Public License as published by
* the Free Software Foundation, either version 3 of the License, or
* any later version.
*
* This program is distributed in the hope that it will be useful,
* but WITHOUT ANY WARRANTY; without even the implied warranty of
* MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
* GNU General Public License for more details.
*
* You should have received a copy of the GNU General Public License
* along with this program.  If not, see <http://www.gnu.org/licenses/>.
*********************************************************************************/
package file

import (
	"fmt"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/helper"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/component-configuration"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/widget"
)

func RegisterFileEndpoints(r *gin.RouterGroup) {
	r.GET("", getFiles)
	r.POST("", addFile)
	r.GET("/:fileID", getFile)
	r.PUT("/:fileID", updateFile)
	r.DELETE("/:fileID", deleteFile)
}

// getFiles godoc
// @Summary Get all files of a specific component configuration or widget
// @ID getFiles
// @Tags files
// @Produce json
// @Success 200 {object} docs.ResponseFiles "Files which belong to config or widget"
// @Failure 404 {object} docs.ResponseError "Not found"
// @Failure 422 {object} docs.ResponseError "Unprocessable entity"
// @Failure 500 {object} docs.ResponseError "Internal server error"
// @Param Authorization header string true "Authorization token"
// @Param objectType query string true "Set to config for files of component configuration, set to widget for files of widget"
// @Param objectID query int true "ID of either config or widget of which files are requested"
// @Router /files [get]
func getFiles(c *gin.Context) {

	var err error
	objectType := c.Request.URL.Query().Get("objectType")
	if objectType != "config" && objectType != "widget" {
		helper.BadRequestError(c, fmt.Sprintf("Object type not supported for files: %s", objectType))
		return
	}
	objectID_s := c.Request.URL.Query().Get("objectID")
	objectID, err := strconv.Atoi(objectID_s)
	if err != nil {
		helper.BadRequestError(c, fmt.Sprintf("Error on ID conversion: %s", err.Error()))
		return
	}

	//Check access
	var ok bool
	var m component_configuration.ComponentConfiguration
	var w widget.Widget
	if objectType == "config" {
		ok, m = component_configuration.CheckPermissions(c, database.Read, "body", objectID)
	} else {
		ok, w = widget.CheckPermissions(c, database.Read, objectID)
	}
	if !ok {
		return
	}

	// get meta data of files
	db := database.GetDB()

	var files []database.File

	if objectType == "config" {
		err = db.Order("ID asc").Model(&m).Related(&files, "Files").Error
	} else {
		err = db.Order("ID asc").Model(&w).Related(&files, "Files").Error
	}

	if !helper.DBError(c, err) {
		c.JSON(http.StatusOK, gin.H{"files": files})
	}

}

// addFile godoc
// @Summary Add a file to a specific component config or widget
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
// @Param Authorization header string true "Authorization token"
// @Param inputFile formData file true "File to be uploaded"
// @Param objectType query string true "Set to config for files of component config, set to widget for files of widget"
// @Param objectID query int true "ID of either config or widget of which files are requested"
// @Router /files [post]
func addFile(c *gin.Context) {

	objectType := c.Request.URL.Query().Get("objectType")
	if objectType != "config" && objectType != "widget" {
		helper.BadRequestError(c, fmt.Sprintf("Object type not supported for files: %s", objectType))
		return
	}
	objectID_s := c.Request.URL.Query().Get("objectID")
	objectID, err := strconv.Atoi(objectID_s)
	if err != nil {
		helper.BadRequestError(c, fmt.Sprintf("Error on ID conversion: %s", err.Error()))
		return
	}

	// Check access
	var ok bool
	if objectType == "config" {
		ok, _ = component_configuration.CheckPermissions(c, database.Update, "body", objectID)
		if !ok {
			return
		}
	} else {
		ok, _ = widget.CheckPermissions(c, database.Update, objectID)
		if !ok {
			return
		}
	}

	// Extract file from POST request form
	file_header, err := c.FormFile("file")
	if err != nil {
		helper.BadRequestError(c, fmt.Sprintf("Get form error: %s", err.Error()))
		return
	}

	var newFile File
	err = newFile.register(file_header, objectType, uint(objectID))
	if !helper.DBError(c, err) {
		c.JSON(http.StatusOK, gin.H{"file": newFile.File})
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
// @Success 200 {object} docs.ResponseFile "File that was requested"
// @Failure 400 {object} docs.ResponseError "Bad request"
// @Failure 404 {object} docs.ResponseError "Not found"
// @Failure 422 {object} docs.ResponseError "Unprocessable entity"
// @Failure 500 {object} docs.ResponseError "Internal server error"
// @Param Authorization header string true "Authorization token"
// @Param fileID path int true "ID of the file to download"
// @Router /files/{fileID} [get]
func getFile(c *gin.Context) {

	// check access
	ok, f := checkPermissions(c, database.Read)
	if !ok {
		return
	}

	err := f.download(c)
	helper.DBError(c, err)
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
// @Param Authorization header string true "Authorization token"
// @Param inputFile formData file true "File to be uploaded"
// @Param fileID path int true "ID of the file to update"
// @Router /files/{fileID} [put]
func updateFile(c *gin.Context) {

	// check access
	ok, f := checkPermissions(c, database.Update)
	if !ok {
		return
	}

	// Extract file from PUT request form
	fileHeader, err := c.FormFile("file")
	if err != nil {
		helper.BadRequestError(c, fmt.Sprintf("Get form error: %s", err.Error()))
		return
	}

	err = f.update(fileHeader)
	if !helper.DBError(c, err) {
		c.JSON(http.StatusOK, gin.H{"file": f.File})
	}
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
// @Param Authorization header string true "Authorization token"
// @Param fileID path int true "ID of the file to update"
// @Router /files/{fileID} [delete]
func deleteFile(c *gin.Context) {

	// check access
	ok, f := checkPermissions(c, database.Delete)
	if !ok {
		return
	}

	err := f.delete()
	if !helper.DBError(c, err) {
		c.JSON(http.StatusOK, gin.H{"file": f.File})
	}

}
