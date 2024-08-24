package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/lib/pq"
)

var db *sql.DB

func connectDB() (*sql.DB, error) {
	connStr := "user=user password=password host=postgres dbname=supply_chain sslmode=disable"
	return sql.Open("postgres", connStr)
}

func main() {
	var err error
	db, err = connectDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	http.HandleFunc("/", homePage)
	http.HandleFunc("/add", addItem)
	http.HandleFunc("/view", viewItems)

	log.Println("Server starting on port 8082...")
	if err := http.ListenAndServe(":8082", nil); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func homePage(w http.ResponseWriter, r *http.Request) {
	statusMessage := r.URL.Query().Get("status")

	tmpl, err := loadTemplate("templates/home.html")
	if err != nil {
		http.Error(w, "Failed to load template: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		StatusMessage string
	}{
		StatusMessage: statusMessage,
	}

	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, "Failed to execute template: "+err.Error(), http.StatusInternalServerError)
	}
}

func addItem(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		err := addItemToDB("Best item", "Description of new item")
		if err != nil {
			http.Redirect(w, r, "/?status=Failed%20to%20add%20item", http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, "/?status=Item%20added%20successfully", http.StatusSeeOther)
	}
}

func addItemToDB(name, description string) error {
	_, err := db.Exec("INSERT INTO items (name, description) VALUES ($1, $2)", name, description)
	return err
}

func viewItems(w http.ResponseWriter, r *http.Request) {
	// Read SQL query from file
	sqlFile := "sql/get_items.sql"
	sqlStmt, err := os.ReadFile(sqlFile)
	if err != nil {
		http.Error(w, "Failed to read SQL file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Execute the SQL query
	rows, err := db.Query(string(sqlStmt))
	if err != nil {
		http.Error(w, "Failed to execute query: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Prepare HTML output
	var result string
	result += "<!DOCTYPE html><html><head><title>View Items</title></head><body><h1>Items</h1>"
	result += "<table border='1'><tr><th>ID</th><th>Name</th><th>Description</th><th>Created At</th></tr>"

	for rows.Next() {
		var id int
		var name, description, createdAt string
		if err := rows.Scan(&id, &name, &description, &createdAt); err != nil {
			http.Error(w, "Failed to scan item: "+err.Error(), http.StatusInternalServerError)
			return
		}
		result += "<tr><td>" + strconv.Itoa(id) + "</td><td>" + name + "</td><td>" + description + "</td><td>" + createdAt + "</td></tr>"
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Error occurred while processing rows: "+err.Error(), http.StatusInternalServerError)
		return
	}

	result += "</table><a href='/'>Back to Home</a></body></html>"
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(result))
}

func loadTemplate(path string) (*template.Template, error) {
	log.Printf("Loading template: %s", path)
	tmplData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return template.New("").Parse(string(tmplData))
}
