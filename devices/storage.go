package devices

import "time"

const (
	APIKeyLength = 32
)

type StorageProvider interface {
	ListByOwnerID(uint) ([]Device, error)
	GetByID(string) (*Device, error)
	Save(*Device) error
	Patch(*Device) error
	Delete(string) error
	SaveAPIKey(*DeviceApiKey) error
	GetByAPIToken(string) (*Device, error)
}

type Storage struct {
	provider StorageProvider
}

func NewStorage(provider StorageProvider) *Storage {
	return &Storage{
		provider: provider,
	}
}

func (s *Storage) Save(d *Device) error {
	return s.provider.Save(d)
}

func (s *Storage) ListByOwnerID(ownerID uint) ([]Device, error) {
	return s.provider.ListByOwnerID(ownerID)
}

func (s *Storage) GetByID(id string) (*Device, error) {
	return s.provider.GetByID(id)
}

func (s *Storage) Rename(id, newName string) error {
	return s.provider.Patch(&Device{
		ID:   id,
		Name: newName,
	})
}

func (s *Storage) Delete(id string) error {
	return s.provider.Delete(id)
}

func (s *Storage) GenNewDevAPIKey(deviceID string) error {
	token, err := GenerateAPIKey(APIKeyLength)
	if err != nil {
		return err
	}
	k := DeviceApiKey{
		DeviceID:  deviceID,
		Token:     token,
		CreatedAt: time.Now(),
	}

	if err := s.provider.SaveAPIKey(&k); err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetByAPIToken(token string) (*Device, error) {
	return s.provider.GetByAPIToken(token)
}
