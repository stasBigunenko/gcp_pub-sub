package postgresql

import (
	"Intern/gcp_pub-sub/modules/subscriber/pkg/app/config"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"sync"
)

type PostgresDB struct {
	Pdb *sql.DB
	mu  sync.Mutex
}

func New(connStr *config.StorageConfiguration) (*PostgresDB, error) {

	db, err := sql.Open("postgres", connStr.ConnString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect database %w\n", err)
	}

	database := &PostgresDB{Pdb: db}

	database.Pdb.Exec(`DROP TABLE products`)
	database.Pdb.Exec(`DROP TABLE categories`)
	database.Pdb.Exec(`DROP TABLE actions`)
	database.Pdb.Exec(`DROP TABLE user_activities`)

	database.Pdb.Exec(`CREATE TABLE products (
    					id CHAR(36) PRIMARY KEY NOT NULL, 
    					name VARCHAR(50) NOT NULL, 
    					description VARCHAR(500) NOT NULL,
    					price FLOAT(8) NOT NULL,
    					category_id CHAR(36) NOT NULL
    				);`)

	database.Pdb.Exec(`CREATE TABLE categories (
    					id CHAR(36) PRIMARY KEY NOT NULL, 
    					name VARCHAR(50) NOT NULL
    				);`)

	database.Pdb.Exec(`CREATE TABLE actions (
    					id CHAR(36) PRIMARY KEY NOT NULL, 
    					name VARCHAR(50) NOT NULL
    				);`)

	database.Pdb.Exec(`CREATE TABLE user_activities (
    					id CHAR(36) PRIMARY KEY NOT NULL, 
    					action_id CHAR(36) NOT NULL,
    					product_id CHAR(36) NOT NULL,
    					created_at TIMESTAMP,
    					FOREIGN KEY (action_id) REFERENCES actions(id),
    					FOREIGN KEY (product_id) REFERENCES products(id)                           
    				);`)

	return database, nil
}

func (db *PostgresDB) AddSomeDataToDB() error {
	_, err := db.Pdb.Exec(`INSERT INTO products (id, name, description, price, category_id) VALUES ($1, $2, $3, $4, $5)`,
		"00000000-0000-0000-0000-000000000001", "Shampoo", "Gel", 100.00, "00000000-0000-0000-1000-000000000000")
	if err != nil {
		return fmt.Errorf("db problems: %v\n", err)
	}
	_, err = db.Pdb.Exec(`INSERT INTO products (id, name, description, price, category_id) VALUES ($1, $2, $3, $4, $5)`,
		"00000000-0000-0000-0000-000000000002", "Soap", "Soft", 130.00, "00000000-0000-0000-1000-000000000000")
	if err != nil {
		return fmt.Errorf("db problems: %v\n", err)
	}
	_, err = db.Pdb.Exec(`INSERT INTO products (id, name, description, price, category_id) VALUES ($1, $2, $3, $4, $5)`,
		"00000000-0000-0000-0000-000000000003", "Toothpaste", "Colgate", 15.50, "00000000-0000-0000-1000-000000000000")
	if err != nil {
		return fmt.Errorf("db problems: %v\n", err)
	}
	_, err = db.Pdb.Exec(`INSERT INTO products (id, name, description, price, category_id) VALUES ($1, $2, $3, $4, $5)`,
		"00000000-0000-0000-0000-000000000004", "Wheel", "Round", 1500.00, "00000000-0000-0000-2000-000000000000")
	if err != nil {
		return fmt.Errorf("db problems: %v\n", err)
	}
	_, err = db.Pdb.Exec(`INSERT INTO products (id, name, description, price, category_id) VALUES ($1, $2, $3, $4, $5)`,
		"00000000-0000-0000-0000-000000000005", "Car", "Big", 10000.00, "00000000-0000-0000-3000-000000000000")
	if err != nil {
		return fmt.Errorf("db problems: %v\n", err)
	}

	_, err = db.Pdb.Exec(`INSERT INTO categories (id, name) VALUES ($1, $2)`,
		"00000000-0000-0000-1000-000000000000", "household chemicals")
	if err != nil {
		return fmt.Errorf("db problems: %v\n", err)
	}
	_, err = db.Pdb.Exec(`INSERT INTO categories (id, name) VALUES ($1, $2)`,
		"00000000-0000-0000-2000-000000000000", "Wheels")
	if err != nil {
		return fmt.Errorf("db problems: %v\n", err)
	}
	_, err = db.Pdb.Exec(`INSERT INTO categories (id, name) VALUES ($1, $2)`,
		"00000000-0000-0000-3000-000000000000", "Cars")
	if err != nil {
		return fmt.Errorf("db problems: %v\n", err)
	}

	_, err = db.Pdb.Exec(`INSERT INTO actions (id, name) VALUES ($1, $2)`,
		"00000000-0000-1000-0000-000000000000", "Put in the bucket")
	if err != nil {
		return fmt.Errorf("db problems: %v\n", err)
	}
	_, err = db.Pdb.Exec(`INSERT INTO actions (id, name) VALUES ($1, $2)`,
		"00000000-0000-2000-0000-000000000000", "Take off from the bucket")
	if err != nil {
		return fmt.Errorf("db problems: %v\n", err)
	}
	_, err = db.Pdb.Exec(`INSERT INTO actions (id, name) VALUES ($1, $2)`,
		"00000000-0000-3000-0000-000000000000", "Looked description")
	if err != nil {
		return fmt.Errorf("db problems: %v\n", err)
	}

	return nil
}
