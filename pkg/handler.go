package handlers

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"

	"Yamy-Gin/pkg/auth"
)

type Handler struct{
	database *sql.DB
}

func NewHandler(database *sql.DB) *Handler{
	return &Handler{database: database}
}

func (a *Handler) SignUp(ctx *gin.Context){
	var signUp auth.InputSignUp

	if err := ctx.ShouldBindJSON(&signUp); err != nil{
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	errSign := signUp.SignUp(a.database)

	if errSign != nil{
		ctx.JSON(400, gin.H{"error": errSign.Error()})
		return
	}

	ctx.JSON(200, gin.H{"status": "succesfully created account"})
}

func (a *Handler) SignIn(ctx *gin.Context){
	var signIn auth.InputSignIn

	if err := ctx.ShouldBindJSON(&signIn); err != nil{
		ctx.JSON(401, gin.H{"error": err.Error()})
		return
	}

	res, err := signIn.SignIn(a.database)
	if err != nil{
		ctx.JSON(500, gin.H{"error": err.Error()})
		return 
	}

	ctx.SetCookie("access_token", res, 3600*24, "/", "", false, true)
	ctx.JSON(200, gin.H{"status": "succesfully login"})
	return 
	
}