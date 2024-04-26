package main

import (
    "encoding/csv"
    "fmt"
    "os"
    "strconv"
    "strings"
    "time"
    "sort"
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
            temp := log_record[0]
            var tags []string
            tags = append(tags, wf_name)
            tags = append(tags, "E2E Workflow Latency")
            tags = append(tags, "Peak Used CPU")
            tags = append(tags, "Peak Used Memory")
            tags = append(tags, "Peak Bytes Sent")
            for _, tag := range temp {
                tags = append(tags, tag)
            }
            results.Write(tags)
            var avg []string
            var percentile99 []string
            avg = append(avg, "average")
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
            timestamp_start, _ := time.Parse("2006-01-02 15:04:05", timestamp_record[0][0])
            timestamp_end, _ := time.Parse("2006-01-02 15:04:05", timestamp_record[0][1])
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
                timestamp, _ := time.Parse("Jan 2 2006 15:04:05", line[0])
                cpu_load_1, _ := strconv.ParseFloat(line[9], 64)
                total_memory, _ := strconv.ParseFloat(line[12], 64)
                available_memory, _ := strconv.ParseFloat(line[13], 64)
                bytes_sent_entry, _ := strconv.ParseFloat(line[18], 64)
                if timestamp.After(timestamp_start) && timestamp.Before(timestamp_end) {
                    cpu = append(cpu, cpu_load_1)
                    memory = append(memory, total_memory - available_memory)
                    bytes_sent = append(bytes_sent, bytes_sent_entry)
                }
            }
            sort.Float64s(wf_latency)
            sum := 0.0
            for _, entry := range wf_latency {
                sum += entry
            }
            if len(wf_latency) > 0 {
                avg = append(avg, strconv.FormatFloat(sum / float64(len(wf_latency)), 'f', -1, 64))
                percentile_idx := int(float64(len(wf_latency)) * 0.99)
                percentile99 = append(percentile99, strconv.FormatFloat(wf_latency[percentile_idx], 'f', -1, 64))
            } else {
                avg = append(avg, "0")
                percentile99 = append(percentile99, "0")
            }
            sort.Float64s(cpu)
            sum = 0.0
            for _, entry := range cpu {
                sum += entry
            }
            if len(cpu) > 0 {
                avg = append(avg, strconv.FormatFloat(sum / float64(len(cpu)), 'f', -1, 64))
                percentile_idx := int(float64(len(cpu)) * 0.99)
                percentile99 = append(percentile99, strconv.FormatFloat(cpu[percentile_idx], 'f', -1, 64))
            } else {
                avg = append(avg, "0")
                percentile99 = append(percentile99, "0")
            }
            sort.Float64s(memory)
            sum = 0.0
            for _, entry := range memory {
                sum += entry
            }
            if len(memory) > 0 {
                avg = append(avg, strconv.FormatFloat(sum / float64(len(memory)), 'f', -1, 64))
                percentile_idx := int(float64(len(memory)) * 0.99)
                percentile99 = append(percentile99, strconv.FormatFloat(memory[percentile_idx], 'f', -1, 64))
            } else {
                avg = append(avg, "0")
                percentile99 = append(percentile99, "0")
            }
            sort.Float64s(bytes_sent)
            sum = 0.0
            for _, entry := range bytes_sent {
                sum += entry
            }
            if len(bytes_sent) > 0 {
                avg = append(avg, strconv.FormatFloat(sum / float64(len(bytes_sent)), 'f', -1, 64))
                percentile_idx := int(float64(len(bytes_sent)) * 0.99)
                percentile99 = append(percentile99, strconv.FormatFloat(bytes_sent[percentile_idx], 'f', -1, 64))
            } else {
                avg = append(avg, "0")
                percentile99 = append(percentile99, "0")
            }
            for _, line := range log_record[1:] {
                var newline []int64
                for _, entry := range line {
                    new_entry, _ := time.ParseDuration(entry)
                    newline = append(newline, new_entry.Nanoseconds())
                }
                log_record_nano = append(log_record_nano, newline)
            }
            for i, _ := range temp {
                var temp2 []float64
                sum = 0.0
                for _, line := range log_record_nano {
                    temp2 = append(temp2, float64(line[i]))
                    sum += float64(line[i])
                }
                sort.Float64s(temp2)
                avg = append(avg, strconv.FormatFloat(sum / float64(len(temp2)), 'f', -1, 64))
                percentile_idx := int(float64(len(temp2)) * 0.99)
                percentile99 = append(percentile99, strconv.FormatFloat(temp2[percentile_idx], 'f', -1, 64))
            }
            results.Write(avg)
            results.Write(percentile99)
        }
    }
}
