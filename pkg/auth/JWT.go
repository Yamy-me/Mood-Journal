package auth

import (
	"os"

	"github.com/joho/godotenv"
)

var SECRET_KEY string

func init() {
	godotenv.Load(`c:\Users\Atai\Desktop\newGoLang\.env`)
	SECRET_KEY = os.Getenv("TOKEN")
}

