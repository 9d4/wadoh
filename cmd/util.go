package cmd

import (
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/9d4/wadoh/storage"
	"github.com/9d4/wadoh/storage/mysqlstore"
	wadohbe "github.com/9d4/wadoh/wadoh-be"
)

func setupLogger() {
	writer := zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
		w.TimeFormat = time.RFC3339
	})
	log.Logger = log.Output(writer)
}

func setLogLevel(l zerolog.Level) {
	log.Logger = log.Level(l)
}

func getStorage() (*storage.Storage, error) {
	rpc, err := wadohbe.NewClient(global.WadohBeAddress)
	if err != nil {
		log.Fatal().Caller().Err(err).Send()
	}

	switch global.Storage.Provider {
	case "mysql":
		store, err := mysqlstore.New(global.Storage.DSN)
		if err != nil {
			return nil, err
		}
		return storage.NewStorage(rpc.Service, store.Users, store.Devices), nil
	default:
		log.Fatal().Msg("unsupported storage provider: " + global.Storage.Provider)
	}
	return nil, nil
}
