package models

import (
	"be_nms/models/modelsMember"

	"github.com/jinzhu/gorm"
)

type Role struct {
	gorm.Model
	RoleName string `gorm:"not null" json:"rolename"`
	Require  bool   `gorm:"not null" json:"require"`
	SystemID string `gorm:"not null" json:"system_id"`

	Member []modelsMember.Member `gorm:"foreignkey:RoleID" json:"-"`
}
