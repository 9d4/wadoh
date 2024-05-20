package mysqlstore

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"github.com/9d4/wadoh/devices"
	"github.com/9d4/wadoh/users"
)

type Store struct {
	Users   users.StorageProvider
	Devices devices.StorageProvider

	db *sql.DB
}

func New(dsn string) (*Store, error) {
	store, err := NewStore(dsn)
	if err != nil {
		return nil, err
	}

	if err := Upgrade(store.db); err != nil {
		return nil, err
	}
	return store, nil
}

func NewStore(dsn string) (*Store, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	return &Store{
		db:      db,
		Users:   newUsersStore(db),
		Devices: newDevicesStore(db),
	}, nil
}
