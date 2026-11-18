package main

import (
    pb "app/macropod_pb"

    "os"
    "fmt"
    "strconv"
    "encoding/json"
    "encoding/csv"
    "io/ioutil"
    "strings"
    "time"
    "os/exec"
    "github.com/montanaflynn/stats"

    "net"
    "net/http"

    "google.golang.org/grpc"
    "golang.org/x/net/context"
    "google.golang.org/grpc/connectivity"
)

type EvalService struct {
    pb.UnimplementedMacroPodEvalServer
}

type LatencyResult struct {
    Latency float64 // latency in seconds between start and end timestamps
    Start   string  // string representation of invocation start time
    End     string  // string representation of invocation end time
}

var (
    metrics_dir          string                                         // string name of directory for storing all metrics data from the evaluations
    latency_dir          string                                         // string name of directory for storing all latency data from the evaluations
    summary_dir          string                                         // string name of directory for storing all summary statistics data from the evaluations
    ingress_address      string                                         // string IP:port for the ingress gRPC service
    ingress_address_http string                                         // string IP:port for the ingress HTTP service
    nodes                []string                                       // string array of all IPs for the nodes
    debug                int                                            // debug level (higher means more verbose logs)
    ttl                  string                                         // time-to-live of the invocation requests

    ingress_channel      *grpc.ClientConn                               // gRPC channel for creating and deleting workflows
    ingress_stub         pb.MacroPodIngressClient                       // gRPC stub for creating and deleting workflows

    metrics_channel      = make(map[string]*grpc.ClientConn)            // gRPC channels for collecting all metrics data from the worker nodes
    metrics_stub         = make(map[string]pb.MacroPodMetricsClient)    // gRPC stubs for collecting all metrics data from the worker nodes
)

// purpose: print out message if debug level is greater than or equal to the global debug level
// input:
//  message: the message to be printed
//  debug_level: the value that must be exceeded or matched for the message to be printed
func Debug(message string, debug_level int) {
    if debug >= debug_level {
        fmt.Println(time.Now().UTC().Format("2006-01-02 15:04:05.000000 UTC") + ": " + message + "\n")
    }
}

// purpose: checks if the ingress controller is alive, and if not iterates attempting to recreate the channel until its able to connect
func Ingress_Check() {
    var err error
    for ingress_channel == nil {
        ingress_channel, err = grpc.Dial(ingress_address, grpc.WithInsecure())
        if err != nil {
            Debug(err.Error(), 0)
        }
        ingress_stub = pb.NewMacroPodIngressClient(ingress_channel)
        Debug("attempting rebuild ingress stub...", 5)
        time.Sleep(10 * time.Millisecond)
    }
    for ingress_channel.GetState() != connectivity.Ready {
        if ingress_channel.GetState() == connectivity.Connecting {
            Debug("waiting for ingress stub to be finish connecting", 5)
            time.Sleep(10 * time.Millisecond)
        } else if ingress_channel.GetState() == connectivity.TransientFailure || ingress_channel.GetState() == connectivity.Shutdown {
            ingress_channel, err = grpc.Dial(ingress_address, grpc.WithInsecure())
            if err != nil {
                Debug(err.Error(), 0)
            }
            ingress_stub = pb.NewMacroPodIngressClient(ingress_channel)
            Debug("attempting rebuild ingress stub due to transient failure or shutdown...", 5)
            time.Sleep(10 * time.Millisecond)
        } else {
            break
        }
    }
}

// purpose: checks if the metrics servers are alive, and if not iterates attempting to recreate the channels for each until its able to connect all
func Metrics_Check() {
    var err error
    for _, worker_node := range nodes {
        _, exists := metrics_channel[worker_node]
        if !exists {
            metrics_channel[worker_node], err = grpc.Dial(worker_node + ":10000", grpc.WithInsecure())
            if err != nil {
                Debug(err.Error(), 0)
            }
            metrics_stub[worker_node] = pb.NewMacroPodMetricsClient(metrics_channel[worker_node])
            Debug("attempting rebuild metrics stub...", 5)
        }
        for metrics_channel[worker_node].GetState() != connectivity.Ready {
            if ingress_channel.GetState() == connectivity.Connecting {
                Debug("waiting for metrics stub to be finish connecting", 5)
                time.Sleep(10 * time.Millisecond)
            } else if metrics_channel[worker_node].GetState() == connectivity.TransientFailure || metrics_channel[worker_node].GetState() == connectivity.Shutdown {
                metrics_channel[worker_node], err = grpc.Dial(worker_node + ":10000", grpc.WithInsecure())
                if err != nil {
                    Debug(err.Error(), 0)
                }
                metrics_stub[worker_node] = pb.NewMacroPodMetricsClient(metrics_channel[worker_node])
                Debug("attempting rebuild metrics stub due to transient failure or shutdown...", 5)
                time.Sleep(10 * time.Millisecond)
            } else {
                break
            }
        }
    }
}

