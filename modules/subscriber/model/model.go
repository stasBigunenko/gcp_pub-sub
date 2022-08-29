package model

import "time"

type Message struct {
	ProductID string `json:"productID,omitempty"`
	ActionID  string `json:"actionID,omitempty"`
}

type Action struct {
	ID        string    `json:"id,omitempty"`
	ActionID  string    `json:"actionID"`
	ProductID string    `json:"productID"`
	CreatedAt time.Time `json:"createdAt"`
}

type DBResponse struct {
	ActionID    string    `json:"actionID"`
	CreatedAt   time.Time `json:"createdAt"`
	ProductID   string    `json:"productID"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float32   `json:"price"`
	Category    string    `json:"category"`
}

type InputWithDate struct {
	ActionID      string `json:"actionID"`
	ActionID2     string `json:"actionID2"`
	ToDateYear    string `json:"toDateYear"`
	ToDateMonth   string `json:"toDateMonth"`
	ToDateDay     string `json:"toDateDay"`
	FromDateYear  string `json:"fromDateYear"`
	FromDateMonth string `json:"fromDateMonth"`
	FromDateDay   string `json:"fromDateDay"`
}
