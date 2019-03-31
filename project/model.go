package project

import (
	"git.rwth-aachen.de/acs/public/villas/villasweb-backend-go/visualization"
	"github.com/jinzhu/gorm"
	// XXX: We don't need also User and Simulation???
)

type Project struct {
	gorm.Model
	Name              string `gorm:"not null"`
	ProjectUser       User   `gorm:"not null"`
	Visualizations    []Visualization
	ProjectSimulation Simulation `gorm:"not null"`
}

// TODO: execute before project.delete()
