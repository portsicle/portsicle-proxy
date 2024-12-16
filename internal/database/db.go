package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

const dbFile = "blocked_domains.db"

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open SQLite database: %v", err)
	}

	createTableQuery := `CREATE TABLE IF NOT EXISTS blocked_domains (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		domain TEXT NOT NULL UNIQUE
	);`

	_, err = db.Exec(createTableQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %v", err)
	}

	return db, nil
}

// adds a domain to the blocked list in the database
func AddBlockedDomain(db *sql.DB, domain string) error {
	insertQuery := "INSERT OR IGNORE INTO blocked_domains (domain) VALUES (?);"
	_, err := db.Exec(insertQuery, domain)
	if err != nil {
		return fmt.Errorf("failed to add domain: %v", err)
	}
	log.Printf("Domain blocked: %s", domain)
	return nil
}

// removes a domain from the blocked list in the database
func RemoveBlockedDomain(db *sql.DB, domain string) error {
	deleteQuery := "DELETE FROM blocked_domains WHERE domain = ?;"
	_, err := db.Exec(deleteQuery, domain)
	if err != nil {
		return fmt.Errorf("failed to remove domain: %v", err)
	}
	log.Printf("Domain unblocked: %s", domain)
	return nil
}

// retrieves all blocked domains from the database
func GetBlockedDomains(db *sql.DB) (map[string]struct{}, error) {
	rows, err := db.Query("SELECT domain FROM blocked_domains;")
	if err != nil {
		return nil, fmt.Errorf("failed to query blocked domains: %v", err)
	}
	defer rows.Close()

	domains := make(map[string]struct{})

	for rows.Next() {
		var domain string
		if err := rows.Scan(&domain); err != nil {
			return nil, fmt.Errorf("failed to fetch a blocked domain: %v", err)
		}
		domains[domain] = struct{}{}
	}

	return domains, nil
}
