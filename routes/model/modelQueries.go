package model

import (
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
)

func FindAllModels() ([]common.Model, int, error) {
	db := common.GetDB()
	var models []common.Model
	err := db.Find(&models).Error
	return models, len(models), err
}
