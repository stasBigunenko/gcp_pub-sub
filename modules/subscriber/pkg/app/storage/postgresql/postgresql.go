package postgresql

import (
	"Intern/gcp_pub-sub/modules/subscriber/pkg/app/config"
	"database/sql"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"io"
	"log"
	"os"
	"sync"
)

// PostgesDB is the object that define db connection
type PostgresDB struct {
	Pdb *sql.DB
	mu  sync.Mutex
}

// New is the constructor that return PostgresDB entity and open db connection
func New(connStr *config.StorageConfiguration) (*PostgresDB, error) {
	db, err := sql.Open("postgres", connStr.ConnString)
	if err != nil {
		return nil, fmt.Errorf("failed to connect database %w\n", err)
	}

	database := &PostgresDB{Pdb: db}

	return database, nil
}

// AddSomeDataIntoTable adding some data to the db in table products
func (db *PostgresDB) AddSomeDataIntoTable(tableName, path string, fieldsQty int) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	if _, err := reader.Read(); err == io.EOF {
		return fmt.Errorf("error: empty file")
	}

	//loop of reading
	for i := 0; ; i++ {
		record, err := reader.Read()
		if err != nil {
			if err != io.EOF {
				return fmt.Errorf("Error reading file: %w", err)
			}
			if i == 0 {
				return fmt.Errorf("Malformed csv file: there's only headers and no values")
			} else {
				log.Print("end of file")
				//end of the file
				return nil
			}
		}

		//malformed csv handling
		if len(record) != fieldsQty {
			return fmt.Errorf("Malformed csv file: wrong number of fields")
		}
		for _, v := range record {
			if v == "" {
				return fmt.Errorf("malformed csv file: empty fields")
			}
		}
		_, err = uuid.Parse(record[0])
		if err != nil {
			return fmt.Errorf("malformed id, should be a uuid: %w", err)
		}

		switch tableName {
		case "products":
			res, err := exists(tableName, record[0], db.Pdb)
			if err != nil {
				return fmt.Errorf("failed to connect database %w\n", err)
			}
			if res {
				break
			}
			_, err = db.Pdb.Exec(`INSERT INTO products (id, name, description, price, category_id) VALUES ($1, $2, $3, $4, $5)`,
				record[0], record[1], record[2], record[3], record[4])
			if err != nil {
				return fmt.Errorf("internal db problems: %w", err)
			}
		case "categories":
			res, err := exists(tableName, record[0], db.Pdb)
			if err != nil {
				return fmt.Errorf("failed to connect database %w\n", err)
			}
			if res {
				break
			}
			_, err = db.Pdb.Exec(`INSERT INTO categories (id, name) VALUES ($1, $2)`,
				record[0], record[1])
			if err != nil {
				return fmt.Errorf("internal db problems: %w", err)
			}
		case "actions":
			res, err := exists(tableName, record[0], db.Pdb)
			if err != nil {
				return fmt.Errorf("failed to connect database %w\n", err)
			}
			if res {
				break
			}
			_, err = db.Pdb.Exec(`INSERT INTO actions (id, name) VALUES ($1, $2)`,
				record[0], record[1])
			if err != nil {
				return fmt.Errorf("internal db problems: %w", err)
			}
		default:
			return fmt.Errorf("cannot find such table")
		}
	}

	return nil
}

func exists(tableName, id string, db *sql.DB) (bool, error) {
	var data string

	err := db.QueryRow(fmt.Sprintf(`SELECT 1 FROM %s WHERE id = $1;`, tableName), id).Scan(&data)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		} else {
			return false, err
		}
	}

	return true, nil
}
