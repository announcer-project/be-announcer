package models

import (
	"be_nms/models/modelsMember"
	"time"
)

type Role struct {
	ID        uint       `gorm:"primary_key"`
	CreatedAt time.Time  `json:"-"`
	UpdatedAt time.Time  `json:"-"`
	DeletedAt *time.Time `sql:"index" json:"-"`
	RoleName  string     `json:"rolename"`
	Require   bool       `json:"require"`
	SystemID  string     `json:"system_id"`

	Member []modelsMember.Member `gorm:"foreignkey:RoleID" json:"-"`
}
