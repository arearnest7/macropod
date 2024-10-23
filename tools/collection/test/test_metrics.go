package main

import (
    "encoding/csv"
    "fmt"
    "os"
    "time"
    "strconv"
    "github.com/mackerelio/go-osstat/cpu"
    "github.com/mackerelio/go-osstat/memory"

)

func main() {
    if len(os.Args) != 3 {
        fmt.Println("go run metrics.go [net interface] [metrics.csv]")
    } else {
        o, _ := os.Create(os.Args[2])
        writer := csv.NewWriter(o)
        headers := []string{"timestamp", "cpu_steal", "used_memory", "bytes_sent", "cpu_user", "cpu_nice", "cpu_system", "cpu_idle", "cpu_iowait", "cpu_irq", "cpu_softirq", "cpu_steal", "cpu_guest", "cpu_guestnice", "cpu_total"}
        writer.Write(headers)
        for true {
            cpu, _ := cpu.Get()
            mem, _ := memory.Get()
            c := make([]string, 0)
            c = append(c, time.Now().UTC().Format("2006-01-02 15:04:05 UTC"))
            c = append(c, strconv.FormatFloat(float64(cpu.Steal), 'f', -1, 64))
            fmt.Print(cpu.Steal)
            c = append(c, strconv.FormatFloat(float64(mem.Used), 'f', -1, 64))
            c = append(c, strconv.FormatFloat(float64(cpu.Iowait), 'f', -1, 64))
            c = append(c, strconv.FormatFloat(float64(cpu.Irq), 'f', -1, 64))
            c = append(c, strconv.FormatFloat(float64(cpu.Softirq), 'f', -1, 64))
            c = append(c, strconv.FormatFloat(float64(cpu.Steal), 'f', -1, 64))
            c = append(c, strconv.FormatFloat(float64(cpu.Guest), 'f', -1, 64))
            c = append(c, strconv.FormatFloat(float64(cpu.GuestNice), 'f', -1, 64))
            writer.Write(c)
            writer.Flush()
            time.Sleep(1 * time.Second)
        }
    }
}