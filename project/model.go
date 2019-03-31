package project

import (
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/simulation"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/user"
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/visualization"
	"github.com/jinzhu/gorm"
)

type Project struct {
	gorm.Model
	Name              string `gorm:"not null"`
	ProjectUser       User   `gorm:"not null"`
	Visualizations    []Visualization
	ProjectSimulation Simulation `gorm:"not null"`
}

// TODO: execute before project.delete()
