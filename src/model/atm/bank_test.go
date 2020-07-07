package atm_test

import (
	"streelity/v1/model"
	"streelity/v1/model/atm"
	"testing"
)

func TestCreateBank(t *testing.T) {
	model.ConnectSync()
	names := []string{"Agribank",
		"ABBank",
		"ACB",
		"AgriBank",
		"ANZ",
		"Bac A Bank",
		"Bangkok Bank",
		"Bank of China",
		"Bank of Coummunication",
		"Bank of India HCM Branch",
		"Bank of Tokyo-Mitsubishi",
		"Bao Viet Bank",
		"BIDC",
		"BIDV",
		"BNp Paribas",
		"Cathay Bank",
		"China Construction Bank",
		"Chinatrust Commercial Bank",
		"Chinfon Bank",
		"CIMD Viet Nam",
		"CitiBank",
		"Commonwealth Bank Vietnam",
		"Credit Agricole CIB",
		"Dai A Bank",
		"DBS Bank",
		"Deutsche Bank",
		"Dong A Bank",
		"EximBank",
		"Far East National Bank",
		"First Commercial Bank",
		"GPBank",
		"HDBank",
		"Hong Leong Bank",
		"HSBC",
		"Hua Nan Commercial Bank",
		"Indovina Bank",
		"JP Morgan Chase Bank",
		"Kho bac nha nuoc VN",
		"Kien Long Bank",
		"Kookmin Bank",
		"Korea Exchange Bank",
		"LaoVietBank",
		"LienVietPostBank",
		"Malayan Banking Berhad",
		"Maritime Bank",
		"MayBank",
		"MBBank",
		"Mega ICBC",
		"Mizuho Corporate Bank",
		"Nam A Bank",
		"Natixis Bank",
		"NCB",
		"NH Chinh sach Xa hoi (VBSP)",
		"NH Cong Thuong HQ",
		"NH Cong Thuong TQ (ICBC)",
		"NH Hop Tac Xa VN (Co-opBank)",
		"NH Lien doanh Viet Nga (VRB)",
		"NH Phat trien Viet Nam (VDB)",
		"NH Quoc Dan (NCB)",
		"NH Xay Dung VN (CB)",
		"Nong Hyup Bank",
		"OCB",
		"OCBC",
		"Ocean Bank",
		"PGBank",
		"PVComBank",
		"Saigon Comercial Bank",
		"SBV",
		"SeA",
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

	t.Logf("Completed")
}
