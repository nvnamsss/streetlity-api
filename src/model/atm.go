package model

import (
	"github.com/golang/geo/r2"
	"github.com/jinzhu/gorm"
)

type ATM struct {
	gorm.Model
	Location r2.Point
}

func (ATM) TableName() string {
	return "atm"
}
