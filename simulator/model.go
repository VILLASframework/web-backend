package simulator

import (
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
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
	Properties    pq.StringArray `gorm:"type:varchar(128)[]"` // TODO: mixed type
	RawProperties pq.StringArray `gorm:"type:varchar(128)[]"` // TODO: mixed type
}
