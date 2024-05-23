package devices

import "time"

type Device struct {
	ID       string       `json:"id"`
	Name     string       `json:"name"`
	OwnerID  uint         `json:"owner_id"`
	LinkedAt time.Time    `json:"linked_at"`
	ApiKey   DeviceApiKey `json:"api_key"`
}

type DeviceApiKey struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
}
