package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"

	"github.com/gocolly/colly"
)

type Statistic struct {
	Label string
	Value string
}

type Recommendation struct {
	Title string 
	Body string
}

type StockHistory struct {
	Symbol string `json:"symbol"`
	Date string `json:"date"`
	Open int `json:"open"`
	Close int `json:"close"`
	High int `json:"high"`
	Low int `json:"low"`
	Volume float64 `json:"volume"`
}

type StockHistoryApiResponse struct {
	Status string `json:"status"`
	Message string `json:"message"`
	Data struct {
		Results []*StockHistory `json:"results"`
	} `json:"data"`
}

func DataScrapper()  {	
	c := colly.NewCollector(
		colly.AllowedDomains("id.tradingview.com"),
	)

	statistic := make(map[string]string)

	c.OnHTML(".block-GgmpMpKr", func(e *colly.HTMLElement) {

		label := e.DOM.Find(".label-GgmpMpKr").Text()
		value := e.DOM.Find(".value-GgmpMpKr").Contents().Not(".measureUnit-lQwbiR8R").Text()

		statistic[label] = value
		
	})

	var recommendation []*Recommendation

	c.OnHTML(".card-exterior-Us1ZHpvJ", func(e *colly.HTMLElement) {
		
		title := e.DOM.Find(".title-tkslJwxl").Text()
		body := e.DOM.Find(".line-clamp-content-t3qFZvNN").Text()

		recommendation = append(recommendation, &Recommendation{Title: title, Body: body})
	})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
		fmt.Println("Visiting", r.URL.String())
	})

	c.Visit("https://id.tradingview.com/symbols/IDX-TLKM/")

	for k, s := range statistic {
		fmt.Println(k, ":", s)
	}

	for _, r := range recommendation {
		fmt.Println(*r)
	}

}

func GetStockHistory(symbol string, fromDate string, toDate string) ([]*StockHistory, error) {
	res, err := http.Get("https://api.goapi.io/stock/idx/" + symbol + "/historical?from=" + fromDate + "&to=" + toDate + "&api_key=cd818a59-52d0-51cd-bd66-fa8c6e45")
	fmt.Println("https://api.goapi.io/stock/idx/" + symbol + "/historical?from=" + fromDate + "&to=" + toDate + "&api_key=cd818a59-52d0-51cd-bd66-fa8c6e45")

	if err != nil {
		return nil, err
	}

	streamData, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	convertedData := StockHistoryApiResponse{}
	err = json.Unmarshal(streamData, &convertedData)
	if err != nil {
		return nil, err
	}

	// for _, s := range convertedData.Data.Results {
	// 	fmt.Println(s)
	// }

	result := convertedData.Data.Results

	slices.Reverse(result)

	return result, nil
}