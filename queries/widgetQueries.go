package queries

import (
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
)

func FindVisualizationWidgets(visualization *common.Visualization) ([]common.Widget, int, error) {
	db := common.GetDB()
	var widgets []common.Widget
	err := db.Model(visualization).Related(&widgets, "Widgets").Error
	return widgets, len(widgets), err
}


func FindWidget(widgetID int) (common.Widget, error){
	db := common.GetDB()
	var w common.Widget
	err := db.First(&w, widgetID).Error
	return w, err
}


func FindWidgetsOfVisualization(vis * common.Visualization) ([]common.Widget, int, error) {
	db := common.GetDB()
	var widgets []common.Widget
	err := db.Model(vis).Related(&vis, "Widgets").Error
	return widgets, len(widgets), err
}

func FindWidgetOfVisualization(visualization *common.Visualization, widgetID int) (common.Widget, error){
	db := common.GetDB()
	var widget common.Widget
	err := db.Model(visualization).Where("ID = ?", widgetID).Related(&widget, "Widgets").Error
	return widget, err
}


func AddWidgetToVisualization(vis *common.Visualization, widget_input * common.Widget) error {

	db := common.GetDB()

	// Add widget to DB
	err := db.Create(widget_input).Error
	if err != nil {
		return err
	}

	// Add association with visualization
	err = db.Model(vis).Association("Widgets").Append(widget_input).Error
	return err
}

func UpdateWidgetOfVisualization(vis * common.Visualization, widget_input common.Widget, widgetID int) error {
	db := common.GetDB()

	// Get widget of visualization that matches with ID (= widget to be updated)
	var widget_old common.Widget
	err := db.Model(vis).Where("ID = ?", widgetID).Related(&widget_old, "Widgets").Error
	if err != nil {
		return err
	}

	// Update widget in DB
	err = db.Model(&widget_old).Updates(map[string]interface{}{"Name": widget_input.Name, "Type": widget_input.Type, "MinHeight": widget_input.MinHeight, "MinWidth": widget_input.MinWidth, "Height": widget_input.Height, "Width": widget_input.Width, "X": widget_input.X, "Y": widget_input.Y, "Z": widget_input.Z, "CustomProperties": widget_input.CustomProperties}).Error
	return err

}