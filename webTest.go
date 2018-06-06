package main

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
	"os"
	"time"
)

var out string
var appIP = "40.115.40.151:80"
var ipFile = "./outboundIP.json"
var fullchain = "/etc/letsencrypt/live/webapp.millasays.com/fullchain.pem"
var privKey = "/etc/letsencrypt/live/webapp.millasays.com/privkey.pem"

type tlsT struct {
	Fullchain string
	PrivKey   string
}

type serverIPT struct {
	IP      string
	Version string
}

func main() {

	getIP, err := loadPrivateIP()
	if err != nil {
		msg := fmt.Sprintf("Error loading IP config - exiting %v\n", err)
		CloseApp(msg, false)
	}

	appIP = getIP

	server := &http.Server{
		Addr:         appIP,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	http.HandleFunc("/", tester)

	TLS, live := checkTLSExists()

	if live {
		out = fmt.Sprintf("TLS Certs loaded - running over https\n")
		fmt.Printf(out)
		fmt.Printf("Server IP: %s\n", appIP)
		log.Fatal(server.ListenAndServeTLS(TLS.Fullchain, TLS.PrivKey))
	} else {
		out = fmt.Sprintf("No TLS Certs loaded - running over http\n")
		fmt.Printf(out)
		fmt.Printf("Server IP: %s\n", appIP)
		log.Fatal(server.ListenAndServe())
	}

}

func tester(w http.ResponseWriter, r *http.Request) {
	clientIP := r.URL.Path
	fmt.Fprintf(w, "%q", html.EscapeString(out))
	fmt.Printf("Inbound from     : %s\n", clientIP)
	fmt.Printf("Response from    : %s\n", appIP)
	fmt.Fprintf(w, "Inbound from     : %q\nResponse from    : %q", html.EscapeString(clientIP), html.EscapeString(appIP))
}

func checkTLSExists() (tlsT, bool) {
	var TLS tlsT
	var exists bool

	if _, err := os.Stat("./live.txt"); err != nil {
		exists = false
	} else {
		exists = true
		TLS.Fullchain = fullchain
		TLS.PrivKey = privKey
	}
	return TLS, exists
}

func loadPrivateIP() (string, error) {
	var ipOut string
	var ipIn serverIPT

	f, err := os.Open(ipFile)
	if err != nil {
		return "", err
	}
	defer f.Close()

	ipJSON := json.NewDecoder(f)
	if err = ipJSON.Decode(&ipIn); err != nil {
		return "", err
	}
	ipOut = fmt.Sprintf("%s:%s", ipIn.serverIP, ipIn.serverPort)
	return ipOut, nil
}
