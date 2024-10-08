package model

import (
	"time"

	"gorm.io/gorm"
)

type State struct {
	Id           uint `gorm:"primaryKey"`
	DeploymentId uint
	Name         string
	Version      uint
	Data         []byte
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}
