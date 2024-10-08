package model

import (
	"time"
)

type Module struct {
	Id          uint `gorm:"primaryKey"`
	Name        string
	Source      ModuleSource  `gorm:"foreignKey:ModuleId;references:Id"`
	Deployments *[]Deployment `gorm:"foreignKey:ModuleId;references:Id"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ModuleSource struct {
	ModuleId   uint
	Repository string
	Branch     string
	Path       string
}
