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
		log.Println("[Database]", "Adding new history:", e.Error())
	}

	return
}

func RemoveMaintenanceHistory(h MaintenanceHistory) (e error) {
	if e = Db.Delete(h).Error; e != nil {
		log.Println("[Database]", "Removing history:", e.Error())
	}

	return
}

func MaintenanceHistoryById(id int64) (h MaintenanceHistory, e error) {
	if e := Db.Find(&h, id).Error; e != nil {
		log.Println("[Database]", "Maintenance history with id:", id, ":", e.Error())
	}

	return
}

func UpdateMaintenanceHistory(id int64, maintenanceUser string, timestamp int64) error {
	h, e := MaintenanceHistoryById(id)

	if e != nil {
		return e
	}

	h.MaintenanceUser = maintenanceUser
	h.Timestamp = timestamp

	if e = Db.Save(&h).Error; e != nil {
		log.Println("[Database]", "Update maintenance history with id:", id, ":", e.Error())
	}

	return e
}
