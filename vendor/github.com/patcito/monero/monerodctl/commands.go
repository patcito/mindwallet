package main

import (
	"fmt"
	"github.com/patcito/monero/rpc"
	"github.com/spf13/cobra"
	"os"
	"strconv"
)

var heightCmd = &cobra.Command{
	Use:   "height",
	Short: "Gets the current block height",
	Run:   height,
}

func height(cmd *cobra.Command, args []string) {
	h, err := node.BlockchainHeight()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Fprintln(os.Stdout, h)
}

var infoCmd = &cobra.Command{
	Use:   "info",
	Short: "Gets daemon information",
	Run:   info,
}

func info(cmd *cobra.Command, args []string) {
	node = rpc.NewNode(nodeAddr + ":" + nodePort)

	i, err := node.Info()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stdout, ""+
		"Height: % 17d\n"+
		"Difficulty: % 13d\n"+
		"TxCount: % 16d\n"+
		"TxPoolSize: % 13d\n"+
		"AltBlocksCount: % 9d\n"+
		"OutgoingConnections: % 4d\n"+
		"IncomingConnections: % 4d\n"+
		"WhitePeerlistSize: % 6d\n"+
		"GreyPeerlistSize:  % 6d\n",
		i.Height,
		i.Difficulty,
		i.TxCount,
		i.TxPoolSize,
		i.AltBlocksCount,
		i.OutgoingConnections,
		i.IncomingConnections,
		i.WhitePeerlistSize,
		i.GreyPeerlistSize,
	)
}

var miningCmd = &cobra.Command{
	Use:   "mining",
	Short: "control mining",
	Run:   mining,
}

func mining(cmd *cobra.Command, args []string) {
	cmd.Usage()
	os.Exit(1)
}

var startMiningCmd = &cobra.Command{
	Use:   "start <address> <threads>",
	Short: "Start solo mining to given address",
	Run:   startMining,
}

func startMining(cmd *cobra.Command, args []string) {
	var (
		addr    string
		threads int
		err     error
	)
	switch len(args) {
	case 1:
		addr = args[0]
		threads = 1
	case 2:
		addr = args[0]
		threads, err = strconv.Atoi(args[1])
		if err != nil {
			fmt.Fprintln(os.Stderr, "invalid thread count,", err)
			cmd.Usage()
			os.Exit(1)
		}
	default:
		cmd.Usage()
		os.Exit(1)
	}

	node = rpc.NewNode(nodeAddr + ":" + nodePort)
	err = node.StartMining(addr, threads)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to start mining,", err)
		os.Exit(1)
	}

}

var stopMiningCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop mining",
	Run:   stopMining,
}

func stopMining(cmd *cobra.Command, args []string) {
	node = rpc.NewNode(nodeAddr + ":" + nodePort)
	err := node.StopMining()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var blockCountCmd = &cobra.Command{
	Use:   "blockcount",
	Short: "block count",
	Run:   blockCount,
}

func blockCount(cmd *cobra.Command, args []string) {
	node = rpc.NewNode(nodeAddr + ":" + nodePort)
	n, err := node.BlockCount()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Fprintln(os.Stdout, n)
}
