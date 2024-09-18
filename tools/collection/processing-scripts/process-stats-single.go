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
            var tags []string
            tags = append(tags, "Workflow")
            tags = append(tags, "E2E Workflow Latency median")
            tags = append(tags, "E2E Workflow Latency 99 percentile")
            for _, tag := range metrics_record[0][1:] {
                tags = append(tags, "Peak " + tag)
            }
            s := strings.Split(wf_raw, ";")
            wf_name := s[0]
            log_name := s[1]
            log_file, _ := os.Open(log_name)
            log_reader := csv.NewReader(log_file)
            log_record, _ := log_reader.ReadAll()
            if len(log_record) > 0 {
                temp := log_record[0]
                for _, tag := range temp {
                    tags = append(tags, "Peak " + tag)
                }
                var median string
                var percentile99 string
                var peaks []string
                var wf_latency []float64
                var metrics [][]float64
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
                    timestamp, _ := time.Parse("2006-01-02 15:04:05 +0000 UTC", line[0])
                    var metrics_line []float64
                    for _, v := range line[1:] {
                        v_p, _ := strconv.ParseFloat(v, 64)
                        metrics_line = append(metrics_line, v_p)
                    }
                    if timestamp.After(timestamp_start) && timestamp.Before(timestamp_end) {
                        metrics = append(metrics, metrics_line)
                    }
                }
                if len(wf_latency) > 0 {
                    len_idx := len(wf_latency)
                    for i, latency := range wf_latency {
                        if latency == 0 {
                            len_idx = i
                            break
                        }
                    }
                    m, _ := stats.Median(wf_latency[:len_idx])
                    median = strconv.FormatFloat(m, 'f', -1, 64)
                    p, _ := stats.Percentile(wf_latency[:len_idx], 99)
                    percentile99 = strconv.FormatFloat(p, 'f', -1, 64)
                } else {
                    median = "0"
                    percentile99 = "0"
                }
                for i := range(len(metrics[0])) {
                    var values []float64
                    for _, entry := range metrics {
                        values = append(values, entry[i])
                    }
                    peak, _ := stats.Max(values)
                    peaks = append(peaks, strconv.FormatFloat(peak, 'f', -1, 64))
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
                    peak, _ := stats.Max(temp2)
                    peaks = append(peaks, strconv.FormatFloat(peak, 'f', -1, 64))
                }
                results.Write(tags)
                stats_line := make([]string, 0)
                stats_line = append(stats_line, wf_name)
                stats_line = append(stats_line, median)
                stats_line = append(stats_line, percentile99)
                for _, peak := range peaks {
                    stats_line = append(stats_line, peak)
                }
                results.Write(stats_line)
            }
        }
    }
}
