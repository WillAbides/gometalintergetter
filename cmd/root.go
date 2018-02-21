package cmd

import (
	"fmt"
	"os"

	"github.com/WillAbides/gometalintergetter/getter"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"runtime"
)

var (
	version, installDir, targetOS, targetArch string
	force                                     bool
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gometalintergetter [version=latest]",
	Short: "install gometalinter",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Inside rootCmd Run with args: %v\n", args)

		if len(args) == 1 {
			version = args[0]
		}

		fmt.Printf("version: %v\n", version)
		fmt.Printf("installPath: %v\n", installDir)
		opts := []getter.Option{
			getter.WithOS(targetOS),
			getter.WithArch(targetArch),
		}
		if force {
			opts = append(opts, getter.WithForce())
		}
		err := getter.DownloadMetalinter(version, installDir, opts...)
		if err != nil {
			panic(err)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.Flags().StringVarP(&installDir, "installdir", "i", ".", "directory where gometalinter will be installed")
	rootCmd.Flags().StringVarP(&targetOS, "os", "o", runtime.GOOS, "target operating system for gometalinter")
	rootCmd.Flags().StringVarP(&targetArch, "arch", "a", runtime.GOARCH, "target system architecture for gometalinter")
	rootCmd.Flags().BoolVarP(&force, "force", "f", false, "force download even if we already have the specified version")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.AutomaticEnv() // read in environment variables that match
}
