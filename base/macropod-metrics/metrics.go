package main

import (
    "encoding/csv"
    "fmt"
    "os"
    "time"
    "strconv"
    "github.com/mackerelio/go-osstat/cpu"
    "github.com/mackerelio.go-osstat/disk"
    "github.com/mackerelio/go-osstat/loadavg"
    "github.com/mackerelio/go-osstat/memory"
    "github.com/mackerelio/go-osstat/network"
    "github.com/mackerelio/go-osstat/uptime"
)

func main() {
    if len(os.Args) != 3 {
        fmt.Println("go run metrics.go [interface] [metrics.csv]")
    } else {
        o, _ := os.Create(os.Args[2])
        writer := csv.NewWriter(o)
        headers := []string{"timestamp"}
        writer.Write(headers)
        for true {
            c, _ := cpu.Get()
            d, _ := disk.Get()
            l, _ := loadavg.Get()
            var n network.Stats
            nets, _ := network.Get()
            m, _ := memory.Get()
            l, _ := loadavg.Get()
            u, _ := uptime.get()
            for _, n := range nets {
                if n.Name == os.Args[1] {
                    net = n
                    break
                }
            }
            line := make([]string, 0)

            // timestamp
            line = append(line, time.Now().UTC().Format("2006-01-02 15:04:05 UTC"))

            // cpu
            line = append(line, strconv.FormatFloat(float64(c.User), 'f', -1, 64))
            line = append(line, strconv.FormatFloat(float64(c.Nice), 'f', -1, 64))
            line = append(line, strconv.FormatFloat(float64(c.System), 'f', -1, 64))
            line = append(line, strconv.FormatFloat(float64(c.Idle), 'f', -1, 64))
            line = append(line, strconv.FormatFloat(float64(c.Iowait), 'f', -1, 64))
            line = append(line, strconv.FormatFloat(float64(c.Irq), 'f', -1, 64))
            line = append(line, strconv.FormatFloat(float64(c.Softirq), 'f', -1, 64))
            line = append(line, strconv.FormatFloat(float64(c.Steal), 'f', -1, 64))
            line = append(line, strconv.FormatFloat(float64(c.Guest), 'f', -1, 64))
            line = append(line, strconv.FormatFloat(float64(c.GuestNice), 'f', -1, 64))
            line = append(line, strconv.FormatFloat(float64(c.Total), 'f', -1, 64))

            // memory
            line = append(line, strconv.FormatFloat(float64(m.Total), 'f', -1, 64))
            line = append(line, strconv.FormatFloat(float64(m.Used), 'f', -1, 64))
            line = append(line, strconv.FormatFloat(float64(m.Buffers), 'f', -1, 64))
            line = append(line, strconv.FormatFloat(float64(m.Cached), 'f', -1, 64))
            line = append(line, strconv.FormatFloat(float64(m.Free), 'f', -1, 64))
            line = append(line, strconv.FormatFloat(float64(m.Available), 'f', -1, 64))
            line = append(line, strconv.FormatFloat(float64(m.Active), 'f', -1, 64))
            line = append(line, strconv.FormatFloat(float64(m.Inactive), 'f', -1, 64))
            line = append(line, strconv.FormatFloat(float64(m.SwapTotal), 'f', -1, 64))
            line = append(line, strconv.FormatFloat(float64(m.SwapUsed), 'f', -1, 64))
            line = append(line, strconv.FormatFloat(float64(m.SwapCached), 'f', -1, 64))
            line = append(line, strconv.FormatFloat(float64(m.SwapFree), 'f', -1, 64))
            line = append(line, strconv.FormatFloat(float64(m.Mapped), 'f', -1, 64))
            line = append(line, strconv.FormatFloat(float64(m.Shmem), 'f', -1, 64))
            line = append(line, strconv.FormatFloat(float64(m.Slab), 'f', -1, 64))
            line = append(line, strconv.FormatFloat(float64(m.PageTables), 'f', -1, 64))
            line = append(line, strconv.FormatFloat(float64(m.Committed), 'f', -1, 64))
            line = append(line, strconv.FormatFloat(float64(m.VmallocUsed), 'f', -1, 64))

            // loadavg
            line = append(line, strconv.FormatFloat(float64(l.Loadavg1), 'f', -1, 64))
            line = append(line, strconv.FormatFloat(float64(l.Loadavg5), 'f', -1, 64))
            line = append(line, strconv.FormatFloat(float64(l.Loadavg15), 'f', -1, 64))

            // uptime
            line = append(line, strconv.FormatFloat(float64(u.Seconds()), 'f', -1, 64))

            // network
            

            // disk
            

            writer.Write(line)
            writer.Flush()
            time.Sleep(1 * time.Second)
        }
    }
}
