package main

import (
    "encoding/csv"
    "fmt"
    "os"
    "strconv"
    "path/filepath"
    "strings"
)

func main() {
    if len(os.Args) <= 3 {
        fmt.Println("go run process-latency [N] [results.csv] [file1] ... [filen]")
    } else {
        n, _ := strconv.Atoi(os.Args[1])
        var file_names []string
        var file_names_stripped []string
        r, _ := os.Create(os.Args[2])
        defer r.Close()
        results := csv.NewWriter(r)
        defer results.Flush()
        for _, file_name := range os.Args[3:] {
            _, f := filepath.Split(file_name)
            idx := strings.Index(f, ".")
            file_names = append(file_names, file_name)
            file_names_stripped = append(file_names_stripped, f[:idx])
        }
        results.Write(file_names_stripped)
        var records [][][]string
        for _, file_name := range file_names {
            file, _ := os.Open(file_name)
            reader := csv.NewReader(file)
            record, _ := reader.ReadAll()
            records = append(records, record)
        }
        for i := range n {
            var l []string
            for _, record := range records {
                if len(record) > i + 1 {
                    l = append(l, record[i+1][0])
                } else {
                    l = append(l, "")
                }
            }
            results.Write(l)
        }
    }
}
