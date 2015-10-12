package apns

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/binary"
	"net"
	"strings"
	"time"
)

const (
	GATEWAY          = "gateway.push.apple.com:2195"
	FEEDBACK         = "feedback.push.apple.com:2196"
	SANDBOX_GATEWAY  = "gateway.sandbox.push.apple.com:2195"
	SANDBOX_FEEDBACK = "feedback.sandbox.push.apple.com:2196"
)

type Client struct {
	Gateway  string
	KeyFile  string
	CertFile string

	connection    net.Conn
	tlsConnection net.Conn
}

// MARK: Struct's constructors
func CreatePushClient(certFile string, keyFile string, isSandbox bool) *Client {
	// Decide which endpoint to use
	var gateway string
	if !isSandbox {
		gateway = GATEWAY
	} else {
		gateway = SANDBOX_GATEWAY
	}

	client := Client{
		Gateway:  gateway,
		KeyFile:  keyFile,
		CertFile: certFile,
	}
	return &client
}
func CreateFeedbackClient(certFile string, keyFile string, isSandbox bool) *Client {
	// Decide which endpoint to use
	var gateway string
	if !isSandbox {
		gateway = FEEDBACK
	} else {
		gateway = SANDBOX_FEEDBACK
	}

	client := Client{
		Gateway:  gateway,
		KeyFile:  keyFile,
		CertFile: certFile,
	}
	return &client
}

/**
 * Read feedback from Apple.
 */
func (c *Client) GetFeedback() ([]*Feedback, error) {
	// Connect to apple push server
	err := c.dial()
	defer c.close()
	if err != nil {
		return nil, err
	}

	var feedbacks []*Feedback

	// Read feedback
	timesSamp := uint32(0)
	tokenLength := uint16(0)
	tokenBuffer := make([]byte, 32)
	messageBuffer := make([]byte, 38)
	for {
		binaryLength, err := c.tlsConnection.Read(messageBuffer)
		if binaryLength == 38 {
			reader := bytes.NewReader(messageBuffer)
			binary.Read(reader, binary.BigEndian, &timesSamp)
			binary.Read(reader, binary.BigEndian, &tokenLength)
			binary.Read(reader, binary.BigEndian, &tokenBuffer)

			feedback := &Feedback{
				Timestamp:   time.Unix(int64(timesSamp), 0),
				DeviceToken: base64.StdEncoding.EncodeToString(tokenBuffer),
			}
			feedbacks = append(feedbacks, feedback)
		}

		// Terminate loop if there is nothing else to do
		if err != nil {
			break
		}
	}
	return feedbacks, nil
}

/**
 * Send single or multiple messages.
 */
func (c *Client) SendMessage(apns []*Message) *Response {
	/* Condition validation */
	if apns == nil || len(apns) == 0 {
		return &Response{
			Success:     false,
			Description: SHUTDOWN,
		}
	}

	// Connect to apple push server
	err := c.dial()
	defer c.close()
	if err != nil {
		return &Response{
			Success:     false,
			Description: SHUTDOWN,
		}
	}

	// Write message
	buffer := make([]byte, int(apns[0].Length()))
	for _, apn := range apns {
		reader, _ := apn.Encode()
		reader.Read(buffer)

		_, err := c.tlsConnection.Write(buffer)
		if err != nil {
			return &Response{
				Success:     false,
				Description: PROCESSING_ERROR,
			}
		}
	}

	// Read result
	responseChannel := make(chan []byte, 1)
	timeoutChannel := make(chan bool, 1)
	go func() {
		end, _ := c.tlsConnection.Read(buffer)
		responseChannel <- buffer[:end]
	}()
	go func() {
		time.Sleep(time.Second * 5)
		timeoutChannel <- true
	}()

	/**
	 * First one back wins! The data structure for an APN response is as follows:
	 * command    -> 1 byte (will always be set to 8)
	 * status     -> 1 byte
	 * identifier -> 4 bytes
	 */
	response := &Response{}
	select {
	case r := <-responseChannel:
		response.Success = false
		response.Description = ResponseCodes[r[1]]

	case <-timeoutChannel:
		response.Success = true
		response.Description = NO_ERRORS
	}

	return response
}

// MARK: Struct's private functions
/**
 * Close connection with apple push server.
 */
func (c *Client) close() {
	if c.tlsConnection != nil {
		c.tlsConnection.Close()
		c.tlsConnection = nil
	}

	if c.connection != nil {
		c.connection.Close()
		c.connection = nil
	}
}

/**
 * Open connection with apple push server.
 */
func (c *Client) dial() error {
	// Load keypair
	kp, err := tls.LoadX509KeyPair(c.CertFile, c.KeyFile)
	if err != nil {
		return err
	}

	// Config TLS
	tokens := strings.Split(c.Gateway, ":")
	config := &tls.Config{
		Certificates: []tls.Certificate{kp},
		ServerName:   tokens[0],
	}

	// Connect to Apple gateway
	conn, err := net.Dial("tcp", c.Gateway)
	if err != nil {
		return err
	} else {
		c.connection = conn
	}

	// Handshake
	tlsConn := tls.Client(conn, config)
	err = tlsConn.Handshake()
	if err != nil {
		return err
	} else {
		c.tlsConnection = tlsConn
		return nil
	}
}
