package main

import (
    "encoding/csv"
    "fmt"
    "os"
    "strconv"
    "path/filepath"
    "io/ioutil"
)

func main() {
    if len(os.Args) <= 3 {
        fmt.Println("go run clean-logs.go [directory] [results directory] [function1] ... [functionn]")
    } else {
        patterns := os.Args[3:]
        files, _ := ioutil.ReadDir(os.Args[1])
        var logs []string
        for _, pattern := range patterns {
            for _, file := range files {
                match, _ := filepath.Match("*" + pattern + "*.csv", file.Name())
                if match {
                    logs = append(logs, file.Name())
                }
            }
        }
        for _, log := range logs {
            i, _ := os.Open(os.Args[1] + "/" + log)
            reader := csv.NewReader(i)
            reader.FieldsPerRecord = -1
            reader.LazyQuotes = true
            record, _ := reader.ReadAll()
            o, _ := os.Create(os.Args[2] + "/" + log)
            writer := csv.NewWriter(o)
            defer writer.Flush()
            for _, line := range record {
                if len(line) == 6 {
                    if _, err := strconv.Atoi(line[len(line)-1]); err == nil {
                        writer.Write(line)
                    }
                }
            }
        }
    }
}