// purpose: collects and tabulates peak values for metrics across all worker nodes until signaled to stop; results are appended to the metrics results file for the evaluation under ID with the benchmark, concurrency and phase marked
// input:
//  eval_id: unique string identifier for the evaluation
//  collect: signals when to end the collection
//  benchmark: signals what is the currently running benchmark variant
//  concurrency: signals what is the current workflow concurrency for the benchmark evaluation
//  phase: relays if its currently in cold or warm start
//  peak: multi layered map defining the peak values observed for all metrics for each benchmark, concurrency and phase
//  metrics: labels of all collected metrics
func Collect_Metrics(eval_id string, collect *bool, benchmark *string, concurrency *string, phase *string, peak map[string]map[string]map[string]map[string]float64, metrics []string) {
    r_m, _ := os.Create(metrics_dir + eval_id)
    metrics_out := csv.NewWriter(r_m)
    labels := []string{"timestamp", "benchmark", "concurrency", "phase"}
    for _, metric := range metrics {
        labels = append(labels, metric)
    }
    for i, _ := range nodes {
        for _, metric := range metrics {
            labels = append(labels, metric + "_" + strconv.Itoa(i+1))
        }
    }
    metrics_out.Write(labels)
    for *collect {
        line := make(map[string]float64)
        cpu_percent := []float64{}
        cpu_cores := []float64{}
        var cpu_total float64
        for _, metric := range metrics {
            line[metric] = float64(0)
        }
        for i, worker_node := range nodes {
            res, err := metrics_stub[worker_node].GetMetrics(context.Background(), &pb.MacroPodRequest{})
            if err != nil {
                Debug(err.Error(), 0)
            }
            cpu_percent = append(cpu_percent, res.GetCPUUsed())
            cpu_cores = append(cpu_cores, res.GetCPUCountPhysical())
            line["memory_used"] += res.GetMemoryUsed()
            line["network_bytes_sent"] += res.GetNetworkBytesSent()

            line["memory_used_" + strconv.Itoa(i+1)] = res.GetMemoryUsed()
            line["network_bytes_sent_" + strconv.Itoa(i+1)] = res.GetNetworkBytesSent()
        }
        for i, _ := range cpu_percent {
            cpu_total += cpu_percent[i] * cpu_cores[i]
        }
        line["cpu"] = cpu_total / line["cpu_count_physical"]
        timestamp := time.Now().UTC().Format("2006-01-02 15:04:05 UTC")
        line_s := []string{timestamp, *benchmark, *concurrency, *phase}
        benchmark_temp := *benchmark
        concurrency_temp := *concurrency
        phase_temp := *phase
        _, benchmark_exists := peak[benchmark_temp]
        if !benchmark_exists && benchmark_temp != "" {
            peak[benchmark_temp] = make(map[string]map[string]map[string]float64)
        }
        _, concurrency_exists := peak[benchmark_temp][concurrency_temp]
        if !concurrency_exists && concurrency_temp != ""{
            peak[benchmark_temp][concurrency_temp] = make(map[string]map[string]float64)
        }
        _, phase_exists := peak[benchmark_temp][concurrency_temp][phase_temp]
        if !phase_exists && phase_temp != "" {
            peak[benchmark_temp][concurrency_temp][phase_temp] = make(map[string]float64)
        }
        for _, metric := range metrics {
            val := strconv.FormatFloat(line[metric], 'f', -1, 64)
            line_s = append(line_s, val)
            if benchmark_temp != "" && concurrency_temp != "" && phase_temp != "" && peak[benchmark_temp][concurrency_temp][phase_temp][metric] < line[metric] {
                peak[benchmark_temp][concurrency_temp][phase_temp][metric] = line[metric]
            }
        }
        for i, _ := range nodes {
            for _, metric := range metrics {
                val := strconv.FormatFloat(line[metric + "_" + strconv.Itoa(i+1)], 'f', -1, 64)
                line_s = append(line_s, val)
                if benchmark_temp != "" && concurrency_temp != "" && phase_temp != "" && peak[benchmark_temp][concurrency_temp][phase_temp][metric + "_" + strconv.Itoa(i+1)] < line[metric + "_" + strconv.Itoa(i+1)] {
                    peak[benchmark_temp][concurrency_temp][phase_temp][metric + "_" + strconv.Itoa(i+1)] = line[metric + "_" + strconv.Itoa(i+1)]
                }
            }
        }
        metrics_out.Write(line_s)
        metrics_out.Flush()
        time.Sleep(1 * time.Second)
    }
}

