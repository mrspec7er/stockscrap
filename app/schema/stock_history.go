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
				
				result := service.GetStockHistory(symbol, fromDate, toDate)
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

var StockHistorySchema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query:    rootQuery,
 })
