package handler

import "github.com/gin-gonic/gin"

type Handler interface {
	ProductsInBucket(*gin.Context)
	ProductsOutFromBucket(*gin.Context)
	ProductsDescription(*gin.Context)
	ProductsBucketAndDescription(*gin.Context)
}
