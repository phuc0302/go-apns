package apns

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"testing"
	"time"
)

func TestCreateClient(t *testing.T) {
	client := CreatePushClient("cert.pem", "key.pem", false)
	if client.Gateway != GATEWAY {
		t.Errorf("Expect %s but found %s", GATEWAY, client.Gateway)
	}

	client = CreatePushClient("cert.pem", "key.pem", true)
	if client.Gateway != SANDBOX_GATEWAY {
		t.Errorf("Expect %s but found %s", SANDBOX_GATEWAY, client.Gateway)
	}

	client = CreateFeedbackClient("cert.pem", "key.pem", false)
	if client.Gateway != FEEDBACK {
		t.Errorf("Expect %s but found %s", FEEDBACK, client.Gateway)
	}

	client = CreateFeedbackClient("cert.pem", "key.pem", true)
	if client.Gateway != SANDBOX_FEEDBACK {
		t.Errorf("Expect %s but found %s", SANDBOX_FEEDBACK, client.Gateway)
	}
}

func TestSendMessage(t *testing.T) {
	client := CreatePushClient("cert.pem", "key.pem", true)

	// Call send message with nil value
	response := client.SendMessage(nil)
	if response == nil {
		t.Error("Expect not nil but found nil")
	} else {
		if response.Success || response.Description != "SHUTDOWN" {
			t.Errorf("Expect %s but found %s", "SHUTDOWN", response.Description)
		}
	}

	// Call send message with empty slide
	response = client.SendMessage([]*Message{})
	if response == nil {
		t.Error("Expect not nil but found nil")
	} else {
		if response.Success || response.Description != "SHUTDOWN" {
			t.Errorf("Expect %s but found %s", "SHUTDOWN", response.Description)
		}
	}

	// Call send message
	message1 := CreateMessage("5tWyQwxLip+3HjGiov3aB9DM+KGPOYc4VTVgp6s7u3c=", "")
	message2 := CreateMessage("ob5T/152seL4KHn+Nj0KSXNar4euFbCFgQnbWx0KymI=", "")
	message1.SetPayload(&Payload{
		Alert: "This is a really long message from Sunbox! ^_^",
	})
	message2.SetPayload(&Payload{
		Alert: "This is a really long message from Sunbox! ^_^",
	})
	response = client.SendMessage([]*Message{message1, message2})

	if !response.Success {
		t.Errorf("Expect %t but found %t", true, response.Success)
	}

	if response.Description != "NO_ERRORS" {
		t.Errorf("Expect %s but found %s", "NO_ERRORS", response.Description)
	}
}

func TestFeedback(t *testing.T) {
	client := CreateFeedbackClient("Pro_Push_Sunbox_Crt.pem", "Pro_Push_Sunbox_Key.pem", false)
	feedbacks, err := client.GetFeedback()

	if err != nil {
		t.Error("Expect nil but found not nil")
	}

	if len(feedbacks) != 2 {
		t.Errorf("Expect 2 but found %d", len(feedbacks))
	}
	fmt.Println(feedbacks)
}

func TestFeedbackBusiness(t *testing.T) {
	messageBuffer := []byte{86, 27, 64, 14, 0, 32, 233, 251, 228, 250, 192, 242, 138, 67, 177, 20, 157, 53, 194, 234, 93, 87, 106, 212, 255, 106, 178, 181, 74, 141, 251, 75, 187, 219, 206, 141, 216, 5}
	var feedbacks []*Feedback

	// Begin process /////////////////////////////////////////////////////////
	timesSamp := uint32(0)
	tokenLength := uint16(0)
	tokenBuffer := make([]byte, 32)

	reader := bytes.NewReader(messageBuffer)
	binary.Read(reader, binary.BigEndian, &timesSamp)
	binary.Read(reader, binary.BigEndian, &tokenLength)
	binary.Read(reader, binary.BigEndian, &tokenBuffer)

	feedback := &Feedback{
		Timestamp:   time.Unix(int64(timesSamp), 0),
		DeviceToken: base64.StdEncoding.EncodeToString(tokenBuffer),
	}
	feedbacks = append(feedbacks, feedback)
	// End process ///////////////////////////////////////////////////////////

	if len(feedbacks) != 1 {
		t.Errorf("Expect 1 but found %d", len(feedbacks))
	}

	if tokenLength != 32 {
		t.Errorf("Expect 32 but found %d", tokenLength)
	}

	if feedback.DeviceToken != "6fvk+sDyikOxFJ01wupdV2rU/2qytUqN+0u7286N2AU=" {
		t.Errorf("Expect %s but found %d", "6fvk+sDyikOxFJ01wupdV2rU/2qytUqN+0u7286N2AU=", feedback.DeviceToken)
	}
}
