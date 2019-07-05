package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
)

const configFile = "./config.json"
const version = "v1.1.0"

type configT struct {
	FullChain string
	PrivKey   string
	Local     bool
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
	var con configT
	tlsOK := con.loadConfig()
	var out outPutT
	out.getPrivateIP(con.Local, tlsOK)

	if tlsOK {
		out.message = fmt.Sprintf("TLS webTest %s\nTLS Certs loaded - running over https\n", version)
		fmt.Printf(out.message)
		fmt.Printf("Server IP: %s\n", out.serverIP)
		err := http.ListenAndServeTLS(out.serverIP, con.FullChain, con.PrivKey, out.handler())
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

func (o *outPutT) getPrivateIP(local, tlsOK bool) {
	var ip string
	if !local {
		conn, err := net.Dial("udp", "8.8.8.8:80")
		defer conn.Close()
		if err != nil {
			log.Printf("No internet, local only: %v\n", err)
			ip = "127.0.0.1:8080"
		} else {
			port := "80"
			if tlsOK {
				port = "443"
			}
			localIP := conn.LocalAddr().(*net.UDPAddr)
			ip = fmt.Sprintf("%v:%s", localIP.IP, port)
		}
	} else {
		ip = "127.0.0.1:8080"
	}
	o.serverIP = ip
}

func (c *configT) loadConfig() bool {
	ok := true
	f, err := os.Open(configFile)
	if err != nil {
		log.Printf("Failed to load TLS certs - no https for you!\n%v\n", err)
		ok = false
	}
	defer f.Close()

	tlsJSON := json.NewDecoder(f)
	if err = tlsJSON.Decode(&c); err != nil {
		log.Printf("Failed to decode TLS certs JSON - no https for you!\n%v\n", err)
		ok = false
	}
	// get if tls certs exist on server
	if _, err := os.Stat(c.FullChain); err != nil {
		log.Printf("Failed to find FullChain cert - no https for you!\n%v\n", err)
		ok = false
	}
	if _, err := os.Stat(c.PrivKey); err != nil {
		log.Printf("Failed to find Private Key - no https for you!\n%v\n", err)
		ok = false
	}
	return ok
}