// purpose: collects latency for a singular invocation using a text/string payload
// input:
//  target: string of the HTTP address for the target
//  t: text/string payload
// output:
//  LatencyResults: collection of latency in seconds, start and end times
func Collect_Latency_HTTP_Text(target string, t string) (LatencyResult) {
    s := time.Now()
    cmd := exec.Command("curl", "-m", ttl, "-X", "POST", "-d", t, target)
    Debug(cmd.String(), 3)
    out, err := cmd.Output()
    if err != nil {
        Debug(err.Error(), 0)
    }
    e := time.Now()
    Debug("http reply: " + string(out), 2)
    latency := e.Sub(s).Seconds()
    start_s := s.UTC().Format("2006-01-02 15:04:05 UTC")
    end_s := e.UTC().Format("2006-01-02 15:04:05 UTC")
    return LatencyResult{Latency: latency, Start: start_s, End: end_s}
}

// purpose: collects latency for a singular invocation using a JSON payload
// input:
//  target: string of the HTTP address for the target
//  j: JSON payload
// output:
//  LatencyResults: collection of latency in seconds, start and end times
func Collect_Latency_HTTP_JSON(target string, j map[string]interface{}) (LatencyResult) {
    s := time.Now()
    payload, _ := json.Marshal(j)
    cmd := exec.Command("curl", "-m", ttl, "-X", "POST", "-d", string(payload), "-H", "Content-Type: application/json", target)
    Debug(cmd.String(), 3)
    out, err := cmd.Output()
    if err != nil {
        Debug(err.Error(), 0)
    }
    e := time.Now()
    Debug("http reply: " + string(out), 2)
    latency := e.Sub(s).Seconds()
    start_s := s.UTC().Format("2006-01-02 15:04:05 UTC")
    end_s := e.UTC().Format("2006-01-02 15:04:05 UTC")
    return LatencyResult{Latency: latency, Start: start_s, End: end_s}
}

// purpose: collects latency for a singular invocation using a bytes payload
// input:
//  target: string of the HTTP address for the target
//  d: bytes payload
// output:
//  LatencyResults: collection of latency in seconds, start and end times
func Collect_Latency_HTTP_Data(target string, d []byte) (LatencyResult) {
    s := time.Now()
    cmd := exec.Command("curl", "-m", ttl, "-X", "POST", "-d", string(d), "-H", "Content-Type: application/octet-stream", target)
    Debug(cmd.String(), 3)
    out, err := cmd.Output()
    if err != nil {
        Debug(err.Error(), 0)
    }
    e := time.Now()
    Debug("http reply: " + string(out), 2)
    latency := e.Sub(s).Seconds()
    start_s := s.UTC().Format("2006-01-02 15:04:05 UTC")
    end_s := e.UTC().Format("2006-01-02 15:04:05 UTC")
    return LatencyResult{Latency: latency, Start: start_s, End: end_s}
}

