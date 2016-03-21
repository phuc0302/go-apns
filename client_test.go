package apns

import (
	"testing"

	"gopkg.in/mgo.v2/bson"
)

func TestCreateClient(t *testing.T) {
	client, _ := CreateClient("", "", true)
	if client != nil {
		t.Error("Expected nil when define invalid certificate path.")
	}

	client, _ = CreateClient("cert.p12", "", true)
	if client != nil {
		t.Error("Expected nil when certificate file is not available.")
	}

	client, _ = CreateClient("cert.p12", "", true)
	if client != nil {
		t.Error("Expected nil when define invalid certificate's password.")
	}

	client, _ = CreateClient("cert.p12", "P@ssw0rd", true)
	if client != nil {
		t.Error("Expected nil when could not decoded certificate.")
	}

	client, _ = CreateClient("certificate.p12", "P@ssw0rd", true)
	if client == nil {
		t.Error("Expected not nil when everything is okay.")
	}
	if client.client == nil {
		t.Error("Expected http client not nil.")
	}
	if client.gateway != SandboxGateway {
		t.Errorf("Expected %s but found %s.", SandboxGateway, client.gateway)
	}

	client, _ = CreateClient("certificate.p12", "P@ssw0rd", false)
	if client.gateway != Gateway {
		t.Errorf("Expected %s but found %s.", Gateway, client.gateway)
	}
}

func TestSendMessage(t *testing.T) {
	client, _ := CreateClient("certificate.p12", "P@ssw0rd", true)
	deviceID := bson.NewObjectId().Hex()

	message, _ := CreateMessage(deviceID, "C4XOCR6kmbH4XJ9fMRm1hyt1iL7f0wqfJENdgTDdx+A=", "com.example.appID",
		Payload{
			Alert: "Sample alert",
		})

	// Call send message with nil value
	responses := client.SendMessages([]*Message{message})

	for _, response := range responses {
		if response.Status != 200 {
			t.Errorf("Expected status code is %d but found: %d", 200, response.Status)
		}
	}
}
