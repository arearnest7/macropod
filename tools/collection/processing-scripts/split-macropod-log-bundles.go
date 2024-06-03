package main

import (
    "fmt"
    "os"
    "bufio"
    "strings"
)

func main() {
    if len(os.Args) <= 3 {
        fmt.Println("go run split-macropod-log-bundles.go [directory] [workflow1;workflow_log_bundle.csv;function1;...;functionn] ... [workflown;workflow_log_bundle.csv;function1;...;functionn]")
    } else {
        for _, workflow_arr_s := range os.Args[2:] {
            workflow_arr := strings.Split(workflow_arr_s, ";")
            i, _ := os.Open(os.Args[1] + "/" + workflow_arr[1])
            fscanner := bufio.NewScanner(i)
            var o *os.File
            for fscanner.Scan() {
                l := fscanner.Text()
                func_title_found := false
                for _, func_name := range workflow_arr[2:] {
                    if strings.Contains(l, func_name) {
                        func_title_found = true
                        break
                    }
                }
                if func_title_found {
                    o, _ = os.Create(os.Args[1] + "/" + workflow_arr[0] + "-" + l + ".csv")
                } else {
                    o.Write([]byte(l))
                }
            }
        }
    }
}
