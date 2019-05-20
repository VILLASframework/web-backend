package file

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/gin-gonic/gin"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
)

func FindAllFiles() ([]common.File, int, error) {
	db := common.GetDB()
	var files []common.File
	err := db.Find(&files).Error
	if err != nil {
		// print error message to screen
		fmt.Println(fmt.Errorf("DB Error in FindAllFiles(): %q", err).Error())
	}
	return files, len(files), err
}

func FindUserFiles(user *common.User) ([]common.File, int, error) {
	db := common.GetDB()
	var files []common.File
	err := db.Model(user).Related(&files, "Files").Error
	return files, len(files), err
}

func FindFile(userID int, fileID string) ( common.File, error) {
	var file common.File
	db := common.GetDB()
	fileID_i, _ := strconv.Atoi(fileID)

	err := db.First(&file, fileID_i).Error

	return file, err

}

func FindFileByPath(path string) (common.File, error) {
	var file common.File
	db := common.GetDB()
	err := db.Where("Path = ?", path).Find(file).Error

	return file, err
}

func RegisterFile(c *gin.Context, widgetID int, simulationmodelID int, simulationID int){

	// Extract file from PUT request form
	file_header, err := c.FormFile("file")
	if err != nil {
		errormsg := fmt.Sprintf("Bad request. Get form error: %s", err.Error())
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errormsg,
		})
		return;
	}

	// Obtain properties of file
	filetype := file_header.Header.Get("Content-Type") // TODO make sure this is properly set in file header
	filename := filepath.Base(file_header.Filename)
	foldername := getFolderName(simulationID, simulationmodelID, widgetID)
	size := file_header.Size

	// Save file to local disc (NOT DB!)
	err = modifyFileOnDisc(file_header, filename, foldername, uint(size), true)
	if err != nil {
		errormsg := fmt.Sprintf("Internal Server Error. Error saving file: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errormsg,
		})
		return
	}

	// Add File object with parameters to DB
	saveFileInDB(c, filename, foldername, filetype, uint(size), widgetID, simulationmodelID, true)

}

func UpdateFile(c *gin.Context, widgetID int, simulationmodelID int, simulationID int){

	// Extract file from PUT request form
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
	foldername := getFolderName(simulationID, simulationmodelID, widgetID)

	err = modifyFileOnDisc(file_header, filename, foldername, uint(size), false)
	if err != nil {
		errormsg := fmt.Sprintf("Internal Server Error. Error saving file: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errormsg,
		})
		return
	}

	saveFileInDB(c, filename, foldername, filetype, uint(size), widgetID, simulationmodelID, false)
}

func ReadFile(c *gin.Context, widgetID int, simulationmodelID int, simulationID int){

	contentType := c.GetHeader("Content-Type")

	db := common.GetDB()
	var fileInDB common.File
	if widgetID != -1 {
		// get associated Widget
		var wdgt common.Widget
		err := db.First(&wdgt, simulationmodelID).Error
		if common.ProvideErrorResponse(c, err) {
			return
		}
		err = db.Model(&wdgt).Related(&fileInDB).Where("Type = ?", contentType).Error
		if common.ProvideErrorResponse(c, err) {
			return
		}

	} else if simulationmodelID != -1 {

		// get associated Simulation Model
		var model common.SimulationModel
		err := db.First(&model, simulationmodelID).Error
		if common.ProvideErrorResponse(c, err) {
			return
		}
		err = db.Model(&model).Related(&fileInDB).Where("Type = ?", contentType).Error
		if common.ProvideErrorResponse(c, err) {
			return
		}
	}

	//Seems this headers needed for some browsers (for example without this headers Chrome will download files as txt)
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", "attachment; filename="+fileInDB.Name )
	c.Header("Content-Type", contentType)
	c.File(fileInDB.Path)

	c.JSON(http.StatusOK, gin.H{
		"message": "OK.",
	})
}

func DeleteFile(c *gin.Context, widgetID int, simulationmodelID int, simulationID int){
	// TODO
}


