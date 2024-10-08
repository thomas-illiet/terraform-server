package router

import (
	"crypto/tls"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thomas-illiet/terrapi/controllers"
	"github.com/thomas-illiet/terrapi/pkg/config"
	"github.com/thomas-illiet/terrapi/pkg/middleware/header"
	"github.com/thomas-illiet/terrapi/pkg/middleware/prometheus"
)

// Load initializes the routing of the application.
func Load(cfg *config.Config) *gin.Engine {
	// Creates a router without any middleware by default
	r := gin.Default()

	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	r.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.Recovery())

	// Custom middleware
	r.Use(header.Version())
	r.Use(header.Cache())
	r.Use(header.Secure())
	r.Use(header.Options())

	v1 := r.Group("/v1")

	remote := v1.Group("/remote")
	remote.GET("/*wildcard", controllers.GetDeployments)
	remote.POST("/*wildcard", controllers.GetDeployments)
	remote.DELETE("/*wildcard", controllers.GetDeployments)
	remote.Handle("LOCK", "/*wildcard", controllers.GetDeployments)
	remote.Handle("UNLOCK", "/*wildcard", controllers.GetDeployments)

	return r
}

// Metrics initializes the routing of metrics and health.
func Metrics(cfg *config.Config) *gin.Engine {
	// Creates a router without any middleware by default
	r := gin.Default()

	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	r.Use(gin.Logger())

	// Custom middleware
	r.Use(header.Version())
	r.Use(header.Cache())
	r.Use(header.Secure())
	r.Use(header.Options())

	r.GET("/metrics", prometheus.Handler(cfg.Metrics.Token))

	r.GET("/healthz", func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "text/plain")
		c.Status(http.StatusOK)
		c.Writer.WriteHeaderNow()
		c.Abort()
	})

	r.GET("/readyz", func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "text/plain")
		c.Status(http.StatusOK)
		c.Writer.WriteHeaderNow()
		c.Abort()
	})

	return r
}

// Curves provides optionally a list of secure curves.
func Curves(cfg *config.Config) []tls.CurveID {
	if cfg.Server.StrictCurves {
		return []tls.CurveID{
			tls.CurveP521,
			tls.CurveP384,
			tls.CurveP256,
		}
	}

	return nil
}

// Ciphers provides optionally a list of secure ciphers.
func Ciphers(cfg *config.Config) []uint16 {
	if cfg.Server.StrictCiphers {
		return []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		}
	}

	return nil
}
