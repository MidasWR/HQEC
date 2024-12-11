package transactions_db

import (
	"gorm.io/gorm"
	"time"
)

type Transactions struct {
	Id     int       `gorm:"primaryKey"`
	UserId int       `gorm:"index"`
	Date   time.Time `gorm:"type:datetime"`
	Money  float32   `gorm:"type:decimal(10,2)"`
}
type Params struct {
	DateStart *time.Time `gorm:"type:datetime"`
	DateEnd   *time.Time `gorm:"type:datetime"`
	Date      *time.Time `gorm:"type:datetime"`
	Money     *float32   `gorm:"type:decimal(10,2)"`
	Page      int        `gorm:"default:1"`
	PerPage   int        `gorm:"default:1"`
}

func AddTransactionToDb(transaction Transactions, db *gorm.DB, id int) error {
	transaction.UserId = id
	return db.Create(&transaction).Error
}
func GetTransactionsFromDb(params Params, db *gorm.DB, id int) ([]Transactions, error) {
	query := db.Model(&Transactions{}).Where("user_id = ?", id)
	if params.Date != nil {
		query = query.Where("date = ?", params.Date)
	}
	if params.Money != nil {
		query = query.Where("money = ?", params.Money)
	}
	if params.DateStart != nil {
		query = query.Where("date >= ?", params.DateStart)
	}
	if params.DateEnd != nil {
		query = query.Where("date <= ?", params.DateEnd)
	}
	if params.PerPage == 0 {
		params.PerPage = 10
	}
	if params.Page == 0 {
		params.Page = 1
	}
	query = query.Limit(params.PerPage).Offset((params.Page - 1) * params.PerPage)
	var transactions []Transactions
	if err := query.Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}
