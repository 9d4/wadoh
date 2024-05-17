package mysqlstore

import (
	"database/sql"

	"github.com/9d4/wadoh/users"
)

type usersStore struct {
	db *sql.DB
}

func newUsersStore(db *sql.DB) *usersStore {
	return &usersStore{db: db}
}

func (s *usersStore) GetByID(uint) (*users.User, error) {
	panic("not implemented")
}

func (s *usersStore) Page(limit int, after int) ([]users.User, error) {
	const query = `SELECT * FROM wadoh_users WHERE id > ? LIMIT ?`

	rows, err := s.db.Query(query, after, limit)
	if err != nil {
		return nil, err
	}

	var usrs []users.User

	defer rows.Close()
	for rows.Next() {
		var u users.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Username, &u.Password, &u.CreatedAt); err != nil {
			return nil, err
		}
		usrs = append(usrs, u)
	}
	return usrs, nil
}

func (s *usersStore) Save(u *users.User) error {
	const query = `INSERT INTO wadoh_users (name, username, password, created_at) VALUES (?, ?, ?, ?)`

	result, err := s.db.Exec(query, u.Name, u.Username, u.Password, u.CreatedAt)
	if err != nil {
		return err
	}

	// Get the last inserted ID and assign it to the user
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	u.ID = uint(id)
	return nil
}

var _ users.StorageProvider = (*usersStore)(nil)
