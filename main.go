package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sql.DB

type Company struct {
	Cik  string `json:"cik"`
	Name string `json:"name"`
}

func main() {
	// load .env file
	err1 := godotenv.Load(".env")

	if err1 != nil {
		log.Fatalf("Error loading .env file")
	}

	// setup and connect to postgres database
	dbuser := os.Getenv("DBUSER")
	dbpasswd := os.Getenv("DBPASSWORD")
	var err error
	s := fmt.Sprintf("user=%s dbname=goapi_db sslmode=disable password=%s host=localhost", dbuser, dbpasswd)
	db, err = sql.Open("postgres", s)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	if pingErr := db.Ping(); err != nil {
		log.Fatal(pingErr)
	} else {
		log.Println("Successfuly connected")
	}

	// define API paths and handle functions
	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /company/{cik}", deleteCompany)
	mux.HandleFunc("GET /company/{cik}", getCompany)
	mux.HandleFunc("POST /company", createCompany)
	mux.HandleFunc("/", handleRoot)
	mux.HandleFunc("GET /company/all", getAllCompanies)

	fmt.Println("Server at :8080")
	http.ListenAndServe(":8080", mux)
}

// handle company/all to request all existing companies
func getAllCompanies(w http.ResponseWriter, r *http.Request) {
	var result []Company
	var current Company
	rows, err := db.Query("SELECT cik, name FROM company")
	if err != nil {
		log.Println("rows error 1")
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&current.Cik, &current.Name); err != nil {
			log.Println("Scan error")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		result = append(result, current)
	}
	if err := rows.Err(); err != nil {
		log.Println("rows error")
	}
	j, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	log.Println(result)
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

// handle root
func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Starting Page")
}

// handle DELETE requests for companies
func deleteCompany(w http.ResponseWriter, r *http.Request) {
	var c Company
	cik := r.PathValue("cik")
	//if extraction of cik fails
	if cik == "" {
		http.Error(w, "invalid/empty cik", http.StatusBadRequest)
		return
	}
	err := db.QueryRow("DELETE FROM company WHERE cik=$1 RETURNING cik, name", cik).Scan(&c.Cik, &c.Name)
	//if company does not exist in database
	if err != nil {
		http.Error(w, "company not found", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// handle POST requests for companies ---- DATABASE
func createCompany(w http.ResponseWriter, r *http.Request) {
	var comp Company
	err := json.NewDecoder(r.Body).Decode(&comp)
	//if decoding from json request fails
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if comp.Cik == "" {
		http.Error(w, "missing CIK", http.StatusBadRequest)
		return
	}
	if comp.Name == "" {
		http.Error(w, "missing Name", http.StatusBadRequest)
		return
	}
	result, err := db.Exec("INSERT INTO company VALUES ($1,$2)", comp.Cik, comp.Name)
	n, err := result.RowsAffected()
	fmt.Printf("%d companies added! Added %+v", n, comp)
	w.WriteHeader(http.StatusNoContent)
}

// handle GET by cik
func getCompany(w http.ResponseWriter, r *http.Request) {

	var c Company
	cik := r.PathValue("cik")
	//if extraction of cik fails
	if cik == "" {
		http.Error(w, "invalid/empty cik", http.StatusBadRequest)
		return
	}
	row := db.QueryRow("SELECT * FROM company WHERE cik=$1", cik)
	if err := row.Scan(&c.Cik, &c.Name); err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "company not found", http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//if company does not exist in map
	j, err := json.Marshal(c)
	//if conversion to json fails
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// tell client that response is in json
	w.Header().Set("Content-Type", "application/json")
	//return OK and json of comp
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

// write get by name
