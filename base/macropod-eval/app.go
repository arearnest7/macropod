package main

import (
    pb "app/macropod_pb"
    structpb "google.golang.org/protobuf/types/known/structpb"

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
    Latency float64
    Start   string
    End     string
}

var (
    metrics_dir     string
    latency_dir     string
    summary_dir     string
    ingress_address string
    worker_nodes    []string

    ingress_channel *grpc.ClientConn
    ingress_stub    pb.MacroPodIngressClient

    metrics_channel = make(map[string]*grpc.ClientConn)
    metrics_stub    = make(map[string]pb.MacroPodMetricsClient)
)

func Ingress_Check() {
    for ingress_channel == nil || ingress_channel.GetState() != connectivity.Ready {
        ingress_channel, _ = grpc.Dial(ingress_address, grpc.WithInsecure())
        ingress_stub = pb.NewMacroPodIngressClient(ingress_channel)
        time.Sleep(10 * time.Millisecond)
    }
}

func Metrics_Check() {
    for _, worker_node := range worker_nodes {
        _, exists := metrics_channel[worker_node]
        if !exists {
            metrics_channel[worker_node], _ = grpc.Dial(worker_node + ":10000", grpc.WithInsecure())
            metrics_stub[worker_node] = pb.NewMacroPodMetricsClient(metrics_channel[worker_node])
        }
        for metrics_channel[worker_node].GetState() != connectivity.Ready {
            metrics_channel[worker_node], _ = grpc.Dial(worker_node + ":10000", grpc.WithInsecure())
            metrics_stub[worker_node] = pb.NewMacroPodMetricsClient(metrics_channel[worker_node])
            time.Sleep(10 * time.Millisecond)
        }
    }
}

