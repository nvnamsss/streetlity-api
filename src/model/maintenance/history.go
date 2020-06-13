package maintenance

import (
	"log"
	"streelity/v1/model"
)

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
	if e = model.Db.Create(&h).Error; e != nil {
		log.Println("[Database]", "Adding new history:", e.Error())
	}

	return
}

func RemoveMaintenanceHistory(h MaintenanceHistory) (e error) {
	if e = model.Db.Delete(h).Error; e != nil {
		log.Println("[Database]", "Removing history:", e.Error())
	}

	return
}

func RemoveMaintenanceHistoriesById(ids ...int64) (e error) {
	for _, id := range ids {
		if e = model.Db.Where("id=?", id).Delete(&MaintenanceHistory{}).Error; e != nil {
			log.Println("[Database]", "remove M history", e.Error())
		}
	}

	return
}

func queryMaintenanceHistory(h MaintenanceHistory) (history MaintenanceHistory, e error) {
	history = h
	if e = model.Db.Find(&history).Error; e != nil {
		log.Println("[Database]", "query maintenance history", e.Error())
	}

	return
}

func MaintenanceHistoriesByMUser(mUser string) (histories []MaintenanceHistory, e error) {
	if e = model.Db.Where("maintenance_user=?", mUser).Find(&histories).Error; e != nil {
		log.Println("[Database]", "query M history by M user", e.Error())
	}

	return
}

func MaintenanceHistoryById(id int64) (h MaintenanceHistory, e error) {
	if e := model.Db.Find(&h, id).Error; e != nil {
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

	if e = model.Db.Save(&h).Error; e != nil {
		log.Println("[Database]", "Update maintenance history with id:", id, ":", e.Error())
	}

	return e
}
