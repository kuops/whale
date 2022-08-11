package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"whale/pkg/config"
	"whale/pkg/server"

	"github.com/spf13/pflag"
)

func main() {
	whaleConfig := &config.Config{}
	pflag.StringVarP(&whaleConfig.ListenAddress, "listen-address", "l", ":8080", "listen address.")
	pflag.IntVarP(&whaleConfig.MaxRequests, "max-requests", "", 40, "max http requests.")
	pflag.BoolVarP(&whaleConfig.Help, "help", "h", false, "display help for whale.")
	pflag.StringVarP(&whaleConfig.NameSpace, "namespace", "n", "default", "collector namespace pods.")
	pflag.BoolVarP(&whaleConfig.AllNamespaces, "all-namespaces", "A", false, "all namespace pods.")
	pflag.StringVarP(&whaleConfig.SocketPath, "socket-path", "", "unix:///run/containerd/containerd.sock", "container runtime interface socket path.")
	pflag.StringVarP(&whaleConfig.NodeIP, "node-ip", "", os.Getenv("NODE_IP"), "running node ip.")
	pflag.StringSliceVarP(&whaleConfig.MountPaths, "mount-paths", "p", []string{"/data/logs"}, "collector container mount paths.")
	pflag.Parse()

	if whaleConfig.NodeIP == "" {
		log.Println("node-ip is empty.")
		os.Exit(5)
	}

	if whaleConfig.Help {
		pflag.PrintDefaults()
		os.Exit(0)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	err := server.Run(ctx, whaleConfig)
	stop()

	if err != nil {
		log.Fatalln(err)
	}
}
