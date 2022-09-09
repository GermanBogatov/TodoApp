package handler

import (
	"github.com/GermanBogatov/TodoApp/app/internal/service"
	"github.com/GermanBogatov/TodoApp/app/pkg/jwt"
	"github.com/GermanBogatov/TodoApp/app/pkg/logging"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	Service *service.Service
	Logger  logging.Logger
	Helper  jwt.Helper
}

func NewHandler(services *service.Service, logger logging.Logger, helper jwt.Helper) (*Handler, error) {
	return &Handler{
		services,
		logger,
		helper,
	}, nil
}

func (h *Handler) InitRoutes() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}
	api := router.Group("/api")
	{
		api.Use(jwt.Middleware())
		lists := api.Group("/lists")
		{
			lists.POST("/", h.createList)
			lists.GET("/", h.getAllList)
			lists.GET("/:id", h.getListById)
			lists.PUT("/:id", h.updateList)
			lists.DELETE("/:id", h.deleteList)

			items := lists.Group(":id/items")
			{
				items.POST("/", h.createItem)
				items.GET("/", h.getAllItem)
			}
		}
		items := api.Group("/items")
		{
			items.GET("/:id", h.getItemById)
			items.PUT("/:id", h.updateItem)
			items.DELETE("/:id", h.deleteItem)
		}
	}
	return router
}
