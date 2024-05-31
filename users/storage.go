package users

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
)

type StorageProvider interface {
	Page(limit, after int) ([]User, error)
	GetByID(uint) (*User, error)
	GetByUsername(string) (*User, error)
	Save(u *User) error
	Update(u *User) error
	Delete(id uint) error
}

type Storage struct {
	provider StorageProvider
}

func NewStorage(provider StorageProvider) *Storage {
	return &Storage{
		provider: provider,
	}
}

func (s *Storage) List(limit, after int) ([]User, error) {
	if limit <= 0 {
		limit = 10
	}

	return s.provider.Page(limit, after)
}

func (s *Storage) GetBy(v interface{}) (*User, error) {
	switch reflect.TypeOf(v).String() {
	case "string":
		return s.provider.GetByUsername(v.(string))
	case "int":
		return s.provider.GetByID(uint(v.(int)))
	case "uint":
		return s.provider.GetByID(v.(uint))
	}
	return nil, errors.New("invalid lookup type")
}

func (s *Storage) Save(u *User) (err error) {
	u.ID = 0
	u.Name = strings.TrimSpace(u.Name)
	u.Username = strings.TrimSpace(u.Username)
	u.Password = strings.TrimSpace(u.Password)
	u.CreatedAt = time.Now()
	u.Perm.UpdatedAt = time.Now()

	if u.Username == "" {
		return errors.New("username should not be empty")
	}

	u.Password, err = hashPwd(u.Password)
	if err != nil {
		return fmt.Errorf("unable to hash password: %w", err)
	}

	if err := s.provider.Save(u); err != nil {
		return err
	}

	return nil
}

func (s *Storage) Update(u *User) (err error) {
	oldUser, err := s.GetBy(u.ID)
	if err != nil {
		return err
	}

	u.Name = strings.TrimSpace(u.Name)
	u.Username = strings.TrimSpace(u.Username)
	u.Password = strings.TrimSpace(u.Password)
	u.Perm.UpdatedAt = time.Now()

	if u.Username == "" {
		return errors.New("username should not be empty")
	}

	// Update password only if not empty
	if u.Password != "" {
		u.Password, err = hashPwd(u.Password)
		if err != nil {
			return fmt.Errorf("unable to hash password: %w", err)
		}
	} else {
		u.Password = oldUser.Password
	}

	if err := s.provider.Update(u); err != nil {
		log.Err(err).Send()
		return err
	}
	return nil
}

func (s *Storage) Delete(id uint) error {
	return s.provider.Delete(id)
}

func hashPwd(pwd string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func ComparePwd(hash, pwd string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd))
}
