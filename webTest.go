package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
)

const ipFile = "./outboundIP.json"
const tlsNameFile = "./tlsName.json"
const version = "v1.0.2"

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

	var tls tlsT
	live := tls.loadTLS()

	if live {
		out.message = fmt.Sprintf("TLS webTest %s\nTLS Certs loaded - running over https\n", version)
		fmt.Printf(out.message)
		fmt.Printf("Server IP: %s\n", out.serverIP)
		err := http.ListenAndServeTLS(out.serverIP, tls.FullChain, tls.PrivKey, out.handler())
		if err != nil {
			log.Fatal(err)
		}
	} else {
		out.message = fmt.Sprintf("TLS webTest %s\nNo TLS Certs loaded - running over http\n", version)
		fmt.Printf(out.message)
		fmt.Printf("Server IP: %s\n", out.serverIP)
		err := http.ListenAndServe(out.serverIP, out.handler())
		if err != nil {
			log.Fatal(err)
		}
	}
}

func (o *outPutT) handler() http.Handler {
	r := http.NewServeMux()
	r.HandleFunc("/", o.testPage)
	return r
}

func (o *outPutT) testPage(w http.ResponseWriter, r *http.Request) {
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
	ok := true
	f, err := os.Open(tlsNameFile)
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
	// get if tls certs exist on server
	if _, err := os.Stat(t.FullChain); err != nil {
		log.Printf("Failed to find FullChain cert - no https for you!\n%v\n", err)
		ok = false
	}
	if _, err := os.Stat(t.PrivKey); err != nil {
		log.Printf("Failed to find Private Key - no https for you!\n%v\n", err)
		ok = false
	}
	return ok
}
