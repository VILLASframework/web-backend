package simulation

import (
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/project"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/simulationModel"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/user"
	"github.com/jinzhu/gorm"
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
