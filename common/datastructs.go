package common

import (
	//"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
	"time"
)

type Simulator struct {
	//gorm.Model
	ID            uint   `gorm:"primary_key;auto_increment"`
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
	//gorm.Model
	ID          uint   `gorm:"primary_key;auto_increment"`
	Name        string `gorm:"not null"`
	Path        string `gorm:"not null"`
	Type        string `gorm:"not null"`
	Size        uint   `gorm:"not null"`
	ImageHeight uint   // only required in case file is an image
	ImageWidth  uint   // only required in case file is an image
	Date        time.Time

	//remove belongs to User relation
	//User   User `gorm:"not null;association_autoupdate:false"`
	UserID uint `gorm:""`
	SimulationModelID uint `gorm:""`
}

type Project struct {
	//gorm.Model
	ID   uint   `gorm:"primary_key;auto_increment"`
	Name string `gorm:"not null"`

	User   User `gorm:"not null;association_autoupdate:false"`
	UserID uint `gorm:"not null"`

	Simulation   Simulation `gorm:"not null;association_autoupdate:false"`
	SimulationID uint       `gorm:"not null"`

	Visualizations []Visualization `gorm:"association_autoupdate:false"`
}

type Simulation struct {
	//gorm.Model
	ID              uint           `gorm:"primary_key;auto_increment"`
	Name            string         `gorm:"not null"`
	Running         bool           `gorm:"default:false"`
	StartParameters postgres.Jsonb // TODO default value

	User   User `gorm:"not null;association_autoupdate:false"`
	UserID uint `gorm:"not null"`

	Models   []SimulationModel `gorm:"foreignkey:BelongsToSimulationID;association_autoupdate:false"`
	Projects []Project         `gorm:"association_autoupdate:false"`
}

type SimulationModel struct {
	//gorm.Model
	ID              uint           `gorm:"primary_key;auto_increment"`
	Name            string         `gorm:"not null"`
	OutputLength    int            `gorm:"default:1"`
	InputLength     int            `gorm:"default:1"`
	StartParameters postgres.Jsonb // TODO: default value?

	BelongsToSimulation   Simulation `gorm:"not null;association_autoupdate:false"`
	BelongsToSimulationID uint       `gorm:"not null"`

	BelongsToSimulator   Simulator `gorm:"not null;association_autoupdate:false"`
	BelongsToSimulatorID uint      `gorm:"not null"`
	//new in villasweb 2.0
	Files 			[]File 			`gorm:""`

	// NOTE: order of signals is important
	OutputMapping []Signal `gorm:""`
	InputMapping  []Signal `gorm:""`
}

type User struct {
	//gorm.Model
	ID       uint   `gorm:"primary_key;auto_increment"`
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Mail     string `gorm:"default:''"`
	Role     string `gorm:"default:'user'"`

	Projects    []Project    `gorm:"association_autoupdate:false"`
	Simulations []Simulation `gorm:"association_autoupdate:false"`
	Files       []File       `gorm:""`
}

type Visualization struct {
	//gorm.Model
	ID   uint   `gorm:"primary_key;auto_increment"`
	Name string `gorm:"not null"`
	Grid int    `gorm:"default:15"`

	Project   Project `gorm:"not null;association_autoupdate:false"`
	ProjectID uint    `gorm:"not null"`

	User   User `gorm:"not null;association_autoupdate:false"`
	UserID uint `gorm:"not null"`

	Widgets []Widget `gorm:""`
}

type Signal struct {
	//gorm.Model
	ID                uint   `gorm:"primary_key;auto_increment"`
	Name              string `gorm:"not null"`
	Unit              string `gorm:"not null"`
	SimulationModelID uint
	//IsRecorded			bool 			`gorm:"default:false"`
}

type Widget struct {
	//gorm.Model
	ID               uint           `gorm:"primary_key;auto_increment"`
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
	CustomProperties postgres.Jsonb // TODO: default value?
	VisualizationID  uint
}
