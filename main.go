package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
)

type Company struct {
	Cik  string `json:"cik"`
	Name string `json:"name"`
}

// temporary database replacement until database is connected
var companyCache = make(map[int]Company)

var cacheMutex sync.RWMutex

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("DELETE /company/{id}", deleteCompany)
	mux.HandleFunc("GET /company/{id}", getCompany)
	mux.HandleFunc("POST /company", createCompany)
	mux.HandleFunc("/", handleRoot)
	mux.HandleFunc("GET /company/all", getAllCompanies)

	fmt.Println("Server at :8080")
	http.ListenAndServe(":8080", mux)
	companyCache[0] = Company{"1", "ABC"}
}

// handle company/all to request all existing companies
func getAllCompanies(w http.ResponseWriter, r *http.Request) {
	j, err := json.Marshal(companyCache)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(j)
}

// handle root
func handleRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Starting Page")
}

// handle DELETE requests for companies
func deleteCompany(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	//if extraction of id fails
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	//if company does not exist in map
	if _, ok := companyCache[id]; !ok {
		http.Error(w, "company not found", http.StatusBadRequest)
		return
	}
	cacheMutex.Lock()
	delete(companyCache, id)
	cacheMutex.Unlock()

	w.WriteHeader(http.StatusNoContent)
}

// handle POST requests for companies
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
	cacheMutex.Lock()
	companyCache[len(companyCache)] = comp
	cacheMutex.Unlock()
	fmt.Println("company added!")
	w.WriteHeader(http.StatusNoContent)
}

// handle GET requests for companies
func getCompany(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	//if extraction of id fails
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	cacheMutex.RLock()
	comp, ok := companyCache[id]
	cacheMutex.RUnlock()
	//if company does not exist in map
	if !ok {
		http.Error(w, "company not found", http.StatusNotFound)
		return
	}
	j, err := json.Marshal(comp)
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