// purpose: begins evaluation based on the evlaution definition provided, cycling over all benchmark, concurrency and phases variants
// input:
//  request: EvalStruct that contains all of the workflow definitions, extra targets and collection parameters
// output:
//  string: evaluation ID for collecting resulting data
func Serve_Eval(request *pb.EvalStruct) (string) {
    Debug("Eval_Start", 0)
    Ingress_Check()
    if _, err := os.Stat(metrics_dir); os.IsNotExist(err) { // creating directories if not present
        os.Mkdir(metrics_dir, 0777)
    }
    if _, err := os.Stat(latency_dir); os.IsNotExist(err) {
        os.Mkdir(latency_dir, 0777)
    }
    if _, err := os.Stat(summary_dir); os.IsNotExist(err) {
        os.Mkdir(summary_dir, 0777)
    }
    collect := true
    benchmark := ""
    concurrency := ""
    phase := ""
    metrics := []string{"cpu", "memory_used", "network_bytes_sent"}
    peak := make(map[string]map[string]map[string]map[string]float64, 0)
    latency := make(map[string]map[string]map[string][]float64, 0)
    eval_id := time.Now().UTC().Format("2006-01-02-15-04-05-UTC")
    Metrics_Check()
    go Collect_Metrics(eval_id, &collect, &benchmark, &concurrency, &phase, peak, metrics)
    r_l, _ := os.Create(latency_dir + eval_id)
    r_s, _ := os.Create(summary_dir + eval_id)
    latency_out := csv.NewWriter(r_l)
    latency_out.Write([]string{"benchmark", "concurrency", "phase", "latency", "start_time", "end_time"})
    latency_out.Flush()
    summary_out := csv.NewWriter(r_s)
    summary_labels := []string{"benchmark", "concurrency", "phase", "latency_median", "latency_p99"}
    for _, metric := range metrics {
        summary_labels = append(summary_labels, metric + "_" + "peak")
    }
    for i, _ := range nodes {
        for _, metric := range metrics {
            summary_labels = append(summary_labels, metric + "_" + strconv.Itoa(i+1) + "_peak")
        }
    }
    summary_out.Write(summary_labels) // setting all labels as first entry
    summary_out.Flush()
    for workflow_name, workflow := range request.GetWorkflows() { // iterating over all MacroPod benchmarks
        for _, c := range request.GetWorkflowConcurrency() { // iterating over all concurrencies
            ingress_stub.CreateWorkflow(context.Background(), workflow)
            benchmark = workflow_name
            concurrency = strconv.Itoa(int(c))
            phase = "cold_start" // starting with cold start phase
            _, benchmark_exists := latency[benchmark]
            if !benchmark_exists && benchmark != "" {
                latency[benchmark] = make(map[string]map[string][]float64)
            }
            _, concurrency_exists := latency[benchmark][concurrency]
            if !concurrency_exists && concurrency != ""{
                latency[benchmark][concurrency] = make(map[string][]float64)
            }
            inv_c := make(chan LatencyResult)
            cnt := 0
            total_sent := 0
            var concurrency_latency []float64
            for len(concurrency_latency) < int(request.GetInvocations()) { // loop until all invocations have been processed
                if cnt < int(c) && total_sent < int(request.GetInvocations()) { // if current count < max concurrency start a new thread for invocation
                    go func() { // starting all concurrent threads
                        var t string
                        var j map[string]interface{}
                        var d []byte
                        if workflow.GetPayload() != nil { // based on payload type in definition the appropriate caller is used
                            if workflow.GetPayload().GetType() == "Text" {
                                t = workflow.GetPayload().GetText()
                                inv_c <- Collect_Latency_HTTP_Text("http://" + ingress_address_http + "/workflow/invoke/" + workflow.GetName(), t)
                            } else if workflow.GetPayload().GetType() == "JSON" {
                                j = workflow.GetPayload().GetJSON().AsMap()
                                inv_c <- Collect_Latency_HTTP_JSON("http://" + ingress_address_http + "/workflow/invoke/" + workflow.GetName(), j)
                            } else if workflow.GetPayload().GetType() == "Data" {
                                d = workflow.GetPayload().GetData()
                                inv_c <- Collect_Latency_HTTP_Data("http://" + ingress_address_http + "/workflow/invoke/" + workflow.GetName(), d)
                            }
                        } else {
                            inv_c <- Collect_Latency_HTTP_Text("http://" + ingress_address_http + "/workflow/invoke/" + workflow.GetName(), t)
                        }
                    }()
                    cnt += 1
                    total_sent += 1
                } else { // gather returning requests; defaults here if current invocation count == max concurrency
                    l := <-inv_c
                    concurrency_latency = append(concurrency_latency, l.Latency)
                    latency_out.Write([]string{benchmark, concurrency, phase, strconv.FormatFloat(l.Latency, 'f', -1, 64), l.Start, l.End})
                    latency_out.Flush()
                    cnt -= 1
                    Debug(benchmark + "," + concurrency + "," + phase + ": " + strconv.Itoa(len(concurrency_latency)) + " processed", 2)
                }
            }
            latency[benchmark][concurrency][phase] = concurrency_latency
            phase = "warm_start" // starting warm start phase
            inv_c = make(chan LatencyResult)
            cnt = 0
            total_sent = 0
            concurrency_latency = make([]float64, 0)
            for len(concurrency_latency) < int(request.GetInvocations()) { // loop until all invocations have been processed
                if cnt < int(c) && total_sent < int(request.GetInvocations()) { // if current count < max concurrency start a new thread for invocation
                    go func() { // starting all concurrent threads
                        var t string
                        var j map[string]interface{}
                        var d []byte
                        if workflow.GetPayload() != nil { // based on payload type in definition the appropriate caller is used
                            if workflow.GetPayload().GetType() == "Text" {
                                t = workflow.GetPayload().GetText()
                                inv_c <- Collect_Latency_HTTP_Text("http://" + ingress_address_http + "/workflow/invoke/" + workflow.GetName(), t)
                            } else if workflow.GetPayload().GetType() == "JSON" {
                                j = workflow.GetPayload().GetJSON().AsMap()
                                inv_c <- Collect_Latency_HTTP_JSON("http://" + ingress_address_http + "/workflow/invoke/" + workflow.GetName(), j)
                            } else if workflow.GetPayload().GetType() == "Data" {
                                d = workflow.GetPayload().GetData()
                                inv_c <- Collect_Latency_HTTP_Data("http://" + ingress_address_http + "/workflow/invoke/" + workflow.GetName(), d)
                            }
                        } else {
                            inv_c <- Collect_Latency_HTTP_Text("http://" + ingress_address_http + "/workflow/invoke/" + workflow.GetName(), t)
                        }
                    }()
                    cnt += 1
                    total_sent += 1
                } else { // gather returning requests; defaults here if current invocation count == max concurrency
                    l := <-inv_c
                    concurrency_latency = append(concurrency_latency, l.Latency)
                    latency_out.Write([]string{benchmark, concurrency, phase, strconv.FormatFloat(l.Latency, 'f', -1, 64), l.Start, l.End})
                    latency_out.Flush()
                    cnt -= 1
                    Debug(benchmark + "," + concurrency + "," + phase + ": " + strconv.Itoa(len(concurrency_latency)) + " processed", 2)
                }
            }
            latency[benchmark][concurrency][phase] = concurrency_latency
            Debug(benchmark + "-" + concurrency + " has been completed... deleting from ingress now...", 1)
            benchmark = ""
            concurrency = ""
            phase = ""
            wf_delete := workflow.GetName()
            ingress_stub.DeleteWorkflow(context.Background(), &pb.MacroPodRequest{Workflow: &wf_delete})
            Debug(benchmark + "-" + concurrency + " has been deleted", 1) // waits until benchmark is fully deleted from K8s before starting next collection (no noise from previous)
            time.Sleep(60 * time.Second)
        }
    }
    for workflow_name, target := range request.GetExtraTargets() { // iterating over all extra target benchmarks (such as Knative)
        for _, c := range request.GetWorkflowConcurrency() { // iterating over all concurrencies
            benchmark = workflow_name
            concurrency = strconv.Itoa(int(c))
            phase = "cold_start" // starting with cold start phase
            _, benchmark_exists := latency[benchmark]
            if !benchmark_exists && benchmark != "" {
                latency[benchmark] = make(map[string]map[string][]float64)
            }
            _, concurrency_exists := latency[benchmark][concurrency]
            if !concurrency_exists && concurrency != ""{
                latency[benchmark][concurrency] = make(map[string][]float64)
            }
            inv_c := make(chan LatencyResult)
            cnt := 0
            total_sent := 0
            var concurrency_latency []float64
            for len(concurrency_latency) < int(request.GetInvocations()) { // loop until all invocations have been processed
                if cnt < int(c) && total_sent < int(request.GetInvocations()) { // if current count < max concurrency start a new thread for invocation
                    go func() { // starting all concurrent threads
                        var t string
                        var j map[string]interface{}
                        var d []byte
                        if request.GetExtraTargetsPayload()[workflow_name] != nil { // based on payload type in definition the appropriate caller is used
                            if request.GetExtraTargetsPayload()[workflow_name].GetType() == "Text" {
                                t = request.GetExtraTargetsPayload()[workflow_name].GetText()
                                inv_c <- Collect_Latency_HTTP_Text(target, t)
                            } else if request.GetExtraTargetsPayload()[workflow_name].GetType() == "JSON" {
                                j = request.GetExtraTargetsPayload()[workflow_name].GetJSON().AsMap()
                                inv_c <- Collect_Latency_HTTP_JSON(target, j)
                            } else if request.GetExtraTargetsPayload()[workflow_name].GetType() == "Data" {
                                d = request.GetExtraTargetsPayload()[workflow_name].GetData()
                                inv_c <- Collect_Latency_HTTP_Data(target, d)
                            }
                        } else {
                            inv_c <- Collect_Latency_HTTP_Text(target, t)
                        }
                    }()
                    cnt += 1
                    total_sent += 1
                } else { // gather returning requests; defaults here if current invocation count == max concurrency
                    l := <-inv_c
                    concurrency_latency = append(concurrency_latency, l.Latency)
                    latency_out.Write([]string{benchmark, concurrency, phase, strconv.FormatFloat(l.Latency, 'f', -1, 64), l.Start, l.End})
                    latency_out.Flush()
                    cnt -= 1
                    Debug(benchmark + "," + concurrency + "," + phase + ": " + strconv.Itoa(len(concurrency_latency)) + " processed", 2)
                }
            }
            latency[benchmark][concurrency][phase] = concurrency_latency
            phase = "warm_start" // starting warm start phase
            inv_c = make(chan LatencyResult)
            cnt = 0
            total_sent = 0
            concurrency_latency = make([]float64, 0)
            for len(concurrency_latency) < int(request.GetInvocations()) { // loop until all invocations have been processed
                if cnt < int(c) && total_sent < int(request.GetInvocations()) { // if current count < max concurrency start a new thread for invocation
                    go func() { // starting all concurrent threads
                        var t string
                        var j map[string]interface{}
                        var d []byte
                        if request.GetExtraTargetsPayload()[workflow_name] != nil { // based on payload type in definition the appropriate caller is used
                            if request.GetExtraTargetsPayload()[workflow_name].GetType() == "Text" {
                                t = request.GetExtraTargetsPayload()[workflow_name].GetText()
                                inv_c <- Collect_Latency_HTTP_Text(target, t)
                            } else if request.GetExtraTargetsPayload()[workflow_name].GetType() == "JSON" {
                                j = request.GetExtraTargetsPayload()[workflow_name].GetJSON().AsMap()
                                inv_c <- Collect_Latency_HTTP_JSON(target, j)
                            } else if request.GetExtraTargetsPayload()[workflow_name].GetType() == "Data" {
                                d = request.GetExtraTargetsPayload()[workflow_name].GetData()
                                inv_c <- Collect_Latency_HTTP_Data(target, d)
                            }
                        } else {
                            inv_c <- Collect_Latency_HTTP_Text(target, t)
                        }
                    }()
                    cnt += 1
                    total_sent += 1
                } else { // gather returning requests; defaults here if current invocation count == max concurrency
                    l := <-inv_c
                    concurrency_latency = append(concurrency_latency, l.Latency)
                    latency_out.Write([]string{benchmark, concurrency, phase, strconv.FormatFloat(l.Latency, 'f', -1, 64), l.Start, l.End})
                    latency_out.Flush()
                    cnt -= 1
                    Debug(benchmark + "," + concurrency + "," + phase + ": " + strconv.Itoa(len(concurrency_latency)) + " processed", 2)
                }
            }
            latency[benchmark][concurrency][phase] = concurrency_latency
            Debug(benchmark + "-" + concurrency + " has been completed... sleeping 60 seconds", 1)
            benchmark = ""
            concurrency = ""
            phase = ""
            time.Sleep(60 * time.Second)
            Debug(benchmark + "-" + concurrency + " sleep cycle done", 1) // waits until benchmark is fully deleted from K8s before starting next collection (no noise from previous)
        }
    }
    collect = false
    Debug("collection ended... waiting 30 seconds...", 0) // ending collection, now processing summary results
    time.Sleep(30 * time.Second)
    Debug("collection wait cycle ended. processing results to summary...", 0)
    phases := []string{"cold_start", "warm_start"}
    for workflow_name, _ := range request.GetWorkflows() { // iterating over all benchmarks, concurrencies and phases and tabulating peak metrics and p99 latency for MacroPod variants
        for _, c := range request.GetWorkflowConcurrency() {
            for _, p := range phases {
                line := []string{workflow_name, strconv.Itoa(int(c)), p}
                median_latency, _ := stats.Median(latency[workflow_name][strconv.Itoa(int(c))][p])
                p99_latency, _ := stats.Percentile(latency[workflow_name][strconv.Itoa(int(c))][p], 99)
                line = append(line, strconv.FormatFloat(median_latency, 'f', -1, 64))
                line = append(line, strconv.FormatFloat(p99_latency, 'f', -1, 64))
                for _, metric := range metrics {
                    val := strconv.FormatFloat(peak[workflow_name][strconv.Itoa(int(c))][p][metric], 'f', -1, 64)
                    line = append(line, val)
                }
                for i, _ := range nodes {
                    for _, metric := range metrics {
                        val := strconv.FormatFloat(peak[workflow_name][strconv.Itoa(int(c))][p][metric + "_" + strconv.Itoa(i+1)], 'f', -1, 64)
                        line = append(line, val)
                    }
                }
                summary_out.Write(line)
                summary_out.Flush()
            }
        }
    }
    for workflow_name, _ := range request.GetExtraTargets() { // iterating over all benchmarks, concurrencies and phases and tabulating peak metrics and p99 latency for extra target variants
        for _, c := range request.GetWorkflowConcurrency() {
            for _, p := range phases {
                line := []string{workflow_name, strconv.Itoa(int(c)), p}
                median_latency, _ := stats.Median(latency[workflow_name][strconv.Itoa(int(c))][p])
                p99_latency, _ := stats.Percentile(latency[workflow_name][strconv.Itoa(int(c))][p], 99)
                line = append(line, strconv.FormatFloat(median_latency, 'f', -1, 64))
                line = append(line, strconv.FormatFloat(p99_latency, 'f', -1, 64))
                for _, metric := range metrics {
                    val := strconv.FormatFloat(peak[workflow_name][strconv.Itoa(int(c))][p][metric], 'f', -1, 64)
                    line = append(line, val)
                }
                for i, _ := range nodes {
                    for _, metric := range metrics {
                        val := strconv.FormatFloat(peak[workflow_name][strconv.Itoa(int(c))][p][metric + "_" + strconv.Itoa(i+1)], 'f', -1, 64)
                        line = append(line, val)
                    }
                }
                summary_out.Write(line)
                summary_out.Flush()
            }
        }
    }
    Debug("collection processing ended", 0)
    return eval_id
}

