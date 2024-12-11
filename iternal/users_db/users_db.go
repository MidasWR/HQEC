package users_db

import "gorm.io/gorm"

type Users struct {
	Id       int    `gorm:"primaryKey"`
	Login    string `gorm:"type:text;unique"`
	Password string `gorm:"type:text"`
}

func FindLoginFromDB(Db *gorm.DB, login string) error {
	if err := Db.Model(&Users{}).Where("login=?", login).First(&Users{}).Error; err != nil {
		return err
	}
	return nil
}
func AddToDB(Db *gorm.DB, login string, password string) error {
	if err := Db.Model(&Users{}).Create(&Users{Id: 0, Login: login, Password: password}).Error; err != nil {
		return err
	}
	return nil
}
func GetHashPassword(Db *gorm.DB, login string) (string, error) {
	var user Users
	if err := Db.Model(&Users{}).Where("login=?", login).First(&user).Error; err != nil {
		return "", err
	}
	return user.Password, nil
}
func GetIdByLogin(Db *gorm.DB, login string) (int, error) {
	var user Users
	if err := Db.Model(&Users{}).Where("login=?", login).First(&user).Error; err != nil {
		return 0, err
	}
	return user.Id, nil
}
