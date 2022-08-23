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
