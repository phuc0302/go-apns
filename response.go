package apns

// ReasonCodes defines keys & descriptions.
var ReasonCodes = map[string]string{
	"PayloadEmpty":              "The message payload was empty.",
	"PayloadTooLarge":           "The message payload was too large. The maximum payload size is 4096 bytes.",
	"BadTopic":                  "The apns-topic was invalid.",
	"TopicDisallowed":           "Pushing to this topic is not allowed.",
	"BadMessageId":              "The apns-id value is bad.",
	"BadExpirationDate":         "The apns-expiration value is bad.",
	"BadPriority":               "The apns-priority value is bad.",
	"MissingDeviceToken":        "The device token is not specified in the request :path. Verify that the :path header contains the device token.",
	"BadDeviceToken":            "The specified device token was bad. Verify that the request contains a valid token and that the token matches the environment.",
	"DeviceTokenNotForTopic":    "The device token does not match the specified topic.",
	"Unregistered":              "The device token is inactive for the specified topic.",
	"DuplicateHeaders":          "One or more headers were repeated.",
	"BadCertificateEnvironment": "The client certificate was for the wrong environment.",
	"BadCertificate":            "The certificate was bad.",
	"Forbidden":                 "The specified action is not allowed.",
	"BadPath":                   "The request contained a bad :path value.",
	"MethodNotAllowed":          "The specified :method was not POST.",
	"TooManyRequests":           "Too many requests were made consecutively to the same device token.",
	"IdleTimeout":               "Idle time out.",
	"Shutdown":                  "The server is shutting down.",
	"InternalServerError":       "An internal server error occurred.",
	"ServiceUnavailable":        "The service is unavailable.",
	"MissingTopic":              "The apns-topic header of the request was not specified and was required. The apns-topic header is mandatory when the client is connected using a certificate that supports multiple topics.",
}

// StatusCodes defines keys & descriptions.
var StatusCodes = map[int]string{
	200: "Success.",
	400: "Bad request.",
	403: "There was an error with the certificate.",
	405: "The request used a bad method value. Only \"POST\" requests are supported.",
	410: "The device token is no longer active for the topic.",
	413: "The notification payload was too large.",
	429: "The server received too many requests for the same device token.",
	500: "Internal server error.",
	503: "The server is shutting down and unavailable.",
}

// Response describes a response from Apple push server.
type Response struct {
	ApnsID      string
	deviceID    string
	deviceToken string

	Status            int
	StatusDescription string
	ReasonDescription string

	Reason    string `json:"reason,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
}
