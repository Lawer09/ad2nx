package panel

// NodeReportRequest 是 /api/v2/server/report 的请求体
type NodeReportRequest struct {
	// traffic: map[userID][upload, download]
	Traffic map[int][]int64 `json:"traffic,omitempty"`
	// alive: map[userID][]ipIdentifier
	Alive map[int][]string `json:"alive,omitempty"`
	// online: map[userID]connectionCount
	Online map[int]int `json:"online,omitempty"`
	// status: node load
	Status *NodeStatus `json:"status,omitempty"`
	// metrics: runtime metrics
	Metrics *NodeMetrics `json:"metrics,omitempty"`
}

type NodeStatus struct {
	CPU          float64     `json:"cpu"`
	Mem          MemInfo     `json:"mem"`
	Swap         MemInfo     `json:"swap"`
	Disk         MemInfo     `json:"disk"`
	KernelStatus interface{} `json:"kernel_status,omitempty"`
	// 当前上报间隔内的带宽
	InboundSpeed  int64 `json:"inbound_speed"`  // 入站速率 (bytes/s)，即用户上传聚合
	OutboundSpeed int64 `json:"outbound_speed"` // 出站速率 (bytes/s)，即用户下载聚合
}

type MemInfo struct {
	Total int64 `json:"total"`
	Used  int64 `json:"used"`
}

type NodeMetrics struct {
	Uptime            int64       `json:"uptime"`
	Goroutines        int         `json:"goroutines"`
	ActiveConnections int         `json:"active_connections"`
	TotalConnections  int64       `json:"total_connections"`
	TcpConnections    int         `json:"tcp_connections"`
	TotalUsers        int         `json:"total_users"`
	ActiveUsers       int         `json:"active_users"`
	InboundSpeed      int64       `json:"inbound_speed"`
	OutboundSpeed     int64       `json:"outbound_speed"`
	CpuPerCore        []float64   `json:"cpu_per_core,omitempty"`
	Load              []float64   `json:"load,omitempty"`
	SpeedLimiter      interface{} `json:"speed_limiter,omitempty"`
	GC                interface{} `json:"gc,omitempty"`
	API               interface{} `json:"api,omitempty"`
	WS                interface{} `json:"ws,omitempty"`
	Limits            interface{} `json:"limits,omitempty"`
	KernelStatus      bool        `json:"kernel_status"`
}
