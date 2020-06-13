package atm

import (
	"log"
	"streelity/v1/model"
)

type Bank struct {
	Id   int64
	Name string `gorm:"column:name"`
}

func (Bank) TableName() string {
	return "bank"
}

func AllBanks() []Bank {
	var banks []Bank
	model.Db.Find(&banks)

	return banks
}

func AddBank(s Bank) (e error) {
	if e = model.Db.Create(&s).Error; e != nil {
		log.Println("[Database]", "Add bank", e.Error())
		return
	}

	return nil
}

func BankByName(name string) (bank Bank, e error) {
	bank.Name = name
	if e = model.Db.Find(&bank).Error; e != nil {
		log.Println("[Database]", "Get bank", e.Error())
	}

	return
}
