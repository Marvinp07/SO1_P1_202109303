package docker

import (
    "context"
    "encoding/json"
    "fmt"
    "io"

    "github.com/docker/docker/api/types"
    "github.com/docker/docker/client"
    "proyecto/go-daemon/internal/metrics"
)

// NewClient crea un cliente Docker usando las variables de entorno
func NewClient() (*client.Client, error) {
    return client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
}

// ListContainers devuelve una lista de contenedores filtrados por imagen/nombre
func ListContainers(ctx context.Context, cli *client.Client, images []string, names []string) ([]types.Container, error) {
    opts := types.ContainerListOptions{All: false}
    cts, err := cli.ContainerList(ctx, opts)
    if err != nil {
        return nil, err
    }

    filtered := make([]types.Container, 0, len(cts))
    for _, c := range cts {
        if len(images) > 0 && !contains(images, c.Image) {
            continue
        }
        if len(names) > 0 && !anyNameMatch(names, c.Names) {
            continue
        }
        filtered = append(filtered, c)
    }
    return filtered, nil
}

// GetContainerMetrics obtiene métricas de un contenedor específico
func GetContainerMetrics(ctx context.Context, cli *client.Client, containerID string) (metrics.ContainerMetrics, error) {
    resp, err := cli.ContainerStats(ctx, containerID, false)
    if err != nil {
        return metrics.ContainerMetrics{}, err
    }
    defer resp.Body.Close()

    b, err := io.ReadAll(resp.Body)
    if err != nil {
        return metrics.ContainerMetrics{}, err
    }

    var s metrics.DockerStatsJSON
    if err := json.Unmarshal(b, &s); err != nil {
        return metrics.ContainerMetrics{}, err
    }

    m := metrics.FromDockerStats(s)
    return m, nil
}

// contains verifica si un valor está en un slice
func contains(xs []string, v string) bool {
    for _, x := range xs {
        if x == v {
            return true
        }
    }
    return false
}

// anyNameMatch verifica si algún nombre coincide con los filtros
func anyNameMatch(filters []string, names []string) bool {
    for _, f := range filters {
        for _, n := range names {
            // nombres vienen con prefijo "/" desde Docker
            if n == f || n == "/"+f {
                return true
            }
        }
    }
    return false
}

// DebugPrintStats imprime métricas en consola
func DebugPrintStats(m metrics.ContainerMetrics) {
    fmt.Printf("CPU=%.2f%% RAM=%s/%s RAM%%=%.2f NET=%s BLOCK=%s PIDS=%d\n",
        m.CPUPercent, m.MemUsageHuman, m.MemLimitHuman, m.MemPercent, m.NetIOHuman, m.BlockIOHuman, m.PIDs)
}
