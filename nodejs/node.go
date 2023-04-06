package nodejs

import (
	"bytes"
	"log"
	"fmt"
	"os"
	"os/exec"
	"encoding/json"
)

func RunNode() {
	cmd := exec.Command("node", "./nodejs/yfinance.js", "AAPL")

	// Capture standard output and standard error separately
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		log.Printf("Error running nodejs script: %s\n", err)
		os.Exit(1)
	}
	// Print the command output and error for debugging
	fmt.Printf("Standard output: %s\n", stdout.String())
	fmt.Printf("Standard error: %s\n", stderr.String())
}

func RunNodeData() {
	cmd := exec.Command("node", "./nodejs/yfinance.js", "AAPL")

	output, err := cmd.Output()
	if err != nil {
		log.Printf("Error running nodejs script: %s\n", err)
		os.Exit(1)
	}

	// Unmarshal the output into a struct
	var result Stock
	err = json.Unmarshal(output, &result)
	if err != nil {
		fmt.Printf("Error unmarshaling JSON: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Short Name: %s, Symbol: %s, Price: %f\n", result.ShortName, result.Symbol, result.RegularMarketPrice)
}

type Stock struct {
	Language                   string    `json:"language"`
	Region                     string    `json:"region"`
	QuoteType                  string    `json:"quoteType"`
	TypeDisp                   string    `json:"typeDisp"`
	QuoteSourceName            string    `json:"quoteSourceName"`
	Triggerable                bool      `json:"triggerable"`
	CustomPriceAlertConfidence string    `json:"customPriceAlertConfidence"`
	Currency                   string    `json:"currency"`
	Exchange                   string    `json:"exchange"`
	ShortName                  string    `json:"shortName"`
	LongName                   string    `json:"longName"`
	MessageBoardId             string    `json:"messageBoardId"`
	ExchangeTimezoneName       string    `json:"exchangeTimezoneName"`
	ExchangeTimezoneShortName  string    `json:"exchangeTimezoneShortName"`
	GmtOffSetMilliseconds      int       `json:"gmtOffSetMilliseconds"`
	Market                     string    `json:"market"`
	EsgPopulated               bool      `json:"esgPopulated"`
	MarketState                string    `json:"marketState"`
	RegularMarketChangePercent float64   `json:"regularMarketChangePercent"`
	RegularMarketPrice         float64   `json:"regularMarketPrice"`
	EpsCurrentYear             float64   `json:"epsCurrentYear"`
	PriceEpsCurrentYear        float64   `json:"priceEpsCurrentYear"`
	SharesOutstanding          int64     `json:"sharesOutstanding"`
	BookValue                  float64   `json:"bookValue"`
	FiftyDayAverage            float64   `json:"fiftyDayAverage"`
	FiftyDayAverageChange      float64   `json:"fiftyDayAverageChange"`
	FiftyDayAverageChangePercent float64 `json:"fiftyDayAverageChangePercent"`
	TwoHundredDayAverage       float64   `json:"twoHundredDayAverage"`
	TwoHundredDayAverageChange float64   `json:"twoHundredDayAverageChange"`
	MarketCap                  int64     `json:"marketCap"`
	ForwardPE                  float64   `json:"forwardPE"`
	PriceToBook                float64   `json:"priceToBook"`
	SourceInterval             int       `json:"sourceInterval"`
	ExchangeDataDelayedBy      int       `json:"exchangeDataDelayedBy"`
	AverageAnalystRating       string    `json:"averageAnalystRating"`
	Tradeable                  bool      `json:"tradeable"`
	CryptoTradeable            bool      `json:"cryptoTradeable"`
	FirstTradeDateMilliseconds string    `json:"firstTradeDateMilliseconds"`
	PriceHint                  int       `json:"priceHint"`
	PreMarketChange            float64   `json:"preMarketChange"`
	PreMarketChangePercent     float64   `json:"preMarketChangePercent"`
	PreMarketTime              string    `json:"preMarketTime"`
	PreMarketPrice             float64   `json:"preMarketPrice"`
	RegularMarketChange        float64   `json:"regularMarketChange"`
	RegularMarketTime          string    `json:"regularMarketTime"`
	RegularMarketDayHigh       float64   `json:"regularMarketDayHigh"`
	RegularMarketDayRange      struct {
		Low  float64 `json:"low"`
		High float64 `json:"high"`
	} `json:"regularMarketDayRange"`
	RegularMarketDayLow          float64 `json:"regularMarketDayLow"`
	RegularMarketVolume          int64   `json:"regularMarketVolume"`
	RegularMarketPreviousClose   float64 `json:"regularMarketPreviousClose"`
	Bid                          float64 `json:"bid"`
	Ask                          float64 `json:"ask"`
	BidSize                      int     `json:"bidSize"`
	AskSize                      int     `json:"askSize"`
	FullExchangeName             string  `json:"fullExchangeName"`
	FinancialCurrency            string  `json:"financialCurrency"`
	RegularMarketOpen            float64 `json:"regularMarketOpen"`
	AverageDailyVolume3Month     int64   `json:"averageDailyVolume3Month"`
	AverageDailyVolume10Day      int64   `json:"averageDailyVolume10Day"`
	FiftyTwoWeekLowChange        float64 `json:"fiftyTwoWeekLowChange"`
	FiftyTwoWeekLowChangePercent float64 `json:"fiftyTwoWeekLowChangePercent"`
	FiftyTwoWeekRange            struct {
		Low  float64 `json:"low"`
		High float64 `json:"high"`
	} `json:"fiftyTwoWeekRange"`
	FiftyTwoWeekHighChange        float64 `json:"fiftyTwoWeekHighChange"`
	FiftyTwoWeekHighChangePercent float64 `json:"fiftyTwoWeekHighChangePercent"`
	FiftyTwoWeekLow               float64 `json:"fiftyTwoWeekLow"`
	FiftyTwoWeekHigh              float64 `json:"fiftyTwoWeekHigh"`
	DividendDate                  string  `json:"dividendDate"`
	EarningsTimestamp             string  `json:"earningsTimestamp"`
	EarningsTimestampStart        string  `json:"earningsTimestampStart"`
	EarningsTimestampEnd          string  `json:"earningsTimestampEnd"`
	TrailingAnnualDividendRate    float64 `json:"trailingAnnualDividendRate"`
	TrailingPE                    float64 `json:"trailingPE"`
	TrailingAnnualDividendYield   float64 `json:"trailingAnnualDividendYield"`
	EpsTrailingTwelveMonths       float64 `json:"epsTrailingTwelveMonths"`
	EpsForward                    float64 `json:"epsForward"`
	TwoHundredDayAverageChangePercent float64 `json:"twoHundredDayAverageChangePercent"`
	DisplayName                  string  `json:"displayName"`
	Symbol                       string  `json:"symbol"`
}
