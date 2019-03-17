package visualization

import (
	"github.com/jinzhu/gorm"
)

type Visualization struct {
	gorm.Model
	Name                 string
	VisualizationProject Project
	Widgets              string
	Grid                 int `gorm:"default:15"`
	VisualizationUser    User
}
