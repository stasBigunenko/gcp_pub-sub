package product

import (
	"Intern/gcp_pub-sub/modules/subscriber/model"
	"Intern/gcp_pub-sub/modules/subscriber/repo"
	"time"
)

type ServiceProduct struct {
	repo repo.ProductsRepo
}

func New(r repo.ProductsRepo) *ServiceProduct {
	return &ServiceProduct{
		repo: r,
	}
}

func checkDate(input model.InputWithDate) (string, string, error) {
	fromDate := input.FromDateYear + "-" + input.FromDateMonth + "-" + input.FromDateDay
	checkFromDate, err := time.Parse("2006-01-02", fromDate)
	if err != nil {
		return "", "", err
	}

	toDate := input.ToDateYear + "-" + input.ToDateMonth + "-" + input.ToDateDay
	checkToDate, err := time.Parse("2006-01-02", toDate)
	if err != nil {
		return "", "", err
	}

	return checkFromDate.Format("2006-01-02"), checkToDate.Format("2006-01-02"), nil
}

func (s *ServiceProduct) ActionIDWithInterval(input model.InputWithDate) ([]model.DBResponse, error) {
	fromDate, toDate, err := checkDate(input)
	if err != nil {
		return []model.DBResponse{}, err
	}

	result, err := s.repo.ActionWithInterval(input.ActionID, fromDate, toDate)
	if err != nil {
		return []model.DBResponse{}, err
	}

	return result, nil
}

func (s *ServiceProduct) TwoActionsIDWithInterval(input model.InputWithDate) ([]model.DBResponse2Actions, error) {
	fromDate, toDate, err := checkDate(input)
	if err != nil {
		return []model.DBResponse2Actions{}, err
	}

	result, err := s.repo.TwoActionsWithInterval(input.ActionID, input.ActionID2, fromDate, toDate)
	if err != nil {
		return []model.DBResponse2Actions{}, err
	}

	return result, nil
}
