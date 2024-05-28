package mysqlstore

import (
	"database/sql"
	"errors"

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

func (s *usersStore) GetByUsername(username string) (*users.User, error) {
	row := s.db.QueryRow(`SELECT * FROM wadoh_users WHERE username=?`, username)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var u users.User
	if err := row.Scan(&u.ID, &u.Name, &u.Username, &u.Password, &u.CreatedAt); err != nil {
		return nil, err
	}

	perm, err := s.firstOrInitPermission(u.ID)
	if err != nil {
		return nil, err
	}
	u.Perm = *perm

	return &u, nil
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

func (s *usersStore) firstOrInitPermission(userID uint) (*users.Permissions, error) {
	row := s.db.QueryRow(
		`SELECT admin, updated_at FROM wadoh_user_perms WHERE user_id=?`,
		userID,
	)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var perm users.Permissions
	if err := row.Scan(&perm.Admin, &perm.UpdatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_, err := s.db.Exec(`INSERT INTO wadoh_user_perms(user_id) VALUES (?)`, userID)
			if err != nil {
				return nil, err
			}
			return &perm, nil
		}
		return nil, err
	}
	return &perm, nil
}

var _ users.StorageProvider = (*usersStore)(nil)
