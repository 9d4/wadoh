package mysqlstore

import (
	"context"
	"database/sql"

	"github.com/rs/zerolog/log"
)

var versions = []func(*sql.Tx) error{upgradeV1, upgradeV2, upgradeV3}

func Upgrade(db *sql.DB) error {
	version, err := getVersion(db)
	if err != nil {
		return err
	}

	for ; version < len(versions); version++ {
		log.Info().Msgf("%d", version)
		tx, err := db.BeginTx(context.Background(), nil)
		if err != nil {
			return err
		}

		upgradeFn := versions[version]
		log.Info().Msgf("upgrading database to version %d", version+1)
		err = upgradeFn(tx)
		if err != nil {
			tx.Rollback()
			return err
		}

		if err := setVersion(tx, version+1); err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.Commit(); err != nil {
			tx.Rollback()
			return err
		}
	}

	return nil
}

func getVersion(db *sql.DB) (int, error) {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS wadoh_version (version int)")
	if err != nil {
		return 0, err
	}

	version := 0
	row := db.QueryRow("SELECT version FROM wadoh_version LIMIT 1")
	if row.Err() == nil {
		row.Scan(&version)
	}
	return version, nil
}

func setVersion(tx *sql.Tx, v int) error {
	_, err := tx.Exec("DELETE FROM wadoh_version")
	if err != nil {
		return err
	}

	_, err = tx.Exec("INSERT INTO wadoh_version (version) VALUES (?)", v)
	if err != nil {
		return err
	}
	return nil
}

func upgradeV1(tx *sql.Tx) error {
	_, err := tx.Exec(`CREATE TABLE wadoh_users (
        id INT AUTO_INCREMENT PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
		username VARCHAR(255) NOT NULL,
		password VARCHAR(255) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		UNIQUE idx_username (username)
    )`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`CREATE TABLE wadoh_refresh_tokens (
        id INT AUTO_INCREMENT PRIMARY KEY,
        user_id INT NOT NULL,
		token VARCHAR(255) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES wadoh_users(id) ON DELETE CASCADE
    )`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`CREATE TABLE wadoh_devices (
		id VARCHAR(255) PRIMARY KEY,
        name VARCHAR(255) NOT NULL,
        user_id INT NOT NULL,
        linked_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (user_id) REFERENCES wadoh_users(id) ON DELETE RESTRICT
	)`)
	if err != nil {
		return err
	}

	return nil
}

func upgradeV2(tx *sql.Tx) error {
	_, err := tx.Exec(`CREATE TABLE wadoh_device_api_keys (
	        id INT AUTO_INCREMENT PRIMARY KEY,
		jid VARCHAR(255) NOT NULL,
		name VARCHAR(255) NOT NULL,
		token VARCHAR(255) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (jid) REFERENCES wadoh_devices(id) ON DELETE CASCADE,
		UNIQUE idx_token (token)
    )`)
	if err != nil {
		return err
	}

	return nil
}

func upgradeV3(tx *sql.Tx) error {
	_, err := tx.Exec(`CREATE TABLE wadoh_user_perms (
        user_id INT NOT NULL,
        admin TINYINT(1) NOT NULL DEFAULT 0,
		updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES wadoh_users(id) ON DELETE CASCADE
    )`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`INSERT INTO wadoh_user_perms (user_id)
        SELECT wadoh_users.id
        FROM wadoh_users
        WHERE NOT EXISTS (
            SELECT 1
            FROM wadoh_user_perms
            WHERE wadoh_user_perms.user_id = wadoh_users.id
    );`)
	if err != nil {
		return err
	}

	return nil
}
