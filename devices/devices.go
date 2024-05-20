package devices

import "time"

type Device struct {
	ID       string    `json:"id"`
	Name     string    `json:"name"`
	OwnerID  uint      `json:"owner_id"`
	LinkedAt time.Time `json:"linked_at"`
}