// purpose: returns a string of all contents from the metrics results for specified ID
// input:
//  request: wrapper for the evluation ID used
// output:
//  string: CSV file
func Serve_EvalMetrics(request *pb.MacroPodRequest) (string) {
    metrics_out, _ := os.ReadFile(metrics_dir + request.GetTarget())
    return string(metrics_out)
}

// purpose: returns a string of all contents from the latency results for specified ID
// input:
//  request: wrapper for the evluation ID used
// output:
//  string: CSV file
func Serve_EvalLatency(request *pb.MacroPodRequest) (string) {
    latency_out, _ := os.ReadFile(latency_dir + request.GetTarget())
    return string(latency_out)
}

// purpose: returns a string of all contents from the summary results for specified ID
// input:
//  request: wrapper for the evluation ID used
// output:
//  string: CSV file
func Serve_EvalSummary(request *pb.MacroPodRequest) (string) {
    summary_out, _ := os.ReadFile(summary_dir + request.GetTarget())
    return string(summary_out)
}

// purpose: gRPC handler for Evaluate
// input:
//  ctx: context object utilized in every gRPC request handler
//  req: EvalStruct object defined by protobuf
// output:
//  MacroPodReply: reply protobuf struct
//  error: always will be nil
func (s *EvalService) Eval(ctx context.Context, req *pb.EvalStruct) (*pb.MacroPodReply, error) {
    id := Serve_Eval(req)
    results := pb.MacroPodReply{Reply: &id}
    return &results, nil
}

