package model

import "gorm.io/gorm"

type Host struct {
	gorm.Model
	Name    string
	GroupID uint
	Content string
	Groups  []*Group `gorm:"many2many:group_hosts"`
}

func (h *Host) TableName() string {
	return "host"
}
