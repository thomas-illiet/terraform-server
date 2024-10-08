package command

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// Configures the global logger based on settings from the
// configuration file or environment variables.
func setupLogger() error {
	// Set the log level based on the configuration value
	switch strings.ToLower(viper.GetString("log.level")) {
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	// If pretty logging is enabled, configure the logger to use a console writer
	if viper.GetBool("log.pretty") {
		log.Logger = log.Output(
			zerolog.ConsoleWriter{
				Out:     os.Stderr,
				NoColor: !viper.GetBool("log.color"),
			},
		)
	}

	return nil
}

// Loads the application configuration using Viper.
// The configuration is loaded from a file specified by "config.file",
// or from default paths if no file is specified. It also reads environment
// variables with the prefix "terrapi" and replaces "." with "_" for compatibility.
func setupConfig() {
	if viper.GetString("config.file") != "" {
		viper.SetConfigFile(viper.GetString("config.file"))
	} else {
		viper.SetConfigName("config")
		viper.AddConfigPath("/etc/terrapi")
		viper.AddConfigPath("$HOME/.terrapi")
		viper.AddConfigPath("./terrapi")
	}

	viper.SetEnvPrefix("terrapi")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if err := readConfig(); err != nil {
		log.Error().
			Err(err).
			Msg("Failed to read config file")
	}

	if err := viper.Unmarshal(cfg); err != nil {
		log.Error().
			Err(err).
			Msg("Failed to parse config file")
	}
}

// Attempts to read the configuration from the file specified.
// It returns nil if the config is read successfully, or if the file is not found.
// Other errors, such as file read errors, are returned.
func readConfig() error {
	err := viper.ReadInConfig()

	// Return nil if the config was read successfully
	if err == nil {
		return nil
	}

	// Return nil if the config file was not found
	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		return nil
	}

	// Return nil if there was a file path error
	if _, ok := err.(*os.PathError); ok {
		return nil
	}

	// Return the error for other issues
	return err
}
