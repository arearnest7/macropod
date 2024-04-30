package main

import (
	"github.com/go-ping/ping"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	"fmt"
)
var maxRTT = 0.0
func GetRTTHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "%f", maxRTT)
}

func ReadPingResults() {

	for {
	namespace := os.Getenv("NAMESPACE")
	endpoints := os.Getenv("ENDPOINTS")
	log.Print(endpoints)
	endpointList := strings.Split(endpoints, ",")
	log.Print(len(endpointList))
	if len(endpointList) == 0 {
		maxRTT = 0.0
	}
	total_max_rtt := 0.0
	for _, endpoint := range endpointList {
		if endpoint == ""{
			maxRTT = 0.0
			continue
		}
		serviceEndpoint := endpoint + "." + namespace + "." + "svc.cluster.local"
		log.Print(serviceEndpoint)
		pinger, err := ping.NewPinger(serviceEndpoint)
		if err != nil {
			maxRTT = 0.0
		}
		pinger.Count = 3
		pinger.Timeout = time.Second * 10
        pinger.OnRecv = func(pkt *ping.Packet) {
            log.Printf("Received packet from %s: time=%v", pkt.IPAddr, pkt.Rtt)
        }

        pinger.OnFinish = func(stats *ping.Statistics) {
            log.Printf("Ping statistics for %s: %+v", serviceEndpoint, stats)
            total_max_rtt += stats.MaxRtt.Seconds()
        }

        pinger.Run()
	}
	maxRTT = total_max_rtt
	time.Sleep(20)
}

}

func main() {
	go ReadPingResults()
	http.HandleFunc("/get_rtt", GetRTTHandler)
	log.Println("Server listening on port 8003...")
	http.ListenAndServe(":8003", nil)
}
