package common

import (
	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
	"time"
)

type Simulator struct {
	gorm.Model
	UUID          string `gorm:"unique;not null"`
	Host          string `gorm:"default:''"`
	Modeltype     string `gorm:"default:''"`
	Uptime        int    `gorm:"default:0"`
	State         string `gorm:"default:''"`
	StateUpdateAt time.Time
	Properties    postgres.Jsonb // TODO: default value?
	RawProperties postgres.Jsonb // TODO: default value?
}

type File struct {
	gorm.Model
	Name        string    `gorm:"not null"`
	Path        string    `gorm:"not null"`
	Type        string    `gorm:"not null"`
	Size        uint      `gorm:"not null"`
	ImageHeight uint      //only required in case file is an image
	ImageWidth  uint      //only required in case file is an image
	FileUser    User      `gorm:"not null"`
	Date        time.Time `gorm:"default:Time.Now"`
}

type Project struct {
	gorm.Model
	Name              string `gorm:"not null"`
	ProjectUser       User   `gorm:"not null"`
	Visualizations    []Visualization
	ProjectSimulation Simulation `gorm:"not null"`
}

type Simulation struct {
	gorm.Model
	Name            string `gorm:"not null"`
	Running         bool   `gorm:"default:false"`
	Models          []SimulationModel
	Projects        []Project
	SimulationUser  User           `gorm:"not null"`
	StartParameters postgres.Jsonb // TODO default value
}

type SimulationModel struct {
	gorm.Model
	Name                string         `gorm:"not null"`
	OutputLength        int            `gorm:"default:1"`
	InputLength         int            `gorm:"default:1"`
	OutputMapping       []Signal       //order of signals is important
	InputMapping        []Signal       //order of signals is important
	StartParameters     postgres.Jsonb // TODO: default value?
	BelongsToSimulation Simulation     `gorm:"not null"`
	BelongsToSimulator  Simulator      `gorm:"not null"`
}

type User struct {
	gorm.Model
	Username    string `gorm:"unique;not null"`
	Password    string `gorm:"not null"`
	Mail        string `gorm:"default:"`
	Role        string `gorm:"default:user"`
	Projects    []Project
	Simulations []Simulation
	Files       []File
}

type Visualization struct {
	gorm.Model
	Name                 string  `gorm:"not null"`
	VisualizationProject Project `gorm:"not null"`
	Widgets              []Widget
	Grid                 int  `gorm:"default:15"`
	VisualizationUser    User `gorm:"not null"`
}

type Signal struct {
	gorm.Model
	Name string `gorm:"not null"`
	Unit string `gorm:"not null"`
	//IsRecorded			bool 			`gorm:"default:false"`
}

type Widget struct {
	gorm.Model
	Name             string         `gorm:"not null"`
	Type             string         `gorm:"not null"`
	Width            uint           `gorm:"not null"`
	Height           uint           `gorm:"not null"`
	MinWidth         uint           `gorm:"not null"`
	MinHeight        uint           `gorm:"not null"`
	X                int            `gorm:"not null"`
	Y                int            `gorm:"not null"`
	Z                int            `gorm:"not null"`
	IsLocked         bool           `gorm:"default:false"`
	CustomProperties postgres.Jsonb //TODO: default value?
}
