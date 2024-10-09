package model

import (
	"time"

	"gorm.io/gorm"
)

type State struct {
	ID           uint `gorm:"primaryKey"`
	DeploymentID uint
	Name         string
	Version      uint
	Data         []byte
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}
