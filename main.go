package main

import (
	"net/http"

	"github.com/graphql-go/handler"
	"github.com/mrspec7er/stockscrap/app/schema"
)

func main()  {
	h := handler.New(&handler.Config{
		Schema: &schema.StockHistorySchema,
		Pretty: true,
		GraphiQL: false,
	})

	http.Handle("/graphql", h)

	http.ListenAndServe(":8080", nil)
}