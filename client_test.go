package yellowcard

import (
	"context"
	"fmt"
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

	ctx := context.Background()

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

func TestClient_GetNetworks(t *testing.T) {
	var (
		httpClient = newMockHttpClient()
		client     = New("key", "secret", WithHttpClient(httpClient))
		respBody   = `
		{
		   "networks":[
			  {
				 "code":"589000",
				 "updatedAt":"2023-09-25T14:48:05.578Z",
				 "status":"active",
				 "channelIds":[
					"81018280-e320-4c81-9b2f-6f636c2239d8"
				 ],
				 "createdAt":"2023-09-25T14:48:05.578Z",
				 "accountNumberType":"bank",
				 "id":"41109c18-9604-4389-8472-44ff4378c6cb",
				 "country":"ZA",
				 "name":"Finbond Mutual Bank",
				 "countryAccountNumberType":"ZABANK"
			  },
			  {
				 "code":"450905",
				 "updatedAt":"2023-09-25T14:48:05.578Z",
				 "status":"active",
				 "channelIds":[
					"81018280-e320-4c81-9b2f-6f636c2239d8"
				 ],
				 "createdAt":"2023-09-25T14:48:05.578Z",
				 "accountNumberType":"bank",
				 "id":"0d19d67e-5946-4289-bac1-ad147d7c84ad",
				 "country":"ZA",
				 "name":"Mercantile Bank Limited",
				 "countryAccountNumberType":"ZABANK"
			  },
              {
				 "code":"450905",
				 "updatedAt":"2023-09-25T14:48:05.578Z",
				 "status":"inactive",
				 "channelIds":[
					"81018280-e320-4c81-9b2f-6f636c2239d8"
				 ],
				 "createdAt":"2023-09-25T14:48:05.578Z",
				 "accountNumberType":"bank",
				 "id":"0d19d67e-5946-4289-bac1-ad147d7c84ad",
				 "country":"ZA",
				 "name":"Mercantile Bank Limited",
				 "countryAccountNumberType":"ZABANK"
			  }
		   ]
		}`
	)

	httpClient.MockRequest(client.config.baseURL+"/business/networks", func() (status int, body string) {
		return http.StatusOK, respBody
	})

	ctx := context.Background()

	// Test gets all networks
	networks, err := client.GetNetworks(ctx, "")
	assert.NoError(t, err)
	assert.NotNil(t, networks)
	assert.Len(t, networks, 2)

	// Test ensures country code is valid
	networks, err = client.GetNetworks(ctx, "MARS")
	assert.EqualError(t, ErrCountryNotSupported, err.Error())
	assert.Nil(t, networks)

	httpClient.MockRequest(client.config.baseURL+"/business/networks?country=ZA", func() (status int, body string) {
		return http.StatusOK, respBody
	})

	// Test returns only active networks for the country code
	networks, err = client.GetNetworks(ctx, CountryCodeZA)
	assert.NoError(t, err)
	assert.NotNil(t, networks)
	assert.Len(t, networks, 2)
}

func TestClient_GetRates(t *testing.T) {
	var (
		httpClient = newMockHttpClient()
		client     = New("key", "secret", WithHttpClient(httpClient))
		respBody   = `
		{
		   "rates":[
			  {
				 "buy":2615,
				 "sell":2615,
				 "locale":"TZ",
				 "rateId":"tanzanian-shilling",
				 "code":"TZS",
				 "updatedAt":"2024-06-10T13:52:24.739Z"
			  },
			  {
				 "locale":"crypto",
				 "rateId":"ethereum",
				 "code":"ETH",
				 "updatedAt":"2024-05-03T09:26:00.570Z"
			  },
			  {
				 "buy":13,
				 "sell":-13.07,
				 "locale":"Bw",
				 "rateId":"pula",
				 "code":"BWP",
				 "updatedAt":"2024-06-10T13:52:35.713Z"
			  }
		   ]
		}`
	)

	httpClient.MockRequest(client.config.baseURL+"/business/rates", func() (status int, body string) {
		return http.StatusOK, respBody
	})

	ctx := context.Background()

	// Test gets all rates
	rates, err := client.GetRates(ctx, "")
	assert.NoError(t, err)
	assert.NotNil(t, rates)
	assert.Len(t, rates, 3)

	// Test ensures country code is valid
	rates, err = client.GetRates(ctx, "MARS")
	assert.EqualError(t, ErrCurrencyCodeNotSupported, err.Error())
	assert.Nil(t, rates)

	httpClient.MockRequest(client.config.baseURL+"/business/rates?currency=TZS", func() (status int, body string) {
		return http.StatusOK, `
		{
		   "rates":[
			  {
				 "buy":2615,
				 "sell":2615,
				 "locale":"TZ",
				 "rateId":"tanzanian-shilling",
				 "code":"TZS",
				 "updatedAt":"2024-06-10T13:52:24.739Z"
			  }
		   ]
		}`
	})

	// Test returns only rates for the currency
	rates, err = client.GetRates(ctx, CurrencyCodeTZS)
	assert.NoError(t, err)
	assert.NotNil(t, rates)
	assert.Len(t, rates, 1)
	assert.Equal(t, float64(2615), rates[0].Buy)
	assert.Equal(t, "TZ", rates[0].Locale)
}

func TestClient_ResolveBankAccount(t *testing.T) {
	var (
		httpClient = newMockHttpClient()
		client     = New("key", "secret", WithHttpClient(httpClient))
	)

	httpClient.MockRequest(client.config.baseURL+"/business/details/bank", func() (status int, body string) {
		return http.StatusOK, `
		{
		   "accountNumber":"589000",
		   "accountName":"Ken Adams",
		   "accountBank":"Finbond Mutual Bank"
		}`
	})

	ctx := context.Background()

	bankAccountDetails, err := client.ResolveBankAccount(ctx, &ResolveBankAccountRequest{
		AccountNumber: "589000",
		NetworkID:     "41109c18-9604-4389-8472-44ff4378c6cb",
	})

	assert.NoError(t, err)
	assert.NotNil(t, bankAccountDetails)
	assert.Equal(t, "589000", bankAccountDetails.AccountNumber)
	assert.Equal(t, "Ken Adams", bankAccountDetails.AccountName)
}

func TestClient_MakePayment(t *testing.T) {
	var (
		httpClient = newMockHttpClient()
		client     = New("key", "secret", WithHttpClient(httpClient))
		uri        = client.config.baseURL + "/business/payments"
	)

	httpClient.MockRequest(uri, func() (status int, body string) {
		return http.StatusOK, `
		{
		   "amount":7491.65,
		   "channelId":"81018280-e320-4c81-9b2f-6f636c2239d8",
		   "destination":{
			  "accountBank":"589000",
			  "accountName":"Ken Adams",
			  "accountNumber":"+12222222222",
			  "accountType":"momo",
			  "networkId":"41109c18-9604-4389-8472-44ff4378c6cb",
			  "networkName":"Finbond Mutual Bank"
		   },
		   "reason":"entertainment",
		   "sender":{
			  "address":"Sample Address",
			  "country":"US",
			  "dob":"10/10/1950",
			  "email":"email@domain.com",
			  "idNumber":"0123456789",
			  "idType":"license",
			  "name":"Sample Name",
			  "phone":"+12222222222"
		   },
		   "sequenceId":"nsahHJODjx",
		   "partnerId":"deb55c03-9961-417a-9550-f5ba7fe258e9",
		   "requestSource":"api",
		   "id":"0aa5bd35-b969-5d1d-ae7b-dfc0c4abbaf7",
		   "status":"created",
		   "currency":"ZAR",
		   "country":"ZA",
		   "convertedAmount":139119.94,
		   "rate":18.57,
		   "forceAccept":false,
		   "expiresAt":"2024-06-15T07:04:25.573Z",
		   "settlementInfo":{
			  
		   },
		   "createdAt":"2024-06-15T06:54:25.576Z",
		   "updatedAt":"2024-06-15T06:54:25.576Z",
		   "directSettlement":false
		}`
	})

	var (
		paymentRequest = &PaymentRequest{
			Amount:    7491.65,
			ChannelID: "81018280-e320-4c81-9b2f-6f636c2239d8",
			Destination: Destination{
				AccountBank:   "589000",
				AccountName:   "Ken Adams",
				AccountNumber: "+12222222222",
				AccountType:   AccountTypeMobileMoney,
				Country:       "ZA",
				NetworkID:     "41109c18-9604-4389-8472-44ff4378c6cb",
			},
			Reason: "entertainment",
			Sender: Sender{
				Address:  "Sample Address",
				Country:  "US",
				Dob:      "10/10/1950",
				Email:    "email@domain.com",
				IDNumber: "0123456789",
				IDType:   "license",
				Name:     "Sample Name",
				Phone:    "+12222222222",
			},
			SequenceID: "nsahHJODjx",
		}
		ctx = context.Background()
	)

	payment, err := client.MakePayment(ctx, paymentRequest, false)
	assert.NoError(t, err)
	assert.NotNil(t, payment)
	assert.Equal(t, paymentRequest.Amount, payment.Amount)
	assert.Equal(t, paymentRequest.Destination.AccountBank, payment.Destination.AccountBank)
}

func TestClient_AcceptPaymentRequest(t *testing.T) {
	var (
		httpClient = newMockHttpClient()
		client     = New("key", "secret", WithHttpClient(httpClient))
		paymentID  = "d83011e8-341f-5e3e-b908-84cb4a552fcc"
		uri        = fmt.Sprintf("%s/business/payments/%s/accept", client.config.baseURL, paymentID)
	)

	httpClient.MockRequest(uri, func() (status int, body string) {
		return http.StatusOK, `
		{
		   "partnerId":"deb55c03-9961-417a-9550-f5ba7fe258e9",
		   "currency":"ZAR",
		   "rate":18.57,
		   "settlementInfo":{
			  
		   },
		   "status":"process",
		   "createdAt":"2024-06-15T07:09:36.607Z",
		   "forceAccept":false,
		   "serviceFeeAmountUSD":37.46,
		   "sequenceId":"AovQYlKGkz",
		   "country":"ZA",
		   "reason":"entertainment",
		   "sender":{
			  "country":"US",
			  "address":"Sample Address",
			  "idType":"license",
			  "phone":"+12222222222",
			  "dob":"10/10/1950",
			  "name":"Sample Name",
			  "idNumber":"0123456789",
			  "email":"email@domain.com"
		   },
		   "convertedAmount":139119.94,
		   "channelId":"81018280-e320-4c81-9b2f-6f636c2239d8",
		   "expiresAt":"2024-06-15T07:19:36.607Z",
		   "serviceFeeAmountLocal":695.63,
		   "requestSource":"api",
		   "updatedAt":"2024-06-15T07:10:34.394Z",
		   "directSettlement":false,
		   "amount":7491.65,
		   "destination":{
			  "networkName":"Finbond Mutual Bank",
			  "accountBank":"589000",
			  "networkId":"41109c18-9604-4389-8472-44ff4378c6cb",
			  "accountNumber":"+12222222222",
			  "accountName":"Ken Adams",
			  "accountType":"momo"
		   },
		   "id":"d83011e8-341f-5e3e-b908-84cb4a552fcc"
		}`
	})

	ctx := context.Background()

	payment, err := client.AcceptPaymentRequest(ctx, paymentID)
	assert.NoError(t, err)
	assert.NotNil(t, payment)
	assert.Equal(t, paymentID, payment.ID)
	assert.Equal(t, "process", payment.Status)
}

func TestClient_DenyPaymentRequest(t *testing.T) {
	var (
		httpClient = newMockHttpClient()
		client     = New("key", "secret", WithHttpClient(httpClient))
		paymentID  = "c1de8da5-c11a-5cff-a17a-3e7c7085044c"
		uri        = fmt.Sprintf("%s/business/payments/%s/deny", client.config.baseURL, paymentID)
	)

	httpClient.MockRequest(uri, func() (status int, body string) {
		return http.StatusOK, `
		{
		   "partnerId":"deb55c03-9961-417a-9550-f5ba7fe258e9",
		   "currency":"ZAR",
		   "rate":18.57,
		   "settlementInfo":{
			  
		   },
		   "status":"denied",
		   "createdAt":"2024-06-15T07:40:57.944Z",
		   "forceAccept":false,
		   "serviceFeeAmountUSD":37.46,
		   "sequenceId":"ZEmcaXRAPc",
		   "country":"ZA",
		   "reason":"entertainment",
		   "sender":{
			  "country":"US",
			  "address":"Sample Address",
			  "idType":"license",
			  "phone":"+12222222222",
			  "dob":"10/10/1950",
			  "name":"Sample Name",
			  "idNumber":"0123456789",
			  "email":"email@domain.com"
		   },
		   "convertedAmount":139119.94,
		   "channelId":"81018280-e320-4c81-9b2f-6f636c2239d8",
		   "expiresAt":"2024-06-15T07:50:57.943Z",
		   "serviceFeeAmountLocal":695.63,
		   "requestSource":"api",
		   "updatedAt":"2024-06-15T07:41:24.624Z",
		   "directSettlement":false,
		   "amount":7491.65,
		   "destination":{
			  "networkName":"Finbond Mutual Bank",
			  "accountBank":"589000",
			  "networkId":"41109c18-9604-4389-8472-44ff4378c6cb",
			  "accountNumber":"+12222222222",
			  "accountName":"Ken Adams",
			  "accountType":"momo"
		   },
		   "id":"c1de8da5-c11a-5cff-a17a-3e7c7085044c"
		}`
	})

	ctx := context.Background()

	payment, err := client.DenyPaymentRequest(ctx, paymentID)
	assert.NoError(t, err)
	assert.NotNil(t, payment)
	assert.Equal(t, paymentID, payment.ID)
	assert.Equal(t, "denied", payment.Status)
}

func TestClient_AcceptPaymentRequestInvalidState(t *testing.T) {
	var (
		httpClient = newMockHttpClient()
		client     = New("key", "secret", WithHttpClient(httpClient))
		paymentID  = "c1de8da5-c11a-5cff-a17a-3e7c7085044c"
		uri        = fmt.Sprintf("%s/business/payments/%s/accept", client.config.baseURL, paymentID)
	)

	httpClient.MockRequest(uri, func() (status int, body string) {
		return http.StatusBadRequest, `
		{
		   "code":"PaymentInvalidState",
		   "message":"payment is not in pending_approval state"
		}`
	})

	ctx := context.Background()

	payment, err := client.AcceptPaymentRequest(ctx, paymentID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "PaymentInvalidState")
	assert.Contains(t, err.Error(), "payment is not in pending_approval state")
	assert.Nil(t, payment)
}

func TestClient_DeclinePaymentRequestInvalidState(t *testing.T) {
	var (
		httpClient = newMockHttpClient()
		client     = New("key", "secret", WithHttpClient(httpClient))
		paymentID  = "c1de8da5-c11a-5cff-a17a-3e7c7085044c"
		uri        = fmt.Sprintf("%s/business/payments/%s/deny", client.config.baseURL, paymentID)
	)

	httpClient.MockRequest(uri, func() (status int, body string) {
		return http.StatusBadRequest, `
		{
		   "code":"PaymentInvalidState",
		   "message":"payment is not in pending_approval state"
		}`
	})

	ctx := context.Background()

	payment, err := client.DenyPaymentRequest(ctx, paymentID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "PaymentInvalidState")
	assert.Contains(t, err.Error(), "payment is not in pending_approval state")
	assert.Nil(t, payment)
}

func TestClient_LookupPayment(t *testing.T) {
	var (
		httpClient = newMockHttpClient()
		client     = New("key", "secret", WithHttpClient(httpClient))
		paymentID  = "c1de8da5-c11a-5cff-a17a-3e7c7085044c"
		uri        = fmt.Sprintf("%s/business/payments/%s", client.config.baseURL, paymentID)
	)

	httpClient.MockRequest(uri, func() (status int, body string) {
		return http.StatusOK, `
		{
		   "partnerId":"deb55c03-9961-417a-9550-f5ba7fe258e9",
		   "currency":"ZAR",
		   "rate":18.57,
		   "settlementInfo":{
			  
		   },
		   "status":"denied",
		   "createdAt":"2024-06-15T07:40:57.944Z",
		   "forceAccept":false,
		   "serviceFeeAmountUSD":37.46,
		   "sequenceId":"ZEmcaXRAPc",
		   "country":"ZA",
		   "reason":"entertainment",
		   "sender":{
			  "country":"US",
			  "address":"Sample Address",
			  "idType":"license",
			  "phone":"+12222222222",
			  "dob":"10/10/1950",
			  "name":"Sample Name",
			  "idNumber":"0123456789",
			  "email":"email@domain.com"
		   },
		   "convertedAmount":139119.94,
		   "channelId":"81018280-e320-4c81-9b2f-6f636c2239d8",
		   "expiresAt":"2024-06-15T07:50:57.943Z",
		   "serviceFeeAmountLocal":695.63,
		   "requestSource":"api",
		   "updatedAt":"2024-06-15T07:41:24.624Z",
		   "directSettlement":false,
		   "amount":7491.65,
		   "destination":{
			  "networkName":"Finbond Mutual Bank",
			  "accountBank":"589000",
			  "networkId":"41109c18-9604-4389-8472-44ff4378c6cb",
			  "accountNumber":"+12222222222",
			  "accountName":"Ken Adams",
			  "accountType":"momo"
		   },
		   "id":"c1de8da5-c11a-5cff-a17a-3e7c7085044c"
		}`
	})

	ctx := context.Background()

	payment, err := client.LookupPayment(ctx, paymentID)
	assert.NoError(t, err)
	assert.NotNil(t, payment)
	assert.Equal(t, paymentID, payment.ID)
	assert.Equal(t, "denied", payment.Status)
}

func TestClient_LookupPaymentNotFound(t *testing.T) {
	var (
		httpClient = newMockHttpClient()
		client     = New("key", "secret", WithHttpClient(httpClient))
		paymentID  = "c1de8da5-c11a-5cff-a17a-3e7c7085044c"
		uri        = fmt.Sprintf("%s/business/payments/%s", client.config.baseURL, paymentID)
	)

	httpClient.MockRequest(uri, func() (status int, body string) {
		return http.StatusNotFound, `
		{
		   "code": "PaymentNotFound:",
		   "message": "payment not found"
		}`
	})

	ctx := context.Background()

	payment, err := client.LookupPayment(ctx, paymentID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "PaymentNotFound")
	assert.Nil(t, payment)
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
