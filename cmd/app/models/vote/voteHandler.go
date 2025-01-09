package vote

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
)

type VoteHandler struct {
	db *sql.DB
}

func NewVoteHandler(db *sql.DB) *VoteHandler {
	return &VoteHandler{db: db}
}

func (h *VoteHandler) InsertVote(w http.ResponseWriter, r *http.Request) {
	var vote Vote
	// decode json payload into vote
	err := json.NewDecoder(r.Body).Decode(&vote)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// validate itemID and nonempty user field
	if err := vote.Validate(h.db); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// save new vote
	serr := vote.Save(h.db)
	if serr != nil {
		if serr.Error() == "Inserting vote failed." {
			http.Error(w, "insert failed", http.StatusBadRequest)
			return
		}
		http.Error(w, serr.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *VoteHandler) DeleteVoteHF(w http.ResponseWriter, r *http.Request) {
	ItemID, err := strconv.Atoi(r.PathValue("itemId"))
	if err != nil {
		http.Error(w, "invalid ItemId provided", http.StatusBadRequest)
		return
	}
	User := r.URL.Query().Get("user")
	vote := Vote{ItemID: int64(ItemID), User: User}

	// delete vote (and check if existed
	derr := vote.Delete(h.db)
	if derr != nil {
		if derr.Error() == "Tried to delete nonexistent vote" {
			http.Error(w, "Vote does not exist", http.StatusBadRequest)
			return
		}
		http.Error(w, "delete operation failed", http.StatusInternalServerError)
		return
	}
	// return status
	w.WriteHeader(http.StatusNoContent)
}
