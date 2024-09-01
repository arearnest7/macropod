package main

import (
    "encoding/csv"
    "bufio"
    "fmt"
    "os"
    "strconv"
    "strings"
    "time"
    "github.com/montanaflynn/stats"
)

func main() {
    if len(os.Args) <= 3 {
        fmt.Println("go run process-stats [results.csv] [metrics.csv] [directory] [logs1.csv;file1.out;...;filen.out] ... [logsn.csv;file1.out;...;filen.out]")
    } else {
        r, _ := os.Create(os.Args[1])
        defer r.Close()
        results := csv.NewWriter(r)
        defer results.Flush()
        metrics_file, _ := os.Open(os.Args[2])
        metrics_reader := csv.NewReader(metrics_file)
        metrics_reader.FieldsPerRecord = -1
        metrics_record, _ := metrics_reader.ReadAll()
        directory := os.Args[3]
        var tags []string
        tags = append(tags, "Workflow")
        tags = append(tags, "E2E Workflow Latency median")
        tags = append(tags, "E2E Workflow Latency 99 percentile")
        tags = append(tags, "Peak CPU Contention")
        tags = append(tags, "Peak Used Memory")
        tags = append(tags, "Peak Bytes Sent")
        for i := range((len(metrics_record[0]) - 4) / 3) {
            tags = append(tags, "Peak Used CPU " + strconv.Itoa(i+1))
            tags = append(tags, "Peak Used Memory " + strconv.Itoa(i+1))
            tags = append(tags, "Peak Bytes Sent " + strconv.Itoa(i+1))
        }
        results.Write(tags)
        for _, latency_pair := range os.Args[4:] {
            s := strings.Split(latency_pair, ";")
            latency_name := s[0]
            out_files := s[1:]
            latency_file, _ := os.Open(latency_name)
            latency_reader := csv.NewReader(latency_file)
            latency_record, _ := latency_reader.ReadAll()
            for _, out_file := range out_files {
                timestamp_file, _ := os.Open(directory + "/" + out_file)
                timestamp_reader := bufio.NewScanner(timestamp_file)
                timestamps := make([]time.Time, 0)
                for timestamp_reader.Scan() {
                    line := timestamp_reader.Text()
                    timestamp, _ := time.Parse("2006-01-02 15:04:05.000000 UTC", line)
                    if !timestamp.IsZero() {
                        timestamps = append(timestamps, timestamp)
                    }
                }
                for idx, wf_prefix := range [3]string{"kn-full", "kn-original", "kn-partial"} {
                    wf_name := wf_prefix + "-" + out_file[:len(out_file)-4]
                    var median string
                    var percentile99 string
                    var peaks []string
                    var wf_latency []float64
                    var metrics [][]float64
                    timestamp_start := timestamps[idx]
                    timestamp_end := timestamps[idx+1]
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
}
