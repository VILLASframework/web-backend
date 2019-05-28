package widget

import (
	"fmt"

	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/routes/visualization"
)

type Widget struct {
	common.Widget
}

func (w *Widget) save() error {
	db := common.GetDB()
	err := db.Create(w).Error
	return err
}

func (w *Widget) ByID(id uint) error {
	db := common.GetDB()
	err := db.Find(w, id).Error
	if err != nil {
		return fmt.Errorf("Widget with id=%v does not exist", id)
	}
	return nil
}

func (w *Widget) addToVisualization(visID uint) error {
	db := common.GetDB()
	var vis visualization.Visualization
	err := vis.ByID(uint(visID))
	if err != nil {
		return err
	}

	// save visualization to DB
	err = w.save()
	if err != nil {
		return err
	}

	// associate visualization with simulation
	err = db.Model(&vis).Association("Widgets").Append(w).Error

	return err
}

func (w *Widget) update(modifiedWidget Widget) error {
	db := common.GetDB()
	err := db.Model(w).Update(modifiedWidget).Error
	if err != nil {
		return err
	}

	return err
}
