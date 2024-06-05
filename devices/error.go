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
	return err
}

func wrapError(parent error, err error, data interface{}) error {
    return fmt.Errorf("devices: %v %w: %w", data, err, parent)
}
