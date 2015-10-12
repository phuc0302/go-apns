package apns

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestCreateMessage(t *testing.T) {
	message := CreateMessage("", "")
	if message != nil {
		t.Errorf("Expect nil but found %s when asign invalid token", message)
	}

	message = CreateMessage("YWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYW", "")
	if message != nil {
		t.Errorf("Expect nil but found %s when asign not base64 token", message)
	}

	message = CreateMessage("YWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWFhYWE=", "")
	if message.OsVersion != "7.0" {
		t.Errorf("Expect 7.0 but found %s when asign empty os version value", message.OsVersion)
	}
}

func TestLength(t *testing.T) {
	message := CreateMessage("YWJjZGVmZ2hpamtsbW5vcHFyc3R1dnd4eXpBQkNERUY=", "")
	payload, _ := json.Marshal(message.payload)

	if message.Length() != uint32(61+len(payload)) {
		t.Errorf("Expect message length is %d but found %d", 61+len(payload), message.Length())
	}

	// Encode message with payload
	message.SetPayload(&Payload{
		Alert: "Alert message",
		Badge: 0,
		Sound: "Default",
	})
	reader, _ := message.Encode()

	bytes := make([]byte, 512)
	length, _ := reader.Read(bytes)

	if uint32(length) != message.Length() {
		t.Errorf("Expect message length is %d but found %d", length, message.Length())
	}
}

func TestEncodeWithoutPayload(t *testing.T) {
	message := CreateMessage("YWJjZGVmZ2hpamtsbW5vcHFyc3R1dnd4eXpBQkNERUY=", "")
	reader, _ := message.Encode()

	if reader != nil {
		t.Error("Expect nil when encode no payload message but found not nil")
	}

	// Asign nil also do the same thing
	message.SetPayload(nil)
	reader, _ = message.Encode()

	if reader != nil {
		t.Error("Expect nil when encode no payload message but found not nil")
	}

	// Asign empty payload also do the same thing
	message.SetPayload(&Payload{})
	reader, _ = message.Encode()

	if reader != nil {
		t.Error("Expect nil when encode no payload message but found not nil")
	}
}

func TestEncodeWithDefaultPayload(t *testing.T) {
	message := CreateMessage("YWJjZGVmZ2hpamtsbW5vcHFyc3R1dnd4eXpBQkNERUY=", "")
	message.SetPayload(&Payload{
		Alert: "Alert message",
	})
	reader, _ := message.Encode()

	if reader == nil {
		t.Error("Expect not nil when encode but found nil")
	} else {
		payload := message.payload["aps"].(*Payload)

		if payload == nil {
			t.Error("Expect payload not nil but found nil")
		}

		if payload.Alert != "Alert message" {
			t.Errorf("Expect %s but found %s", "Alert message", payload.Alert)
		}

		if payload.Badge != 1 {
			t.Errorf("Expect %d but found %d", 1, payload.Badge)
		}

		if payload.Sound != "Default" {
			t.Errorf("Expect %s but found %s", "Default", payload.Sound)
		}
	}
}

func TestEncodeExceedLengthForOS7(t *testing.T) {
	message := CreateMessage("YWJjZGVmZ2hpamtsbW5vcHFyc3R1dnd4eXpBQkNERUY=", "")
	message.SetPayload(&Payload{
		Alert: "Alert message",
		Badge: 0,
		Sound: "Default",
	})
	for i := 0; i < 100; i++ {
		message.SetField(fmt.Sprintf("Key %d", i), "Value")
	}
	reader, _ := message.Encode()

	if reader != nil {
		t.Error("Expect nil when encode no exceeded message but found not nil")
	}
}
