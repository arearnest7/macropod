package main

import (
    "encoding/csv"
    "fmt"
    "os"
    "strconv"
    "sort"
)

func main() {
    if len(os.Args) <= 1 {
        fmt.Println("go run process-metrics [results.csv] [metrics1] ... [metricsn]")
    } else {
        var records [][][]string
        for _, file_name := range os.Args[2:] {
            file, _ := os.Open(file_name)
            reader := csv.NewReader(file)
            record, _ := reader.ReadAll()
            records = append(records, record)
        }
        metrics_total := make(map[string][]string)
        var keys []string
        for _, record := range records {
            for _, line := range record[1:] {
                _, exists := metrics_total[line[0]]
                if exists {
                    for i, entry := range line[1:] {
                        temp1, _ := strconv.ParseFloat(metrics_total[line[0]][i], 64)
                        temp2, _ := strconv.ParseFloat(entry, 64)
                        metrics_total[line[0]][i] = strconv.FormatFloat(temp1 + temp2, 'f', -1, 64)
                    }
                } else {
                    keys = append(keys, line[0])
                    metrics_total[line[0]] = line[1:]
                }
            }
        }
        sort.Strings(keys)
        r, _ := os.Create(os.Args[1])
        defer r.Close()
        results := csv.NewWriter(r)
        defer results.Flush()
        headers := []string{"timestamp", "user", "idle", "iowait", "steal", "guest", "current", "min", "max", "cpu_load_1", "cpu_load_5", "cpu_load_15", "total", "available", "used", "free", "active", "inactive", "bytes_sent", "bytes_recv", "packets_sent", "packets_recv"}
        results.Write(headers)
        for _, timestamp := range keys {
            line := metrics_total[timestamp]
            var output []string
            output = append(output, timestamp)
            for _, entry := range line {
                output = append(output, entry)
            }
            results.Write(output)
        }
    }
}
