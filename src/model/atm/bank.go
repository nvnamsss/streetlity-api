package atm

import (
	"errors"
	"log"
	"streelity/v1/model"
)

const BankTableName = "bank"

type Bank struct {
	Id   int64
	Name string `gorm:"column:name"`
}

func (Bank) TableName() string {
	return BankTableName
}

func AllBanks() []Bank {
	var banks []Bank
	model.Db.Find(&banks)

	return banks
}

func CreateBank(s Bank) (e error) {
	if _, e = BankByName(s.Name); e == nil {
		e = errors.New("Bank was existed")
		log.Println("[Database]", "Create new bank", e.Error())
		return
	}

	if e = model.Db.Create(&s).Error; e != nil {
		log.Println("[Database]", "Add bank", e.Error())
		return
	}

	return nil
}

func BankByName(name string) (bank Bank, e error) {
	bank.Name = name
	db := model.Db.Find(&bank)

	if e := db.Error; e != nil {
		log.Println("[Database]", "Get bank", e.Error())
	}

	if db.RowsAffected == 0 {
		e := errors.New("Bank was not found")
		log.Println("[Database]", e.Error())
	}

	return
}