func saveFileInDB(c *gin.Context, filename string, foldername string, filetype string, size uint, widgetID int, simulationmodelID int, createObj bool) {

	filesavepath := filepath.Join(foldername, filename)

	// get last modify time of target file
	fileinfo, err := os.Stat(filesavepath)
	if err != nil {
		errormsg := fmt.Sprintf("Internal Server Error. Error stat on file: %s", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errormsg,
		})
		return
	}
	modTime := fileinfo.ModTime()

	// Create file object for Database
	var fileObj common.File
	fileObj.Size = uint(size)
	fileObj.Type = filetype
	fileObj.Path = filesavepath
	fileObj.Date = modTime
	fileObj.Name = filename
	fileObj.ImageHeight = 0
	fileObj.ImageWidth = 0

	// Check if file shall be associated with Widget or Simulation Model
	db := common.GetDB()
	if widgetID != -1 {

		if createObj {
			// associate to Widget
			var wdgt common.Widget
			err := db.First(&wdgt, widgetID).Error
			if common.ProvideErrorResponse(c, err) {
				return
			}
			err = db.Model(&wdgt).Association("Files").Append(&fileObj).Error
		} else {
			// update file obj in DB
			fileInDB, err := FindFileByPath(filesavepath)
			if common.ProvideErrorResponse(c, err){
				return
			}

			err = db.Model(&fileInDB).Where("Path = ?", filesavepath).Updates(map[string]interface{}{"Size": fileObj.Size, "Date": fileObj.Date, "ImageHeight": fileObj.ImageHeight, "ImageWidth": fileObj.ImageWidth}).Error

		}

		if common.ProvideErrorResponse(c, err) == false {
			c.JSON(http.StatusOK, gin.H{
				"message": "OK.",
				"fileID": fileObj.ID,
			})
			return
		}

	}
	if simulationmodelID != -1 {

		if createObj {
			// associate to Simulation Model
			db := common.GetDB()
			var model common.SimulationModel
			err := db.First(&model, simulationmodelID).Error
			if common.ProvideErrorResponse(c, err) {
				return
			}
			err = db.Model(&model).Association("Files").Append(&fileObj).Error
		} else {
			// update file obj in DB
			fileInDB, err := FindFileByPath(filesavepath)
			if common.ProvideErrorResponse(c, err){
				return
			}

			err = db.Model(&fileInDB).Where("Path = ?", filesavepath).Updates(map[string]interface{}{"Size": fileObj.Size, "Date": fileObj.Date, "ImageHeight": fileObj.ImageHeight, "ImageWidth": fileObj.ImageWidth}).Error
		}

		if common.ProvideErrorResponse(c, err) == false {
			c.JSON(http.StatusOK, gin.H{
				"message": "OK.",
				"fileID": fileObj.ID,
			})
			return
		}
	}
}

func modifyFileOnDisc(file_header *multipart.FileHeader, filename string, foldername string, size uint, createFile bool) error {

	filesavepath := filepath.Join(foldername, filename)
	var err error

	if createFile {
		// Ensure folder with name foldername exists
		err = os.MkdirAll(foldername, os.ModePerm)
	} else {
		// test if file exists
		_, err = os.Stat(filesavepath)
	}
	if err != nil {
		return err
	}

	var open_options int
	if createFile {
		// create file it not exists, file MUST not exist
		open_options = os.O_RDWR|os.O_CREATE|os.O_EXCL
	} else {
		open_options = os.O_RDWR
	}

	fileTarget , err := os.OpenFile(filesavepath, open_options, 0666)
	if err != nil {
		return err
	}
	defer fileTarget.Close()

	// Save file to target path
	uploadedFile, err := file_header.Open()
	if err != nil {
		return err
	}
	defer uploadedFile.Close()

	var uploadContent = make([]byte, size)
	for {

		n, err := uploadedFile.Read(uploadContent)
		if err != nil && err != io.EOF {
			return err
		}

		if n == 0 {
			break
		}

		_, err = fileTarget.Write(uploadContent[:n])
		if err != nil {
			return err
		}

	}
	return err
}


func getFolderName(simulationID int, simulationmodelID int, widgetID int) string {
	base_foldername := "files/"
	elementname := ""
	elementid := 0
	if simulationmodelID == -1{
		elementname = "/widget_"
		elementid = widgetID
	} else {
		elementname = "/simulationmodel_"
		elementid = simulationmodelID
	}


	foldername := base_foldername + "simulation_"+ string(simulationID) + elementname + string(elementid) + "/"
	return foldername
}