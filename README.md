# Currency Converter API

A GoLang wrapper of [Currency Converter API](https://www.currencyconverterapi.com).

You can use [Currency Converter API](https://www.currencyconverterapi.com) to easily query conversion between difference
currencies and check historical rate with a given date range.

## Installation

Use go get

```bash
go get github.com/kitloong/go-currency-converter-api
```

Then import

```go
import "github.com/kitloong/go-currency-converter-api/api"
```

## Usage

Firstly, create an API instance with

| Name      | Description                                                                        |
|-----------|------------------------------------------------------------------------------------|
| `BaseURL` | The API server URL, refer to https://www.currencyconverterapi.com/docs for details |
| `Version` | The API version number, latest is `v7`.                                            |
| `APIKey`  | Your secret API key.                                                               |

```go
currAPI := api.NewAPI(api.Config{
    BaseURL: "https://free.currconv.com",
    Version: "v7",
    APIKey:  "[KEY]",
})
```

Available methods:

- [Convert](#convert)
- [ConvertCompact](#convertcompact)
- [ConvertHistorical](#converthistorical)
- [ConvertHistoricalCompact](#converthistoricalcompact)
- [Currencies](#currencies)
- [Countries](#countries)
- [Usage](#usage)

### `Convert`

Returns the currency conversion rate with `[FROM]_[TO]` request.

To convert currency from `USD` to `MYR`, construct a request struct and set `USD_MYR` to the `Q` field:

```go
convert, err := currAPI.Convert(api.ConvertRequest{
    Q: []string{"USD_MYR"},
})

// convert
// &{
//     Query: {
//         Count: 1
//     }
//     Results: map[
//         "USD_MYR": {
//             ID:  "USD_MYR"
//             Val: 4.348493
//             To:  "MYR"
//             Fr:  "USD"
//         }
//     ]
// }
```

Since `Q` is a string slice, you can append more currencies to request multiple conversion in a single request:

```go
convert, err := currAPI.Convert(api.ConvertRequest{
    Q: []string{"USD_MYR", "MYR_USD"},
})

// convert
// &{
//     Query: {
//         Count: 2
//     }
//     Results: map[
//         "MYR_USD": {
//             ID:  "MYR_USD"
//             Val: 0.229964
//             To:  "USD"
//             Fr:  "MYR"
//         }
//         "USD_MYR": {
//             ID:  "USD_MYR"
//             Val: 4.348493
//             To:  "MYR"
//             Fr:  "USD"
//         }
//     ]
// }
```

### `ConvertCompact`

Returns conversion result with compact mode:

```go
convert, err := currAPI.ConvertCompact(api.ConvertRequest{
    Q: []string{"USD_MYR", "MYR_USD"},
})

// convert
// map[
//     "MYR_USD": 0.229964
//     "USD_MYR": 4.348493
// ]
```

### `ConvertHistorical`

Returns historical currency conversion rate data:

```go
convert, err := currAPI.ConvertHistorical(api.ConvertHistoricalRequest{
    Q:    []string{"USD_MYR", "MYR_USD"},
    Date: time.Date(2023, 2, 14, 0, 0, 0, 0, time.UTC),
})

// convert
// &{
//     Query: {
//         Count: 2
//     }
//     Date: "2023-02-14"
//     Results: map[
//         "MYR_USD": {
//             ID:  "MYR_USD"
//             To:  "USD"
//             Fr:  "MYR"
//             Val: map[
//                 "2023-02-14": 0.229965
//             ]}
//         "USD_MYR": {
//             ID:  "USD_MYR"
//             To:  "MYR"
//             Fr:  "USD"
//             Val: map[
//                 "2023-02-14": 4.348497
//             ]
//         }
//     ]
// }
```

Set `EndDate` to request historical data with date range:

```go
convert, err := currAPI.ConvertHistorical(api.ConvertHistoricalRequest{
    Q:       []string{"USD_MYR", "MYR_USD"},
    Date:    time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
    EndDate: time.Date(2023, 2, 5, 0, 0, 0, 0, time.UTC),
})

// convert
// &{
//     Query:{
//         Count: 2
//     }
//     Date: "2023-02-01"
//     EndDate: "2023-02-05"
//     Results: map[
//         "MYR_USD": {
//             ID:  "MYR_USD"
//             To:  "USD"
//             Fr:  "MYR"
//             Val: map[
//                 "2023-02-01": 0.234411
//                 "2023-02-02": 0.235513
//                 "2023-02-03": 0.23485
//                 "2023-02-04": 0.23485
//                 "2023-02-05": 0.234851
//             ]}
//         "USD_MYR": {
//             ID:  "USD_MYR"
//             To:  "MYR"
//             Fr:  "USD"
//             Val: map[
//                 "2023-02-01": 4.266011
//                 "2023-02-02": 4.246055
//                 "2023-02-03": 4.258039
//                 "2023-02-04": 4.258039
//                 "2023-02-05": 4.258023
//             ]
//         }
//     ]
// }
```

### `ConvertHistoricalCompact`

Returns historical data with compact mode:

```go
convert, err := currAPI.ConvertHistoricalCompact(api.ConvertHistoricalRequest{
    Q:       []string{"USD_MYR", "MYR_USD"},
    Date:    time.Date(2023, 2, 1, 0, 0, 0, 0, time.UTC),
    EndDate: time.Date(2023, 2, 5, 0, 0, 0, 0, time.UTC),
})

// convert
// map[
//     "2023-02-01": 4.266011
//     "2023-02-02": 4.246055
//     "2023-02-03": 4.258039
//     "2023-02-04": 4.258039
//     "2023-02-05": 4.258023
// ]
```

### `Currencies`

Returns a list of currencies:

```go
currencies, err := currAPI.Currencies()

// currencies
// &{
//     Results:map[
//         "MYR": {
//             ID:             "MYR"
//             CurrencyName:   "Malaysian Ringgit"
//             CurrencySymbol: "RM"
//         }
//         "USD": {
//             ID:             "USD"
//             CurrencyName:   "United States Dollar"
//             CurrencySymbol: "$"
//         }
//         ...
//     ]
// }
```

### `Countries`

Returns a list of countries:

```go
countries, err := currAPI.Countries()

// countries
// &{
//     Results: map[
//         "MY": {
//             ID:             "MY"
//             Alpha3:         "MYS"
//             CurrencyID:     "MYR"
//             CurrencyName:   "Malaysian ringgit"
//             CurrencySymbol: "RM"
//             Name:           "Malaysia"
//         }
//         "US": {
//             ID:             "US"
//             Alpha3:         "USA"
//             CurrencyID:     "USD"
//             CurrencyName:   "United States dollar"
//             CurrencySymbol: "$"
//             Name:           "United States of America"
//         }
//         ...
//     ]
// }
```

### `Usage`

Returns your current API usage:

```go
usage, err := currAPI.Usage()

// usage
// &{
//     Timestamp: "2023-02-15 02:05:49.988 +0000 UTC"
//     Usage: 1
// }
```

## License

The project is open-sourced software licensed under the [MIT](LICENSE) license
