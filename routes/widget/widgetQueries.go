package widget

import (
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
)

func FindVisualizationWidgets(visualization *common.Visualization) ([]common.Widget, int, error) {
	db := common.GetDB()
	var widgets []common.Widget
	err := db.Model(visualization).Related(&widgets, "Widgets").Error
	return widgets, len(widgets), err
}


