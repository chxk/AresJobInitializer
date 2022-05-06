package manual

import (
	"fmt"

	"aresjob-initializer/cache"
	"aresjob-initializer/utils"

	"github.com/golang/glog"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type Config struct {
	cache.Config
	IPs []string
}

func (c *Config) InstallFlags(fs *pflag.FlagSet) {
	c.Config.InstallFlags(fs)
	fs.StringSliceVar(&c.IPs, "ips", nil, "array of IP addresses, e.g. 127.0.0.1,127.0.0.2")
}

var (
	config Config
)

// MpiCmd: 生成MPI框架的hostfile文件
var MpiCmd = &cobra.Command{
	Use:   "mpi",
	Short: "generate mpi host file",
	RunE: func(cmd *cobra.Command, args []string) error {
		glog.Infof("params: config=%+v", config)

		if len(config.MPIHostFilePath) > 0 {
			if err := utils.WriteMPIHostFile(config.MPIImplementation, config.MPIHostFilePath, config.IPs, config.MPIHostReplicas); err != nil {
				return fmt.Errorf("failed to write hostfile: %v", err)
			}
		}
		if err := utils.CreateFile(config.HackFilePath); err != nil {
			return fmt.Errorf("failed to create hackfile: %v", err)
		}

		glog.Infof("generate completed")
		return nil
	},
}

func init() {
	config.InstallFlags(MpiCmd.Flags())
}
