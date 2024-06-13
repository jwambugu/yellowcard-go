package yellowcard

import (
	"fmt"
	"time"
)

// errorResponse represents an error response received from the API.
type errorResponse struct {
	Code       string `json:"code"`
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

func (e errorResponse) String() string {
	return fmt.Sprintf("yellowcard: request failed with status [%d] %s: %s", e.StatusCode, e.Code, e.Message)
}

// Channel is specific financial mechanism used to facilitate a payment.
type Channel struct {
	ApiStatus               string    `json:"apiStatus"`
	Balancer                any       `json:"balancer"`
	ChannelType             string    `json:"channelType"`
	Country                 string    `json:"country"`
	CountryCurrency         string    `json:"countryCurrency"`
	CreatedAt               time.Time `json:"createdAt"`
	Currency                string    `json:"currency"`
	EstimatedSettlementTime int       `json:"estimatedSettlementTime"`
	FeeLocal                int       `json:"feeLocal"`
	FeeUSD                  int       `json:"feeUSD"`
	Id                      string    `json:"id"`
	Max                     int       `json:"max"`
	Min                     float64   `json:"min"`
	RampType                string    `json:"rampType"`
	SettlementType          string    `json:"settlementType"`
	Status                  string    `json:"status"`
	SuccessThreshold        int       `json:"successThreshold,omitempty"`
	UpdatedAt               time.Time `json:"updatedAt"`
	VendorId                string    `json:"vendorId"`
	WidgetStatus            string    `json:"widgetStatus,omitempty"`
}

// ChannelResponse represents the response returned by the get channels request
type ChannelResponse struct {
	Channels []*Channel `json:"channels"`
}

// Network is a company, bank, or service that the end-user interfaces financially with.
// There can be multiple Channel(s) linked to a Network.
type Network struct {
	AccountNumberType        string    `json:"accountNumberType"`
	ChannelIds               []string  `json:"channelIds"`
	Code                     string    `json:"code,omitempty"`
	Country                  string    `json:"country"`
	CountryAccountNumberType string    `json:"countryAccountNumberType"`
	CreatedAt                time.Time `json:"createdAt,omitempty"`
	Id                       string    `json:"id"`
	Name                     string    `json:"name"`
	Status                   string    `json:"status"`
	UpdatedAt                time.Time `json:"updatedAt"`
}

// NetworksResponse represents the response returned by the get networks request
type NetworksResponse struct {
	Networks []*Network `json:"networks"`
}

// Rate represents currency exchange rate information.
type Rate struct {
	Buy       float64   `json:"buy,omitempty"`
	Code      string    `json:"code"`
	Locale    string    `json:"locale"`
	RateId    string    `json:"rateId"`
	Sell      float64   `json:"sell,omitempty"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// RatesResponse represents the response returned by the get rates request
type RatesResponse struct {
	Rates []*Rate `json:"rates"`
}
