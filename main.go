package main

import (
    "database/sql"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "runtime"

    _ "github.com/lib/pq"
)

// User struct represents a user in the database
type User struct {
    ID        string `json:"id"`
    Username  string `json:"username"`
    Email     string `json:"email"`
    CreatedAt string `json:"created_at"`
    UpdatedAt string `json:"updated_at"`
}

func main() {
    // PostgreSQL connection string (update with your credentials)
    connStr := "host=my-postgres user=farkhan password=h3xpharm dbname=mydb sslmode=disable"

    // Connect to the database
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Test the database connection
    err = db.Ping()
    if err != nil {
        log.Fatal("Cannot connect to the database:", err)
    }
    fmt.Println("Connected to the database!")

    // Root endpoint to show Go version and link to /users API
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        goVersion := runtime.Version()

        // Set the content type to HTML
        w.Header().Set("Content-Type", "text/html")

        // Write HTML response with Go version and a link to the /users API
        fmt.Fprintf(w, `
            <html>
            <head><title>Go API</title></head>
            <body>
            <h1>Welcome to the Go API</h1>
            <p>Go Version: %s</p>
            <p><a href="/users">Click here to view the list of users</a></p>
            </body>
            </html>
        `, goVersion)
    })

    // Endpoint to return the list of users as JSON
    http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
        users, err := getUsers(db)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        // Set the response content type to JSON
        w.Header().Set("Content-Type", "application/json")
        
        // Return the JSON-encoded list of users
        json.NewEncoder(w).Encode(users)
    })

    // Start the server on port 8080
    log.Println("Server starting on port 8080...")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

// Function to query the users table and return a list of users
func getUsers(db *sql.DB) ([]User, error) {
    rows, err := db.Query("SELECT id, username, email, created_at, updated_at FROM users")
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var users []User
    for rows.Next() {
        var user User
        err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt, &user.UpdatedAt)
        if err != nil {
            return nil, err
        }
        users = append(users, user)
    }

    // Check for any errors during row iteration
    if err = rows.Err(); err != nil {
        return nil, err
    }

    return users, nil
}

