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
    "github.com/mackerelio/go-osstat/loadavg"
)

func main() {
    if len(os.Args) != 3 {
        fmt.Println("go run metrics.go [net interface] [metrics.csv]")
    } else {
        o, _ := os.Create(os.Args[2])
        writer := csv.NewWriter(o)
        headers := []string{"timestamp", "bytes_sent", "loadavg_1", "memory_used", "memory_buffers", "memory_cached", "memory_free", "memory_available", "memory_active", "memory_inactive", "memory_swap_total", "memory_swap_used", "memory_swap_cached", "memory_swap_free", "memory_mapped", "memory_shmem", "memory_slab", "memory_page_tables", "memory_committed", "memory_v_malloc_used", "cpu_user", "cpu_nice", "cpu_system", "cpu_idle", "cpu_iowait", "cpu_irq", "cpu_softirq", "cpu_steal", "cpu_guest", "cpu_guestnice", "cpu_total"}
        writer.Write(headers)
        normalization_set := false
        var bytes_sent_norm float64
        for true {
            cpu, _ := cpu.Get()
            mem, _ := memory.Get()
            nets, _ := network.Get()
            loadavg, _ := loadavg.Get()
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
            c = append(c, strconv.FormatFloat(float64(net.TxBytes) - bytes_sent_norm, 'f', -1, 64))
            c = append(c, strconv.FormatFloat(float64(loadavg.Loadavg1), 'f', -1, 64))
            c = append(c, strconv.FormatFloat(float64(mem.Used), 'f', -1, 64))
            c = append(c, strconv.FormatFloat(float64(mem.Buffers), 'f', -1, 64))
            c = append(c, strconv.FormatFloat(float64(mem.Cached), 'f', -1, 64))
            c = append(c, strconv.FormatFloat(float64(mem.Free), 'f', -1, 64))
            c = append(c, strconv.FormatFloat(float64(mem.Available), 'f', -1, 64))
            c = append(c, strconv.FormatFloat(float64(mem.Active), 'f', -1, 64))
            c = append(c, strconv.FormatFloat(float64(mem.Inactive), 'f', -1, 64))
            c = append(c, strconv.FormatFloat(float64(mem.SwapTotal), 'f', -1, 64))
            c = append(c, strconv.FormatFloat(float64(mem.SwapUsed), 'f', -1, 64))
            c = append(c, strconv.FormatFloat(float64(mem.SwapCached), 'f', -1, 64))
            c = append(c, strconv.FormatFloat(float64(mem.SwapFree), 'f', -1, 64))
            c = append(c, strconv.FormatFloat(float64(mem.Mapped), 'f', -1, 64))
            c = append(c, strconv.FormatFloat(float64(mem.Shmem), 'f', -1, 64))
            c = append(c, strconv.FormatFloat(float64(mem.Slab), 'f', -1, 64))
            c = append(c, strconv.FormatFloat(float64(mem.PageTables), 'f', -1, 64))
            c = append(c, strconv.FormatFloat(float64(mem.Committed), 'f', -1, 64))
            c = append(c, strconv.FormatFloat(float64(mem.VmallocUsed), 'f', -1, 64))
            c = append(c, strconv.FormatFloat(float64(cpu.User), 'f', -1, 64))
            c = append(c, strconv.FormatFloat(float64(cpu.Nice), 'f', -1, 64))
            c = append(c, strconv.FormatFloat(float64(cpu.System), 'f', -1, 64))
            c = append(c, strconv.FormatFloat(float64(cpu.Idle), 'f', -1, 64))
            c = append(c, strconv.FormatFloat(float64(cpu.Iowait), 'f', -1, 64))
            c = append(c, strconv.FormatFloat(float64(cpu.Irq), 'f', -1, 64))
            c = append(c, strconv.FormatFloat(float64(cpu.Softirq), 'f', -1, 64))
            c = append(c, strconv.FormatFloat(float64(cpu.Steal), 'f', -1, 64))
            c = append(c, strconv.FormatFloat(float64(cpu.Guest), 'f', -1, 64))
            c = append(c, strconv.FormatFloat(float64(cpu.GuestNice), 'f', -1, 64))
            c = append(c, strconv.FormatFloat(float64(cpu.Total), 'f', -1, 64))
            writer.Write(c)
            writer.Flush()
            bytes_sent_norm = float64(net.TxBytes)
            time.Sleep(1 * time.Second)
        }
    }
}
