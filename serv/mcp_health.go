package serv

// registerHealthTools is retained as an inert compatibility hook.
func (ms *mcpServer) registerHealthTools() {
}

// HealthResult represents the health check response
type HealthResult struct {
	Status       string         `json:"status"`
	DatabaseType string         `json:"database_type,omitempty"`
	PingLatency  string         `json:"ping_latency,omitempty"`
	SchemaReady  bool           `json:"schema_ready"`
	TableCount   int            `json:"table_count"`
	PoolStats    *PoolStatsInfo `json:"pool_stats,omitempty"`
	Error        string         `json:"error,omitempty"`
}

// PoolStatsInfo represents database connection pool statistics
type PoolStatsInfo struct {
	MaxOpenConnections int    `json:"max_open_connections"`
	OpenConnections    int    `json:"open_connections"`
	InUse              int    `json:"in_use"`
	Idle               int    `json:"idle"`
	WaitCount          int64  `json:"wait_count"`
	WaitDuration       string `json:"wait_duration"`
}
