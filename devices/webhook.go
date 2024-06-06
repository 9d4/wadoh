package devices

const (
	webhookMethod = "POST"
)

type DeviceWebhook struct {
	DeviceID string
	URL      string
}

type WebhookData struct {
	From    string `json:"from"`
	Message string `json:"message"`
}
