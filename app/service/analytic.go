package service

import (
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/mrspec7er/stockscrap/app/dto"
)

type AnalyticService struct {
	Utils UtilService
}

func (s AnalyticService) GetQuarterAnalytic(symbol string, fromYear int) (dto.StockQuarterHistories, error) {
	quarters := []*dto.QuarterHistory{}

	ctx := make(chan []*dto.QuarterHistory, (time.Now().Year())-fromYear+1)
	wg := &sync.WaitGroup{}
	wg.Add(time.Now().Year() - fromYear + 1)

	for i := fromYear; i <= time.Now().Year(); i++ {

		go s.GetQuarterHistories(symbol, i, ctx, wg)

	}

	wg.Wait()
	close(ctx)

	for quarter := range ctx {
		quarters = append(quarters, quarter...)
	}

	averageSupport := 0
	averageResistance := 0
	for _, q := range quarters {
		averageSupport = averageSupport + q.Low.Price
		averageResistance = averageResistance + q.High.Price
	}

	return dto.StockQuarterHistories{AverageResistance: averageResistance / len(quarters), AverageSupport: averageSupport / len(quarters), Quarters: quarters}, nil
}

func (s AnalyticService) GetQuarterHistories(symbol string, year int, ctx chan []*dto.QuarterHistory, wg *sync.WaitGroup) {
	quarters := []*dto.QuarterHistory{}
	histories, err := s.Utils.GetStockHistory(symbol, strconv.Itoa(year)+"-01-02", strconv.Itoa(year+1)+"-01-01")

	dataLength := len(histories)

	if err != nil {
		log.Println(err.Error())
		wg.Done()
		return
	}

	if dataLength < 61 {
		wg.Done()
		return
	}
	Q1Low, Q1High := s.GetQuarterSupportResistance(histories, 0, 62)
	quarters = append(quarters, &dto.QuarterHistory{Quarter: strconv.Itoa(year) + "-Q1", High: Q1High, Low: Q1Low})

	if dataLength < 124 {
		wg.Done()
		return
	}
	Q2Low, Q2High := s.GetQuarterSupportResistance(histories, 63, 124)
	quarters = append(quarters, &dto.QuarterHistory{Quarter: strconv.Itoa(year) + "-Q2", High: Q2High, Low: Q2Low})

	if dataLength < 188 {
		wg.Done()
		return
	}
	Q3Low, Q3High := s.GetQuarterSupportResistance(histories, 125, 188)
	quarters = append(quarters, &dto.QuarterHistory{Quarter: strconv.Itoa(year) + "-Q3", High: Q3High, Low: Q3Low})

	if dataLength < 240 {
		wg.Done()
		return
	}
	Q4Low, Q4High := s.GetQuarterSupportResistance(histories, 189, len(histories))
	quarters = append(quarters, &dto.QuarterHistory{Quarter: strconv.Itoa(year) + "-Q4", High: Q4High, Low: Q4Low})

	ctx <- quarters
	wg.Done()
	return
}

func (s AnalyticService) GetQuarterSupportResistance(histories []*dto.StockHistory, startRange int, endRange int) (support dto.QuarterDetail, resistance dto.QuarterDetail) {
	supportPrice := 9999999999
	supportDate := "2000-01-02"
	supportVolume := float64(0)

	resistancePrice := 0
	resistanceDate := "2000-01-02"
	resistanceVolume := float64(0)

	for _, h := range histories[startRange:endRange] {
		if h.Close < supportPrice {
			supportPrice = h.Close
			supportDate = h.Date
			supportVolume = h.Volume
		}

		if h.Close > resistancePrice {
			resistancePrice = h.Close
			resistanceDate = h.Date
			resistanceVolume = h.Volume
		}
	}

	return dto.QuarterDetail{Price: supportPrice, Date: supportDate, Volume: supportVolume}, dto.QuarterDetail{Price: resistancePrice, Date: resistanceDate, Volume: resistanceVolume}
}
