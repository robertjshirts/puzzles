package dal

import (
	"database/sql"
	"fmt"
	"log"
)

func main() {
	// Connect to the database
	db, err := sql.Open("postgres", "user=myuser password=mypassword dbname=mydb sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Ping the database to verify connection
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to the database!")

	// Execute a query
	rows, err := db.Query("SELECT id, name FROM users")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	// Iterate over the results
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("ID: %d, Name: %s\n", id, name)
	}

	// Check for errors after iterating
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}
