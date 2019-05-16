package file

import (
	"fmt"
	"strconv"

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
	//TODO Check here if user owns the file
	var file common.File
	db := common.GetDB()
	fileID_i, _ := strconv.Atoi(fileID)

	err := db.First(&file, fileID_i).Error

	return file, err

}

func AddFile(m map[string]interface{}) error {

	// TODO deserialize m (JSON file object) to data struct File

	// TODO we need the user here as well to be able to create the association in the DB

	// TODO add deserialized File to DB

	var err error
	return err

}