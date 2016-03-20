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
func (c *Client) SendMessages(apns []*Message) []*Response {
	/* Condition validation */
	if len(apns) == 0 {
		return nil
	}

	responses := make([]*Response, len(apns))
	for idx, apn := range apns {
		request := apn.Encode(c.gateway)

		res, err := c.client.Do(request)
		defer res.Body.Close()
		fmt.Println(err)

		response := Response{}
		response.ApnsID = res.Header.Get("apns-id")
		response.deviceID = apn.deviceID
		response.deviceToken = base64.StdEncoding.EncodeToString(apn.deviceToken)

		response.Status = res.StatusCode
		response.StatusDescription = StatusCodes[res.StatusCode]

		if res.StatusCode != 200 {
			decoder := json.NewDecoder(res.Body)
			err = decoder.Decode(&response)

			if err == nil {
				response.ReasonDescription = ReasonCodes[response.Reason]
			}
		}
		responses[idx] = &response
	}
	return responses
}
