package cache

import (
	"errors"
	"fmt"
	"math/rand"
	"net/url"
	"time"

	"aresjob-initializer/utils"

	"github.com/golang/glog"
)

func NewCacheClient(uri string, config *Config) (CacheClient, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, err
	}
	switch u.Scheme {
	case "nodePort", "nodeport":
		return NewNodePortClient("http", u.Host, config)
	case "http", "https":
		return NewObserverClient(u.String(), config)
	default:
		return nil, fmt.Errorf("unknown backend for cache: %s", u.Scheme)
	}
}

func GetDependedRoleCaches(cache CacheClient, key string, roles []string) (JobCache, error) {
	glog.Infof("fetching: key=%v roles=%v", key, roles)
	jobCache, err := cache.GetJobCache(key)
	if err != nil {
		return nil, err
	}
	if jobCache == nil {
		glog.Warningf("no role cache was found: %v", key)
		return nil, nil
	}

	currentStatus := make(map[string]RoleStatus)
	result := JobCache{}
	for _, role := range roles {
		if jobCache[role] == nil {
			currentStatus[role] = RoleStatusNotFound
		} else if jobCache[role].Running {
			currentStatus[role] = RoleStatusRunning
			result[role] = jobCache[role]
		} else {
			currentStatus[role] = RoleStatusNotRunning
		}
	}
	glog.Infof("current status: %+v", currentStatus)
	if len(roles) == len(result) {
		return result, nil
	}
	return nil, nil

}

func WaitDependedRoles(cache CacheClient, key string, config Config) (JobCache, error) {
	start := time.Now()
	for {
		consumed := time.Now().Sub(start)
		if config.Timeout > 0 && consumed >= config.Timeout {
			msg := fmt.Sprintf("timeout(%v): %v", config.Timeout, consumed)
			glog.Error(msg)
			return nil, errors.New(msg)
		} else {
			glog.Infof("checking(%v)...", consumed)
		}

		jobCache, err := GetDependedRoleCaches(cache, key, config.DependedRoles)
		if err != nil {
			glog.Errorf("failed to check roles(%v) ready or not: %v", config.DependedRoles, err)
			return nil, err
		}
		if jobCache != nil {
			glog.Infof("ready!")
			return jobCache, nil
		}
		interval := randDuration(config.Interval)
		glog.Infof("going to wait %v", interval)
		time.Sleep(interval)
	}
	return nil, nil
}

func Process(cache CacheClient, key string, config Config) error {
	err := HandleDependedRoleCaches(cache, key, config)
	if err == nil {
		return nil
	}
	// hack: 人为干预手段，如果创建该文件，则正常退出
	if utils.Exists(config.HackFilePath) {
		glog.Warningf("hack: going to exit.")
		return nil
	} else {
		glog.Infof("hacking command: \n\t"+
			"./initializer mpi --logtostderr --write-mpi-hostfile-to='%s' --mpi-host-replicas=%d --mpi-implementation=%s --hack-file=%s --ips=xxx",
			config.MPIHostFilePath, config.MPIHostReplicas, config.MPIImplementation, config.HackFilePath,
		)
	}
	return err
}

func HandleDependedRoleCaches(cache CacheClient, key string, config Config) error {
	var (
		jobCache JobCache
		err      error
	)
	if len(config.DependedRoles) == 0 {
		glog.Infof("no depended roles, will skip waiting")
	} else {
		jobCache, err = WaitDependedRoles(cache, key, config)
		if err != nil {
			return err
		}
	}

	if len(config.MPIHostFilePath) == 0 {
		return nil
	}
	ips := []string{}
	for _, cache := range jobCache {
		for _, ip := range cache.IPs {
			if len(ip) == 0 {
				continue
			}
			ips = append(ips, ip)
		}
	}
	if err := utils.WriteMPIHostFile(config.MPIImplementation, config.MPIHostFilePath, ips, config.MPIHostReplicas); err != nil {
		return fmt.Errorf("failed to write hostfile: %v", err)
	}
	return nil
}

func randDuration(interval time.Duration) time.Duration {
	if interval <= 0 {
		return 0
	}
	disturbance := interval / 2
	min, max := interval-disturbance, interval+disturbance
	return min + time.Duration(rand.Int63n(int64(max-min)))
}
