package httpserver

import (
    "encoding/json"
    "net/http"

    "proyecto/go-daemon/internal/metrics"
)

func Run(addr string, getSnapshot func() []metrics.ContainerMetrics) error {
    return http.ListenAndServe(addr, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path != "/metrics" {
            http.NotFound(w, r)
            return
        }
        ms := getSnapshot()
        w.Header().Set("Content-Type", "application/json")
        if err := json.NewEncoder(w).Encode(ms); err != nil {
            http.Error(w, "Error serializando métricas", http.StatusInternalServerError)
        }
    }))
}
