package simulation

import (
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/visualization"
	"github.com/jinzhu/gorm"
	// XXX: What about other models??
)

type Simulation struct {
	gorm.Model
	Name            string
	Running         bool              `gorm:"default:false"`
	Models          []SimulationModel // TODO: association & foreign key
	Projects        []Project         // TODO: association & foreign key
	User            []Users           // TODO: association & foreign key
	StartParameters []string          // TODO: Mixed Type
}