// purpose: gRPC handler for EvaluatorMetrics
// input:
//  ctx: context object utilized in every gRPC request handler
//  req: MacroPodRequest object defined by protobuf
// output:
//  MacroPodReply: reply protobuf struct
//  error: always will be nil
func (s *EvalService) EvalMetrics(ctx context.Context, req *pb.MacroPodRequest) (*pb.MacroPodReply, error) {
    metrics := Serve_EvalMetrics(req)
    results := pb.MacroPodReply{Reply: &metrics}
    return &results, nil
}

// purpose: gRPC handler for EvalLatency
// input:
//  ctx: context object utilized in every gRPC request handler
//  req: MacroPodRequest object defined by protobuf
// output:
//  MacroPodReply: reply protobuf struct
//  error: always will be nil
func (s *EvalService) EvalLatency(ctx context.Context, req *pb.MacroPodRequest) (*pb.MacroPodReply, error) {
    latency := Serve_EvalLatency(req)
    results := pb.MacroPodReply{Reply: &latency}
    return &results, nil
}

// purpose: gRPC handler for EvalSummary
// input:
//  ctx: context object utilized in every gRPC request handler
//  req: MacroPodRequest object defined by protobuf
// output:
//  MacroPodReply: reply protobuf struct
//  error: always will be nil
func (s *EvalService) EvalSummary(ctx context.Context, req *pb.MacroPodRequest) (*pb.MacroPodReply, error) {
    summary := Serve_EvalSummary(req)
    results := pb.MacroPodReply{Reply: &summary}
    return &results, nil
}

