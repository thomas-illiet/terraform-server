package model

import "time"

type Deployment struct {
	Id        uint `gorm:"primaryKey"`
	ModuleId  uint
	Name      string
	Variables *[]DeploymentVariable `gorm:"foreignKey:DeploymentId;references:Id"`
	States    *[]State              `gorm:"foreignKey:DeploymentId;references:Id"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type DeploymentVariable struct {
	DeploymentId int
	Name         string
	Value        string
}
