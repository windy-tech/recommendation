package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	language "cloud.google.com/go/language/apiv1"
	"github.com/windy-tech/recommendation/pkg/text"
)

type Book struct {
	ID         int64   `json:"id,omitempty"`
	Name       string  `json:"name,omitempty"`
	Category   string  `json:"category,omitempty"`
	Content    string  `json:"content,omitempty"`
	Confidence float32 `json:"confidence,omitempty"`
}

func GetAllBooks(db *sql.DB, _ *language.Client, w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT * FROM books")
	if err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	books := make([]Book, 0)
	for rows.Next() {
		var book Book
		err = rows.Scan(&book.ID, &book.Name, &book.Category, &book.Confidence, &book.Content)
		if err != nil {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}
		books = append(books, book)
	}
	respondJSON(w, 200, books)
}

func getCatagory(content string, lang *language.Client) (string, error) {
	resp, err := text.ClassifyText(context.Background(), lang, content)
	if err != nil {
		return "", err
	}
	return resp.Categories[0].Name, nil
}

func CreateBook(db *sql.DB, lang *language.Client, w http.ResponseWriter, r *http.Request) {
	book := Book{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&book); err != nil {
		log.Println("Failed to decode book")
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()
	cata, err := getCatagory(book.Content, lang)
	if err != nil {
		log.Println("Failed to use ML")
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	book.Category = cata
	stmt, err := db.Prepare("INSERT INTO books(name , catagory , confidence, content) values(?,?,?,?)")
	if err != nil {
		log.Println("Failed to prepare insterting  to database")
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	res, err := stmt.Exec(book.Name, book.Category, book.Confidence, book.Content)
	if err != nil {
		log.Println("Failed to instert to database")
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	book.ID, _ = res.LastInsertId()
	respondJSON(w, http.StatusCreated, book)
}

// respondJSON makes the response with payload as json format
func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(response))
}

// respondError makes the error response with payload as json format
func respondError(w http.ResponseWriter, code int, message string) {
	respondJSON(w, code, map[string]string{"error": message})
}
