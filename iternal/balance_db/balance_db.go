package balance_db

import (
	"gorm.io/gorm"
	"time"
)

type Balance struct {
	Id      int       `gorm:"primaryKey"`
	UserId  int       `gorm:"foreignKey:UserId;references:Id"`
	Balance float32   `gorm:"type:decimal(10,2)"`
	Date    time.Time `gorm:"type:datetime"`
}

func GetBalance(db *gorm.DB, id int) (float32, error) {
	var balance Balance
	err := db.Model(&Balance{}).Select("balance").Where("user_id = ?", id).Scan(&balance).Error
	if err != nil {
		return 0, err
	}
	return balance.Balance, nil
}
func ReBalance(db *gorm.DB, money float32, id int) error {
	if err := db.Model(&Balance{}).Where("user_id = ?", id).UpdateColumn("balance", gorm.Expr("balance + ?", money)).Error; err != nil {
		return err
	}
	return nil
}
func GetBalanceStructById(db *gorm.DB, id int) (*Balance, error) {
	var balance Balance
	if err := db.Model(&Balance{}).Select("balance").Where("id = ?", id).Scan(&balance).Error; err != nil {
		return nil, err
	}
	return &balance, nil
}
