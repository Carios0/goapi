package main

import (
	"check24/cmd/app/data"
	"check24/cmd/app/models/comment"
	"check24/cmd/app/models/item"
	"check24/cmd/app/models/vote"
	"crypto/sha256"
	"crypto/subtle"
	"database/sql"
	"log"
	"net/http"
)

type application struct {
	auth struct {
		username string
		password string
	}
}

func main() {
	// initialize db and handler
	var db *sql.DB
	db = data.ConnectDB()
	// init Handlers
	ih := item.NewItemHandler(db)
	vh := vote.NewVoteHandler(db)
	ch := comment.NewCommentHandler(db)
	log.Println("Handler created")
	// Create tables
	queries := []string{`
	CREATE TABLE IF NOT EXISTS items(
            id INT AUTO_INCREMENT PRIMARY KEY,
            title VARCHAR(100),
            description TEXT,
	    author VARCHAR(100),
	    time DATETIME DEFAULT CURRENT_TIMESTAMP,
	    done BOOLEAN DEFAULT FALSE,
	    priority INT DEFAULT 2
	);
	`,
		`
	CREATE TABLE IF NOT EXISTS votes(
	    itemId INT,
            user VARCHAR(100),
            PRIMARY KEY (itemId, user),
            FOREIGN KEY (itemId) REFERENCES items(id) ON DELETE CASCADE
	);
	`,
		`
	CREATE TABLE IF NOT EXISTS comments(
	    id INT AUTO_INCREMENT PRIMARY KEY,
	    itemId INT,
	    description TEXT,
	    user VARCHAR(100),
	    FOREIGN KEY (itemId) REFERENCES items(id) ON DELETE CASCADE
	);
	`}
	for _, query := range queries {
		_, err := db.Exec(query)
		if err != nil {
			log.Fatalf("Error creating tables: %v", err)
		}
	}
	start(ih, vh, ch)
	data.CloseDB(db)
}

// start server and call handler functions
func start(ih *item.ItemHandler, vh *vote.VoteHandler, ch *comment.CommentHandler) {
	// set acceptable username and password
	app := new(application)
	app.auth.username = "leo"
	app.auth.password = "123"
	// route requests for ItemHandler
	mux := http.NewServeMux()
	mux.HandleFunc("/", ih.HelloWorld) // no authentication required to visit home
	mux.HandleFunc("POST /items", app.BasicAuth(ih.InsertItem))
	mux.HandleFunc("GET /items/all", app.BasicAuth(ih.GetAllItemsHF))
	mux.HandleFunc("GET /items/all/author/{author}", app.BasicAuth(ih.GetAllItemsByAuthorHF))
	mux.HandleFunc("GET /items/all/done/{done}", app.BasicAuth(ih.GetAllItemsByDoneHF))
	mux.HandleFunc("GET /items/all/sort/priority/{inv}", app.BasicAuth(ih.SortByPrioHF))
	mux.HandleFunc("GET /items/all/sort/time/{inv}", app.BasicAuth(ih.SortByTimeHF))
	mux.HandleFunc("GET /items/all/sort/votes/{inv}", app.BasicAuth(ih.SortByVotesHF))
	mux.HandleFunc("GET /items/{id}", app.BasicAuth(ih.GetItemById))
	mux.HandleFunc("PATCH /items/{id}", app.BasicAuth(ih.UpdateItem))
	mux.HandleFunc("DELETE /items/{id}", app.BasicAuth(ih.DeleteItemHF))
	// route requests for VoteHandler
	mux.HandleFunc("POST /votes", app.BasicAuth(vh.InsertVote))
	mux.HandleFunc("DELETE /votes/{itemId}", app.BasicAuth(vh.DeleteVoteHF))
	// route requests for commentHandler
	mux.HandleFunc("POST /comments", app.BasicAuth(ch.InsertComment))
	mux.HandleFunc("DELETE /comments/{id}", app.BasicAuth(ch.DeleteCommentHF))
	mux.HandleFunc("GET /comments/{itemId}", app.BasicAuth(ch.GetCommentsByItem))

	log.Println("Starting server on 8080")
	http.ListenAndServe(":8080", mux)
}

func (app *application) BasicAuth(nextFunc http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if ok {
			// hashing all credentials to hide length
			usernameHash := sha256.Sum256([]byte(username))
			passwordHash := sha256.Sum256([]byte(password))
			expectedUsernameHash := sha256.Sum256([]byte(app.auth.username))
			expectedPasswordHash := sha256.Sum256([]byte(app.auth.password))

			nameMatch := (subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1)
			pwMatch := (subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1)

			if nameMatch && pwMatch {
				nextFunc.ServeHTTP(w, r)
				return
			}
		}
		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Authentication failed", http.StatusUnauthorized)
	})
}
