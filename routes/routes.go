package routes

import (
	"database/sql"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"

	handlers "Yamy-Gin/pkg"
	"Yamy-Gin/pkg/middleware"
)

type Route struct{}

func New() *Route {
	return &Route{}

}

func (r *Route) InitRoutes(db *sql.DB) *gin.Engine {
	router := gin.New()
	router.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:5173"},
		AllowMethods: []string{"GET", "POST", "DELETE"},
		AllowCredentials: true,
	}))


	Handler := handlers.NewHandler(db)

	auth := router.Group("/auth")
	{
		auth.POST("/sign-in", Handler.SignIn)
		auth.POST("/sign-up", Handler.SignUp)
	}
	middlewareFunc := middleware.AuthMiddleware()
	entries := router.Group("/entries")
	entries.Use(middlewareFunc)
	
	{
		
		entries.POST("/", Handler.EntryCreate)
		entries.GET("/", Handler.GetAllById)
		entries.GET("/:id", Handler.GetById)
		entries.DELETE("/:id", Handler.DeleteByID)
	}

	return router
}
