package devices

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"strings"
	"time"
)

type Device struct {
	ID       string       `json:"id"`
	Name     string       `json:"name"`
	OwnerID  uint         `json:"owner_id"`
	LinkedAt time.Time    `json:"linked_at"`
	ApiKey   DeviceApiKey `json:"api_key"`
}

func (d Device) Phone() string {
	s := strings.SplitN(d.ID, ":", 2)
	if len(s) == 2 {
		return s[0]
	}
	return d.ID
}

type DeviceApiKey struct {
	ID        uint      `json:"id"`
	DeviceID  string    `json:"device_id"`
	Name      string    `json:"name"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
}

func GenerateAPIKey(length int) (key string, err error) {
	b := make([]byte, length)
	_, err = rand.Read(b)
	if err != err {
		return
	}
	out := bytes.Buffer{}
	_, err = base64.NewEncoder(base64.StdEncoding, &out).Write(b)
	if err != nil {
		return
	}
	return out.String(), nil
}
