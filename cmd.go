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

var serverEchoCmd = &cobra.Command{
	Use:     "echo",
	Aliases: []string{"e"},
	Short:   "Generate and echo server on the given ip address:port.",
	Args:    cobra.MaximumNArgs(0),
	RunE: func(cmd *cobra.Command, args []string) error {
		return startEchoServer(listeningAddress)
	},
}

func init() {
	serverEchoCmd.Flags().StringVarP(&listeningAddress, "listen", "l", listeningAddress, "The address which the server will listen")
	serverCmd.AddCommand(serverEchoCmd)
	rootCmd.AddCommand(serverCmd)
}
