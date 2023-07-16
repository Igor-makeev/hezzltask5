package handler

import (
	"hezzltask5/internal/service"

	"github.com/gin-gonic/gin"
)

// Стркуктура Обработчика
type Handler struct {
	Router  *gin.Engine
	service *service.Service
}

// Конструктор обработчика
func NewHandler(s *service.Service) *Handler {

	handler := &Handler{
		Router:  gin.New(),
		service: s,
	}

	items := handler.Router.Group("/items")
	{

		items.POST("/create", handler.createHandler)
		items.GET("/list", handler.listHandler)
		items.PATCH("/update", handler.updateHandler)
		items.PATCH("/remove", handler.removeHandler)

	}

	return handler
}
