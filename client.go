package yellowcard

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// _defaultHTTPTimeout is the default timeout on the http.Client used by the library.
const _defaultHTTPTimeout = 60 * time.Second

const (
	_prodBaseURL    = "https://api.yellowcard.io"
	_sandboxBaseURL = "https://sandbox.api.yellowcard.io"
)

var httpClient = &http.Client{
	Timeout: _defaultHTTPTimeout,
}

// Environment represents the type of environment for API servers.
type Environment uint8

const (
	// EnvironmentSandbox indicates the sandbox environment.
	EnvironmentSandbox Environment = iota + 1
	// EnvironmentProduction indicates the production environment.
	EnvironmentProduction
)

// Client represents a client for interacting with the API.
type Client struct {
	config *ClientConfig
	key    string
	secret string
}

// ClientConfig is used to configure a new Client backend.
type ClientConfig struct {
	baseURL    string
	env        Environment
	HTTPClient *http.Client
}

// DefaultConfig returns a default configuration for creating a ClientConfig instance.
// The default base URL for the config is set to
func DefaultConfig() *ClientConfig {
	return &ClientConfig{
		baseURL:    _prodBaseURL,
		env:        EnvironmentProduction,
		HTTPClient: httpClient,
	}
}

// WithEnvironment configures the ClientConfig based on the specified environment.
func WithEnvironment(env Environment) func(config *ClientConfig) {
	return func(config *ClientConfig) {
		switch env {
		case EnvironmentSandbox:
			config.env = env
			config.baseURL = _sandboxBaseURL
		default:
			config.env = EnvironmentProduction
			config.baseURL = _prodBaseURL
		}
	}
}

// WithHttpClient configures the ClientConfig based on the specified http client.
func WithHttpClient(cl *http.Client) func(config *ClientConfig) {
	return func(config *ClientConfig) {
		if cl != nil {
			config.HTTPClient = cl
		}
	}
}

// httpAuthHeaders generates HTTP headers required for authentication using HMAC with SHA-256.
func (cl *Client) httpAuthHeaders(method string, path string, body []byte) map[string]string {
	var (
		now = time.Now().UTC().Format(time.RFC3339)
		mac = hmac.New(sha256.New, []byte(cl.secret))
	)

	mac.Write([]byte(now))
	mac.Write([]byte(path))
	mac.Write([]byte(method))

	if body != nil {
		var (
			bodyHash    = sha256.Sum256(body)
			bodyHmacStr = base64.StdEncoding.EncodeToString(bodyHash[:])
		)

		mac.Write([]byte(bodyHmacStr))
	}

	var (
		sum       = mac.Sum(nil)
		signature = base64.StdEncoding.EncodeToString(sum)
	)

	return map[string]string{
		"Accept":         "application/json",
		"Authorization":  fmt.Sprintf("YcHmacV1 %s:%s", cl.key, signature),
		"X-YC-Timestamp": now,
	}
}

func (cl *Client) call(
	ctx context.Context,
	method string,
	path string,
	body io.Reader,
	params map[string]string,
) ([]byte, error) {
	uri := cl.config.baseURL + path

	req, err := http.NewRequestWithContext(ctx, method, uri, body)
	if err != nil {
		return nil, fmt.Errorf("yellowcard: create request - %v", err)
	}

	var bodyBytes []byte

	if body != nil {
		if _, err = body.Read(bodyBytes); err != nil {
			return nil, fmt.Errorf("yellowcard: read body - %v", err)
		}
	}

	headers := cl.httpAuthHeaders(method, path, bodyBytes)
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if len(params) != 0 {
		q := req.URL.Query()
		for key, value := range params {
			q.Add(key, value)
		}

		req.URL.RawQuery = q.Encode()

		fmt.Println(req.URL.String())
	}

	resp, err := cl.config.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("yellowcard: do request - %v", err)
	}

	defer func(r io.ReadCloser) {
		_ = r.Close()
	}(resp.Body)

	resBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("yellowcard: read response body - %v", err)
	}

	if statusCode := resp.StatusCode; statusCode != http.StatusOK {
		errResp := &errorResponse{StatusCode: statusCode}

		if err = json.Unmarshal(resBody, errResp); err != nil {
			return nil, fmt.Errorf("yellowcard: deserialize error response - %v", err)
		}

		return nil, fmt.Errorf("%s", errResp)
	}

	return resBody, nil
}

// GetChannels retrieves all supported payment ramps (Bank Transfer, Mobile Money, E-Wallets transfers)
// By default only active channels are returned.
func (cl *Client) GetChannels(ctx context.Context, country CountryCode) ([]*Channel, error) {
	params := make(map[string]string)
	if country != "" {
		if _, ok := CountryCodes[country]; !ok {
			return nil, ErrCountryNotSupported
		}

		params["country"] = country.String()
	}

	resBody, err := cl.call(ctx, http.MethodGet, "/business/channels", nil, params)
	if err != nil {
		return nil, err
	}

	var resp *ChannelResponse
	if err = json.Unmarshal(resBody, &resp); err != nil {
		return nil, fmt.Errorf("yellowcard: deserialize channels response - %v", err)
	}

	var activeChannels []*Channel
	for _, channel := range resp.Channels {
		if channel.Status == "active" {
			activeChannels = append(activeChannels, channel)
		}
	}

	return activeChannels, nil
}

// GetNetworks retrieves all supported end financial interfaces (Banks, Mobile Money Networks, E-Wallets)
// By default only active networks are returned.
func (cl *Client) GetNetworks(ctx context.Context, country CountryCode) ([]*Network, error) {
	params := make(map[string]string)
	if country != "" {
		if _, ok := CountryCodes[country]; !ok {
			return nil, ErrCountryNotSupported
		}

		params["country"] = country.String()
	}

	resBody, err := cl.call(ctx, http.MethodGet, "/business/networks", nil, params)
	if err != nil {
		return nil, err
	}

	var resp *NetworksResponse
	if err = json.Unmarshal(resBody, &resp); err != nil {
		return nil, fmt.Errorf("yellowcard: deserialize networks response - %v", err)
	}

	var networks []*Network
	for _, network := range resp.Networks {
		if network.Status == "active" {
			networks = append(networks, network)
		}
	}

	return networks, nil
}

// GetRates retrieves rates for supported countries
func (cl *Client) GetRates(ctx context.Context, currency CurrencyCode) ([]*Rate, error) {
	params := make(map[string]string)
	if currency != "" {
		if _, ok := CurrencyCodes[currency]; !ok {
			return nil, ErrCountryNotSupported
		}

		params["country"] = currency.String()
	}

	resBody, err := cl.call(ctx, http.MethodGet, "/business/rates", nil, params)
	if err != nil {
		return nil, err
	}

	var resp *RatesResponse
	if err = json.Unmarshal(resBody, &resp); err != nil {
		return nil, fmt.Errorf("yellowcard: deserialize rates response - %v", err)
	}

	return resp.Rates, nil
}

// New creates and initializes a new instance of API.
func New(key string, secret string, opts ...func(*ClientConfig)) *Client {
	config := DefaultConfig()

	for _, opt := range opts {
		opt(config)
	}

	return &Client{
		config: config,
		key:    key,
		secret: secret,
	}
}
