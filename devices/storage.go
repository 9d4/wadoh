package devices

type StorageProvider interface {
}

type Storage struct {
	provider StorageProvider
}

func NewStorage(provider StorageProvider) *Storage {
	return &Storage{
		provider: provider,
	}
}
