package storage

import (
	"github.com/9d4/wadoh/devices"
	"github.com/9d4/wadoh/users"
)

type Storage struct {
	Users   *users.Storage
	Devices *devices.Storage
}

func NewStorage(usersP users.StorageProvider, devicesP devices.StorageProvider) *Storage {
	return &Storage{
		Users:   users.NewStorage(usersP),
		Devices: devices.NewStorage(devicesP),
	}
}
