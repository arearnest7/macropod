package main

import (
	"context"
	"fmt"
	"github.com/go-ping/ping"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"os"
	"strconv"
    "encoding/json"
	"time"
)

type PingResult struct {
	Name            string            `json:"name"`
	Namespace       string            `json:"namespace"`
	IP              string            `json:"ip_address"`
	ReceivedPackets []string          `json:"received_packets"`
	Statistics      PingStatistics  `json:"statistics"`
}

type PingStatistics struct {
	PacketsTransmitted int     `json:"packets_transmitted"`
	PacketsReceived    int     `json:"packets_received"`
	PacketLoss         float64 `json:"packet_loss"`
	MinRtt             float64 `json:"min_rtt_seconds"`
	AvgRtt             float64 `json:"avg_rtt_seconds"`
	MaxRtt             float64 `json:"max_rtt_seconds"`
	StdDevRtt          float64 `json:"stddev_rtt_seconds"`
}

type OverallResult struct {
	PingResult []PingResult `json:"ping_result"`
}

func main() {
	intervalStr := os.Getenv("INTERVAL_SECONDS")
	intervalSec, err := strconv.Atoi(intervalStr)
	interval := time.Duration(intervalSec) * time.Second
	if err != nil {
		intervalSec = 10
	}

	config, err := rest.InClusterConfig()
	if err != nil {
		print("error")
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	namespace := os.Getenv("NAMESPACE")
    ip := os.Getenv("IP_SELF")
	for {
		pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
        result := OverallResult{}
		for _, pod := range pods.Items {
            pod_Result := PingResult{}
            if ip == pod.Status.PodIP{
                continue
            }
			fmt.Printf("Pod name: %s, Pod IP: %s\n", pod.Name, pod.Status.PodIP)
			pinger, err := ping.NewPinger(pod.Status.PodIP)
			if err != nil {
				continue
			}
            pod_Result = PingResult{
				Name:            pod.Name,
				Namespace:       pod.Namespace,
				IP:              pod.Status.PodIP,
				ReceivedPackets: []string{},
			}
			pinger.Count = 3

			pinger.Timeout = time.Second * 10

			pinger.OnRecv = func(pkt *ping.Packet) {
                packetInfo := fmt.Sprintf("Received packet from %s: time=%v", pkt.IPAddr, pkt.Rtt)
				pod_Result.ReceivedPackets = append(pod_Result.ReceivedPackets, packetInfo)
			}
			pinger.OnFinish = func(stats *ping.Statistics) {
                pod_Result.Statistics = PingStatistics{
                    PacketsTransmitted: stats.PacketsSent,
                    PacketsReceived:    stats.PacketsRecv,
                    PacketLoss:         stats.PacketLoss,
                    MinRtt:             stats.MinRtt.Seconds(),
                    AvgRtt:             stats.AvgRtt.Seconds(),
                    MaxRtt:             stats.MaxRtt.Seconds(),
                    StdDevRtt:          stats.StdDevRtt.Seconds(),
                }
			}

			fmt.Println("Pinging", pod.Status.PodIP)
			if err := pinger.Run(); err != nil {
				panic(err)
			}
            result.PingResult = append(result.PingResult,pod_Result)

		}

        jsonData, err := json.MarshalIndent(result, "", "    ")
		if err != nil {
			fmt.Printf("ERROR: Could not serialize results: %s\n", err.Error())
			continue
		}

		fileName := "ping_result.json"
		if err := os.WriteFile(fileName, jsonData, 0644); err != nil {
			fmt.Printf("ERROR: Could not write to file: %s\n", err.Error())
			continue 
		}

		fmt.Printf("Ping results saved to %s\n", fileName)
		time.Sleep(interval)
	}
}
