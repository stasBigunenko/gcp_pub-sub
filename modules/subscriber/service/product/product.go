package product

import (
	"Intern/gcp_pub-sub/modules/subscriber/model"
	"Intern/gcp_pub-sub/modules/subscriber/repo"
	"fmt"
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

func (s *ServiceProduct) InBucketInterval(input model.InputWithDate) ([]model.DBResponse, error) {
	fromDate := input.FromDateYear + "-" + input.FromDateMonth + "-" + input.FromDateDay
	checkFromDate, err := time.Parse("2006-01-02", fromDate)
	if err != nil {
		fmt.Errorf("wrong format date: %w\n", err)
		return []model.DBResponse{}, err
	}

	toDate := input.ToDateYear + "-" + input.ToDateMonth + "-" + input.ToDateDay
	checkToDate, err := time.Parse("2006-01-02", toDate)
	if err != nil {
		fmt.Errorf("wrong format date: %w\n", err)
		return []model.DBResponse{}, err
	}

	result, err := s.repo.InBucketsWithInterval(input.ActionID, checkFromDate.Format("2006-01-02"), checkToDate.Format("2006-01-02"))
	if err != nil {
		return []model.DBResponse{}, err
	}

	return result, nil
}
