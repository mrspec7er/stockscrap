package schema

import (
	"github.com/graphql-go/graphql"
	"github.com/mrspec7er/stockscrap/app/service"
)

var rootQuery = graphql.NewObject(graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{
		"stockHistories": &graphql.Field{
			Type: graphql.NewList(stockHistoryType),
			Args: graphql.FieldConfigArgument{
				"symbol": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
						
				},
				"fromDate": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"toDate": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				symbol := p.Args["symbol"].(string)
				fromDate := p.Args["fromDate"].(string)
				toDate := p.Args["toDate"].(string)
				
				result, err := service.GetStockHistory(symbol, fromDate, toDate)
				if err != nil {
					return nil, err
				}
				return result, nil
			},
		},

		"quarterHistories": &graphql.Field{
			Type: stockQuarterHistoriesType,
			Args: graphql.FieldConfigArgument{
				"symbol": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
						
				},
				"fromYear": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				symbol := p.Args["symbol"].(string)
				fromYear := p.Args["fromYear"].(int)
				
				result, err := service.GetQuarterHistories(symbol, fromYear)
				
				if err != nil {
					return nil, err
				}
				return result, nil
			},
		},
	},
})

var stockHistoryType = graphql.NewObject(graphql.ObjectConfig{
	Name: "StockHistory",
	Fields: graphql.Fields{
		"symbol": &graphql.Field{
			Type: graphql.String,
		},
		"date": &graphql.Field{
			Type: graphql.String,
		},
		"open": &graphql.Field{
			Type: graphql.Int,
		},
		"close": &graphql.Field{
			Type: graphql.Int,
		},
		"high": &graphql.Field{
			Type: graphql.Int,
		},
		"low": &graphql.Field{
			Type: graphql.Int,
		},
		"volume": &graphql.Field{
			Type: graphql.Float,
		},
	},
})

var stockQuarterDetailType = graphql.NewObject(graphql.ObjectConfig{
	Name: "StockQuarterDetail",
	Fields: graphql.Fields{
		"date": &graphql.Field{
			Type: graphql.String,
		},
		"price": &graphql.Field{
			Type: graphql.Int,
		},
		"volume": &graphql.Field{
			Type: graphql.Float,
		},
	},
})

var stockQuarterType = graphql.NewObject(graphql.ObjectConfig{
	Name: "StockQuarter",
	Fields: graphql.Fields{
		"quarter": &graphql.Field{
			Type: graphql.String,
		},
		"high": &graphql.Field{
			Type: stockQuarterDetailType,
		},
		"low": &graphql.Field{
			Type: stockQuarterDetailType,
		},
	},
})

var stockQuarterHistoriesType = graphql.NewObject(graphql.ObjectConfig{
	Name: "StockQuarterHistories",
	Fields: graphql.Fields{
		"averageResistance": &graphql.Field{
			Type: graphql.Int,
		},
		"averageSupport": &graphql.Field{
			Type: graphql.Int,
		},
		"quarters": &graphql.Field{
			Type: graphql.NewList(stockQuarterType),
		},
	},
})

var StockTechnicalAnalysisSchema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query:    rootQuery,
})
