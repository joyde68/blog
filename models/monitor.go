package models

import (
	"fmt"
	"github.com/joyde68/blog/pkg"
	"runtime"
	"time"
)

type monitorStats struct {
	NumGoroutine int
	MemAllocated string
	MemMalloc    string
	MemTotal     string
	MemSys       string
	MemHeap      string
	MemGc        string
	LastGcTime   string
}

func ReadMemStats() *monitorStats {
	m := new(runtime.MemStats)
	runtime.ReadMemStats(m)
	ms := new(monitorStats)
	ms.NumGoroutine = runtime.NumGoroutine()
	ms.MemAllocated = pkg.FileSize(int64(m.Alloc))
	ms.MemTotal = pkg.FileSize(int64(m.TotalAlloc))
	ms.MemSys = pkg.FileSize(int64(m.Sys))
	ms.MemHeap = pkg.FileSize(int64(m.HeapAlloc))
	ms.MemMalloc = pkg.FileSize(int64(m.Mallocs))
	ms.LastGcTime = fmt.Sprintf("%.1fs", float64(time.Now().UnixNano()-int64(m.LastGC))/1000/1000/1000)
	ms.MemGc = pkg.FileSize(int64(m.NextGC))
	return ms
}
