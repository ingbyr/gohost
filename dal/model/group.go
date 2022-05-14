package model

import "gorm.io/gorm"

type Group struct {
	gorm.Model
	Name  string
	Hosts []*Host `gorm:"many2many:host_groups"`
}

func (g *Group) TableName() string {
	return "group"
}
