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
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"git.rwth-aachen.de/acs/public/villas/web-backend-go/database"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/component-configuration"
	"git.rwth-aachen.de/acs/public/villas/web-backend-go/routes/widget"
)

type File struct {
	database.File
}

func (f *File) byID(id uint) error {
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

func (f *File) register(fileHeader *multipart.FileHeader, objectType string, objectID uint) error {

	// Obtain properties of file
	f.Type = fileHeader.Header.Get("Content-Type")
	f.Name = filepath.Base(fileHeader.Filename)
	//f.Path = filepath.Join(getFolderName(objectType, objectID), f.Name)
	f.Size = uint(fileHeader.Size)
	f.Date = time.Now().String()

	var m component_configuration.ComponentConfiguration
	var w widget.Widget
	var err error
	if objectType == "config" {
		// check if config exists
		err = m.ByID(objectID)
		f.WidgetID = 0
		f.ConfigID = objectID
		if err != nil {
			return err
		}

	} else {
		// check if widget exists
		f.WidgetID = objectID
		f.ConfigID = 0
		err = w.ByID(uint(objectID))
		if err != nil {
			return err
		}

	}

	// set file data
	fileContent, err := fileHeader.Open()
	if err != nil {
		return err
	}

	f.FileData, err = ioutil.ReadAll(fileContent)
	defer fileContent.Close()

	// Add File object with parameters to DB
	err = f.save()
	if err != nil {
		return err
	}

	// Create association to config or widget
	if objectType == "config" {
		db := database.GetDB()
		err := db.Model(&m).Association("Files").Append(f).Error
		if err != nil {
			return err
		}
	} else {
		db := database.GetDB()
		err := db.Model(&w).Association("Files").Append(f).Error
		if err != nil {
			return err
		}
	}

	return nil
}

func (f *File) update(fileHeader *multipart.FileHeader) error {

	// set file data
	fileContent, err := fileHeader.Open()
	if err != nil {
		return err
	}

	fileData, err := ioutil.ReadAll(fileContent)
	defer fileContent.Close()

	db := database.GetDB()
	err = db.Model(f).Updates(map[string]interface{}{"Size": fileHeader.Size,
		"FileData": fileData,
		"Date":     time.Now().String()}).Error
	return err
}

func (f *File) delete() error {

	db := database.GetDB()

	if f.WidgetID > 0 {
		// remove association between file and widget
		var w widget.Widget
		err := w.ByID(f.WidgetID)
		if err != nil {
			return err
		}
		err = db.Model(&w).Association("Files").Delete(f).Error
		if err != nil {
			return err
		}
	} else {
		// remove association between file and config
		var m component_configuration.ComponentConfiguration
		err := m.ByID(f.ConfigID)
		if err != nil {
			return err
		}
		err = db.Model(&m).Association("Files").Delete(f).Error
		if err != nil {
			return err
		}
	}

	// delete file from DB
	err := db.Delete(f).Error

	return err
}
