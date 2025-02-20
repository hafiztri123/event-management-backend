package health

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)


type Status string

const (
	StatusUp  Status = "UP"
	StatusDown Status = "DOWN"
)

type Component struct {
	Name    string      `json:"name"`
	Status  Status      `json:"status"`
	Message string      `json:"message,omitempty"`
	Details interface{} `json:"details,omitempty"`
}



type Response struct {
	Status     Status               `json:"status"`
	Components map[string]Component `json:"components"`
	Timestamp  time.Time            `json:"timestamp"`
	Version    string               `json:"version"`
	Uptime     string               `json:"uptime"`
}

type Checker interface {
	Name() string
	Check(ctx context.Context) Component
}

type HealthChecker struct {
	checkers 	[]Checker
	uptime 		time.Time
	version 	string
	mu 			sync.RWMutex
}

func NewHealthChecker(version string) *HealthChecker {
	return &HealthChecker{
		checkers: 	[]Checker{},
		uptime: 	time.Now(),
		version:	version, 
	}
}


func (h *HealthChecker) AddChecker(checker Checker) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.checkers = append(h.checkers, checker)
}

func (h *HealthChecker) CheckHealth(ctx context.Context) Response {
	h.mu.RLock()
	checkers := make([]Checker, len(h.checkers))
	copy(checkers, 	h.checkers)
	h.mu.RUnlock()

	components := make(map[string]Component)
	overallStatus := StatusUp

	var wg sync.WaitGroup
	var mu sync.Mutex

	timeout := 5 * time.Second
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	for _, checker := range checkers {
		wg.Add(1)
		go func(c Checker){
			defer wg.Done()

			component := c.Check(ctx)
			mu.Lock()
			components[c.Name()] = component
			if component.Status == StatusDown {
				overallStatus = StatusDown
			}
			mu.Unlock()

		}(checker)
	}

	wg.Wait()

	return Response{
		Status: overallStatus,
		Components: components,
		Timestamp: time.Now(),
		Version: h.version,
		Uptime: time.Since(h.uptime).String(),
	}
}

func (	h *HealthChecker) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){
		resp := h.CheckHealth(r.Context())
		w.Header().Set("Content-Type", "application/json")

		if resp.Status == StatusDown {
			w.WriteHeader(http.StatusServiceUnavailable)
		} else {
			w.WriteHeader(http.StatusOK)
		}

		json.NewEncoder(w).Encode(resp)
	}
	
}


type DatabaseChecker struct {
	db *gorm.DB
}

func NewDatabaseChecker(db *gorm.DB) *DatabaseChecker {
	return &DatabaseChecker{
		db: db,
	}
}

func (c *DatabaseChecker) Name() string {
	return "database"
}

func (c *DatabaseChecker) Check(ctx context.Context) Component {
	sqlDB, err := c.db.DB()
	if err != nil {
		return Component {
			Name: c.Name(),
			Status: StatusDown,
			Message: fmt.Sprintf("Failed to get database connection: %v", err),
		}
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return Component{
			Name: c.Name(),
			Status: StatusDown,
			Message: fmt.Sprintf("Database ping failed: %v", err),
		}
	}

	stats := sqlDB.Stats()
	return Component{
		Name: c.Name(),
		Status: StatusUp,
		Details: map[string]interface{}{
			"open_connections": stats.OpenConnections,
			"in_use": stats.InUse,
			"idle": stats.Idle,
			"max_open_connections": stats.MaxOpenConnections,
		},
	}
}

type RedisChecker struct {
	client *redis.Client
}

func NewRedisChecker(client  *redis.Client) *RedisChecker {
	return &RedisChecker{
		client: client,
	}
}

func (c *RedisChecker) Name() string {
	return "redis"
}

