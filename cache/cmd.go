package cache

import (
	"time"

	"github.com/avast/retry-go"
	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type Config struct {
	DependedRoles     []string
	Interval          time.Duration
	Timeout           time.Duration
	DialTimeout       time.Duration
	ReadTimeout       time.Duration
	MaxRetries        int
	MPIHostFilePath   string // MPI的hostfile文件路径，可用逗号分割以写入多个文件
	MPIHostReplicas   int    // 每个node ip的重复次数
	MPIImplementation string // mpi实现
	HackFilePath      string // hack文件路径
}

func (c *Config) InstallFlags(fs *pflag.FlagSet) {
	fs.StringVar(&c.MPIHostFilePath, "write-mpi-hostfile-to", "", "write mpi worker deployment hostIP to write-mpi-hostfile-to, and will do nothing when get null string")
	fs.IntVar(&c.MPIHostReplicas, "mpi-host-replicas", 1, "mpi ranks per node")
	fs.StringVar(&c.MPIImplementation, "mpi-implementation", "", "mpi implementation")
	fs.StringVar(&c.HackFilePath, "hack-file", "/etc/ares-init", "filepath for hacking")
}

var (
	redisKey string
	cacheURI string
	config   Config
)

// CacheCmd 通过roleCache检测角色是否ready
var CacheCmd = &cobra.Command{
	Use:   "cache",
	Short: "initializer based on role cache",
	RunE: func(cmd *cobra.Command, args []string) error {
		glog.Infof("params: key=%v, cacheURI=%v, config=%+v", redisKey, cacheURI, config)

		var (
			c   CacheClient
			err error
		)
		if err = retry.Do(func() error {
			c, err = NewCacheClient(cacheURI, &config)
			if err != nil {
				return err
			}
			return Process(c, redisKey, config)
		}, retry.Attempts(uint(config.MaxRetries)), retry.Delay(config.Interval)); err != nil {
			glog.Errorf("failed to process: %+v", err)
			return err
		}
		defer c.Close()

		glog.Infof("initialize completed\n\n>>>")
		return nil
	},
}

func init() {
	CacheCmd.Flags().StringVar(&redisKey, "key", "", "redis key, e.g., rsdev:ares-task-2380-record-15706-dev")
	CacheCmd.Flags().StringVar(&cacheURI, "cacheURI", "", "cache URI")
	CacheCmd.Flags().StringSliceVar(&config.DependedRoles, "roles", nil, "roles to check, e.g. worker,launcher")
	CacheCmd.Flags().DurationVar(&config.Interval, "interval", time.Minute, "check interval")
	CacheCmd.Flags().DurationVar(&config.Timeout, "timeout", time.Hour, "wait timeout")
	CacheCmd.Flags().DurationVar(&config.DialTimeout, "dial-timeout", time.Minute, "dial timeout of connection")
	CacheCmd.Flags().DurationVar(&config.ReadTimeout, "read-timeout", time.Minute, "read timeout of connection")
	CacheCmd.Flags().IntVar(&config.MaxRetries, "max-retries", 10, "max number of retries")
	config.InstallFlags(CacheCmd.Flags())

	CacheCmd.MarkFlagRequired("key")
	CacheCmd.MarkFlagRequired("cacheURI")
}
