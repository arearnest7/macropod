package main

import (
    pb "app/macropod_pb"

    "os"
    "fmt"
    "strconv"
    "time"
    cpu "github.com/shirou/gopsutil/v4/cpu"
    mem "github.com/shirou/gopsutil/v4/mem"
    network "github.com/shirou/gopsutil/v4/net"

    "net"
    "net/http"

    "google.golang.org/grpc"
    "golang.org/x/net/context"
)

type MetricsService struct {
    pb.UnimplementedMacroPodMetricsServer
}

var (
    metrics pb.MetricsStruct
)

// purpose: Indefinitely collects all system metrics every second and stores the results in metrics structure
func Collect_Metrics() {
    prev_set := false
    metrics = pb.MetricsStruct{}
    var prev_net network.IOCountersStat
    for true {
        // cpu
        cpu_used, _ := cpu.Percent(time.Second, false)
        cpu_physical, _ := cpu.Counts(false)
        cpu_physical_f := float64(cpu_physical)
        metrics.CPUUsed = &cpu_used[0]
        metrics.CPUCountPhysical = &cpu_physical_f

        // mem
        virt_mem, _ := mem.VirtualMemory()
        mem_used := float64(virt_mem.Used)
        metrics.MemoryUsed = &mem_used

        // net
        if !prev_set {
            prev_n, _ := network.IOCounters(false)
            prev_net = prev_n[0]
            prev_set = true
        }
        n, _ := network.IOCounters(false)
        net := n[0]
        net_bytes_sent := float64(net.BytesSent - prev_net.BytesSent)
        metrics.NetworkBytesSent = &net_bytes_sent

        prev_net = net
        time.Sleep(1 * time.Second)
    }
}

// purpose: returns metrics structure for handlers to return
// output:
//  metrics: metrics structure defined by protobuf
func Serve_GetMetrics() (pb.MetricsStruct) {
    return metrics
}

// purpose: gRPC handler for GetMetrics
// input:
//  ctx: context object utilized in every gRPC request handler
//  req: empty but required MacroPodRequest object
// output:
//  MetricsStruct: metrics protobuf struct that contains all metrics stats for 1 second
//  error: always will be nil
func (s *MetricsService) GetMetrics(ctx context.Context, req *pb.MacroPodRequest) (*pb.MetricsStruct, error) {
    results := Serve_GetMetrics()
    return &results, nil
}

// purpose: http handler for GetMetrics
// input:
//  res: writer to provide HTTP reply data
//  req: GET/POST object; for purposes nothing is used in it for this function
func HTTP_GetMetrics(res http.ResponseWriter, req *http.Request) {
    fmt.Print("request processing")
    results := Serve_GetMetrics()
    res_txt := ""
    res_txt += strconv.FormatFloat(results.GetCPUUsed(), 'f', -1, 64) + ","
    res_txt += strconv.FormatFloat(results.GetCPUCountPhysical(), 'f', -1, 64) + ","
    res_txt += strconv.FormatFloat(results.GetMemoryUsed(), 'f', -1, 64) + ","
    res_txt += strconv.FormatFloat(results.GetNetworkBytesSent(), 'f', -1, 64)
    fmt.Printf("%s\n",res_txt)
    fmt.Fprint(res, res_txt)
}

// input:
//  SERVICE_PORT: string defining gRPC service binding; default 8000
//  HTTP_PORT: string defining HTTP service binding; default 9000
func main() {
    service_port := os.Getenv("SERVICE_PORT")
    if service_port == "" {
        service_port = "8000"
    }
    http_port := os.Getenv("HTTP_PORT")
    if http_port == "" {
        http_port = "9000"
    }

    l, _ := net.Listen("tcp", ":" + service_port)
    s := grpc.NewServer(grpc.MaxSendMsgSize(1024*1024*200), grpc.MaxRecvMsgSize(1024*1024*200))
    pb.RegisterMacroPodMetricsServer(s, &MetricsService{})

    go Collect_Metrics()
    go s.Serve(l)

    h := http.NewServeMux()
    h.HandleFunc("/", HTTP_GetMetrics)
    http.ListenAndServe(":" + http_port, h)
}
