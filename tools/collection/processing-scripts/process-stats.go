package main

import (
    "encoding/csv"
    "fmt"
    "os"
    "strconv"
    "strings"
    "time"
    "github.com/montanaflynn/stats"
)

func main() {
    if len(os.Args) <= 3 {
        fmt.Println("go run process-stats [results.csv] [latency.csv] [metrics.csv] [wf_name1;logs1.csv] ... [wf_namen;logsn.csv]")
    } else {
        r, _ := os.Create(os.Args[1])
        defer r.Close()
        results := csv.NewWriter(r)
        defer results.Flush()
        latency_file, _ := os.Open(os.Args[2])
        latency_reader := csv.NewReader(latency_file)
        latency_record, _ := latency_reader.ReadAll()
        metrics_file, _ := os.Open(os.Args[3])
        metrics_reader := csv.NewReader(metrics_file)
        metrics_record, _ := metrics_reader.ReadAll()
        for _, wf_raw := range os.Args[4:] {
            s := strings.Split(wf_raw, ";")
            wf_name := s[0]
            log_name := s[1]
            log_file, _ := os.Open(log_name)
            log_reader := csv.NewReader(log_file)
            log_record, _ := log_reader.ReadAll()
            if len(log_record) > 0 {
                temp := log_record[0]
                tags := []string{wf_name, "E2E Workflow Latency", "Peak CPU Contention", "Peak Used Memory", "Peak Bytes Sent"}
                for _, tag := range temp {
                    tags = append(tags, tag)
                }
                results.Write(tags)
                var median []string
                var percentile99 []string
                median = append(median, "median")
                percentile99 = append(percentile99, "99 percentile")
                var wf_latency []float64
                var cpu []float64
                var memory []float64
                var bytes_sent []float64
                var log_record_nano [][]int64
                timestamp_name := s[1] + ".mts"
                timestamp_file, _ := os.Open(timestamp_name)
                timestamp_reader := csv.NewReader(timestamp_file)
                timestamp_record, _ := timestamp_reader.ReadAll()
                timestamp_start, _ := time.Parse("2006-01-02 15:04:05.000000 UTC", timestamp_record[0][0])
                timestamp_end, _ := time.Parse("2006-01-02 15:04:05.000000 UTC", timestamp_record[0][1])
                timestamp_end = timestamp_end.Add(time.Second * 2)
                var latency_idx int
                for i, entry := range latency_record[0] {
                    if entry == wf_name {
                        latency_idx = i
                        break
                    }
                }
                for _, line := range latency_record[1:] {
                    l, _ := strconv.ParseFloat(line[latency_idx], 64)
                    wf_latency = append(wf_latency, l)
                }
                for _, line := range metrics_record[1:] {
                    timestamp, _ := time.Parse("2006-01-02 15:04:05 UTC", line[0])
                    cpu_load_1, _ := strconv.ParseFloat(line[1], 64)
                    used_memory, _ := strconv.ParseFloat(line[2], 64)
                    bytes_sent_entry, _ := strconv.ParseFloat(line[3], 64)
                    if timestamp.After(timestamp_start) && timestamp.Before(timestamp_end) {
                       cpu = append(cpu, cpu_load_1)
                        memory = append(memory, used_memory)
                        bytes_sent = append(bytes_sent, bytes_sent_entry)
                    }
                }
                if len(wf_latency) > 0 {
                    m, _ := stats.Median(wf_latency)
                    median = append(median, strconv.FormatFloat(m, 'f', -1, 64))
                    p, _ := stats.Percentile(wf_latency, 99)
                    percentile99 = append(percentile99, strconv.FormatFloat(p, 'f', -1, 64))
                } else {
                    median = append(median, "0")
                    percentile99 = append(percentile99, "0")
                }
                if len(cpu) > 0 {
                    m, _ := stats.Median(cpu)
                    median = append(median, strconv.FormatFloat(m, 'f', -1, 64))
                    p, _ := stats.Percentile(cpu, 99)
                    percentile99 = append(percentile99, strconv.FormatFloat(p, 'f', -1, 64))
                } else {
                    median = append(median, "0")
                    percentile99 = append(percentile99, "0")
                }
                if len(memory) > 0 {
                    m, _ := stats.Median(memory)
                    median = append(median, strconv.FormatFloat(m, 'f', -1, 64))
                    p, _ := stats.Percentile(memory, 99)
                    percentile99 = append(percentile99, strconv.FormatFloat(p, 'f', -1, 64))
                } else {
                    median = append(median, "0")
                    percentile99 = append(percentile99, "0")
                }
                if len(bytes_sent) > 0 {
                    m, _ := stats.Median(bytes_sent)
                    median = append(median, strconv.FormatFloat(m, 'f', -1, 64))
                    p, _ := stats.Percentile(bytes_sent, 99)
                    percentile99 = append(percentile99, strconv.FormatFloat(p, 'f', -1, 64))
                } else {
                    median = append(median, "0")
                    percentile99 = append(percentile99, "0")
                }
                for _, line := range log_record[1:] {
                    var newline []int64
                    for _, entry := range line {
                        new_entry, _ := strconv.ParseInt(entry, 10, 64)
                        newline = append(newline, new_entry)
                    }
                    log_record_nano = append(log_record_nano, newline)
                }
                for i, _ := range temp {
                    var temp2 []float64
                    for _, line := range log_record_nano {
                        val := float64(line[i])
                        if val > 0 {
                            temp2 = append(temp2, val)
                        } else {
                            temp2 = append(temp2, 0)
                        }
                    }
                    m, _ := stats.Median(temp2)
                    median = append(median, strconv.FormatFloat(m, 'f', -1, 64))
                    p, _ := stats.Percentile(temp2, 99)
                    percentile99 = append(percentile99, strconv.FormatFloat(p, 'f', -1, 64))
                }
                results.Write(median)
                results.Write(percentile99)
            }
        }
    }
}
