package simulator

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Simulator struct {
	gorm.Model
	UUID          string   `gorm:"unique;not null"`
	Host          string   `gorm:"default:"`
	Model         string   `gorm:"default:"`
	Uptime        int      `gorm:"default:0"`
	State         string   `gorm:"default:"`
	StateUpdateAt Time     `gorm:"default:Time.Now"`
	Properties    []string // TODO: mixed type
	RawProperties []string // TODO: mixed type
}
