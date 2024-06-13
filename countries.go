package yellowcard

import "errors"

// A CountryCode is an ISO 3166-2 short alphanumeric identification code for countries.
type CountryCode string

const (
	CountryCodeBF CountryCode = "BF"
	CountryCodeBJ CountryCode = "BJ"
	CountryCodeBW CountryCode = "BW"
	CountryCodeCD CountryCode = "CD"
	CountryCodeCG CountryCode = "CG"
	CountryCodeCI CountryCode = "CI"
	CountryCodeCM CountryCode = "CM"
	CountryCodeGA CountryCode = "GA"
	CountryCodeGH CountryCode = "GH"
	CountryCodeKE CountryCode = "KE"
	CountryCodeML CountryCode = "ML"
	CountryCodeMW CountryCode = "MW"
	CountryCodeNG CountryCode = "NG"
	CountryCodeRW CountryCode = "RW"
	CountryCodeSN CountryCode = "SN"
	CountryCodeTG CountryCode = "TG"
	CountryCodeTZ CountryCode = "TZ"
	CountryCodeUG CountryCode = "UG"
	CountryCodeZA CountryCode = "ZA"
	CountryCodeZM CountryCode = "ZM"
)

// CurrencyCode represents ISO 4217 currency code.
type CurrencyCode string

func (c CurrencyCode) String() string {
	return string(c)
}

const (
	CurrencyCodeBWP CurrencyCode = "BWP"
	CurrencyCodeCDF CurrencyCode = "CDF"
	CurrencyCodeGHS CurrencyCode = "GHS"
	CurrencyCodeKES CurrencyCode = "KES"
	CurrencyCodeMWK CurrencyCode = "MWK"
	CurrencyCodeNGN CurrencyCode = "NGN"
	CurrencyCodeRWF CurrencyCode = "RWF"
	CurrencyCodeTZS CurrencyCode = "TZS"
	CurrencyCodeUGX CurrencyCode = "UGX"
	CurrencyCodeXAF CurrencyCode = "XAF"
	CurrencyCodeXOF CurrencyCode = "XOF"
	CurrencyCodeZAR CurrencyCode = "ZAR"
	CurrencyCodeZMW CurrencyCode = "ZMW"
)

// CurrencyCodes is a list of currently supported countries currency codes
var CurrencyCodes = map[CurrencyCode]struct{}{
	CurrencyCodeBWP: {},
	CurrencyCodeCDF: {},
	CurrencyCodeGHS: {},
	CurrencyCodeKES: {},
	CurrencyCodeMWK: {},
	CurrencyCodeNGN: {},
	CurrencyCodeRWF: {},
	CurrencyCodeTZS: {},
	CurrencyCodeUGX: {},
	CurrencyCodeXAF: {},
	CurrencyCodeXOF: {},
	CurrencyCodeZAR: {},
	CurrencyCodeZMW: {},
}

// Country holds information about a country.
type Country struct {
	Code         string
	Name         string
	CurrencyCode CurrencyCode
}

func (c CountryCode) String() string {
	return string(c)
}

// ErrCountryNotSupported is returned when the provided country is not supported.
var ErrCountryNotSupported = errors.New("yellowcard: country is not supported")

// CountryCodes is a list of currently supported countries. See https://docs.yellowcard.engineering/docs/coverage-api
var CountryCodes = map[CountryCode]Country{
	CountryCodeBF: {Code: "BF", CurrencyCode: CurrencyCodeXOF, Name: "Burkina Faso"},
	CountryCodeBJ: {Code: "BJ", CurrencyCode: CurrencyCodeXOF, Name: "Benin"},
	CountryCodeBW: {Code: "BW", CurrencyCode: CurrencyCodeBWP, Name: "Botswana"},
	CountryCodeCD: {Code: "CD", CurrencyCode: CurrencyCodeCDF, Name: "Democratic Republic of the Congo"},
	CountryCodeCG: {Code: "CG", CurrencyCode: CurrencyCodeXAF, Name: "Congo Brazzaville"},
	CountryCodeCI: {Code: "CI", CurrencyCode: CurrencyCodeXOF, Name: "Ivory Coast"},
	CountryCodeCM: {Code: "CM", CurrencyCode: CurrencyCodeXAF, Name: "Cameroon"},
	CountryCodeGA: {Code: "GA", CurrencyCode: CurrencyCodeXAF, Name: "Gabon"},
	CountryCodeGH: {Code: "GH", CurrencyCode: CurrencyCodeGHS, Name: "Ghana"},
	CountryCodeKE: {Code: "KE", CurrencyCode: CurrencyCodeKES, Name: "Kenya"},
	CountryCodeML: {Code: "ML", CurrencyCode: CurrencyCodeXOF, Name: "Mali"},
	CountryCodeMW: {Code: "MW", CurrencyCode: CurrencyCodeMWK, Name: "Malawi"},
	CountryCodeNG: {Code: "NG", CurrencyCode: CurrencyCodeNGN, Name: "Nigeria"},
	CountryCodeRW: {Code: "RW", CurrencyCode: CurrencyCodeRWF, Name: "Rwanda"},
	CountryCodeSN: {Code: "SN", CurrencyCode: CurrencyCodeXOF, Name: "Senegal"},
	CountryCodeTG: {Code: "TG", CurrencyCode: CurrencyCodeXOF, Name: "Togo"},
	CountryCodeTZ: {Code: "TZ", CurrencyCode: CurrencyCodeTZS, Name: "Tanzania"},
	CountryCodeUG: {Code: "UG", CurrencyCode: CurrencyCodeUGX, Name: "Uganda"},
	CountryCodeZA: {Code: "ZA", CurrencyCode: CurrencyCodeZAR, Name: "South Africa"},
	CountryCodeZM: {Code: "ZM", CurrencyCode: CurrencyCodeZMW, Name: "Zambia"},
}
