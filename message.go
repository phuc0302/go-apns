package apns

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

const (
	iOS7_MaxPayload    = 255
	iOS8_MaxPayload    = 2048
	itemId_DeviceToken = 1
	itemId_Payload     = 2
	itemId_Identifier  = 3
	itemId_Expired     = 4
	itemId_Priority    = 5
	length_DeviceToken = 32
	length_Identifier  = 4
	length_Expired     = 4
	length_Priority    = 1
)

type Alert struct {
	Body         string   `json:"body,omitempty"`
	LaunchImage  string   `json:"launch-image,omitempty"`
	ActionLocKey string   `json:"action-loc-key,omitempty"`
	LocKey       string   `json:"loc-key,omitempty"`
	LocArgs      []string `json:"loc-args,omitempty"`
}

type Payload struct {
	Alert    interface{} `json:"alert,omitempty"`
	Badge    uint        `json:"badge,omitempty"`
	Sound    string      `json:"sound,omitempty"`
	Category string      `json:"category,omitempty"`

	ContentAvailable int `json:"content-available,omitempty"`
}

//////////////////////////////////////////////////////////////////////////////
//-- Message ---------------------------------------------------------------//
type Message struct {
	Id          int32
	Expired     uint32
	OsVersion   string
	DeviceToken []byte

	priority uint8
	payload  map[string]interface{}
}

// MARK: Struct's constructors
func CreateMessage(deviceToken string, osVersion string) *Message {
	/* Condition validation: Validate decoded device's token */
	token, err := base64.StdEncoding.DecodeString(deviceToken)
	if err != nil || len(token) != length_DeviceToken {
		return nil
	}

	/* Condition validation: Revert to classic size */
	if len(osVersion) == 0 {
		osVersion = "7.0"
	}

	apns := &Message{
		Id:          rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(9999),
		DeviceToken: token,
		OsVersion:   osVersion,
		priority:    10,
		payload:     make(map[string]interface{}),
	}
	return apns
}

/**
 * Calculate length of the message.
 */
func (m *Message) Length() uint32 {
	length := 61

	payload, err := json.Marshal(m.payload)
	if err == nil {
		length += len(payload)
	}
	return uint32(length)
}

/**
 * Encode the message before send to Apple.
 */
func (m *Message) Encode() (io.Reader, error) {
	/* Condition validation: Ignore message without payload */
	if m.payload["aps"] == nil {
		return nil, errors.New("Payload must not empty.")
	}

	// Encode payload
	payload, err := json.Marshal(m.payload)
	if err != nil {
		return nil, err
	}

	/* Condition validation: Validate payload's token */
	strings := strings.Split(m.OsVersion, ".")
	version, err := strconv.Atoi(strings[0])

	if err == nil && version >= 8 && len(payload) > iOS8_MaxPayload {
		return nil, errors.New(fmt.Sprintf("Payload is larger than: %i bytes.", iOS8_MaxPayload))
	} else if len(payload) > iOS7_MaxPayload {
		return nil, errors.New(fmt.Sprintf("Payload is larger than: %i bytes.", iOS7_MaxPayload))
	}

	/** Encode message. */
	buffer := bytes.NewBuffer(nil)

	// Write message header
	binary.Write(buffer, binary.BigEndian, uint8(2))
	binary.Write(buffer, binary.BigEndian, uint32(56+len(payload))) // Other content beside payload will take up 56 bytes

	// Write device token
	binary.Write(buffer, binary.BigEndian, uint8(itemId_DeviceToken))
	binary.Write(buffer, binary.BigEndian, uint16(len(m.DeviceToken)))
	binary.Write(buffer, binary.BigEndian, m.DeviceToken)
	// Write payload
	binary.Write(buffer, binary.BigEndian, uint8(itemId_Payload))
	binary.Write(buffer, binary.BigEndian, uint16(len(payload)))
	binary.Write(buffer, binary.BigEndian, payload)
	// Write identifier
	binary.Write(buffer, binary.BigEndian, uint8(itemId_Identifier))
	binary.Write(buffer, binary.BigEndian, uint16(length_Identifier))
	binary.Write(buffer, binary.BigEndian, m.Id)
	// Write expire time
	binary.Write(buffer, binary.BigEndian, uint8(itemId_Expired))
	binary.Write(buffer, binary.BigEndian, uint16(length_Expired))
	binary.Write(buffer, binary.BigEndian, m.Expired)
	// Write priority
	binary.Write(buffer, binary.BigEndian, uint8(itemId_Priority))
	binary.Write(buffer, binary.BigEndian, uint16(length_Priority))
	binary.Write(buffer, binary.BigEndian, m.priority)

	return buffer, nil
}

/**
 * Add payload to message.
 */
func (a *Message) SetPayload(p *Payload) {
	/* Condition validation */
	if p == nil || p.Alert == nil {
		return
	}

	if p.Badge == 0 {
		p.Badge = 1
	}

	if len(p.Sound) == 0 {
		p.Sound = "Default"
	}

	a.payload["aps"] = p
}

/**
 * Add custom key-value pair to the message's payload.
 */
func (a *Message) SetField(key string, value interface{}) {
	/* Condition validation */
	if len(key) == 0 || key == "aps" {
		return
	}
	a.payload[key] = value
}
