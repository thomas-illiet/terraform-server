package database

import (
	"github.com/thomas-illiet/terrapi/pkg/model"
	"gorm.io/gorm"
)

func CreateModel(db *gorm.DB) *gorm.DB {
	db.AutoMigrate(&model.Module{})
	db.AutoMigrate(&model.ModuleSource{})
	db.AutoMigrate(&model.Deployment{})
	db.AutoMigrate(&model.DeploymentVariable{})
	db.AutoMigrate(&model.State{})
	return db
}
