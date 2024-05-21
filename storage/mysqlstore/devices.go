package mysqlstore

import (
	"database/sql"
	"strings"

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
	const query = `SELECT id, name, user_id, linked_at FROM wadoh_devices WHERE id=?`

	row := s.db.QueryRow(query, id)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var dev devices.Device
	if err := row.Scan(&dev.ID, &dev.Name, &dev.OwnerID, &dev.LinkedAt); err != nil {
		return nil, err
	}
	return &dev, nil
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
