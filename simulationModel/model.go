package simulationModel

import (
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/simulation"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/simulator"
	"github.com/jinzhu/gorm"
)

type SimulationModel struct {
	gorm.Model
	Name            string
	OutputLength    int        `gorm:"default:1"`
	InputLength     int        `gorm:"default:1"`
	OutputMapping   []string   // TODO: Mixed Type
	InputMapping    []string   // TODO: Mixed Type
	StartParameters []string   // TODO: Mixed Type
	ModelSimulation Simulation `gorm:"not null"`
	ModelSimulator  Simulator  `gorm:"not null"`
}
