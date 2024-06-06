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
	wadohbe "github.com/9d4/wadoh/wadoh-be"
)

func init() {
	cobra.OnInitialize(initialize)

	persistent := rootCmd.PersistentFlags()
	persistent.StringVarP(&customConfigFile, "config", "c", "", "config file to read")
}

func initialize() {
	setupConfig()
	setLogLevel(global.LogLevel)
	setupLogger()
}

var rootCmd = &cobra.Command{
	Use:   "wadoh",
	Short: "Start wadoh web server",
	Run: run(func(cmd *cobra.Command, args []string, storage *storage.Storage) {
		rpc, err := wadohbe.NewClient(global.WadohBeAddress)
		if err != nil {
			log.Fatal().Err(err).Send()
		}

		srv := http.NewServer(storage, rpc.Service, func(c *http.Config) {
			*c = global.HTTP
		})

		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

		// Start http server
		if err := srv.Serve(); err != nil {
			log.Err(err).Send()
			stop()
		} else {
			log.Info().Str("addr", srv.Address()).Msg("http server listening")
		}

		// Start device message listener
		go func() {
			for a := range rpc.ReceiveMessage() {
				log.Debug().Any("cmd:recvmessage", a.GetFrom()).Send()
			}
		}()

		<-ctx.Done()

		log.Info().Msg("shutting down")
		if err := srv.ShutDown(context.Background()); err != nil {
			log.Err(err).Send()
		}
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
			_ = storage.Users.Save(&users.User{
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
