package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"

	handlers "Yamy-Gin/pkg"

)

type Route struct{}

func New() *Route{
	return &Route{}
	
}

func (r *Route) InitRoutes() *gin.Engine {
	db, _ := sql.Open("sqlite3", `c:\Users\Atai\Desktop\newGoLang\database.db`)
	router := gin.New()
	Handler := handlers.NewHandler(db)

	auth := router.Group("/auth")
		{
			auth.POST("/sign-in", Handler.SignIn)
			auth.POST("/sing-up", Handler.SignUp)
		}

	return router
}
