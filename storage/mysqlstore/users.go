package mysqlstore

import (
	"database/sql"
	"errors"
	"time"

	"github.com/9d4/wadoh/users"
	"github.com/rs/zerolog/log"
)

type usersStore struct {
	db *sql.DB
}

func newUsersStore(db *sql.DB) *usersStore {
	return &usersStore{db: db}
}

func (s *usersStore) GetByID(id uint) (*users.User, error) {
	row := s.db.QueryRow(`SELECT * FROM wadoh_users WHERE id=?`, id)
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
	const query = "SELECT u.id, u.name, u.username, u.password, u.created_at," +
		"p.admin, p.updated_at " +
		"FROM wadoh_users AS u " +
		"LEFT JOIN wadoh_user_perms AS p ON p.user_id = u.id " +
		"WHERE id > ? LIMIT ?"

	log.Printf("%q", query)
	rows, err := s.db.Query(query, after, limit)
	if err != nil {
		return nil, err
	}

	var usrs []users.User

	defer rows.Close()
	for rows.Next() {
		var u users.User

		var permAdmin *bool
		var permUpdatedAt *time.Time

		if err := rows.Scan(
			&u.ID, &u.Name, &u.Username, &u.Password, &u.CreatedAt,
			&permAdmin, &permUpdatedAt,
		); err != nil {
			return nil, err
		}

		if permAdmin != nil {
			u.Perm.Admin = *permAdmin
		}
		if permUpdatedAt != nil {
			u.Perm.UpdatedAt = *permUpdatedAt
		}
		usrs = append(usrs, u)
	}
	return usrs, nil
}

func (s *usersStore) Save(u *users.User) error {
	const query = `INSERT INTO wadoh_users (name, username, password, created_at) VALUES (?, ?, ?, ?)`

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	result, err := tx.Exec(query, u.Name, u.Username, u.Password, u.CreatedAt)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Get the last inserted ID and assign it to the user
	id, err := result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}

	const permQ = "INSERT INTO wadoh_user_perms (user_id, admin, updated_at)" +
		"VALUES (?, ?, ?)"
	_, err = tx.Exec(permQ, id, u.Perm.Admin, u.Perm.UpdatedAt)
	if err != nil {
		tx.Rollback()
		return err
	}

	u.ID = uint(id)
	return tx.Commit()
}

func (s *usersStore) Update(u *users.User) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	const userQ = "UPDATE wadoh_users " +
		"SET name=?,username=?,password=? " +
		"WHERE id=?"
	_, err = tx.Exec(userQ, u.Name, u.Username, u.Password, u.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	const permQ = "UPDATE wadoh_user_perms " +
		"SET admin=?,updated_at=? " +
		"WHERE user_id=?"
	_, err = tx.Exec(permQ, u.Perm.Admin, u.Perm.UpdatedAt, u.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (s *usersStore) Delete(id uint) error {
    tx, err := s.db.Begin()
    if err != nil {
        return err
    }

    const deviceQ = "DELETE FROM wadoh_devices WHERE user_id=?"
    _, err = tx.Exec(deviceQ, id)
    if err != nil {
        tx.Rollback()
        return err
    }

    const userQ = "DELETE FROM wadoh_users WHERE id=?"
    _, err = tx.Exec(userQ, id)
    if err != nil {
        tx.Rollback()
        return err
    }
    return tx.Commit()
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
