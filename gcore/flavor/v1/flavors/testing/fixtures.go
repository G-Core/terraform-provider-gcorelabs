package testing

import (
	"bitbucket.gcore.lu/gcloud/gcorecloud-go"
	"bitbucket.gcore.lu/gcloud/gcorecloud-go/gcore/flavor/v1/flavors"

	"github.com/shopspring/decimal"
)

const ListResponse = `
{
  "count": 1,
  "results": [
    {
      "flavor_id": "g1-cpu-16-16",
      "flavor_name": "g1-cpu-16-16",
      "price_status": "show",
      "currency_code": "USD",
      "price_per_hour": 0.42,
      "price_per_month": 303.6,
      "ram": 16384,
      "vcpus": 16
    }
  ]
}
`

const ListResponseMalformedCurrency = `
{
  "count": 1,
  "results": [
    {
      "flavor_id": "g1-cpu-16-16",
      "flavor_name": "g1-cpu-16-16",
      "price_status": "show",
      "currency_code": "XXXXXX",
      "price_per_hour": 0.42,
      "price_per_month": 303.6,
      "ram": 16384,
      "vcpus": 16
    }
  ]
}
`

const ListSliceResponse = `
[
	{
	  "flavor_id": "g1-cpu-16-16",
	  "flavor_name": "g1-cpu-16-16",
	  "price_status": "show",
	  "currency_code": "USD",
	  "price_per_hour": 0.42,
	  "price_per_month": 303.6,
	  "ram": 16384,
	  "vcpus": 16
}
]
`

var (
	priceStatus     = "show"
	pricePerHour    = decimal.NewFromFloat(0.42)
	pricePerMonth   = decimal.NewFromFloat(303.6)
	currencyCode, _ = gcorecloud.ParseCurrency("USD")
	flavor          = "g1-cpu-16-16"
	Flavor1         = flavors.Flavor{
		FlavorID:      flavor,
		FlavorName:    flavor,
		PriceStatus:   &priceStatus,
		CurrencyCode:  currencyCode,
		PricePerHour:  &pricePerHour,
		PricePerMonth: &pricePerMonth,
		RAM:           16384,
		VCPUS:         16,
	}

	ExpectedFlavorSlice = []flavors.Flavor{Flavor1}
)
