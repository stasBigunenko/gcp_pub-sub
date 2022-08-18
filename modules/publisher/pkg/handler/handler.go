package handler

import "github.com/gin-gonic/gin"

type Handler interface {
	Index(c *gin.Context)
	SendData(c *gin.Context)
}
