package models

import (
	"be_nms/models/modelsMember"

	"github.com/jinzhu/gorm"
)

type Role struct {
	gorm.Model
	RoleName string `json:"rolename"`
	Require  bool   `json:"require"`
	SystemID string `json:"system_id"`

	Member []modelsMember.Member `gorm:"foreignkey:RoleID"`
}
