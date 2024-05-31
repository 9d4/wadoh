package users

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type StorageProvider interface {
	Page(limit, after int) ([]User, error)
	GetByID(uint) (*User, error)
	GetByUsername(string) (*User, error)
	Save(u *User) error
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
