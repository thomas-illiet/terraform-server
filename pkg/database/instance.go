package database

import (
	"fmt"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func ConnectDb(host string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(host), &gorm.Config{})
	if err != nil {
		fmt.Println("Failed to connect database")
	}
	return db
}
