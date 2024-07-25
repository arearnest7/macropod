package main

import (
    "encoding/csv"
    "fmt"
    "os"
    "strconv"
    "sort"
    "time"
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
        init_time, _ := time.Parse("2006-01-02 15:04:05 UTC", records[0][1][0])
        for j, record := range records {
            var offset time.Duration
            prev := make([]string, 0)
            if j != 0 {
                metric_init, _ := time.Parse("2006-01-02 15:04:05 UTC", record[1][0])
                offset = init_time.Sub(metric_init)
            }
            for _, line := range record[1:] {
                metric_timestamp, _ := time.Parse("2006-01-02 15:04:05 UTC", line[0])
                metric_timestamp = metric_timestamp.Add(offset)
                timestamp_s := metric_timestamp.String()
                _, exists := metrics_total[timestamp_s]
                prev_set := false
                for _, prev_entry := range prev {
                    if prev_entry == timestamp_s {
                        prev_set = true
                        break
                    }
                }
                if exists && !prev_set {
                    for i, entry := range line[1:4] {
                        temp1, _ := strconv.ParseFloat(metrics_total[timestamp_s][i], 64)
                        temp2, _ := strconv.ParseFloat(entry, 64)
                        metrics_total[timestamp_s][i] = strconv.FormatFloat(temp1 + temp2, 'f', -1, 64)
                        metrics_total[timestamp_s] = append(metrics_total[timestamp_s], entry)
                    }
                    prev = append(prev, timestamp_s)
                } else if !prev_set {
                    keys = append(keys, timestamp_s)
                    metrics_total[timestamp_s] = line[1:4]
                    for range(j) {
                        metrics_total[timestamp_s] = append(metrics_total[timestamp_s], "0")
                        metrics_total[timestamp_s] = append(metrics_total[timestamp_s], "0")
                        metrics_total[timestamp_s] = append(metrics_total[timestamp_s], "0")
                    }
                    for _, entry := range line[1:4] {
                        metrics_total[timestamp_s] = append(metrics_total[timestamp_s], entry)
                    }
                    prev = append(prev, timestamp_s)
                }
            }
        }
        sort.Strings(keys)
        r, _ := os.Create(os.Args[1])
        defer r.Close()
        results := csv.NewWriter(r)
        defer results.Flush()
        headers := []string{"timestamp", "cpu_total_wo_idle", "used_memory_total", "bytes_sent_total"}
        for i := range(len(os.Args[2:])) {
            headers = append(headers, "cpu_total_wo_idle_" + strconv.Itoa(i+1))
            headers = append(headers, "user_memory_" + strconv.Itoa(i+1))
            headers = append(headers, "bytes_sent_" + strconv.Itoa(i+1))
        }
        results.Write(headers)
        for _, timestamp := range keys {
            line := metrics_total[timestamp]
            if len(line) < 3 + (3*len(records)) {
                for range(3 + (3*len(records)) - len(line)) {
                    line = append(line, "0")
                }
            }
            var output []string
            output = append(output, timestamp)
            for _, entry := range line {
                output = append(output, entry)
            }
            results.Write(output)
        }
        results.Flush()
    }
}
