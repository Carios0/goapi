package vote

import (
	"database/sql"
	"errors"
	"fmt"
)

type Vote struct {
	ItemID int64  `json:"itemid"`
	User   string `json:"user"`
}

func (v *Vote) Save(db *sql.DB) error {
	result, err := db.Exec("INSERT INTO votes VALUES (?,?)", v.ItemID, v.User)
	if err != nil {
		return err
	}
	rowsAff, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAff == 0 {
		return errors.New("Inserting vote failed.")
	}
	return nil
}

func (v *Vote) Validate(db *sql.DB) error {
	// check if item with id=itemID exists
	var count int
	query := "SELECT COUNT(*) FROM items WHERE id = ?"
	err := db.QueryRow(query, v.ItemID).Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New(fmt.Sprintf("Item with Id %d does not exist", v.ItemID))
	}
	// check if User has a name
	if v.User == "" {
		return errors.New("Empty username")
	}
	return nil
}

func (v *Vote) Delete(db *sql.DB) error {
	result, err := db.Exec("DELETE FROM votes WHERE itemId=? AND user=?", v.ItemID, v.User)
	if err != nil {
		return err
	}
	rowsAff, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAff == 0 {
		return errors.New("Tried to delete nonexistent vote")
	}
	return nil
}
