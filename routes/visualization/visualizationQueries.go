package visualization

import (
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
)

func FindAllVisualizations() ([]common.Visualization, int, error) {
	db := common.GetDB()
	var visualization []common.Visualization
	err := db.Find(&visualization).Error
	return visualization, len(visualization), err
}
