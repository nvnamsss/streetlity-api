package main_test

import (
	"io/ioutil"
	"net/http"
	"testing"
)

func TestPing(t *testing.T) {
	res, err := http.Get("http://localhost:9000/ping/")
	if err != nil {
		t.Fatal(err)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)

	result := string(body)
	expect := "Ping"

	if result != expect {
		t.Errorf("Hi mom failed, expected %v, got %v", expect, result)
	} else {
		t.Logf("Hi mom passed, expected %v, got %v", expect, result)
	}
}
