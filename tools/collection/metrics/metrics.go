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
            c := make([]string, 0)
            c = append(c, time.Now().UTC().Format("2006-01-02 15:04:05 UTC"))
            c = append(c, strconv.FormatFloat(load_avg.Loadavg1, 'f', -1, 64))
            c = append(c, strconv.FormatFloat(float64(mem.Used), 'f', -1, 64))
            c = append(c, strconv.FormatFloat(float64(net.TxBytes), 'f', -1, 64))
            writer.Write(c)
            writer.Flush()
            time.Sleep(1 * time.Second)
        }
    }
}
