package main

import (
	"log"
	"runtime"

	"github.com/PuerkitoBio/throttled"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/rover/go/services/orbit/internal"
	hlog "github.com/rover/go/services/orbit/internal/log"
)

var app *orbit.App
var config orbit.Config

var rootCmd *cobra.Command

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	rootCmd.Execute()
}

func init() {
	viper.SetDefault("port", 8000)
	viper.SetDefault("history-retention-count", 0)

	viper.BindEnv("port", "PORT")
	viper.BindEnv("db-url", "DATABASE_URL")
	viper.BindEnv("rover-core-db-url", "STELLAR_CORE_DATABASE_URL")
	viper.BindEnv("rover-core-url", "STELLAR_CORE_URL")
	viper.BindEnv("per-hour-rate-limit", "PER_HOUR_RATE_LIMIT")
	viper.BindEnv("redis-url", "REDIS_URL")
	viper.BindEnv("ruby-orbit-url", "RUBY_HORIZON_URL")
	viper.BindEnv("dakibot-url", "DAKIBOT_URL")
	viper.BindEnv("log-level", "LOG_LEVEL")
	viper.BindEnv("sentry-dsn", "SENTRY_DSN")
	viper.BindEnv("loggly-token", "LOGGLY_TOKEN")
	viper.BindEnv("loggly-host", "LOGGLY_HOST")
	viper.BindEnv("tls-cert", "TLS_CERT")
	viper.BindEnv("tls-key", "TLS_KEY")
	viper.BindEnv("ingest", "INGEST")
	viper.BindEnv("network-passphrase", "NETWORK_PASSPHRASE")
	viper.BindEnv("history-retention-count", "HISTORY_RETENTION_COUNT")
	viper.BindEnv("history-stale-threshold", "HISTORY_STALE_THRESHOLD")
	viper.BindEnv("skip-cursor-update", "SKIP_CURSOR_UPDATE")

	rootCmd = &cobra.Command{
		Use:   "orbit",
		Short: "client-facing api server for the rover network",
		Long:  "client-facing api server for the rover network",
		Run: func(cmd *cobra.Command, args []string) {
			initApp(cmd, args)
			app.Serve()
		},
	}

	rootCmd.PersistentFlags().String(
		"db-url",
		"",
		"orbit postgres database to connect with",
	)

	rootCmd.PersistentFlags().String(
		"rover-core-db-url",
		"",
		"rover-core postgres database to connect with",
	)

	rootCmd.PersistentFlags().String(
		"rover-core-url",
		"",
		"rover-core to connect with (for http commands)",
	)

	rootCmd.PersistentFlags().Int(
		"port",
		8000,
		"tcp port to listen on for http requests",
	)

	rootCmd.PersistentFlags().Int(
		"per-hour-rate-limit",
		3600,
		"max count of requests allowed in a one hour period, by remote ip address",
	)

	rootCmd.PersistentFlags().String(
		"redis-url",
		"",
		"redis to connect with, for rate limiting",
	)

	rootCmd.PersistentFlags().String(
		"dakibot-url",
		"",
		"dakibot service to redirect to",
	)

	rootCmd.PersistentFlags().String(
		"log-level",
		"info",
		"Minimum log severity (debug, info, warn, error) to log",
	)

	rootCmd.PersistentFlags().String(
		"sentry-dsn",
		"",
		"Sentry URL to which panics and errors should be reported",
	)

	rootCmd.PersistentFlags().String(
		"loggly-token",
		"",
		"Loggly token, used to configure log forwarding to loggly",
	)

	rootCmd.PersistentFlags().String(
		"loggly-host",
		"",
		"Hostname to be added to every loggly log event",
	)

	rootCmd.PersistentFlags().String(
		"tls-cert",
		"",
		"The TLS certificate file to use for securing connections to orbit",
	)

	rootCmd.PersistentFlags().String(
		"tls-key",
		"",
		"The TLS private key file to use for securing connections to orbit",
	)

	rootCmd.PersistentFlags().Bool(
		"ingest",
		false,
		"causes this orbit process to ingest data from rover-core into orbit's db",
	)

	rootCmd.PersistentFlags().String(
		"network-passphrase",
		"",
		"Override the network passphrase",
	)

	rootCmd.PersistentFlags().Uint(
		"history-retention-count",
		0,
		"the minimum number of ledgers to maintain within orbit's history tables.  0 signifies an unlimited number of ledgers will be retained",
	)

	rootCmd.PersistentFlags().Uint(
		"history-stale-threshold",
		0,
		"the maximum number of ledgers the history db is allowed to be out of date from the connected rover-core db before orbit considers history stale",
	)

	rootCmd.AddCommand(dbCmd)

	viper.BindPFlags(rootCmd.PersistentFlags())
}

func initApp(cmd *cobra.Command, args []string) {
	initConfig()

	var err error
	app, err = orbit.NewApp(config)

	if err != nil {
		log.Fatal(err.Error())
	}
}

func initConfig() {
	if viper.GetString("db-url") == "" {
		log.Fatal("Invalid config: db-url is blank.  Please specify --db-url on the command line or set the DATABASE_URL environment variable.")
	}

	if viper.GetString("rover-core-db-url") == "" {
		log.Fatal("Invalid config: rover-core-db-url is blank.  Please specify --rover-core-db-url on the command line or set the STELLAR_CORE_DATABASE_URL environment variable.")
	}

	if viper.GetString("rover-core-url") == "" {
		log.Fatal("Invalid config: rover-core-url is blank.  Please specify --rover-core-url on the command line or set the STELLAR_CORE_URL environment variable.")
	}

	ll, err := logrus.ParseLevel(viper.GetString("log-level"))

	if err != nil {
		log.Fatalf("Could not parse log-level: %v", viper.GetString("log-level"))
	}

	hlog.DefaultLogger.Level = ll

	cert, key := viper.GetString("tls-cert"), viper.GetString("tls-key")

	switch {
	case cert != "" && key == "":
		log.Fatal("Invalid TLS config: key not configured")
	case cert == "" && key != "":
		log.Fatal("Invalid TLS config: cert not configured")
	}

	config = orbit.Config{
		DatabaseURL:            viper.GetString("db-url"),
		StellarCoreDatabaseURL: viper.GetString("rover-core-db-url"),
		StellarCoreURL:         viper.GetString("rover-core-url"),
		Port:                   viper.GetInt("port"),
		RateLimit:              throttled.PerHour(viper.GetInt("per-hour-rate-limit")),
		RedisURL:               viper.GetString("redis-url"),
		DakibotURL:             viper.GetString("dakibot-url"),
		LogLevel:               ll,
		SentryDSN:              viper.GetString("sentry-dsn"),
		LogglyToken:            viper.GetString("loggly-token"),
		LogglyHost:             viper.GetString("loggly-host"),
		TLSCert:                cert,
		TLSKey:                 key,
		Ingest:                 viper.GetBool("ingest"),
		HistoryRetentionCount:  uint(viper.GetInt("history-retention-count")),
		StaleThreshold:         uint(viper.GetInt("history-stale-threshold")),
		SkipCursorUpdate:       viper.GetBool("skip-cursor-update"),
	}
}
