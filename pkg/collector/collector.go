package collector

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "container"
)

func Register(ctx context.Context)  {
	prometheus.MustRegister(newMountCollector(ctx))
}