package mysqlstore

import (
	"database/sql"

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

func (s *devicesStore) Save(d *devices.Device) error {
	const query = `INSERT INTO wadoh_devices (id, name, user_id, linked_at) VALUES (?, ?, ?, ?)`

	result, err := s.db.Exec(query, d.ID, d.Name, d.OwnerID, d.LinkedAt)
	if err != nil {
		return err
	}
	_, err = result.RowsAffected()
	if err != nil {
		return err
	}

	return nil
}
