package devices

import (
	"context"
	"net/url"
	"time"

	"github.com/9d4/wadoh/internal"
	"github.com/9d4/wadoh/wadoh-be/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	APIKeyLength = 32
)

type StorageProvider interface {
	ListByOwnerID(uint) ([]Device, error)
	ListAll(uint) ([]Device, error)
	GetByID(string) (*Device, error)
	Save(*Device) error
	Patch(*Device) error
	Delete(string) error
	SaveAPIKey(*DeviceApiKey) error
	GetByAPIToken(string) (*Device, error)
	SaveWebhook(*DeviceWebhook) error // Deprecated
	ChangeID(id, newID string) error
}

type Storage struct {
	provider StorageProvider
	pbCli    pb.ControllerServiceClient
}

func NewStorage(provider StorageProvider, pbCli pb.ControllerServiceClient) *Storage {
	return &Storage{
		provider: provider,
		pbCli:    pbCli,
	}
}

func (s *Storage) Save(d *Device) error {
	return s.provider.Save(d)
}

func (s *Storage) ListByOwnerID(ownerID uint) ([]Device, error) {
	return s.provider.ListByOwnerID(ownerID)
}

func (s *Storage) ListAll(ownerID uint) ([]Device, error) {
	return s.provider.ListAll(ownerID)
}

func (s *Storage) GetByID(id string) (*Device, error) {
	device, err := s.provider.GetByID(id)
	if err != nil {
		return nil, parseError(err, id)
	}
	device.Webhook = &DeviceWebhook{
		DeviceID: device.ID,
	}

	res, err := s.pbCli.GetWebhook(context.Background(), &pb.GetWebhookRequest{
		Jid: device.ID,
	})
	if err == nil {
		device.Webhook.URL = res.GetUrl()
	} else if status.Code(err) != codes.NotFound {
		return nil, parseError(err, id)
	}

	return device, nil
}

// On reconnect we need to change the device id to the new one
func (s *Storage) ChangeID(id, newID string) error {
	return s.provider.ChangeID(id, newID)
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

func (s *Storage) SaveWebhook(wh *DeviceWebhook) error {
	if wh.URL != "" {
		url, err := url.ParseRequestURI(wh.URL)
		if err != nil {
			newErr := internal.NewError(internal.EBADINPUT, "Unable to parse url, please check before try again", "url.parse_err")
			return wrapError(err, newErr, wh.URL)
		}
		wh.URL = url.String()

		_, err = s.pbCli.SaveWebhook(context.Background(), &pb.SaveWebhookRequest{
			Jid: wh.DeviceID,
			Url: wh.URL,
		})
		if err != nil {
			return parseError(err, wh.DeviceID)
		}
		return nil
	}

	_, err := s.pbCli.DeleteWebhook(context.Background(), &pb.DeleteWebhookRequest{
		Jid: wh.DeviceID,
	})
	if err != nil {
		return parseError(err, wh.DeviceID)
	}
	return nil
}
