package product

import (
	"Intern/gcp_pub-sub/modules/subscriber/model"
	"Intern/gcp_pub-sub/modules/subscriber/repo"
	"errors"
	"github.com/google/uuid"
	"time"
)

// ServiceProduct is the objest that define a business logic of the app
type ServiceProduct struct {
	repo repo.ProductsRepo
}

// New is the constructor of the ServiceProduct entity
func New(r repo.ProductsRepo) *ServiceProduct {
	return &ServiceProduct{
		repo: r,
	}
}

// CheckDate checks if the date is correct in the input and return it in the needed format
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

// isValidUUID check if string has a valid uuid format
func isValidUUID(s string) bool {
	_, err := uuid.Parse(s)
	if err != nil {
		return false
	}

	return true
}

// ActionIDWithInterval prepare data from the input to the needed format and send to the repositoriy
func (s *ServiceProduct) ActionIDWithInterval(input model.InputWithDate) ([]model.DBResponse, error) {
	fromDate, toDate, err := checkDate(input)
	if err != nil {
		return []model.DBResponse{}, err
	}

	if !isValidUUID(input.ActionID) {
		return []model.DBResponse{}, errors.New("invalid input")
	}

	result, err := s.repo.ActionWithInterval(input.ActionID, fromDate, toDate)
	if err != nil {
		return []model.DBResponse{}, err
	}

	return result, nil
}

// TwoActionsIDWithInterval prepare data from the input to the needed format and send to the repositoriy
func (s *ServiceProduct) TwoActionsIDWithInterval(input model.InputWithDate) ([]model.DBResponse2Actions, error) {
	fromDate, toDate, err := checkDate(input)
	if err != nil {
		return []model.DBResponse2Actions{}, err
	}

	if !(isValidUUID(input.ActionID) && isValidUUID(input.ActionID2)) {
		return []model.DBResponse2Actions{}, errors.New("invalid input")
	}

	result, err := s.repo.TwoActionsWithInterval(input.ActionID, input.ActionID2, fromDate, toDate)
	if err != nil {
		return []model.DBResponse2Actions{}, err
	}

	return result, nil
}
