package main

import (
    "encoding/csv"
    "fmt"
    "os"
    "time"
    "strconv"
    "github.com/mackerelio/go-osstat/cpu"
    "github.com/mackerelio/go-osstat/memory"
    "github.com/mackerelio/go-osstat/network"
)

func main() {
    if len(os.Args) != 3 {
        fmt.Println("go run metrics.go [net interface] [metrics.csv]")
    } else {
        o, _ := os.Create(os.Args[2])
        writer := csv.NewWriter(o)
        headers := []string{"timestamp", "cpu_total", "used_memory", "bytes_sent"}
        writer.Write(headers)
        normalization_set := false
        var bytes_sent_norm float64
        for true {
            cpu, _ := cpu.Get()
            mem, _ := memory.Get()
            nets, _ := network.Get()
            var net network.Stats
            for _, n := range nets {
                if n.Name == os.Args[1] {
                    net = n
                    break
                }
            }
            if !normalization_set {
                normalization_set = true
                bytes_sent_norm = float64(net.TxBytes)
            }
            c := make([]string, 0)
            c = append(c, time.Now().UTC().Format("2006-01-02 15:04:05 UTC"))
            c = append(c, strconv.FormatFloat(float64(cpu.User), 'f', -1, 64))
            c = append(c, strconv.FormatFloat(float64(mem.Used), 'f', -1, 64))
            c = append(c, strconv.FormatFloat(float64(net.TxBytes) - bytes_sent_norm, 'f', -1, 64))
            writer.Write(c)
            writer.Flush()
            bytes_sent_norm = float64(net.TxBytes)
            time.Sleep(1 * time.Second)
        }
    }
}
