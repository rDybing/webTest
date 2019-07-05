package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

var out string
var appIP string
var ipFile = "./outboundIP.json"
var tlsNameFile = "./tlsName.json"

type tlsT struct {
	Fullchain string
	PrivKey   string
}

type serverIPT struct {
	IP   string `json:"serverIP"`
	Port string `json:"serverPort"`
}

func main() {

	getIP, err := loadPrivateIP()
	if err != nil {
		msg := fmt.Sprintf("Error loading IP config - exiting %v\n", err)
		closeApp(msg, false)
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
	fmt.Fprintf(w, "%s", out)
	fmt.Printf("Inbound from     : %s\n", clientIP)
	fmt.Printf("Response from    : %s\n", appIP)
	fmt.Fprintf(w, "Inbound from     : %s\nResponse from    : %s", clientIP, appIP)
}

func checkTLSExists() (tlsT, bool) {
	var TLS tlsT
	var exists bool

	if _, err := os.Stat("./live.txt"); err != nil {
		exists = false
	} else {
		exists = true
		
		TLS.Fullchain, TLS.PrivKey = loadTLSName()
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
	ipOut = fmt.Sprintf("%s:%s", ipIn.IP, ipIn.Port)
	return ipOut, nil
}

func loadTLSName() (string, string) {
	var tls tlsT

	f, err := os.Open(tlsNameFile)
	if err != nil {
		log.Fatalf("Failed to load TLS certs - exiting...\n%v\n", err)
		return "", ""
	}
	defer f.Close()

	tlsJSON := json.NewDecoder(f)
	if err = tlsJSON.Decode(&tls); err != nil {
		log.Fatalf("Failed to decode TLS certs JSON - exiting...\n%v\n", err)
		return "", ""
	}
	return tls.Fullchain, tls.PrivKey
}

func closeApp(in string, save bool) {
	log.Printf(in)
	os.Exit(0)
}
