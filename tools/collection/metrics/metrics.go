package main

import (
    "encoding/csv"
    "fmt"
    "os"
    "time"
    "strconv"
    "github.com/mackerelio/go-osstat/loadavg"
    "github.com/mackerelio/go-osstat/memory"
    "github.com/mackerelio/go-osstat/network"
)

func main() {
    if len(os.Args) != 3 {
        fmt.Println("go run metrics.go [net interface] [metrics.csv]")
    } else {
        o, _ := os.Create(os.Args[2])
        writer := csv.NewWriter(o)
        headers := []string{"timestamp", "cpu_load_1", "used_memory", "bytes_sent"}
        writer.Write(headers)
        normalization_set := false
        var load_avg_norm float64
        var mem_used_norm float64
        var bytes_sent_norm float64
        for true {
            load_avg, _ := loadavg.Get()
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
                load_avg_norm = load_avg.Loadavg1
                mem_used_norm = float64(mem.Used)
                bytes_sent_norm = float64(net.TxBytes)
            }
            c := make([]string, 0)
            c = append(c, time.Now().UTC().Format("2006-01-02 15:04:05 UTC"))
            c = append(c, strconv.FormatFloat(load_avg.Loadavg1 - load_avg_norm, 'f', -1, 64))
            c = append(c, strconv.FormatFloat(float64(mem.Used) - mem_used_norm, 'f', -1, 64))
            c = append(c, strconv.FormatFloat(float64(net.TxBytes) - bytes_sent_norm, 'f', -1, 64))
            writer.Write(c)
            writer.Flush()
            time.Sleep(1 * time.Second)
        }
    }
}
