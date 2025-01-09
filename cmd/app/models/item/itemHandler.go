package item

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
)

type ItemHandler struct {
	db *sql.DB
}

func NewItemHandler(db *sql.DB) *ItemHandler {
	return &ItemHandler{db: db}
}

func (h *ItemHandler) HelloWorld(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Starting Page, Welcome!"))
}

func (h *ItemHandler) InsertItem(w http.ResponseWriter, r *http.Request) {
	var item Item
	// decode json payload into item
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// dont allow empty title or author fields
	if err := item.ValidateNewItem(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// save new item and get id of new item
	id, err := item.Save(h.db)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "insert failed", http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// check if INSERT was successful
	itErr := item.GetItemById(h.db, id)
	if itErr != nil {
		if itErr == sql.ErrNoRows {
			http.Error(w, "Failed to insert new item", http.StatusBadRequest)
			return
		}
		http.Error(w, "Failed to identify new Item", http.StatusInternalServerError)
		return
	}

	log.Printf("New item added! Added %+v", item)
	w.WriteHeader(http.StatusNoContent)
}

func (h *ItemHandler) GetItemById(w http.ResponseWriter, r *http.Request) {
	var item Item
	//check valid id and scan item to return
	id, idErr := h.GetId(w, r)
	if idErr != nil {
		return
	}
	itErr := item.GetItemById(h.db, id)
	if itErr != nil {
		if itErr == sql.ErrNoRows {
			http.Error(w, "Item does not exist", http.StatusBadRequest)
			return
		}
		http.Error(w, "Failed to identify Item", http.StatusInternalServerError)
		return
	}

	// convert to json
	j, err := json.Marshal(item)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// tell client that response is in json and return payload
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
}

func (h *ItemHandler) GetAllItemsHF(w http.ResponseWriter, r *http.Request) {

	// query for all items
	result, err := GetAllItems(h.db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//turning result into json format
	j, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// tell client that response is in json and return
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
	log.Printf("Returned all %d items", len(result))
}

func (h *ItemHandler) GetAllItemsByAuthorHF(w http.ResponseWriter, r *http.Request) {
	//get author
	author := r.PathValue("author")
	if author == "" {
		http.Error(w, "Invalid Author input.", http.StatusBadRequest)
		return
	}
	// query for all items by author
	result, err := GetAllItemsByAuthor(h.db, author)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//turning result into json format
	j, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// tell client that response is in json and return
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
	log.Printf("Returned all %d items", len(result))
}

func (h *ItemHandler) GetAllItemsByDoneHF(w http.ResponseWriter, r *http.Request) {
	// get done value
	done, err := strconv.ParseBool(r.PathValue("done"))
	if err != nil {
		http.Error(w, "Invalid input for done. Boolean required", http.StatusBadRequest)
		return
	}
	// query for all items that are done/not done
	result, err := GetAllItemsByDone(h.db, done)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//turning result into json format
	j, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// tell client that response is in json and return
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
	log.Printf("Returned all %d items", len(result))
}

func (h *ItemHandler) SortByPrioHF(w http.ResponseWriter, r *http.Request) {
	// get inv to see if sorting should be inverted
	inv, err := strconv.ParseBool(r.PathValue("inv"))
	if err != nil {
		http.Error(w, "Invalid input for inv. Boolean required", http.StatusBadRequest)
		return
	}
	// query for all items
	result, err := GetAllItems(h.db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// sort result
	result = SortByPrio(result, inv)
	//turning result into json format
	j, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// tell client that response is in json and return
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
	log.Printf("Returned all %d items", len(result))
}

func (h *ItemHandler) SortByVotesHF(w http.ResponseWriter, r *http.Request) {
	// get inv to see if sorting should be inverted
	inv, err := strconv.ParseBool(r.PathValue("inv"))
	if err != nil {
		http.Error(w, "Invalid input for inv. Boolean required", http.StatusBadRequest)
		return
	}
	// query for all items
	items, err := GetAllItems(h.db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// sort result
	result, err := SortByVotes(h.db, items, inv)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	//turning result into json format
	j, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// tell client that response is in json and return
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
	log.Printf("Returned all %d items", len(result))
}

func (h *ItemHandler) SortByTimeHF(w http.ResponseWriter, r *http.Request) {

	// get inv to see if sorting should be inverted
	inv, err := strconv.ParseBool(r.PathValue("inv"))
	if err != nil {
		http.Error(w, "Invalid input for inv. Boolean required", http.StatusBadRequest)
		return
	}
	// query for all items
	result, err := GetAllItems(h.db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// sort result
	result = SortByTime(result, inv)
	//turning result into json format
	j, err := json.Marshal(result)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// tell client that response is in json and return
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)
	log.Printf("Returned all %d items", len(result))
}

func (h *ItemHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	var item Item
	// retrieve and check valid id
	id, idErr := h.GetId(w, r)
	if idErr != nil {
		return
	}
	itErr := item.GetItemById(h.db, id)
	if itErr != nil {
		if itErr == sql.ErrNoRows {
			http.Error(w, "Item to update does not exist", http.StatusBadRequest)
			return
		}
		http.Error(w, "Failed to identify Item for update", http.StatusInternalServerError)
		return
	}

	field := r.URL.Query().Get("field")
	value := r.URL.Query().Get("value")

	err := item.UpdateItem(h.db, field, value)
	if err != nil {
		if err.Error() == "Update failed" {
			http.Error(w, "Invalid field specified", http.StatusBadRequest)
			return
		} else if err.Error() == "Invalid Priority" {
			http.Error(w, "Invalid priority specified", http.StatusBadRequest)
			return
		} else if err.Error() == "Ivalid Done" {
			http.Error(w, "Invalid done value specified", http.StatusBadRequest)
			return
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	// return status
	log.Println("Update successful!")
	w.WriteHeader(http.StatusOK)
}

func (h *ItemHandler) DeleteItemHF(w http.ResponseWriter, r *http.Request) {
	var item Item
	//check valid id and scan item to delete
	id, idErr := h.GetId(w, r)
	if idErr != nil {
		return
	}
	itErr := item.GetItemById(h.db, id)
	if itErr != nil {
		if itErr == sql.ErrNoRows {
			http.Error(w, "Item to delete does not exist", http.StatusBadRequest)
			return
		}
		http.Error(w, "Failed to identify Item for deletion", http.StatusInternalServerError)
		return
	}
	// delete item
	err := item.DeleteItem(h.db)
	if err != nil {
		http.Error(w, "delete operation failed", http.StatusInternalServerError)
		return
	}
	// return status
	w.WriteHeader(http.StatusNoContent)
}

func (h *ItemHandler) GetId(w http.ResponseWriter, r *http.Request) (int64, error) {
	// parse id and convert to int
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "Missing/invalid id", http.StatusBadRequest)
		return -1, errors.New("Failed extracting Id")
	}
	return int64(id), nil
}
