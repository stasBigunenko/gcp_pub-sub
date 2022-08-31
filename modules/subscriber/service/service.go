package service

import "Intern/gcp_pub-sub/modules/subscriber/model"

type Service interface {
	ActionIDWithInterval(model.InputWithDate) ([]model.DBResponse, error)
	TwoActionsIDWithInterval(model.InputWithDate) ([]model.DBResponse2Actions, error)
}
