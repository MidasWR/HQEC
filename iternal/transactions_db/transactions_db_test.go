package transactions_db

import (
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
	"time"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.Nil(t, err)
	err = db.AutoMigrate(&Transactions{})
	assert.Nil(t, err)
	return db
}

func TestAddTransactionToDb(t *testing.T) {
	db := setupTestDB(t)
	transaction := Transactions{
		Date:   time.Now(),
		Money:  100.50,
		UserId: 1,
	}

	err := AddTransactionToDb(transaction, db, 1)
	assert.Nil(t, err)
	var result Transactions
	err = db.First(&result).Error
	assert.Nil(t, err)
	assert.Equal(t, transaction.Money, result.Money)
}

func TestGetTransactionsFromDb(t *testing.T) {
	db := setupTestDB(t)
	transactions := []Transactions{
		{Date: time.Now().Add(-48 * time.Hour), Money: 50.0, UserId: 1},
		{Date: time.Now().Add(-24 * time.Hour), Money: 100.0, UserId: 1},
		{Date: time.Now().Truncate(time.Second), Money: 200.0, UserId: 1},
	}
	for _, transaction := range transactions {
		err := AddTransactionToDb(transaction, db, 1)
		assert.Nil(t, err)
	}

	params := Params{
		Page:    1,
		PerPage: 2,
	}
	result, err := GetTransactionsFromDb(params, db, 1)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(result))
	start := time.Now().Add(-36 * time.Hour)
	end := time.Now().Add(-12 * time.Hour)
	params = Params{
		DateStart: &start,
		DateEnd:   &end,
	}
	result, err = GetTransactionsFromDb(params, db, 1)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, float32(100.0), result[0].Money)
	date := time.Now().Truncate(time.Second)
	params = Params{
		Date: &date,
	}
	result, err = GetTransactionsFromDb(params, db, 1)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(result))
	assert.Equal(t, float32(200.0), result[0].Money)
}
