package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/9d4/wadoh/http"
	"github.com/9d4/wadoh/storage"
	"github.com/9d4/wadoh/users"
)

func init() {
	cobra.OnInitialize(initialize)
	setupLogger()

	persistent := rootCmd.PersistentFlags()
	persistent.StringVarP(&configFile, "config", "c", "", "config file to read")
}

func initialize() {
	setupConfig()
	setLogLevel(global.LogLevel)
}

var rootCmd = &cobra.Command{
	Use:   "wadoh",
	Short: "Start wadoh web server",
	Run: run(func(cmd *cobra.Command, args []string, storage *storage.Storage) {
		srv := http.NewServer(storage, func(c *http.Config) {
			c.Address = global.HTTP.Address
		})

		log.Debug().Any("a", global).Send()

		if err := srv.Serve(); err != nil {
			log.Err(err).Send()
		} else {
			log.Info().Str("addr", srv.Address()).Msg("http server listening")
		}

		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
		<-interrupt

		log.Info().Msg("shutting down")
		srv.ShutDown(context.Background())
		log.Info().Msg("exited")
	}),
}

func run(runFn func(*cobra.Command, []string, *storage.Storage)) func(cmd *cobra.Command, args []string) {
	fn := func(cmd *cobra.Command, args []string) {
		storage, err := getStorage()
		if err != nil {
			log.Fatal().Err(err).Send()
		}

		usrs, err := storage.Users.List(1, 0)
		if err == nil && len(usrs) < 1 {
			storage.Users.Save(&users.User{
				Name:      "admin",
				Username:  "Admin",
				Password:  "admin",
				CreatedAt: time.Now(),
			})
		}

		runFn(cmd, args, storage)
	}
	return fn
}