// purpose: prints out a help message for quick reference on how to operate service
// input:
//  res: writer to provide HTTP reply data
//  req: GET/POST object
func HTTP_Help(res http.ResponseWriter, req *http.Request) {
    help_print := "usage:\n"
    help_print += "\t/eval/start\n"
    help_print += "\t\tpurpose: start an evaluation for all configurations defined in the input evaluation definition\n"
    help_print += "\t\tinput: EvaluatorStruct JSON file (see GitHub for more details)\n"
    help_print += "\t\toutput: evaluation ID for retrieving metrics, latency and summary files\n"
    help_print += "\t/eval/metrics/{id}\n"
    help_print += "\t\tpurpose: retrieve metrics results of the evaluation with said ID\n"
    help_print += "\t\toutput: CSV file\n"
    help_print += "\t/eval/latency/{id}\n"
    help_print += "\t\tpurpose: retrieve latency results of the evaluation with said ID\n"
    help_print += "\t\toutput: CSV file\n"
    help_print += "\t/eval/summary/{id}\n"
    help_print += "\t\tpurpose: retrieve summary results of the evaluation with said ID\n"
    help_print += "\t\toutput: retrieve metrics results of the evaluation with said ID\n"
    fmt.Fprint(res, help_print)
}

// purpose: http handler for Eval
// input:
//  res: writer to provide HTTP reply data
//  req: GET/POST object
func HTTP_Eval(res http.ResponseWriter, req *http.Request) {
    body, _ := ioutil.ReadAll(req.Body)
    request := pb.EvalStruct{}
    json.Unmarshal(body, &request)
    id := Serve_Eval(&request)
    fmt.Fprint(res, id)
}

