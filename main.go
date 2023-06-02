package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type Person struct {
	ID   int
	Name string
	City string
}

var db *sql.DB

func main() {
	var err error
	// Create a new database connection
	db, err = sql.Open("mysql", "root:localhost@tcp(localhost:3306)/mydb")
	if err != nil {
		panic(err)
	}

	// Create a new router
	http.HandleFunc("/people", GetPeopleHandler)
	http.HandleFunc("/people/new", CreatePersonHandler)
	http.HandleFunc("/people/update", UpdatePersonHandler)
	http.HandleFunc("/people/delete", DeletePersonHandler)

	// Start the web server
	http.ListenAndServe(":8000", nil)
}

func GetPeopleHandler(w http.ResponseWriter, r *http.Request) {
	// Get all people from the database.
	rows, err := db.Query("SELECT * FROM people")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	// Iterate over the rows and create a new Person object for each one
	people := []Person{}
	for rows.Next() {
		var person Person
		err := rows.Scan(&person.ID, &person.Name, &person.City)
		if err != nil {
			panic(err)
		}
		people = append(people, person)
	}

	// Convert the people to a JSON string
	jsonString, err := json.Marshal(people)
	if err != nil {
		panic(err)
	}

	// Write the people to the response
	fmt.Fprintf(w, "%s", jsonString)
}

func CreatePersonHandler(w http.ResponseWriter, r *http.Request) {
	// Get the name and city from the request body
	name := r.FormValue("name")
	city := r.FormValue("city")

	if name == "" {
		name = "Akash"
	}
	if city == "" {
		city = "Haridwar"
	}

	// Insert a new person into the database
	stmt, err := db.Prepare("INSERT INTO people (name, city) VALUES (?, ?)")
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(name, city)
	if err != nil {
		panic(err)
	}

	// Redirect the user to the list of people
	http.Redirect(w, r, "/people", http.StatusFound)
}

func UpdatePersonHandler(w http.ResponseWriter, r *http.Request) {
	// Get the person ID from the request URI
	personID := r.URL.Query().Get("id")

	// Get the name and city from the request body
	name := r.FormValue("name")
	city := r.FormValue("city")

	// Update the person in the database
	stmt, err := db.Prepare("UPDATE people SET name = ?, city = ? WHERE id = ?")
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(name, city, personID)
	if err != nil {
		panic(err)
	}

	// Redirect the user to the list of people
	http.Redirect(w, r, "/people", http.StatusFound)
}

func DeletePersonHandler(w http.ResponseWriter, r *http.Request) {
	// Get the person ID from the request URI
	personID := r.URL.Query().Get("id")

	// Delete the person from the database.
	stmt, err := db.Prepare("DELETE FROM people WHERE id = ?")
	if err != nil {
		panic(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(personID)
	if err != nil {
		if err == sql.ErrNoRows {
			// The person does not exist in the database.
			fmt.Fprintf(w, "The person does not exist.")
		} else {
			// There was an error deleting the person.
			panic(err)
		}
	}
	http.Redirect(w, r, "/people", http.StatusFound)
}
