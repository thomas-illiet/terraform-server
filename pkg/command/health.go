package command

import (
	"fmt"
	"net/http"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	// Defines the command for performing health checks.
	healthCmd = &cobra.Command{
		Use:   "health",
		Short: "Perform health checks on the application",
		Run:   healthAction,
		Args:  cobra.NoArgs, // No arguments are required for this command.
	}
)

func init() {
	rootCmd.AddCommand(healthCmd)
	healthCmd.PersistentFlags().String("metrics-addr", defaultMetricsAddr, "Address to bind the metrics")
	viper.SetDefault("metrics.addr", defaultMetricsAddr)
	viper.BindPFlag("metrics.addr", healthCmd.PersistentFlags().Lookup("metrics-addr"))
}

func healthAction(_ *cobra.Command, _ []string) {
	// Making an HTTP GET request to the health endpoint.
	resp, err := http.Get(
		fmt.Sprintf(
			"http://%s/healthz", // Health endpoint URL
			cfg.Metrics.Addr,    // Address of the metrics service
		),
	)

	// Handling errors that might occur during the HTTP request.
	if err != nil {
		log.Error().
			Err(err).
			Msg("Failed to request health check")
		os.Exit(1)
	}

	// Ensuring the response body is closed when we're done.
	defer resp.Body.Close()

	// Checking if the response status code indicates failure.
	if resp.StatusCode != http.StatusOK {
		log.Error().
			Int("code", resp.StatusCode).
			Msg("Health check failed, service might be unhealthy")
		os.Exit(1)
	}
}
