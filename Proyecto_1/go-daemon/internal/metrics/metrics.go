package metrics

import (
    "fmt"
    "sync"
)

type DockerStatsJSON struct {
    CPUStats struct {
        CPUUsage struct {
            TotalUsage uint64 `json:"total_usage"`
        } `json:"cpu_usage"`
        SystemCPUUsage uint64 `json:"system_cpu_usage"`
        OnlineCPUs     uint64 `json:"online_cpus"`
    } `json:"cpu_stats"`

    PreCPUStats struct {
        CPUUsage struct {
            TotalUsage uint64 `json:"total_usage"`
        } `json:"cpu_usage"`
        SystemCPUUsage uint64 `json:"system_cpu_usage"`
        OnlineCPUs     uint64 `json:"online_cpus"`
    } `json:"precpu_stats"`

    MemoryStats struct {
        Usage uint64 `json:"usage"`
        Limit uint64 `json:"limit"`
    } `json:"memory_stats"`

    Networks map[string]struct {
        RxBytes uint64 `json:"rx_bytes"`
        TxBytes uint64 `json:"tx_bytes"`
    } `json:"networks"`

    BlkioStats struct {
        IoServiceBytesRecursive []struct {
            Major uint64 `json:"major"`
            Minor uint64 `json:"minor"`
            Op    string `json:"op"`
            Value uint64 `json:"value"`
        } `json:"io_service_bytes_recursive"`
    } `json:"blkio_stats"`

    PidsStats struct {
        Current uint64 `json:"current"`
    } `json:"pids_stats"`
}

type ContainerMetrics struct {
    ContainerID   string
    Name          string
    Image         string
    CPUPercent    float64
    MemUsage      uint64
    MemLimit      uint64
    MemPercent    float64
    NetRx         uint64
    NetTx         uint64
    BlockRead     uint64
    BlockWrite    uint64
    PIDs          int

    MemUsageHuman string
    MemLimitHuman string
    NetIOHuman    string
    BlockIOHuman  string
}

func FromDockerStats(s DockerStatsJSON) ContainerMetrics {
    cpuPerc := calcCPUPercent(s)
    memPerc := 0.0
    if s.MemoryStats.Limit > 0 {
        memPerc = float64(s.MemoryStats.Usage) / float64(s.MemoryStats.Limit) * 100.0
    }

    rx, tx := uint64(0), uint64(0)
    for _, n := range s.Networks {
        rx += n.RxBytes
        tx += n.TxBytes
    }

    var br, bw uint64
    for _, e := range s.BlkioStats.IoServiceBytesRecursive {
        switch e.Op {
        case "Read":
            br += e.Value
        case "Write":
            bw += e.Value
        }
    }

    m := ContainerMetrics{
        CPUPercent:    cpuPerc,
        MemUsage:      s.MemoryStats.Usage,
        MemLimit:      s.MemoryStats.Limit,
        MemPercent:    memPerc,
        NetRx:         rx,
        NetTx:         tx,
        BlockRead:     br,
        BlockWrite:    bw,
        PIDs:          int(s.PidsStats.Current),
        MemUsageHuman: humanBytes(s.MemoryStats.Usage),
        MemLimitHuman: humanBytes(s.MemoryStats.Limit),
        NetIOHuman:    fmt.Sprintf("%s / %s", humanBytes(rx), humanBytes(tx)),
        BlockIOHuman:  fmt.Sprintf("%s / %s", humanBytes(br), humanBytes(bw)),
    }
    return m
}

// Fórmula basada en stats de Docker
func calcCPUPercent(s DockerStatsJSON) float64 {
    cpuDelta := float64(s.CPUStats.CPUUsage.TotalUsage - s.PreCPUStats.CPUUsage.TotalUsage)
    sysDelta := float64(s.CPUStats.SystemCPUUsage - s.PreCPUStats.SystemCPUUsage)
    onlineCPUs := float64(s.CPUStats.OnlineCPUs)
    if onlineCPUs == 0 {
        onlineCPUs = 1
    }
    if sysDelta > 0 && cpuDelta > 0 {
        return (cpuDelta / sysDelta) * 100.0 * onlineCPUs
    }
    return 0.0
}

func humanBytes(b uint64) string {
    const kb = 1024
    const mb = kb * 1024
    const gb = mb * 1024
    switch {
    case b >= gb:
        return fmt.Sprintf("%.2fGiB", float64(b)/float64(gb))
    case b >= mb:
        return fmt.Sprintf("%.2fMiB", float64(b)/float64(mb))
    case b >= kb:
        return fmt.Sprintf("%.2fKiB", float64(b)/float64(kb))
    default:
        return fmt.Sprintf("%dB", b)
    }
}

var (
    snapshot     []ContainerMetrics
    snapshotLock sync.RWMutex
)

func UpdateSnapshot(ms []ContainerMetrics) {
    snapshotLock.Lock()
    defer snapshotLock.Unlock()
    snapshot = ms
}

func SnapshotStore() func() []ContainerMetrics {
    return func() []ContainerMetrics {
        snapshotLock.RLock()
        defer snapshotLock.RUnlock()
        out := make([]ContainerMetrics, len(snapshot))
        copy(out, snapshot)
        return out
    }
}
