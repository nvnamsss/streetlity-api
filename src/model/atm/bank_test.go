package atm_test

import (
	"streelity/v1/model"
	"streelity/v1/model/atm"
	"testing"
)

func TestCreateBank(t *testing.T) {
	model.ConnectSync()
	names := []string{"Agribank",
		"Vietcombank",
		"Shinhanbank",
		"VPBank",
		"OceanBank",
		"VietinBank",
		"HDBank",
		"VIBank",
		"EximBank",
		"Sacombank",
		"DongABank",
		"NamABank",
		"SaigonBank",
		"TPBank"}

	for _, name := range names {
		atm.CreateBank(atm.Bank{Name: name})
	}
}
