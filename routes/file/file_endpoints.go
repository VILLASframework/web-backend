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
	"net/http"

	"git.rwth-aachen.de/acs/public/villas/web-backend-go/helper"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/scenario"
	"github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
)

func RegisterFileEndpoints(r *gin.RouterGroup) {
	r.GET("", getFiles)
	r.POST("", addFile)
	r.GET("/:fileID", getFile)
	r.PUT("/:fileID", updateFile)
	r.DELETE("/:fileID", deleteFile)
}

// getFiles godoc
// @Summary Get all files of a specific scenario
// @ID getFiles
// @Tags files
// @Produce json
// @Success 200 {object} api.ResponseFiles "Files which belong to scenario"
// @Failure 404 {object} api.ResponseError "Not found"
// @Failure 422 {object} api.ResponseError "Unprocessable entity"
// @Failure 500 {object} api.ResponseError "Internal server error"
// @Param scenarioID query int true "Scenario ID"
// @Router /files [get]
// @Security Bearer
func getFiles(c *gin.Context) {

	ok, so := scenario.CheckPermissions(c, database.Read, "query", -1)
	if !ok {
		return
	}

	// get meta data of files
	db := database.GetDB()
	var files []database.File
	err := db.Order("ID asc").Model(so).Related(&files, "Files").Error
	if !helper.DBError(c, err) {
		c.JSON(http.StatusOK, gin.H{"files": files})
	}

}

// addFile godoc
// @Summary Add a file to a specific scenario
// @ID addFile
// @Tags files
// @Produce json
// @Accept text/plain
// @Accept png
// @Accept jpeg
// @Accept gif
// @Accept model/x-cim
// @Accept model/x-cim.zip
// @Accept multipart/form-data
// @Success 200 {object} api.ResponseFile "File that was added"
// @Failure 400 {object} api.ResponseError "Bad request"
// @Failure 404 {object} api.ResponseError "Not found"
// @Failure 422 {object} api.ResponseError "Unprocessable entity"
// @Failure 500 {object} api.ResponseError "Internal server error"
// @Param inputFile formData file true "File to be uploaded"
// @Param scenarioID query int true "ID of scenario to which file shall be added"
// @Router /files [post]
// @Security Bearer
func addFile(c *gin.Context) {

	ok, so := scenario.CheckPermissions(c, database.Read, "query", -1)
	if !ok {
		return
	}

	// Extract file from POST request form
	fileHeader, err := c.FormFile("file")
	if err != nil {
		helper.BadRequestError(c, fmt.Sprintf("Get form error: %s", err.Error()))
		return
	}

	var newFile File
	err = newFile.Register(fileHeader, so.ID)
	if !helper.DBError(c, err) {
		c.JSON(http.StatusOK, gin.H{"file": newFile.File})
	}

}

// getFile godoc
// @Summary Download a file
// @ID getFile
// @Tags files
// @Produce text/plain
// @Produce text/csv
// @Produce application/gzip
// @Produce application/x-gtar
// @Produce application/x-tar
// @Produce application/x-ustar
// @Produce application/zip
// @Produce application/msexcel
// @Produce application/xml
// @Produce application/x-bag
// @Produce png
// @Produce jpeg
// @Produce gif
// @Produce model/x-cim
// @Produce model/x-cim.zip
// @Success 200 {object} api.ResponseFile "File that was requested"
// @Failure 400 {object} api.ResponseError "Bad request"
// @Failure 404 {object} api.ResponseError "Not found"
// @Failure 422 {object} api.ResponseError "Unprocessable entity"
// @Failure 500 {object} api.ResponseError "Internal server error"
// @Param fileID path int true "ID of the file to download"
// @Router /files/{fileID} [get]
// @Security Bearer
func getFile(c *gin.Context) {

	// check access
	ok, f := CheckPermissions(c, database.Read)
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
// @Accept multipart/form-data
// @Success 200 {object} api.ResponseFile "File that was updated"
// @Failure 400 {object} api.ResponseError "Bad request"
// @Failure 404 {object} api.ResponseError "Not found"
// @Failure 422 {object} api.ResponseError "Unprocessable entity"
// @Failure 500 {object} api.ResponseError "Internal server error"
// @Param inputFile formData file true "File to be uploaded"
// @Param fileID path int true "ID of the file to update"
// @Router /files/{fileID} [put]
// @Security Bearer
func updateFile(c *gin.Context) {

	// check access
	ok, f := CheckPermissions(c, database.Update)
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
// @Success 200 {object} api.ResponseFile "File that was deleted"
// @Failure 400 {object} api.ResponseError "Bad request"
// @Failure 404 {object} api.ResponseError "Not found"
// @Failure 422 {object} api.ResponseError "Unprocessable entity"
// @Failure 500 {object} api.ResponseError "Internal server error"
// @Param fileID path int true "ID of the file to update"
// @Router /files/{fileID} [delete]
// @Security Bearer
func deleteFile(c *gin.Context) {

	// check access
	ok, f := CheckPermissions(c, database.Delete)
	if !ok {
		return
	}

	err := f.Delete()
	if !helper.DBError(c, err) {
		c.JSON(http.StatusOK, gin.H{"file": f.File})
	}

}
