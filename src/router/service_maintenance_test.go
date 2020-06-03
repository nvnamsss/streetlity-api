package router_test

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func TestSendOrder(t *testing.T) {
	host := "http://" + "35.240.232.218" + "/user/notify"
	ids := []string{}
	ids = append(ids, "thosua")
	resp, err := http.PostForm(host, url.Values{
		"id":            ids,
		"notify-tittle": {"Customer is on service"},
		"notify-body":   {"A customer is looking for maintaning"},
		"data":          {"score:sss", "id:1"},
	})

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resp)
	fmt.Println("Done")
}

type MaintenanceData struct {
	Location []string
	Address  string
	Note     string
	Images   []string
}

func prepareHeader(req *http.Request) {
	req.Header.Set("Auth", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1OTQ0NDUzMzIsImlkIjoibXJYWVpaemFidiJ9.0Hd4SpIELulSuTxGAeuCPl_A33X-KoPUpRmgK4dTphk")
	req.Header.Set("Version", "1.0.0")
}

func TestAddMaintenanceService(t *testing.T) {
	host := "http://localhost:9000/service/maintenance/add"
	var arr []MaintenanceData = []MaintenanceData{}
	arr = append(arr, MaintenanceData{Location: []string{"1", "2"}, Address: "1150 sidney", Note: "Its a note", Images: []string{}})
	arr = append(arr, MaintenanceData{Location: []string{"5", "2"}, Address: "1150 sidney", Note: "Its a note", Images: []string{}})
	arr = append(arr, MaintenanceData{Location: []string{"7", "3"}, Address: "1150 sidney", Note: "Its a note", Images: []string{}})
	arr = append(arr, MaintenanceData{Location: []string{"1", "4"}, Address: "1150 sidney", Note: "Its a note", Images: []string{}})

	client := &http.Client{}

	for _, data := range arr {
		form := url.Values{
			"location": data.Location,
			"address":  {data.Address},
			"note":     {data.Note},
			"images":   data.Images,
		}

		req, _ := http.NewRequest("POST", host, strings.NewReader(form.Encode()))
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		prepareHeader(req)

		_, err := client.Do(req)
		if err != nil {
			log.Println(err.Error())
		}
	}

	t.Logf("Good")
}
