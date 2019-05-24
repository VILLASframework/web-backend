package visualization

import (
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/common"
)

func FindAllVisualizationsOfSim(sim *common.Simulation) ([]common.Visualization, int, error) {
	db := common.GetDB()
	var visualizations []common.Visualization
	err := db.Order("ID asc").Model(sim).Related(&visualizations, "Visualizations").Error
	return visualizations, len(visualizations), err
}

func FindVisualizationOfSim(sim *common.Simulation, visID int) (common.Visualization, error) {
	db := common.GetDB()
	var vis common.Visualization
	err := db.Model(sim).Where("ID = ?", visID).Related(&vis, "Visualizations").Error
	return vis, err
}

func AddVisualizationToSim(sim * common.Simulation, vis * common.Visualization) error {
	db := common.GetDB()

	// Add visualization to DB
	err := db.Create(vis).Error
	if err != nil {
		return err
	}

	// Add association with simulation
	err = db.Model(sim).Association("Visualizations").Append(vis).Error
	return err

}

func UpdateVisualizationOfSim(sim * common.Simulation, vis common.Visualization, visID int) error {
	db := common.GetDB()

	// Get visualization of simulation that matches with ID (= visualization to be updated)
	var vis_old common.Visualization
	err := db.Model(sim).Where("ID = ?", visID).Related(&vis_old, "Visualizations").Error
	if err != nil {
		return err
	}

	// Update visualization in DB (only name and grid can be updated)
	err = db.Model(&vis_old).Updates(map[string]interface{}{"Name": vis.Name, "Grid": vis.Grid}).Error
	return err
}