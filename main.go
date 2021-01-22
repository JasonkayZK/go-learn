package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/graphql-go/graphql"

	"github.com/jasonkayzk/go-graphql-api/gql"
	"github.com/jasonkayzk/go-graphql-api/handler"
	"github.com/jasonkayzk/go-graphql-api/mapper"
)

func main() {
	// Initialize our API and return a pointer to our router for http.ListenAndServe
	// and a pointer to our db to defer its closing when main() is finished
	router, db := initializeAPI()
	defer db.Close()

	// Listen on port 4000 and if there's an error log it and exit
	if err := http.ListenAndServe(":4000", router); err != nil {
		log.Fatal(err)
	}
}

func initializeAPI() (*chi.Mux, *mapper.Db) {
	// Create a new connection to our pg database
	db, err := mapper.New(
		mapper.ConnString("127.0.0.1", 3306, "root", "123456", "go_graphql_db"),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Create our root query for graphql
	rootQuery := gql.NewRoot(db)
	// Create a new graphql schema, passing in the the root query
	sc, err := graphql.NewSchema(
		graphql.SchemaConfig{Query: rootQuery.Query},
	)
	if err != nil {
		fmt.Println("Error creating schema: ", err)
	}

	// Create a handler struct that holds a pointer to our database as well
	// as the address of our graphql schema
	s := handler.Handler{
		GqlSchema: &sc,
	}

	// Create a new router
	router := chi.NewRouter()

	// Add some middleware to our router
	router.Use(
		render.SetContentType(render.ContentTypeJSON), // set content-type headers as application/json
		middleware.Logger,          // log API request calls
		middleware.DefaultCompress, // compress results, mostly gzip assets and json
		middleware.StripSlashes,    // match paths with a trailing slash, strip it, and continue routing through the mux
		middleware.Recoverer,       // recover from panics without crashing handler
	)

	// Create the graphql route with a Handler method to handle it
	router.Post("/graphql", s.GraphQL())

	return router, db
}
