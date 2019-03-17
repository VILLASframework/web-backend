package simulator

import (
	"github.com/jinzhu/gorm"
)

type Simulator struct {
	gorm.Model
	UUID          string
	Host          string `gorm:"default:"`
	Model         string `gorm:"default:"`
	Uptime        int    `gorm:"default:0"`
	State         string `gorm:"default:"`
	Properties    []string
	RawProperties []string
}
