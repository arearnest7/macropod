package main

import (
    pb "app/macropod_pb"

    "os"
    "fmt"
    "strconv"
    "time"
    host "github.com/shirou/gopsutil/v4/host"
    load "github.com/shirou/gopsutil/v4/load"
    cpu "github.com/shirou/gopsutil/v4/cpu"
    mem "github.com/shirou/gopsutil/v4/mem"
    disk "github.com/shirou/gopsutil/v4/disk"
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

func Collect_Metrics() {
    prev_set := false
    metrics = pb.MetricsStruct{}
    var prev_net network.IOCountersStat
    for true {
        // host
        uptime_i, _ := host.Uptime()
        uptime := float64(uptime_i)
        metrics.Uptime = &uptime

        // load
        avg, _ := load.Avg()
        metrics.LoadAvg1 = &avg.Load1
        metrics.LoadAvg5 = &avg.Load5
        metrics.LoadAvg15 = &avg.Load15

        // cpu
        cpu_used, _ := cpu.Percent(time.Second, false)
        cpu_logical, _ := cpu.Counts(true)
        cpu_physical, _ := cpu.Counts(false)
        cpu_logical_f := float64(cpu_logical)
        cpu_physical_f := float64(cpu_physical)
        metrics.CPUUsed = &cpu_used[0]
        metrics.CPUCountLogical = &cpu_logical_f
        metrics.CPUCountPhysical = &cpu_physical_f

        // mem
        virt_mem, _ := mem.VirtualMemory()
        mem_total := float64(virt_mem.Total)
        mem_available := float64(virt_mem.Available)
        mem_used := float64(virt_mem.Used)
        mem_buffers := float64(virt_mem.Buffers)
        mem_cached := float64(virt_mem.Cached)
        mem_writeback := float64(virt_mem.WriteBack)
        mem_dirty := float64(virt_mem.Dirty)
        mem_writebacktmp := float64(virt_mem.WriteBackTmp)
        mem_shared := float64(virt_mem.Shared)
        mem_slab := float64(virt_mem.Slab)
        mem_sreclaimable := float64(virt_mem.Sreclaimable)
        mem_sunreclaim := float64(virt_mem.Sunreclaim)
        mem_pagetables := float64(virt_mem.PageTables)
        mem_swapcached := float64(virt_mem.SwapCached)
        mem_commitlimit := float64(virt_mem.CommitLimit)
        mem_committedas := float64(virt_mem.CommittedAS)
        mem_hightotal := float64(virt_mem.HighTotal)
        mem_highfree := float64(virt_mem.HighFree)
        mem_lowtotal := float64(virt_mem.LowTotal)
        mem_lowfree := float64(virt_mem.LowFree)
        mem_swaptotal := float64(virt_mem.SwapTotal)
        mem_swapfree := float64(virt_mem.SwapFree)
        mem_mapped := float64(virt_mem.Mapped)
        mem_vmalloctotal := float64(virt_mem.VmallocTotal)
        mem_vmallocused := float64(virt_mem.VmallocUsed)
        mem_vmallocchunk := float64(virt_mem.VmallocChunk)
        mem_hugepagestotal := float64(virt_mem.HugePagesTotal)
        mem_hugepagesfree := float64(virt_mem.HugePagesFree)
        mem_hugepagesrsvd := float64(virt_mem.HugePagesRsvd)
        mem_hugepagessurp := float64(virt_mem.HugePagesSurp)
        mem_hugepagesize := float64(virt_mem.HugePageSize)
        mem_anonhugepages := float64(virt_mem.AnonHugePages)
        metrics.MemoryUsed = &mem_used
        metrics.MemoryAvailable = &mem_available
        metrics.MemoryTotal = &mem_total
        metrics.MemoryBuffers = &mem_buffers
        metrics.MemoryCached = &mem_cached
        metrics.MemoryWriteBack = &mem_writeback
        metrics.MemoryDirty = &mem_dirty
        metrics.MemoryWriteBackTmp = &mem_writebacktmp
        metrics.MemoryShared = &mem_shared
        metrics.MemorySlab = &mem_slab
        metrics.MemorySreclaimable = &mem_sreclaimable
        metrics.MemorySunreclaim = &mem_sunreclaim
        metrics.MemoryPageTables = &mem_pagetables
        metrics.MemorySwapCached = &mem_swapcached
        metrics.MemoryCommitLimit = &mem_commitlimit
        metrics.MemoryCommittedAS = &mem_committedas
        metrics.MemoryHighTotal = &mem_hightotal
        metrics.MemoryHighFree = &mem_highfree
        metrics.MemoryLowTotal = &mem_lowtotal
        metrics.MemoryLowFree = &mem_lowfree
        metrics.MemorySwapTotal = &mem_swaptotal
        metrics.MemorySwapFree = &mem_swapfree
        metrics.MemoryMapped = &mem_mapped
        metrics.MemoryVmallocTotal = &mem_vmalloctotal
        metrics.MemoryVmallocUsed = &mem_vmallocused
        metrics.MemoryVmallocChunk = &mem_vmallocchunk
        metrics.MemoryHugePagesTotal = &mem_hugepagestotal
        metrics.MemoryHugePagesFree = &mem_hugepagesfree
        metrics.MemoryHugePagesRsvd = &mem_hugepagesrsvd
        metrics.MemoryHugePagesSurp = &mem_hugepagessurp
        metrics.MemoryHugePageSize = &mem_hugepagesize
        metrics.MemoryAnonHugePages = &mem_anonhugepages

        // disk
        disk_usage, _ := disk.Usage("/")
        disk_total := float64(disk_usage.Total)
        disk_free := float64(disk_usage.Free)
        disk_used := float64(disk_usage.Used)
        disk_inodes_total := float64(disk_usage.InodesTotal)
        disk_inodes_used := float64(disk_usage.InodesUsed)
        disk_inodes_free := float64(disk_usage.InodesFree)
        metrics.DiskTotal = &disk_total
        metrics.DiskFree = &disk_free
        metrics.DiskUsed = &disk_used
        metrics.DiskInodesTotal = &disk_inodes_total
        metrics.DiskInodesUsed = &disk_inodes_used
        metrics.DiskInodesFree = &disk_inodes_free

        // net
        if !prev_set {
            prev_n, _ := network.IOCounters(false)
            prev_net = prev_n[0]
            prev_set = true
        }
        n, _ := network.IOCounters(false)
        net := n[0]
        net_bytes_sent := float64(net.BytesSent - prev_net.BytesSent)
        net_bytes_recv := float64(net.BytesRecv - prev_net.BytesRecv)
        net_packets_sent := float64(net.PacketsSent - prev_net.PacketsSent)
        net_packets_recv := float64(net.PacketsRecv - prev_net.PacketsRecv)
        net_errin := float64(net.Errin - prev_net.Errin)
        net_errout := float64(net.Errout - prev_net.Errout)
        net_dropin := float64(net.Dropin - prev_net.Dropin)
        net_dropout := float64(net.Dropout - prev_net.Dropout)
        net_fifoin := float64(net.Fifoin - prev_net.Fifoin)
        net_fifoout := float64(net.Fifoout - prev_net.Fifoout)
        metrics.NetworkBytesSent = &net_bytes_sent
        metrics.NetworkBytesRecv = &net_bytes_recv
        metrics.NetworkPacketsSent = &net_packets_sent
        metrics.NetworkPacketsRecv = &net_packets_recv
        metrics.NetworkErrin = &net_errin
        metrics.NetworkErrout = &net_errout
        metrics.NetworkDropin = &net_dropin
        metrics.NetworkDropout = &net_dropout
        metrics.NetworkFifoin = &net_fifoin
        metrics.NetworkFifoout = &net_fifoout

        prev_net = net
        time.Sleep(1 * time.Second)
    }
}

func Serve_GetMetrics() (pb.MetricsStruct) {
    return metrics
}

func (s *MetricsService) GetMetrics(ctx context.Context, req *pb.MacroPodRequest) (*pb.MetricsStruct, error) {
    results := Serve_GetMetrics()
    return &results, nil
}

func HTTP_GetMetrics(res http.ResponseWriter, req *http.Request) {
    fmt.Print("request processing")
    results := Serve_GetMetrics()
    res_txt := ""
    res_txt += strconv.FormatFloat(results.GetUptime(), 'f', -1, 64) + ","
    res_txt += strconv.FormatFloat(results.GetLoadAvg1(), 'f', -1, 64) + ","
    res_txt += strconv.FormatFloat(results.GetLoadAvg5(), 'f', -1, 64) + ","
    res_txt += strconv.FormatFloat(results.GetLoadAvg15(), 'f', -1, 64) + ","
    res_txt += strconv.FormatFloat(results.GetCPUUsed(), 'f', -1, 64) + ","
    res_txt += strconv.FormatFloat(results.GetCPUCountLogical(), 'f', -1, 64) + ","
    res_txt += strconv.FormatFloat(results.GetCPUCountPhysical(), 'f', -1, 64) + ","
    res_txt += strconv.FormatFloat(results.GetMemoryUsed(), 'f', -1, 64) + ","
    res_txt += strconv.FormatFloat(results.GetMemoryAvailable(), 'f', -1, 64) + ","
    res_txt += strconv.FormatFloat(results.GetMemoryTotal(), 'f', -1, 64) + ","
    res_txt += strconv.FormatFloat(results.GetMemoryBuffers(), 'f', -1, 64) + ","
    res_txt += strconv.FormatFloat(results.GetMemoryCached(), 'f', -1, 64) + ","
    res_txt += strconv.FormatFloat(results.GetMemoryWriteBack(), 'f', -1, 64) + ","
    res_txt += strconv.FormatFloat(results.GetMemoryDirty(), 'f', -1, 64) + ","
    res_txt += strconv.FormatFloat(results.GetMemoryWriteBackTmp(), 'f', -1, 64) + ","
    res_txt += strconv.FormatFloat(results.GetMemoryShared(), 'f', -1, 64) + ","
    res_txt += strconv.FormatFloat(results.GetMemorySlab(), 'f', -1, 64) + ","
    res_txt += strconv.FormatFloat(results.GetMemorySreclaimable(), 'f', -1, 64) + ","
    res_txt += strconv.FormatFloat(results.GetMemorySunreclaim(), 'f', -1, 64) + ","
    res_txt += strconv.FormatFloat(results.GetMemoryPageTables(), 'f', -1, 64) + ","
    res_txt += strconv.FormatFloat(results.GetMemorySwapCached(), 'f', -1, 64) + ","
    res_txt += strconv.FormatFloat(results.GetMemoryCommitLimit(), 'f', -1, 64) + ","
    res_txt += strconv.FormatFloat(results.GetMemoryCommittedAS(), 'f', -1, 64) + ","
    res_txt += strconv.FormatFloat(results.GetMemoryHighTotal(), 'f', -1, 64) + ","
    res_txt += strconv.FormatFloat(results.GetMemoryHighFree(), 'f', -1, 64) + ","
    res_txt += strconv.FormatFloat(results.GetMemoryLowTotal(), 'f', -1, 64) + ","
    res_txt += strconv.FormatFloat(results.GetMemoryLowFree(), 'f', -1, 64) + ","
    res_txt += strconv.FormatFloat(results.GetMemorySwapTotal(), 'f', -1, 64) + ","
    res_txt += strconv.FormatFloat(results.GetMemorySwapFree(), 'f', -1, 64) + ","
    res_txt += strconv.FormatFloat(results.GetMemoryMapped(), 'f', -1, 64) + ","
    res_txt += strconv.FormatFloat(results.GetMemoryVmallocTotal(), 'f', -1, 64) + ","
    res_txt += strconv.FormatFloat(results.GetMemoryVmallocUsed(), 'f', -1, 64) + ","
    res_txt += strconv.FormatFloat(results.GetMemoryVmallocChunk(), 'f', -1, 64) + ","
    res_txt += strconv.FormatFloat(results.GetMemoryHugePagesTotal(), 'f', -1, 64) + ","
    res_txt += strconv.FormatFloat(results.GetMemoryHugePagesFree(), 'f', -1, 64) + ","
    res_txt += strconv.FormatFloat(results.GetMemoryHugePagesRsvd(), 'f', -1, 64) + ","
    res_txt += strconv.FormatFloat(results.GetMemoryHugePagesSurp(), 'f', -1, 64) + ","
    res_txt += strconv.FormatFloat(results.GetMemoryHugePageSize(), 'f', -1, 64) + ","
    res_txt += strconv.FormatFloat(results.GetMemoryAnonHugePages(), 'f', -1, 64) + ","
    res_txt += strconv.FormatFloat(results.GetDiskUsed(), 'f', -1, 64) + ","
    res_txt += strconv.FormatFloat(results.GetDiskFree(), 'f', -1, 64) + ","
    res_txt += strconv.FormatFloat(results.GetDiskTotal(), 'f', -1, 64) + ","
    res_txt += strconv.FormatFloat(results.GetDiskInodesUsed(), 'f', -1, 64) + ","
    res_txt += strconv.FormatFloat(results.GetDiskInodesFree(), 'f', -1, 64) + ","
    res_txt += strconv.FormatFloat(results.GetDiskInodesTotal(), 'f', -1, 64) + ","
    res_txt += strconv.FormatFloat(results.GetNetworkBytesSent(), 'f', -1, 64) + ","
    res_txt += strconv.FormatFloat(results.GetNetworkBytesRecv(), 'f', -1, 64) + ","
    res_txt += strconv.FormatFloat(results.GetNetworkPacketsSent(), 'f', -1, 64) + ","
    res_txt += strconv.FormatFloat(results.GetNetworkPacketsRecv(), 'f', -1, 64) + ","
    res_txt += strconv.FormatFloat(results.GetNetworkErrin(), 'f', -1, 64) + ","
    res_txt += strconv.FormatFloat(results.GetNetworkErrout(), 'f', -1, 64) + ","
    res_txt += strconv.FormatFloat(results.GetNetworkDropin(), 'f', -1, 64) + ","
    res_txt += strconv.FormatFloat(results.GetNetworkDropout(), 'f', -1, 64) + ","
    res_txt += strconv.FormatFloat(results.GetNetworkFifoin(), 'f', -1, 64) + ","
    res_txt += strconv.FormatFloat(results.GetNetworkFifoout(), 'f', -1, 64) + "\n"
    fmt.Printf("%s\n",res_txt)
    fmt.Fprint(res, res_txt)
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
    l, _ := net.Listen("tcp", ":" + service_port)
    s := grpc.NewServer(grpc.MaxSendMsgSize(1024*1024*200), grpc.MaxRecvMsgSize(1024*1024*200))
    pb.RegisterMacroPodMetricsServer(s, &MetricsService{})

    go Collect_Metrics()
    go s.Serve(l)

    h := http.NewServeMux()
    h.HandleFunc("/", HTTP_GetMetrics)
    http.ListenAndServe(":" + http_port, h)
}
