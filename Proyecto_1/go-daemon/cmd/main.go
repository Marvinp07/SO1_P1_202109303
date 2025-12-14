package main

import (
    "context"
    "fmt"
    "os"
    "os/signal"
    "syscall"
    "time"

    "proyecto/go-daemon/internal/config"
    "proyecto/go-daemon/internal/docker"
    "proyecto/go-daemon/internal/httpserver"
    "proyecto/go-daemon/internal/logger"
    "proyecto/go-daemon/internal/metrics"
    "proyecto/go-daemon/internal/kernelproc"
    "proyecto/go-daemon/internal/storage"
)

func main() {
    // Cargar configuración
    cfg, err := config.Load("config.yaml")
    if err != nil {
        fmt.Println("Error cargando config:", err)
        os.Exit(1)
    }

    // Inicializar logger
    log, elog, err := logger.New(cfg.Daemon.LogPath, cfg.Daemon.ErrorLogPath)
    if err != nil {
        fmt.Println("Error inicializando logger:", err)
        os.Exit(1)
    }
    log.Info("Go-Daemon iniciado")

    // Inicializar cliente Docker
    cli, err := docker.NewClient()
    if err != nil {
        elog.Error("Error creando cliente Docker: %v", err)
        os.Exit(1)
    }

    // Inicializar SQLite
    db, err := storage.InitDB(cfg.Daemon.DBPath) // usa la ruta del config.yaml
    if err != nil {
        elog.Error("Error inicializando SQLite: %v", err)
        os.Exit(1)
    }
    defer db.Close()

    // Contexto y señales
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // HTTP server opcional
    if cfg.HTTP.Enabled {
        go func() {
            if err := httpserver.Run(cfg.HTTP.Addr, metrics.SnapshotStore()); err != nil {
                elog.Error("HTTP server error: %v", err)
            }
        }()
        log.Info("HTTP server escuchando en %s", cfg.HTTP.Addr)
    }

    // Loop principal
    ticker := time.NewTicker(time.Duration(cfg.Daemon.IntervalSeconds) * time.Second)
    defer ticker.Stop()

    // Manejo de señales
    sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

    for {
        select {
        case <-ticker.C:
            // Listar contenedores según filtros
            containers, err := docker.ListContainers(ctx, cli, cfg.Filters.Images, cfg.Filters.Names)
            if err != nil {
                elog.Error("Error listando contenedores: %v", err)
                continue
            }

            // Obtener métricas y registrar
            mset := make([]metrics.ContainerMetrics, 0, len(containers))
            for _, c := range containers {
                m, err := docker.GetContainerMetrics(ctx, cli, c.ID)
                if err != nil {
                    elog.Error("Error metrics cont=%s: %v", c.ID[:12], err)
                    continue
                }
                m.ContainerID = c.ID
                m.Name = c.Names[0]
                m.Image = c.Image

                log.Info("[%s] Image=%s CPU=%.2f%% RAM=%s/%s RAM%%=%.2f NET=%s BLOCK=%s PIDS=%d",
                    m.Name, m.Image, m.CPUPercent, m.MemUsageHuman, m.MemLimitHuman, m.MemPercent,
                    m.NetIOHuman, m.BlockIOHuman, m.PIDs)

                mset = append(mset, m)
            }

            // Actualizar snapshot para HTTP
            metrics.UpdateSnapshot(mset)

            // Leer sysinfo
            sysSnap, err := kernelproc.ReadProcJSON("/proc/sysinfo_so1_202109303")
            if err == nil {
                log.Info("Kernel SYS: Total=%dKB Free=%dKB Used=%dKB Proc=%d",
                    sysSnap.Memory.TotalKB, sysSnap.Memory.FreeKB, sysSnap.Memory.UsedKB, len(sysSnap.Processes))

                if err := storage.InsertSysInfo(sysSnap); err != nil {
                    elog.Error("Error insertando sysinfo: %v", err)
                }
            } else {
                elog.Error("Error leyendo sysinfo: %v", err)
            }


            // Leer continfo
            contSnap, err := kernelproc.ReadProcJSON("/proc/continfo_so1_202109303")
            if err == nil {
                log.Info("Kernel CONT: Proc=%d", len(contSnap.Processes))

                if err := storage.InsertContInfo(contSnap); err != nil {
                    elog.Error("Error insertando continfo: %v", err)
                }
            } else {
                elog.Error("Error leyendo continfo: %v", err)
            }


            // Guardar en SQLite
            for _, m := range mset {
                if err := storage.InsertDockerMetrics(m); err != nil {
                    elog.Error("Error insertando docker metrics: %v", err)
                }
            }

            if err := storage.InsertSysInfo(sysSnap); err != nil {
                elog.Error("Error insertando sysinfo: %v", err)
            }

            if err := storage.InsertContInfo(contSnap); err != nil {
                elog.Error("Error insertando continfo: %v", err)
            }

        case s := <-sigs:
            log.Info("Recibida señal: %s. Deteniendo daemon...", s.String())
            return
        }
    }
}
