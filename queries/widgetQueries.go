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


