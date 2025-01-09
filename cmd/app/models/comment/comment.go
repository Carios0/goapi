package comment

import (
	"database/sql"
	"errors"
	"fmt"
)

type Comment struct {
	Id          int    `json:"id"`
	ItemID      int64  `json:"itemid"`
	Description string `json:"description"`
	User        string `json:"user"`
}

func (c *Comment) Delete(db *sql.DB) error {
	result, err := db.Exec("DELETE FROM comments WHERE id=?", c.Id)
	if err != nil {
		return err
	}
	rowsAff, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAff == 0 {
		return errors.New("Tried to delete nonexistent comment")
	}
	return nil
}

func (c *Comment) Save(db *sql.DB) error {
	result, err := db.Exec("INSERT INTO comments (itemId, description, user) VALUES (?,?,?)", c.ItemID, c.Description, c.User)
	if err != nil {
		return err
	}
	rowsAff, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAff == 0 {
		return errors.New("Inserting comment failed.")
	}
	return nil
}

func (c *Comment) Validate(db *sql.DB) error {
	// check if item with id=itemID exists
	var count int
	query := "SELECT COUNT(*) FROM items WHERE id = ?"
	err := db.QueryRow(query, c.ItemID).Scan(&count)
	if err != nil {
		return err
	}
	if count == 0 {
		return errors.New(fmt.Sprintf("Item with Id %d does not exist", c.ItemID))
	}
	// check if User has a name
	if c.User == "" {
		return errors.New("Empty username")
	}
	return nil
}

func GetCommentsByItem(db *sql.DB, itemId int64) ([]Comment, error) {
	rows, err := db.Query("SELECT * FROM comments WHERE itemId=?", itemId)
	if err != nil {
		return nil, errors.New("Error fetching all items")
	}
	defer rows.Close()

	var comms []Comment
	for rows.Next() {
		var comm Comment
		err := rows.Scan(&comm.Id, &comm.ItemID, &comm.Description, &comm.User)
		if err != nil {
			return nil, err
		}
		comms = append(comms, comm)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return comms, nil
}
