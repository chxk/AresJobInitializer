package cache

type RoleCache struct {
	Status  RoleStatus `json:"status"`
	Running bool       `json:"running"`
	Phase   string     `json:"phase"`
	IPs     []string   `json:"ips"`
}

type JobCache map[string]*RoleCache

type RoleStatus string

const (
	RoleStatusNotFound   RoleStatus = "NotFound"
	RoleStatusNotRunning RoleStatus = "NotRunning"
	RoleStatusRunning    RoleStatus = "Running"
)

type CacheClient interface {
	GetJobCache(key string) (JobCache, error)
	Close()
}
