/** File package, methods.
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
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/scenario"
	"github.com/gin-gonic/gin"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
)

type File struct {
	database.File
}

func (f *File) ByID(id uint) error {
	db := database.GetDB()
	err := db.Find(f, id).Error
	if err != nil {
		return err
	}
	return nil
}

func (f *File) save() error {
	db := database.GetDB()
	err := db.Create(f).Error

	return err
}

func (f *File) download(c *gin.Context) error {

	// create unique file name
	filename := "file_" + strconv.FormatUint(uint64(f.ID), 10) + "_" + f.Name
	// detect the content type of the file
	contentType := http.DetectContentType(f.FileData)
	//Seems this headers needed for some browsers (for example without this headers Chrome will download files as txt)
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Data(http.StatusOK, contentType, f.FileData)

	return nil
}

func (f *File) Register(fileHeader *multipart.FileHeader, scenarioID uint) error {

	// Obtain properties of file
	f.Type = fileHeader.Header.Get("Content-Type")
	f.Name = filepath.Base(fileHeader.Filename)
	f.Size = uint(fileHeader.Size)
	f.Date = time.Now().String()
	f.ScenarioID = scenarioID

	// set file data
	fileContent, err := fileHeader.Open()
	if err != nil {
		return err
	}

	f.FileData, err = ioutil.ReadAll(fileContent)
	defer fileContent.Close()

	// Add image dimensions in case the file is an image
	if strings.Contains(f.Type, "image") || strings.Contains(f.Type, "Image") {
		// set the file reader back to the start of the file
		_, err := fileContent.Seek(0, 0)
		if err == nil {

			imageConfig, _, err := image.DecodeConfig(fileContent)
			if err != nil {
				log.Println("Unable to decode image configuration: Dimensions of image file are not set: ", err)
			} else {
				f.ImageHeight = imageConfig.Height
				f.ImageWidth = imageConfig.Width
			}
		} else {
			log.Println("Error on setting file reader back to start of file, dimensions not updated:", err)
		}
	}

	// Add File object with parameters to DB
	err = f.save()
	if err != nil {
		return err
	}

	// Create association to scenario
	db := database.GetDB()

	var so scenario.Scenario
	err = so.ByID(scenarioID)
	if err != nil {
		return err
	}

	err = db.Model(&so).Association("Files").Append(f).Error

	return err
}

func (f *File) update(fileHeader *multipart.FileHeader) error {

	// set file data
	fileContent, err := fileHeader.Open()
	if err != nil {
		return err
	}

	fileData, err := ioutil.ReadAll(fileContent)
	defer fileContent.Close()

	fileType := fileHeader.Header.Get("Content-Type")
	imageHeight := f.ImageHeight
	imageWidth := f.ImageWidth

	// Update image dimensions in case the file is an image
	if strings.Contains(fileType, "image") || strings.Contains(fileType, "Image") {
		// set the file reader back to the start of the file
		_, err := fileContent.Seek(0, 0)
		if err == nil {
			imageConfig, _, err := image.DecodeConfig(fileContent)
			if err != nil {
				log.Println("Unable to decode image configuration: Dimensions of image file are not updated.", err)
			} else {
				imageHeight = imageConfig.Height
				imageWidth = imageConfig.Width
			}
		} else {
			log.Println("Error on setting file reader back to start of file, dimensions not updated::", err)
		}
	} else {
		imageWidth = 0
		imageHeight = 0
	}

	db := database.GetDB()
	err = db.Model(f).Updates(map[string]interface{}{
		"Size":        uint(fileHeader.Size),
		"FileData":    fileData,
		"Date":        time.Now().String(),
		"Name":        filepath.Base(fileHeader.Filename),
		"Type":        fileType,
		"ImageHeight": imageHeight,
		"ImageWidth":  imageWidth,
	}).Error

	return err
}

func (f *File) Delete() error {

	db := database.GetDB()

	// remove association between file and scenario
	var so scenario.Scenario
	err := so.ByID(f.ScenarioID)
	if err != nil {
		return err
	}
	err = db.Model(&so).Association("Files").Delete(f).Error
	if err != nil {
		return err
	}

	// delete file from DB
	err = db.Delete(f).Error

	return err
}
