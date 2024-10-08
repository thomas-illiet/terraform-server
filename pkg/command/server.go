package command

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"
	"time"

	"github.com/oklog/run"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/thomas-illiet/terrapi/pkg/database"
	"github.com/thomas-illiet/terrapi/pkg/router"
)

// Command-line flags default values
var (
	serverCmd = &cobra.Command{
		Use:   "server",
		Short: "Start integrated server",
		Run:   serverAction,
		Args:  cobra.NoArgs,
	}

	defaultMetricsAddr         = "0.0.0.0:8085"
	defaultServerAddr          = "0.0.0.0:8080"
	defaultServerCert          = ""
	defaultServerKey           = ""
	defaultServerStrictCurves  = false
	defaultServerStrictCiphers = false
	defaultEncryptionSecret    = ""
)

// Initialization of CLI flags and viper config binding
func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.PersistentFlags().String("metrics-addr", defaultMetricsAddr, "Address to bind the metrics")
	viper.SetDefault("metrics.addr", defaultMetricsAddr)
	viper.BindPFlag("metrics.addr", serverCmd.PersistentFlags().Lookup("metrics-addr"))

	serverCmd.PersistentFlags().String("metrics-token", "", "Token to secure metrics")
	viper.SetDefault("metrics.token", "")
	viper.BindPFlag("metrics.token", serverCmd.PersistentFlags().Lookup("metrics-token"))

	serverCmd.PersistentFlags().String("server-addr", defaultServerAddr, "Address to bind the server")
	viper.SetDefault("server.addr", defaultServerAddr)
	viper.BindPFlag("server.addr", serverCmd.PersistentFlags().Lookup("server-addr"))

	serverCmd.PersistentFlags().String("server-cert", defaultServerCert, "Path to certificate for SSL encryption")
	viper.SetDefault("server.cert", defaultServerCert)
	viper.BindPFlag("server.cert", serverCmd.PersistentFlags().Lookup("server-cert"))

	serverCmd.PersistentFlags().String("server-key", defaultServerKey, "Path to private key for SSL encryption")
	viper.SetDefault("server.key", defaultServerKey)
	viper.BindPFlag("server.key", serverCmd.PersistentFlags().Lookup("server-key"))

	serverCmd.PersistentFlags().Bool("strict-curves", defaultServerStrictCurves, "Enforce strict SSL curves")
	viper.SetDefault("server.strict_curves", defaultServerStrictCurves)
	viper.BindPFlag("server.strict_curves", serverCmd.PersistentFlags().Lookup("strict-curves"))

	serverCmd.PersistentFlags().Bool("strict-ciphers", defaultServerStrictCiphers, "Enforce strict SSL ciphers")
	viper.SetDefault("server.strict_ciphers", defaultServerStrictCiphers)
	viper.BindPFlag("server.strict_ciphers", serverCmd.PersistentFlags().Lookup("strict-ciphers"))

	serverCmd.PersistentFlags().String("encryption-secret", defaultEncryptionSecret, "Secret key for file encryption")
	viper.SetDefault("encryption.secret", defaultEncryptionSecret)
	viper.BindPFlag("encryption.secret", serverCmd.PersistentFlags().Lookup("encryption-secret"))
}

// Starts the server based on configuration and manages graceful shutdown
func serverAction(_ *cobra.Command, _ []string) {
	var gr run.Group

	// Connect to the database
	db := database.ConnectDb("library.db")
	database.CreateModel(db)
	database.InitDatabase(db)

	// Setup HTTP server
	var server http.Server
	if cfg.Server.Cert != "" && cfg.Server.Key != "" {
		server = createTLSServer() // Setup HTTPS server
	} else {
		server = createHTTPServer() // Setup HTTP server
	}

	// Add server to the run group for graceful shutdown handling
	gr.Add(func() error {
		log.Info().
			Str("addr", cfg.Server.Addr).
			Msg("Starting web server")

		// Start the appropriate server (HTTP or HTTPS)
		if cfg.Server.Cert != "" && cfg.Server.Key != "" {
			return server.ListenAndServeTLS("", "")
		} else {
			return server.ListenAndServe()
		}
	}, func(reason error) {
		// Gracefully shutdown the server with a timeout
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			log.Error().
				Err(err).
				Msg("Failed to shutdown web server gracefully")
		}

		log.Info().Err(reason).Msg("Web server shutdown gracefully")
	})

	metricServer := CreateMetricServer()
	gr.Add(func() error {
		log.Info().
			Str("addr", cfg.Metrics.Addr).
			Msg("Starting metrics server")

		return metricServer.ListenAndServe()
	}, func(reason error) {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		if err := metricServer.Shutdown(ctx); err != nil {
			log.Error().
				Err(err).
				Msg("Failed to shutdown metrics gracefully")

			return
		}

		log.Info().
			Err(reason).
			Msg("Shutdown metrics gracefully")
	})

	// Start the run group
	if err := gr.Run(); err != nil {
		log.Fatal().Err(err).Msg("Error running the server")
		os.Exit(1)
	}
}

// Sets up an HTTPS server with TLS configurations
func createTLSServer() http.Server {
	// Load the TLS certificate and private key
	cert, err := tls.LoadX509KeyPair(cfg.Server.Cert, cfg.Server.Key)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load certificates")
	}

	// Return HTTPS server configuration
	return http.Server{
		Addr:         cfg.Server.Addr,
		Handler:      router.Load(cfg),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		TLSConfig: &tls.Config{
			PreferServerCipherSuites: true,
			MinVersion:               tls.VersionTLS12,
			CurvePreferences:         router.Curves(cfg),
			CipherSuites:             router.Ciphers(cfg),
			Certificates:             []tls.Certificate{cert},
		},
	}
}

// Sets up a standard HTTP server
func createHTTPServer() http.Server {
	// Return HTTP server configuration
	return http.Server{
		Addr:         cfg.Server.Addr,
		Handler:      router.Load(cfg),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}

func CreateMetricServer() http.Server {
	return http.Server{
		Addr:         cfg.Metrics.Addr,
		Handler:      router.Metrics(cfg),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
}
