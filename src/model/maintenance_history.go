package model

import "log"

type MaintenanceHistory struct {
	Id              int64
	MaintenanceUser string `gorm:"column:maintenance_user"`
	CommonUser      string `gorm:"column:common_user"`
	Timestamp       int64  `gorm:"type:datetime"`
}

func (MaintenanceHistory) TableName() string {
	return "maintenance_history"
}

func AddMaintenanceHistory(h MaintenanceHistory) (e error) {
	if e = Db.Create(&h).Error; e != nil {
		log.Println("[Database]", e.Error())
	}

	return
}

func RemoveMaintenanceHistory(h MaintenanceHistory) (e error) {
	if e = Db.Delete(h).Error; e != nil {
		log.Println("[Database]", e.Error())
	}

	return
}
