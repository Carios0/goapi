package comment

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type CommentHandler struct {
	db *sql.DB
}

func NewCommentHandler(db *sql.DB) *CommentHandler {
	return &CommentHandler{db: db}
}

func (h *CommentHandler) InsertComment(w http.ResponseWriter, r *http.Request) {
	var comm Comment
	// decode json payload into vote
	err := json.NewDecoder(r.Body).Decode(&comm)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// validate itemID and nonempty user field
	if err := comm.Validate(h.db); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// save new vote
	serr := comm.Save(h.db)
	if serr != nil {
		if serr.Error() == "Inserting comment failed." {
			http.Error(w, "insert failed", http.StatusBadRequest)
			return
		}
		http.Error(w, serr.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *CommentHandler) DeleteCommentHF(w http.ResponseWriter, r *http.Request) {
	var comm Comment
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid comment id provided", http.StatusBadRequest)
		return
	}
	comm.Id = id

	// delete comment (and check if existed)
	derr := comm.Delete(h.db)
	if derr != nil {
		if derr.Error() == "Tried to delete nonexistent comment" {
			http.Error(w, "Comment does not exist", http.StatusBadRequest)
			return
		}
		http.Error(w, "delete operation failed", http.StatusInternalServerError)
		return
	}
	// return status
	w.WriteHeader(http.StatusNoContent)
}

func (h *CommentHandler) GetCommentsByItem(w http.ResponseWriter, r *http.Request) {
	//get itemID and confirm item exists (in Validate)
	var comm Comment
	itemId, err := strconv.Atoi(r.PathValue("itemId"))
	if err != nil {
		http.Error(w, "invalid item id provided", http.StatusBadRequest)
		return
	}
	comm.ItemID = int64(itemId)
	comm.User = "placeholderName" // else the validation would fail because of empty name
	if err := comm.Validate(h.db); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// query for all comments on this item
	result, err := GetCommentsByItem(h.db, int64(itemId))
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
