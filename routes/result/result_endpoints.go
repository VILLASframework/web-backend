package result

import (
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/helper"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/scenario"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RegisterResultEndpoints(r *gin.RouterGroup) {
	r.GET("", getResults)
	r.POST("", addResult)
	r.PUT("/:resultID", updateResult)
	r.GET("/:resultID", getResult)
	r.DELETE("/:scenarioID", deleteResult)
	r.POST("/:resultID/file", addResultFile)
	r.GET("/:resultID/file/:fileID", getResultFile)
	r.DELETE("/:resultID/file/:fileID", deleteResultFile)
}

// getResults godoc
// @Summary Get all results of scenario
// @ID getResults
// @Produce  json
// @Tags results
// @Success 200 {object} docs.ResponseResults "Results which belong to scenario"
// @Failure 404 {object} docs.ResponseError "Not found"
// @Failure 422 {object} docs.ResponseError "Unprocessable entity"
// @Failure 500 {object} docs.ResponseError "Internal server error"
// @Param scenarioID query int true "Scenario ID"
// @Router /results [get]
// @Security Bearer
func getResults(c *gin.Context) {

	ok, scenario := scenario.CheckPermissions(c, database.Read, "query", -1)
	if !ok {
		return
	}

	db := database.GetDB()
	var result []database.Result
	err := db.Order("ID asc").Model(scenario).Related(&result, "Results").Error
	if !helper.DBError(c, err) {
		c.JSON(http.StatusOK, gin.H{"result": result})
	}
}

// addResult godoc
// @Summary Add a result to a scenario
// @ID addResult
// @Accept json
// @Produce json
// @Tags results
// @Success 200 {object} docs.ResponseResult "Result that was added"
// @Failure 400 {object} docs.ResponseError "Bad request"
// @Failure 404 {object} docs.ResponseError "Not found"
// @Failure 422 {object} docs.ResponseError "Unprocessable entity"
// @Failure 500 {object} docs.ResponseError "Internal server error"
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
	ok, _ := scenario.CheckPermissions(c, database.Update, "body", int(newResult.ScenarioID))
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
// @Success 200 {object} docs.ResponseResult "Result that was updated"
// @Failure 400 {object} docs.ResponseError "Bad request"
// @Failure 404 {object} docs.ResponseError "Not found"
// @Failure 422 {object} docs.ResponseError "Unprocessable entity"
// @Failure 500 {object} docs.ResponseError "Internal server error"
// @Param inputResult body result.updateResultRequest true "Result to be updated"
// @Param resultID path int true "Result ID"
// @Router /results/{resultID} [put]
// @Security Bearer
func updateResult(c *gin.Context) {

	ok, oldResult := CheckPermissions(c, database.Update, "path", -1)
	if !ok {
		return
	}

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
// @Success 200 {object} docs.ResponseResult "Result that was requested"
// @Failure 400 {object} docs.ResponseError "Bad request"
// @Failure 404 {object} docs.ResponseError "Not found"
// @Failure 422 {object} docs.ResponseError "Unprocessable entity"
// @Failure 500 {object} docs.ResponseError "Internal server error"
// @Param resultID path int true "Result ID"
// @Router /results/{resultID} [get]
// @Security Bearer
func getResult(c *gin.Context) {

	ok, result := CheckPermissions(c, database.Read, "path", -1)
	if !ok {
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": result.Result})
}

// deleteResult godoc
// @Summary Delete a Result
// @ID deleteResult
// @Tags results
// @Produce json
// @Success 200 {object} docs.ResponseResult "Result that was deleted"
// @Failure 400 {object} docs.ResponseError "Bad request"
// @Failure 404 {object} docs.ResponseError "Not found"
// @Failure 422 {object} docs.ResponseError "Unprocessable entity"
// @Failure 500 {object} docs.ResponseError "Internal server error"
// @Param resultID path int true "Result ID"
// @Router /results/{resultID} [delete]
// @Security Bearer
func deleteResult(c *gin.Context) {
	ok, result := CheckPermissions(c, database.Delete, "path", -1)
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
// @Success 200 {object} docs.ResponseResult "Result that was updated"
// @Failure 400 {object} docs.ResponseError "Bad request"
// @Failure 404 {object} docs.ResponseError "Not found"
// @Failure 422 {object} docs.ResponseError "Unprocessable entity"
// @Failure 500 {object} docs.ResponseError "Internal server error"
// @Param inputFile formData file true "File to be uploaded"
// @Param resultID path int true "Result ID"
// @Router /results/{resultID}/file [post]
// @Security Bearer
func addResultFile(c *gin.Context) {
	ok, _ := CheckPermissions(c, database.Update, "path", -1)
	if !ok {
		return
	}

	// TODO check permissions of scenario first (file will be added to scenario)

	// TODO add file to DB, associate with scenario and add file ID to result

}

// getResultFile godoc
// @Summary Download a result file
// @ID getResultFile
// @Tags results
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
// @Success 200 {object} docs.ResponseFile "File that was requested"
// @Failure 400 {object} docs.ResponseError "Bad request"
// @Failure 404 {object} docs.ResponseError "Not found"
// @Failure 422 {object} docs.ResponseError "Unprocessable entity"
// @Failure 500 {object} docs.ResponseError "Internal server error"
// @Param resultID path int true "Result ID"
// @Param fileID path int true "ID of the file to download"
// @Router /results/{resultID}/file/{fileID} [get]
// @Security Bearer
func getResultFile(c *gin.Context) {

	// check access
	ok, _ := CheckPermissions(c, database.Read, "path", -1)
	if !ok {
		return
	}

	// TODO download result file
}

// deleteResultFile godoc
// @Summary Delete a result file
// @ID deleteResultFile
// @Tags results
// @Produce json
// @Success 200 {object} docs.ResponseResult "Result for which file was deleted"
// @Failure 400 {object} docs.ResponseError "Bad request"
// @Failure 404 {object} docs.ResponseError "Not found"
// @Failure 422 {object} docs.ResponseError "Unprocessable entity"
// @Failure 500 {object} docs.ResponseError "Internal server error"
// @Param resultID path int true "Result ID"
// @Param fileID path int true "ID of the file to delete"
// @Router /results/{resultID}/file/{fileID} [delete]
// @Security Bearer
func deleteResultFile(c *gin.Context) {
	// TODO check access to scenario (file deletion) first

	// check access
	ok, _ := CheckPermissions(c, database.Update, "path", -1)
	if !ok {
		return
	}

}
