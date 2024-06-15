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
	ID                      string    `json:"id"`
	Max                     int       `json:"max"`
	Min                     float64   `json:"min"`
	RampType                string    `json:"rampType"`
	SettlementType          string    `json:"settlementType"`
	Status                  string    `json:"status"`
	SuccessThreshold        int       `json:"successThreshold,omitempty"`
	UpdatedAt               time.Time `json:"updatedAt"`
	VendorID                string    `json:"vendorId"`
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
	ChannelIDs               []string  `json:"channelIds"`
	Code                     string    `json:"code,omitempty"`
	Country                  string    `json:"country"`
	CountryAccountNumberType string    `json:"countryAccountNumberType"`
	CreatedAt                time.Time `json:"createdAt,omitempty"`
	ID                       string    `json:"id"`
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
	RateID    string    `json:"rateId"`
	Sell      float64   `json:"sell,omitempty"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// RatesResponse represents the response returned by the get rates request
type RatesResponse struct {
	Rates []*Rate `json:"rates"`
}

// ResolveBankAccountRequest validates a bank account before sending the money.
type ResolveBankAccountRequest struct {
	AccountNumber string `json:"accountNumber"`
	NetworkID     string `json:"networkId"`
}

type ResolveBankAccountResponse struct {
	AccountBank   string `json:"accountBank"`
	AccountName   string `json:"accountName"`
	AccountNumber string `json:"accountNumber"`
}

type Destination struct {
	AccountBank   string `json:"accountBank"`
	AccountName   string `json:"accountName"`
	AccountNumber string `json:"accountNumber"`
	AccountType   string `json:"accountType"`
	Country       string `json:"country"`
	NetworkID     string `json:"networkId"`
	NetworkName   string `json:"networkName"`
}

type Sender struct {
	Address  string `json:"address"`
	Country  string `json:"country"`
	Dob      string `json:"dob"`
	Email    string `json:"email"`
	IdNumber string `json:"idNumber"`
	IDType   string `json:"idType"`
	Name     string `json:"name"`
	Phone    string `json:"phone"`
}

type PaymentRequest struct {
	Amount      float64     `json:"amount"`
	ChannelID   string      `json:"channelId"`
	Destination Destination `json:"destination"`
	Reason      string      `json:"reason"`
	Sender      Sender      `json:"sender"`
	SequenceID  string      `json:"sequenceId"`
}

type Payment struct {
	Amount                float64     `json:"amount"`
	ChannelID             string      `json:"channelId"`
	ConvertedAmount       float64     `json:"convertedAmount"`
	Country               string      `json:"country"`
	CreatedAt             time.Time   `json:"createdAt"`
	Currency              string      `json:"currency"`
	Destination           Destination `json:"destination"`
	DirectSettlement      bool        `json:"directSettlement"`
	ExpiresAt             time.Time   `json:"expiresAt"`
	ForceAccept           bool        `json:"forceAccept"`
	ID                    string      `json:"id"`
	PartnerID             string      `json:"partnerId"`
	Rate                  float64     `json:"rate"`
	Reason                string      `json:"reason"`
	RequestSource         string      `json:"requestSource"`
	Sender                Sender      `json:"sender"`
	SequenceID            string      `json:"sequenceId"`
	ServiceFeeAmountLocal float64     `json:"serviceFeeAmountLocal"`
	ServiceFeeAmountUSD   float64     `json:"serviceFeeAmountUSD"`
	SettlementInfo        any         `json:"settlementInfo"`
	Status                string      `json:"status"`
	UpdatedAt             time.Time   `json:"updatedAt"`
}
