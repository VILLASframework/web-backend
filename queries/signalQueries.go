package queries

import (
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
)

func ReplaceSignals(model *common.Model, signals []common.Signal, direction string) error {

	db := common.GetDB()
	var err error

	if direction == "in" {
		err = db.Model(model).Select("InputMapping").Update("InputMapping", signals).Error
	} else {
		err = db.Model(model).Select("OutputMapping").Update("OutputMapping", signals).Error
	}

	return err

}