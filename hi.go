package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
)

func main() {
	//Set listening port
	port := 9899
	if len(os.Args) > 1 {
		var err error
		port, err = strconv.Atoi(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
	}
	//Find default output IP and hostname
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}
	//Register handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ClientAddress := r.Header.Get("X-Real-Ip")
		if ClientAddress == "" {
			ClientAddress = r.Header.Get("X-Forwarded-For")
		}
		if ClientAddress == "" {
			ClientAddress = r.RemoteAddr
		}
		log.Println("Request from " + ClientAddress)
		w.Write([]byte(hostname + " " + localAddr.IP.String() + "\n"))
	})
	//Serve
	log.Println(hostname + " " + localAddr.IP.String() + " " + strconv.Itoa(port))
	log.Fatal(http.ListenAndServe("0.0.0.0:"+strconv.Itoa(port), nil))
}

