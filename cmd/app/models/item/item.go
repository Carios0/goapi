package item

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"slices"
	"strconv"
	"time"
)

type Item struct {
	Id          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Author      string    `json:"author"`
	Time        time.Time `json:"time"`
	Done        bool      `json:"done"`
	Priority    int       `json:"priority"`
}

func (item *Item) ValidateNewItem() error {
	if item.Title == "" || item.Author == "" {
		return errors.New("Missing value(s): Title and Author may not be empty.")
	}
	return nil
}

func (item *Item) Save(db *sql.DB) (int64, error) {
	query := "INSERT INTO items (title, description, author) VALUES (?,?,?)"
	result, err := db.Exec(query, item.Title, item.Description, item.Author)
	if err != nil {
		return -1, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (item *Item) GetItemById(db *sql.DB, id int64) error {

	// query for and scan item
	err := db.QueryRow("SELECT * FROM items WHERE id = ?", id).Scan(
		&item.Id,
		&item.Title,
		&item.Description,
		&item.Author,
		&item.Time,
		&item.Done,
		&item.Priority)
	if err != nil {
		return err
	}
	return nil
}

func GetAllItems(db *sql.DB) ([]Item, error) {
	rows, err := db.Query("SELECT * FROM items")
	if err != nil {
		return nil, errors.New("Error fetching all items")
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var item Item
		err := rows.Scan(&item.Id, &item.Title, &item.Description, &item.Author, &item.Time, &item.Done, &item.Priority)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func GetAllItemsByAuthor(db *sql.DB, value string) ([]Item, error) {
	rows, err := db.Query("SELECT * FROM items WHERE author=?", value)
	if err != nil {
		return nil, errors.New("Error fetching all items")
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var item Item
		err := rows.Scan(&item.Id, &item.Title, &item.Description, &item.Author, &item.Time, &item.Done, &item.Priority)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

func GetAllItemsByDone(db *sql.DB, value bool) ([]Item, error) {
	rows, err := db.Query("SELECT * FROM items WHERE done=?", value)
	if err != nil {
		return nil, errors.New("Error fetching all items")
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var item Item
		err := rows.Scan(&item.Id, &item.Title, &item.Description, &item.Author, &item.Time, &item.Done, &item.Priority)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
func (item *Item) UpdateItem(db *sql.DB, field string, value string) error {

	log.Printf("Changing item with id=%d to %s=%s", item.Id, field, value)
	// map of allowed fields to update with types
	validFields := map[string]string{
		"title":       "string",
		"description": "string",
		"done":        "bool",
		"priority":    "int",
	}
	dataType, valid := validFields[field] // check if field exists in map
	if !valid {
		return errors.New("Invalid field")
	}
	// convert value into right type and preconstruct query
	var query string
	var args []interface{}
	switch dataType {
	case "int":
		intVal, err := strconv.Atoi(value)
		if err != nil || intVal < 1 || intVal > 4 {
			return errors.New("Invalid Priority")
		}
		query = "UPDATE items SET priority=? WHERE id=?"
		args = []interface{}{intVal, item.Id}
	case "string":
		query = fmt.Sprintf("UPDATE items SET %s=? WHERE id=?", field)
		args = []interface{}{value, item.Id}
	case "bool":
		boolVal, err := strconv.ParseBool(value)
		if err != nil {
			return errors.New("Invalid Done")
		}
		query = "UPDATE items SET done=? WHERE id=?"
		args = []interface{}{boolVal, item.Id}
	default:
		return errors.New("Invalid field")
	}
	// execute query
	_, qerr := db.Exec(query, args...)
	if qerr != nil {
		log.Printf("Update failed: %v", qerr)
		return errors.New("Update failed")
	}
	return nil
}

func (item *Item) DeleteItem(db *sql.DB) error {

	//query delete
	_, qerr := db.Exec("DELETE FROM items WHERE id=?", item.Id)
	if qerr != nil {
		return qerr
	}
	log.Printf("Deleted item with id=%d, title=%s, author=%s created %s\n", item.Id, item.Title, item.Author, item.Time.Format("2006-01-02 15:04:05"))
	return nil
}

func SortByTime(items []Item, inv bool) []Item {
	if !inv {
		slices.SortFunc(items, cmpTime)
	}
	if inv {
		slices.SortFunc(items, func(a, b Item) int { return -cmpTime(a, b) })
	}
	return items
}

func cmpTime(a, b Item) int {
	if a.Time.After(b.Time) {
		return 1
	}
	if a.Time.Before(b.Time) {
		return -1
	}
	return 0
}

func SortByPrio(items []Item, inv bool) []Item {
	if !inv {
		slices.SortFunc(items, func(a, b Item) int { return a.Priority - b.Priority })
	}
	if inv {
		slices.SortFunc(items, func(a, b Item) int { return -(a.Priority - b.Priority) })
	}
	return items
}

func SortByVotes(db *sql.DB, items []Item, inv bool) ([]Item, error) {
	// map itemId to voteCount
	itemVote := make(map[int]int)
	// count votes on all items
	query := "SELECT COUNT(*) FROM votes WHERE itemId=?"
	for _, item := range items {
		var count int
		err := db.QueryRow(query, item.Id).Scan(&count)
		if err != nil {
			return nil, err
		}
		itemVote[item.Id] = count
	}
	// sort items
	if !inv {
		slices.SortFunc(items, func(a, b Item) int { return itemVote[a.Id] - itemVote[b.Id] })
	}
	if inv {
		slices.SortFunc(items, func(a, b Item) int { return -(itemVote[a.Id] - itemVote[b.Id]) })
	}
	return items, nil
}
