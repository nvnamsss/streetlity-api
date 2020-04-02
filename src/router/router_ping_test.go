package router_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Ping(w http.ResponseWriter, req *http.Request) {
	log.Println("Ping")
	w.Write([]byte("Ping"))

}
func TestPingNa(t *testing.T) {
	req, err := http.NewRequest("GET", "/health-check", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Ping)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	handler.ServeHTTP(rr, req)

	result := rr.Body.String()
	expect := "Ping"

	if result != expect {
		t.Errorf("Hi mom failed, expected %v, got %v", expect, result)
	} else {
		t.Logf("Hi mom passed, expected %v, got %v", expect, result)
	}
}
