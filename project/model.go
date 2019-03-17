package project

import (
	"git.rwth-aachen.de/stemavros/villasweb-backend-go/visualization"
	"github.com/jinzhu/gorm"
	// XXX: We don't need also User and Simulation???
)

type Project struct {
	gorm.Model
	Name              string
	UserProject       User            // XXX: association?
	Visualizations    []Visualization // TODO: association & foreign key
	SimulationProject Simulation      // XXX: association?n
}

// TODO: execute before project.delete()
