package dto

type QuarterDetail struct {
	Date   string
	Price  int
	Volume float64
}

type QuarterHistory struct {
	Quarter string
	High    QuarterDetail
	Low     QuarterDetail
}

type StockQuarterHistories struct {
	AverageResistance int
	AverageSupport    int
	Quarters          []*QuarterHistory
}

type Statistic struct {
	Label string
	Value string
}

type Recommendation struct {
	Title string
	Body  string
}

type StockHistory struct {
	Symbol string  `json:"symbol"`
	Date   string  `json:"date"`
	Open   int     `json:"open"`
	Close  int     `json:"close"`
	High   int     `json:"high"`
	Low    int     `json:"low"`
	Volume float64 `json:"volume"`
}

type StockHistoryApiResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    struct {
		Results []*StockHistory `json:"results"`
	} `json:"data"`
}
