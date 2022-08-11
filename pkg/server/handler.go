package server

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
)

func index() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(`
			<html>
				<head><title>HostPath Exporter</title></head>
				<body>
					<h1>HostPath Exporter</h1>
					<p><a href="` + "/metrics" + `">Metrics</a></p>
				</body>
			</html>`))
	}
}

func metrics() gin.HandlerFunc {
	return gin.WrapH(promhttp.Handler())
}