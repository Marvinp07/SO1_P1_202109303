package storage

import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"

    "proyecto/go-daemon/internal/metrics"
    "proyecto/go-daemon/internal/kernelproc"
)

var db *sql.DB

// Inicializar la base y crear tablas si no existen
func InitDB(path string) (*sql.DB, error) {
    var err error
    db, err = sql.Open("sqlite3", path)
    if err != nil {
        return nil, err
    }

    // Crear tablas si no existen
    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS sysinfo (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
            total_kb INTEGER,
            free_kb INTEGER,
            used_kb INTEGER,
            pid INTEGER,
            name TEXT,
            vsz_kb INTEGER,
            rss_kb INTEGER,
            mem_percent REAL
        );
        CREATE TABLE IF NOT EXISTS continfo (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
            total_kb INTEGER,
            free_kb INTEGER,
            used_kb INTEGER,
            pid INTEGER,
            name TEXT,
            vsz_kb INTEGER,
            rss_kb INTEGER,
            mem_percent REAL
        );
        CREATE TABLE IF NOT EXISTS docker_metrics (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
            container_id TEXT,
            name TEXT,
            image TEXT,
            cpu_percent REAL,
            mem_usage INTEGER,
            mem_limit INTEGER,
            mem_percent REAL,
            net_rx INTEGER,
            net_tx INTEGER,
            block_read INTEGER,
            block_write INTEGER,
            pids INTEGER
        );
    `)
    if err != nil {
        return nil, err
    }

    return db, nil
}

// Insertar métricas de Docker
func InsertDockerMetrics(m metrics.ContainerMetrics) error {
    _, err := db.Exec(`
        INSERT INTO docker_metrics (
            container_id, name, image, cpu_percent, mem_usage, mem_limit, mem_percent,
            net_rx, net_tx, block_read, block_write, pids
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
        m.ContainerID, m.Name, m.Image, m.CPUPercent, m.MemUsage, m.MemLimit, m.MemPercent,
        m.NetRx, m.NetTx, m.BlockRead, m.BlockWrite, m.PIDs,
    )
    return err
}

// Insertar snapshot de sysinfo
func InsertSysInfo(s *kernelproc.ProcSnapshot) error {
    for _, p := range s.Processes {
        _, err := db.Exec(`
            INSERT INTO sysinfo (total_kb, free_kb, used_kb, pid, name, vsz_kb, rss_kb, mem_percent)
            VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
            s.Memory.TotalKB, s.Memory.FreeKB, s.Memory.UsedKB,
            p.PID, p.Name, p.VSZKB, p.RSSKB, p.MemPercent,
        )
        if err != nil {
            return err
        }
    }
    return nil
}

// Insertar snapshot de continfo
func InsertContInfo(c *kernelproc.ProcSnapshot) error {
    for _, p := range c.Processes {
        _, err := db.Exec(`
            INSERT INTO continfo (total_kb, free_kb, used_kb, pid, name, vsz_kb, rss_kb, mem_percent)
            VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
            c.Memory.TotalKB, c.Memory.FreeKB, c.Memory.UsedKB,
            p.PID, p.Name, p.VSZKB, p.RSSKB, p.MemPercent,
        )
        if err != nil {
            return err
        }
    }
    return nil
}


