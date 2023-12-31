/*
Copyright © 2023 Mahdi Lotfi mahdilotfi167@gmail.com
*/

package cmd

import (
	"fmt"
	"nsproxy/config"
	"nsproxy/internal/cache"
	"nsproxy/internal/server"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "nsproxy [host[:port]]",
	Args:  cobra.MatchAll(cobra.MaximumNArgs(1)),
	Short: "A simple and lightweight DNS proxy",
	Long: `Nowadays, due to the huge increase in websites and people’s use of them,
from a point-of-view of a network engineer, a Huge number of domains
need to be resolved by DNS servers.
As an interesting fact, a typical home connected to the internet makes ~10k
DNS queries per day!
So there is a huge load on DNS servers.
nsproxy forwards DNS requests(only if needed) and replies between DNS clients and
DNS servers.`,
	Run: func(cmd *cobra.Command, args []string) {
		addr := "0.0.0.0:53"
		if len(args) > 0 {
			addr = args[0]
		}

		var serverConfig config.ServerConfig
		if err := viper.Unmarshal(&serverConfig); err != nil {
			fmt.Fprintln(os.Stderr, "Bad config file")
			return
		}

		var cacheConfig config.CacheConfig
		if err := viper.Unmarshal(&cacheConfig); err != nil {
			fmt.Fprintln(os.Stderr, "Bad config file")
			return
		}

		cacheManager := cache.NewCacheManager(&cacheConfig)

		server := server.NewServer(addr, &serverConfig, cacheManager)
		server.Run()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is /etc/nsproxy.json)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.SetDefault("cache-expiration-time", 60*60*24)
	viper.SetDefault("external-dns-servers", []string{"8.8.8.8:53"})

	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath("/etc")
		viper.SetConfigType("json")
		viper.SetConfigName("nsproxy")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
