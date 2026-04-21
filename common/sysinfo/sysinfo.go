package sysinfo

import (
	"runtime"
	"time"

	"ad2nx/api/panel"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
)

var startTime = time.Now()

// GetUptime 返回进程运行时长（秒）
func GetUptime() int64 {
	return int64(time.Since(startTime).Seconds())
}

// GetGoroutines 返回当前 goroutine 数量
func GetGoroutines() int {
	return runtime.NumGoroutine()
}

// GetNodeStatus 采集 CPU/Mem/Swap/Disk 系统信息
func GetNodeStatus() *panel.NodeStatus {
	status := &panel.NodeStatus{}

	// CPU 使用率（1 秒采样）
	if percents, err := cpu.Percent(time.Second, false); err == nil && len(percents) > 0 {
		status.CPU = percents[0]
	}

	// 内存
	if v, err := mem.VirtualMemory(); err == nil {
		status.Mem = panel.MemInfo{
			Total: int64(v.Total),
			Used:  int64(v.Used),
		}
	}

	// Swap
	if s, err := mem.SwapMemory(); err == nil {
		status.Swap = panel.MemInfo{
			Total: int64(s.Total),
			Used:  int64(s.Used),
		}
	}

	// 磁盘（根分区）
	if d, err := disk.Usage("/"); err == nil {
		status.Disk = panel.MemInfo{
			Total: int64(d.Total),
			Used:  int64(d.Used),
		}
	}

	return status
}

// GetNodeMetrics 采集运行时指标
// totalUsers: 当前用户总数
// activeUsers: 本周期有流量的活跃用户数
// inboundSpeed / outboundSpeed: 本周期的入站/出站速率 (bytes/s)
func GetNodeMetrics(totalUsers, activeUsers int, inboundSpeed, outboundSpeed int64) *panel.NodeMetrics {
	metrics := &panel.NodeMetrics{
		Uptime:        GetUptime(),
		Goroutines:    GetGoroutines(),
		TotalUsers:    totalUsers,
		ActiveUsers:   activeUsers,
		InboundSpeed:  inboundSpeed,
		OutboundSpeed: outboundSpeed,
	}

	// 每核 CPU 使用率
	if perCore, err := cpu.Percent(0, true); err == nil {
		metrics.CpuPerCore = perCore
	}

	// 系统负载（仅 Linux/macOS 有效，Windows 返回空）
	if avg, err := load.Avg(); err == nil {
		metrics.Load = []float64{avg.Load1, avg.Load5, avg.Load15}
	}

	return metrics
}
