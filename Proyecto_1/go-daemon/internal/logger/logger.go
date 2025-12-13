package logger

import (
    "log"
    "os"
    "time"
)

type Logger struct {
    l *log.Logger
}

func New(logPath, errPath string) (*Logger, *Logger, error) {
    // asegurar directorio
    if err := os.MkdirAll(dirOf(logPath), 0o755); err != nil {
        return nil, nil, err
    }
    if err := os.MkdirAll(dirOf(errPath), 0o755); err != nil {
        return nil, nil, err
    }

    lf, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
    if err != nil {
        return nil, nil, err
    }
    ef, err := os.OpenFile(errPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
    if err != nil {
        return nil, nil, err
    }

    return &Logger{l: log.New(lf, "", 0)}, &Logger{l: log.New(ef, "", 0)}, nil
}

func (lg *Logger) Info(format string, args ...any) {
    lg.l.Printf("%s INFO  "+format, append([]any{ts()}, args...)...)
}

func (lg *Logger) Error(format string, args ...any) {
    lg.l.Printf("%s ERROR "+format, append([]any{ts()}, args...)...)
}

func ts() string {
    return time.Now().Format("2006-01-02 15:04:05")
}

func dirOf(path string) string {
    i := len(path) - 1
    for i >= 0 && path[i] != '/' {
        i--
    }
    if i <= 0 {
        return "."
    }
    return path[:i]
}
