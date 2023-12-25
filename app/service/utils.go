package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"slices"

	"github.com/gocolly/colly"
	"github.com/mrspec7er/stockscrap/app/dto"
	"github.com/mrspec7er/stockscrap/app/repository"
)

type UtilService struct {
	Redis repository.Redis
}

func (UtilService) DataScrapper()  {	
	c := colly.NewCollector(
		colly.AllowedDomains("id.tradingview.com"),
	)

	statistic := make(map[string]string)

	c.OnHTML(".block-GgmpMpKr", func(e *colly.HTMLElement) {

		label := e.DOM.Find(".label-GgmpMpKr").Text()
		value := e.DOM.Find(".value-GgmpMpKr").Contents().Not(".measureUnit-lQwbiR8R").Text()

		statistic[label] = value
		
	})

	var recommendation []*dto.Recommendation

	c.OnHTML(".card-exterior-Us1ZHpvJ", func(e *colly.HTMLElement) {
		
		title := e.DOM.Find(".title-tkslJwxl").Text()
		body := e.DOM.Find(".line-clamp-content-t3qFZvNN").Text()

		recommendation = append(recommendation, &dto.Recommendation{Title: title, Body: body})
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

func (u UtilService) GetStockHistory(symbol string, fromDate string, toDate string) ([]*dto.StockHistory, error) {
	key := symbol + "-" + fromDate + "-" + toDate
	result := []*dto.StockHistory{}
	result, err := u.Redis.Retrieve(key)

	if err != nil {	
		result, err = u.GoApiGetHistories(symbol, fromDate, toDate)

		if len(result) == 0 || err != nil {
			return nil, errors.New("Invalid symbol or date type")
		}

		u.Redis.CacheHistory(key, result)
	}

	slices.Reverse(result)

	return result, nil
}

func (u UtilService) GoApiGetHistories(symbol string, fromDate string, toDate string) ([]*dto.StockHistory, error) {
	res, err := http.Get("https://api.goapi.io/stock/idx/" + symbol + "/historical?from=" + fromDate + "&to=" + toDate + "&api_key=cd818a59-52d0-51cd-bd66-fa8c6e45")
	fmt.Println("https://api.goapi.io/stock/idx/" + symbol + "/historical?from=" + fromDate + "&to=" + toDate + "&api_key=cd818a59-52d0-51cd-bd66-fa8c6e45")

	if err != nil {
		return nil, err
	}

	streamData, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	convertedData := dto.StockHistoryApiResponse{}
	err = json.Unmarshal(streamData, &convertedData)
	if err != nil {
		return nil, err
	}

	return convertedData.Data.Results, nil
}