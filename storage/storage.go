package storage

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"time"
)

type Users struct {
	Id           int            `gorm:"primaryKey"`
	Login        string         `gorm:"type:text;unique"`
	Password     string         `gorm:"type:text"`
	Transactions []Transactions `gorm:"foreignKey:UserId;references:Id"`
}

type Transactions struct {
	Id     int       `gorm:"primaryKey"`
	UserId int       `gorm:"index"`
	Date   time.Time `gorm:"type:datetime"`
	Money  float32   `gorm:"type:decimal(10,2)"`
}

type Budget struct {
	Id      int       `gorm:"primaryKey"`
	UserId  int       `gorm:"foreignKey:UserId;references:Id"`
	Balance float32   `gorm:"type:decimal(10,2)"`
	Date    time.Time `gorm:"type:datetime"`
}

func NewDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("all.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&Users{}); err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&Transactions{}); err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(&Budget{}); err != nil {
		return nil, err
	}
	return db, nil
}
