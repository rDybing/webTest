package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTestPage(t *testing.T) {
	tt := []struct {
		name  string
		value string
		err   string
	}{
		{name: "nothing added to URL", value: ""},
		{name: "something added to url", value: "test"},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			req, err := http.NewRequest("GET", "localhost:8080/"+tc.value, nil)
			if err != nil {
				t.Fatalf("Could not create new request:%v\n", err)
			}
			rec := httptest.NewRecorder()
			var o outPutT

			o.message = "ok"
			o.serverIP = "server"

			answer := o.message + "Inbound from     : :\nResponse from    : " + o.serverIP

			o.testPage(rec, req)
			res := rec.Result()
			defer res.Body.Close()

			b, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("Could not read response: %v\n", err)
			}
			if string(b) != answer {
				t.Errorf("Got wrong response. Expected:\n%s\nGot\n%s\n", answer, string(b))
			}

			if tc.err != "" {
				if res.StatusCode != http.StatusBadRequest {
					t.Errorf("Expected status bad request, got:\n%v\n", res.StatusCode)
				}
				if msg := string(bytes.TrimSpace(b)); msg != tc.err {
					t.Errorf("Expected message %q, got %q\n", tc.err, msg)
				}
				return
			}

			if res.StatusCode != http.StatusOK {
				t.Errorf("Expected status OK, got: %v\n", res.StatusCode)
			}
		})
	}
}

func TestRouting(t *testing.T) {
	var o outPutT
	srv := httptest.NewServer(o.handler())
	defer srv.Close()

	res, err := http.Get(fmt.Sprintf("%s", srv.URL))

	if err != nil {
		t.Fatalf("Could not send GET request: %v\n", err)
	}

	if res.StatusCode != http.StatusOK {
		t.Errorf("Expected status OK, got: %v\n", res.StatusCode)
	}
}

/*
func TestLoadTLS(t *testing.T) {
	var tls tlsT

}
*/
