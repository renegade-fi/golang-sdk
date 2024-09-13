package client

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/renegade-fi/golang-sdk/wallet"
)

const (
	contentTypeHeader   = "Content-Type"
	contentTypeJSON     = "application/json"
	signatureHeader     = "renegade-auth"
	expirationHeader    = "renegade-auth-expiration"
	signatureExpiration = 5 * time.Second
)

// HttpClient represents an HTTP client with a base URL and auth key
type HttpClient struct {
	baseURL    string
	httpClient *http.Client
	authKey    *wallet.HmacKey
}

// NewHttpClient creates a new HttpClient with the given base URL and auth key
func NewHttpClient(baseURL string, authKey *wallet.HmacKey) *HttpClient {
	return &HttpClient{
		baseURL:    baseURL,
		httpClient: &http.Client{},
		authKey:    authKey,
	}
}

// Get performs a GET request to the specified path
func (c *HttpClient) Get(path string, body interface{}) ([]byte, error) {
	return c.doRequest(http.MethodGet, path, body, false /* withAuth */)
}

// Post performs a POST request to the specified path
func (c *HttpClient) Post(path string, body interface{}) ([]byte, error) {
	return c.doRequest(http.MethodPost, path, body, false /* withAuth */)
}

// GetJSON performs a GET request and unmarshals the response into the provided interface
func (c *HttpClient) GetJSON(path string, body interface{}, response interface{}) error {
	respBody, err := c.doRequest(http.MethodGet, path, body, false /* withAuth */)
	if err != nil {
		return err
	}
	return json.Unmarshal(respBody, response)
}

// PostJSON performs a POST request and unmarshals the response into the provided interface
func (c *HttpClient) PostJSON(path string, body interface{}, response interface{}) error {
	respBody, err := c.doRequest(http.MethodPost, path, body, false /* withAuth */)
	if err != nil {
		return err
	}
	return json.Unmarshal(respBody, response)
}

// GetWithAuth performs an authenticated GET request
func (c *HttpClient) GetWithAuth(path string, body interface{}, response interface{}) error {
	respBody, err := c.doRequest(http.MethodGet, path, body, true /* withAuth */)
	if err != nil {
		return err
	}
	return json.Unmarshal(respBody, response)
}

// PostWithAuth performs an authenticated POST request
func (c *HttpClient) PostWithAuth(path string, body interface{}, response interface{}) error {
	respBody, err := c.doRequest(http.MethodPost, path, body, true /* withAuth */)
	if err != nil {
		return err
	}
	return json.Unmarshal(respBody, response)
}

// doRequest performs an HTTP request with optional authentication
func (c *HttpClient) doRequest(method, path string, body interface{}, withAuth bool) ([]byte, error) {
	url := fmt.Sprintf("%s%s", c.baseURL, path)

	// Marshal the body
	var bodyBytes []byte
	var err error
	if body != nil {
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
	}

	// Create the request
	req, err := http.NewRequest(method, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set(contentTypeHeader, contentTypeJSON)
	if withAuth {
		c.addAuth(req, bodyBytes)
	}

	// Send the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read and check the response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// addAuth adds authentication headers to the request
func (c *HttpClient) addAuth(req *http.Request, bodyBytes []byte) {
	// Compute the expiration time
	expiration := time.Now().Add(signatureExpiration * time.Second).UnixMilli()
	expirationBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(expirationBytes, uint64(expiration))

	// Create the hmac
	h := hmac.New(sha256.New, c.authKey[:])
	h.Write(append(bodyBytes, expirationBytes...))
	signature := base64.RawStdEncoding.EncodeToString(h.Sum(nil))

	req.Header.Set(signatureHeader, signature)
	req.Header.Set(expirationHeader, strconv.FormatInt(expiration, 10))
}
