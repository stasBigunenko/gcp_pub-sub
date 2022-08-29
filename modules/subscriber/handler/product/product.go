package product

import (
	"Intern/gcp_pub-sub/modules/subscriber/model"
	"Intern/gcp_pub-sub/modules/subscriber/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Product struct {
	service service.Service
}

func New(service service.Service) *Product {
	return &Product{
		service: service,
	}
}

func (p *Product) ProductsInBucket(c *gin.Context) {
	var input model.InputWithDate

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	products, err := p.service.ActionIDWithInterval(input)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"products": products})
}

func (p *Product) ProductsOutFromBucket(c *gin.Context) {
	var input model.InputWithDate

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	products, err := p.service.ActionIDWithInterval(input)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"products": products})
}

func (p *Product) ProductsDescription(c *gin.Context) {
	var input model.InputWithDate

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	products, err := p.service.ActionIDWithInterval(input)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"products": products})
}

func (p *Product) ProductsBucketAndDescription(c *gin.Context) {
	var input model.InputWithDate

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	products, err := p.service.ActionIDWithInterval(input)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"products": products})
}
