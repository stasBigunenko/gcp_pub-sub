package repo

import "Intern/gcp_pub-sub/modules/subscriber/model"

type ProductsRepo interface {
	AddAction(string, string) error

	ActionWithInterval(string, string, string) ([]model.DBResponse, error)
	TwoActionsWithInterval(string, string, string, string) ([]model.DBResponse, error)
}
