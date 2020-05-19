package router_test

import (
	"fmt"
	"net/http"
	"net/url"
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
