package service

import (
	"strconv"
	"time"
)

type QuarterDetail struct {
	Date string
	Price int
	Volume float64
}

type QuarterHistory struct {
	Quarter string
	High QuarterDetail
	Low QuarterDetail
}

type StockQuarterHistories struct {
	AverageResistance int
	AverageSupport    int
	Quarters  []*QuarterHistory
}

func GetQuarterHistories(symbol string, fromYear int) (StockQuarterHistories, error) {
	quarters := []*QuarterHistory{}

	for i := fromYear; i < time.Now().Year(); i++ {

		histories, err := GetStockHistory(symbol, strconv.Itoa(i) + "-01-02", strconv.Itoa(i + 1) + "-01-01")

		if err != nil {
			return StockQuarterHistories{}, err
		}

		Q1Low, Q1High := GetQuarterSupportResistance(histories, 0, 62)
		quarters = append(quarters, &QuarterHistory{Quarter: strconv.Itoa(i) + "-Q1", High: Q1High, Low: Q1Low})

		Q2Low, Q2High := GetQuarterSupportResistance(histories, 63, 124)
		quarters = append(quarters, &QuarterHistory{Quarter: strconv.Itoa(i) + "-Q2", High: Q2High, Low: Q2Low})

		Q3Low, Q3High := GetQuarterSupportResistance(histories, 125, 188)
		quarters = append(quarters, &QuarterHistory{Quarter: strconv.Itoa(i) + "-Q3", High: Q3High, Low: Q3Low})

		Q4Low, Q4High := GetQuarterSupportResistance(histories, 189, len(histories))
		quarters = append(quarters, &QuarterHistory{Quarter: strconv.Itoa(i) + "-Q4", High: Q4High, Low: Q4Low})
		
	}

	averageSupport := 0
	averageResistance := 0
	for _, q := range quarters {
		averageSupport = averageSupport + q.Low.Price
		averageResistance = averageResistance + q.High.Price
	}

	return StockQuarterHistories{AverageResistance: averageResistance / len(quarters), AverageSupport: averageSupport / len(quarters), Quarters: quarters}, nil
}

func GetQuarterSupportResistance(histories []*StockHistory, startRange int, endRange int) (support QuarterDetail, resistance QuarterDetail) {
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

	return QuarterDetail{Price: supportPrice, Date: supportDate, Volume: supportVolume}, QuarterDetail{Price: resistancePrice, Date: resistanceDate, Volume: resistanceVolume}
}