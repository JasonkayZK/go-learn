package gql

import (
	"github.com/graphql-go/graphql"
	"github.com/jasonkayzk/go-graphql-api/mapper"
)

// Resolver struct holds a connection to our database
type Resolver struct {
	db *mapper.Db
}

// UserResolver resolves our user query through a db call to GetUserByName
func (r *Resolver) UserResolver(p graphql.ResolveParams) (interface{}, error) {
	// Strip the name from arguments and assert that it's a string
	name, ok := p.Args["name"].(string)
	if ok {
		users := r.db.GetUsersByName(name)
		return users, nil
	}

	return nil, nil
}
