/** Result package, endpoints.
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

package result

import (
	"fmt"
	"net/http"

	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/helper"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/file"
	"github.com/gin-gonic/gin"
)

func RegisterResultEndpoints(r *gin.RouterGroup) {
	r.GET("", getResults)
	r.POST("", addResult)
	r.PUT("/:resultID", updateResult)
	r.GET("/:resultID", getResult)
	r.DELETE("/:resultID", deleteResult)
	r.POST("/:resultID/file", addResultFile)
	r.DELETE("/:resultID/file/:fileID", deleteResultFile)
}

// getResults godoc
// @Summary Get all results of scenario
// @ID getResults
// @Produce  json
// @Tags results
// @Success 200 {object} api.ResponseResults "Results which belong to scenario"
// @Failure 404 {object} api.ResponseError "Not found"
// @Failure 422 {object} api.ResponseError "Unprocessable entity"
// @Failure 500 {object} api.ResponseError "Internal server error"
// @Param scenarioID query int true "Scenario ID"
// @Router /results [get]
// @Security Bearer
func getResults(c *gin.Context) {

	ok, sco := database.CheckScenarioPermissions(c, database.Read, "query", -1)
	if !ok {
		return
	}

	db := database.GetDB()
	var results []database.Result
	err := db.Order("ID asc").Model(sco).Related(&results, "Results").Error
	if !helper.DBError(c, err) {
		c.JSON(http.StatusOK, gin.H{"results": results})
	}
}

// addResult godoc
// @Summary Add a result to a scenario
// @ID addResult
// @Accept json
// @Produce json
// @Tags results
// @Success 200 {object} api.ResponseResult "Result that was added"
// @Failure 400 {object} api.ResponseError "Bad request"
// @Failure 404 {object} api.ResponseError "Not found"
// @Failure 422 {object} api.ResponseError "Unprocessable entity"
// @Failure 500 {object} api.ResponseError "Internal server error"
// @Param inputResult body result.addResultRequest true "Result to be added incl. ID of Scenario"
// @Router /results [post]
// @Security Bearer
func addResult(c *gin.Context) {

	// bind request to JSON
	var req addResultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.BadRequestError(c, err.Error())
		return
	}

	// Validate the request
	if err := req.validate(); err != nil {
		helper.UnprocessableEntityError(c, err.Error())
		return
	}

	// Create the new result from the request
	newResult := req.createResult()

	// Check if user is allowed to modify scenario specified in request
	ok, _ := database.CheckScenarioPermissions(c, database.Update, "body", int(newResult.ScenarioID))
	if !ok {
		return
	}

	// add result to DB and add association to scenario
	err := newResult.addToScenario()
	if !helper.DBError(c, err) {
		c.JSON(http.StatusOK, gin.H{"result": newResult.Result})
	}

}

// updateResult godoc
// @Summary Update a result
// @ID updateResult
// @Tags results
// @Accept json
// @Produce json
// @Success 200 {object} api.ResponseResult "Result that was updated"
// @Failure 400 {object} api.ResponseError "Bad request"
// @Failure 404 {object} api.ResponseError "Not found"
// @Failure 422 {object} api.ResponseError "Unprocessable entity"
// @Failure 500 {object} api.ResponseError "Internal server error"
// @Param inputResult body result.updateResultRequest true "Result to be updated"
// @Param resultID path int true "Result ID"
// @Router /results/{resultID} [put]
// @Security Bearer
func updateResult(c *gin.Context) {

	ok, oldResult_r := database.CheckResultPermissions(c, database.Update, "path", -1)
	if !ok {
		return
	}

	var oldResult Result
	oldResult.Result = oldResult_r

	var req updateResultRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.BadRequestError(c, err.Error())
		return
	}

	// Validate the request
	if err := req.Result.validate(); err != nil {
		helper.BadRequestError(c, err.Error())
		return
	}
	// Create the updatedResult from oldResult
	updatedResult := req.updatedResult(oldResult)

	// update the Result in the DB
	err := oldResult.update(updatedResult)
	if !helper.DBError(c, err) {
		c.JSON(http.StatusOK, gin.H{"result": updatedResult.Result})
	}

}

// getResult godoc
// @Summary Get a Result
// @ID getResult
// @Tags results
// @Produce json
// @Success 200 {object} api.ResponseResult "Result that was requested"
// @Failure 400 {object} api.ResponseError "Bad request"
// @Failure 404 {object} api.ResponseError "Not found"
// @Failure 422 {object} api.ResponseError "Unprocessable entity"
// @Failure 500 {object} api.ResponseError "Internal server error"
// @Param resultID path int true "Result ID"
// @Router /results/{resultID} [get]
// @Security Bearer
func getResult(c *gin.Context) {

	ok, result := database.CheckResultPermissions(c, database.Read, "path", -1)
	if !ok {
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": result})
}

// deleteResult godoc
// @Summary Delete a Result incl. all result files
// @ID deleteResult
// @Tags results
// @Produce json
// @Success 200 {object} api.ResponseResult "Result that was deleted"
// @Failure 400 {object} api.ResponseError "Bad request"
// @Failure 404 {object} api.ResponseError "Not found"
// @Failure 422 {object} api.ResponseError "Unprocessable entity"
// @Failure 500 {object} api.ResponseError "Internal server error"
// @Param resultID path int true "Result ID"
// @Router /results/{resultID} [delete]
// @Security Bearer
func deleteResult(c *gin.Context) {
	ok, result_r := database.CheckResultPermissions(c, database.Delete, "path", -1)
	if !ok {
		return
	}

	var result Result
	result.Result = result_r

	// Check if user is allowed to modify scenario associated with result
	ok, _ = database.CheckScenarioPermissions(c, database.Update, "body", int(result.ScenarioID))
	if !ok {
		return
	}

	err := result.delete()
	if !helper.DBError(c, err) {
		c.JSON(http.StatusOK, gin.H{"result": result.Result})
	}

}

// addResultFile godoc
// @Summary Upload a result file to the DB and associate it with scenario and result
// @ID addResultFile
// @Tags results
// @Accept text/plain
// @Accept text/csv
// @Accept application/gzip
// @Accept application/x-gtar
// @Accept application/x-tar
// @Accept application/x-ustar
// @Accept application/zip
// @Accept application/msexcel
// @Accept application/xml
// @Accept application/x-bag
// @Produce json
// @Success 200 {object} api.ResponseResult "Result that was updated"
// @Failure 400 {object} api.ResponseError "Bad request"
// @Failure 404 {object} api.ResponseError "Not found"
// @Failure 422 {object} api.ResponseError "Unprocessable entity"
// @Failure 500 {object} api.ResponseError "Internal server error"
// @Param inputFile formData file true "File to be uploaded"
// @Param resultID path int true "Result ID"
// @Router /results/{resultID}/file [post]
// @Security Bearer
func addResultFile(c *gin.Context) {
	ok, result_r := database.CheckResultPermissions(c, database.Update, "path", -1)
	if !ok {
		return
	}

	var result Result
	result.Result = result_r

	// Check if user is allowed to modify scenario associated with result
	ok, sco := database.CheckScenarioPermissions(c, database.Update, "body", int(result.ScenarioID))
	if !ok {
		return
	}

	// Extract file from POST request form
	file_header, err := c.FormFile("file")
	if err != nil {
		helper.BadRequestError(c, fmt.Sprintf("Get form error: %s", err.Error()))
		return
	}

	// save result file to DB and associate it with scenario
	var newFile file.File
	err = newFile.Register(file_header, sco.ID)
	if helper.DBError(c, err) {
		return
	}

	// add file ID to ResultFileIDs of Result
	err = result.addResultFileID(newFile.File.ID)
	if !helper.DBError(c, err) {
		c.JSON(http.StatusOK, gin.H{"result": result.Result})
	}

}

// deleteResultFile godoc
// @Summary Delete a result file
// @ID deleteResultFile
// @Tags results
// @Produce json
// @Success 200 {object} api.ResponseResult "Result for which file was deleted"
// @Failure 400 {object} api.ResponseError "Bad request"
// @Failure 404 {object} api.ResponseError "Not found"
// @Failure 422 {object} api.ResponseError "Unprocessable entity"
// @Failure 500 {object} api.ResponseError "Internal server error"
// @Param resultID path int true "Result ID"
// @Param fileID path int true "ID of the file to delete"
// @Router /results/{resultID}/file/{fileID} [delete]
// @Security Bearer
func deleteResultFile(c *gin.Context) {

	// check access
	ok, result_r := database.CheckResultPermissions(c, database.Update, "path", -1)
	if !ok {
		return
	}

	var result Result
	result.Result = result_r

	ok, f_r := database.CheckFilePermissions(c, database.Delete)
	if !ok {
		return
	}

	var f file.File
	f.File = f_r

	// Check if user is allowed to modify scenario associated with result
	ok, _ = database.CheckScenarioPermissions(c, database.Update, "body", int(result.ScenarioID))
	if !ok {
		return
	}

	// remove file ID from ResultFileIDs of Result
	err := result.removeResultFileID(f.ID)
	if helper.DBError(c, err) {
		return
	}

	// Delete the file
	err = f.Delete()
	if !helper.DBError(c, err) {
		c.JSON(http.StatusOK, gin.H{"result": result.Result})
	}

}
