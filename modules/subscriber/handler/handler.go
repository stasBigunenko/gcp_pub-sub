package handler

import "github.com/gin-gonic/gin"

type Handler interface {
	ProductsInBucket(*gin.Context)
}
