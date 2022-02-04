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
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"git.rwth-aachen.de/acs/public/villas/web-backend-go/configuration"
	"github.com/gin-gonic/gin"

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

	if f.Key == "" {
		// Seems this headers needed for some browsers (for example without this headers Chrome will download files as txt)
		c.Header("Content-Description", "File Transfer")
		c.Header("Content-Disposition", "attachment; filename="+f.Name)
		c.Header("Expires", "")
		c.Header("Cache-Control", "")
		c.Data(http.StatusOK, f.Type, f.FileData)
	} else {
		url, err := f.getS3Url()
		if err != nil {
			return fmt.Errorf("failed to presign S3 request: %s", err)
		}
		c.Redirect(http.StatusFound, url)
	}

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
	defer fileContent.Close()

	bucket, err := configuration.GlobalConfig.String("s3.bucket")
	if err != nil || bucket == "" {
		// s3 object storage not used, s3.bucket param is empty
		// save file to postgres DB
		f.FileData, err = ioutil.ReadAll(fileContent)
		if err != nil {
			return err
		}
		f.Key = ""
	} else {
		err := f.putS3(fileContent)
		if err != nil {
			return fmt.Errorf("failed to upload to S3 bucket: %s", err)
		}
		log.Println("Saved new file in S3 object storage")
	}

	// Add image dimensions in case the file is an image
	if strings.Contains(f.Type, "image") || strings.Contains(f.Type, "Image") {
		// set the file reader back to the start of the file
		_, err := fileContent.Seek(0, 0)
		if err == nil {

			imageConfig, _, err := image.DecodeConfig(fileContent)
			if err != nil {
				log.Println("unable to decode image configuration: Dimensions of image file are not set, using default size 512x512, error:", err)
				f.ImageWidth = 512
				f.ImageHeight = 512
			} else {
				f.ImageHeight = imageConfig.Height
				f.ImageWidth = imageConfig.Width
			}

		} else {
			return fmt.Errorf("error on setting file reader back to start of file, dimensions not updated: %v", err)
		}
	}

	// Add File object with parameters to DB
	err = f.save()
	if err != nil {
		return err
	}

	// Create association to scenario
	db := database.GetDB()

	var so database.Scenario
	err = db.Find(&so, scenarioID).Error
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
	defer fileContent.Close()

	bucket, err := configuration.GlobalConfig.String("s3.bucket")
	if err != nil || bucket == "" {
		// s3 object storage not used, s3.bucket param is empty
		// save file to postgres DB
		f.FileData, err = ioutil.ReadAll(fileContent)
		if err != nil {
			return err
		}
		f.Key = ""
	} else {
		err := f.putS3(fileContent)
		if err != nil {
			return fmt.Errorf("failed to upload to S3 bucket: %s", err)
		}

		log.Println("Updated file in S3 object storage")
	}

	f.Type = fileHeader.Header.Get("Content-Type")
	f.Size = uint(fileHeader.Size)
	f.Date = time.Now().String()
	f.Name = filepath.Base(fileHeader.Filename)

	// Update image dimensions in case the file is an image
	if strings.Contains(f.Type, "image") || strings.Contains(f.Type, "Image") {
		// set the file reader back to the start of the file
		_, err := fileContent.Seek(0, 0)
		if err == nil {
			imageConfig, _, err := image.DecodeConfig(fileContent)
			if err != nil {
				log.Println("Unable to decode image configuration: Dimensions of image file are not updated.", err)
			}

			f.ImageHeight = imageConfig.Height
			f.ImageWidth = imageConfig.Width
		} else {
			log.Println("Error on setting file reader back to start of file, dimensions not updated::", err)
		}
	} else {
		f.ImageWidth = 0
		f.ImageHeight = 0
	}

	// Add File object with parameters to DB
	db := database.GetDB()
	err = db.Model(f).Updates(map[string]interface{}{
		"Size":        f.Size,
		"FileData":    f.FileData,
		"Date":        f.Date,
		"Name":        f.Name,
		"Type":        f.Type,
		"ImageHeight": f.ImageHeight,
		"ImageWidth":  f.ImageWidth,
		"Key":         f.Key,
	}).Error

	return err
}

func (f *File) Delete() error {

	db := database.GetDB()

	// remove association between file and scenario
	var so database.Scenario
	err := db.Find(&so, f.ScenarioID).Error
	if err != nil {
		return err
	}

	// delete file from s3 bucket
	if f.Key != "" {
		// TODO we do not delete the file from s3 object storage
		// to ensure that no data is lost if multiple File objects reference the same S3 data object
		// This behavior should be replaced by a different file handling in the future
		//err = f.deleteS3()
		//if err != nil {
		//	return err
		//}
		//log.Println("Deleted file in S3 object storage")
		log.Printf("Did NOT delete file with Key %v in S3 object storage!\n", f.Key)
	}

	err = db.Model(&so).Association("Files").Delete(f).Error
	if err != nil {
		return err
	}

	// delete file from DB
	err = db.Delete(f).Error

	return err
}
