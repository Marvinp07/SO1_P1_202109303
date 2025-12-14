package kernelproc

import (
    "encoding/json"
    "os"
)

type ProcMem struct {
    TotalKB uint64 `json:"total_kb"`
    FreeKB  uint64 `json:"free_kb"`
    UsedKB  uint64 `json:"used_kb"`
}

type ProcEntry struct {
    PID        int    `json:"pid"`
    Name       string `json:"name"`
    VSZKB      uint64 `json:"vsz_kb"`
    RSSKB      uint64 `json:"rss_kb"`
    MemPercent int    `json:"mem_percent"`
}

type ProcSnapshot struct {
    Memory    ProcMem     `json:"memory"`
    Processes []ProcEntry `json:"processes"`
}

func ReadProcJSON(path string) (*ProcSnapshot, error) {
    b, err := os.ReadFile(path)
    if err != nil {
        return nil, err
    }
    var s ProcSnapshot
    if err := json.Unmarshal(b, &s); err != nil {
        return nil, err
    }
    return &s, nil
}
