package handlers

import (
	"database/sql"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"

	"Yamy-Gin/pkg/auth"
	"Yamy-Gin/pkg/entry"
	"Yamy-Gin/pkg/llm"
)

type Handler struct {
	database *sql.DB
}

func NewHandler(database *sql.DB) *Handler {
	return &Handler{database: database}
}

func (a *Handler) SignUp(ctx *gin.Context) {
	var signUp auth.InputSignUp

	if err := ctx.ShouldBindJSON(&signUp); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}

	errSign := signUp.SignUp(a.database)

	if errSign != nil {
		ctx.JSON(400, gin.H{"error": errSign.Error()})
		return
	}

	ctx.JSON(200, gin.H{"status": "succesfully created account"})
}

func (a *Handler) SignIn(ctx *gin.Context) {
	var signIn auth.InputSignIn

	if err := ctx.ShouldBindJSON(&signIn); err != nil {
		ctx.JSON(401, gin.H{"error": err.Error()})
		return
	}

	res, err := signIn.SignIn(a.database)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.SetCookie("access_token", res, 3600*24, "/", "", false, true)
	ctx.JSON(200, gin.H{"status": "succesfully login",
		"access_token": res})
	return

}


func (a *Handler) EntryCreate(ctx *gin.Context) {
	userId, ok := ctx.Get("user_id")
	if !ok {
		ctx.JSON(401, gin.H{"error": "unauthorized"})
		return
	}

	var entryData entry.Entry
	if err := ctx.ShouldBindBodyWithJSON(&entryData); err != nil {
		ctx.JSON(400, gin.H{"error": "need post with entry data"})
		return
	}

	newRepo := entry.NewRepository(a.database)
	USERID, _ := userId.(int)

	newEntryID, err := newRepo.Create(USERID, entryData.Content, entryData.MoodScore, entryData.Sentiment)
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	go func(content string, entryID int64) {
		result, err := llm.AnalyzeEntry(content)
		if err != nil {
			log.Printf("[ERROR] LLM analyze failed: %v", err)
			return
		}
		newRepo.UpdateSentiment(entryID, result.Sentiment)
	}(entryData.Content, newEntryID)

	ctx.JSON(200, gin.H{"status": "created"})
}

func (a *Handler) GetById(ctx *gin.Context){
	userid, ok := ctx.Get("user_id")
	useridINT := userid.(int)

	id := ctx.Param("id")
	idINT, err := strconv.Atoi(id)
	if err != nil{
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if !ok {
		ctx.JSON(500, gin.H{"error": "Can't find user_id in json"})
		return
	}
	newRepo := entry.NewRepository(a.database)


	entryGET, errEntry := newRepo.GetByID(idINT, useridINT)
	if errEntry != nil{
		ctx.JSON(500, gin.H{"error": errEntry.Error()})
		return
	}

	ctx.JSON(200, gin.H{"Entry": *entryGET})
}

func (a *Handler) DeleteByID(ctx *gin.Context){
	userid, ok := ctx.Get("user_id")
	useridINT := userid.(int)

	id := ctx.Param("id")
	idINT, err := strconv.Atoi(id)
	if err != nil{
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if !ok {
		ctx.JSON(500, gin.H{"error": "Can't find user_id in json"})
		return
	}
	newRepo := entry.NewRepository(a.database)

	errDelete, _ := newRepo.Delete(idINT, useridINT)
	if errDelete != nil{
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{"status": "deleted"})
}

func (a *Handler) GetAllById(ctx *gin.Context){
	userid, ok := ctx.Get("user_id")
	useridINT := userid.(int)
	if !ok {
		ctx.JSON(500, gin.H{"error": "Can't find user_id in json"})
		return
	}


	newRepo := entry.NewRepository(a.database)
	ListOfEntrys, err := newRepo.GetAllByUser(useridINT)
	if err != nil{
		ctx.JSON(500, gin.H{"error": "Can't find user_id in json"})
		return
	}

	ctx.JSON(200, gin.H{"ok": ListOfEntrys})
}