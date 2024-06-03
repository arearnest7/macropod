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

func convert_timestring(timestring string) (time.Time) {
    if strings.Index(timestring, " m=+") > -1 {
        t, err := time.Parse("2006-01-02 15:04:05.000000000 +0000 UTC", strings.Split(timestring, " m=+")[0])
        if err != nil {
            t, err = time.Parse("2006-01-02 15:04:05.00000000 +0000 UTC", strings.Split(timestring, " m=+")[0])
            if err != nil {
                t, err = time.Parse("2006-01-02 15:04:05.0000000 +0000 UTC", strings.Split(timestring, " m=+")[0])
                if err != nil {
                    t, err = time.Parse("2006-01-02 15:04:05.000000 +0000 UTC", strings.Split(timestring, " m=+")[0])
                    if err != nil {
                        t, err = time.Parse("2006-01-02 15:04:05.00000 +0000 UTC", strings.Split(timestring, " m=+")[0])
                        if err != nil {
                            t, err = time.Parse("2006-01-02 15:04:05.0000 +0000 UTC", strings.Split(timestring, " m=+")[0])
                            if err != nil {
                                t, err = time.Parse("2006-01-02 15:04:05.000 +0000 UTC", strings.Split(timestring, " m=+")[0])
                                if err != nil {
                                    t, err = time.Parse("2006-01-02 15:04:05.00 +0000 UTC", strings.Split(timestring, " m=+")[0])
                                    if err != nil {
                                        t, err = time.Parse("2006-01-02 15:04:05.0 +0000 UTC", strings.Split(timestring, " m=+")[0])
                                        if err != nil {
                                            t, err = time.Parse("2006-01-02 15:04:05 +0000 UTC", strings.Split(timestring, " m=+")[0])
                                        }
                                    }
                                }
                            }
                        }
                    }
                }
            }
        }
        ts := t
        return ts
    } else if strings.Index(timestring, "pm") == -1 && strings.Index(timestring, "am") == -1 {
        t, err := time.Parse("2006-01-02 15:04:05.000000", timestring)
        if err != nil {
            t, err = time.Parse("2006-01-02 15:04:05.00000", timestring)
            if err != nil {
                t, err = time.Parse("2006-01-02 15:04:05.0000", timestring)
                if err != nil {
                    t, err = time.Parse("2006-01-02 15:04:05.000", timestring)
                    if err != nil {
                        t, err = time.Parse("2006-01-02 15:04:05.00", timestring)
                        if err != nil {
                            t, err = time.Parse("2006-01-02 15:04:05.0", timestring)
                            if err != nil {
                                t, err = time.Parse("2006-01-02 15:04:05", timestring)
                            }
                        }
                    }
                }
            }
        }
        ts := t
        return ts
    } else {
        days := map[string]string{"1st": "1", "2nd": "2", "3rd": "3", "4th": "4", "5th": "5",
                                  "6th": "6", "7th": "7", "8th": "8", "9th": "9", "10th": "10",
                                  "11th": "11", "12th": "12", "13th": "13", "14th": "14", "15th": "15",
                                  "16th": "16", "17th": "17", "18th": "18", "19th": "19", "20th": "20",
                                  "21st": "21", "22nd": "22", "23rd": "23", "24th": "24", "25th": "25",
                                  "26th": "26", "27th": "27", "28th": "28", "29th": "29", "30th": "30",
                                  "31st": "31"}
        var temp string
        for suffix, day := range days {
            if strings.Index(timestring, suffix) > -1 {
                temp = strings.Replace(timestring, suffix, day, 1)
                break
            }
        }
        t, _ := time.Parse("January 2 2006 3:04:05", temp[:len(temp)-5])
        ts := t
        return ts
    }
}

func main() {
    if len(os.Args) <= 5 {
        fmt.Println("go run process-logs [pattern] [directory] [results.csv] [prefix] [file1] ... [filen]")
    } else {
        headers := strings.Split(os.Args[1], ",")
        var pattern [][]string
        for _, header := range headers {
            pattern = append(pattern, strings.Split(header, ":"))
        }
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
                l += len(file_dump)
                li = append(li, 0)
                file_dumps = append(file_dumps, file_dump)
            }
            var prev_ts time.Time
            for _ = range l {
                for i, file_dump := range file_dumps {
                    if prev_ts.Before(convert_timestring(file_dump[li[i]][0])) {
                        record = append(record, file_dump[li[i]])
                        li[i] += 1
                        break
                    }
                }
            }
            records = append(records, record)
        }
        r, _ := os.Create(os.Args[3])
        defer r.Close()
        results := csv.NewWriter(r)
        defer results.Flush()
        var tags []string
        for i, header := range headers[1:] {
            tags = append(tags, headers[i] + "-" + header)
        }
        results.Write(tags)
        var export_timestamps []string
        idx_00, _ := strconv.Atoi(pattern[0][0])
        ts := convert_timestring(records[idx_00][0][0])
        export_timestamps = append(export_timestamps, ts.Format("2006-01-02 15:04:05"))
        ts = convert_timestring(records[idx_00][len(records[idx_00])-1][0])
        export_timestamps = append(export_timestamps, ts.Format("2006-01-02 15:04:05"))
        r2, _ := os.Create(os.Args[3] + ".mts")
        defer r2.Close()
        results2 := csv.NewWriter(r2)
        defer results2.Flush()
        results2.Write(export_timestamps)
        for i, line := range records[idx_00] {
            if line[len(line)-1] == pattern[0][1] {
                tracker = append(tracker, i)
            }
        }
        for _, idx := range tracker {
            var diff []string
            workflow_id := records[idx_00][idx][1]
            prev_timestamp := convert_timestring(records[idx_00][idx][0])
            for _, checkpoint := range pattern[1:] {
                record_id, _ := strconv.Atoi(checkpoint[0])
                count, _ := strconv.Atoi(checkpoint[2])
                var d time.Duration
                var first_timestamp time.Time
                var set bool
                for _, line := range records[record_id] {
                    if convert_timestring(line[0]).After(prev_timestamp) && line[len(line)-1] == checkpoint[1] && workflow_id == line[1] {
                        timestamp := convert_timestring(line[0])
                        if first_timestamp.IsZero() {
                            first_timestamp = timestamp
                        }
                        d = timestamp.Sub(prev_timestamp)
                        count -= 1
                        if count == 0 {
                            set = true
                            prev_timestamp = first_timestamp
                            diff = append(diff, d.String())
                            break
                        }
                    }
                }
                if !set {
                    diff = append(diff, "0")
                }
            }
            results.Write(diff)
        }
    }
}
