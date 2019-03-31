package simulation

import (
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/visualization"
	"github.com/jinzhu/gorm"
	// XXX: What about other models??
)

type Simulation struct {
	gorm.Model
	Name            string `gorm:"not null"`
	Running         bool   `gorm:"default:false"`
	Models          []SimulationModel
	Projects        []Project
	SimulationUser  User     `gorm:"not null"`
	StartParameters []string // TODO: Mixed Type
}
