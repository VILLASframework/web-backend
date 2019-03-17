package simulationModel

import (
	"github.com/jinzhu/gorm"
	// XXX: other models
)

type SimulationModel struct {
	gorm.Model
	Name            string
	OutputLength    int        `gorm:"default:1"`
	InputLength     int        `gorm:"default:1"`
	OutputMapping   []string   // TODO: Mixed Type
	InputMapping    []string   // TODO: Mixed Type
	StartParameters []string   // TODO: Mixed Type
	ModelSimulation Simulation // TODO: association & foreign key
	ModelSimulator  Simulator  // TODO: association & foreign key
}
