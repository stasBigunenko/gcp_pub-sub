package product

import (
	"Intern/gcp_pub-sub/modules/subscriber/model"
	"database/sql"
	"fmt"
	"github.com/google/uuid"
	"sync"
	"time"
)

type Repository struct {
	db *sql.DB
	mu sync.Mutex
}

func New(db *sql.DB) *Repository {
	return &Repository{db: db}
}

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

func (r *Repository) ShowAllActions() ([]model.Action, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	rows, err := r.db.Query(
		`SELECT * FROM user_activities`)
	if err != nil {
		return []model.Action{}, fmt.Errorf("internal error: %v\n", err)
	}
	defer rows.Close()

	actions := []model.Action{}

	for rows.Next() {
		action := model.Action{}
		err = rows.Scan(&action.ID, &action.ActionID, &action.ProductID, &action.CreatedAt)
		if err != nil {
			return []model.Action{}, fmt.Errorf("internal error: %v\n", err)
		}

		actions = append(actions, action)
	}

	return actions, nil
}

func (r *Repository) InBucketsWithInterval(actionID, fromDate, toDate string) ([]model.DBResponse, error) {
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
