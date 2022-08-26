package service

import "Intern/gcp_pub-sub/modules/subscriber/model"

type Service interface {
	InBucketInterval(date model.InputWithDate) ([]model.DBResponse, error)
}