// purpose: http handler for EvalMetrics
// input:
//  res: writer to provide HTTP reply data
//  req: GET/POST object
func HTTP_EvalMetrics(res http.ResponseWriter, req *http.Request) {
    id := req.PathValue("id")
    request := pb.MacroPodRequest{Target: &id}
    results := Serve_EvalMetrics(&request)
    fmt.Fprint(res, results)
}

// purpose: http handler for EvalLatency
// input:
//  res: writer to provide HTTP reply data
//  req: GET/POST object
func HTTP_EvalLatency(res http.ResponseWriter, req *http.Request) {
    id := req.PathValue("id")
    request := pb.MacroPodRequest{Target: &id}
    results := Serve_EvalLatency(&request)
    fmt.Fprint(res, results)
}

// purpose: http handler for EvalSummary
// input:
//  res: writer to provide HTTP reply data
//  req: GET/POST object
func HTTP_EvalSummary(res http.ResponseWriter, req *http.Request) {
    id := req.PathValue("id")
    request := pb.MacroPodRequest{Target: &id}
    results := Serve_EvalSummary(&request)
    fmt.Fprint(res, results)
}

// input:
//  SERVICE_PORT: string defining gRPC service binding; default 8000
//  HTTP_PORT: string defining HTTP service binding; default 8000
//  INGRESS_ADDRESS: string IP:port for the ingress gRPC service; default 127.0.0.1:8001
//  INGRESS_ADDRESS_HTTP: string IP:port for the ingress HTTP service; default 127.0.0.1:9001
//  METRICS_DIR: directory for storing all metrics results; default /app/metrics/
//  LATENCY_DIR: directory for storing all latency results; default /app/latency/
//  SUMMARY_DIR: directory for storing all summary results; default /app/summary/
//  WORKER_NODES: array of all worker node IPs as a string (for pulling stats); default "192.168.0.2 192.168.0.3 192.168.0.4 192.168.0.5"
//  DEBUG: debug level (higher is more verbose); default 0 (off)
//  TTL: time to live for all evaluation requests; default 360 (360s)
func main() {
    service_port := os.Getenv("SERVICE_PORT")
    if service_port == "" {
        service_port = "8000"
    }
    http_port := os.Getenv("HTTP_PORT")
    if http_port == "" {
        http_port = "9000"
    }
    ingress_address = os.Getenv("INGRESS_ADDRESS")
    if ingress_address == "" {
        ingress_address = "127.0.0.1:8001"
    }
    ingress_address_http = os.Getenv("INGRESS_ADDRESS_HTTP")
    if ingress_address_http == "" {
        ingress_address_http = "127.0.0.1:9001"
    }
    metrics_dir = os.Getenv("METRICS_DIR")
    if metrics_dir == "" {
        metrics_dir = "/app/metrics/"
    }
    latency_dir = os.Getenv("LATENCY_DIR")
    if latency_dir == "" {
        latency_dir = "/app/latency/"
    }
    summary_dir = os.Getenv("SUMMARY_DIR")
    if summary_dir == "" {
        summary_dir = "/app/summary/"
    }
    worker_nodes_str := os.Getenv("WORKER_NODES")
    if worker_nodes_str == "" {
        worker_nodes_str = "192.168.0.2 192.168.0.3 192.168.0.4 192.168.0.5"
    }
    nodes = strings.Split(worker_nodes_str, " ")

    var err error
    debug, err = strconv.Atoi(os.Getenv("DEBUG"))
    if err != nil {
        debug = 0
    }
    ttl = os.Getenv("TTL")
    if ttl == "" {
        ttl = "360"
    }
    l, _ := net.Listen("tcp", ":" + service_port)
    s := grpc.NewServer(grpc.MaxSendMsgSize(1024*1024*200), grpc.MaxRecvMsgSize(1024*1024*200))
    pb.RegisterMacroPodEvalServer(s, &EvalService{})
    Ingress_Check()
    Metrics_Check()
    go s.Serve(l)
    h := http.NewServeMux()
    h.HandleFunc("/", HTTP_Help)
    h.HandleFunc("/eval/start", HTTP_Eval)
    h.HandleFunc("/eval/metrics/{id}", HTTP_EvalMetrics)
    h.HandleFunc("/eval/latency/{id}", HTTP_EvalLatency)
    h.HandleFunc("/eval/summary/{id}", HTTP_EvalSummary)
    http.ListenAndServe(":" + http_port, h)
}
