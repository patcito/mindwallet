package main

import (
	"github.com/patcito/monero/rpc"
	"github.com/spf13/cobra"
)

const Version = "0.1.0"

var rootCmd = &cobra.Command{Use: "monerodctl"}

var (
	nodeAddr, nodePort string

	node *rpc.Node
)

func init() {
	miningCmd.AddCommand(startMiningCmd, stopMiningCmd)

	rootCmd.AddCommand(
		heightCmd,
		infoCmd,
		miningCmd,
		blockCountCmd,
	)

	rootCmd.PersistentFlags().StringVarP(&nodeAddr, "daemon-address", "a", "127.0.0.1", "daemon host or ip address")
	rootCmd.PersistentFlags().StringVarP(&nodePort, "daemon-port", "p", "18081", "daemon port number")
}

func main() {
	rootCmd.Execute()
}
