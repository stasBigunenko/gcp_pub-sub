package repo

import "Intern/gcp_pub-sub/modules/subscriber/model"

type ProductsRepo interface {
	AddAction(string, string) error
	ShowAllActions() ([]model.Action, error)

	InBucketsWithInterval(string, string, string) ([]model.DBResponse, error)
}