func Collect_Metrics(eval_id string, collect *bool, benchmark *string, concurrency *string, phase *string, peak map[string]map[string]map[string]map[string]float64, metrics []string) {
    r_m, _ := os.Create(metrics_dir + eval_id)
    metrics_out := csv.NewWriter(r_m)
    labels := []string{"timestamp", "benchmark", "concurrency", "phase"}
    for _, metric := range metrics {
        labels = append(labels, metric)
    }
    for i, _ := range worker_nodes {
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
        for i, worker_node := range worker_nodes {
            res, err := metrics_stub[worker_node].GetMetrics(context.Background(), &pb.MacroPodRequest{})
	    if err != nil {
			fmt.Printf("%s\n",err)
	    }
	    //fmt.Printf("%v\n",res)
            cpu_percent = append(cpu_percent, res.GetCPUUsed())
            cpu_cores = append(cpu_cores, res.GetCPUCountPhysical())
            line["cpu_count_logical"] += res.GetCPUCountLogical()
            line["cpu_count_physical"] += res.GetCPUCountPhysical()
            line["uptime"] += res.GetUptime()
            line["loadavg1"] += res.GetLoadAvg1()
            line["loadavg5"] += res.GetLoadAvg5()
            line["loadavg15"] += res.GetLoadAvg15()
            line["memory_used"] += res.GetMemoryUsed()
            line["memory_available"] += res.GetMemoryAvailable()
            line["memory_total"] += res.GetMemoryTotal()
            line["memory_buffers"] += res.GetMemoryBuffers()
            line["memory_cached"] += res.GetMemoryCached()
            line["memory_writeback"] += res.GetMemoryWriteBack()
            line["memory_dirty"] += res.GetMemoryDirty()
            line["memory_writeback_tmp"] += res.GetMemoryWriteBackTmp()
            line["memory_shared"] += res.GetMemoryShared()
            line["memory_slab"] += res.GetMemorySlab()
            line["memory_sreclaimable"] += res.GetMemorySreclaimable()
            line["memory_sunreclaim"] += res.GetMemorySunreclaim()
            line["memory_page_tables"] += res.GetMemoryPageTables()
            line["memory_swap_cached"] += res.GetMemorySwapCached()
            line["memory_commit_limit"] += res.GetMemoryCommitLimit()
            line["memory_committed_as"] += res.GetMemoryCommittedAS()
            line["memory_high_total"] += res.GetMemoryHighTotal()
            line["memory_high_free"] += res.GetMemoryHighFree()
            line["memory_low_total"] += res.GetMemoryLowTotal()
            line["memory_low_free"] += res.GetMemoryLowFree()
            line["memory_swap_total"] += res.GetMemorySwapTotal()
            line["memory_swap_free"] += res.GetMemorySwapFree()
            line["memory_mapped"] += res.GetMemoryMapped()
            line["memory_vmalloc_total"] += res.GetMemoryVmallocTotal()
            line["memory_vmalloc_used"] += res.GetMemoryVmallocUsed()
            line["memory_vmalloc_chunk"] += res.GetMemoryVmallocChunk()
            line["memory_huge_pages_total"] += res.GetMemoryHugePagesTotal()
            line["memory_huge_pages_free"] += res.GetMemoryHugePagesFree()
            line["memory_huge_pages_rsvd"] += res.GetMemoryHugePagesRsvd()
            line["memory_huge_pages_surp"] += res.GetMemoryHugePagesSurp()
            line["memory_huge_page_size"] += res.GetMemoryHugePageSize()
            line["memory_anon_huge_pages"] += res.GetMemoryAnonHugePages()
            line["disk_used"] += res.GetDiskUsed()
            line["disk_free"] += res.GetDiskFree()
            line["disk_total"] += res.GetDiskTotal()
            line["disk_inodes_used"] += res.GetDiskInodesUsed()
            line["disk_inodes_free"] += res.GetDiskInodesFree()
            line["disk_inodes_total"] += res.GetDiskInodesTotal()
            line["network_bytes_sent"] += res.GetNetworkBytesSent()
            line["network_bytes_recv"] += res.GetNetworkBytesRecv()
            line["network_packets_sent"] += res.GetNetworkPacketsSent()
            line["network_packets_recv"] += res.GetNetworkPacketsRecv()
            line["network_err_in"] += res.GetNetworkErrin()
            line["network_err_out"] += res.GetNetworkErrout()
            line["network_drop_in"] += res.GetNetworkDropin()
            line["network_drop_out"] += res.GetNetworkDropout()
            line["network_fifo_in"] += res.GetNetworkFifoin()
            line["network_fifo_out"] += res.GetNetworkFifoout()

            line["cpu_" + strconv.Itoa(i+1)] = res.GetCPUUsed()
            line["cpu_count_logical_" + strconv.Itoa(i+1)] = res.GetCPUCountLogical()
            line["cpu_count_physical_" + strconv.Itoa(i+1)] = res.GetCPUCountPhysical()
            line["uptime_" + strconv.Itoa(i+1)] = res.GetUptime()
            line["loadavg1_" + strconv.Itoa(i+1)] = res.GetLoadAvg1()
            line["loadavg5_" + strconv.Itoa(i+1)] = res.GetLoadAvg5()
            line["loadavg15_" + strconv.Itoa(i+1)] = res.GetLoadAvg15()
            line["memory_used_" + strconv.Itoa(i+1)] = res.GetMemoryUsed()
            line["memory_available_" + strconv.Itoa(i+1)] = res.GetMemoryAvailable()
            line["memory_total_" + strconv.Itoa(i+1)] = res.GetMemoryTotal()
            line["memory_buffers_" + strconv.Itoa(i+1)] = res.GetMemoryBuffers()
            line["memory_cached_" + strconv.Itoa(i+1)] = res.GetMemoryCached()
            line["memory_writeback_" + strconv.Itoa(i+1)] = res.GetMemoryWriteBack()
            line["memory_dirty_" + strconv.Itoa(i+1)] = res.GetMemoryDirty()
            line["memory_writeback_tmp_" + strconv.Itoa(i+1)] = res.GetMemoryWriteBackTmp()
            line["memory_shared_" + strconv.Itoa(i+1)] = res.GetMemoryShared()
            line["memory_slab_" + strconv.Itoa(i+1)] = res.GetMemorySlab()
            line["memory_sreclaimable_" + strconv.Itoa(i+1)] = res.GetMemorySreclaimable()
            line["memory_sunreclaim_" + strconv.Itoa(i+1)] = res.GetMemorySunreclaim()
            line["memory_page_tables_" + strconv.Itoa(i+1)] = res.GetMemoryPageTables()
            line["memory_swap_cached_" + strconv.Itoa(i+1)] = res.GetMemorySwapCached()
            line["memory_commit_limit_" + strconv.Itoa(i+1)] = res.GetMemoryCommitLimit()
            line["memory_committed_as_" + strconv.Itoa(i+1)] = res.GetMemoryCommittedAS()
            line["memory_high_total_" + strconv.Itoa(i+1)] = res.GetMemoryHighTotal()
            line["memory_high_free_" + strconv.Itoa(i+1)] = res.GetMemoryHighFree()
            line["memory_low_total_" + strconv.Itoa(i+1)] = res.GetMemoryLowTotal()
            line["memory_low_free_" + strconv.Itoa(i+1)] = res.GetMemoryLowFree()
            line["memory_swap_total_" + strconv.Itoa(i+1)] = res.GetMemorySwapTotal()
            line["memory_swap_free_" + strconv.Itoa(i+1)] = res.GetMemorySwapFree()
            line["memory_mapped_" + strconv.Itoa(i+1)] = res.GetMemoryMapped()
            line["memory_vmalloc_total_" + strconv.Itoa(i+1)] = res.GetMemoryVmallocTotal()
            line["memory_vmalloc_used_" + strconv.Itoa(i+1)] = res.GetMemoryVmallocUsed()
            line["memory_vmalloc_chunk_" + strconv.Itoa(i+1)] = res.GetMemoryVmallocChunk()
            line["memory_huge_pages_total_" + strconv.Itoa(i+1)] = res.GetMemoryHugePagesTotal()
            line["memory_huge_pages_free_" + strconv.Itoa(i+1)] = res.GetMemoryHugePagesFree()
            line["memory_huge_pages_rsvd_" + strconv.Itoa(i+1)] = res.GetMemoryHugePagesRsvd()
            line["memory_huge_pages_surp_" + strconv.Itoa(i+1)] = res.GetMemoryHugePagesSurp()
            line["memory_huge_page_size_" + strconv.Itoa(i+1)] = res.GetMemoryHugePageSize()
            line["memory_anon_huge_pages_" + strconv.Itoa(i+1)] = res.GetMemoryAnonHugePages()
            line["disk_used_" + strconv.Itoa(i+1)] = res.GetDiskUsed()
            line["disk_free_" + strconv.Itoa(i+1)] = res.GetDiskFree()
            line["disk_total_" + strconv.Itoa(i+1)] = res.GetDiskTotal()
            line["disk_inodes_used_" + strconv.Itoa(i+1)] = res.GetDiskInodesUsed()
            line["disk_inodes_free_" + strconv.Itoa(i+1)] = res.GetDiskInodesFree()
            line["disk_inodes_total_" + strconv.Itoa(i+1)] = res.GetDiskInodesTotal()
            line["network_bytes_sent_" + strconv.Itoa(i+1)] = res.GetNetworkBytesSent()
            line["network_bytes_recv_" + strconv.Itoa(i+1)] = res.GetNetworkBytesRecv()
            line["network_packets_sent_" + strconv.Itoa(i+1)] = res.GetNetworkPacketsSent()
            line["network_packets_recv_" + strconv.Itoa(i+1)] = res.GetNetworkPacketsRecv()
            line["network_err_in_" + strconv.Itoa(i+1)] = res.GetNetworkErrin()
            line["network_err_out_" + strconv.Itoa(i+1)] = res.GetNetworkErrout()
            line["network_drop_in_" + strconv.Itoa(i+1)] = res.GetNetworkDropin()
            line["network_drop_out_" + strconv.Itoa(i+1)] = res.GetNetworkDropout()
            line["network_fifo_in_" + strconv.Itoa(i+1)] = res.GetNetworkFifoin()
            line["network_fifo_out_" + strconv.Itoa(i+1)] = res.GetNetworkFifoout()
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
                //fmt.Println(benchmark_temp + " " + concurrency_temp + " " + phase_temp + " " + metric)
                //fmt.Println("val: " + val)
                peak[benchmark_temp][concurrency_temp][phase_temp][metric] = line[metric]
            }
        }
        for i, _ := range worker_nodes {
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

func Collect_Latency(workflow string, t string, j map[string]interface{}, d []byte) (LatencyResult) {
    //fmt.Println("collecting grpc")
    ctx, cancel := context.WithTimeout(context.Background(), time.Second * 180)
    s := time.Now()
    jstruct, _ := structpb.NewStruct(j)
    ingress_stub.WorkflowInvoke(ctx, &pb.MacroPodRequest{Workflow: &workflow, Text: &t, JSON: jstruct, Data: d})
    e := time.Now()
    cancel()
    latency := e.Sub(s).Seconds()
    start_s := s.UTC().Format("2006-01-02 15:04:05 UTC")
    end_s := e.UTC().Format("2006-01-02 15:04:05 UTC")
    return LatencyResult{Latency: latency, Start: start_s, End: end_s}
}

func Collect_Latency_HTTP_Text(target string, t string) (LatencyResult) {
    //fmt.Println("collecting http text")
    s := time.Now()
    cmd := exec.Command("curl", "-m", "180", "-X", "POST", "-d", t, "-H", "Content-Type: plain/txt", target)
    cmd.Run()
    e := time.Now()
    latency := e.Sub(s).Seconds()
    start_s := s.UTC().Format("2006-01-02 15:04:05 UTC")
    end_s := e.UTC().Format("2006-01-02 15:04:05 UTC")
    return LatencyResult{Latency: latency, Start: start_s, End: end_s}
}

func Collect_Latency_HTTP_JSON(target string, j map[string]interface{}) (LatencyResult) {
    //fmt.Println("collecting http json")
    s := time.Now()
    payload, _ := json.Marshal(j)
    cmd := exec.Command("curl", "-m", "180", "-X", "POST", "-d", string(payload), "-H", "Content-Type: application/json", target)
    cmd.Run()
    e := time.Now()
    latency := e.Sub(s).Seconds()
    start_s := s.UTC().Format("2006-01-02 15:04:05 UTC")
    end_s := e.UTC().Format("2006-01-02 15:04:05 UTC")
    return LatencyResult{Latency: latency, Start: start_s, End: end_s}
}

func Collect_Latency_HTTP_Data(target string, d []byte) (LatencyResult) {
    //fmt.Println("collecting http data")
    s := time.Now()
    cmd := exec.Command("curl", "-m", "180", "-X", "POST", "-d", string(d), "-H", "Content-Type: application/octet-stream", target)
    cmd.Run()
    e := time.Now()
    latency := e.Sub(s).Seconds()
    start_s := s.UTC().Format("2006-01-02 15:04:05 UTC")
    end_s := e.UTC().Format("2006-01-02 15:04:05 UTC")
    return LatencyResult{Latency: latency, Start: start_s, End: end_s}
}

func Serve_Eval(request *pb.EvalStruct) (string) {
    Ingress_Check()
    if _, err := os.Stat(metrics_dir); os.IsNotExist(err) {
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
    metrics := []string{"cpu", "cpu_count_logical", "cpu_count_physical", "uptime", "loadavg1", "loadavg5", "loadavg15", "memory_used", "memory_available", "memory_total", "memory_buffers", "memory_cached", "memory_writeback", "memory_dirty", "memory_writeback_tmp", "memory_shared", "memory_slab", "memory_sreclaimable", "memory_sunreclaim", "memory_page_tables", "memory_swap_cached", "memory_commit_limit", "memory_committed_as", "memory_high_total", "memory_high_free", "memory_low_total", "memory_low_free", "memory_swap_total", "memory_swap_free", "memory_mapped", "memory_vmalloc_total", "memory_vmalloc_used", "memory_vmalloc_chunk", "memory_huge_pages_total", "memory_huge_pages_free", "memory_huge_pages_rsvd", "memory_huge_pages_surp", "memory_huge_page_size", "memory_anon_huge_pages", "disk_used", "disk_free", "disk_total", "disk_inodes_used", "disk_inodes_free", "disk_inodes_total", "network_bytes_sent", "network_bytes_recv", "network_packets_sent", "network_packets_recv", "network_err_in", "network_err_out", "network_drop_in", "network_drop_out", "network_fifo_in", "network_fifo_out"}
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
    summary_labels := []string{"benchmark", "concurrency", "phase", "latency_p99"}
    for _, metric := range metrics {
        summary_labels = append(summary_labels, metric + "_" + "peak")
    }
    for i, _ := range worker_nodes {
        for _, metric := range metrics {
            summary_labels = append(summary_labels, metric + "_" + strconv.Itoa(i+1) + "_peak")
        }
    }
    summary_out.Write(summary_labels)
    summary_out.Flush()
    //fmt.Println(summary_labels)
    for workflow_name, workflow := range request.GetWorkflows() {
        for _, c := range request.GetWorkflowConcurrency() {
            ingress_stub.CreateWorkflow(context.Background(), workflow)
            benchmark = workflow_name
            concurrency = strconv.Itoa(int(c))
            phase = "cold_start"
            //fmt.Println(benchmark + " " + concurrency + " " + phase)
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
            for len(concurrency_latency) < int(request.GetInvocations()) {
                if cnt < int(c) && total_sent < int(request.GetInvocations()) {
                    go func() {
                        var t string
                        var j map[string]interface{}
                        var d []byte
                        if workflow.GetPayload() != nil {
                            if workflow.GetPayload().GetType() == "Text" {
                                t = workflow.GetPayload().GetText()
                            } else if workflow.GetPayload().GetType() == "JSON" {
                                j = workflow.GetPayload().GetJSON().AsMap()
                            } else if workflow.GetPayload().GetType() == "Data" {
                                d = workflow.GetPayload().GetData()
                            }
                        }
                        inv_c <- Collect_Latency(workflow.GetName(), t, j, d)
                    }()
                    cnt += 1
                    total_sent += 1
                } else {
                    l := <-inv_c
                    concurrency_latency = append(concurrency_latency, l.Latency)
                    latency_out.Write([]string{benchmark, concurrency, phase, strconv.FormatFloat(l.Latency, 'f', -1, 64), l.Start, l.End})
                    latency_out.Flush()
                    cnt -= 1
                    fmt.Println("processed " + benchmark + "," + concurrency + "," + phase + ": " + strconv.Itoa(cnt) + " remain unprocessed")
                }
            }
            latency[benchmark][concurrency][phase] = concurrency_latency
            phase = "warm_start"
            //fmt.Println(benchmark + " " + concurrency + " " + phase)
            inv_c = make(chan LatencyResult)
            cnt = 0
            total_sent = 0
            concurrency_latency = make([]float64, 0)
            for len(concurrency_latency) < int(request.GetInvocations()) {
                if cnt < int(c) && total_sent < int(request.GetInvocations()) {
                    go func() {
                        var t string
                        var j map[string]interface{}
                        var d []byte
                        if workflow.GetPayload() != nil {
                            if workflow.GetPayload().GetType() == "Text" {
                                t = workflow.GetPayload().GetText()
                            } else if workflow.GetPayload().GetType() == "JSON" {
                                j = workflow.GetPayload().GetJSON().AsMap()
                            } else if workflow.GetPayload().GetType() == "Data" {
                                d = workflow.GetPayload().GetData()
                            }
                        }
                        inv_c <- Collect_Latency(workflow.GetName(), t, j, d)
                    }()
                    cnt += 1
                    total_sent += 1
                } else {
                    l := <-inv_c
                    concurrency_latency = append(concurrency_latency, l.Latency)
                    latency_out.Write([]string{benchmark, concurrency, phase, strconv.FormatFloat(l.Latency, 'f', -1, 64), l.Start, l.End})
                    latency_out.Flush()
                    cnt -= 1
                    fmt.Println("processed " + benchmark + "," + concurrency + "," + phase + ": " + strconv.Itoa(cnt) + " remain unprocessed")
                }
            }
            latency[benchmark][concurrency][phase] = concurrency_latency
            fmt.Println(benchmark + "-" + concurrency + " has been completed... deleting from ingress now...")
            benchmark = ""
            concurrency = ""
            phase = ""
            wf_delete := workflow.GetName()
            ingress_stub.DeleteWorkflow(context.Background(), &pb.MacroPodRequest{Workflow: &wf_delete})
            time.Sleep(300 * time.Second)
        }
    }
    for workflow_name, target := range request.GetExtraTargets() {
        for _, c := range request.GetWorkflowConcurrency() {
            benchmark = workflow_name
            concurrency = strconv.Itoa(int(c))
            phase = "cold_start"
            //fmt.Println(benchmark + " " + concurrency + " " + phase)
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
            for len(concurrency_latency) < int(request.GetInvocations()) {
                if cnt < int(c) && total_sent < int(request.GetInvocations()) {
                    go func() {
                        var t string
                        var j map[string]interface{}
                        var d []byte
                        if request.GetExtraTargetsPayload()[workflow_name] != nil {
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
                } else {
                    l := <-inv_c
                    concurrency_latency = append(concurrency_latency, l.Latency)
                    latency_out.Write([]string{benchmark, concurrency, phase, strconv.FormatFloat(l.Latency, 'f', -1, 64), l.Start, l.End})
                    latency_out.Flush()
                    cnt -= 1
                    fmt.Println("processed " + benchmark + "," + concurrency + "," + phase + ": " + strconv.Itoa(cnt) + " remain unprocessed")
                }
            }
            latency[benchmark][concurrency][phase] = concurrency_latency
            phase = "warm_start"
            //fmt.Println(benchmark + " " + concurrency + " " + phase)
            inv_c = make(chan LatencyResult)
            cnt = 0
            total_sent = 0
            concurrency_latency = make([]float64, 0)
            for len(concurrency_latency) < int(request.GetInvocations()) {
                if cnt < int(c) && total_sent < int(request.GetInvocations()) {
                    go func() {
                        var t string
                        var j map[string]interface{}
                        var d []byte
                        if request.GetExtraTargetsPayload()[workflow_name] != nil {
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
                } else {
                    l := <-inv_c
                    concurrency_latency = append(concurrency_latency, l.Latency)
                    latency_out.Write([]string{benchmark, concurrency, phase, strconv.FormatFloat(l.Latency, 'f', -1, 64), l.Start, l.End})
                    latency_out.Flush()
                    cnt -= 1
                    fmt.Println("processed " + benchmark + "," + concurrency + "," + phase + ": " + strconv.Itoa(cnt) + " remain unprocessed")
                }
            }
            latency[benchmark][concurrency][phase] = concurrency_latency
            fmt.Println(benchmark + "-" + concurrency + " has been completed... sleeping 300 seconds")
            benchmark = ""
            concurrency = ""
            phase = ""
            time.Sleep(300 * time.Second)
        }
    }
    //fmt.Println("done collection, waiting 30 seconds....")
    collect = false
    time.Sleep(30 * time.Second)
    //fmt.Println("30 seconds are up....")
    phases := []string{"cold_start", "warm_start"}
    for workflow_name, _ := range request.GetWorkflows() {
        for _, c := range request.GetWorkflowConcurrency() {
            for _, p := range phases {
                line := []string{workflow_name, strconv.Itoa(int(c)), p}
                p99_latency, _ := stats.Percentile(latency[workflow_name][strconv.Itoa(int(c))][p], 99)
                line = append(line, strconv.FormatFloat(p99_latency, 'f', -1, 64))
                for _, metric := range metrics {
                    val := strconv.FormatFloat(peak[workflow_name][strconv.Itoa(int(c))][p][metric], 'f', -1, 64)
                    line = append(line, val)
                }
                for i, _ := range worker_nodes {
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
    for workflow_name, _ := range request.GetExtraTargets() {
        for _, c := range request.GetWorkflowConcurrency() {
            for _, p := range phases {
                line := []string{workflow_name, strconv.Itoa(int(c)), p}
                p99_latency, _ := stats.Percentile(latency[workflow_name][strconv.Itoa(int(c))][p], 99)
                line = append(line, strconv.FormatFloat(p99_latency, 'f', -1, 64))
                for _, metric := range metrics {
                    val := strconv.FormatFloat(peak[workflow_name][strconv.Itoa(int(c))][p][metric], 'f', -1, 64)
                    line = append(line, val)
                }
                for i, _ := range worker_nodes {
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
    //fmt.Println(eval_id)
    return eval_id
}

func Serve_EvalMetrics(request *pb.MacroPodRequest) (string) {
    metrics_out, _ := os.ReadFile(metrics_dir + request.GetTarget())
    return string(metrics_out)
}

func Serve_EvalLatency(request *pb.MacroPodRequest) (string) {
    latency_out, _ := os.ReadFile(latency_dir + request.GetTarget())
    return string(latency_out)
}

func Serve_EvalSummary(request *pb.MacroPodRequest) (string) {
    summary_out, _ := os.ReadFile(summary_dir + request.GetTarget())
    return string(summary_out)
}

func (s *EvalService) Eval(ctx context.Context, req *pb.EvalStruct) (*pb.MacroPodReply, error) {
    id := Serve_Eval(req)
    results := pb.MacroPodReply{Reply: &id}
    return &results, nil
}

func (s *EvalService) EvalMetrics(ctx context.Context, req *pb.MacroPodRequest) (*pb.MacroPodReply, error) {
    metrics := Serve_EvalMetrics(req)
    results := pb.MacroPodReply{Reply: &metrics}
    return &results, nil
}

func (s *EvalService) EvalLatency(ctx context.Context, req *pb.MacroPodRequest) (*pb.MacroPodReply, error) {
    latency := Serve_EvalLatency(req)
    results := pb.MacroPodReply{Reply: &latency}
    return &results, nil
}

func (s *EvalService) EvalSummary(ctx context.Context, req *pb.MacroPodRequest) (*pb.MacroPodReply, error) {
    summary := Serve_EvalSummary(req)
    results := pb.MacroPodReply{Reply: &summary}
    return &results, nil
}

func HTTP_Help(res http.ResponseWriter, req *http.Request) {
    help_print := "Macropod Eval function\n"
    fmt.Fprint(res, help_print)
}

func HTTP_Eval(res http.ResponseWriter, req *http.Request) {
    body, _ := ioutil.ReadAll(req.Body)
    request := pb.EvalStruct{}
    json.Unmarshal(body, &request)
    fmt.Printf("%v\n", request)
    id := Serve_Eval(&request)
    fmt.Fprint(res, id)
}

func HTTP_EvalMetrics(res http.ResponseWriter, req *http.Request) {
    id := req.PathValue("id")
    request := pb.MacroPodRequest{Target: &id}
    results := Serve_EvalMetrics(&request)
    fmt.Fprint(res, results)
}

func HTTP_EvalLatency(res http.ResponseWriter, req *http.Request) {
    id := req.PathValue("id")
    request := pb.MacroPodRequest{Target: &id}
    results := Serve_EvalLatency(&request)
    fmt.Fprint(res, results)
}

func HTTP_EvalSummary(res http.ResponseWriter, req *http.Request) {
    id := req.PathValue("id")
    request := pb.MacroPodRequest{Target: &id}
    results := Serve_EvalSummary(&request)
    fmt.Fprint(res, results)
}

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
        worker_nodes_str = "192.168.56.21 192.168.56.22 192.168.56.23 192.168.56.24"
    }
    worker_nodes = strings.Split(worker_nodes_str, " ")
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
