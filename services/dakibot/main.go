package main

import (
	"database/sql"
	"net/http"
	"os"
	"runtime"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"
	"github.com/rs/cors"
	"github.com/spf13/cobra"
	"github.com/rover/go/services/dakibot/internal"
	"github.com/rover/go/support/config"
	"github.com/rover/go/support/http/server"
	"github.com/rover/go/support/log"
	"github.com/rover/go/support/render/problem"
)

// Config represents the configuration of a dakibot server
type Config struct {
	Port              int               `toml:"port" valid:"required"`
	DakibotSecret   string            	`toml:"daki_secret" valid:"required"`
	NetworkPassphrase string            `toml:"network_passphrase" valid:"required"`
	HorizonURL        string            `toml:"horizon_url" valid:"required"`
	StartingBalance   string            `toml:"starting_balance" valid:"required"`
	TLS               *server.TLSConfig `valid:"optional"`
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	rootCmd := &cobra.Command{
		Use:   "dakibot",
		Short: "dakibot for the Rover Test Network",
		Long:  "client-facing api server for the dakibot service on the Rover Test Network",
		Run:   run,
	}

	rootCmd.PersistentFlags().String("conf", "./dakibot.cfg", "config file path")
	rootCmd.Execute()
}

func run(cmd *cobra.Command, args []string) {
	var (
		cfg     Config
		cfgPath = cmd.PersistentFlags().Lookup("conf").Value.String()
	)
	log.SetLevel(log.InfoLevel)
	err := config.Read(cfgPath, &cfg)
	if err != nil {
		switch cause := errors.Cause(err).(type) {
		case *config.InvalidConfigError:
			log.Error("config file: ", cause)
		default:
			log.Error(err)
		}
		os.Exit(1)
	}

	fb := initDakibot(cfg.DakibotSecret, cfg.NetworkPassphrase, cfg.HorizonURL, cfg.StartingBalance)
	router := initRouter(fb)
	registerProblems()

	server.Serve(router, cfg.Port, cfg.TLS)
}

func initRouter(fb *internal.Bot) *chi.Mux {
	routerConfig := server.EmptyConfig()

	// middleware
	server.AddBasicMiddleware(routerConfig)
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"*"},
	})
	routerConfig.Middleware(func(h http.Handler) http.Handler {
		return c.Handler(h)
	})

	// endpoints
	handler := &internal.DakibotHandler{Dakibot: fb}
	routerConfig.Route(http.MethodGet, "/", http.HandlerFunc(handler.Handle))
	routerConfig.Route(http.MethodPost, "/", http.HandlerFunc(handler.Handle))
	// not found handler
	routerConfig.NotFound(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		problem.Render(r.Context(), w, problem.NotFound)
	}))

	return server.NewRouter(routerConfig)
}

func registerProblems() {
	problem.RegisterError(sql.ErrNoRows, problem.NotFound)
}
