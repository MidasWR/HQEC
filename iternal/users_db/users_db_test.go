package users_db

import (
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.Nil(t, err)
	err = db.AutoMigrate(&Users{})
	assert.Nil(t, err)
	return db
}
func TestAddToDB(t *testing.T) {
	db := setupTestDB(t)
	login := "testuser"
	password := "password123"
	err := AddToDB(db, login, password)
	assert.Nil(t, err)
	var user Users
	err = db.Where("login = ?", login).First(&user).Error
	assert.Nil(t, err)
	assert.Equal(t, login, user.Login)
	assert.Equal(t, password, user.Password)
}

func TestFindLoginFromDB(t *testing.T) {
	db := setupTestDB(t)
	login := "testuser"
	password := "password123"
	err := AddToDB(db, login, password)
	assert.Nil(t, err)
	err = FindLoginFromDB(db, login)
	assert.Nil(t, err)
	err = FindLoginFromDB(db, "nonexistent")
	assert.NotNil(t, err)
}

func TestGetHashPassword(t *testing.T) {
	db := setupTestDB(t)
	login := "testuser"
	password := "password123"
	err := AddToDB(db, login, password)
	assert.Nil(t, err)
	hash, err := GetHashPassword(db, login)
	assert.Nil(t, err)
	assert.Equal(t, password, hash)
	_, err = GetHashPassword(db, "nonexistent")
	assert.NotNil(t, err)
}

func TestGetIdByLogin(t *testing.T) {
	db := setupTestDB(t)
	login := "testuser"
	password := "password123"
	err := AddToDB(db, login, password)
	assert.Nil(t, err)
	id, err := GetIdByLogin(db, login)
	assert.Nil(t, err)
	assert.NotZero(t, id)
	_, err = GetIdByLogin(db, "nonexistent")
	assert.NotNil(t, err)
}
