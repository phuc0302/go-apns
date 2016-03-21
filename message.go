package apns

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/satori/go.uuid"
)

const (
	// PriorityHigh indicates that the push message must be sent immediately.
	PriorityHigh = 10
	// PriorityLow indicates that the push message will be sent at a time that takes into account power considerations for the device.
	PriorityLow = 5
	// PayloadSize indicates the maximum size allowed for a notification payload when using the HTTP/2.
	PayloadSize = 4096
)

// Alert describes an alert message to be displayed when sending a push notification.
type Alert struct {
	Body         string   `json:"body,omitempty"`           // The text of the alert message.
	LaunchImage  string   `json:"launch-image,omitempty"`   // The image is used as the launch image when users tap the action button or move the action slider.
	LocKey       string   `json:"loc-key,omitempty"`        // The key to an alert-message string in a Localizable.strings. The key string can be formatted to take the variables specified in the loc-args.
	LocArgs      []string `json:"loc-args,omitempty"`       // Variable string values to appear in place of the format specifiers in loc-key.
	ActionLocKey string   `json:"action-loc-key,omitempty"` // If a string is specified, the system displays an alert that includes the Close and View buttons. The string is used as a key to get a localized string in the current localization to use for the right button’s title instead of “View”.

	// Available from iOS 8.2
	Title        string   `json:"title,omitempty"`          // A short string describing the purpose of the notification. Apple Watch displays this string as part of the notification interface.
	TitleLocKey  string   `json:"title-loc-key,omitempty"`  // The key to a title string in the Localizable.strings. The key string can be formatted to take the variables specified in the title-loc-args.
	TitleLocArgs []string `json:"title-loc-args,omitempty"` // Variable string values to appear in place of the format specifiers in title-loc-key.
}

// Payload describes an aps dictionary within a push notification message.
type Payload struct {
	Alert interface{} `json:"alert,omitempty"` // Either string or Alert struct.
	Badge uint        `json:"badge,omitempty"`
	Sound string      `json:"sound,omitempty"`

	Category         string `json:"category,omitempty"`
	ContentAvailable uint   `json:"content-available,omitempty"`
}

// Message describes an apns2 message.
type Message struct {
	ApnsID         string // A canonical UUID that identifies the notification.
	ApnsTopic      string // The topic of the remote notification, which is typically the bundleID for your app.
	ApnsPriority   int8   // High or low
	ApnsExpiration int64  // A UNIX epoch date expressed in seconds (UTC).

	deviceID    string
	deviceToken []byte
	payload     map[string]interface{}
}

// CreateMessage returns a default message
func CreateMessage(deviceID string, deviceToken string, apnsTopic string, payload Payload) (*Message, error) {
	/* Condition validation: validate deviceToken */
	if len(deviceToken) == 0 {
		return nil, fmt.Errorf("MissingDeviceToken")
	}

	/* Condition validation: validate apnsTopic */
	if len(apnsTopic) == 0 {
		return nil, fmt.Errorf("MissingTopic")
	}

	// /* Condition validation: validate payload */
	// alertString, ok := payload.Alert.(string)
	// if ok && len(alertString) == 0 {
	// 	return nil
	// } else if !ok {
	// 	_, ok := payload.Alert.(Alert)
	// 	if !ok {
	// 		return nil
	// 	}
	// }

	if payload.Badge == 0 {
		payload.Badge = 1
	}
	if len(payload.Sound) == 0 {
		payload.Sound = "Default"
	}

	/* Condition validation: Validate decoded device's token */
	token, _ := base64.StdEncoding.DecodeString(deviceToken)
	if len(token) != 32 {
		return nil, fmt.Errorf("BadDeviceToken")
	}

	// Finalize
	message := Message{
		ApnsID:         uuid.NewV4().String(),
		ApnsTopic:      apnsTopic,
		ApnsPriority:   PriorityHigh,
		ApnsExpiration: time.Now().UTC().Unix() + 86400,

		deviceID:    deviceID,
		deviceToken: token,
		payload:     make(map[string]interface{}),
	}

	message.payload["aps"] = payload
	return &message, nil
}

// Encode prepares the HTTP/2 request to send to Apple.
func (m *Message) Encode(gateway string) (*http.Request, error) {
	/* Condition validation: validate gateway */
	if len(gateway) == 0 {
		return nil, fmt.Errorf("ServiceUnavailable")
	}

	// Encode payload
	payload, err := json.Marshal(m.payload)
	if err != nil {
		return nil, fmt.Errorf("PayloadEmpty")
	} else if len(payload) > PayloadSize {
		return nil, fmt.Errorf("PayloadTooLarge")
	}

	// Prepare request
	urlString := fmt.Sprintf("https://%s/3/device/%x", gateway, m.deviceToken)
	request, _ := http.NewRequest("POST", urlString, bytes.NewBuffer(payload))
	request.Header.Set("Content-Type", "application/json; charset=utf-8")
	request.Header.Set("apns-id", m.ApnsID)
	request.Header.Set("apns-topic", m.ApnsTopic)
	request.Header.Set("apns-priority", fmt.Sprintf("%d", m.ApnsPriority))
	request.Header.Set("apns-expiration", fmt.Sprintf("%d", m.ApnsExpiration))

	return request, nil
}

// SetField adds custom key-value pair to the message's payload.
func (m *Message) SetField(key string, value interface{}) {
	/* Condition validation */
	if len(key) == 0 || key == "aps" {
		return
	}
	m.payload[key] = value
}
