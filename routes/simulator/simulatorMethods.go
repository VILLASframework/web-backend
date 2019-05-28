package simulator

import (
	"fmt"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
)

type Simulator struct {
	common.Simulator
}

func (s *Simulator) save() error {
	db := common.GetDB()
	err := db.Create(s).Error
	return err
}

func (s *Simulator) ByID(id uint) error {
	db := common.GetDB()
	err := db.Find(s, id).Error
	if err != nil {
		return fmt.Errorf("Simulator with id=%v does not exist", id)
	}
	return nil
}
