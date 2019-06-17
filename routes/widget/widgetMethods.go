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

func (w *Widget) addToVisualization() error {
	db := common.GetDB()
	var vis visualization.Visualization
	err := vis.ByID(uint(w.VisualizationID))
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
	err := db.Model(w).Updates(map[string]interface{}{
		"Name":             modifiedWidget.Name,
		"Type":             modifiedWidget.Type,
		"Width":            modifiedWidget.Width,
		"Height":           modifiedWidget.Height,
		"MinWidth":         modifiedWidget.MinWidth,
		"MinHeight":        modifiedWidget.MinHeight,
		"X":                modifiedWidget.X,
		"Y":                modifiedWidget.Y,
		"Z":                modifiedWidget.Z,
		"IsLocked":         modifiedWidget.IsLocked,
		"CustomProperties": modifiedWidget.CustomProperties,
	}).Error

	return err
}

func (w *Widget) delete() error {

	db := common.GetDB()
	var vis visualization.Visualization
	err := vis.ByID(w.VisualizationID)
	if err != nil {
		return err
	}

	// remove association between Visualization and Widget
	// Widget itself is not deleted from DB, it remains as "dangling"
	err = db.Model(&vis).Association("Widgets").Delete(w).Error

	// TODO: What about files that belong to a widget? Keep them or remove them here?

	return err
}
