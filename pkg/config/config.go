package config

type Config struct {
	ListenAddress string
	MaxRequests   int
	Help          bool
	NameSpace     string
	AllNamespaces  bool
	MountPaths     []string
	NodeIP       string
	SocketPath   string
}
