# Go Yellowcard

Submission for Yellow Card Technical Support Engineer role.

## Requirements

- Go 1.19 or later

## Documentation

For a comprehensive guide, checkout the [API documentation](https://docs.yellowcard.engineering/docs/getting-started)

### Available APIs

| API                                                                                       | Description                                                                                |
|-------------------------------------------------------------------------------------------|--------------------------------------------------------------------------------------------|
| [Get Channels](https://docs.yellowcard.engineering/reference/get-channels)                | Retrieve all supported payment ramps (Bank Transfer, Mobile Money, E-Wallets transfers).   |
| [Get Networks](https://docs.yellowcard.engineering/reference/get-networks)                | Retrieve all supported end financial interfaces (Banks, Mobile Money Networks, E-Wallets). |
| [Get Rates](https://docs.yellowcard.engineering/reference/get-rates)                      | Retrieve rates for supported countries.                                                    |
| [Resolve Bank Account](https://sandbox.api.yellowcard.io/business/details/bank)           | Validate a bank account before sending.                                                    |
| [Submit Payment Request](https://sandbox.api.yellowcard.io/business/payments)             | Submit a disbursement payment request. This will lock in a rate and await approval.        |
| [Accept Payment Request](https://sandbox.api.yellowcard.io/business/payments/{id}/accept) | Accept a payment request for execution.                                                    |
| [Deny Payment Request](https://sandbox.api.yellowcard.io/business/payments/{id}/deny)     | Deny a payment request.                                                                    |
| [Lookup Payment](https://sandbox.api.yellowcard.io/business/payments/{id})                | Retrieve information about a specific payment.                                             |

### Usage

The SDK provided optional configuration options that can be tweaked based on preference.
By default, these are the default options

- Base URL is set to [production URL](https://api.yellowcard.io/business)
- Environment is set to [Production Server](https://docs.yellowcard.engineering/docs/environments-api)
- The Http Client is set to ```&http.Client{}```

```go

import (
    yellowcard "github.com/jwambugu/yellowcard-go"
)

client := yellowcard.New("API_KEY", "SECRET_KEY")
```

#### With a http Client

```go

import (
    "net/http"
    yellowcard "github.com/jwambugu/yellowcard-go"
)

httpClient := http.DefaultClient
client := yellowcard.New("API_KEY", "SECRET_KEY", yellowcard.WithHttpClient(httpClient))
```

#### With an environment

When this option is set, the base URL will also be updated to match the correct environment, if an invalid environment
is provided, it defaults to `yellowcard.EnvironmentProduction`

```go

import (
    yellowcard "github.com/jwambugu/yellowcard-go"
)

client := yellowcard.New("API_KEY", "SECRET_KEY", yellowcard.WithEnvironment(yellowcard.EnvironmentSandbox))
```

#### API usage

Some APIs provide a way to filter data based on countries and currency code. Check
the [list of coverage map](https://docs.yellowcard.engineering/docs/coverage-widget)

```go

import (
    "context"
    yellowcard "github.com/jwambugu/yellowcard-go"
)

var(
    client = yellowcard.New("API_KEY", "SECRET_KEY")
    ctx = context.Background()
)

// Get all active channels
channels, err := client.GetChannels(ctx, "")

// Get all active channels in a country
channels, err = client.GetChannels(ctx, yellowcard.CountryCodeKE)

// Get all active networks
networks, err := client.GetNetworks(ctx, "")

// Get all active networks in a country
networks, err := client.GetNetworks(ctx, yellowcard.CountryCodeKE)

// Get all rates
rates, err := client.GetRates(ctx, "")

// Get all rates for a currency
rates, err = client.GetRates(ctx, yellowcard.CurrencyCodeKES)

// Resolve bank account
bankAccountDetails, err := client.ResolveBankAccount(ctx, &yellowcard.ResolveBankAccountRequest{
    AccountNumber: "589000",
    NetworkID:     "41109c18-9604-4389-8472-44ff4378c6cb",
})

// Submit payment request 
paymentRequest := &yellowcard.PaymentRequest{
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

payment, err := client.MakePayment(ctx, paymentRequest, false)

// Skip the accept payment request
payment, err = client.MakePayment(ctx, paymentRequest, true)

// Accept payment request
payment, err = client.AcceptPaymentRequest(ctx, "d83011e8-341f-5e3e-b908-84cb4a552fcc")

// Deny payment request
payment, err := client.DenyPaymentRequest(ctx, "d83011e8-341f-5e3e-b908-84cb4a552fcc")

// Lookup payment
payment, err := client.LookupPayment(ctx, "d83011e8-341f-5e3e-b908-84cb4a552fcc")
```

## Test
The test suite needs testify's `assert` package to run:
```go
github.com/stretchr/testify/assert
```

Before running the tests, make sure to grab all of the package's dependencies:
```shell
go get -t -v
```

Run all tests:
```shell
go test ./... -v
```

