package router

import (
	"crypto/tls"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/thomas-illiet/terrapi/pkg/config"
	"github.com/thomas-illiet/terrapi/pkg/handler"
	"github.com/thomas-illiet/terrapi/pkg/middleware/header"
	"github.com/thomas-illiet/terrapi/pkg/middleware/prometheus"
	"gorm.io/gorm"
)

// Load initializes the routing of the application.
func Load(db *gorm.DB, cfg *config.Config) *gin.Engine {
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

	//
	r.Use(ApiMiddleware(db, cfg))

	v1 := r.Group("/v1")

	remote := v1.Group("/states")
	remote.GET("/*wildcard", handler.StateFetch)
	remote.POST("/*wildcard", handler.StateUpdate)
	remote.DELETE("/*wildcard", handler.StateDelete)
	remote.Handle("LOCK", "/*wildcard", handler.StateLock)
	remote.Handle("UNLOCK", "/*wildcard", handler.StateUnlock)

	module := v1.Group("/modules")
	module.GET("", handler.ModuleFetchs)
	module.POST("", handler.ModuleCreate)
	module.GET("/:id", handler.ModuleFetch)
	module.DELETE("/:id", handler.ModuleDelete)

	deployment := v1.Group("/deployments")
	deployment.GET("", handler.DeploymentFetchs)
	deployment.POST("", handler.DeploymentCreate)
	deployment.GET("/:id", handler.DeploymentFetch)
	deployment.DELETE("/:id", handler.DeploymentDelete)

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

// ApiMiddleware will add the db connection to the context
func ApiMiddleware(db *gorm.DB, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("databaseConn", db)
		c.Set("config", cfg)
		c.Next()
	}
}
