package product

import (
	"Intern/gcp_pub-sub/modules/subscriber/service"
)

type Product struct {
	service service.Service
}

func New(service service.Service) *Product {
	return &Product{
		service: service,
	}
}
