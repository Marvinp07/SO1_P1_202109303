package config

import (
    "os"

    "gopkg.in/yaml.v3"
)

type DaemonCfg struct {
    IntervalSeconds int    `yaml:"interval_seconds"`
    LogPath         string `yaml:"log_path"`
    ErrorLogPath    string `yaml:"error_log_path"`
    DBPath          string `yaml:"db_path"`
}

type FiltersCfg struct {
    Images []string `yaml:"images"`
    Names  []string `yaml:"names"`
}

type HTTPCfg struct {
    Enabled bool   `yaml:"enabled"`
    Addr    string `yaml:"addr"`
}

type Config struct {
    Daemon  DaemonCfg  `yaml:"daemon"`
    Filters FiltersCfg `yaml:"filters"`
    HTTP    HTTPCfg    `yaml:"http"`
}

func Load(path string) (*Config, error) {
    b, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }
    var c Config
    if err := yaml.Unmarshal(b, &c); err != nil {
        return nil, err
    }
    // Defaults mínimos
    if c.Daemon.IntervalSeconds <= 0 {
        c.Daemon.IntervalSeconds = 5
    }
    if c.Daemon.LogPath == "" {
        c.Daemon.LogPath = "../logs/daemon.log"
    }
    if c.Daemon.ErrorLogPath == "" {
        c.Daemon.ErrorLogPath = "../logs/daemon_error.log"
    }
    if c.HTTP.Addr == "" {
        c.HTTP.Addr = "0.0.0.0:8080"
    }
    return &c, nil
}
