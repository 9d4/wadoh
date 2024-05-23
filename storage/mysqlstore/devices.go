package mysqlstore

import (
	"database/sql"
	"strings"
	"time"

	"github.com/9d4/wadoh/devices"
)

type devicesStore struct {
	db *sql.DB
}

func newDevicesStore(db *sql.DB) *devicesStore {
	return &devicesStore{db: db}
}

func (s *devicesStore) ListByOwnerID(ownerID uint) ([]devices.Device, error) {
	const query = `SELECT id, name, user_id, linked_at FROM wadoh_devices WHERE user_id=?`

	rows, err := s.db.Query(query, ownerID)
	if err != nil {
		return nil, err
	}

	var devs []devices.Device
	for rows.Next() {
		var d devices.Device
		if err := rows.Scan(&d.ID, &d.Name, &d.OwnerID, &d.LinkedAt); err != nil {
			return nil, err
		}
		devs = append(devs, d)
	}

	return devs, nil
}

func (s *devicesStore) GetByID(id string) (*devices.Device, error) {
	const query = "SELECT devices.id, devices.name, devices.user_id, devices.linked_at, " +
		"`keys`.id, `keys`.name, `keys`.token, `keys`.created_at FROM wadoh_devices devices " +
		"LEFT JOIN wadoh_device_api_keys `keys` ON `keys`.jid=devices.id " +
		"WHERE devices.id=?"

	row := s.db.QueryRow(query, id)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var dev devices.Device
	var keyID *uint
	var keyName, keyToken *string
	var keyCreatedAt *time.Time

	if err := row.Scan(&dev.ID, &dev.Name, &dev.OwnerID, &dev.LinkedAt,
		&keyID, &keyName, &keyToken, &keyCreatedAt,
	); err != nil {
		return nil, err
	}

    if keyID != nil {
        dev.ApiKey.ID = *keyID
    }
    if keyName != nil {
        dev.ApiKey.Name = *keyName
    }
    if keyToken != nil {
        dev.ApiKey.Token = *keyToken
    }
    if keyCreatedAt != nil {
        dev.ApiKey.CreatedAt = *keyCreatedAt
    }

	return &dev, nil
}

func (s *devicesStore) Save(d *devices.Device) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	const query = `INSERT INTO wadoh_devices (id, name, user_id, linked_at) VALUES (?, ?, ?, ?)`

	result, err := tx.Exec(query, d.ID, d.Name, d.OwnerID, d.LinkedAt)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}

	result, err = tx.Exec(`INSERT INTO wadoh_device_api_keys (jid, name, token, created_at) VALUES (?, ?, ?, ?)`)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = result.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (s *devicesStore) Patch(d *devices.Device) error {
	query := `UPDATE wadoh_devices SET $cols WHERE id=?`
	cols := ""
	args := []any{}

	i := 0
	if d.Name != "" {
		cols += "name=?"
		args = append(args, d.Name)
		i++
	}
	if d.OwnerID != 0 {
		if i != 0 {
			cols += ",user_id=?"
		} else {
			cols += "user_id=?"
		}
		args = append(args, d.OwnerID)
	}
	// append the id
	args = append(args, d.ID)
	query = strings.Replace(query, "$cols", cols, 1)

	_, err := s.db.Exec(query, args...)
	if err != nil {
		return err
	}
	return nil
}
