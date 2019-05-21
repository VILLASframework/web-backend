package common

import (
	"time"

	// "github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
)

type User struct {
	//gorm.Model
	ID       uint   `gorm:"primary_key;auto_increment"`
	Username string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Mail     string `gorm:"default:''"`
	Role     string `gorm:"default:'user'"`

	Simulations []Simulation `gorm:"many2many:user_simulations"`
}

type Simulation struct {
	//gorm.Model
	ID              uint           `gorm:"primary_key;auto_increment"`
	Name            string         `gorm:"not null"`
	Running         bool           `gorm:"default:false"`
	StartParameters postgres.Jsonb

	Users   		[]User 			`gorm:"not null;many2many:user_simulations"`
	Models   		[]Model 		`gorm:"foreignkey:SimulationID"`
	Visualizations 	[]Visualization `gorm:"foreignkey:SimulationID"`
}

type Model struct {
	//gorm.Model
	ID              uint           `gorm:"primary_key;auto_increment"`
	Name            string         `gorm:"not null"`
	OutputLength    int            `gorm:"default:1"`
	InputLength     int            `gorm:"default:1"`
	StartParameters postgres.Jsonb

	SimulationID 	uint
	Simulator  		Simulator
	SimulatorID 	uint

	// NOTE: order of signals is important
	OutputMapping 	[]Signal 		`gorm:"foreignkey:ModelID"`
	InputMapping  	[]Signal 		`gorm:"foreignkey:ModelID"`

	//new in villasweb 2.0 (for CIM file of simulation model and other model file formats)
	Files 		  	[]File 			`gorm:"foreignkey:ModelID"`

}

type Signal struct {
	//gorm.Model
	ID                uint   `gorm:"primary_key;auto_increment"`
	Name              string `gorm:"not null"`
	Unit              string `gorm:"not null"`
	Index			  uint 	 `gorm:"not null"`
	Direction		  string `gorm:"not null"`
	ModelID uint
	//IsRecorded			bool 			`gorm:"default:false"`
}

type Simulator struct {
	//gorm.Model
	ID            uint   `gorm:"primary_key;auto_increment"`
	UUID          string `gorm:"unique;not null"`
	Host          string `gorm:"default:''"`
	Modeltype     string `gorm:"default:''"`
	Uptime        int    `gorm:"default:0"`
	State         string `gorm:"default:''"`
	StateUpdateAt time.Time
	Properties    postgres.Jsonb
	RawProperties postgres.Jsonb
}

type Visualization struct {
	//gorm.Model
	ID   uint   `gorm:"primary_key;auto_increment"`
	Name string `gorm:"not null"`
	Grid int    `gorm:"default:15"`

	SimulationID uint    `gorm:"not null"`

	Widgets []Widget `gorm:"foreignkey:VisualizationID"`
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
	CustomProperties postgres.Jsonb
	VisualizationID  uint			`gorm:"not null"`
	//new in villasweb 2.0
	Files 			 []File			`gorm:"foreignkey:WidgetID"`
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

	//new in villasweb 2.0
	ModelID uint `gorm:""`
	WidgetID uint `gorm:""`
}
