package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const ipFile = "./outboundIP.json"
const tlsNameFile = "./tlsName.json"

type tlsT struct {
	FullChain string
	PrivKey   string
}

type serverIPT struct {
	IP   string `json:"serverIP"`
	Port string `json:"serverPort"`
}

func main() {
	appIP, err := loadPrivateIP()
	if err != nil {
		msg := fmt.Sprintf("Error loading IP config - exiting %v\n", err)
		closeApp(msg, false)
	}

	server := &http.Server{
		Addr:         appIP,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	http.HandleFunc("/", tester)

	var tls tlsT
	live := tls.loadTLS()

	if live {
		out := fmt.Sprintf("TLS Certs loaded - running over https\n")
		fmt.Printf(out)
		fmt.Printf("Server IP: %s\n", appIP)
		log.Fatal(server.ListenAndServeTLS(tls.Fullchain, tls.PrivKey))
	} else {
		out := fmt.Sprintf("No TLS Certs loaded - running over http\n")
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

func (t *tlsT) loadTLS() bool {
	f, err := os.Open(tlsNameFile)
	ok := true
	if err != nil {
		log.Printf("Failed to load TLS certs - no https for you!\n%v\n", err)
		ok = false
	}
	defer f.Close()

	tlsJSON := json.NewDecoder(f)
	if err = tlsJSON.Decode(&t); err != nil {
		log.Printf("Failed to decode TLS certs JSON - no https for you!\n%v\n", err)
		ok = false
	}
	return ok
}

func closeApp(in string, save bool) {
	log.Printf(in)
	os.Exit(0)
}
