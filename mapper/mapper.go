package mapper

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// Db is our database struct used for interacting with the database
type Db struct {
	*sql.DB
}

// User shape
type User struct {
	ID         int
	Name       string
	Age        int
	Profession string
	Friendly   bool
}

// New makes a new database using the connection string and
// returns it, otherwise returns the error
func New(connString string) (*Db, error) {
	db, err := sql.Open("mysql", connString)
	if err != nil {
		return nil, err
	}

	// Check that our connection is good
	if err := db.Ping(); err != nil {
		return nil, err
	}

	fmt.Println("mysql connected!")

	return &Db{db}, nil
}

// ConnString returns a connection string based on the parameters it's given
// This would normally also contain the password, however we're not using one
func ConnString(host string, port int, user, passwd, dbName string) string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s",
		user, passwd, host, port, dbName,
	)
}

// GetUsersByName is called within our user query for graphql
func (d *Db) GetUsersByName(name string) []User {
	// Prepare query, takes a name argument, protects from sql injection
	stmt, err := d.Prepare("SELECT * FROM go_graphql_db.`users` WHERE name=?")
	if err != nil {
		fmt.Println("GetUserByName preparation Err: ", err)
	}

	// Make query with our stmt, passing in name argument
	rows, err := stmt.Query(name)
	if err != nil {
		fmt.Println("GetUserByName Query Err: ", err)
	}

	// Create User struct for holding each row's data
	var r User
	// Create slice of Users for our response
	var users []User
	// Copy the columns from row into the values pointed at by r (User)
	for rows.Next() {
		if err := rows.Scan(
			&r.ID,
			&r.Name,
			&r.Age,
			&r.Profession,
			&r.Friendly,
		); err != nil {
			fmt.Println("Error scanning rows: ", err)
		}
		users = append(users, r)
	}

	return users
}
