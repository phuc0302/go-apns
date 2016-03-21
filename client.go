package apns

import (
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/crypto/pkcs12"
	"golang.org/x/net/http2"
)

const (
	// Gateway defines Apple's production server to send apns2
	Gateway = "api.push.apple.com"
	// SandboxGateway defines Apple's development server to send apns2
	SandboxGateway = "api.development.push.apple.com"
)

// Client describes a wrapper to send apns2 message
type Client struct {
	gateway string
	client  *http.Client
}

// CreateClient creates HTTP/2 apns client
func CreateClient(certFile string, password string, isSandbox bool) (*Client, error) {
	/* Condition validation: validate certificate's path */
	if len(certFile) == 0 {
		return nil, fmt.Errorf("Invalid certificate's path.")
	}

	/* Condition validation: validate certificate's password */
	if len(password) == 0 {
		return nil, fmt.Errorf("Invalid certificate's password.")
	}

	// Load certificate
	bytes, _ := ioutil.ReadFile(certFile)
	key, cert, err := pkcs12.Decode(bytes, password)

	/* Condition validation: validate the correctness of loading process */
	if err != nil {
		return nil, err
	}

	// Create certificate
	certificate := tls.Certificate{
		Certificate: [][]byte{cert.Raw},
		PrivateKey:  key,
		Leaf:        cert,
	}

	// Define gateway
	var gateway string
	if !isSandbox {
		gateway = Gateway
	} else {
		gateway = SandboxGateway
	}

	// Config TLS
	config := &tls.Config{
		Certificates: []tls.Certificate{certificate},
		ServerName:   gateway,
	}

	// Config transport
	transport := &http2.Transport{
		TLSClientConfig: config,
	}

	// Finalize
	client := Client{
		gateway: gateway,
		client:  &http.Client{Transport: transport},
	}

	return &client, nil
}

// SendMessages delivers push message to Apple.
func (c *Client) SendMessages(messages []*Message) []*Response {
	/* Condition validation */
	if len(messages) == 0 {
		return nil
	}

	responses := make([]*Response, len(messages))
	for idx, message := range messages {
		// Encode message
		request, err := message.Encode(c.gateway)

		// Create response nomatter what
		response := &Response{}
		responses[idx] = response

		// Define required info
		response.ApnsID = message.ApnsID
		response.deviceID = message.deviceID
		response.deviceToken = base64.StdEncoding.EncodeToString(message.deviceToken)

		/* Condition validation: validate encoding process */
		if err != nil {
			response.Status = 400
			response.StatusDescription = StatusCodes[400]

			response.Reason = err.Error()
			response.ReasonDescription = ReasonCodes[err.Error()]
			continue
		}

		// Send response to Apple server
		res, err := c.client.Do(request)
		if err == nil {
			defer res.Body.Close()

			// Define response status
			response.Status = res.StatusCode
			response.StatusDescription = StatusCodes[res.StatusCode]

			if res.StatusCode != 200 {
				decoder := json.NewDecoder(res.Body)
				err = decoder.Decode(&response)

				if err == nil {
					response.ReasonDescription = ReasonCodes[response.Reason]
				}
			}
		} else {
			response.Status = 503
			response.StatusDescription = StatusCodes[503]

			response.Reason = "ServiceUnavailable"
			response.ReasonDescription = ReasonCodes["ServiceUnavailable"]
		}
	}
	return responses
}
