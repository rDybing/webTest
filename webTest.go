package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
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

type outPutT struct {
	message  string
	clientIP string
	serverIP string
}

func main() {
	var out outPutT
	err := out.loadPrivateIP()

	if err != nil {
		log.Fatalf("Error loading IP config - exiting %v\n", err)
	}

	server := &http.Server{
		Addr:         out.serverIP,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	http.HandleFunc("/", out.tester)

	var tls tlsT
	live := tls.loadTLS()

	if live {
		out.message = fmt.Sprintf("TLS Certs loaded - running over https\n")
		fmt.Printf(out.message)
		fmt.Printf("Server IP: %s\n", out.serverIP)
		log.Fatal(server.ListenAndServeTLS(tls.FullChain, tls.PrivKey))
	} else {
		out.message = fmt.Sprintf("No TLS Certs loaded - running over http\n")
		fmt.Printf(out.message)
		fmt.Printf("Server IP: %s\n", out.serverIP)
		log.Fatal(server.ListenAndServe())
	}

}

func (o *outPutT) tester(w http.ResponseWriter, r *http.Request) {
	ip, port, _ := net.SplitHostPort(r.RemoteAddr)
	o.clientIP = fmt.Sprintf("%s:%s", ip, port)
	fmt.Fprintf(w, "%s", o.message)
	fmt.Printf("Inbound from     : %s\n", o.clientIP)
	fmt.Printf("Response from    : %s\n", o.serverIP)
	fmt.Fprintf(w, "Inbound from     : %s\nResponse from    : %s", o.clientIP, o.serverIP)
}

func (o *outPutT) loadPrivateIP() error {
	var ipIn serverIPT

	f, err := os.Open(ipFile)
	if err != nil {
		return err
	}
	defer f.Close()

	ipJSON := json.NewDecoder(f)
	if err := ipJSON.Decode(&ipIn); err != nil {
		return err
	}
	o.serverIP = fmt.Sprintf("%s:%s", ipIn.IP, ipIn.Port)
	return nil
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
