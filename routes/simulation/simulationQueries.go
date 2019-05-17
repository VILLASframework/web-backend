package simulation

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
)

func FindAllSimulations() ([]common.Simulation, int, error) {
	db := common.GetDB()
	var simulations []common.Simulation
	err := db.Find(&simulations).Error
	return simulations, len(simulations), err
}

func FindUserSimulations(user *common.User) ([]common.Simulation, int, error) {
	db := common.GetDB()
	var simulations []common.Simulation
	err := db.Model(user).Related(&simulations, "Simulations").Error
	return simulations, len(simulations), err
}


func AddFile(filename string, foldername string, filetype string, size uint, widgetID int, simulationmodelID int ) error {

	filesavepath := filepath.Join(foldername, filename)

	// get last modify time of target file
	fileinfo, err := os.Stat(filesavepath)
	if err != nil {
		return err
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
		// associate to Widget
		var widget common.Widget
		err := db.First(&widget,  widgetID).Error

		if err != nil {
			return err
		}
		err = db.Model(&widget).Association("Files").Append(&fileObj).Error
	}
	if simulationmodelID != -1 {
		// associate to Simulation Model
		var model common.SimulationModel
		err := db.First(&model, simulationmodelID).Error
		if err != nil {
			return err
		}
		err = db.Model(&model).Association("Files").Append(&fileObj).Error
	}

	return err
}

func SaveFile(file *multipart.FileHeader, filename string, foldername string, size uint, ) error {

	filesavepath := filepath.Join(foldername, filename)

	// Ensure folder with name foldername exists
	err := os.MkdirAll(foldername, os.ModePerm)
	if err != nil {
		// TODO error handling
		return err
	}

	fileTarget , errcreate := os.OpenFile(filesavepath, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
	if errcreate != nil {
		// TODO error handling: File could not be created
		return errcreate
	}
	defer fileTarget.Close()


	// Save file to target path
	uploadedFile, erropen := file.Open()
	if erropen != nil {
		// TODO error handling
		return erropen
	}
	defer uploadedFile.Close()

	var uploadContent = make([]byte, size)
	for {

		n, err := uploadedFile.Read(uploadContent)
		if err != nil && err != io.EOF {
			// TODO error handling
			return err
		}

		if n == 0 {
			break
		}

		_, err = fileTarget.Write(uploadContent[:n])
		if err != nil {
			// TODO error handling
			return err
		}

	}

	return err

}