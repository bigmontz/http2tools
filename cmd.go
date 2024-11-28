package main

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "http2tools",
	Short: "Toolbox for HTTP/2.0",
	Long:  "Toolbox with client and server tools for HTTP/2.0.",
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "HTTP/2.0 servers.",
	Long:  "Commands for starting HTTP/2.0 servers",
}

var listeningAddress string = "0.0.0.0:1010"
var useTcp bool = false

var serverEchoCmd = &cobra.Command{
	Use:     "echo",
	Aliases: []string{"e"},
	Short:   "Generate an echo server on the given ip address:port.",
	Args:    cobra.MaximumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		if useTcp {
			return startTcpEchoServer(listeningAddress)
		}
		return startEchoServer(listeningAddress)
	},
}

var targetAddress string = "localhost:1011"
var useTls bool = false

var serverProxyCmd = &cobra.Command{
	Use:     "proxy",
	Aliases: []string{"p"},
	Short:   "Generate an proxy server on the given ip address:port.",
	Args:    cobra.MaximumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		return startTcpProxyServer(listeningAddress, targetAddress, useTls)
	},
}

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "HTTP/2.0 clients.",
	Long:  "Commands for connection to servers using HTTP/2.0 and trying out capabilities",
}

var numberOfBytesToSend int = 4048
var batchSize int = 1024
var timeBetweenBatches int = 0

var clientRandomCmd = &cobra.Command{
	Use:     "random",
	Aliases: []string{"r"},
	Short:   "Send random data in the body of a HTTP/2.0 request.",
	Long:    "Send random data in the body of a HTTP/2.0 request.",
	Args:    cobra.MinimumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		return connectionUsingRandomClient(args[0], numberOfBytesToSend, batchSize, timeBetweenBatches)
	},
}

func init() {
	serverEchoCmd.Flags().StringVarP(&listeningAddress, "listen", "l", listeningAddress, "The address which the server will listen")
	serverEchoCmd.Flags().BoolVarP(&useTcp, "tcp", "t", useTcp, "If server should use TCP instead instead HTTP/2.0")
	serverCmd.AddCommand(serverEchoCmd)
	serverProxyCmd.Flags().StringVarP(&listeningAddress, "listen", "l", listeningAddress, "The address which the server will listen")
	serverProxyCmd.Flags().StringVarP(&targetAddress, "target", "t", targetAddress, "The address which the server will connect to")
	serverProxyCmd.Flags().BoolVarP(&useTls, "tls", "s", useTls, "If the proxy should connect use tls to the target address")
	serverCmd.AddCommand(serverProxyCmd)
	rootCmd.AddCommand(serverCmd)

	clientRandomCmd.Flags().IntVarP(&numberOfBytesToSend, "bytes", "b", numberOfBytesToSend, "The number of bytes to be send (-1 to inf)")
	clientRandomCmd.Flags().IntVarP(&batchSize, "batch", "s", batchSize, "The number of bytes per batch")
	clientRandomCmd.Flags().IntVarP(&timeBetweenBatches, "interval", "i", timeBetweenBatches, "Interval between batches in ms")
	clientCmd.AddCommand(clientRandomCmd)
	rootCmd.AddCommand(clientCmd)
}
