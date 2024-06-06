package devices

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/9d4/wadoh/internal"
)

var (
	ErrDeviceNotFound = internal.NewError(internal.ENOTFOUND, "Device not found", "device.not_found")
)

func parseError(err error, data interface{}) error {
	if errors.Is(err, sql.ErrNoRows) {
		return wrapError(err, ErrDeviceNotFound, data)
	}
	return wrapError(err, nil, data)
}

func wrapError(parent error, err error, data interface{}) error {
	if err == nil {
		return fmt.Errorf("devices: %v: %w", data, parent)
	}
	return fmt.Errorf("devices: %v %w: %w", data, err, parent)
}
