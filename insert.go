package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

// Define a struct to represent the user data
type User struct {
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	AccountType string `json:"accountType"`
	UserProfile string `json:"userProfile"`
}

func main() {
	// Define a handler function to handle the incoming requests
	http.HandleFunc("/insert", func(w http.ResponseWriter, r *http.Request) {
		// Check if the request method is POST
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Decode the JSON data from the "user_data" parameter in the request body
		var user User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error decoding JSON data: %v", err), http.StatusBadRequest)
			return
		}

		// Database connection parameters
		dbUser := "root"
		dbPassword := ""
		dbName := "krutrim"

		// Connect to the MySQL database
		db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@/%s", dbUser, dbPassword, dbName))
		if err != nil {
			http.Error(w, fmt.Sprintf("Error connecting to database: %v", err), http.StatusInternalServerError)
			return
		}
		defer db.Close()

		// Prepare the INSERT statement
		stmt, err := db.Prepare("INSERT INTO new_records (firstname, lastname, email, password, account_type, user_profile) VALUES (?, ?, ?, ?, ?, ?)")
		if err != nil {
			http.Error(w, fmt.Sprintf("Error preparing INSERT statement: %v", err), http.StatusInternalServerError)
			return
		}
		defer stmt.Close()

		// Execute the INSERT statement with the values from the User struct
		_, err = stmt.Exec(user.FirstName, user.LastName, user.Email, user.Password, user.AccountType, user.UserProfile)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error inserting data into database: %v", err), http.StatusInternalServerError)
			return
		}

		// Write a success response
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintln(w, "Data inserted successfully!")
	})

	// Start the HTTP server
	fmt.Println("Server is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}