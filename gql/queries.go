package gql

import (
	"github.com/graphql-go/graphql"
	"github.com/jasonkayzk/go-graphql-api/mapper"
)

// Root holds a pointer to a graphql object
type Root struct {
	Query *graphql.Object
}

// User describes a graphql object containing a User
var User = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"age": &graphql.Field{
				Type: graphql.Int,
			},
			"profession": &graphql.Field{
				Type: graphql.String,
			},
			"friendly": &graphql.Field{
				Type: graphql.Boolean,
			},
		},
	},
)

// NewRoot returns base query type. This is where we add all the base queries
func NewRoot(db *mapper.Db) *Root {
	resolver := Resolver{db: db}

	// Create a new Root that describes our base query set up. In this
	// example we have a user query that takes one argument called name
	root := Root{
		Query: graphql.NewObject(
			graphql.ObjectConfig{
				Name: "Query",
				Fields: graphql.Fields{
					"users": &graphql.Field{
						Type: graphql.NewList(User),
						Args: graphql.FieldConfigArgument{
							"name": &graphql.ArgumentConfig{
								Type: graphql.String,
							},
						},
						// Create a resolver holding our database. Resolver can be found in resolvers.go
						Resolve: resolver.UserResolver,
					},
				},
			},
		),
	}
	return &root
}
