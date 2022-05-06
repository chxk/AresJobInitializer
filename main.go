package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"

	"aresjob-initializer/cache"
	"aresjob-initializer/manual"

	"github.com/golang/glog"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Short: "Initializer for a panda record",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// For cobra + glog flags. Available to all subcommands.
			flag.Parse()
		},
	}
)

func main() {
	rand.Seed(time.Now().UnixNano())
	defer glog.Flush()
	rootCmd.AddCommand(cache.CacheCmd)
	rootCmd.AddCommand(manual.MpiCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)
}
