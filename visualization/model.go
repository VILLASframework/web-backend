package visualization

import (
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/user"
	"github.com/jinzhu/gorm"
)

type Visualization struct {
	gorm.Model
	Name                 string   `gorm:"not null"`
	VisualizationProject Project  `gorm:"not null"`
	Widgets              []string // XXX: array of what type?
	Grid                 int      `gorm:"default:15"`
	VisualizationUser    User     `gorm:"not null"`
}
