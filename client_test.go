package yellowcard

import (
	"context"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()
	assert.NotNil(t, config)
	assert.Equal(t, _prodBaseURL, config.baseURL)
	assert.Equal(t, EnvironmentProduction, config.env)
	assert.NotNil(t, _httpClient)
}

func TestClient_HttpAuthHeaders(t *testing.T) {
	tests := []struct {
		name   string
		method string
		path   string
		body   []byte
		want   map[string]string
	}{
		{
			name:   "Tests GET request",
			method: http.MethodGet,
			path:   "/test",
			body:   nil,
			want: map[string]string{
				"Accept":         "application/json",
				"Authorization":  "YcHmacV1 key:QEwrY3n9vsnI53x07zW6+XVNWy+933g/zksnwHaKfsU=",
				"X-YC-Timestamp": "2024-06-14T16:20:00Z",
			},
		},
		{
			name:   "Tests POST request",
			method: http.MethodPost,
			path:   "/test",
			body:   []byte(`{"data":"value"}`),
			want: map[string]string{
				"Accept":         "application/json",
				"Content-Type":   "application/json charset=utf-8",
				"Authorization":  "YcHmacV1 key:X/J/HnTGlVDhHEuu1XlEy3Fsy2tpRM+WHduSha2wvbw=",
				"X-YC-Timestamp": "2024-06-14T16:20:00Z",
			},
		},
	}

	var (
		fixedTimeUTC = time.Date(2024, time.June, 14, 16, 20, 0, 0, time.UTC)
		client       = New("key", "secret")
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHeaders := client.getHeaders(tt.method, tt.path, tt.body, fixedTimeUTC)
			assert.Equal(t, tt.want, gotHeaders)
			assert.Len(t, gotHeaders, len(tt.want))
		})
	}
}

func TestClient_GetChannels(t *testing.T) {

	var (
		httpClient = newMockHttpClient()
		client     = New("key", "secret", WithHttpClient(httpClient))
		ctx        = context.Background()
		respBody   = `
		{
		   "channels":[
			  {
				 "max":1000000,
				 "currency":"XAF",
				 "countryCurrency":"CMXAF",
				 "status":"active",
				 "feeLocal":0,
				 "createdAt":"2022-08-10T09:36:35.818Z",
				 "vendorId":"c339d25c-308a-42bf-ba51-d3d37e098264",
				 "country":"CM",
				 "feeUSD":0,
				 "min":1000,
				 "channelType":"momo",
				 "rampType":"deposit",
				 "updatedAt":"2024-03-22T10:00:20.732Z",
				 "apiStatus":"active",
				 "settlementType":"instant",
				 "estimatedSettlementTime":60,
				 "id":"79da4d6e-1c42-4aac-ae7d-422730528f96",
				 "balancer":{
					
				 }
			  },
			  {
				 "max":1000000,
				 "currency":"XAF",
				 "countryCurrency":"CMXAF",
				 "status":"inactive",
				 "feeLocal":0,
				 "createdAt":"2022-08-10T11:52:55.767Z",
				 "vendorId":"0f0683d2-7036-4e01-82b4-89a9b1a9914e",
				 "country":"CM",
				 "feeUSD":0,
				 "min":1000,
				 "channelType":"momo",
				 "rampType":"withdraw",
				 "updatedAt":"2024-03-22T10:00:20.733Z",
				 "apiStatus":"inactive",
				 "settlementType":"instant",
				 "estimatedSettlementTime":60,
				 "id":"402fd1e6-935e-45ff-a39d-2a5b7a57f2cc",
				 "balancer":{}
			  }
		   ]
		}`
	)

	httpClient.MockRequest(client.config.baseURL+"/business/channels", func() (status int, body string) {
		return http.StatusOK, respBody
	})

	// Test gets all channels
	channels, err := client.GetChannels(ctx, "")
	assert.NoError(t, err)
	assert.NotNil(t, channels)
	assert.Len(t, channels, 1)

	// Test ensures country code is valid
	channels, err = client.GetChannels(ctx, "MARS")
	assert.EqualError(t, ErrCountryNotSupported, err.Error())
	assert.Nil(t, channels)

	httpClient.MockRequest(client.config.baseURL+"/business/channels?country=CM", func() (status int, body string) {
		return http.StatusOK, respBody
	})

	// Test returns only active channels for the country code
	channels, err = client.GetChannels(ctx, CountryCodeCM)
	assert.NoError(t, err)
	assert.NotNil(t, channels)
}

func TestNewClient_WithOpts(t *testing.T) {
	client := New("key", "secret")
	assert.NotNil(t, client)
	assert.Equal(t, "key", client.key)
	assert.Equal(t, "secret", client.secret)

	client = New("key", "secret", WithEnvironment(EnvironmentSandbox))
	assert.Equal(t, EnvironmentSandbox, client.config.env)
	assert.Equal(t, _sandboxBaseURL, client.config.baseURL)

	// Test unknown environment
	client = New("key", "secret", WithEnvironment(Environment(2)))
	assert.Equal(t, EnvironmentProduction, client.config.env)
	assert.Equal(t, _prodBaseURL, client.config.baseURL)

	withHttpClient := &http.Client{
		Timeout: 10 * time.Second,
	}

	client = New("key", "secret", WithHttpClient(withHttpClient))
	assert.NotNil(t, client.config.httpClient)
	assert.NotNil(t, client.config.httpClient)
}
