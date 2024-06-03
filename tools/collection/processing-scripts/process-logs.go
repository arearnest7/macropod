package main

import (
    "encoding/csv"
    "fmt"
    "os"
    "strings"
    "path/filepath"
    "io/ioutil"
    "time"
    "strconv"
)

func main() {
    if len(os.Args) <= 5 {
        fmt.Println("go run process-logs [entrypoint formatted as file_idx:checkpoint:depth] [directory] [results.csv] [prefix] [file1] ... [filen]")
    } else {
        entrypoint := strings.Split(os.Args[1], ":")
        functions := os.Args[5:]
        files, _ := ioutil.ReadDir(os.Args[2])
        var logs [][]string
        for _, function := range functions {
            var function_files []string
            for _, file := range files {
                match, _ := filepath.Match(os.Args[4] + "-" + function + "*.csv", file.Name())
                if match {
                    function_files = append(function_files, file.Name())
                }
            }
            logs = append(logs, function_files)
        }
        var records [][][]string
        var tracker []int
        for _, function_files := range logs {
            var record [][]string
            var file_dumps [][][]string
            l := 0
            var li []int
            for _, log := range function_files {
                file, _ := os.Open(os.Args[2] + "/" + log)
                reader := csv.NewReader(file)
                file_dump, _ := reader.ReadAll()
                if len(file_dump) > 0 {
                    l += len(file_dump)
                    li = append(li, 0)
                    file_dumps = append(file_dumps, file_dump)
                }
            }
            for _ = range l {
                idx := -1
                var prev_ts time.Time
                for i, file_dump := range file_dumps {
                    if li[i] < len(file_dump) {
                        ts, _ := time.Parse("2006-01-02 15:04:05.000000 UTC", file_dump[li[i]][0])
                        if prev_ts.Before(ts) {
                            idx = i
                            prev_ts = ts
                        }
                    }
                }
                record = append(record, file_dumps[idx][li[idx]])
                li[idx] += 1
            }
            records = append(records, record)
        }
        idx_0, _ := strconv.Atoi(entrypoint[0])
        for i, line := range records[idx_0] {
            if line[5] == entrypoint[1] && line[2] == entrypoint[2] {
                tracker = append(tracker, i)
            }
        }
        r, _ := os.Create(os.Args[3])
        defer r.Close()
        results := csv.NewWriter(r)
        defer results.Flush()
        var ckpt_int []string
        var ts_int []time.Time
        var set_final []bool
        wfid_0 := records[idx_0][tracker[0]][1]
        for i, record := range records {
            for _, line := range record {
                temp_ckpt := strconv.Itoa(i) + ":" + line[5] + ":" + line[2]
                set_ckpt := false
                if line[1] == wfid_0 {
                    for _, checkpoint := range ckpt_int {
                        if checkpoint == temp_ckpt {
                            set_ckpt = true
                            break
                        }
                    }
                    if !set_ckpt {
                        ckpt_int = append(ckpt_int, temp_ckpt)
                        ts, _ := time.Parse("2006-01-02 15:04:05.000000 UTC", line[0])
                        ts_int = append(ts_int, ts)
                        set_final = append(set_final, false)
                    }
                }
            }
        }
        var checkpoints []string
        for _ = range len(ckpt_int) {
            idx := -1
            var ts_prev time.Time
            for j := range len(ckpt_int) {
                if !set_final[j] && (idx == -1 || ts_prev.After(ts_int[j])) {
                    ts_prev = ts_int[j]
                    idx = j
                }
            }
            set_final[idx] = true
            checkpoints = append(checkpoints, ckpt_int[idx])
        }
        var tags []string
        for i, checkpoint := range checkpoints[1:] {
            tags = append(tags, checkpoints[i] + "-" + checkpoint)
        }
        results.Write(tags)
        var export_timestamps []string
        export_timestamps = append(export_timestamps, records[idx_0][0][0])
        export_timestamps = append(export_timestamps, records[idx_0][len(records[idx_0])-1][0])
        r2, _ := os.Create(os.Args[3] + ".mts")
        defer r2.Close()
        results2 := csv.NewWriter(r2)
        defer results2.Flush()
        results2.Write(export_timestamps)

        for _, idx := range tracker {
            var diff []string
            workflow_id := records[idx_0][idx][1]
            prev_timestamp, _ := time.Parse("2006-01-02 15:04:05.000000 UTC", records[idx_0][idx][0])
            for _, val := range checkpoints[1:] {
                checkpoint := strings.Split(val, ":")
                record_id, _ := strconv.Atoi(checkpoint[0])
                set_val := false
                for _, line := range records[record_id] {
                    if line[5] == checkpoint[1] && workflow_id == line[1] && line[2] == checkpoint[2] {
                        timestamp, _ := time.Parse("2006-01-02 15:04:05.000000 UTC", line[0])
                        d := timestamp.Sub(prev_timestamp)
                        prev_timestamp = timestamp
                        diff = append(diff, strconv.FormatInt(d.Nanoseconds(), 10))
                        set_val = true
                        break
                    }
                }
                if !set_val {
                    diff = append(diff, "0")
                }
            }
            results.Write(diff)
        }
    }
}
