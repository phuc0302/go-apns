package apns

import (
	"encoding/base64"
	"fmt"
	"strings"
	"testing"
	"time"

	"gopkg.in/mgo.v2/bson"
)

func TestCreateMessage(t *testing.T) {
	message := CreateMessage("", "", "", Payload{})
	if message != nil {
		t.Error("Expected nil when define invalid device's id.")
	}

	message = CreateMessage(bson.NewObjectId().Hex(), "", "", Payload{})
	if message != nil {
		t.Error("Expected nil when define invalid device's token.")
	}

	message = CreateMessage(bson.NewObjectId().Hex(), "C4XOCR6kmbH4XJ9fMRm1hyt1iL7f0wqfJENdgTDdx+A=", "", Payload{})
	if message != nil {
		t.Error("Expected nil when define invalid message's topic.")
	}

	message = CreateMessage(bson.NewObjectId().Hex(), "C4XOCR6kmbH4XJ9fMRm1hyt1iL7f0wqfJENdgTDdx+A=", "com.example.appID", Payload{})
	if message != nil {
		t.Error("Expected nil when define invalid message's payload.")
	}

	deviceID := bson.NewObjectId().Hex()
	message = CreateMessage(deviceID, "C4XOCR6kmbH4XJ9fMRm1hyt1iL7f0wqfJENdgTDdx+A=", "com.example.appID",
		Payload{
			Alert: "Sample alert",
		})

	if message == nil {
		t.Error("Expected not nil when everything had been defined.")
	}
	if message.ApnsID == "" {
		t.Error("Expected ApnsID not nil.")
	}
	if message.ApnsTopic != "com.example.appID" {
		t.Errorf("Expected ApnsTopic is: %s but found: %s", "com.example.appID", message.ApnsTopic)
	}
	if message.ApnsPriority != PriorityHigh {
		t.Errorf("Expected ApnsPriority is: %d but found: %d", PriorityHigh, message.ApnsPriority)
	}
	if message.ApnsExpiration <= time.Now().UTC().Unix() {
		t.Errorf("Expected ApnsExpiration is greater than: %d but found: %d", time.Now().UTC().Unix(), message.ApnsExpiration)
	}

	if message.deviceID != deviceID {
		t.Errorf("Expected deviceID: %s but found: %s.", deviceID, message.deviceID)
	}
	if base64.StdEncoding.EncodeToString(message.deviceToken) != "C4XOCR6kmbH4XJ9fMRm1hyt1iL7f0wqfJENdgTDdx+A=" {
		t.Errorf("Expected deviceToken: %s but found: %s.", "C4XOCR6kmbH4XJ9fMRm1hyt1iL7f0wqfJENdgTDdx+A=", base64.StdEncoding.EncodeToString(message.deviceToken))
	}

	if message.payload["aps"] == nil {
		t.Error("Expected aps not nil.")

		aps, ok := message.payload["aps"].(Payload)
		if ok {
			if aps.Badge != 1 {
				t.Errorf("Expected Badge is: %d but found: %d", 1, aps.Badge)
			}
			if aps.Sound != "Default" {
				t.Errorf("Expected Sound is: %s but found: %s", "Default", aps.Sound)
			}
		} else {
			t.Error("Expected aps must be instance of Payload")
		}

	}
}

func TestEncode(t *testing.T) {
	deviceID := bson.NewObjectId().Hex()
	message := CreateMessage(deviceID, "C4XOCR6kmbH4XJ9fMRm1hyt1iL7f0wqfJENdgTDdx+A=", "com.example.appID",
		Payload{
			Alert: "Sample alert",

			ContentAvailable: 1,
		})

	request := message.Encode("")
	if request != nil {
		t.Error("Expected nil when define empty gateway.")
	}

	request = message.Encode(SandboxGateway)
	if request == nil {
		t.Error("Expected not nil.")
	} else {
		if strings.ToLower(request.Method) != "post" {
			t.Errorf("Expected method is: %s but found: %s", "POST", request.Method)
		}
		if request.URL.Host != SandboxGateway {
			t.Errorf("Expected host is: %s but found: %s", SandboxGateway, request.URL.Host)
		}
		if request.URL.Path != fmt.Sprintf("/3/device/%x", message.deviceToken) {
			t.Errorf("Expected host is: %s but found: %s", fmt.Sprintf("/3/device/%x", message.deviceToken), request.URL.Path)
		}

		if request.Header.Get("content-type") != "application/json; charset=utf-8" {
			t.Errorf("Expected host is: %s but found: %s", "application/json; charset=utf-8", request.Header.Get("content-type"))
		}
		if request.Header.Get("apns-id") != message.ApnsID {
			t.Errorf("Expected host is: %s but found: %s", message.ApnsID, request.Header.Get("apns-id"))
		}
		if request.Header.Get("apns-topic") != message.ApnsTopic {
			t.Errorf("Expected host is: %s but found: %s", message.ApnsTopic, request.Header.Get("apns-topic"))
		}
		if request.Header.Get("apns-priority") != fmt.Sprintf("%d", message.ApnsPriority) {
			t.Errorf("Expected host is: %s but found: %s", fmt.Sprintf("%d", message.ApnsPriority), request.Header.Get("apns-priority"))
		}
		if request.Header.Get("apns-expiration") != fmt.Sprintf("%d", message.ApnsExpiration) {
			t.Errorf("Expected host is: %s but found: %s", fmt.Sprintf("%d", message.ApnsExpiration), request.Header.Get("apns-expiration"))
		}
	}
}
