package srpc

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"streelity/v1/config"
)

type MaintenanceOrder struct {
	Id              int64  `json:"Id"`
	CommonUser      string `json:"CommonUser"`
	MaintenanceUser string `json:"MaintenanceUser"`
	Timestamp       int64  `json:"Timestamp"`
	Receiver        string `json:"Receiver"`
	Reason          string `json:"Reason"`
	Note            string `json:"Note"`
	Status          int    `json:"Status"`
}

func RequestOrder(values url.Values) (res struct {
	Status  bool             `json:"Status"`
	Message string           `json:"Message"`
	Order   MaintenanceOrder `json:"Order"`
}, e error) {
	host := "http://" + config.Config.MaintenanceHost + "/order/"
	resp, e := http.PostForm(host, values)

	if e != nil {
		log.Println("[RPC]", "request order", e.Error())
		return
	}

	defer resp.Body.Close()

	body, e := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &res)

	return
}

func AcceptORder() (res struct {
	Status  bool             `json:"Status"`
	Message string           `json:"Message"`
	Order   MaintenanceOrder `json:"Order"`
}) {
	return
}
