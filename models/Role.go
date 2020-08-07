package models

import (
	"be_nms/models/modelsMember"

	"github.com/jinzhu/gorm"
)

type Role struct {
	gorm.Model
	RoleName string
	SystemID string

	Member []modelsMember.Member `gorm:"foreignkey:RoleID"`
}
