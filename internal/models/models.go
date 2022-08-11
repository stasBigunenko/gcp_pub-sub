package models

type Message struct {
	CategoryID string `json:"categoryID,omitempty"`
	ProductID  string `json:"productID,omitempty"`
	ActionID   string `json:"actionID,omitempty"`
}
