package devices

type StorageProvider interface {
	ListByOwnerID(uint) ([]Device, error)
	Save(*Device) error
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
