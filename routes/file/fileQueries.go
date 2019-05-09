package file

import (
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
)

func FindAllFiles() ([]common.File, int, error) {
	db := common.GetDB()
	var files []common.File
	err := db.Find(&files).Error
	return files, len(files), err
}

func FindUserFiles(user *common.User) ([]common.File, int, error) {
	db := common.GetDB()
	var files []common.File
	err := db.Model(user).Related(&files, "Files").Error
	return files, len(files), err
}