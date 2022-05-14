package model

import "gorm.io/gorm"

type Host struct {
	gorm.Model
	Name    string
	Content string
	Groups  []*Group `gorm:"many2many:host_group"`
}

func (h *Host) TableName() string {
	return "host"
}