func (c *RedisChecker) Check(ctx context.Context) Component {
    ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
    defer cancel()

    // Ping Redis to check connectivity
    pong, err := c.client.Ping(ctx).Result()
    if err != nil || pong != "PONG" {
        return Component{
            Name:    c.Name(),
            Status:  StatusDown,
            Message: fmt.Sprintf("Redis ping failed: %v", err),
        }
    }

    // Fetch memory info
    info, err := c.client.Info(ctx, "memory").Result()
    if err != nil {
        return Component{
            Name:   c.Name(),
            Status: StatusUp,
        }
    }

    // Parse and format the memory info
    parsed := parseRedisInfo(info)
    formatted := formatMemoryInfo(parsed)

    return Component{
        Name:   c.Name(),
        Status: StatusUp,
        Details: map[string]interface{}{
            "info": formatted,
        },
    }
}

type DiskChecker struct {
	path string
}

func NewDiskChecker(path string) *DiskChecker {
	return &DiskChecker{path: path}
}

func (c *DiskChecker) Name() string {
	return "disk"
}

func (c *DiskChecker) Check(ctx context.Context) Component {
	// Get disk usage
	// This is a simplified version - in production you would use syscall or exec to get actual disk usage
	return Component{
		Name:   c.Name(),
		Status: StatusUp,
		Details: map[string]interface{}{
			"path": c.path,
			// In a real implementation, include:
			// "total": totalSpace,
			// "free": freeSpace,
			// "used": usedSpace,
		},
	}
}

// MemoryChecker checks memory usage
type MemoryChecker struct{}

// NewMemoryChecker creates a new memory health checker
func NewMemoryChecker() *MemoryChecker {
	return &MemoryChecker{}
}

func (c *MemoryChecker) Name() string {
	return "memory"
}

func (c *MemoryChecker) Check(ctx context.Context) Component {
	// Get memory stats
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	return Component{
		Name:   c.Name(),
		Status: StatusUp,
		Details: map[string]interface{}{
			"alloc":      m.Alloc,
			"total_alloc": m.TotalAlloc,
			"sys":        m.Sys,
			"num_gc":     m.NumGC,
		},
	}
}

func parseRedisInfo(info string) map[string]string {
    parsed := make(map[string]string)
    lines := strings.Split(info, "\r\n") // Split by lines
    for _, line := range lines {
        if strings.Contains(line, ":") { // Only process lines with key-value pairs
            parts := strings.SplitN(line, ":", 2) // Split into key and value
            if len(parts) == 2 {
                key := strings.TrimSpace(parts[0])
                value := strings.TrimSpace(parts[1])
                parsed[key] = value
            }
        }
    }
    return parsed
}

func formatMemoryInfo(parsed map[string]string) map[string]interface{} {
    formatted := make(map[string]interface{})

    // Helper function to convert bytes to human-readable format
    formatBytes := func(bytesStr string) string {
        bytes, err := strconv.ParseUint(bytesStr, 10, 64)
        if err != nil {
            return bytesStr // Return original value if parsing fails
        }
        const (
            KB = 1 << 10
            MB = 1 << 20
            GB = 1 << 30
        )
        switch {
        case bytes >= GB:
            return fmt.Sprintf("%.2fGB", float64(bytes)/GB)
        case bytes >= MB:
            return fmt.Sprintf("%.2fMB", float64(bytes)/MB)
        case bytes >= KB:
            return fmt.Sprintf("%.2fKB", float64(bytes)/KB)
        default:
            return fmt.Sprintf("%dB", bytes)
        }
    }

    // Memory Usage
    formatted["Used Memory"] = formatBytes(parsed["used_memory"])
    formatted["Used Memory RSS"] = formatBytes(parsed["used_memory_rss"])
    formatted["Memory Fragmentation Ratio"] = parsed["mem_fragmentation_ratio"]
    formatted["Max Memory"] = formatBytes(parsed["maxmemory"])
    formatted["Max Memory Policy"] = parsed["maxmemory_policy"]

    // Overhead
    formatted["Overhead Memory"] = formatBytes(parsed["used_memory_overhead"])
    formatted["Dataset Memory"] = formatBytes(parsed["used_memory_dataset"])

    // System Memory
    formatted["Total System Memory"] = formatBytes(parsed["total_system_memory"])

    return formatted
}