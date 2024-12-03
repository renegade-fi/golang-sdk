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
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/renegade-fi/golang-sdk/wallet"
)

const (
	contentTypeHeader       = "Content-Type"
	contentTypeJSON         = "application/json"
	renegadeHeaderNamespace = "x-renegade"
	signatureHeader         = "x-renegade-auth"
	expirationHeader        = "x-renegade-auth-expiration"
	signatureExpiration     = 5 * time.Second
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
	return c.doRequest(http.MethodGet, path, nil /* headers */, body, false /* withAuth */)
}

// Post performs a POST request to the specified path
func (c *HttpClient) Post(path string, body interface{}) ([]byte, error) {
	return c.doRequest(http.MethodPost, path, nil /* headers */, body, false /* withAuth */)
}

// GetJSON performs a GET request and unmarshals the response into the provided interface
func (c *HttpClient) GetJSON(path string, body interface{}, response interface{}) error {
	respBody, err := c.doRequest(http.MethodGet, path, nil /* headers */, body, false /* withAuth */)
	if err != nil {
		return err
	}
	return json.Unmarshal(respBody, response)
}

// PostJSON performs a POST request and unmarshals the response into the provided interface
func (c *HttpClient) PostJSON(path string, body interface{}, response interface{}) error {
	respBody, err := c.doRequest(http.MethodPost, path, nil /* headers */, body, false /* withAuth */)
	if err != nil {
		return err
	}
	return json.Unmarshal(respBody, response)
}

// GetWithAuth performs an authenticated GET request
func (c *HttpClient) GetWithAuth(path string, body interface{}, response interface{}) error {
	return c.GetWithAuthAndHeaders(path, nil /* headers */, body, response)
}

// GetWithAuthAndHeaders performs an authenticated GET request with additional headers
func (c *HttpClient) GetWithAuthAndHeaders(
	path string,
	headers *http.Header,
	body interface{},
	response interface{},
) error {
	respBody, err := c.doRequest(http.MethodGet, path, headers, body, true /* withAuth */)
	if err != nil {
		return err
	}
	return json.Unmarshal(respBody, response)
}

// PostWithAuth performs an authenticated POST request
func (c *HttpClient) PostWithAuth(
	path string,
	body interface{},
	response interface{},
) error {
	return c.PostWithAuthAndHeaders(path, nil /* headers */, body, response)
}

// PostWithAuthAndHeaders performs an authenticated POST request with additional headers
func (c *HttpClient) PostWithAuthAndHeaders(
	path string,
	headers *http.Header,
	body interface{},
	response interface{},
) error {
	respBody, err := c.doRequest(http.MethodPost, path, headers, body, true /* withAuth */)
	if err != nil {
		return err
	}
	return json.Unmarshal(respBody, response)
}

// PostWithAuthRaw performs an authenticated POST request and returns the raw response
func (c *HttpClient) PostWithAuthRaw(
	path string,
	headers *http.Header,
	body interface{},
) (int, []byte, error) {
	return c.doRequestWithStatus(http.MethodPost, path, headers, body, true /* withAuth */)
}

// doRequest performs an HTTP request with optional authentication
func (c *HttpClient) doRequest(
	method,
	path string,
	headers *http.Header,
	body interface{},
	withAuth bool,
) ([]byte, error) {
	_, respBody, err := c.doRequestWithStatus(method, path, headers, body, withAuth)
	return respBody, err
}

// doRequestWithStatus performs an HTTP request with optional authentication and
// returns the raw response with the status code
func (c *HttpClient) doRequestWithStatus(
	method,
	path string,
	headers *http.Header,
	body interface{},
	withAuth bool,
) (int, []byte, error) {
	url := fmt.Sprintf("%s%s", c.baseURL, path)

	// Marshal the body
	var bodyBytes []byte
	var err error
	if body != nil {
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return 0, nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
	}

	// Create the request
	req, err := http.NewRequest(method, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return 0, nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	if headers != nil {
		req.Header = *headers
	}
	req.Header.Set(contentTypeHeader, contentTypeJSON)
	if withAuth {
		c.addAuth(req, bodyBytes)
	}

	// Send the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Read and check the response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check the status code
	statusCode := resp.StatusCode
	if statusCode < 200 || statusCode >= 300 {
		return statusCode, respBody, fmt.Errorf(
			"unexpected status code: %d, body: %s",
			statusCode, string(respBody),
		)
	}

	return statusCode, respBody, nil
}

// addAuth adds authentication headers to the request
func (c *HttpClient) addAuth(req *http.Request, bodyBytes []byte) {
	// Compute the expiration time
	expiration := time.Now().Add(signatureExpiration * time.Second).UnixMilli()
	expirationBytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(expirationBytes, uint64(expiration))
	req.Header.Set(expirationHeader, strconv.FormatInt(expiration, 10))

	// Create the hmac
	h := hmac.New(sha256.New, c.authKey[:])
	hmacPayload := c.getHmacPayload(req.URL.Path, req.Header, bodyBytes)
	h.Write(hmacPayload)

	signature := base64.RawStdEncoding.EncodeToString(h.Sum(nil))
	req.Header.Set(signatureHeader, signature)
}

// getHmacPayload creates the payload for the hmac
func (c *HttpClient) getHmacPayload(path string, headers http.Header, bodyBytes []byte) []byte {
	// Add the path
	payload := []byte(path)

	// Add the headers; filtered only for renegade headers
	var validKeys []string
	for key := range headers {
		lowerKey := strings.ToLower(key)
		if !strings.HasPrefix(lowerKey, renegadeHeaderNamespace) || lowerKey == signatureHeader {
			continue
		}

		validKeys = append(validKeys, key)
	}

	// Add headers in sorted order
	sort.Strings(validKeys)
	for _, key := range validKeys {
		lowerKey := strings.ToLower(key)
		for _, value := range headers[key] {
			payload = append(payload, lowerKey...)
			payload = append(payload, value...)
		}
	}

	// Add the body
	payload = append(payload, bodyBytes...)
	return payload
}
