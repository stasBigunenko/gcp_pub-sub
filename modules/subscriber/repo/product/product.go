package product

import (
	"Intern/gcp_pub-sub/modules/subscriber/model"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"sync"
	"time"
)

// Repository is object that define storage business logic layer
type Repository struct {
	db *sql.DB
	mu sync.Mutex
}

// New is the constructor of the Repository entity
func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// Addaction insert in the db action that makes user
func (r *Repository) AddAction(actionId, productId string) error {
	id := uuid.New().String()
	createdAt := time.Now()

	_, err := r.db.Exec(`INSERT INTO user_activities (id, action_id, product_id, created_at) VALUES ($1, $2, $3, $4)`,
		id, actionId, productId, createdAt)
	if err != nil {
		return fmt.Errorf("db problems: %v\n", err)
	}

	return nil
}

// ActionWithInterval select from db data according the actionID between interval of dates
func (r *Repository) ActionWithInterval(actionID, fromDate, toDate string) ([]model.DBResponse, error) {
	var results []model.DBResponse

	rows, err := r.db.Query(
		`SELECT products.*, user_activities.action_id, user_activities.created_at FROM products JOIN user_activities on (
    	user_activities.action_id=$1 AND 
    	user_activities.product_id=products.id) WHERE 
    	user_activities.created_at between $2 and $3;`,
		actionID, fromDate, toDate)
	if err != nil {
		return []model.DBResponse{}, err
	}
	defer rows.Close()

	for rows.Next() {
		res := model.DBResponse{}

		err = rows.Scan(&res.ProductID, &res.Name, &res.Description, &res.Price, &res.Category, &res.ActionID, &res.CreatedAt)
		if err != nil {
			return []model.DBResponse{}, err
		}

		results = append(results, res)
	}

	return results, nil
}

// TwoActionsWithInterval select from db data according the two actions between interval of dates
func (r *Repository) TwoActionsWithInterval(actionID, actionID2, fromDate, toDate string) ([]model.DBResponse2Actions, error) {
	var results []model.DBResponse2Actions

	rows, err := r.db.Query(
		`SELECT products.* FROM products JOIN user_activities as ua1 on (
    	ua1.action_id=$1 AND ua1.product_id=products.id)
    	JOIN user_activities as ua2 on (
    	ua2.action_id=$2 AND ua2.product_id=products.id) WHERE 
    	ua1.created_at between $3 and $4;`,
		actionID, actionID2, fromDate, toDate)
	if err != nil {
		return []model.DBResponse2Actions{}, err
	}
	defer rows.Close()

	for rows.Next() {
		res := model.DBResponse2Actions{}

		err = rows.Scan(&res.ProductID, &res.Name, &res.Description, &res.Price, &res.Category)
		if err != nil {
			return []model.DBResponse2Actions{}, err
		}

		results = append(results, res)
	}

	return results, nil
}
