package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
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

func TestGetPrivateIP(t *testing.T) {
	var o outPutT

	const local = "127.0.0.1"

	tt := []struct {
		name  string
		local bool
		tlsOK bool
	}{
		{name: "run local, no tls loaded", local: true, tlsOK: false},
		{name: "run local, with tls loaded", local: true, tlsOK: true},
		{name: "run WAN, no tls loaded", local: false, tlsOK: false},
		{name: "run WAN, with tls loaded", local: false, tlsOK: true},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			port := "80"
			if tc.tlsOK {
				port = "443"
			}
			if tc.local {
				port = "8080"
			}
			o.getPrivateIP(tc.local, tc.tlsOK)
			ip := strings.Split(o.serverIP, ":")
			if ip[1] != port {
				t.Fatalf("Expected port %s, got %s\n", port, ip[1])
			}
			if tc.local && ip[0] != local {
				t.Fatalf("Expected loopback IP, got %s\n", ip[0])
			}
			if !tc.local && ip[0] == local {
				t.Fatalf("Expected non loopback IP, got %s\n", ip[0])
			}
		})
	}
}
