package apns

import "time"

const (
	NO_ERRORS            = "NO_ERRORS"
	PROCESSING_ERROR     = "PROCESSING_ERROR"
	MISSING_DEVICE_TOKEN = "MISSING_DEVICE_TOKEN"
	MISSING_TOPIC        = "MISSING_TOPIC"
	MISSING_PAYLOAD      = "MISSING_PAYLOAD"
	INVALID_TOKEN_SIZE   = "INVALID_TOKEN_SIZE"
	INVALID_TOPIC_SIZE   = "INVALID_TOPIC_SIZE"
	INVALID_PAYLOAD_SIZE = "INVALID_PAYLOAD_SIZE"
	INVALID_TOKEN        = "INVALID_TOKEN"
	SHUTDOWN             = "SHUTDOWN"
	UNKNOWN              = "UNKNOWN"
)

var ResponseCodes = map[uint8]string{
	0:   NO_ERRORS,
	1:   PROCESSING_ERROR,
	2:   MISSING_DEVICE_TOKEN,
	3:   MISSING_TOPIC,
	4:   MISSING_PAYLOAD,
	5:   INVALID_TOKEN_SIZE,
	6:   INVALID_TOPIC_SIZE,
	7:   INVALID_PAYLOAD_SIZE,
	8:   INVALID_TOKEN,
	10:  SHUTDOWN,
	255: UNKNOWN,
}

//////////////////////////////////////////////////////////////////////////////
//-- Feedback --------------------------------------------------------------//
type Feedback struct {
	DeviceToken string
	Timestamp   time.Time
}

//////////////////////////////////////////////////////////////////////////////
//-- Response --------------------------------------------------------------//
type Response struct {
	Success     bool
	Description string
}
