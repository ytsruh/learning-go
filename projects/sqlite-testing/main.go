package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite3", "./main.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	http.HandleFunc("/write", writeHandler)
	http.HandleFunc("/read", readHandler)
	http.HandleFunc("/backup", backupHandler)

	fmt.Println("Server is starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func writeHandler(w http.ResponseWriter, r *http.Request) {
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS dummy (data TEXT); INSERT INTO dummy (data) VALUES ('Dummy data')")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, "Dummy data written to DB")
}

func readHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT data FROM dummy")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var data string
		if err := rows.Scan(&data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		fmt.Fprintln(w, "Data from DB:", data)
	}
}

func backupHandler(w http.ResponseWriter, r *http.Request) {
	backupPath := "./backup.db"

	db, err := sql.Open("sqlite3", "./main.db")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Create the backup file
	_, err = db.Exec(fmt.Sprintf("VACUUM INTO '%s';", backupPath))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Open the backup file
	backupFile, err := os.Open(backupPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer backupFile.Close()

	// Set the headers and copy the backup file to the response writer
	w.Header().Set("Content-Disposition", "attachment; filename="+filepath.Base(backupPath))
	w.Header().Set("Content-Type", "application/octet-stream")
	io.Copy(w, backupFile)

	// Delete the backup file
	err = os.Remove(backupPath)
	if err != nil {
		log.Printf("Failed to delete backup file: %v", err)
	}
}
