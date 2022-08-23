package repo

import "Intern/gcp_pub-sub/modules/subscriber/model"

type ProductsRepo interface {
	AddAction(actionId, productId string) error
	ShowAllActions() ([]model.Action, error)
}
