package apns

import (
	"fmt"
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

	message := CreateMessage(deviceID, "C4XOCR6kmbH4XJ9fMRm1hyt1iL7f0wqfJENdgTDdx+A=", "com.example.appID",
		Payload{
			Alert: "Sample alert",
		})

	// Call send message with nil value
	responses := client.SendMessages([]*Message{message})

	for _, response := range responses {
		fmt.Println(response)
	}

	// 	// Call send message
	// 	message1 := CreateMessage("5tWyQwxLip+3HjGiov3aB9DM+KGPOYc4VTVgp6s7u3c=", "")
	// 	message2 := CreateMessage("ob5T/152seL4KHn+Nj0KSXNar4euFbCFgQnbWx0KymI=", "")
	// 	message1.SetPayload(&Payload{
	// 		Alert: "This is a really long message from Sunbox! ^_^",
	// 	})
	// 	message2.SetPayload(&Payload{
	// 		Alert: "This is a really long message from Sunbox! ^_^",
	// 	})
	// 	response = client.SendMessage([]*Message{message1, message2})

	// 	if !response.Success {
	// 		t.Errorf("Expect %t but found %t", true, response.Success)
	// 	}

	// 	if response.Description != "NO_ERRORS" {
	// 		t.Errorf("Expect %s but found %s", "NO_ERRORS", response.Description)
	// 	}
	// }

	// func TestFeedback(t *testing.T) {
	// 	client := CreateFeedbackClient("Pro_Push_Sunbox_Crt.pem", "Pro_Push_Sunbox_Key.pem", false)
	// 	feedbacks, err := client.GetFeedback()

	// 	if err != nil {
	// 		t.Error("Expect nil but found not nil")
	// 	}

	// 	if len(feedbacks) != 2 {
	// 		t.Errorf("Expect 2 but found %d", len(feedbacks))
	// 	}
	// 	fmt.Println(feedbacks)
}

// func TestFeedbackBusiness(t *testing.T) {
// 	messageBuffer := []byte{86, 27, 64, 14, 0, 32, 233, 251, 228, 250, 192, 242, 138, 67, 177, 20, 157, 53, 194, 234, 93, 87, 106, 212, 255, 106, 178, 181, 74, 141, 251, 75, 187, 219, 206, 141, 216, 5}
// 	var feedbacks []*Feedback

// 	// Begin process /////////////////////////////////////////////////////////
// 	timesSamp := uint32(0)
// 	tokenLength := uint16(0)
// 	tokenBuffer := make([]byte, 32)

// 	reader := bytes.NewReader(messageBuffer)
// 	binary.Read(reader, binary.BigEndian, &timesSamp)
// 	binary.Read(reader, binary.BigEndian, &tokenLength)
// 	binary.Read(reader, binary.BigEndian, &tokenBuffer)

// 	feedback := &Feedback{
// 		Timestamp:   time.Unix(int64(timesSamp), 0),
// 		DeviceToken: base64.StdEncoding.EncodeToString(tokenBuffer),
// 	}
// 	feedbacks = append(feedbacks, feedback)
// 	// End process ///////////////////////////////////////////////////////////

// 	if len(feedbacks) != 1 {
// 		t.Errorf("Expect 1 but found %d", len(feedbacks))
// 	}

// 	if tokenLength != 32 {
// 		t.Errorf("Expect 32 but found %d", tokenLength)
// 	}

// 	if feedback.DeviceToken != "6fvk+sDyikOxFJ01wupdV2rU/2qytUqN+0u7286N2AU=" {
// 		t.Errorf("Expect %s but found %d", "6fvk+sDyikOxFJ01wupdV2rU/2qytUqN+0u7286N2AU=", feedback.DeviceToken)
// 	}
// }
